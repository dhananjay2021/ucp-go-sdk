// Copyright 2026 UCP Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package main demonstrates implementing a UCP-compliant business server.
//
// This example shows how to:
// - Set up a UCP server with capabilities
// - Handle checkout creation, updates, and completion
// - Manage orders
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"sync/atomic"

	"github.com/dhananjay2021/ucp-go-sdk/client"
	"github.com/dhananjay2021/ucp-go-sdk/extensions"
	"github.com/dhananjay2021/ucp-go-sdk/models"
	"github.com/dhananjay2021/ucp-go-sdk/server"
)

// specVersion is the UCP release this demo server advertises and stamps on
// responses. Sourced from the SDK so the example tracks the models it ships with.
const specVersion = models.SpecVersion

// demoSigningKey is a sample EC P-256 public JWK published in the discovery
// profile's canonical keys[] field (UCP 2026-08-25). In a real deployment this
// is your public signing key; the private half stays on your signer.
var demoSigningKey = models.JWK{
	Kid: "demo-key-1",
	Kty: "EC",
	Crv: "P-256",
	X:   "f83OJ3D2xF1Bg8vub9tLe1gHMzV76e8Tus9uPHvRVEU",
	Y:   "x_FEzRu9m36HLN_tue659LNpXW6pCyStikYjKIWI5a0",
	Use: "sig",
	Alg: "ES256",
}

// In-memory product catalog for demo
var productCatalog = map[string]struct {
	Title    string
	Price    int // cents
	ImageURL string
}{
	"PROD-001": {Title: "Wireless Headphones", Price: 14999, ImageURL: "https://example.com/images/headphones.jpg"},
	"PROD-002": {Title: "Phone Case", Price: 2999, ImageURL: "https://example.com/images/case.jpg"},
}

// In-memory storage for demo purposes
var (
	checkouts = make(map[string]*extensions.ExtendedCheckoutResponse)
	orders    = make(map[string]*models.Order)
	carts     = make(map[string]*models.CartResponse)
	mu        sync.RWMutex
	idCounter atomic.Int64
)

func generateID(prefix string) string {
	id := idCounter.Add(1)
	return fmt.Sprintf("%s-%d", prefix, id)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Configure the UCP server
	config := server.Config{
		Version: specVersion,
		Capabilities: []models.CapabilityDiscovery{
			{
				CapabilityBase: models.CapabilityBase{
					Name:    client.CapabilityCheckout,
					Version: specVersion,
					Spec:    "https://ucp.dev/specification/checkout",
					Schema:  "https://ucp.dev/schemas/shopping/checkout.json",
				},
			},
			{
				CapabilityBase: models.CapabilityBase{
					Name:    client.CapabilityOrder,
					Version: specVersion,
					Spec:    "https://ucp.dev/specification/order",
					Schema:  "https://ucp.dev/schemas/shopping/order.json",
				},
			},
			{
				CapabilityBase: models.CapabilityBase{
					Name:    client.CapabilityFulfillment,
					Version: specVersion,
					Spec:    "https://ucp.dev/specification/fulfillment",
					Schema:  "https://ucp.dev/schemas/shopping/fulfillment.json",
					Extends: client.CapabilityCheckout,
				},
			},
			{
				CapabilityBase: models.CapabilityBase{
					Name:    "dev.ucp.shopping.cart",
					Version: specVersion,
					Spec:    "https://ucp.dev/specification/cart",
					Schema:  "https://ucp.dev/schemas/shopping/cart.json",
				},
			},
			// Catalog browsing.
			{
				CapabilityBase: models.CapabilityBase{
					Name:    client.CapabilityCatalogSearch,
					Version: specVersion,
					Schema:  "https://ucp.dev/schemas/shopping/catalog.json",
				},
			},
			{
				CapabilityBase: models.CapabilityBase{
					Name:    client.CapabilityCatalogLookup,
					Version: specVersion,
					Schema:  "https://ucp.dev/schemas/shopping/catalog.json",
				},
			},
			// Per-purpose buyer consent (extends checkout).
			{
				CapabilityBase: models.CapabilityBase{
					Name:    client.CapabilityBuyerConsent,
					Version: specVersion,
					Schema:  "https://ucp.dev/schemas/shopping/buyer_consent.json",
					Extends: client.CapabilityCheckout,
				},
			},
			// Browser-addressable shopping permalinks: config carries the
			// business browser endpoint.
			{
				CapabilityBase: models.CapabilityBase{
					Name:    client.CapabilityPermalink,
					Version: specVersion,
					Schema:  "https://ucp.dev/schemas/shopping/permalink.json",
					Config: map[string]interface{}{
						"endpoint": fmt.Sprintf("http://localhost:%s/permalink", port),
					},
				},
			},
			// Location search + lookup (common namespace).
			{
				CapabilityBase: models.CapabilityBase{
					Name:    client.CapabilityLocationSearch,
					Version: specVersion,
					Schema:  "https://ucp.dev/schemas/common/location_search.json",
				},
			},
			{
				CapabilityBase: models.CapabilityBase{
					Name:    client.CapabilityLocationLookup,
					Version: specVersion,
					Schema:  "https://ucp.dev/schemas/common/location_lookup.json",
				},
			},
			// Loyalty membership + earning forecasts (extends checkout).
			{
				CapabilityBase: models.CapabilityBase{
					Name:    client.CapabilityLoyalty,
					Version: specVersion,
					Schema:  "https://ucp.dev/schemas/common/loyalty.json",
					Extends: client.CapabilityCheckout,
				},
			},
			// Payment terms (alternative payment schedules; extends checkout).
			{
				CapabilityBase: models.CapabilityBase{
					Name:    client.CapabilityPaymentTerms,
					Version: specVersion,
					Schema:  "https://ucp.dev/schemas/common/payment_terms.json",
					Extends: client.CapabilityCheckout,
				},
			},
			// Split payments: config declares valid instrument combinations.
			{
				CapabilityBase: models.CapabilityBase{
					Name:    client.CapabilitySplitPayments,
					Version: specVersion,
					Schema:  "https://ucp.dev/schemas/common/payment_split_payments.json",
					Extends: client.CapabilityCheckout,
					Config: map[string]interface{}{
						"allowed_combinations": [][]models.InstrumentGroup{{
							{Types: []string{"card"}, Min: 1, Max: 1},
							{Types: []string{"gift_card"}, Min: 0, Max: 3},
						}},
					},
				},
			},
			// Payment authentication (3DS device data collection + challenge).
			{
				CapabilityBase: models.CapabilityBase{
					Name:    client.CapabilityPaymentAuth,
					Version: specVersion,
					Schema:  "https://ucp.dev/schemas/common/payment_authentication.json",
					Extends: client.CapabilityCheckout,
				},
			},
		},
		SigningKeys: []models.JWK{demoSigningKey},
		Services: models.Services{
			client.ServiceShopping: models.UCPService{
				Version: specVersion,
				Spec:    "https://ucp.dev/specification/shopping",
				Rest: &models.RestTransport{
					Schema:   "https://ucp.dev/schemas/services/shopping/rest.openapi.json",
					Endpoint: fmt.Sprintf("http://localhost:%s", port),
				},
			},
		},
		PaymentHandlers: []models.PaymentHandlerResponse{
			{
				ID:                "default",
				Name:              "dev.ucp.tokenization",
				Version:           specVersion,
				Spec:              "https://ucp.dev/handlers/tokenization/spec",
				ConfigSchema:      "https://ucp.dev/handlers/tokenization/config.json",
				InstrumentSchemas: []string{"https://ucp.dev/schemas/shopping/types/card_payment_instrument.json"},
				Config:            map[string]interface{}{"gateway": "demo"},
			},
		},
	}

	// Create the server
	srv := server.NewServer(config)

	// Register handlers
	srv.HandleCreateCheckout(handleCreateCheckout)
	srv.HandleGetCheckout(handleGetCheckout)
	srv.HandleUpdateCheckout(handleUpdateCheckout)
	srv.HandleCompleteCheckout(handleCompleteCheckout)
	srv.HandleCancelCheckout(handleCancelCheckout)
	srv.HandleGetOrder(handleGetOrder)

	// Register cart handlers
	srv.HandleCreateCart(handleCreateCart)
	srv.HandleGetCart(handleGetCart)
	srv.HandleUpdateCart(handleUpdateCart)
	srv.HandleDeleteCart(handleDeleteCart)

	// Register catalog handlers
	srv.HandleSearchCatalog(handleSearchCatalog)
	srv.HandleLookupCatalog(handleLookupCatalog)

	// Compose a parent mux: the UCP server handles all its built-in routes as
	// the catch-all, and we mount the newer common/ capabilities (location
	// search/lookup) and the permalink browser endpoint alongside it. Go 1.22's
	// pattern precedence routes the specific patterns here and everything else to
	// the UCP server.
	root := http.NewServeMux()
	root.HandleFunc("POST /locations/search", handleLocationSearch)
	root.HandleFunc("POST /locations/lookup", handleLocationLookup)
	root.HandleFunc("GET /permalink", handlePermalink)
	root.Handle("/", srv)

	// Apply middleware
	handler := server.Chain(root,
		server.LoggingMiddleware,
		server.RequestIDMiddleware,
		server.CORSMiddleware([]string{"*"}),
	)

	log.Printf("Starting UCP business server on port %s", port)
	log.Printf("Discovery endpoint: http://localhost:%s/.well-known/ucp", port)
	log.Printf("Advertising %d capabilities incl. permalink, location, loyalty, payment terms, split payments", len(config.Capabilities))

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func handleCreateCheckout(r *http.Request, req *extensions.ExtendedCheckoutCreateRequest) (*extensions.ExtendedCheckoutResponse, error) {
	checkoutID := generateID("chk")

	// Calculate totals - look up items from catalog
	var subtotal int
	lineItems := make([]models.LineItemResponse, len(req.LineItems))

	for i, li := range req.LineItems {
		product, ok := productCatalog[li.Item.ID]
		if !ok {
			return nil, server.BadRequestError(fmt.Sprintf("unknown product: %s", li.Item.ID))
		}

		itemTotal := product.Price * li.Quantity
		subtotal += itemTotal

		lineItems[i] = models.LineItemResponse{
			ID: generateID("li"),
			Item: models.ItemResponse{
				ID:       li.Item.ID,
				Title:    product.Title,
				Price:    product.Price,
				ImageURL: product.ImageURL,
			},
			Quantity: li.Quantity,
			Totals: []models.TotalResponse{
				{Type: models.TotalTypeSubtotal, Amount: itemTotal},
			},
		}
	}

	tax := subtotal * 875 / 10000 // 8.75% tax

	// Create checkout response
	checkout := &extensions.ExtendedCheckoutResponse{
		UCP: models.ResponseCheckout{
			Version: specVersion,
			Capabilities: []models.CapabilityResponse{
				{CapabilityBase: models.CapabilityBase{Name: client.CapabilityCheckout, Version: specVersion}},
			},
		},
		ID:        checkoutID,
		LineItems: lineItems,
		Status:    models.CheckoutStatusIncomplete,
		Currency:  req.Currency,
		Totals: []models.TotalResponse{
			{Type: models.TotalTypeSubtotal, Amount: subtotal},
			{Type: models.TotalTypeTax, Amount: tax},
			{Type: models.TotalTypeTotal, Amount: subtotal + tax},
		},
		Links: []models.Link{
			{Type: "terms_of_service", URL: "https://example.com/terms", Title: "Terms of Service"},
			{Type: "privacy_policy", URL: "https://example.com/privacy", Title: "Privacy Policy"},
		},
		Payment: models.PaymentResponse{
			Handlers: []models.PaymentHandlerResponse{
				{
					ID:                "default",
					Name:              "dev.ucp.tokenization",
					Version:           specVersion,
					Spec:              "https://ucp.dev/handlers/tokenization/spec",
					ConfigSchema:      "https://ucp.dev/handlers/tokenization/config.json",
					InstrumentSchemas: []string{"https://ucp.dev/schemas/shopping/types/card_payment_instrument.json"},
					Config:            map[string]interface{}{"gateway": "demo"},
				},
			},
		},
		Messages: []models.Message{
			{Type: models.MessageTypeInfo, Content: "Buyer information required", Severity: models.SeverityRecoverable},
		},
		// Policies extension: durable return/warranty terms for these items.
		Policies: []models.Policy{
			{
				Type:        "dev.ucp.shopping.policy.return",
				Description: models.Description{Plain: "Free returns within 30 days of delivery."},
				URL:         "https://example.com/returns",
			},
		},
		// Payment Terms extension: offer "Pay now" vs "Pay in 4".
		PaymentTerms: buildPaymentTerms(subtotal + tax),
		// Loyalty extension: forecast rewards for this transaction.
		Loyalty: models.Loyalty{
			"com.example.rewards": {
				ID:   "mem-9021",
				Name: "Example Rewards",
				Rewards: []models.LoyaltyMembershipReward{
					{
						Currency: models.LoyaltyRewardCurrency{Name: "Example Points", Code: "PTS"},
						EarningForecast: &models.LoyaltyEarningForecast{
							Amount: models.LoyaltyRewardAmount((subtotal + tax) / 100),
							Breakdown: []models.LoyaltyEarningBreakdown{
								{ID: "base", Amount: models.LoyaltyRewardAmount((subtotal + tax) / 100), Description: "1 point per $1"},
							},
						},
					},
				},
				Provisional: false,
			},
		},
	}

	// Store checkout
	mu.Lock()
	checkouts[checkoutID] = checkout
	mu.Unlock()

	log.Printf("Created checkout %s with %d items, subtotal: %d cents", checkoutID, len(lineItems), subtotal)
	return checkout, nil
}

// buildPaymentTerms offers two ways to pay for the checkout total: the full
// amount now, or four equal installments. Demonstrates the Payment Terms
// extension (UCP spec change #602).
func buildPaymentTerms(total int) []models.PaymentTerm {
	quarter := total / 4
	remainder := total - quarter*3
	return []models.PaymentTerm{
		{
			ID:    "pay_now",
			Title: "Pay now",
			Schedules: []models.PaymentSchedule{
				{
					ID:          "now",
					Amount:      models.Amount(total),
					Description: models.Description{Plain: "Full amount due at checkout"},
					Type:        string(models.PaymentScheduleTypeImmediate),
				},
			},
		},
		{
			ID:          "pay_in_4",
			Title:       "Pay in 4",
			Description: &models.Description{Plain: "4 interest-free payments"},
			Schedules: []models.PaymentSchedule{
				{ID: "p1", Amount: models.Amount(remainder), Description: models.Description{Plain: "Due today"}, Type: string(models.PaymentScheduleTypeImmediate)},
				{ID: "p2", Amount: models.Amount(quarter), Description: models.Description{Plain: "Due in 2 weeks"}, Type: "deferred"},
				{ID: "p3", Amount: models.Amount(quarter), Description: models.Description{Plain: "Due in 4 weeks"}, Type: "deferred"},
				{ID: "p4", Amount: models.Amount(quarter), Description: models.Description{Plain: "Due in 6 weeks"}, Type: "deferred"},
			},
		},
	}
}

func handleGetCheckout(r *http.Request, id string) (*extensions.ExtendedCheckoutResponse, error) {
	mu.RLock()
	checkout, ok := checkouts[id]
	mu.RUnlock()

	if !ok {
		return nil, server.NotFoundError("checkout not found")
	}

	return checkout, nil
}

func handleUpdateCheckout(r *http.Request, id string, req *extensions.ExtendedCheckoutUpdateRequest) (*extensions.ExtendedCheckoutResponse, error) {
	mu.Lock()
	defer mu.Unlock()

	checkout, ok := checkouts[id]
	if !ok {
		return nil, server.NotFoundError("checkout not found")
	}

	// Update buyer info
	if req.Buyer != nil {
		checkout.Buyer = &models.BuyerWithConsentResponse{
			Email:       req.Buyer.Email,
			PhoneNumber: req.Buyer.PhoneNumber,
			FirstName:   req.Buyer.FirstName,
			LastName:    req.Buyer.LastName,
			FullName:    req.Buyer.FullName,
			Consent:     req.Buyer.Consent,
		}
	}

	// Update fulfillment
	if req.Fulfillment != nil && len(req.Fulfillment.Methods) > 0 {
		methods := make([]models.FulfillmentMethodResponse, len(req.Fulfillment.Methods))
		for i, m := range req.Fulfillment.Methods {
			destinations := make([]models.FulfillmentDestinationResponse, len(m.Destinations))
			for j, d := range m.Destinations {
				destID := generateID("dest")
				destinations[j] = models.FulfillmentDestinationResponse{
					PostalAddress: d.PostalAddress,
					ID:            destID,
				}
			}
			methods[i] = models.FulfillmentMethodResponse{
				ID:           m.ID,
				Type:         models.FulfillmentMethodTypeShipping, // Default to shipping for demo
				LineItemIDs:  m.LineItemIDs,
				Destinations: destinations,
			}
		}
		checkout.Fulfillment = &models.FulfillmentResponse{
			Methods: methods,
		}
	}

	// Update payment
	if req.Payment.SelectedInstrumentID != "" {
		checkout.Payment.SelectedInstrumentID = req.Payment.SelectedInstrumentID
		checkout.Payment.Instruments = req.Payment.Instruments

		// Payment Authentication extension: once an instrument is selected, ask
		// the platform to run a 3DS challenge before completion. Surfaced as an
		// outstanding Action keyed by its reverse-domain Action type.
		challenge := models.PaymentThreeDSChallengeConfig{
			PaymentInstrumentID: req.Payment.SelectedInstrumentID,
			URL:                 fmt.Sprintf("https://example.com/3ds/%s", id),
		}
		checkout.Actions = models.Actions{
			models.ActionTypeThreeDSChallenge: {
				{
					ID: generateID("act"),
					Config: map[string]interface{}{
						"payment_instrument_id": challenge.PaymentInstrumentID,
						"url":                   challenge.URL,
					},
				},
			},
		}
	}

	// Update status based on completeness
	checkout.Messages = nil
	if checkout.Buyer == nil || checkout.Buyer.Email == "" {
		checkout.Messages = append(checkout.Messages, models.Message{
			Type:     models.MessageTypeInfo,
			Content:  "Email required",
			Severity: models.SeverityRecoverable,
			Path:     "$.buyer.email",
		})
	}
	if checkout.Payment.SelectedInstrumentID == "" {
		checkout.Messages = append(checkout.Messages, models.Message{
			Type:     models.MessageTypeInfo,
			Content:  "Payment required",
			Severity: models.SeverityRecoverable,
			Path:     "$.payment",
		})
	}

	if len(checkout.Messages) == 0 {
		checkout.Status = models.CheckoutStatusReadyForComplete
	} else {
		checkout.Status = models.CheckoutStatusIncomplete
	}

	log.Printf("Updated checkout %s, status: %s", id, checkout.Status)
	return checkout, nil
}

func handleCompleteCheckout(r *http.Request, id string) (*extensions.ExtendedCheckoutResponse, error) {
	// Generate order ID before acquiring lock to avoid deadlock
	orderID := generateID("ord")

	mu.Lock()
	defer mu.Unlock()

	checkout, ok := checkouts[id]
	if !ok {
		return nil, server.NotFoundError("checkout not found")
	}

	if checkout.Status != models.CheckoutStatusReadyForComplete {
		return nil, server.BadRequestError("checkout is not ready for completion")
	}

	orderLineItems := make([]models.OrderLineItem, len(checkout.LineItems))
	for i, li := range checkout.LineItems {
		orderLineItems[i] = models.OrderLineItem{
			ID:   li.ID,
			Item: li.Item,
			Quantity: models.OrderLineItemQuantity{
				Total:     li.Quantity,
				Fulfilled: 0,
			},
			Totals: li.Totals,
			Status: models.OrderLineItemStatusProcessing,
		}
	}

	order := &models.Order{
		UCP: models.ResponseOrder{
			Version: specVersion,
			Capabilities: []models.CapabilityResponse{
				{CapabilityBase: models.CapabilityBase{Name: client.CapabilityOrder, Version: specVersion}},
			},
		},
		ID:           orderID,
		CheckoutID:   id,
		Label:        fmt.Sprintf("Order %s", orderID),
		PermalinkURL: fmt.Sprintf("https://example.com/orders/%s", orderID),
		LineItems:    orderLineItems,
		Totals:       checkout.Totals,
		Fulfillment:  models.OrderFulfillment{},
		// Snapshot the return policy that applied at checkout onto the order.
		Policies: []models.Policy{
			{
				Type:        "dev.ucp.shopping.policy.return",
				Description: models.Description{Plain: "Free returns within 30 days of delivery."},
				URL:         "https://example.com/returns",
			},
		},
	}

	orders[orderID] = order

	// Update checkout
	checkout.Status = models.CheckoutStatusCompleted
	checkout.Order = &models.OrderConfirmation{
		ID:           orderID,
		PermalinkURL: order.PermalinkURL,
	}

	log.Printf("Completed checkout %s, created order %s", id, orderID)
	return checkout, nil
}

func handleCancelCheckout(r *http.Request, id string) (*extensions.ExtendedCheckoutResponse, error) {
	mu.Lock()
	defer mu.Unlock()

	checkout, ok := checkouts[id]
	if !ok {
		return nil, server.NotFoundError("checkout not found")
	}

	if checkout.Status == models.CheckoutStatusCompleted {
		return nil, server.BadRequestError("cannot cancel completed checkout")
	}

	checkout.Status = models.CheckoutStatusCanceled

	log.Printf("Canceled checkout %s", id)
	return checkout, nil
}

func handleGetOrder(r *http.Request, id string) (*models.Order, error) {
	mu.RLock()
	order, ok := orders[id]
	mu.RUnlock()

	if !ok {
		return nil, server.NotFoundError("order not found")
	}

	return order, nil
}

// Cart handlers

func handleCreateCart(r *http.Request, req *models.CartCreateRequest) (*models.CartResponse, error) {
	mu.Lock()
	defer mu.Unlock()

	cartID := generateID("cart")

	// Build line items with pricing from catalog
	lineItems := make([]models.LineItemResponse, 0, len(req.LineItems))
	subtotal := 0

	for _, li := range req.LineItems {
		product, ok := productCatalog[li.Item.ID]
		if !ok {
			return nil, server.BadRequestError("unknown product: " + li.Item.ID)
		}

		lineTotal := product.Price * li.Quantity
		subtotal += lineTotal

		lineItems = append(lineItems, models.LineItemResponse{
			ID:       li.Item.ID,
			Quantity: li.Quantity,
			Item: models.ItemResponse{
				ID:       li.Item.ID,
				Title:    product.Title,
				ImageURL: product.ImageURL,
			},
			Totals: []models.TotalResponse{
				{Type: models.TotalTypeSubtotal, Amount: lineTotal},
			},
		})
	}

	// Calculate estimated totals (no tax yet without address)
	cart := &models.CartResponse{
		UCP: &models.ResponseCart{
			Schema: "https://ucp.dev/schemas/shopping/cart.json",
		},
		ID:        cartID,
		LineItems: lineItems,
		Currency:  "USD", // Determined by context or geo-IP
		Totals: []models.TotalResponse{
			{Type: models.TotalTypeSubtotal, Amount: subtotal},
			{Type: models.TotalTypeTotal, Amount: subtotal}, // Estimated, no tax yet
		},
		Messages: []models.Message{
			{
				Type:    models.MessageTypeInfo,
				Content: "Tax will be calculated at checkout with shipping address.",
			},
		},
	}

	// Store context if provided
	if req.Context != nil {
		log.Printf("Cart created with context: country=%s, region=%s, intent=%s",
			req.Context.AddressCountry, req.Context.AddressRegion, req.Context.Intent)
	}

	carts[cartID] = cart
	log.Printf("Created cart %s with %d items, subtotal: %d cents", cartID, len(lineItems), subtotal)

	return cart, nil
}

func handleGetCart(r *http.Request, id string) (*models.CartResponse, error) {
	mu.RLock()
	cart, ok := carts[id]
	mu.RUnlock()

	if !ok {
		return nil, server.NotFoundError("cart not found")
	}

	return cart, nil
}

func handleUpdateCart(r *http.Request, id string, req *models.CartUpdateRequest) (*models.CartResponse, error) {
	mu.Lock()
	defer mu.Unlock()

	cart, ok := carts[id]
	if !ok {
		return nil, server.NotFoundError("cart not found")
	}

	// Rebuild line items with new quantities
	lineItems := make([]models.LineItemResponse, 0, len(req.LineItems))
	subtotal := 0

	for _, li := range req.LineItems {
		product, ok := productCatalog[li.Item.ID]
		if !ok {
			return nil, server.BadRequestError("unknown product: " + li.Item.ID)
		}

		lineTotal := product.Price * li.Quantity
		subtotal += lineTotal

		lineItems = append(lineItems, models.LineItemResponse{
			ID:       li.Item.ID,
			Quantity: li.Quantity,
			Item: models.ItemResponse{
				ID:       li.Item.ID,
				Title:    product.Title,
				ImageURL: product.ImageURL,
			},
			Totals: []models.TotalResponse{
				{Type: models.TotalTypeSubtotal, Amount: lineTotal},
			},
		})
	}

	// Update cart
	cart.LineItems = lineItems
	cart.Totals = []models.TotalResponse{
		{Type: models.TotalTypeSubtotal, Amount: subtotal},
		{Type: models.TotalTypeTotal, Amount: subtotal},
	}

	log.Printf("Updated cart %s, new subtotal: %d cents", id, subtotal)
	return cart, nil
}

func handleDeleteCart(r *http.Request, id string) error {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := carts[id]; !ok {
		return server.NotFoundError("cart not found")
	}

	delete(carts, id)
	log.Printf("Deleted cart %s", id)
	return nil
}

// ---------------------------------------------------------------------------
// Catalog handlers (dev.ucp.shopping.catalog.search / .lookup)
// ---------------------------------------------------------------------------

func catalogProduct(id string) models.Product {
	p := productCatalog[id]
	price := models.Price{Amount: p.Price, Currency: "USD"}
	return models.Product{
		ID:         id,
		Title:      p.Title,
		PriceRange: models.PriceRange{Min: price, Max: price},
		Media:      []models.Media{{Type: "image", URL: p.ImageURL}},
		Variants: []models.Variant{
			{
				ID:           id,
				Title:        p.Title,
				Price:        price,
				Availability: &models.Availability{Available: true, Status: "in_stock"},
			},
		},
	}
}

func handleSearchCatalog(r *http.Request, req *models.CatalogSearchRequest) (*models.CatalogSearchResponse, error) {
	products := make([]models.Product, 0, len(productCatalog))
	for id := range productCatalog {
		products = append(products, catalogProduct(id))
	}
	log.Printf("Catalog search %q returned %d products", req.Query, len(products))
	return &models.CatalogSearchResponse{Products: products}, nil
}

func handleLookupCatalog(r *http.Request, req *models.CatalogLookupRequest) (*models.CatalogLookupResponse, error) {
	products := make([]models.Product, 0, len(req.IDs))
	var messages []models.Message
	for _, id := range req.IDs {
		if _, ok := productCatalog[id]; ok {
			products = append(products, catalogProduct(id))
		} else {
			messages = append(messages, models.Message{Type: models.MessageTypeInfo, Content: "not found: " + id})
		}
	}
	return &models.CatalogLookupResponse{Products: products, Messages: messages}, nil
}

// ---------------------------------------------------------------------------
// Location search + lookup (dev.ucp.common.location.search / .lookup)
// ---------------------------------------------------------------------------

func strPtr(s string) *string { return &s }

// demoLocations are two sample storefronts used by the location capability demo.
var demoLocations = []models.Location{
	{
		ID:   "loc-sf",
		Name: "Example Store — San Francisco",
		Address: &models.PostalAddress{
			StreetAddress: "1 Market St", AddressLocality: "San Francisco",
			AddressRegion: "CA", PostalCode: "94105", AddressCountry: "US",
		},
		Geo:      &models.Geo{Latitude: 37.7936, Longitude: -122.3965},
		Timezone: strPtr("America/Los_Angeles"),
		Amenities: models.LocationAmenities{
			"dev.ucp.amenity.curbside_pickup": {Description: "Curbside pickup"},
			"dev.ucp.amenity.wi_fi":           {Description: "Free Wi-Fi"},
		},
		Hours: []models.DailyHour{
			{Day: models.DailyHourDayMonday, Opens: "09:00", Closes: "18:00"},
			{Day: models.DailyHourDaySaturday, Opens: "10:00", Closes: "16:00"},
		},
	},
	{
		ID:   "loc-mv",
		Name: "Example Store — Mountain View",
		Address: &models.PostalAddress{
			StreetAddress: "600 Amphitheatre Pkwy", AddressLocality: "Mountain View",
			AddressRegion: "CA", PostalCode: "94043", AddressCountry: "US",
		},
		Geo:      &models.Geo{Latitude: 37.4220, Longitude: -122.0841},
		Timezone: strPtr("America/Los_Angeles"),
		Amenities: models.LocationAmenities{
			"dev.ucp.amenity.curbside_pickup": {Description: "Curbside pickup"},
		},
	},
}

func handleLocationSearch(w http.ResponseWriter, r *http.Request) {
	var req models.LocationSearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err.Error() != "EOF" {
		server.WriteError(w, http.StatusBadRequest, "invalid_request", "Failed to parse request body")
		return
	}
	q := ""
	if req.Query != nil {
		q = *req.Query
	}
	log.Printf("Location search %q returned %d locations", q, len(demoLocations))
	server.WriteJSON(w, http.StatusOK, &models.LocationSearchResponse{
		UCP: &models.ResponseLocation{
			Version:      specVersion,
			Capabilities: []models.CapabilityResponse{{CapabilityBase: models.CapabilityBase{Name: client.CapabilityLocationSearch, Version: specVersion}}},
		},
		Locations:  demoLocations,
		Pagination: &models.LocationPaginationResponse{HasNextPage: false},
	})
}

func handleLocationLookup(w http.ResponseWriter, r *http.Request) {
	var req models.LocationLookupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.WriteError(w, http.StatusBadRequest, "invalid_request", "Failed to parse request body")
		return
	}

	byID := make(map[string]models.Location, len(demoLocations))
	for _, loc := range demoLocations {
		byID[loc.ID] = loc
	}

	results := make([]models.LocationLookupResult, 0, len(req.IDs))
	var messages []models.Message
	for _, id := range req.IDs {
		if loc, ok := byID[id]; ok {
			results = append(results, models.LocationLookupResult{
				Location: loc,
				Inputs:   []models.LocationLookupInput{{ID: id}},
			})
		} else {
			messages = append(messages, models.Message{Type: models.MessageTypeInfo, Content: "not found: " + id})
		}
	}

	server.WriteJSON(w, http.StatusOK, &models.LocationLookupResponse{
		UCP: &models.ResponseLocation{
			Version:      specVersion,
			Capabilities: []models.CapabilityResponse{{CapabilityBase: models.CapabilityBase{Name: client.CapabilityLocationLookup, Version: specVersion}}},
		},
		Locations: results,
		Messages:  messages,
	})
}

// ---------------------------------------------------------------------------
// Permalink browser endpoint (dev.ucp.shopping.permalink)
// ---------------------------------------------------------------------------

// handlePermalink resolves a shopping permalink into a buyer-facing destination
// and issues a 303 redirect, per the permalink browser binding. The `items`
// query is a compact comma-separated list of item_id:quantity pairs; an optional
// `continue_to` root-relative path steers the post-resolution destination.
func handlePermalink(w http.ResponseWriter, r *http.Request) {
	items := r.URL.Query().Get("items")
	continueTo := r.URL.Query().Get("continue_to")

	dest := "/cart"
	if continueTo != "" {
		dest = continueTo
	}
	if items != "" {
		dest = fmt.Sprintf("%s?items=%s", dest, items)
	}

	log.Printf("Permalink resolved items=%q -> %s", items, dest)
	http.Redirect(w, r, dest, http.StatusSeeOther)
}

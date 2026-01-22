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
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/dhananjay2021/ucp-go-sdk/client"
	"github.com/dhananjay2021/ucp-go-sdk/extensions"
	"github.com/dhananjay2021/ucp-go-sdk/models"
	"github.com/dhananjay2021/ucp-go-sdk/server"
)

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
	mu        sync.RWMutex
	idCounter = 0
)

func generateID(prefix string) string {
	mu.Lock()
	defer mu.Unlock()
	idCounter++
	return fmt.Sprintf("%s-%d", prefix, idCounter)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Configure the UCP server
	config := server.Config{
		Version: "2026-01-11",
		Capabilities: []models.CapabilityDiscovery{
			{
				CapabilityBase: models.CapabilityBase{
					Name:    client.CapabilityCheckout,
					Version: "2026-01-11",
					Spec:    "https://ucp.dev/specification/checkout",
					Schema:  "https://ucp.dev/schemas/shopping/checkout.json",
				},
			},
			{
				CapabilityBase: models.CapabilityBase{
					Name:    client.CapabilityOrder,
					Version: "2026-01-11",
					Spec:    "https://ucp.dev/specification/order",
					Schema:  "https://ucp.dev/schemas/shopping/order.json",
				},
			},
			{
				CapabilityBase: models.CapabilityBase{
					Name:    client.CapabilityFulfillment,
					Version: "2026-01-11",
					Spec:    "https://ucp.dev/specification/fulfillment",
					Schema:  "https://ucp.dev/schemas/shopping/fulfillment.json",
					Extends: client.CapabilityCheckout,
				},
			},
		},
		Services: models.Services{
			client.ServiceShopping: models.UCPService{
				Version: "2026-01-11",
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
				Version:           "2026-01-11",
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

	// Apply middleware
	handler := server.Chain(srv,
		server.LoggingMiddleware,
		server.RequestIDMiddleware,
		server.CORSMiddleware([]string{"*"}),
	)

	log.Printf("Starting UCP business server on port %s", port)
	log.Printf("Discovery endpoint: http://localhost:%s/.well-known/ucp", port)

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
			Version: "2026-01-11",
			Capabilities: []models.CapabilityResponse{
				{CapabilityBase: models.CapabilityBase{Name: client.CapabilityCheckout, Version: "2026-01-11"}},
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
					Version:           "2026-01-11",
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
	}

	// Store checkout
	mu.Lock()
	checkouts[checkoutID] = checkout
	mu.Unlock()

	log.Printf("Created checkout %s with %d items, subtotal: %d cents", checkoutID, len(lineItems), subtotal)
	return checkout, nil
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
			Version: "2026-01-11",
			Capabilities: []models.CapabilityResponse{
				{CapabilityBase: models.CapabilityBase{Name: client.CapabilityOrder, Version: "2026-01-11"}},
			},
		},
		ID:           orderID,
		CheckoutID:   id,
		PermalinkURL: fmt.Sprintf("https://example.com/orders/%s", orderID),
		LineItems:    orderLineItems,
		Totals:       checkout.Totals,
		Fulfillment:  models.OrderFulfillment{},
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

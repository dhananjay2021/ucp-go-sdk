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
	"time"

	"github.com/Universal-Commerce-Protocol/go-sdk/client"
	"github.com/Universal-Commerce-Protocol/go-sdk/extensions"
	"github.com/Universal-Commerce-Protocol/go-sdk/models"
	"github.com/Universal-Commerce-Protocol/go-sdk/server"
)

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
				ID:              "default",
				Type:            "tokenization",
				Name:            "Demo Payment Handler",
				TokenizationURL: fmt.Sprintf("http://localhost:%s/tokenize", port),
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

	// Calculate totals
	var subtotal float64
	lineItems := make([]models.LineItemResponse, len(req.LineItems))

	for i, li := range req.LineItems {
		var price float64
		fmt.Sscanf(li.Item.Price, "%f", &price)
		itemTotal := price * float64(li.Quantity)
		subtotal += itemTotal

		lineItems[i] = models.LineItemResponse{
			ID: generateID("li"),
			Item: models.ItemResponse{
				ID:          li.Item.ID,
				Name:        li.Item.Name,
				Description: li.Item.Description,
				Price:       li.Item.Price,
				ImageURL:    li.Item.ImageURL,
			},
			Quantity: li.Quantity,
			Totals: []models.TotalResponse{
				{Type: models.TotalTypeSubtotal, Amount: fmt.Sprintf("%.2f", itemTotal)},
			},
		}
	}

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
			{Type: models.TotalTypeSubtotal, Amount: fmt.Sprintf("%.2f", subtotal)},
			{Type: models.TotalTypeTax, Amount: fmt.Sprintf("%.2f", subtotal*0.0875)},
			{Type: models.TotalTypeTotal, Amount: fmt.Sprintf("%.2f", subtotal*1.0875)},
		},
		Links: []models.Link{
			{Rel: "terms", Href: "https://example.com/terms", Title: "Terms of Service"},
			{Rel: "privacy", Href: "https://example.com/privacy", Title: "Privacy Policy"},
		},
		Payment: models.PaymentResponse{
			Handlers: []models.PaymentHandlerResponse{
				{ID: "default", Type: "tokenization", Name: "Demo Payment Handler"},
			},
		},
		Messages: []models.Message{
			{Type: models.MessageTypeInfo, Title: "Buyer information required", Severity: models.SeverityRecoverable},
		},
	}

	// Store checkout
	mu.Lock()
	checkouts[checkoutID] = checkout
	mu.Unlock()

	log.Printf("Created checkout %s with %d items", checkoutID, len(lineItems))
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
		checkout.Buyer = &models.BuyerConsentResponse{
			Email:          req.Buyer.Email,
			Phone:          req.Buyer.Phone,
			FirstName:      req.Buyer.FirstName,
			LastName:       req.Buyer.LastName,
			BillingAddress: req.Buyer.BillingAddress,
		}
	}

	// Update fulfillment
	if req.Fulfillment != nil && req.Fulfillment.Destination != nil {
		checkout.Fulfillment = &models.FulfillmentResponse{
			Destination: &models.FulfillmentDestinationResponse{},
		}
		if req.Fulfillment.Destination.Shipping != nil {
			checkout.Fulfillment.Destination.Shipping = &models.ShippingDestinationResponse{
				Address: req.Fulfillment.Destination.Shipping.Address,
			}
		}
	}

	// Update payment
	if req.Payment != nil {
		checkout.Payment.HandlerID = req.Payment.HandlerID
		checkout.Payment.Status = "pending"
	}

	// Update status based on completeness
	checkout.Messages = nil
	if checkout.Buyer == nil || checkout.Buyer.Email == "" {
		checkout.Messages = append(checkout.Messages, models.Message{
			Type:     models.MessageTypeInfo,
			Title:    "Email required",
			Severity: models.SeverityRecoverable,
			Field:    "buyer.email",
		})
	}
	if checkout.Payment.HandlerID == "" {
		checkout.Messages = append(checkout.Messages, models.Message{
			Type:     models.MessageTypeInfo,
			Title:    "Payment required",
			Severity: models.SeverityRecoverable,
			Field:    "payment",
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
	mu.Lock()
	defer mu.Unlock()

	checkout, ok := checkouts[id]
	if !ok {
		return nil, server.NotFoundError("checkout not found")
	}

	if checkout.Status != models.CheckoutStatusReadyForComplete {
		return nil, server.BadRequestError("checkout is not ready for completion")
	}

	// Create order
	orderID := generateID("ord")
	now := time.Now()

	orderLineItems := make([]models.OrderLineItem, len(checkout.LineItems))
	for i, li := range checkout.LineItems {
		orderLineItems[i] = models.OrderLineItem{
			ID:       li.ID,
			Item:     li.Item,
			Quantity: li.Quantity,
			Totals:   li.Totals,
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
		Status:       models.OrderStatusConfirmed,
		LineItems:    orderLineItems,
		Currency:     checkout.Currency,
		Totals:       checkout.Totals,
		CreatedAt:    &now,
		PermalinkURL: fmt.Sprintf("https://example.com/orders/%s", orderID),
	}

	orders[orderID] = order

	// Update checkout
	checkout.Status = models.CheckoutStatusCompleted
	checkout.Order = &models.OrderConfirmation{
		ID:           orderID,
		PermalinkURL: order.PermalinkURL,
		CreatedAt:    &now,
	}
	checkout.Payment.Status = "captured"

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

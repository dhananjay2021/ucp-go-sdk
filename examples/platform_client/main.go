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

// Package main demonstrates using the UCP Go SDK as a platform client.
//
// This example shows how to:
// - Discover a merchant's UCP capabilities
// - Create a checkout session
// - Update the checkout with buyer information
// - Complete the checkout
// - Retrieve the order
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Universal-Commerce-Protocol/go-sdk/client"
	"github.com/Universal-Commerce-Protocol/go-sdk/extensions"
	"github.com/Universal-Commerce-Protocol/go-sdk/models"
)

func main() {
	// Get merchant URL from environment or use default
	merchantURL := os.Getenv("MERCHANT_URL")
	if merchantURL == "" {
		merchantURL = "http://localhost:8080"
	}

	// Create a UCP client
	ucpClient := client.NewClient(merchantURL,
		client.WithAPIKey(os.Getenv("API_KEY")),
		client.WithUserAgent("ucp-example-client/1.0"),
	)

	ctx := context.Background()

	// Step 1: Discover merchant capabilities
	fmt.Println("=== Step 1: Discovering merchant capabilities ===")
	profile, err := ucpClient.FetchProfile(ctx)
	if err != nil {
		log.Fatalf("Failed to fetch profile: %v", err)
	}

	fmt.Printf("Merchant UCP Version: %s\n", profile.UCP.Version)
	fmt.Printf("Supported capabilities:\n")
	for _, cap := range profile.UCP.Capabilities {
		fmt.Printf("  - %s (v%s)\n", cap.Name, cap.Version)
	}

	// Check for required capabilities
	if !client.HasCapability(profile, client.CapabilityCheckout) {
		log.Fatal("Merchant does not support checkout capability")
	}

	// Step 2: Create a checkout session
	fmt.Println("\n=== Step 2: Creating checkout session ===")
	checkout, err := ucpClient.CreateCheckout(ctx, &extensions.ExtendedCheckoutCreateRequest{
		LineItems: []models.LineItemCreateRequest{
			{
				Item: models.ItemCreateRequest{
					ID:          "PROD-001",
					Name:        "Wireless Headphones",
					Description: "Premium wireless headphones with noise cancellation",
					Price:       "149.99",
					ImageURL:    "https://example.com/images/headphones.jpg",
				},
				Quantity: 1,
			},
			{
				Item: models.ItemCreateRequest{
					ID:          "PROD-002",
					Name:        "Phone Case",
					Description: "Protective phone case",
					Price:       "29.99",
				},
				Quantity: 2,
			},
		},
		Currency: "USD",
	})
	if err != nil {
		log.Fatalf("Failed to create checkout: %v", err)
	}

	fmt.Printf("Checkout ID: %s\n", checkout.ID)
	fmt.Printf("Status: %s\n", checkout.Status)
	fmt.Printf("Line items: %d\n", len(checkout.LineItems))
	for _, total := range checkout.Totals {
		fmt.Printf("  %s: %s\n", total.Type, total.Amount)
	}

	// Step 3: Update checkout with buyer information
	fmt.Println("\n=== Step 3: Updating checkout with buyer info ===")
	checkout, err = ucpClient.UpdateCheckout(ctx, checkout.ID, &extensions.ExtendedCheckoutUpdateRequest{
		Buyer: &models.BuyerConsentUpdateRequest{
			Email:     "buyer@example.com",
			FirstName: "Jane",
			LastName:  "Doe",
			Phone:     "+1-555-123-4567",
			BillingAddress: &models.PostalAddress{
				AddressLines:       []string{"123 Main St", "Apt 4B"},
				Locality:           "San Francisco",
				AdministrativeArea: "CA",
				PostalCode:         "94102",
				CountryCode:        "US",
			},
		},
		Fulfillment: &models.FulfillmentUpdateRequest{
			Destination: &models.FulfillmentDestinationRequest{
				Shipping: &models.ShippingDestinationRequest{
					Address: &models.PostalAddress{
						AddressLines:       []string{"123 Main St", "Apt 4B"},
						Locality:           "San Francisco",
						AdministrativeArea: "CA",
						PostalCode:         "94102",
						CountryCode:        "US",
					},
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to update checkout: %v", err)
	}

	fmt.Printf("Status after update: %s\n", checkout.Status)

	// Step 4: Add payment (simulated - in real scenario, this would involve tokenization)
	fmt.Println("\n=== Step 4: Adding payment information ===")
	checkout, err = ucpClient.UpdateCheckout(ctx, checkout.ID, &extensions.ExtendedCheckoutUpdateRequest{
		Payment: &models.PaymentUpdateRequest{
			HandlerID: "default",
			Token:     "tok_test_simulated",
		},
	})
	if err != nil {
		log.Fatalf("Failed to add payment: %v", err)
	}

	fmt.Printf("Status after payment: %s\n", checkout.Status)

	// Step 5: Complete the checkout
	if checkout.Status == models.CheckoutStatusReadyForComplete {
		fmt.Println("\n=== Step 5: Completing checkout ===")
		checkout, err = ucpClient.CompleteCheckout(ctx, checkout.ID)
		if err != nil {
			log.Fatalf("Failed to complete checkout: %v", err)
		}

		fmt.Printf("Final status: %s\n", checkout.Status)
		if checkout.Order != nil {
			fmt.Printf("Order ID: %s\n", checkout.Order.ID)
			if checkout.Order.PermalinkURL != "" {
				fmt.Printf("Order URL: %s\n", checkout.Order.PermalinkURL)
			}
		}
	} else {
		fmt.Printf("\nCheckout not ready for completion. Status: %s\n", checkout.Status)
		if len(checkout.Messages) > 0 {
			fmt.Println("Messages:")
			for _, msg := range checkout.Messages {
				fmt.Printf("  [%s] %s: %s\n", msg.Type, msg.Title, msg.Detail)
			}
		}
	}

	// Step 6: Retrieve order (if completed)
	if checkout.Order != nil {
		fmt.Println("\n=== Step 6: Retrieving order ===")
		order, err := ucpClient.GetOrder(ctx, checkout.Order.ID)
		if err != nil {
			log.Printf("Failed to get order: %v", err)
		} else {
			fmt.Printf("Order status: %s\n", order.Status)
			fmt.Printf("Order items: %d\n", len(order.LineItems))
		}
	}

	fmt.Println("\n=== Example complete ===")
}

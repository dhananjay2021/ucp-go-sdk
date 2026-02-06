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

	"github.com/dhananjay2021/ucp-go-sdk/client"
	"github.com/dhananjay2021/ucp-go-sdk/extensions"
	"github.com/dhananjay2021/ucp-go-sdk/models"
)

func main() {
	// Get merchant URL from environment or use default
	merchantURL := os.Getenv("MERCHANT_URL")
	if merchantURL == "" {
		merchantURL = "http://localhost:8080"
	}

	// Create a UCP client with the required UCP-Agent header
	ucpClient := client.NewClient(merchantURL,
		client.WithAPIKey(os.Getenv("API_KEY")),
		client.WithUserAgent("ucp-example-client/1.0"),
		client.WithUCPAgent("https://example-platform.com/.well-known/ucp"), // Required: identifies the calling platform
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

	// Optional: Use Cart for pre-purchase exploration (if supported)
	hasCart := false
	for _, cap := range profile.UCP.Capabilities {
		if cap.Name == "dev.ucp.shopping.cart" {
			hasCart = true
			break
		}
	}

	if hasCart {
		fmt.Println("\n=== Cart Demo: Pre-purchase exploration ===")

		// Create a cart with buyer context
		cart, err := ucpClient.CreateCart(ctx, &models.CartCreateRequest{
			LineItems: []models.LineItemCreateRequest{
				{Item: models.ItemCreateRequest{ID: "PROD-001"}, Quantity: 1},
				{Item: models.ItemCreateRequest{ID: "PROD-002"}, Quantity: 3},
			},
			Context: &models.Context{
				AddressCountry: "US",
				AddressRegion:  "CA",
				Intent:         "comparing prices before buying",
			},
		})
		if err != nil {
			log.Printf("Cart creation failed: %v", err)
		} else {
			fmt.Printf("Cart ID: %s\n", cart.ID)
			fmt.Printf("Estimated total: %d cents\n", cart.Totals[len(cart.Totals)-1].Amount)
			for _, msg := range cart.Messages {
				fmt.Printf("  [%s] %s\n", msg.Type, msg.Content)
			}

			// Update cart quantities
			cart, err = ucpClient.UpdateCart(ctx, cart.ID, &models.CartUpdateRequest{
				ID: cart.ID,
				LineItems: []models.LineItemCreateRequest{
					{Item: models.ItemCreateRequest{ID: "PROD-001"}, Quantity: 2}, // Changed from 1 to 2
					{Item: models.ItemCreateRequest{ID: "PROD-002"}, Quantity: 1}, // Changed from 3 to 1
				},
			})
			if err != nil {
				log.Printf("Cart update failed: %v", err)
			} else {
				fmt.Printf("Updated cart total: %d cents\n", cart.Totals[len(cart.Totals)-1].Amount)
			}

			// Delete cart (we'll use checkout directly in this example)
			if err := ucpClient.DeleteCart(ctx, cart.ID); err != nil {
				log.Printf("Cart delete failed: %v", err)
			} else {
				fmt.Println("Cart deleted (proceeding to checkout)")
			}
		}
	}

	// Step 2: Create a checkout session
	// Note: In UCP, the platform sends only item IDs. The merchant
	// returns the full item details (title, price, etc.) in the response.
	fmt.Println("\n=== Step 2: Creating checkout session ===")
	checkout, err := ucpClient.CreateCheckout(ctx, &extensions.ExtendedCheckoutCreateRequest{
		LineItems: []models.LineItemCreateRequest{
			{
				Item:     models.ItemCreateRequest{ID: "PROD-001"},
				Quantity: 1,
			},
			{
				Item:     models.ItemCreateRequest{ID: "PROD-002"},
				Quantity: 2,
			},
		},
		Currency: "USD",
		Payment:  models.PaymentCreateRequest{},
		// Context provides buyer signals for localization and personalization
		Context: &models.Context{
			AddressCountry: "US",
			AddressRegion:  "CA",
			PostalCode:     "94043",
			Intent:         "looking for electronics accessories",
		},
	})
	if err != nil {
		log.Fatalf("Failed to create checkout: %v", err)
	}

	fmt.Printf("Checkout ID: %s\n", checkout.ID)
	fmt.Printf("Status: %s\n", checkout.Status)
	fmt.Printf("Line items: %d\n", len(checkout.LineItems))
	for _, total := range checkout.Totals {
		fmt.Printf("  %s: %d (cents)\n", total.Type, total.Amount)
	}

	// Step 3: Update checkout with buyer information
	fmt.Println("\n=== Step 3: Updating checkout with buyer info ===")
	checkout, err = ucpClient.UpdateCheckout(ctx, checkout.ID, &extensions.ExtendedCheckoutUpdateRequest{
		ID: checkout.ID,
		LineItems: []models.LineItemUpdateRequest{
			{Item: models.ItemUpdateRequest{ID: "PROD-001"}, Quantity: 1},
			{Item: models.ItemUpdateRequest{ID: "PROD-002"}, Quantity: 2},
		},
		Currency: "USD",
		Payment:  models.PaymentUpdateRequest{},
		Buyer: &models.BuyerWithConsentUpdateRequest{
			Email:       "buyer@example.com",
			FirstName:   "Jane",
			LastName:    "Doe",
			PhoneNumber: "+1-555-123-4567",
		},
		Fulfillment: &models.FulfillmentUpdateRequest{
			Methods: []models.FulfillmentMethodUpdateRequest{
				{
					ID:          "ship-1",
					LineItemIDs: []string{"PROD-001", "PROD-002"},
					Destinations: []models.FulfillmentDestinationRequest{
						{
							PostalAddress: models.PostalAddress{
								StreetAddress:   "123 Main St",
								ExtendedAddress: "Apt 4B",
								AddressLocality: "San Francisco",
								AddressRegion:   "CA",
								PostalCode:      "94102",
								AddressCountry:  "US",
							},
						},
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
		ID: checkout.ID,
		LineItems: []models.LineItemUpdateRequest{
			{Item: models.ItemUpdateRequest{ID: "PROD-001"}, Quantity: 1},
			{Item: models.ItemUpdateRequest{ID: "PROD-002"}, Quantity: 2},
		},
		Currency: "USD",
		Payment: models.PaymentUpdateRequest{
			SelectedInstrumentID: "pi-test-001",
			Instruments: []models.PaymentInstrument{
				{
					ID:         "pi-test-001",
					HandlerID:  "default",
					Type:       models.PaymentInstrumentTypeCard,
					Brand:      "visa",
					LastDigits: "4242",
				},
			},
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
				fmt.Printf("  [%s] %s\n", msg.Type, msg.Content)
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
			fmt.Printf("Order items: %d\n", len(order.LineItems))
		}
	}

	fmt.Println("\n=== Example complete ===")
}

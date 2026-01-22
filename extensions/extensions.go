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

// Package extensions provides extended types that compose UCP base schemas with extensions.
package extensions

import (
	"time"

	"github.com/dhananjay2021/ucp-go-sdk/models"
)

// ExtendedPaymentCredential extends PaymentCredential with an optional token field.
type ExtendedPaymentCredential struct {
	models.CardCredential

	// Token is an optional pre-existing token.
	Token string `json:"token,omitempty"`
}

// PlatformConfig contains platform-specific configuration.
type PlatformConfig struct {
	// WebhookURL is the URL for webhook notifications.
	WebhookURL string `json:"webhook_url,omitempty"`
}

// ExtendedCheckoutResponse combines base checkout with fulfillment, discounts, and buyer consent.
type ExtendedCheckoutResponse struct {
	// UCP contains protocol metadata.
	UCP models.ResponseCheckout `json:"ucp"`

	// ID is the unique identifier of the checkout session.
	ID string `json:"id"`

	// LineItems are the items being checked out.
	LineItems []models.LineItemResponse `json:"line_items"`

	// Buyer contains buyer information with consent.
	Buyer *models.BuyerWithConsentResponse `json:"buyer,omitempty"`

	// Status is the current checkout state.
	Status models.CheckoutStatus `json:"status"`

	// Currency is the ISO 4217 currency code.
	Currency string `json:"currency"`

	// Totals contains the cart totals breakdown.
	Totals []models.TotalResponse `json:"totals"`

	// Messages contains error and info messages.
	Messages []models.Message `json:"messages,omitempty"`

	// Links are URLs to be displayed by the platform.
	Links []models.Link `json:"links"`

	// ExpiresAt is the RFC 3339 expiry timestamp.
	ExpiresAt *time.Time `json:"expires_at,omitempty"`

	// ContinueURL is for checkout handoff and session recovery.
	ContinueURL string `json:"continue_url,omitempty"`

	// Payment contains payment information.
	Payment models.PaymentResponse `json:"payment"`

	// Order contains details about an order created for this checkout.
	Order *models.OrderConfirmation `json:"order,omitempty"`

	// Fulfillment contains fulfillment information (extension).
	Fulfillment *models.FulfillmentResponse `json:"fulfillment,omitempty"`

	// Discounts contains applied discounts (extension).
	Discounts *models.DiscountsResponse `json:"discounts,omitempty"`

	// Platform contains platform configuration.
	Platform *PlatformConfig `json:"platform,omitempty"`
}

// ExtendedCheckoutCreateRequest combines base checkout create with extensions.
type ExtendedCheckoutCreateRequest struct {
	// LineItems are the items to checkout.
	LineItems []models.LineItemCreateRequest `json:"line_items"`

	// Currency is the ISO 4217 currency code.
	Currency string `json:"currency"`

	// Payment contains payment information.
	Payment models.PaymentCreateRequest `json:"payment"`

	// Buyer contains buyer information with consent (extension).
	Buyer *models.BuyerWithConsentCreateRequest `json:"buyer,omitempty"`

	// Fulfillment contains fulfillment information (extension).
	Fulfillment *models.FulfillmentCreateRequest `json:"fulfillment,omitempty"`

	// Discounts contains discount codes to apply (extension).
	Discounts *models.DiscountsCreateRequest `json:"discounts,omitempty"`
}

// ExtendedCheckoutUpdateRequest combines base checkout update with extensions.
type ExtendedCheckoutUpdateRequest struct {
	// ID is the unique identifier of the checkout session.
	ID string `json:"id"`

	// LineItems are the line items.
	LineItems []models.LineItemUpdateRequest `json:"line_items"`

	// Currency is the ISO 4217 currency code.
	Currency string `json:"currency"`

	// Payment contains payment information.
	Payment models.PaymentUpdateRequest `json:"payment"`

	// Buyer contains buyer information with consent (extension).
	Buyer *models.BuyerWithConsentUpdateRequest `json:"buyer,omitempty"`

	// Fulfillment contains fulfillment information (extension).
	Fulfillment *models.FulfillmentUpdateRequest `json:"fulfillment,omitempty"`

	// Discounts contains discount updates (extension).
	Discounts *models.DiscountsUpdateRequest `json:"discounts,omitempty"`
}

// ExtendedOrder combines base order with extensions.
type ExtendedOrder struct {
	models.Order

	// Discounts contains applied discounts.
	Discounts *models.DiscountsResponse `json:"discounts,omitempty"`
}

// CheckoutWithFulfillmentCreateRequest is a checkout create request with fulfillment.
type CheckoutWithFulfillmentCreateRequest struct {
	models.CheckoutCreateRequest

	// Fulfillment contains fulfillment information.
	Fulfillment *models.FulfillmentCreateRequest `json:"fulfillment,omitempty"`
}

// CheckoutWithFulfillmentUpdateRequest is a checkout update request with fulfillment.
type CheckoutWithFulfillmentUpdateRequest struct {
	models.CheckoutUpdateRequest

	// Fulfillment contains fulfillment information.
	Fulfillment *models.FulfillmentUpdateRequest `json:"fulfillment,omitempty"`
}

// CheckoutWithFulfillmentResponse is a checkout response with fulfillment.
type CheckoutWithFulfillmentResponse struct {
	models.CheckoutResponse

	// Fulfillment contains fulfillment information.
	Fulfillment *models.FulfillmentResponse `json:"fulfillment,omitempty"`
}

// CheckoutWithDiscountCreateRequest is a checkout create request with discounts.
type CheckoutWithDiscountCreateRequest struct {
	models.CheckoutCreateRequest

	// Discounts contains discount codes to apply.
	Discounts *models.DiscountsCreateRequest `json:"discounts,omitempty"`
}

// CheckoutWithDiscountUpdateRequest is a checkout update request with discounts.
type CheckoutWithDiscountUpdateRequest struct {
	models.CheckoutUpdateRequest

	// Discounts contains discount updates.
	Discounts *models.DiscountsUpdateRequest `json:"discounts,omitempty"`
}

// CheckoutWithDiscountResponse is a checkout response with discounts.
type CheckoutWithDiscountResponse struct {
	models.CheckoutResponse

	// Discounts contains applied discounts.
	Discounts *models.DiscountsResponse `json:"discounts,omitempty"`
}

// CheckoutWithBuyerConsentCreateRequest is a checkout create request with buyer consent.
type CheckoutWithBuyerConsentCreateRequest struct {
	models.CheckoutCreateRequest

	// Buyer contains buyer consent information.
	Buyer *models.BuyerWithConsentCreateRequest `json:"buyer,omitempty"`
}

// CheckoutWithBuyerConsentUpdateRequest is a checkout update request with buyer consent.
type CheckoutWithBuyerConsentUpdateRequest struct {
	models.CheckoutUpdateRequest

	// Buyer contains buyer consent information.
	Buyer *models.BuyerWithConsentUpdateRequest `json:"buyer,omitempty"`
}

// CheckoutWithBuyerConsentResponse is a checkout response with buyer consent.
type CheckoutWithBuyerConsentResponse struct {
	models.CheckoutResponse

	// Buyer contains buyer consent information.
	Buyer *models.BuyerWithConsentResponse `json:"buyer,omitempty"`
}

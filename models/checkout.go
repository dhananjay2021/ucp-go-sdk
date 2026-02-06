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

package models

import "time"

// CheckoutCreateRequest represents a request to create a checkout session.
type CheckoutCreateRequest struct {
	// LineItems are the items to checkout.
	LineItems []LineItemCreateRequest `json:"line_items"`

	// Currency is the ISO 4217 currency code.
	Currency string `json:"currency"`

	// Payment contains payment information.
	Payment PaymentCreateRequest `json:"payment"`

	// Buyer contains optional buyer information.
	Buyer *Buyer `json:"buyer,omitempty"`

	// Context provides buyer signals for localization (country, region, postal_code, intent).
	Context *Context `json:"context,omitempty"`
}

// CheckoutUpdateRequest represents a request to update a checkout session.
type CheckoutUpdateRequest struct {
	// ID is the unique identifier of the checkout session.
	ID string `json:"id"`

	// LineItems are the line items being checked out.
	LineItems []LineItemUpdateRequest `json:"line_items"`

	// Currency is the ISO 4217 currency code.
	Currency string `json:"currency"`

	// Payment contains payment information.
	Payment PaymentUpdateRequest `json:"payment"`

	// Buyer contains optional buyer information.
	Buyer *Buyer `json:"buyer,omitempty"`

	// Context provides buyer signals for localization.
	Context *Context `json:"context,omitempty"`
}

// CheckoutResponse represents a checkout session response.
type CheckoutResponse struct {
	// UCP contains protocol metadata.
	UCP ResponseCheckout `json:"ucp"`

	// ID is the unique identifier of the checkout session.
	ID string `json:"id"`

	// LineItems are the items being checked out.
	LineItems []LineItemResponse `json:"line_items"`

	// Status is the current checkout state.
	Status CheckoutStatus `json:"status"`

	// Currency is the ISO 4217 currency code.
	Currency string `json:"currency"`

	// Totals contains the cart totals breakdown.
	Totals []TotalResponse `json:"totals"`

	// Links are URLs to be displayed by the platform.
	Links []Link `json:"links"`

	// Payment contains payment information.
	Payment PaymentResponse `json:"payment"`

	// Buyer contains buyer information.
	Buyer *Buyer `json:"buyer,omitempty"`

	// Messages contains error and info messages.
	Messages []Message `json:"messages,omitempty"`

	// ExpiresAt is the RFC 3339 expiry timestamp.
	ExpiresAt *time.Time `json:"expires_at,omitempty"`

	// ContinueURL is for checkout handoff and session recovery.
	ContinueURL string `json:"continue_url,omitempty"`

	// Order contains details about an order created for this checkout.
	Order *OrderConfirmation `json:"order,omitempty"`

	// EmbeddedConfig provides per-checkout configuration for embedded transport binding.
	// Allows businesses to vary ECP availability and delegations.
	EmbeddedConfig *EmbeddedTransportConfig `json:"embedded_config,omitempty"`

	// Context provides buyer signals used for this checkout.
	Context *Context `json:"context,omitempty"`
}

// CheckoutCompleteRequest represents a request to complete a checkout.
type CheckoutCompleteRequest struct {
	// This is intentionally empty for the base checkout.
	// Extensions like AP2 add fields via composition.
}

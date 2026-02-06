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

// CartCreateRequest represents a request to create a new cart session.
// Carts provide lightweight pre-purchase exploration with estimated pricing.
type CartCreateRequest struct {
	// LineItems are the items to add to the cart (required).
	LineItems []LineItemCreateRequest `json:"line_items"`

	// Context provides buyer signals for localization (country, region, postal_code).
	// Used by merchant for pricing, availability, and currency. Falls back to geo-IP if omitted.
	Context *Context `json:"context,omitempty"`

	// Buyer provides optional buyer information for personalized estimates.
	Buyer *Buyer `json:"buyer,omitempty"`
}

// CartUpdateRequest represents a request to update an existing cart.
type CartUpdateRequest struct {
	// ID is the unique cart identifier (required).
	ID string `json:"id"`

	// LineItems are the updated cart items (full replacement).
	LineItems []LineItemCreateRequest `json:"line_items"`

	// Context provides updated buyer signals for localization.
	Context *Context `json:"context,omitempty"`

	// Buyer provides updated buyer information.
	Buyer *Buyer `json:"buyer,omitempty"`
}

// ResponseCart represents UCP metadata for cart responses.
type ResponseCart struct {
	// Schema is the schema URL.
	Schema string `json:"$schema,omitempty"`

	// Capabilities contains information about supported capabilities.
	Capabilities map[string]interface{} `json:"capabilities,omitempty"`
}

// CartResponse represents the response from cart operations.
type CartResponse struct {
	// UCP contains the UCP response metadata.
	UCP *ResponseCart `json:"ucp,omitempty"`

	// ID is the unique cart identifier.
	ID string `json:"id"`

	// LineItems are the cart items with pricing.
	LineItems []LineItemResponse `json:"line_items,omitempty"`

	// Currency is the ISO 4217 currency code determined by the merchant.
	Currency string `json:"currency"`

	// Totals contains the estimated cost breakdown.
	// May be partial if shipping/tax not yet calculable.
	Totals []TotalResponse `json:"totals,omitempty"`

	// Messages contains validation messages, warnings, or informational notices.
	Messages []Message `json:"messages,omitempty"`

	// Links contains optional merchant links (policies, FAQs).
	Links []Link `json:"links,omitempty"`

	// ContinueURL provides a URL for cart handoff and session recovery.
	// Enables sharing and human-in-the-loop flows.
	ContinueURL string `json:"continue_url,omitempty"`

	// ExpiresAt is the cart expiry timestamp (RFC 3339).
	ExpiresAt string `json:"expires_at,omitempty"`
}

// CartWithCheckout extends CheckoutCreateRequest to support cart-to-checkout conversion.
type CartWithCheckout struct {
	CheckoutCreateRequest

	// CartID is the cart ID to convert to checkout.
	// When specified, business MUST use cart contents (line_items, context, buyer)
	// and MUST ignore overlapping fields in checkout payload.
	CartID string `json:"cart_id,omitempty"`
}

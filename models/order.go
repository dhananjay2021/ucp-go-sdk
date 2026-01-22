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

// OrderLineItemStatus represents the fulfillment status of an order line item.
type OrderLineItemStatus string

const (
	// OrderLineItemStatusProcessing indicates the item is being processed.
	OrderLineItemStatusProcessing OrderLineItemStatus = "processing"

	// OrderLineItemStatusPartial indicates partial fulfillment.
	OrderLineItemStatusPartial OrderLineItemStatus = "partial"

	// OrderLineItemStatusFulfilled indicates the item is fully fulfilled.
	OrderLineItemStatusFulfilled OrderLineItemStatus = "fulfilled"
)

// OrderLineItemQuantity represents quantity tracking for an order line item.
type OrderLineItemQuantity struct {
	// Total is the current total quantity.
	Total int `json:"total"`

	// Fulfilled is the quantity fulfilled (sum from fulfillment events).
	Fulfilled int `json:"fulfilled"`
}

// OrderLineItem represents a line item in an order.
type OrderLineItem struct {
	// ID is the line item identifier.
	ID string `json:"id"`

	// Item contains the product data.
	Item ItemResponse `json:"item"`

	// Quantity tracks total and fulfilled quantities.
	Quantity OrderLineItemQuantity `json:"quantity"`

	// Totals contains the line item totals breakdown.
	Totals []TotalResponse `json:"totals"`

	// Status is the derived fulfillment status.
	Status OrderLineItemStatus `json:"status"`

	// ParentID is the parent line item identifier for nested structures.
	ParentID string `json:"parent_id,omitempty"`
}

// ExpectationLineItem represents a line item reference in an expectation.
type ExpectationLineItem struct {
	// ID is the line item ID reference.
	ID string `json:"id"`

	// Quantity is the quantity of this item in this expectation.
	Quantity int `json:"quantity"`
}

// Expectation represents a buyer-facing fulfillment expectation.
type Expectation struct {
	// ID is the expectation identifier.
	ID string `json:"id"`

	// LineItems specifies which line items and quantities are in this expectation.
	LineItems []ExpectationLineItem `json:"line_items"`

	// MethodType is the delivery method type (shipping, pickup, digital).
	MethodType MethodType `json:"method_type"`

	// Destination is the delivery destination address.
	Destination PostalAddress `json:"destination"`

	// Description is a human-readable delivery description.
	Description string `json:"description,omitempty"`

	// FulfillableOn indicates when this expectation can be fulfilled.
	FulfillableOn string `json:"fulfillable_on,omitempty"`
}

// FulfillmentEventLineItem represents a line item reference in a fulfillment event.
type FulfillmentEventLineItem struct {
	// ID is the line item ID reference.
	ID string `json:"id"`

	// Quantity is the quantity fulfilled in this event.
	Quantity int `json:"quantity"`
}

// FulfillmentEvent represents an append-only fulfillment event.
type FulfillmentEvent struct {
	// ID is the fulfillment event identifier.
	ID string `json:"id"`

	// OccurredAt is when this fulfillment event occurred.
	OccurredAt time.Time `json:"occurred_at"`

	// Type is the fulfillment event type (processing, shipped, delivered, etc.).
	Type string `json:"type"`

	// LineItems specifies which line items and quantities are fulfilled.
	LineItems []FulfillmentEventLineItem `json:"line_items"`

	// TrackingNumber is the carrier tracking number.
	TrackingNumber string `json:"tracking_number,omitempty"`

	// TrackingURL is the URL to track this shipment.
	TrackingURL string `json:"tracking_url,omitempty"`

	// Carrier is the carrier name.
	Carrier string `json:"carrier,omitempty"`

	// Description is a human-readable shipment status.
	Description string `json:"description,omitempty"`
}

// AdjustmentLineItem represents a line item reference in an adjustment.
type AdjustmentLineItem struct {
	// ID is the line item ID reference.
	ID string `json:"id"`

	// Quantity is the quantity affected by this adjustment.
	Quantity int `json:"quantity"`
}

// Adjustment represents an append-only adjustment event.
type Adjustment struct {
	// ID is the adjustment event identifier.
	ID string `json:"id"`

	// Type is the adjustment type (refund, return, credit, etc.).
	Type string `json:"type"`

	// OccurredAt is when this adjustment occurred.
	OccurredAt time.Time `json:"occurred_at"`

	// Status is the adjustment status.
	Status AdjustmentStatus `json:"status"`

	// LineItems specifies which line items are affected (optional).
	LineItems []AdjustmentLineItem `json:"line_items,omitempty"`

	// Amount is the amount in minor units (cents) for refunds, credits, etc.
	Amount int `json:"amount,omitempty"`

	// Description is a human-readable reason or description.
	Description string `json:"description,omitempty"`
}

// OrderFulfillment represents fulfillment data in an order.
type OrderFulfillment struct {
	// Expectations are buyer-facing fulfillment expectations.
	Expectations []Expectation `json:"expectations,omitempty"`

	// Events are append-only fulfillment events.
	Events []FulfillmentEvent `json:"events,omitempty"`
}

// Order represents an order.
type Order struct {
	// UCP contains protocol metadata.
	UCP ResponseOrder `json:"ucp"`

	// ID is the order identifier.
	ID string `json:"id"`

	// CheckoutID is the associated checkout session ID.
	CheckoutID string `json:"checkout_id"`

	// PermalinkURL is a permalink to access the order on merchant site.
	PermalinkURL string `json:"permalink_url"`

	// LineItems are the immutable line items in the order.
	LineItems []OrderLineItem `json:"line_items"`

	// Fulfillment contains fulfillment expectations and events.
	Fulfillment OrderFulfillment `json:"fulfillment"`

	// Totals contains the order totals.
	Totals []TotalResponse `json:"totals"`

	// Adjustments lists order adjustments (refunds, returns, etc.).
	Adjustments []Adjustment `json:"adjustments,omitempty"`
}

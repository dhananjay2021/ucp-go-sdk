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

// OrderStatus represents the status of an order.
type OrderStatus string

const (
	// OrderStatusPending indicates the order is pending.
	OrderStatusPending OrderStatus = "pending"

	// OrderStatusConfirmed indicates the order is confirmed.
	OrderStatusConfirmed OrderStatus = "confirmed"

	// OrderStatusProcessing indicates the order is being processed.
	OrderStatusProcessing OrderStatus = "processing"

	// OrderStatusShipped indicates the order has been shipped.
	OrderStatusShipped OrderStatus = "shipped"

	// OrderStatusDelivered indicates the order has been delivered.
	OrderStatusDelivered OrderStatus = "delivered"

	// OrderStatusCanceled indicates the order has been canceled.
	OrderStatusCanceled OrderStatus = "canceled"

	// OrderStatusReturned indicates the order has been returned.
	OrderStatusReturned OrderStatus = "returned"
)

// OrderLineItem represents a line item in an order.
type OrderLineItem struct {
	// ID is the line item identifier.
	ID string `json:"id"`

	// Item contains the item details.
	Item ItemResponse `json:"item"`

	// Quantity is the number of items.
	Quantity int `json:"quantity"`

	// Totals contains the line item totals.
	Totals []TotalResponse `json:"totals,omitempty"`

	// FulfillmentStatus is the fulfillment status for this item.
	FulfillmentStatus string `json:"fulfillment_status,omitempty"`

	// TrackingNumber is the shipping tracking number.
	TrackingNumber string `json:"tracking_number,omitempty"`

	// TrackingURL is the tracking URL.
	TrackingURL string `json:"tracking_url,omitempty"`
}

// Adjustment represents an order adjustment (refund, return, etc).
type Adjustment struct {
	// ID is the adjustment identifier.
	ID string `json:"id"`

	// Type is the adjustment type (refund, return, etc).
	Type string `json:"type"`

	// Status is the adjustment status.
	Status AdjustmentStatus `json:"status"`

	// Amount is the adjustment amount.
	Amount string `json:"amount,omitempty"`

	// Currency is the adjustment currency.
	Currency string `json:"currency,omitempty"`

	// Reason is the reason for the adjustment.
	Reason string `json:"reason,omitempty"`

	// LineItemIDs are the affected line items.
	LineItemIDs []string `json:"line_item_ids,omitempty"`

	// CreatedAt is when the adjustment was created.
	CreatedAt *time.Time `json:"created_at,omitempty"`

	// UpdatedAt is when the adjustment was last updated.
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// Order represents an order.
type Order struct {
	// UCP contains protocol metadata.
	UCP ResponseOrder `json:"ucp"`

	// ID is the order identifier.
	ID string `json:"id"`

	// CheckoutID is the associated checkout session ID.
	CheckoutID string `json:"checkout_id,omitempty"`

	// Status is the order status.
	Status OrderStatus `json:"status"`

	// LineItems are the items in the order.
	LineItems []OrderLineItem `json:"line_items"`

	// Buyer contains buyer information.
	Buyer *Buyer `json:"buyer,omitempty"`

	// Currency is the order currency.
	Currency string `json:"currency"`

	// Totals contains the order totals.
	Totals []TotalResponse `json:"totals"`

	// Fulfillment contains fulfillment details.
	Fulfillment *FulfillmentResponse `json:"fulfillment,omitempty"`

	// Payment contains payment details.
	Payment *PaymentResponse `json:"payment,omitempty"`

	// Adjustments lists order adjustments.
	Adjustments []Adjustment `json:"adjustments,omitempty"`

	// CreatedAt is when the order was created.
	CreatedAt *time.Time `json:"created_at,omitempty"`

	// UpdatedAt is when the order was last updated.
	UpdatedAt *time.Time `json:"updated_at,omitempty"`

	// PermalinkURL is a URL to view the order.
	PermalinkURL string `json:"permalink_url,omitempty"`

	// Messages contains order-related messages.
	Messages []Message `json:"messages,omitempty"`

	// Metadata contains custom metadata.
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

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

// CheckoutStatus represents the state of a checkout session.
type CheckoutStatus string

const (
	// CheckoutStatusIncomplete indicates the checkout is missing required data.
	CheckoutStatusIncomplete CheckoutStatus = "incomplete"

	// CheckoutStatusRequiresEscalation indicates buyer input or review is needed.
	CheckoutStatusRequiresEscalation CheckoutStatus = "requires_escalation"

	// CheckoutStatusReadyForComplete indicates the checkout can be completed.
	CheckoutStatusReadyForComplete CheckoutStatus = "ready_for_complete"

	// CheckoutStatusCompleteInProgress indicates completion is in progress.
	CheckoutStatusCompleteInProgress CheckoutStatus = "complete_in_progress"

	// CheckoutStatusCompleted indicates the checkout has been completed.
	CheckoutStatusCompleted CheckoutStatus = "completed"

	// CheckoutStatusCanceled indicates the checkout has been canceled.
	CheckoutStatusCanceled CheckoutStatus = "canceled"
)

// MessageType represents the type of a checkout message.
type MessageType string

const (
	// MessageTypeError indicates an error message.
	MessageTypeError MessageType = "error"

	// MessageTypeWarning indicates a warning message.
	MessageTypeWarning MessageType = "warning"

	// MessageTypeInfo indicates an informational message.
	MessageTypeInfo MessageType = "info"
)

// Severity indicates who resolves an error.
type Severity string

const (
	// SeverityRecoverable indicates the agent can fix via API.
	SeverityRecoverable Severity = "recoverable"

	// SeverityRequiresBuyerInput indicates merchant requires information their API doesn't support.
	SeverityRequiresBuyerInput Severity = "requires_buyer_input"

	// SeverityRequiresBuyerReview indicates buyer must authorize before order placement.
	SeverityRequiresBuyerReview Severity = "requires_buyer_review"
)

// ContentType represents the content format.
type ContentType string

const (
	// ContentTypePlain indicates plain text content.
	ContentTypePlain ContentType = "plain"

	// ContentTypeMarkdown indicates markdown content.
	ContentTypeMarkdown ContentType = "markdown"
)

// TotalType represents the type of total categorization.
type TotalType string

const (
	// TotalTypeSubtotal is the subtotal before taxes and fees.
	TotalTypeSubtotal TotalType = "subtotal"

	// TotalTypeTax is the tax amount.
	TotalTypeTax TotalType = "tax"

	// TotalTypeFee is a fee amount.
	TotalTypeFee TotalType = "fee"

	// TotalTypeDiscount is a discount amount.
	TotalTypeDiscount TotalType = "discount"

	// TotalTypeFulfillment is the fulfillment/shipping cost.
	TotalTypeFulfillment TotalType = "fulfillment"

	// TotalTypeItemsDiscount is discount on items.
	TotalTypeItemsDiscount TotalType = "items_discount"

	// TotalTypeTotal is the final total.
	TotalTypeTotal TotalType = "total"
)

// MethodType represents the delivery method type.
type MethodType string

const (
	// MethodTypeShipping indicates shipping delivery.
	MethodTypeShipping MethodType = "shipping"

	// MethodTypePickup indicates in-store pickup.
	MethodTypePickup MethodType = "pickup"

	// MethodTypeDigital indicates digital delivery.
	MethodTypeDigital MethodType = "digital"
)

// CardNumberType represents the type of card number.
type CardNumberType string

const (
	// CardNumberTypeFPAN is a Funding Primary Account Number.
	CardNumberTypeFPAN CardNumberType = "fpan"

	// CardNumberTypeDPAN is a Device Primary Account Number.
	CardNumberTypeDPAN CardNumberType = "dpan"

	// CardNumberTypeNetworkToken is a network token.
	CardNumberTypeNetworkToken CardNumberType = "network_token"
)

// AdjustmentStatus represents the status of an adjustment (refund, return, etc).
type AdjustmentStatus string

const (
	// AdjustmentStatusPending indicates the adjustment is pending.
	AdjustmentStatusPending AdjustmentStatus = "pending"

	// AdjustmentStatusCompleted indicates the adjustment is completed.
	AdjustmentStatusCompleted AdjustmentStatus = "completed"

	// AdjustmentStatusFailed indicates the adjustment failed.
	AdjustmentStatusFailed AdjustmentStatus = "failed"
)

// Link represents a link to be displayed by the platform.
type Link struct {
	// Rel is the link relation type.
	Rel string `json:"rel"`

	// Href is the URL of the link.
	Href string `json:"href"`

	// Title is an optional display title.
	Title string `json:"title,omitempty"`
}

// Message represents an error, warning, or info message.
type Message struct {
	// Type is the message type (error, warning, info).
	Type MessageType `json:"type"`

	// Code is a machine-readable error code.
	Code string `json:"code,omitempty"`

	// Title is a short summary of the message.
	Title string `json:"title"`

	// Detail provides additional context.
	Detail string `json:"detail,omitempty"`

	// ContentType indicates the format of the detail (plain, markdown).
	ContentType ContentType `json:"content_type,omitempty"`

	// Severity indicates who can resolve this issue.
	Severity Severity `json:"severity,omitempty"`

	// Field is the field this message relates to.
	Field string `json:"field,omitempty"`

	// ContinueURL is a URL for resolving the issue.
	ContinueURL string `json:"continue_url,omitempty"`
}

// TotalResponse represents a total amount breakdown.
type TotalResponse struct {
	// Type is the categorization of this total.
	Type TotalType `json:"type"`

	// Amount is the monetary value as a string.
	Amount string `json:"amount"`

	// Label is an optional display label.
	Label string `json:"label,omitempty"`
}

// TotalCreateRequest represents a total in a create request.
type TotalCreateRequest struct {
	// Type is the categorization of this total.
	Type TotalType `json:"type"`

	// Amount is the monetary value as a string.
	Amount string `json:"amount"`

	// Label is an optional display label.
	Label string `json:"label,omitempty"`
}

// PostalAddress represents a postal/mailing address.
type PostalAddress struct {
	// AddressLines are the street address lines.
	AddressLines []string `json:"address_lines,omitempty"`

	// Locality is the city/town.
	Locality string `json:"locality,omitempty"`

	// AdministrativeArea is the state/province/region.
	AdministrativeArea string `json:"administrative_area,omitempty"`

	// PostalCode is the ZIP/postal code.
	PostalCode string `json:"postal_code,omitempty"`

	// CountryCode is the ISO 3166-1 alpha-2 country code.
	CountryCode string `json:"country_code,omitempty"`

	// Recipients are the names of recipients.
	Recipients []string `json:"recipients,omitempty"`

	// Organization is the company/organization name.
	Organization string `json:"organization,omitempty"`

	// PhoneNumber is a contact phone number.
	PhoneNumber string `json:"phone_number,omitempty"`
}

// ItemResponse represents an item in a line item response.
type ItemResponse struct {
	// ID is a unique identifier for the item.
	ID string `json:"id,omitempty"`

	// Name is the display name of the item.
	Name string `json:"name"`

	// Description provides additional details.
	Description string `json:"description,omitempty"`

	// Price is the unit price as a string.
	Price string `json:"price"`

	// ImageURL is a URL to an item image.
	ImageURL string `json:"image_url,omitempty"`

	// ProductURL is a URL to the product page.
	ProductURL string `json:"product_url,omitempty"`

	// SKU is the stock keeping unit.
	SKU string `json:"sku,omitempty"`

	// Attributes contains item-specific attributes.
	Attributes map[string]string `json:"attributes,omitempty"`
}

// ItemCreateRequest represents an item in a create request.
type ItemCreateRequest struct {
	// ID is a unique identifier for the item.
	ID string `json:"id,omitempty"`

	// Name is the display name of the item.
	Name string `json:"name"`

	// Description provides additional details.
	Description string `json:"description,omitempty"`

	// Price is the unit price as a string.
	Price string `json:"price"`

	// ImageURL is a URL to an item image.
	ImageURL string `json:"image_url,omitempty"`

	// ProductURL is a URL to the product page.
	ProductURL string `json:"product_url,omitempty"`

	// SKU is the stock keeping unit.
	SKU string `json:"sku,omitempty"`

	// Attributes contains item-specific attributes.
	Attributes map[string]string `json:"attributes,omitempty"`
}

// LineItemResponse represents a line item in a checkout response.
type LineItemResponse struct {
	// ID is a unique identifier for the line item.
	ID string `json:"id"`

	// Item contains the item details.
	Item ItemResponse `json:"item"`

	// Quantity is the number of items.
	Quantity int `json:"quantity"`

	// Totals contains the line item totals breakdown.
	Totals []TotalResponse `json:"totals"`

	// ParentID is the parent line item identifier for nested structures.
	ParentID string `json:"parent_id,omitempty"`
}

// LineItemCreateRequest represents a line item in a create request.
type LineItemCreateRequest struct {
	// ID is an optional client-provided identifier.
	ID string `json:"id,omitempty"`

	// Item contains the item details.
	Item ItemCreateRequest `json:"item"`

	// Quantity is the number of items.
	Quantity int `json:"quantity"`

	// Totals contains optional totals.
	Totals []TotalCreateRequest `json:"totals,omitempty"`

	// ParentID is the parent line item identifier for nested structures.
	ParentID string `json:"parent_id,omitempty"`
}

// LineItemUpdateRequest represents a line item update.
type LineItemUpdateRequest struct {
	// ID is the line item identifier.
	ID string `json:"id"`

	// Item contains updated item details.
	Item *ItemCreateRequest `json:"item,omitempty"`

	// Quantity is the updated quantity.
	Quantity *int `json:"quantity,omitempty"`
}

// Buyer represents information about the buyer.
type Buyer struct {
	// Email is the buyer's email address.
	Email string `json:"email,omitempty"`

	// Phone is the buyer's phone number.
	Phone string `json:"phone,omitempty"`

	// FirstName is the buyer's first name.
	FirstName string `json:"first_name,omitempty"`

	// LastName is the buyer's last name.
	LastName string `json:"last_name,omitempty"`

	// BillingAddress is the buyer's billing address.
	BillingAddress *PostalAddress `json:"billing_address,omitempty"`
}

// OrderConfirmation contains details about an order created for a checkout.
type OrderConfirmation struct {
	// ID is the order identifier.
	ID string `json:"id"`

	// PermalinkURL is a URL to view the order.
	PermalinkURL string `json:"permalink_url,omitempty"`

	// CreatedAt is when the order was created.
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

// Expectation represents a delivery expectation.
type Expectation struct {
	// MinDays is the minimum number of days.
	MinDays *int `json:"min_days,omitempty"`

	// MaxDays is the maximum number of days.
	MaxDays *int `json:"max_days,omitempty"`

	// DeliveryDate is a specific delivery date.
	DeliveryDate string `json:"delivery_date,omitempty"`

	// Description is a human-readable description.
	Description string `json:"description,omitempty"`
}

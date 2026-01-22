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
	// Type is the link type (e.g., privacy_policy, terms_of_service, refund_policy).
	Type string `json:"type"`

	// URL is the actual URL pointing to the content.
	URL string `json:"url"`

	// Title is an optional display text for the link.
	Title string `json:"title,omitempty"`
}

// Message represents an error, warning, or info message.
type Message struct {
	// Type is the message type (error, warning, info).
	Type MessageType `json:"type"`

	// Code is a machine-readable error code.
	Code string `json:"code,omitempty"`

	// Content is the human-readable message.
	Content string `json:"content"`

	// ContentType indicates the format of the content (plain, markdown).
	ContentType ContentType `json:"content_type,omitempty"`

	// Severity indicates who can resolve this issue.
	Severity Severity `json:"severity,omitempty"`

	// Path is the RFC 9535 JSONPath to the component this message refers to.
	Path string `json:"path,omitempty"`
}

// TotalResponse represents a total amount breakdown.
type TotalResponse struct {
	// Type is the categorization of this total.
	Type TotalType `json:"type"`

	// Amount is the monetary value in minor (cents) currency units.
	Amount int `json:"amount"`

	// DisplayText is the text to display against the amount.
	DisplayText string `json:"display_text,omitempty"`
}

// TotalCreateRequest represents a total in a create request.
type TotalCreateRequest struct {
	// Type is the categorization of this total.
	Type TotalType `json:"type"`

	// Amount is the monetary value in minor (cents) currency units.
	Amount int `json:"amount"`

	// DisplayText is the text to display against the amount.
	DisplayText string `json:"display_text,omitempty"`
}

// PostalAddress represents a postal/mailing address using Schema.org naming conventions.
type PostalAddress struct {
	// StreetAddress is the street address.
	StreetAddress string `json:"street_address,omitempty"`

	// ExtendedAddress is an address extension such as apartment number or C/O.
	ExtendedAddress string `json:"extended_address,omitempty"`

	// AddressLocality is the city/town (e.g., Mountain View).
	AddressLocality string `json:"address_locality,omitempty"`

	// AddressRegion is the state/province/region (e.g., California).
	AddressRegion string `json:"address_region,omitempty"`

	// AddressCountry is the country code (ISO 3166-1 alpha-2 recommended, e.g., "US").
	AddressCountry string `json:"address_country,omitempty"`

	// PostalCode is the ZIP/postal code (e.g., 94043).
	PostalCode string `json:"postal_code,omitempty"`

	// FirstName is the first name of the contact.
	FirstName string `json:"first_name,omitempty"`

	// LastName is the last name of the contact.
	LastName string `json:"last_name,omitempty"`

	// FullName is the full name of the contact (first_name/last_name take precedence if present).
	FullName string `json:"full_name,omitempty"`

	// PhoneNumber is a contact phone number.
	PhoneNumber string `json:"phone_number,omitempty"`
}

// ItemResponse represents an item in a line item response.
type ItemResponse struct {
	// ID is a unique identifier for the item.
	ID string `json:"id"`

	// Title is the product title.
	Title string `json:"title"`

	// Price is the unit price in minor (cents) currency units.
	Price int `json:"price"`

	// ImageURL is a URL to an item image.
	ImageURL string `json:"image_url,omitempty"`
}

// ItemCreateRequest represents an item in a create request.
// The platform sends just the ID; the business returns full item details.
type ItemCreateRequest struct {
	// ID is the unique identifier for the item.
	ID string `json:"id"`
}

// ItemUpdateRequest represents an item in an update request.
type ItemUpdateRequest struct {
	// ID is the unique identifier for the item.
	ID string `json:"id"`
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
	// Item contains the item details.
	Item ItemCreateRequest `json:"item"`

	// Quantity is the number of items.
	Quantity int `json:"quantity"`
}

// LineItemUpdateRequest represents a line item update.
type LineItemUpdateRequest struct {
	// ID is the line item identifier.
	ID string `json:"id,omitempty"`

	// Item contains updated item details.
	Item ItemUpdateRequest `json:"item"`

	// Quantity is the updated quantity.
	Quantity int `json:"quantity"`

	// ParentID is the parent line item identifier for nested structures.
	ParentID string `json:"parent_id,omitempty"`
}

// Buyer represents information about the buyer.
type Buyer struct {
	// FirstName is the buyer's first name.
	FirstName string `json:"first_name,omitempty"`

	// LastName is the buyer's last name.
	LastName string `json:"last_name,omitempty"`

	// FullName is the buyer's full name (first_name/last_name take precedence if present).
	FullName string `json:"full_name,omitempty"`

	// Email is the buyer's email address.
	Email string `json:"email,omitempty"`

	// PhoneNumber is the buyer's phone number (E.164 format).
	PhoneNumber string `json:"phone_number,omitempty"`
}

// OrderConfirmation contains details about an order created for a checkout.
type OrderConfirmation struct {
	// ID is the unique order identifier.
	ID string `json:"id"`

	// PermalinkURL is a permalink to access the order on the merchant site.
	PermalinkURL string `json:"permalink_url"`
}

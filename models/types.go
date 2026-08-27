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

	// SeverityUnrecoverable indicates no valid resource exists to act on.
	// Retry with a new resource or inputs.
	SeverityUnrecoverable Severity = "unrecoverable"
)

// ErrorCode represents standard error codes for UCP messages.
// Standard errors have defined semantics; freeform codes are also permitted.
type ErrorCode string

const (
	// ErrorCodeOutOfStock indicates the item is out of stock.
	ErrorCodeOutOfStock ErrorCode = "out_of_stock"

	// ErrorCodeItemUnavailable indicates the item is not available.
	ErrorCodeItemUnavailable ErrorCode = "item_unavailable"

	// ErrorCodeAddressUndeliverable indicates the address cannot be delivered to.
	ErrorCodeAddressUndeliverable ErrorCode = "address_undeliverable"

	// ErrorCodePaymentFailed indicates payment processing failed.
	ErrorCodePaymentFailed ErrorCode = "payment_failed"
)

// InfoCode represents standard info codes for informational messages.
// Standard codes are defined in capability specifications; freeform codes are permitted.
type InfoCode string

const (
	// InfoCodeIdentityOptional indicates buyer identity is optional.
	InfoCodeIdentityOptional InfoCode = "identity_optional"

	// InfoCodeSignal indicates a signal-based informational message.
	InfoCodeSignal InfoCode = "signal"

	// InfoCodeFreeShipping indicates free shipping is available.
	InfoCodeFreeShipping InfoCode = "free_shipping"

	// InfoCodeNotFound indicates the requested resource was not found.
	InfoCodeNotFound InfoCode = "not_found"
)

// WarningCode represents standard warning codes for warning messages.
// Standard codes are defined in capability specifications; freeform codes are permitted.
type WarningCode string

const (
	// WarningCodeFinalSale indicates the item is final sale (no returns).
	WarningCodeFinalSale WarningCode = "final_sale"

	// WarningCodeProp65 indicates a California Proposition 65 warning.
	WarningCodeProp65 WarningCode = "prop65"

	// WarningCodeFulfillmentChanged indicates fulfillment terms have changed.
	WarningCodeFulfillmentChanged WarningCode = "fulfillment_changed"

	// WarningCodeAgeRestricted indicates the item is age-restricted.
	WarningCodeAgeRestricted WarningCode = "age_restricted"
)

// AvailablePaymentInstrument represents an instrument type available from a payment handler.
type AvailablePaymentInstrument struct {
	// Type is the instrument type identifier (e.g., "card", "gift_card").
	// References an instrument schema's type constant.
	Type string `json:"type"`

	// Constraints contains optional constraints on this instrument type.
	// Structure depends on instrument type and active capabilities.
	Constraints map[string]interface{} `json:"constraints,omitempty"`
}

// ContentType represents the content format.
type ContentType string

const (
	// ContentTypePlain indicates plain text content.
	ContentTypePlain ContentType = "plain"

	// ContentTypeMarkdown indicates markdown content.
	ContentTypeMarkdown ContentType = "markdown"
)

// Context represents buyer signals for relevance and localization.
// Context values are provisional hints - businesses SHOULD use them when
// authoritative data (e.g., address) is absent, and MAY ignore unsupported
// values without returning errors.
type Context struct {
	// AddressCountry is the country hint. Recommended to be in 2-letter ISO 3166-1
	// alpha-2 format (e.g., "US").
	AddressCountry string `json:"address_country,omitempty"`

	// AddressRegion is the region/state hint (e.g., "CA" for California).
	AddressRegion string `json:"address_region,omitempty"`

	// PostalCode is the postal/zip code hint (e.g., "94043").
	PostalCode string `json:"postal_code,omitempty"`

	// Intent describes the buyer's purpose (e.g., "looking for a gift under $50").
	// Informs relevance, recommendations, and personalization.
	Intent string `json:"intent,omitempty"`

	// Language is the preferred language using IETF BCP 47 language tags
	// (e.g., "en", "fr-CA", "zh-Hans"). For REST, equivalent to Accept-Language header.
	// When provided, overrides Accept-Language.
	Language string `json:"language,omitempty"`

	// Currency is the preferred currency in ISO 4217 format (e.g., "EUR", "USD").
	// Also serves as the denomination for price filter values.
	Currency string `json:"currency,omitempty"`

	// Eligibility contains buyer claims about eligible benefits (e.g., loyalty
	// membership, payment instrument perks). Values MUST use reverse-domain naming
	// (e.g., "com.example.loyalty_gold", "org.school.student") and MUST be non-identifying.
	Eligibility []string `json:"eligibility,omitempty"`
}

// Attribution represents platform-emitted referral and conversion-event context.
// Includes campaign identifiers, click IDs, source/medium markers, etc. -
// the same parameters platforms communicate via URL query parameters in
// browser-based flows. Values are URL-style parameter strings.
type Attribution map[string]string

// Amount represents a non-negative monetary amount in the currency's minor unit
// as defined by ISO 4217 (e.g., cents for USD, yen for JPY).
type Amount int

// SignedAmount represents a monetary amount that may be negative.
// Discounts are negative, charges are positive.
// Uses the currency's minor unit as defined by ISO 4217.
type SignedAmount int

// Signals represents environment data provided by the platform for authorization
// and abuse prevention. Values MUST NOT be buyer-asserted claims. All signal keys
// use reverse-domain naming (e.g., dev.ucp.buyer_ip).
type Signals struct {
	// BuyerIP is the client's IP address (IPv4 or IPv6).
	BuyerIP string `json:"dev.ucp.buyer_ip,omitempty"`

	// UserAgent is the client's HTTP User-Agent header or equivalent.
	UserAgent string `json:"dev.ucp.user_agent,omitempty"`

	// Additional contains any additional platform-specific signals.
	Additional map[string]interface{} `json:"-"`
}

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

// Policy is a durable business rule about the items in a response (return/refund
// terms, warranty, and the like) at the time of purchase. Every policy carries a
// Type (an open reverse-DNS vocabulary) and a Description so a platform can
// present it without understanding type-specific fields. Policies are reference
// data; the obligation to display a term to the buyer is carried by a messages[]
// warning whose code equals the policy Type. Added by the 2026-08-25 UCP release.
type Policy struct {
	// Type is the policy type discriminator, an open reverse-DNS vocabulary.
	// Well-known values: "dev.ucp.shopping.policy.return",
	// "dev.ucp.shopping.policy.warranty". Platforms MUST tolerate unknown values.
	Type string `json:"type"`

	// Description is a human-readable policy summary in one or more formats.
	// Required so a platform can present it without understanding type-specific
	// fields.
	Description Description `json:"description"`

	// AppliesTo lists RFC 9535 JSONPath expressions identifying the nodes this
	// policy applies to, relative to the embedding response root. When omitted,
	// the policy applies to the entire response.
	AppliesTo []string `json:"applies_to,omitempty"`

	// URL is an optional link to the full policy document.
	URL string `json:"url,omitempty"`
}

// MessagePresentation indicates how a warning should be rendered.
type MessagePresentation string

const (
	// MessagePresentationNotice indicates platform MUST display, MAY dismiss.
	MessagePresentationNotice MessagePresentation = "notice"

	// MessagePresentationDisclosure indicates platform MUST display in proximity
	// to the referenced component, MUST NOT hide or auto-dismiss.
	MessagePresentationDisclosure MessagePresentation = "disclosure"
)

// Message represents an error, warning, or info message.
type Message struct {
	// Type is the message type (error, warning, info).
	Type MessageType `json:"type"`

	// Code is a machine-readable error/warning code.
	Code string `json:"code,omitempty"`

	// Content is the human-readable message.
	Content string `json:"content"`

	// ContentType indicates the format of the content (plain, markdown).
	ContentType ContentType `json:"content_type,omitempty"`

	// Severity indicates who can resolve this issue (for error messages).
	Severity Severity `json:"severity,omitempty"`

	// Path is the RFC 9535 JSONPath to the component this message refers to.
	Path string `json:"path,omitempty"`

	// Presentation indicates how a warning should be rendered (for warning messages).
	// Defaults to "notice".
	Presentation MessagePresentation `json:"presentation,omitempty"`

	// ImageURL is a URL to a required visual element for warnings (e.g., warning symbol).
	ImageURL string `json:"image_url,omitempty"`

	// URL is a reference URL for more information (e.g., regulatory site, policy page).
	URL string `json:"url,omitempty"`
}

// TotalResponse represents a total amount breakdown.
type TotalResponse struct {
	// Type is the categorization of this total.
	Type TotalType `json:"type"`

	// Amount is the monetary value in minor (cents) currency units.
	// May be negative for discount types (discount, items_discount).
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

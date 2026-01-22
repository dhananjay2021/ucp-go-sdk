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

// PaymentInstrumentType represents the type of payment instrument.
type PaymentInstrumentType string

const (
	// PaymentInstrumentTypeCard indicates a card payment instrument.
	PaymentInstrumentTypeCard PaymentInstrumentType = "card"
)

// PaymentHandlerResponse represents a payment handler in a response.
type PaymentHandlerResponse struct {
	// ID is the unique identifier for this handler instance.
	ID string `json:"id"`

	// Name is the specification name using reverse-DNS format (e.g., dev.ucp.delegate_payment).
	Name string `json:"name"`

	// Version is the handler version in YYYY-MM-DD format.
	Version string `json:"version"`

	// Spec is a URL to the technical specification for this handler.
	Spec string `json:"spec"`

	// ConfigSchema is a URL to the JSON Schema for validating the config object.
	ConfigSchema string `json:"config_schema"`

	// InstrumentSchemas is a list of URLs to schemas for validating instrument objects.
	InstrumentSchemas []string `json:"instrument_schemas"`

	// Config contains handler-specific configuration.
	Config map[string]interface{} `json:"config"`
}

// PaymentIdentity represents payment identity information.
type PaymentIdentity struct {
	// AccessToken is the OAuth access token.
	AccessToken string `json:"access_token"`
}

// CardCredential represents card payment credentials.
// CRITICAL: Both parties handling CardCredential MUST be PCI DSS compliant.
type CardCredential struct {
	// Type is always "card" for card credentials.
	Type PaymentInstrumentType `json:"type"`

	// CardNumberType indicates the type of card number (fpan, network_token, dpan).
	CardNumberType CardNumberType `json:"card_number_type"`

	// Number is the card number.
	Number string `json:"number,omitempty"`

	// ExpiryMonth is the card expiration month (1-12).
	ExpiryMonth int `json:"expiry_month,omitempty"`

	// ExpiryYear is the card expiration year.
	ExpiryYear int `json:"expiry_year,omitempty"`

	// Name is the cardholder name.
	Name string `json:"name,omitempty"`

	// CVC is the card verification code.
	CVC string `json:"cvc,omitempty"`

	// Cryptogram is for network tokens.
	Cryptogram string `json:"cryptogram,omitempty"`

	// ECIValue is the electronic commerce indicator for network tokens.
	ECIValue string `json:"eci_value,omitempty"`
}

// PaymentCredential represents a payment credential.
// Currently only card credentials are supported.
type PaymentCredential struct {
	// Type indicates the credential type.
	Type string `json:"type"`

	// CardNumberType indicates the type of card number.
	CardNumberType CardNumberType `json:"card_number_type,omitempty"`

	// Number is the card number.
	Number string `json:"number,omitempty"`

	// ExpiryMonth is the card expiration month.
	ExpiryMonth int `json:"expiry_month,omitempty"`

	// ExpiryYear is the card expiration year.
	ExpiryYear int `json:"expiry_year,omitempty"`

	// Name is the cardholder name.
	Name string `json:"name,omitempty"`

	// CVC is the card verification code.
	CVC string `json:"cvc,omitempty"`

	// Cryptogram is for network tokens.
	Cryptogram string `json:"cryptogram,omitempty"`

	// ECIValue is the electronic commerce indicator.
	ECIValue string `json:"eci_value,omitempty"`
}

// PaymentInstrumentBase represents the base fields for any payment instrument.
type PaymentInstrumentBase struct {
	// ID is a unique identifier for this instrument instance.
	ID string `json:"id"`

	// HandlerID is the handler that produced this instrument.
	HandlerID string `json:"handler_id"`

	// Type is the instrument type (e.g., "card").
	Type PaymentInstrumentType `json:"type"`

	// BillingAddress is the billing address for this payment method.
	BillingAddress *PostalAddress `json:"billing_address,omitempty"`

	// Credential contains payment credential data.
	Credential *PaymentCredential `json:"credential,omitempty"`
}

// CardPaymentInstrument represents a card payment instrument.
type CardPaymentInstrument struct {
	PaymentInstrumentBase

	// Brand is the card brand/network (e.g., visa, mastercard, amex).
	Brand string `json:"brand"`

	// LastDigits is the last 4 digits of the card number.
	LastDigits string `json:"last_digits"`

	// ExpiryMonth is the card expiration month.
	ExpiryMonth int `json:"expiry_month,omitempty"`

	// ExpiryYear is the card expiration year.
	ExpiryYear int `json:"expiry_year,omitempty"`

	// RichTextDescription is an optional rich text description of the card.
	RichTextDescription string `json:"rich_text_description,omitempty"`

	// RichCardArt is an optional URI to card art.
	RichCardArt string `json:"rich_card_art,omitempty"`
}

// PaymentInstrument represents a payment instrument (currently only cards supported).
// For JSON marshaling, this uses the card payment instrument structure.
type PaymentInstrument struct {
	// ID is a unique identifier for this instrument instance.
	ID string `json:"id"`

	// HandlerID is the handler that produced this instrument.
	HandlerID string `json:"handler_id"`

	// Type is the instrument type (e.g., "card").
	Type PaymentInstrumentType `json:"type"`

	// BillingAddress is the billing address for this payment method.
	BillingAddress *PostalAddress `json:"billing_address,omitempty"`

	// Credential contains payment credential data.
	Credential *PaymentCredential `json:"credential,omitempty"`

	// Brand is the card brand/network (for card instruments).
	Brand string `json:"brand,omitempty"`

	// LastDigits is the last 4 digits of the card number (for card instruments).
	LastDigits string `json:"last_digits,omitempty"`

	// ExpiryMonth is the card expiration month (for card instruments).
	ExpiryMonth int `json:"expiry_month,omitempty"`

	// ExpiryYear is the card expiration year (for card instruments).
	ExpiryYear int `json:"expiry_year,omitempty"`

	// RichTextDescription is an optional rich text description.
	RichTextDescription string `json:"rich_text_description,omitempty"`

	// RichCardArt is an optional URI to card art.
	RichCardArt string `json:"rich_card_art,omitempty"`
}

// TokenCredentialCreateRequest represents a request to create a token credential.
type TokenCredentialCreateRequest struct {
	// Token is the credential token.
	Token string `json:"token"`

	// Type indicates the token type.
	Type string `json:"type"`
}

// TokenCredentialUpdateRequest represents a request to update a token credential.
type TokenCredentialUpdateRequest struct {
	// Token is the credential token.
	Token string `json:"token"`

	// Type indicates the token type.
	Type string `json:"type"`
}

// TokenCredentialResponse represents a tokenized credential response.
type TokenCredentialResponse struct {
	// Type indicates the token type.
	Type string `json:"type"`
}

// Binding represents the binding context for tokenization.
type Binding struct {
	// CheckoutID is the checkout session ID.
	CheckoutID string `json:"checkout_id"`

	// Identity contains optional identity information.
	Identity *PaymentIdentity `json:"identity,omitempty"`
}

// PaymentAccountInfo represents payment account information.
type PaymentAccountInfo struct {
	// PaymentAccountReference is a reference to the payment account.
	PaymentAccountReference string `json:"payment_account_reference,omitempty"`
}

// PaymentCreateRequest represents payment in a checkout create request.
type PaymentCreateRequest struct {
	// Instruments is the list of payment instruments.
	Instruments []PaymentInstrument `json:"instruments,omitempty"`

	// SelectedInstrumentID is the ID of the selected payment instrument.
	SelectedInstrumentID string `json:"selected_instrument_id,omitempty"`
}

// PaymentUpdateRequest represents payment in a checkout update request.
type PaymentUpdateRequest struct {
	// Instruments is the list of payment instruments.
	Instruments []PaymentInstrument `json:"instruments,omitempty"`

	// SelectedInstrumentID is the ID of the selected payment instrument.
	SelectedInstrumentID string `json:"selected_instrument_id,omitempty"`
}

// PaymentResponse represents payment information in a checkout response.
type PaymentResponse struct {
	// Handlers lists available payment handlers.
	Handlers []PaymentHandlerResponse `json:"handlers"`

	// Instruments is the list of payment instruments available.
	Instruments []PaymentInstrument `json:"instruments,omitempty"`

	// SelectedInstrumentID is the ID of the currently selected payment instrument.
	SelectedInstrumentID string `json:"selected_instrument_id,omitempty"`
}

// PaymentData represents payment data for complete requests.
type PaymentData struct {
	// PaymentData contains the payment instrument data.
	PaymentData PaymentInstrument `json:"payment_data"`
}

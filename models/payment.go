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

	// AvailableInstruments lists instrument types this handler supports, with optional constraints.
	// When absent, every instrument should be considered available.
	AvailableInstruments []AvailablePaymentInstrument `json:"available_instruments,omitempty"`

	// Config contains handler-specific configuration.
	Config map[string]interface{} `json:"config"`
}

// PaymentIdentity represents payment identity information.
type PaymentIdentity struct {
	// AccessToken is the OAuth access token.
	AccessToken string `json:"access_token"`
}

// CredentialType identifies a payment credential's on-the-wire shape.
type CredentialType string

const (
	// CredentialTypePAN indicates a PAN credential carrying a funding primary
	// account number (FPAN).
	CredentialTypePAN CredentialType = "pan"

	// CredentialTypeNetworkToken indicates a network token credential verified
	// with a transaction cryptogram.
	CredentialTypeNetworkToken CredentialType = "network_token"

	// CredentialTypeCard indicates the deprecated combined card credential.
	//
	// Deprecated: use CredentialTypePAN or CredentialTypeNetworkToken.
	CredentialTypeCard CredentialType = "card"
)

// PanCredential is a card credential carrying a funding primary account number
// (FPAN). Credential selection follows the shape of the value on the wire rather
// than its provenance: a network token surfaced in PAN form and verified with a
// CVC is carried here, while a token verified with a discrete cryptogram uses
// NetworkTokenCredential. This credential type MUST NOT be used for checkout,
// only with payment handlers that tokenize or encrypt credentials.
//
// CRITICAL: Both parties handling a PAN credential (sender and receiver) MUST be
// PCI DSS compliant. Transmission MUST use HTTPS/TLS with strong cipher suites.
//
// Introduced by UCP spec change #424 (split PAN and Network Token credentials).
type PanCredential struct {
	// Type is always "pan" for PAN credentials.
	Type CredentialType `json:"type"`

	// Number is the funding primary account number (FPAN).
	Number string `json:"number"`

	// ExpiryMonth is the card expiration month (1-12).
	ExpiryMonth int `json:"expiry_month,omitempty"`

	// ExpiryYear is the card expiration year.
	ExpiryYear int `json:"expiry_year,omitempty"`

	// Name is the cardholder name.
	Name string `json:"name,omitempty"`

	// CVC is the card verification code (max 4 characters).
	CVC string `json:"cvc,omitempty"`
}

// NetworkTokenCredential is a card-network token credential verified with a
// transaction cryptogram. The Number field carries the network token or
// wallet-provisioned token rather than the underlying FPAN.
//
// Introduced by UCP spec change #424 (split PAN and Network Token credentials).
type NetworkTokenCredential struct {
	// Type is always "network_token" for network token credentials.
	Type CredentialType `json:"type"`

	// Number is the network token or wallet-provisioned token replacing the
	// underlying FPAN.
	Number string `json:"number"`

	// ExpiryMonth is the token expiration month (1-12).
	ExpiryMonth int `json:"expiry_month,omitempty"`

	// ExpiryYear is the token expiration year.
	ExpiryYear int `json:"expiry_year,omitempty"`

	// Name is the cardholder name.
	Name string `json:"name,omitempty"`

	// Cryptogram is the transaction cryptogram or dynamic CVC (dCVV), in the long
	// or short form expected by the card network or processor.
	Cryptogram string `json:"cryptogram"`

	// ECIValue is the Electronic Commerce Indicator / Security Level Indicator
	// associated with the transaction.
	ECIValue string `json:"eci_value,omitempty"`

	// TokenRequestorID is the payment network token requestor identifier, when
	// required by the processor or network-token program.
	TokenRequestorID string `json:"token_requestor_id,omitempty"`
}

// CardCredential represents card payment credentials.
// CRITICAL: Both parties handling CardCredential MUST be PCI DSS compliant.
//
// Deprecated: as of the 2026-08-25 UCP release the combined card credential is
// split into PanCredential (FPAN + CVC) and NetworkTokenCredential (token +
// cryptogram). Use those types for new code; CardCredential is retained for
// backward compatibility.
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

// CardDisplay represents display information for a card payment instrument.
type CardDisplay struct {
	// Brand is the card brand/network (e.g., visa, mastercard, amex).
	Brand string `json:"brand,omitempty"`

	// LastDigits is the last 4 digits of the card number.
	LastDigits string `json:"last_digits,omitempty"`

	// ExpiryMonth is the card expiration month (1-12).
	ExpiryMonth int `json:"expiry_month,omitempty"`

	// ExpiryYear is the card expiration year.
	ExpiryYear int `json:"expiry_year,omitempty"`

	// Description is an optional rich text description of the card.
	Description string `json:"description,omitempty"`

	// CardArt is an optional URI to a rich image representing the card.
	CardArt string `json:"card_art,omitempty"`
}

// CardPaymentInstrument represents a card payment instrument.
type CardPaymentInstrument struct {
	PaymentInstrumentBase

	// Display contains display information for this card.
	Display *CardDisplay `json:"display,omitempty"`

	// Legacy flat fields for backward compatibility
	// Brand is the card brand/network (e.g., visa, mastercard, amex).
	Brand string `json:"brand,omitempty"`

	// LastDigits is the last 4 digits of the card number.
	LastDigits string `json:"last_digits,omitempty"`

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

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
	// ID is the unique identifier for this handler.
	ID string `json:"id"`

	// Type is the handler type (e.g., "tokenization").
	Type string `json:"type"`

	// Name is a display name for the handler.
	Name string `json:"name,omitempty"`

	// Description provides details about the handler.
	Description string `json:"description,omitempty"`

	// TokenizationURL is the URL for tokenization requests.
	TokenizationURL string `json:"tokenization_url,omitempty"`

	// SupportedInstruments lists supported payment instrument schemas.
	SupportedInstruments []string `json:"supported_instruments,omitempty"`

	// Config contains handler-specific configuration.
	Config map[string]interface{} `json:"config,omitempty"`
}

// PaymentHandlerCreateRequest represents a payment handler in a create request.
type PaymentHandlerCreateRequest struct {
	// ID is the handler identifier.
	ID string `json:"id"`

	// Config contains handler-specific configuration.
	Config map[string]interface{} `json:"config,omitempty"`
}

// PaymentIdentity represents payment identity information.
type PaymentIdentity struct {
	// AccessToken is the OAuth access token.
	AccessToken string `json:"access_token,omitempty"`

	// TokenType is the token type (e.g., "Bearer").
	TokenType string `json:"token_type,omitempty"`

	// ExpiresIn is the token expiration time in seconds.
	ExpiresIn int `json:"expires_in,omitempty"`
}

// CardCredential represents card payment credentials.
type CardCredential struct {
	// Type is always "card" for card credentials.
	Type PaymentInstrumentType `json:"type"`

	// CardNumberType indicates the type of card number.
	CardNumberType CardNumberType `json:"card_number_type"`

	// Number is the card number.
	Number string `json:"number"`

	// ExpiryMonth is the card expiration month (1-12).
	ExpiryMonth int `json:"expiry_month"`

	// ExpiryYear is the card expiration year.
	ExpiryYear int `json:"expiry_year"`

	// CVC is the card verification code.
	CVC string `json:"cvc,omitempty"`

	// Name is the cardholder name.
	Name string `json:"name,omitempty"`

	// Cryptogram is for network tokens.
	Cryptogram string `json:"cryptogram,omitempty"`

	// ECIValue is the electronic commerce indicator for network tokens.
	ECIValue string `json:"eci_value,omitempty"`
}

// PaymentCredential represents a payment credential (can be different types).
type PaymentCredential struct {
	// Card is set when Type is "card".
	Card *CardCredential `json:"-"`

	// Raw contains the raw credential data.
	Raw map[string]interface{} `json:"-"`
}

// CardPaymentInstrument represents a card payment instrument.
type CardPaymentInstrument struct {
	// Type is always "card".
	Type PaymentInstrumentType `json:"type"`

	// Last4 is the last 4 digits of the card number.
	Last4 string `json:"last4,omitempty"`

	// Brand is the card brand (e.g., "visa", "mastercard").
	Brand string `json:"brand,omitempty"`

	// ExpiryMonth is the card expiration month.
	ExpiryMonth int `json:"expiry_month,omitempty"`

	// ExpiryYear is the card expiration year.
	ExpiryYear int `json:"expiry_year,omitempty"`

	// Name is the cardholder name.
	Name string `json:"name,omitempty"`
}

// PaymentInstrument represents a payment instrument (union type).
type PaymentInstrument struct {
	// Card is set for card instruments.
	Card *CardPaymentInstrument `json:"-"`

	// Raw contains the raw instrument data.
	Raw map[string]interface{} `json:"-"`
}

// TokenCredentialResponse represents a tokenized credential response.
type TokenCredentialResponse struct {
	// Token is the tokenized credential.
	Token string `json:"token"`

	// ExpiresAt is when the token expires.
	ExpiresAt string `json:"expires_at,omitempty"`

	// Instrument contains instrument details.
	Instrument *PaymentInstrument `json:"instrument,omitempty"`
}

// TokenCredentialCreateRequest represents a request to create a token credential.
type TokenCredentialCreateRequest struct {
	// Credential contains the payment credential to tokenize.
	Credential PaymentCredential `json:"credential"`

	// Binding contains the binding context.
	Binding Binding `json:"binding"`
}

// Binding represents the binding context for tokenization.
type Binding struct {
	// CheckoutID is the checkout session ID.
	CheckoutID string `json:"checkout_id"`

	// Identity contains optional identity information.
	Identity *PaymentIdentity `json:"identity,omitempty"`
}

// PaymentCreateRequest represents a request to create payment.
type PaymentCreateRequest struct {
	// HandlerID is the payment handler to use.
	HandlerID string `json:"handler_id"`

	// Token is the tokenized credential.
	Token string `json:"token,omitempty"`

	// Instrument contains optional instrument details.
	Instrument *PaymentInstrument `json:"instrument,omitempty"`

	// BillingAddress is the billing address.
	BillingAddress *PostalAddress `json:"billing_address,omitempty"`
}

// PaymentUpdateRequest represents a request to update payment.
type PaymentUpdateRequest struct {
	// HandlerID is the payment handler to use.
	HandlerID string `json:"handler_id,omitempty"`

	// Token is the tokenized credential.
	Token string `json:"token,omitempty"`

	// Instrument contains updated instrument details.
	Instrument *PaymentInstrument `json:"instrument,omitempty"`

	// BillingAddress is the billing address.
	BillingAddress *PostalAddress `json:"billing_address,omitempty"`
}

// PaymentResponse represents payment information in a checkout response.
type PaymentResponse struct {
	// Status is the payment status.
	Status string `json:"status,omitempty"`

	// HandlerID is the selected payment handler.
	HandlerID string `json:"handler_id,omitempty"`

	// Instrument contains instrument details (masked).
	Instrument *PaymentInstrument `json:"instrument,omitempty"`

	// Handlers lists available payment handlers.
	Handlers []PaymentHandlerResponse `json:"handlers,omitempty"`

	// Total is the payment amount.
	Total string `json:"total,omitempty"`

	// Currency is the payment currency.
	Currency string `json:"currency,omitempty"`
}

// PaymentData represents payment data in various contexts.
type PaymentData struct {
	// HandlerID is the payment handler ID.
	HandlerID string `json:"handler_id,omitempty"`

	// Token is a payment token.
	Token string `json:"token,omitempty"`

	// Instrument contains instrument details.
	Instrument *PaymentInstrument `json:"instrument,omitempty"`

	// BillingAddress is the billing address.
	BillingAddress *PostalAddress `json:"billing_address,omitempty"`
}

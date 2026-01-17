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

// AccountInfo represents buyer account information.
type AccountInfo struct {
	// ID is the account identifier.
	ID string `json:"id,omitempty"`

	// Email is the account email.
	Email string `json:"email,omitempty"`

	// Phone is the account phone number.
	Phone string `json:"phone,omitempty"`

	// FirstName is the account holder's first name.
	FirstName string `json:"first_name,omitempty"`

	// LastName is the account holder's last name.
	LastName string `json:"last_name,omitempty"`

	// Linked indicates if the account is linked.
	Linked bool `json:"linked,omitempty"`
}

// BuyerConsentCreateRequest represents buyer consent in a create request.
type BuyerConsentCreateRequest struct {
	// Email is the buyer's email.
	Email string `json:"email,omitempty"`

	// Phone is the buyer's phone number.
	Phone string `json:"phone,omitempty"`

	// FirstName is the buyer's first name.
	FirstName string `json:"first_name,omitempty"`

	// LastName is the buyer's last name.
	LastName string `json:"last_name,omitempty"`

	// BillingAddress is the buyer's billing address.
	BillingAddress *PostalAddress `json:"billing_address,omitempty"`

	// MarketingOptIn indicates consent to marketing.
	MarketingOptIn *bool `json:"marketing_opt_in,omitempty"`

	// TermsAccepted indicates acceptance of terms.
	TermsAccepted *bool `json:"terms_accepted,omitempty"`

	// AccountInfo contains linked account info.
	AccountInfo *AccountInfo `json:"account_info,omitempty"`
}

// BuyerConsentUpdateRequest represents buyer consent in an update request.
type BuyerConsentUpdateRequest struct {
	// Email is the buyer's email.
	Email string `json:"email,omitempty"`

	// Phone is the buyer's phone number.
	Phone string `json:"phone,omitempty"`

	// FirstName is the buyer's first name.
	FirstName string `json:"first_name,omitempty"`

	// LastName is the buyer's last name.
	LastName string `json:"last_name,omitempty"`

	// BillingAddress is the buyer's billing address.
	BillingAddress *PostalAddress `json:"billing_address,omitempty"`

	// MarketingOptIn indicates consent to marketing.
	MarketingOptIn *bool `json:"marketing_opt_in,omitempty"`

	// TermsAccepted indicates acceptance of terms.
	TermsAccepted *bool `json:"terms_accepted,omitempty"`

	// AccountInfo contains linked account info.
	AccountInfo *AccountInfo `json:"account_info,omitempty"`
}

// BuyerConsentResponse represents buyer consent in a response.
type BuyerConsentResponse struct {
	// Email is the buyer's email.
	Email string `json:"email,omitempty"`

	// Phone is the buyer's phone number.
	Phone string `json:"phone,omitempty"`

	// FirstName is the buyer's first name.
	FirstName string `json:"first_name,omitempty"`

	// LastName is the buyer's last name.
	LastName string `json:"last_name,omitempty"`

	// BillingAddress is the buyer's billing address.
	BillingAddress *PostalAddress `json:"billing_address,omitempty"`

	// MarketingOptIn indicates consent to marketing.
	MarketingOptIn *bool `json:"marketing_opt_in,omitempty"`

	// TermsAccepted indicates acceptance of terms.
	TermsAccepted *bool `json:"terms_accepted,omitempty"`

	// AccountInfo contains linked account info.
	AccountInfo *AccountInfo `json:"account_info,omitempty"`

	// RequiredFields lists fields that are required.
	RequiredFields []string `json:"required_fields,omitempty"`
}

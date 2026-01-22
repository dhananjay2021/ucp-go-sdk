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

// Consent represents user consent states for data processing.
type Consent struct {
	// Analytics indicates consent for analytics and performance tracking.
	Analytics *bool `json:"analytics,omitempty"`

	// Preferences indicates consent for storing user preferences.
	Preferences *bool `json:"preferences,omitempty"`

	// Marketing indicates consent for marketing communications.
	Marketing *bool `json:"marketing,omitempty"`

	// SaleOfData indicates consent for selling data to third parties (CCPA).
	SaleOfData *bool `json:"sale_of_data,omitempty"`
}

// BuyerWithConsent represents a buyer with consent tracking.
type BuyerWithConsent struct {
	Buyer

	// Consent contains consent tracking fields.
	Consent *Consent `json:"consent,omitempty"`
}

// BuyerWithConsentCreateRequest represents buyer with consent in a create request.
type BuyerWithConsentCreateRequest struct {
	// FirstName is the buyer's first name.
	FirstName string `json:"first_name,omitempty"`

	// LastName is the buyer's last name.
	LastName string `json:"last_name,omitempty"`

	// FullName is the buyer's full name.
	FullName string `json:"full_name,omitempty"`

	// Email is the buyer's email address.
	Email string `json:"email,omitempty"`

	// PhoneNumber is the buyer's phone number.
	PhoneNumber string `json:"phone_number,omitempty"`

	// Consent contains consent tracking fields.
	Consent *Consent `json:"consent,omitempty"`
}

// BuyerWithConsentUpdateRequest represents buyer with consent in an update request.
type BuyerWithConsentUpdateRequest struct {
	// FirstName is the buyer's first name.
	FirstName string `json:"first_name,omitempty"`

	// LastName is the buyer's last name.
	LastName string `json:"last_name,omitempty"`

	// FullName is the buyer's full name.
	FullName string `json:"full_name,omitempty"`

	// Email is the buyer's email address.
	Email string `json:"email,omitempty"`

	// PhoneNumber is the buyer's phone number.
	PhoneNumber string `json:"phone_number,omitempty"`

	// Consent contains consent tracking fields.
	Consent *Consent `json:"consent,omitempty"`
}

// BuyerWithConsentResponse represents buyer with consent in a response.
type BuyerWithConsentResponse struct {
	// FirstName is the buyer's first name.
	FirstName string `json:"first_name,omitempty"`

	// LastName is the buyer's last name.
	LastName string `json:"last_name,omitempty"`

	// FullName is the buyer's full name.
	FullName string `json:"full_name,omitempty"`

	// Email is the buyer's email address.
	Email string `json:"email,omitempty"`

	// PhoneNumber is the buyer's phone number.
	PhoneNumber string `json:"phone_number,omitempty"`

	// Consent contains consent tracking fields.
	Consent *Consent `json:"consent,omitempty"`
}

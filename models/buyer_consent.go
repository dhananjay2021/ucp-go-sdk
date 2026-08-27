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

// Buyer Consent extension (dev.ucp.shopping.buyer_consent).
//
// As of the 2026-08-25 UCP release, consent is modeled as per-purpose decisions
// keyed by reverse-DNS identifier rather than a fixed set of boolean flags. Each
// purpose carries the current granted state, the source of that state, a
// human-readable description, optional links, and optional segments that refine
// the decision to specific channels, vendors, or programs.

// ConsentSource identifies the party that asserted a consent decision.
type ConsentSource string

const (
	// ConsentSourceBusiness means the value reflects the business's default
	// policy.
	ConsentSourceBusiness ConsentSource = "business"

	// ConsentSourcePlatform means the value reflects an explicit buyer decision
	// captured by the platform.
	ConsentSourcePlatform ConsentSource = "platform"
)

// Well-known consent purpose identifiers. Vendors and merchants may define
// additional purposes under their own reverse-DNS namespace.
const (
	// ConsentPurposeMarketing is consent for marketing communications.
	ConsentPurposeMarketing = "dev.ucp.consent.marketing"

	// ConsentPurposeAnalytics is consent for analytics and performance tracking.
	ConsentPurposeAnalytics = "dev.ucp.consent.analytics"

	// ConsentPurposePreferences is consent for storing buyer preferences.
	ConsentPurposePreferences = "dev.ucp.consent.preferences"

	// ConsentPurposeSaleOrSharing is consent for sale or sharing of data (CCPA).
	ConsentPurposeSaleOrSharing = "dev.ucp.consent.sale_or_sharing"
)

// Well-known consent segment identifiers under ConsentPurposeMarketing.
const (
	// ConsentSegmentMarketingEmail refines marketing consent to email.
	ConsentSegmentMarketingEmail = "dev.ucp.consent.marketing.email"

	// ConsentSegmentMarketingSMS refines marketing consent to SMS.
	ConsentSegmentMarketingSMS = "dev.ucp.consent.marketing.sms"
)

// ConsentSegment is a buyer's consent decision for a specific refinement of a
// parent purpose (e.g., email marketing under the marketing purpose). It
// overrides the parent purpose's Granted value for this scope. Segments do not
// nest further.
type ConsentSegment struct {
	// Granted indicates whether consent has been granted for this segment. It
	// overrides the parent purpose's Granted value for this specific scope.
	Granted bool `json:"granted"`

	// Source identifies the party that asserted the current Granted value.
	Source ConsentSource `json:"source"`

	// Description is a human-readable description of what the buyer is consenting
	// to within this segment. Required on responses; omitted on requests.
	Description string `json:"description,omitempty"`

	// Links are optional segment-specific links (e.g., channel terms or privacy
	// disclosures). Omitted on requests.
	Links []Link `json:"links,omitempty"`
}

// ConsentPurpose is a buyer's consent decision for a purpose (e.g., marketing,
// analytics). It carries the current binary state, its source, human-readable
// context, and optional refinements scoping the decision to specific channels,
// vendors, or programs.
type ConsentPurpose struct {
	// Granted indicates whether consent has been granted for this purpose. The
	// Source field identifies who asserted this state.
	Granted bool `json:"granted"`

	// Source identifies the party that asserted the current Granted value.
	Source ConsentSource `json:"source"`

	// Description is a human-readable description of what the buyer is consenting
	// to. Required on responses; omitted on requests.
	Description string `json:"description,omitempty"`

	// Links are optional links providing context (e.g., privacy policy, terms).
	// Omitted on requests.
	Links []Link `json:"links,omitempty"`

	// Segments are optional refinements scoping this purpose to specific
	// channels, vendors, or programs, keyed by reverse-DNS identifier.
	Segments map[string]ConsentSegment `json:"segments,omitempty"`
}

// Consent is per-purpose consent, keyed by reverse-DNS purpose identifier. UCP
// defines four well-known purposes (marketing, analytics, preferences,
// sale_or_sharing); vendors and merchants may define additional purposes under
// their own reverse-DNS namespace.
type Consent map[string]ConsentPurpose

// BuyerWithConsent represents a buyer with per-purpose consent tracking.
type BuyerWithConsent struct {
	Buyer

	// Consent contains per-purpose consent decisions.
	Consent Consent `json:"consent,omitempty"`
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

	// Consent contains per-purpose consent decisions.
	Consent Consent `json:"consent,omitempty"`
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

	// Consent contains per-purpose consent decisions.
	Consent Consent `json:"consent,omitempty"`
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

	// Consent contains per-purpose consent decisions.
	Consent Consent `json:"consent,omitempty"`
}

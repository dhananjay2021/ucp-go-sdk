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

// LoyaltyRewardAmount is a non-negative integer amount denominated in the minor
// unit of the associated reward currency. The reward currency's DecimalPlaces
// defines the minor-to-major ratio and defaults to 0 when omitted.
type LoyaltyRewardAmount int

// LoyaltyRewardCurrency describes the currency of a loyalty reward: a unit of
// value that customers can accumulate through various commercial activities.
type LoyaltyRewardCurrency struct {
	// Name is the human-readable name of the currency (e.g., "LoyaltyStars").
	Name string `json:"name"`

	// Code is the business-specific representation of the currency (e.g., "LST").
	Code string `json:"code"`

	// DecimalPlaces is the position of a digit to the right of a decimal point.
	// Applies to all reward amount fields and defaults to 0 when omitted.
	DecimalPlaces *int `json:"decimal_places,omitempty"`
}

// LoyaltyEarningBreakdown is a single breakdown rule contributing to reward earnings.
type LoyaltyEarningBreakdown struct {
	// ID is the unique rewards breakdown rule identifier.
	ID string `json:"id"`

	// Amount is the rewards earned from this rule.
	Amount LoyaltyRewardAmount `json:"amount"`

	// Description is a display-ready, human-readable rationale for the specific
	// rewards (e.g., "2x on footwear").
	Description string `json:"description"`

	// BenefitID is the optional id of the LoyaltyMembershipTierBenefit that
	// produced this rewards rule. Resolves against LoyaltyMembershipTierBenefit.ID
	// within the same parent loyalty membership.
	BenefitID *string `json:"benefit_id,omitempty"`
}

// LoyaltyEarningForecast is a preview of rewards to be earned from the current transaction.
type LoyaltyEarningForecast struct {
	// Amount is the total rewards to be earned if the transaction completes.
	Amount LoyaltyRewardAmount `json:"amount"`

	// Breakdown is the list of earning breakdowns contributing to the total.
	Breakdown []LoyaltyEarningBreakdown `json:"breakdown,omitempty"`
}

// LoyaltyMembershipReward is a quantifiable reward type and optional earning
// forecast for the current transaction.
type LoyaltyMembershipReward struct {
	// Currency is a unit of value that customers can accumulate through various
	// commercial activities.
	Currency LoyaltyRewardCurrency `json:"currency"`

	// EarningForecast is a preview of rewards to be earned from the current transaction.
	EarningForecast *LoyaltyEarningForecast `json:"earning_forecast,omitempty"`
}

// LoyaltyMembershipTierBenefit is a benefit associated with a membership tier.
type LoyaltyMembershipTierBenefit struct {
	// ID is the unique identifier for the tier benefit.
	ID string `json:"id"`

	// Description is a display-ready, human-readable explanation of this benefit
	// (e.g., "Early access to sales").
	Description string `json:"description"`
}

// LoyaltyMembershipTier is a specific achievement rank or status milestone that
// unlocks escalating value as a member progresses through activity or spend.
type LoyaltyMembershipTier struct {
	// ID is the unique identifier for the membership tier.
	ID string `json:"id"`

	// Name is the human-readable name of the tier (e.g., "Platinum").
	Name string `json:"name"`

	// Benefits is the list of benefits associated with this tier.
	Benefits []LoyaltyMembershipTierBenefit `json:"benefits,omitempty"`
}

// LoyaltyMembership is the loyalty membership the business has accepted for the
// eligibility claim represented by the parent map key. Programs that can be
// joined independently MUST be modeled as separate sibling entries under the
// Loyalty map, distinguished by reverse-domain naming (e.g.,
// "com.example.rewards" and "com.example.rewards.card").
type LoyaltyMembership struct {
	// ID is the unique loyalty membership identifier.
	ID string `json:"id"`

	// Name is the business-specific name of the loyalty membership/program.
	Name string `json:"name"`

	// DisplayID is a masked or partial version of the membership id for user
	// recognition (e.g., "****5678"). MUST NOT be set if the membership has not
	// been verified.
	DisplayID *string `json:"display_id,omitempty"`

	// Tiers is the active or display-safe tier context for this membership. Most
	// programs are single-status (one entry); programs with parallel status
	// dimensions (e.g., current and lifetime) populate one entry per active tier.
	// Omitted when no tier context has been resolved.
	Tiers []LoyaltyMembershipTier `json:"tiers,omitempty"`

	// Rewards is the reward types and earning forecasts associated with this
	// membership. Each entry encapsulates one type of reward.
	Rewards []LoyaltyMembershipReward `json:"rewards,omitempty"`

	// Provisional is true if this membership requires additional verification.
	Provisional bool `json:"provisional"`
}

// Loyalty is a key-value map whose keys represent buyer/platform asserted
// eligibility claims and whose values represent associated membership
// information. All loyalty keys MUST use reverse-domain naming to ensure
// provenance and prevent collisions when multiple extensions contribute to the
// shared namespace.
type Loyalty map[string]LoyaltyMembership

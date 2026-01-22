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

// AllocationMethod represents how a discount is allocated.
type AllocationMethod string

const (
	// AllocationMethodEach applies discount independently per item.
	AllocationMethodEach AllocationMethod = "each"

	// AllocationMethodAcross splits discount proportionally by value.
	AllocationMethodAcross AllocationMethod = "across"
)

// DiscountAllocation represents how a discount amount was allocated to a target.
type DiscountAllocation struct {
	// Path is the JSONPath to the allocation target.
	Path string `json:"path"`

	// Amount is the amount allocated in minor (cents) currency units.
	Amount int `json:"amount"`
}

// AppliedDiscount represents a discount that was successfully applied.
type AppliedDiscount struct {
	// Title is the human-readable discount name.
	Title string `json:"title"`

	// Amount is the total discount amount in minor (cents) currency units.
	Amount int `json:"amount"`

	// Code is the discount code (omitted for automatic discounts).
	Code string `json:"code,omitempty"`

	// Automatic indicates if applied automatically by merchant rules.
	Automatic bool `json:"automatic,omitempty"`

	// Method is the allocation method (each or across).
	Method AllocationMethod `json:"method,omitempty"`

	// Priority is the stacking order (lower numbers applied first).
	Priority int `json:"priority,omitempty"`

	// Allocations is the breakdown of where this discount was allocated.
	Allocations []DiscountAllocation `json:"allocations,omitempty"`
}

// DiscountsCreateRequest represents discounts in a checkout create request.
type DiscountsCreateRequest struct {
	// Codes are discount codes to apply (case-insensitive).
	Codes []string `json:"codes,omitempty"`

	// Applied contains applied discounts (for platform-side pre-application).
	Applied []AppliedDiscount `json:"applied,omitempty"`
}

// DiscountsUpdateRequest represents discounts in a checkout update request.
type DiscountsUpdateRequest struct {
	// Codes are discount codes to apply (replaces previously submitted codes).
	Codes []string `json:"codes,omitempty"`

	// Applied contains applied discounts (for platform-side pre-application).
	Applied []AppliedDiscount `json:"applied,omitempty"`
}

// DiscountsResponse represents discounts in a checkout response.
type DiscountsResponse struct {
	// Codes are the discount codes that were submitted.
	Codes []string `json:"codes,omitempty"`

	// Applied contains discounts successfully applied (code-based and automatic).
	Applied []AppliedDiscount `json:"applied,omitempty"`
}

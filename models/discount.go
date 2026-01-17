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

// DiscountType represents the type of discount.
type DiscountType string

const (
	// DiscountTypeCode is a discount code/coupon.
	DiscountTypeCode DiscountType = "code"

	// DiscountTypeAutomatic is an automatically applied discount.
	DiscountTypeAutomatic DiscountType = "automatic"

	// DiscountTypeLoyalty is a loyalty program discount.
	DiscountTypeLoyalty DiscountType = "loyalty"
)

// DiscountStatus represents the status of a discount.
type DiscountStatus string

const (
	// DiscountStatusPending indicates the discount is pending validation.
	DiscountStatusPending DiscountStatus = "pending"

	// DiscountStatusApplied indicates the discount has been applied.
	DiscountStatusApplied DiscountStatus = "applied"

	// DiscountStatusRejected indicates the discount was rejected.
	DiscountStatusRejected DiscountStatus = "rejected"

	// DiscountStatusExpired indicates the discount has expired.
	DiscountStatusExpired DiscountStatus = "expired"
)

// DiscountCreateRequest represents a request to apply a discount.
type DiscountCreateRequest struct {
	// Code is the discount code.
	Code string `json:"code,omitempty"`

	// Type is the discount type.
	Type DiscountType `json:"type,omitempty"`
}

// DiscountUpdateRequest represents a request to update a discount.
type DiscountUpdateRequest struct {
	// ID is the discount identifier.
	ID string `json:"id"`

	// Code is the updated discount code.
	Code string `json:"code,omitempty"`
}

// DiscountResponse represents a discount in a response.
type DiscountResponse struct {
	// ID is the discount identifier.
	ID string `json:"id"`

	// Code is the discount code.
	Code string `json:"code,omitempty"`

	// Type is the discount type.
	Type DiscountType `json:"type,omitempty"`

	// Status is the discount status.
	Status DiscountStatus `json:"status"`

	// Name is a display name for the discount.
	Name string `json:"name,omitempty"`

	// Description provides details about the discount.
	Description string `json:"description,omitempty"`

	// Amount is the discount amount.
	Amount string `json:"amount,omitempty"`

	// Percentage is the discount percentage.
	Percentage string `json:"percentage,omitempty"`

	// LineItemIDs are the line items this discount applies to.
	LineItemIDs []string `json:"line_item_ids,omitempty"`

	// Message contains any error or info message.
	Message *Message `json:"message,omitempty"`
}

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

import "time"

// FulfillmentMethodType represents the type of fulfillment method.
type FulfillmentMethodType string

const (
	// FulfillmentMethodTypeShipping indicates shipping delivery.
	FulfillmentMethodTypeShipping FulfillmentMethodType = "shipping"

	// FulfillmentMethodTypePickup indicates in-store pickup.
	FulfillmentMethodTypePickup FulfillmentMethodType = "pickup"
)

// ShippingDestinationRequest represents a shipping destination in a request.
type ShippingDestinationRequest struct {
	PostalAddress

	// ID is an optional identifier for this shipping destination.
	ID string `json:"id,omitempty"`
}

// ShippingDestinationResponse represents a shipping destination in a response.
type ShippingDestinationResponse struct {
	PostalAddress

	// ID is a unique identifier for this shipping destination.
	ID string `json:"id"`
}

// RetailLocationRequest represents a retail/pickup location in a request.
type RetailLocationRequest struct {
	// Name is the location name (e.g., store name).
	Name string `json:"name"`

	// Address is the physical address of the location.
	Address *PostalAddress `json:"address,omitempty"`
}

// RetailLocationResponse represents a retail/pickup location in a response.
type RetailLocationResponse struct {
	// ID is a unique identifier for this location.
	ID string `json:"id"`

	// Name is the location name (e.g., store name).
	Name string `json:"name"`

	// Address is the physical address of the location.
	Address *PostalAddress `json:"address,omitempty"`
}

// FulfillmentDestinationRequest represents a fulfillment destination in a request.
// Can be either a shipping address or a pickup location.
type FulfillmentDestinationRequest struct {
	// Shipping destination fields (embedded PostalAddress)
	PostalAddress

	// ID is an optional identifier for this destination.
	ID string `json:"id,omitempty"`

	// Address is used for pickup locations.
	Address *PostalAddress `json:"address,omitempty"`

	// Name is the location name (for pickup locations).
	Name string `json:"name,omitempty"`
}

// FulfillmentDestinationResponse represents a fulfillment destination in a response.
// Can be either a shipping address or a pickup location.
type FulfillmentDestinationResponse struct {
	// Shipping destination fields (embedded PostalAddress)
	PostalAddress

	// ID is a unique identifier for this destination.
	ID string `json:"id"`

	// Address is used for pickup locations.
	Address *PostalAddress `json:"address,omitempty"`

	// Name is the location name (for pickup locations).
	Name string `json:"name,omitempty"`
}

// FulfillmentOptionResponse represents a fulfillment option within a group.
type FulfillmentOptionResponse struct {
	// ID is a unique fulfillment option identifier.
	ID string `json:"id"`

	// Title is a short label (e.g., "Express Shipping", "Curbside Pickup").
	Title string `json:"title"`

	// Description provides complete context for buyer decision.
	Description string `json:"description,omitempty"`

	// Carrier is the carrier name (for shipping).
	Carrier string `json:"carrier,omitempty"`

	// EarliestFulfillmentTime is the earliest fulfillment date.
	EarliestFulfillmentTime *time.Time `json:"earliest_fulfillment_time,omitempty"`

	// LatestFulfillmentTime is the latest fulfillment date.
	LatestFulfillmentTime *time.Time `json:"latest_fulfillment_time,omitempty"`

	// Totals contains the fulfillment option totals breakdown.
	Totals []TotalResponse `json:"totals"`
}

// FulfillmentGroupCreateRequest represents a fulfillment group in a create request.
type FulfillmentGroupCreateRequest struct {
	// SelectedOptionID is the ID of the selected fulfillment option.
	SelectedOptionID *string `json:"selected_option_id,omitempty"`
}

// FulfillmentGroupUpdateRequest represents a fulfillment group in an update request.
type FulfillmentGroupUpdateRequest struct {
	// ID is the group identifier.
	ID string `json:"id"`

	// SelectedOptionID is the ID of the selected fulfillment option.
	SelectedOptionID *string `json:"selected_option_id,omitempty"`
}

// FulfillmentGroupResponse represents a fulfillment group in a response.
type FulfillmentGroupResponse struct {
	// ID is the group identifier.
	ID string `json:"id"`

	// LineItemIDs are the line items in this group.
	LineItemIDs []string `json:"line_item_ids"`

	// Options are the available fulfillment options for this group.
	Options []FulfillmentOptionResponse `json:"options,omitempty"`

	// SelectedOptionID is the ID of the selected fulfillment option.
	SelectedOptionID *string `json:"selected_option_id,omitempty"`
}

// FulfillmentMethodCreateRequest represents a fulfillment method in a create request.
type FulfillmentMethodCreateRequest struct {
	// Type is the fulfillment method type (shipping or pickup).
	Type FulfillmentMethodType `json:"type"`

	// LineItemIDs are the line items fulfilled via this method.
	LineItemIDs []string `json:"line_item_ids,omitempty"`

	// Destinations are the available destinations.
	Destinations []FulfillmentDestinationRequest `json:"destinations,omitempty"`

	// SelectedDestinationID is the ID of the selected destination.
	SelectedDestinationID *string `json:"selected_destination_id,omitempty"`

	// Groups are the fulfillment groups.
	Groups []FulfillmentGroupCreateRequest `json:"groups,omitempty"`
}

// FulfillmentMethodUpdateRequest represents a fulfillment method in an update request.
type FulfillmentMethodUpdateRequest struct {
	// ID is the method identifier.
	ID string `json:"id"`

	// LineItemIDs are the line items fulfilled via this method.
	LineItemIDs []string `json:"line_item_ids"`

	// Destinations are the available destinations.
	Destinations []FulfillmentDestinationRequest `json:"destinations,omitempty"`

	// SelectedDestinationID is the ID of the selected destination.
	SelectedDestinationID *string `json:"selected_destination_id,omitempty"`

	// Groups are the fulfillment groups.
	Groups []FulfillmentGroupUpdateRequest `json:"groups,omitempty"`
}

// FulfillmentMethodResponse represents a fulfillment method in a response.
type FulfillmentMethodResponse struct {
	// ID is a unique fulfillment method identifier.
	ID string `json:"id"`

	// Type is the fulfillment method type (shipping or pickup).
	Type FulfillmentMethodType `json:"type"`

	// LineItemIDs are the line items fulfilled via this method.
	LineItemIDs []string `json:"line_item_ids"`

	// Destinations are the available destinations.
	Destinations []FulfillmentDestinationResponse `json:"destinations,omitempty"`

	// SelectedDestinationID is the ID of the selected destination.
	SelectedDestinationID *string `json:"selected_destination_id,omitempty"`

	// Groups are the fulfillment groups.
	Groups []FulfillmentGroupResponse `json:"groups,omitempty"`
}

// FulfillmentAvailableMethodResponse represents inventory availability for a fulfillment method.
type FulfillmentAvailableMethodResponse struct {
	// Type is the fulfillment method type (shipping or pickup).
	Type FulfillmentMethodType `json:"type"`

	// LineItemIDs are the line items available for this method.
	LineItemIDs []string `json:"line_item_ids"`

	// FulfillableOn is "now" for immediate availability, or ISO 8601 date for future.
	FulfillableOn *string `json:"fulfillable_on,omitempty"`

	// Description provides human-readable availability info.
	Description string `json:"description,omitempty"`
}

// FulfillmentCreateRequest represents fulfillment in a checkout create request.
type FulfillmentCreateRequest struct {
	// Methods are the fulfillment methods.
	Methods []FulfillmentMethodCreateRequest `json:"methods,omitempty"`
}

// FulfillmentUpdateRequest represents fulfillment in a checkout update request.
type FulfillmentUpdateRequest struct {
	// Methods are the fulfillment methods.
	Methods []FulfillmentMethodUpdateRequest `json:"methods,omitempty"`
}

// FulfillmentResponse represents fulfillment information in a response.
type FulfillmentResponse struct {
	// Methods are the fulfillment methods.
	Methods []FulfillmentMethodResponse `json:"methods,omitempty"`

	// AvailableMethods lists inventory availability hints.
	AvailableMethods []FulfillmentAvailableMethodResponse `json:"available_methods,omitempty"`
}

// AllowsMultiDestination represents multi-destination configuration.
type AllowsMultiDestination struct {
	// Shipping indicates if multi-destination shipping is allowed.
	Shipping bool `json:"shipping,omitempty"`

	// Pickup indicates if multi-destination pickup is allowed.
	Pickup bool `json:"pickup,omitempty"`
}

// MerchantFulfillmentConfig represents merchant fulfillment configuration.
type MerchantFulfillmentConfig struct {
	// AllowsMethodCombinations specifies allowed method type combinations.
	AllowsMethodCombinations [][]FulfillmentMethodType `json:"allows_method_combinations,omitempty"`

	// AllowsMultiDestination specifies multi-destination settings.
	AllowsMultiDestination *AllowsMultiDestination `json:"allows_multi_destination,omitempty"`
}

// PlatformFulfillmentConfig represents platform fulfillment configuration.
type PlatformFulfillmentConfig struct {
	// SupportsMultiGroup indicates if the platform supports multiple groups.
	SupportsMultiGroup bool `json:"supports_multi_group,omitempty"`
}

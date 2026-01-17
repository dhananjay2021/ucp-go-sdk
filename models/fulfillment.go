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

// RetailLocationRequest represents a retail location in a request.
type RetailLocationRequest struct {
	// ID is the location identifier.
	ID string `json:"id,omitempty"`

	// Name is the location name.
	Name string `json:"name,omitempty"`

	// Address is the location address.
	Address *PostalAddress `json:"address,omitempty"`
}

// RetailLocationResponse represents a retail location in a response.
type RetailLocationResponse struct {
	// ID is the location identifier.
	ID string `json:"id"`

	// Name is the location name.
	Name string `json:"name"`

	// Address is the location address.
	Address *PostalAddress `json:"address,omitempty"`

	// Hours contains operating hours.
	Hours string `json:"hours,omitempty"`

	// Phone is the location phone number.
	Phone string `json:"phone,omitempty"`
}

// ShippingDestinationRequest represents a shipping destination in a request.
type ShippingDestinationRequest struct {
	// Address is the shipping address.
	Address *PostalAddress `json:"address,omitempty"`
}

// ShippingDestinationResponse represents a shipping destination in a response.
type ShippingDestinationResponse struct {
	// Address is the shipping address.
	Address *PostalAddress `json:"address,omitempty"`
}

// FulfillmentDestinationRequest represents a fulfillment destination in a request.
type FulfillmentDestinationRequest struct {
	// Shipping contains shipping destination details.
	Shipping *ShippingDestinationRequest `json:"shipping,omitempty"`

	// Pickup contains pickup location details.
	Pickup *RetailLocationRequest `json:"pickup,omitempty"`
}

// FulfillmentDestinationResponse represents a fulfillment destination in a response.
type FulfillmentDestinationResponse struct {
	// Shipping contains shipping destination details.
	Shipping *ShippingDestinationResponse `json:"shipping,omitempty"`

	// Pickup contains pickup location details.
	Pickup *RetailLocationResponse `json:"pickup,omitempty"`
}

// FulfillmentMethodRequest represents a fulfillment method in a request.
type FulfillmentMethodRequest struct {
	// ID is the method identifier.
	ID string `json:"id,omitempty"`

	// Type is the delivery method type.
	Type MethodType `json:"type,omitempty"`
}

// FulfillmentMethodResponse represents a fulfillment method in a response.
type FulfillmentMethodResponse struct {
	// ID is the method identifier.
	ID string `json:"id"`

	// Type is the delivery method type.
	Type MethodType `json:"type"`

	// Name is the display name.
	Name string `json:"name,omitempty"`

	// Description provides additional details.
	Description string `json:"description,omitempty"`

	// Price is the cost as a string.
	Price string `json:"price,omitempty"`

	// Expectation contains delivery timing.
	Expectation *Expectation `json:"expectation,omitempty"`

	// Selected indicates if this method is selected.
	Selected bool `json:"selected,omitempty"`
}

// FulfillmentAvailableMethodRequest represents an available method query.
type FulfillmentAvailableMethodRequest struct {
	// Type filters by method type.
	Type MethodType `json:"type,omitempty"`

	// Destination is the delivery destination.
	Destination *FulfillmentDestinationRequest `json:"destination,omitempty"`
}

// FulfillmentAvailableMethodResponse represents an available fulfillment method.
type FulfillmentAvailableMethodResponse struct {
	// ID is the method identifier.
	ID string `json:"id"`

	// Type is the delivery method type.
	Type MethodType `json:"type"`

	// Name is the display name.
	Name string `json:"name"`

	// Description provides additional details.
	Description string `json:"description,omitempty"`

	// Price is the cost as a string.
	Price string `json:"price,omitempty"`

	// Expectation contains delivery timing.
	Expectation *Expectation `json:"expectation,omitempty"`
}

// FulfillmentOptionRequest represents a fulfillment option in a request.
type FulfillmentOptionRequest struct {
	// MethodID is the selected method ID.
	MethodID string `json:"method_id,omitempty"`

	// Destination is the delivery destination.
	Destination *FulfillmentDestinationRequest `json:"destination,omitempty"`
}

// FulfillmentOptionResponse represents a fulfillment option in a response.
type FulfillmentOptionResponse struct {
	// Method is the selected method.
	Method *FulfillmentMethodResponse `json:"method,omitempty"`

	// Destination is the delivery destination.
	Destination *FulfillmentDestinationResponse `json:"destination,omitempty"`

	// AvailableMethods lists available methods.
	AvailableMethods []FulfillmentAvailableMethodResponse `json:"available_methods,omitempty"`
}

// FulfillmentGroupRequest represents a fulfillment group in a request.
type FulfillmentGroupRequest struct {
	// ID is the group identifier.
	ID string `json:"id,omitempty"`

	// LineItemIDs are the line items in this group.
	LineItemIDs []string `json:"line_item_ids,omitempty"`

	// Option is the selected fulfillment option.
	Option *FulfillmentOptionRequest `json:"option,omitempty"`
}

// FulfillmentGroupResponse represents a fulfillment group in a response.
type FulfillmentGroupResponse struct {
	// ID is the group identifier.
	ID string `json:"id"`

	// LineItemIDs are the line items in this group.
	LineItemIDs []string `json:"line_item_ids,omitempty"`

	// Option is the selected fulfillment option.
	Option *FulfillmentOptionResponse `json:"option,omitempty"`
}

// FulfillmentRequest represents fulfillment data in a request.
type FulfillmentRequest struct {
	// Groups are the fulfillment groups.
	Groups []FulfillmentGroupRequest `json:"groups,omitempty"`

	// Destination is a global destination.
	Destination *FulfillmentDestinationRequest `json:"destination,omitempty"`

	// MethodID is a global method selection.
	MethodID string `json:"method_id,omitempty"`
}

// FulfillmentResponse represents fulfillment data in a response.
type FulfillmentResponse struct {
	// Groups are the fulfillment groups.
	Groups []FulfillmentGroupResponse `json:"groups,omitempty"`

	// AvailableMethods lists globally available methods.
	AvailableMethods []FulfillmentAvailableMethodResponse `json:"available_methods,omitempty"`

	// Destination is the global destination.
	Destination *FulfillmentDestinationResponse `json:"destination,omitempty"`

	// Method is the globally selected method.
	Method *FulfillmentMethodResponse `json:"method,omitempty"`
}

// FulfillmentCreateRequest represents fulfillment in a checkout create request.
type FulfillmentCreateRequest struct {
	// Destination is the delivery destination.
	Destination *FulfillmentDestinationRequest `json:"destination,omitempty"`

	// MethodID is the selected method ID.
	MethodID string `json:"method_id,omitempty"`

	// Groups are the fulfillment groups.
	Groups []FulfillmentGroupRequest `json:"groups,omitempty"`
}

// FulfillmentUpdateRequest represents fulfillment in a checkout update request.
type FulfillmentUpdateRequest struct {
	// Destination is the delivery destination.
	Destination *FulfillmentDestinationRequest `json:"destination,omitempty"`

	// MethodID is the selected method ID.
	MethodID string `json:"method_id,omitempty"`

	// Groups are the fulfillment groups.
	Groups []FulfillmentGroupRequest `json:"groups,omitempty"`
}

// FulfillmentEvent represents a fulfillment tracking event.
type FulfillmentEvent struct {
	// Timestamp is when the event occurred.
	Timestamp *time.Time `json:"timestamp,omitempty"`

	// Status is the event status.
	Status string `json:"status,omitempty"`

	// Description provides event details.
	Description string `json:"description,omitempty"`

	// Location is where the event occurred.
	Location string `json:"location,omitempty"`
}

// MerchantFulfillmentConfig represents merchant fulfillment configuration.
type MerchantFulfillmentConfig struct {
	// SupportedMethods lists supported method types.
	SupportedMethods []MethodType `json:"supported_methods,omitempty"`

	// PickupLocations lists available pickup locations.
	PickupLocations []RetailLocationResponse `json:"pickup_locations,omitempty"`
}

// PlatformFulfillmentConfig represents platform fulfillment configuration.
type PlatformFulfillmentConfig struct {
	// PreferredMethods lists preferred method types.
	PreferredMethods []MethodType `json:"preferred_methods,omitempty"`
}

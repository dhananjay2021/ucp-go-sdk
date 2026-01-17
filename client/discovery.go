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

package client

import (
	"github.com/Universal-Commerce-Protocol/go-sdk/models"
)

// HasCapability checks if the profile supports a given capability name.
func HasCapability(profile *models.UCPProfile, capabilityName models.CapabilityName) bool {
	if profile == nil {
		return false
	}
	for _, cap := range profile.UCP.Capabilities {
		if cap.Name == capabilityName {
			return true
		}
	}
	return false
}

// GetCapability returns a capability by name, or nil if not found.
func GetCapability(profile *models.UCPProfile, capabilityName models.CapabilityName) *models.CapabilityDiscovery {
	if profile == nil {
		return nil
	}
	for i, cap := range profile.UCP.Capabilities {
		if cap.Name == capabilityName {
			return &profile.UCP.Capabilities[i]
		}
	}
	return nil
}

// GetServiceEndpoint returns the REST endpoint for a service, or empty string if not found.
func GetServiceEndpoint(profile *models.UCPProfile, serviceName string) string {
	if profile == nil {
		return ""
	}
	if service, ok := profile.UCP.Services[serviceName]; ok {
		if service.Rest != nil {
			return service.Rest.Endpoint
		}
	}
	return ""
}

// GetPaymentHandlers returns the payment handlers from the profile.
func GetPaymentHandlers(profile *models.UCPProfile) []models.PaymentHandlerResponse {
	if profile == nil || profile.Payment == nil {
		return nil
	}
	return profile.Payment.Handlers
}

// GetPaymentHandler returns a payment handler by ID, or nil if not found.
func GetPaymentHandler(profile *models.UCPProfile, handlerID string) *models.PaymentHandlerResponse {
	handlers := GetPaymentHandlers(profile)
	for i, h := range handlers {
		if h.ID == handlerID {
			return &handlers[i]
		}
	}
	return nil
}

// Well-known capability names.
const (
	CapabilityCheckout        models.CapabilityName = "dev.ucp.shopping.checkout"
	CapabilityOrder           models.CapabilityName = "dev.ucp.shopping.order"
	CapabilityIdentityLinking models.CapabilityName = "dev.ucp.identity_linking"
	CapabilityFulfillment     models.CapabilityName = "dev.ucp.shopping.fulfillment"
	CapabilityDiscount        models.CapabilityName = "dev.ucp.shopping.discount"
	CapabilityBuyerConsent    models.CapabilityName = "dev.ucp.shopping.buyer_consent"
	CapabilityPayment         models.CapabilityName = "dev.ucp.shopping.payment"
)

// Well-known service names.
const (
	ServiceShopping = "dev.ucp.shopping"
)

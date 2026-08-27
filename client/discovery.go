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
	"github.com/dhananjay2021/ucp-go-sdk/models"
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
	CapabilityCheckout      models.CapabilityName = "dev.ucp.shopping.checkout"
	CapabilityCart          models.CapabilityName = "dev.ucp.shopping.cart"
	CapabilityCatalogSearch models.CapabilityName = "dev.ucp.shopping.catalog.search"
	CapabilityCatalogLookup models.CapabilityName = "dev.ucp.shopping.catalog.lookup"
	CapabilityOrder         models.CapabilityName = "dev.ucp.shopping.order"
	CapabilityFulfillment   models.CapabilityName = "dev.ucp.shopping.fulfillment"
	CapabilityDiscount      models.CapabilityName = "dev.ucp.shopping.discount"
	CapabilityBuyerConsent  models.CapabilityName = "dev.ucp.shopping.buyer_consent"
	CapabilityPayment       models.CapabilityName = "dev.ucp.shopping.payment"
	CapabilityPermalink     models.CapabilityName = "dev.ucp.shopping.permalink"

	// Capabilities that live in the common/ namespace as of the 2026-08-25 UCP
	// release.
	CapabilityIdentityLinking   models.CapabilityName = "dev.ucp.common.identity_linking"
	CapabilityLocationSearch    models.CapabilityName = "dev.ucp.common.location.search"
	CapabilityLocationLookup    models.CapabilityName = "dev.ucp.common.location.lookup"
	CapabilityLoyalty           models.CapabilityName = "dev.ucp.common.loyalty"
	CapabilityPaymentTerms      models.CapabilityName = "dev.ucp.common.payment.terms"
	CapabilitySplitPayments     models.CapabilityName = "dev.ucp.common.payment.split_payments"
	CapabilityPaymentAuth       models.CapabilityName = "dev.ucp.common.payment.authentication"
	CapabilityPaymentAP2Mandate models.CapabilityName = "dev.ucp.common.payment.ap2_mandate"

	// CapabilityIdentityLinkingLegacy is the pre-2026-08-25 identity linking
	// capability name.
	//
	// Deprecated: use CapabilityIdentityLinking. Retained so profiles published
	// before the namespace move still resolve.
	CapabilityIdentityLinkingLegacy models.CapabilityName = "dev.ucp.identity_linking"
)

// Well-known service names.
const (
	ServiceShopping = "dev.ucp.shopping"
)

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

package validation

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dhananjay2021/ucp-go-sdk/models"
)

// CapabilityNegotiator handles capability negotiation between platform and business.
type CapabilityNegotiator struct {
	platformCapabilities []models.CapabilityDiscovery
}

// NewCapabilityNegotiator creates a new capability negotiator.
func NewCapabilityNegotiator(platformCapabilities []models.CapabilityDiscovery) *CapabilityNegotiator {
	return &CapabilityNegotiator{
		platformCapabilities: platformCapabilities,
	}
}

// NegotiationResult contains the result of capability negotiation.
type NegotiationResult struct {
	// Success indicates if negotiation was successful.
	Success bool

	// CommonCapabilities are the capabilities supported by both parties.
	CommonCapabilities []models.CapabilityDiscovery

	// MissingRequired lists required capabilities not supported by the business.
	MissingRequired []models.CapabilityName

	// VersionMismatches lists capabilities with incompatible versions.
	VersionMismatches []VersionMismatch

	// NegotiatedVersion is the agreed-upon protocol version.
	NegotiatedVersion models.Version
}

// VersionMismatch represents a version mismatch between capabilities.
type VersionMismatch struct {
	Capability      models.CapabilityName
	PlatformVersion models.Version
	BusinessVersion models.Version
}

// Negotiate performs capability negotiation with a business profile.
func (n *CapabilityNegotiator) Negotiate(businessProfile *models.UCPProfile, requiredCapabilities []models.CapabilityName) *NegotiationResult {
	result := &NegotiationResult{
		Success: true,
	}

	// Build map of business capabilities
	businessCaps := make(map[models.CapabilityName]models.CapabilityDiscovery)
	for _, cap := range businessProfile.UCP.Capabilities {
		businessCaps[cap.Name] = cap
	}

	// Find common capabilities
	for _, platformCap := range n.platformCapabilities {
		if businessCap, ok := businessCaps[platformCap.Name]; ok {
			// Check version compatibility
			if !versionsCompatible(platformCap.Version, businessCap.Version) {
				result.VersionMismatches = append(result.VersionMismatches, VersionMismatch{
					Capability:      platformCap.Name,
					PlatformVersion: platformCap.Version,
					BusinessVersion: businessCap.Version,
				})
				continue
			}

			// Use the older version
			negotiatedCap := platformCap
			if compareVersions(businessCap.Version, platformCap.Version) < 0 {
				negotiatedCap.Version = businessCap.Version
			}
			result.CommonCapabilities = append(result.CommonCapabilities, negotiatedCap)
		}
	}

	// Check required capabilities
	for _, required := range requiredCapabilities {
		found := false
		for _, common := range result.CommonCapabilities {
			if common.Name == required {
				found = true
				break
			}
		}
		if !found {
			result.MissingRequired = append(result.MissingRequired, required)
			result.Success = false
		}
	}

	// Version mismatches also fail negotiation
	if len(result.VersionMismatches) > 0 {
		result.Success = false
	}

	// Negotiate protocol version
	result.NegotiatedVersion = negotiateProtocolVersion(
		businessProfile.UCP.Version,
		n.getMinPlatformVersion(),
	)

	return result
}

// HasCapability checks if the result includes a specific capability.
func (r *NegotiationResult) HasCapability(name models.CapabilityName) bool {
	for _, cap := range r.CommonCapabilities {
		if cap.Name == name {
			return true
		}
	}
	return false
}

// GetCapability returns a capability from the result, or nil if not found.
func (r *NegotiationResult) GetCapability(name models.CapabilityName) *models.CapabilityDiscovery {
	for i, cap := range r.CommonCapabilities {
		if cap.Name == name {
			return &r.CommonCapabilities[i]
		}
	}
	return nil
}

// versionsCompatible checks if two versions are compatible.
// UCP versions are in YYYY-MM-DD format.
// Currently, we require exact match for major version (year).
func versionsCompatible(v1, v2 models.Version) bool {
	if !v1.IsValid() || !v2.IsValid() {
		return false
	}

	// Extract years
	year1 := strings.Split(string(v1), "-")[0]
	year2 := strings.Split(string(v2), "-")[0]

	// Same year = compatible
	return year1 == year2
}

// compareVersions compares two versions.
// Returns -1 if v1 < v2, 0 if equal, 1 if v1 > v2.
func compareVersions(v1, v2 models.Version) int {
	// Parse as dates
	parts1 := strings.Split(string(v1), "-")
	parts2 := strings.Split(string(v2), "-")

	for i := 0; i < 3; i++ {
		n1, _ := strconv.Atoi(parts1[i])
		n2, _ := strconv.Atoi(parts2[i])
		if n1 < n2 {
			return -1
		}
		if n1 > n2 {
			return 1
		}
	}
	return 0
}

// negotiateProtocolVersion returns the lower of two versions.
func negotiateProtocolVersion(v1, v2 models.Version) models.Version {
	if compareVersions(v1, v2) < 0 {
		return v1
	}
	return v2
}

// getMinPlatformVersion returns the minimum version from platform capabilities.
func (n *CapabilityNegotiator) getMinPlatformVersion() models.Version {
	if len(n.platformCapabilities) == 0 {
		return ""
	}

	minVersion := n.platformCapabilities[0].Version
	for _, cap := range n.platformCapabilities[1:] {
		if compareVersions(cap.Version, minVersion) < 0 {
			minVersion = cap.Version
		}
	}
	return minVersion
}

// ValidateCapabilityName checks if a capability name is valid.
func ValidateCapabilityName(name models.CapabilityName) error {
	if !name.IsValid() {
		return fmt.Errorf("invalid capability name: %s (must be reverse-domain notation)", name)
	}
	return nil
}

// ValidateVersion checks if a version string is valid.
func ValidateVersion(version models.Version) error {
	if !version.IsValid() {
		return fmt.Errorf("invalid version: %s (must be YYYY-MM-DD format)", version)
	}
	return nil
}

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

import (
	"encoding/json"
	"regexp"
)

// Version represents a UCP protocol version in YYYY-MM-DD format.
type Version string

// VersionPattern is the regex pattern for valid UCP versions.
var VersionPattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

// IsValid checks if the version matches the required pattern.
func (v Version) IsValid() bool {
	return VersionPattern.MatchString(string(v))
}

// CapabilityName represents a stable capability identifier in reverse-domain notation.
type CapabilityName string

// CapabilityNamePattern is the regex pattern for valid capability names.
var CapabilityNamePattern = regexp.MustCompile(`^[a-z][a-z0-9]*(?:\.[a-z][a-z0-9_]*)+$`)

// IsValid checks if the capability name matches the required pattern.
func (c CapabilityName) IsValid() bool {
	return CapabilityNamePattern.MatchString(string(c))
}

// CapabilityBase contains the common fields for all capability declarations.
type CapabilityBase struct {
	// Name is a stable capability identifier in reverse-domain notation.
	// Example: "dev.ucp.shopping.checkout"
	Name CapabilityName `json:"name,omitempty"`

	// Version is the capability version in YYYY-MM-DD format.
	Version Version `json:"version,omitempty"`

	// Spec is a URL to human-readable specification document.
	Spec string `json:"spec,omitempty"`

	// Schema is a URL to JSON Schema for this capability's payload.
	Schema string `json:"schema,omitempty"`

	// Extends is the parent capability this extends.
	// Present for extensions, absent for root capabilities.
	Extends CapabilityName `json:"extends,omitempty"`

	// Config contains capability-specific configuration.
	Config map[string]interface{} `json:"config,omitempty"`

	// AdditionalProperties captures any extra fields.
	AdditionalProperties map[string]interface{} `json:"-"`
}

// MarshalJSON implements custom JSON marshaling to include additional properties.
func (c CapabilityBase) MarshalJSON() ([]byte, error) {
	type Alias CapabilityBase
	data, err := json.Marshal(Alias(c))
	if err != nil {
		return nil, err
	}

	if len(c.AdditionalProperties) == 0 {
		return data, nil
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	for k, v := range c.AdditionalProperties {
		if _, exists := m[k]; !exists {
			m[k] = v
		}
	}

	return json.Marshal(m)
}

// CapabilityDiscovery is a full capability declaration for discovery profiles.
// Includes spec/schema URLs for agent fetching.
type CapabilityDiscovery struct {
	CapabilityBase
}

// CapabilityResponse is a capability reference in responses.
// Only name/version required to confirm active capabilities.
type CapabilityResponse struct {
	CapabilityBase
}

// RestTransport represents a REST transport binding.
type RestTransport struct {
	// Schema is a URL to OpenAPI 3.x specification (JSON format).
	Schema string `json:"schema"`

	// Endpoint is the merchant's REST API endpoint.
	Endpoint string `json:"endpoint"`

	// AdditionalProperties captures any extra fields.
	AdditionalProperties map[string]interface{} `json:"-"`
}

// MCPTransport represents an MCP transport binding.
type MCPTransport struct {
	// Schema is a URL to OpenRPC specification (JSON format).
	Schema string `json:"schema"`

	// Endpoint is the merchant's MCP endpoint.
	Endpoint string `json:"endpoint"`

	// AdditionalProperties captures any extra fields.
	AdditionalProperties map[string]interface{} `json:"-"`
}

// A2ATransport represents an A2A transport binding.
type A2ATransport struct {
	// Endpoint is the merchant's Agent Card endpoint.
	Endpoint string `json:"endpoint"`

	// AdditionalProperties captures any extra fields.
	AdditionalProperties map[string]interface{} `json:"-"`
}

// EmbeddedTransport represents an embedded transport binding (JSON-RPC 2.0 over postMessage).
type EmbeddedTransport struct {
	// Schema is a URL to OpenRPC specification (JSON format).
	Schema string `json:"schema"`

	// AdditionalProperties captures any extra fields.
	AdditionalProperties map[string]interface{} `json:"-"`
}

// EmbeddedTransportConfig represents per-checkout configuration for embedded transport binding.
// Allows businesses to vary ECP availability and delegations based on cart contents,
// agent authorization, or policy.
type EmbeddedTransportConfig struct {
	// Delegate specifies delegations the business allows.
	// At service-level, declares available delegations.
	// In checkout responses, confirms accepted delegations for this session.
	Delegate []string `json:"delegate,omitempty"`
}

// UCPService represents a service definition with transport bindings.
type UCPService struct {
	// Version is the service version in YYYY-MM-DD format.
	Version Version `json:"version"`

	// Spec is a URL to service documentation.
	Spec string `json:"spec"`

	// Rest is the REST transport binding.
	Rest *RestTransport `json:"rest,omitempty"`

	// MCP is the MCP transport binding.
	MCP *MCPTransport `json:"mcp,omitempty"`

	// A2A is the A2A transport binding.
	A2A *A2ATransport `json:"a2a,omitempty"`

	// Embedded is the embedded transport binding.
	Embedded *EmbeddedTransport `json:"embedded,omitempty"`

	// AdditionalProperties captures any extra fields.
	AdditionalProperties map[string]interface{} `json:"-"`
}

// Services is a map of service definitions keyed by reverse-domain service name.
type Services map[string]UCPService

// DiscoveryProfile represents the full UCP metadata for /.well-known/ucp discovery.
type DiscoveryProfile struct {
	// Version is the UCP protocol version.
	Version Version `json:"version"`

	// Services contains service definitions keyed by service name.
	Services Services `json:"services"`

	// Capabilities lists the supported capabilities and extensions.
	Capabilities []CapabilityDiscovery `json:"capabilities"`

	// AdditionalProperties captures any extra fields.
	AdditionalProperties map[string]interface{} `json:"-"`
}

// ResponseCheckout represents UCP metadata for checkout responses.
type ResponseCheckout struct {
	// Version is the UCP protocol version.
	Version Version `json:"version"`

	// Capabilities lists the active capabilities for this response.
	Capabilities []CapabilityResponse `json:"capabilities"`

	// AdditionalProperties captures any extra fields.
	AdditionalProperties map[string]interface{} `json:"-"`
}

// ResponseOrder represents UCP metadata for order responses.
type ResponseOrder struct {
	// Version is the UCP protocol version.
	Version Version `json:"version"`

	// Capabilities lists the active capabilities for this response.
	Capabilities []CapabilityResponse `json:"capabilities"`

	// AdditionalProperties captures any extra fields.
	AdditionalProperties map[string]interface{} `json:"-"`
}

// JWK represents a JSON Web Key for signature verification.
type JWK struct {
	// Kid is the key ID, referenced in signature headers.
	Kid string `json:"kid"`

	// Kty is the key type (e.g., "EC", "RSA").
	Kty string `json:"kty"`

	// Crv is the curve name for EC keys (e.g., "P-256").
	Crv string `json:"crv,omitempty"`

	// X is the X coordinate for EC public keys (base64url encoded).
	X string `json:"x,omitempty"`

	// Y is the Y coordinate for EC public keys (base64url encoded).
	Y string `json:"y,omitempty"`

	// N is the modulus for RSA public keys (base64url encoded).
	N string `json:"n,omitempty"`

	// E is the exponent for RSA public keys (base64url encoded).
	E string `json:"e,omitempty"`

	// Use is the key usage ("sig" for signing keys).
	Use string `json:"use,omitempty"`

	// Alg is the algorithm (e.g., "ES256", "RS256").
	Alg string `json:"alg,omitempty"`
}

// UCPProfile represents the full discovery profile returned from /.well-known/ucp.
type UCPProfile struct {
	// UCP contains the protocol metadata.
	UCP DiscoveryProfile `json:"ucp"`

	// Payment contains payment configuration.
	Payment *PaymentConfig `json:"payment,omitempty"`

	// SigningKeys are public keys for signature verification.
	SigningKeys []JWK `json:"signing_keys,omitempty"`

	// AdditionalProperties captures any extra fields.
	AdditionalProperties map[string]interface{} `json:"-"`
}

// PaymentConfig represents payment configuration in the discovery profile.
type PaymentConfig struct {
	// Handlers contains payment handler definitions.
	Handlers []PaymentHandlerResponse `json:"handlers,omitempty"`
}

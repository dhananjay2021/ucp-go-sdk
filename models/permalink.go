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

// PermalinkConfig is the Business browser endpoint configuration for shopping
// permalinks. Additional business-specific configuration keys are permitted.
type PermalinkConfig struct {
	// Endpoint is the absolute HTTPS browser endpoint with a non-empty
	// authority and without userinfo, query, fragment, whitespace,
	// backslashes, or trailing slash. Optional compact item path and query
	// parameters are appended to this endpoint.
	Endpoint string `json:"endpoint"`

	// AdditionalProperties captures any extra business-specific
	// configuration keys permitted by the config schema.
	AdditionalProperties map[string]interface{} `json:"-"`
}

// PermalinkCapabilityBase contains the capability declaration metadata shared
// by the platform, business, and response permalink capability variants. It
// mirrors the common UCP capability fields for the browser-addressable
// shopping intent capability, which defines no shopping-state fields of its
// own; permalink query parameters address existing UCP field paths.
type PermalinkCapabilityBase struct {
	// Name is the stable capability identifier in reverse-domain notation,
	// for example "dev.ucp.shopping.permalink".
	Name CapabilityName `json:"name,omitempty"`

	// Version is the capability version in YYYY-MM-DD format.
	Version Version `json:"version,omitempty"`

	// Spec is a URL to the human-readable permalink specification.
	Spec string `json:"spec,omitempty"`

	// Schema is a URL to the JSON Schema for this capability's payload.
	Schema string `json:"schema,omitempty"`

	// Extends is the parent capability this declaration extends, when present.
	Extends CapabilityName `json:"extends,omitempty"`
}

// PermalinkPlatformCapability is the platform-level permalink capability
// declaration. Platforms advertise support for shopping permalinks; the
// permalink endpoint is business-specific, so no platform config is required.
type PermalinkPlatformCapability struct {
	PermalinkCapabilityBase
}

// PermalinkBusinessCapability is the Business declaration for the UCP shopping
// permalink browser endpoint.
type PermalinkBusinessCapability struct {
	PermalinkCapabilityBase

	// Config is the Business browser endpoint configuration (required).
	Config PermalinkConfig `json:"config"`
}

// PermalinkResponseCapability is the shopping permalink capability reference
// included in API responses.
type PermalinkResponseCapability struct {
	PermalinkCapabilityBase
}

// PermalinkResolveParams captures the inputs to a permalink browser redirect
// route defined by the permalink binding. The Business resolves these against
// its active schemas and capabilities, then issues a 303 redirect to a
// buyer-facing destination. Permalink query syntax omits the leading '/' for
// UCP field paths, but Businesses may also accept canonical JSON Pointer names.
type PermalinkResolveParams struct {
	// Items is the compact, comma-separated list of item_id_token:quantity
	// pairs taken from the route path. It is absent for the no-item permalink
	// route.
	Items *string `json:"items,omitempty"`

	// ContinueTo is the optional root-relative destination preference applied
	// after the Business resolves permalink inputs. It MUST start with '/' and
	// MUST NOT start with '//', and MUST NOT contain a URL scheme, backslashes,
	// whitespace, or control characters.
	ContinueTo *string `json:"continue_to,omitempty"`
}

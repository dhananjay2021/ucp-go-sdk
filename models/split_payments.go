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

// BusinessSplitPaymentsConfig is the business-level configuration for split
// payments (UCP spec change #409). Declaring the capability means multiple
// payment instruments are supported; this config declares which combinations are
// valid.
type BusinessSplitPaymentsConfig struct {
	// AllowedCombinations is an array of valid instrument combinations. Each
	// combination is an array of instrument groups. A payment is valid if it
	// matches any combination.
	AllowedCombinations [][]InstrumentGroup `json:"allowed_combinations"`
}

// InstrumentGroup is a constraint within an allowed combination that defines
// which instrument types can fill this group and how many are permitted.
type InstrumentGroup struct {
	// Types are the instrument types accepted by this group (OR logic). Any listed
	// type qualifies.
	Types []string `json:"types"`

	// Min is the minimum number of instruments required from this group. Defaults
	// to 0 (optional) when omitted.
	Min int `json:"min,omitempty"`

	// Max is the maximum number of instruments allowed from this group. Defaults
	// to 1 when omitted. MUST be greater than or equal to Min.
	Max int `json:"max,omitempty"`
}

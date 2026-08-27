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

// External-oracle tests: rather than assert field shapes by hand, these tests
// validate hand-written model output against the *actual* UCP JSON Schemas for
// the 2026-08-25 release. The schemas are bundled (all $refs inlined) into
// testdata/ucp_bundle.json by scripts/bundle_schemas.py, so the test is
// hermetic. Each case marshals a hand-written value and validates it against the
// corresponding schema definition at #/$defs/<Name>.
//
// Scope: the leaf/value types this SDK models faithfully from the 2026-08-25
// schemas. Capability envelopes whose idiomatic Go shape intentionally diverges
// from the raw schema are covered by unit tests instead.

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

const oracleBundleID = "https://ucp.dev/schemas/_bundle.json"

// newOracleCompiler loads the bundled UCP schema and returns a compiler with the
// bundle registered under its canonical $id.
func newOracleCompiler(t *testing.T) *jsonschema.Compiler {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "ucp_bundle.json"))
	if err != nil {
		t.Fatalf("read bundle: %v (regenerate with scripts/bundle_schemas.py)", err)
	}
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("parse bundle: %v", err)
	}
	c := jsonschema.NewCompiler()
	if err := c.AddResource(oracleBundleID, doc); err != nil {
		t.Fatalf("add bundle resource: %v", err)
	}
	return c
}

// asJSONValue marshals v and decodes it back into the generic representation the
// validator expects.
func asJSONValue(t *testing.T, v interface{}) interface{} {
	t.Helper()
	raw, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	inst, err := jsonschema.UnmarshalJSON(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf("re-parse: %v", err)
	}
	return inst
}

// validateConforms fails if v does not validate against #/$defs/<def>.
func validateConforms(t *testing.T, c *jsonschema.Compiler, def string, v interface{}) {
	t.Helper()
	sch, err := c.Compile(oracleBundleID + "#/$defs/" + def)
	if err != nil {
		t.Fatalf("compile #/$defs/%s: %v", def, err)
	}
	if err := sch.Validate(asJSONValue(t, v)); err != nil {
		t.Errorf("%s did not conform to schema $defs/%s:\n%v", def, def, err)
	}
}

// validateRejects asserts the oracle actually bites: a deliberately invalid
// payload MUST fail validation.
func validateRejects(t *testing.T, c *jsonschema.Compiler, def string, invalid interface{}) {
	t.Helper()
	sch, err := c.Compile(oracleBundleID + "#/$defs/" + def)
	if err != nil {
		t.Fatalf("compile #/$defs/%s: %v", def, err)
	}
	if err := sch.Validate(asJSONValue(t, invalid)); err == nil {
		t.Errorf("expected schema $defs/%s to reject %#v, but it validated", def, invalid)
	}
}

func TestOracle_Conformance(t *testing.T) {
	c := newOracleCompiler(t)
	due := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)

	cases := []struct {
		def string
		val interface{}
	}{
		{"PanCredential", PanCredential{
			Type: CredentialTypePAN, Number: "4242424242424242",
			ExpiryMonth: 12, ExpiryYear: 2030, Name: "Jane Doe", CVC: "223",
		}},
		{"NetworkTokenCredential", NetworkTokenCredential{
			Type: CredentialTypeNetworkToken, Number: "5204240000004242",
			Cryptogram: "gXc5UCLnM6ckD7pjM1TdPA==", ECIValue: "07",
			TokenRequestorID: "12345678901",
		}},
		{"PaymentSchedule", PaymentSchedule{
			ID: "s1", Amount: 2500,
			Description: Description{Plain: "First payment today"},
			Type:        string(PaymentScheduleTypeImmediate),
		}},
		{"PaymentTerm", PaymentTerm{
			ID: "pay_in_4", Title: "Pay in 4",
			Schedules: []PaymentSchedule{
				{ID: "s1", Amount: 2500, Description: Description{Plain: "Today"}, Type: "immediate"},
				{ID: "s2", Amount: 2500, Description: Description{Plain: "Later"}, DueAt: &due, Type: "deferred"},
			},
		}},
		{"InstrumentGroup", InstrumentGroup{Types: []string{"card"}, Min: 1, Max: 1}},
		{"BusinessSplitPaymentsConfig", BusinessSplitPaymentsConfig{
			AllowedCombinations: [][]InstrumentGroup{{
				{Types: []string{"card"}, Min: 1, Max: 1},
				{Types: []string{"gift_card"}, Max: 3},
			}},
		}},
		{"Policy", Policy{
			Type:        "dev.ucp.shopping.policy.return",
			Description: Description{Plain: "30-day returns"},
			AppliesTo:   []string{"$.line_items[0]"},
			URL:         "https://example.test/returns",
		}},
		{"Geo", Geo{Latitude: 37.42, Longitude: -122.08}},
		{"Locality", Locality{
			AddressCountry: strptr("US"), AddressRegion: strptr("CA"),
			PostalCode: strptr("94043"),
		}},
		{"LocationSummary", LocationSummary{
			ID: "loc_1", Name: "Downtown Store",
			Address: &PostalAddress{AddressLocality: "Mountain View"},
		}},
	}

	for _, tc := range cases {
		t.Run(tc.def, func(t *testing.T) {
			validateConforms(t, c, tc.def, tc.val)
		})
	}
}

func TestOracle_RejectsInvalid(t *testing.T) {
	c := newOracleCompiler(t)

	// PAN credential MUST carry a number.
	validateRejects(t, c, "PanCredential", map[string]interface{}{"type": "pan"})
	// Network token credential MUST carry a cryptogram.
	validateRejects(t, c, "NetworkTokenCredential", map[string]interface{}{
		"type": "network_token", "number": "5204240000004242",
	})
	// Policy MUST carry a description.
	validateRejects(t, c, "Policy", map[string]interface{}{"type": "dev.ucp.shopping.policy.return"})
	// Geo MUST carry both coordinates.
	validateRejects(t, c, "Geo", map[string]interface{}{"latitude": 37.42})
	// Location summary MUST carry a name.
	validateRejects(t, c, "LocationSummary", map[string]interface{}{"id": "loc_1"})
}

func strptr(s string) *string { return &s }

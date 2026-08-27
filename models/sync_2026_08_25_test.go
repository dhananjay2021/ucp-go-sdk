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
	"reflect"
	"testing"
	"time"
)

// jsonKeys marshals v and returns the set of top-level object keys, so tests can
// assert wire shape (property names) against the schema without hand-writing the
// full JSON string.
func jsonKeys(t *testing.T, v interface{}) map[string]bool {
	t.Helper()
	raw, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err != nil {
		t.Fatalf("unmarshal to map: %v", err)
	}
	keys := make(map[string]bool, len(obj))
	for k := range obj {
		keys[k] = true
	}
	return keys
}

func assertHasKeys(t *testing.T, keys map[string]bool, want ...string) {
	t.Helper()
	for _, k := range want {
		if !keys[k] {
			t.Errorf("expected JSON key %q, got keys %v", k, keys)
		}
	}
}

func assertNoKeys(t *testing.T, keys map[string]bool, notWant ...string) {
	t.Helper()
	for _, k := range notWant {
		if keys[k] {
			t.Errorf("did not expect JSON key %q, got keys %v", k, keys)
		}
	}
}

func TestSpecVersion(t *testing.T) {
	if SpecVersion != "2026-08-25" {
		t.Fatalf("SpecVersion = %q, want 2026-08-25", SpecVersion)
	}
}

func TestUCPProfileVerificationKeys(t *testing.T) {
	canonical := JWK{Kid: "a", Kty: "OKP"}
	legacy := JWK{Kid: "b", Kty: "EC"}

	t.Run("prefers canonical keys", func(t *testing.T) {
		p := UCPProfile{Keys: []JWK{canonical}, SigningKeys: []JWK{legacy}}
		got := p.VerificationKeys()
		if len(got) != 1 || got[0].Kid != "a" {
			t.Fatalf("VerificationKeys() = %+v, want canonical keys[]", got)
		}
	})

	t.Run("falls back to signing_keys", func(t *testing.T) {
		p := UCPProfile{SigningKeys: []JWK{legacy}}
		got := p.VerificationKeys()
		if len(got) != 1 || got[0].Kid != "b" {
			t.Fatalf("VerificationKeys() = %+v, want deprecated signing_keys", got)
		}
	})

	t.Run("canonical keys serialize under keys", func(t *testing.T) {
		keys := jsonKeys(t, UCPProfile{
			UCP:  DiscoveryProfile{},
			Keys: []JWK{canonical},
		})
		assertHasKeys(t, keys, "keys")
		assertNoKeys(t, keys, "signing_keys")
	})
}

func TestPanCredentialJSON(t *testing.T) {
	c := PanCredential{
		Type:        CredentialTypePAN,
		Number:      "4242424242424242",
		ExpiryMonth: 12,
		ExpiryYear:  2030,
		Name:        "Jane Doe",
		CVC:         "223",
	}
	keys := jsonKeys(t, c)
	assertHasKeys(t, keys, "type", "number", "expiry_month", "expiry_year", "name", "cvc")
	// PAN credential is verified by CVC, never by a cryptogram.
	assertNoKeys(t, keys, "cryptogram", "eci_value", "token_requestor_id")

	var back PanCredential
	if err := json.Unmarshal([]byte(`{"type":"pan","number":"4242424242424242"}`), &back); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if back.Type != CredentialTypePAN || back.Number != "4242424242424242" {
		t.Fatalf("round-trip mismatch: %+v", back)
	}
}

func TestNetworkTokenCredentialJSON(t *testing.T) {
	c := NetworkTokenCredential{
		Type:             CredentialTypeNetworkToken,
		Number:           "5204240000004242",
		Cryptogram:       "gXc5UCLnM6ckD7pjM1TdPA==",
		ECIValue:         "07",
		TokenRequestorID: "12345678901",
	}
	keys := jsonKeys(t, c)
	assertHasKeys(t, keys, "type", "number", "cryptogram", "eci_value", "token_requestor_id")
	// A network token credential carries no raw CVC.
	assertNoKeys(t, keys, "cvc")
}

func TestPaymentTermJSON(t *testing.T) {
	due := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
	term := PaymentTerm{
		ID:    "pay_in_4",
		Title: "Pay in 4",
		Schedules: []PaymentSchedule{
			{
				ID:          "s1",
				Amount:      2500,
				Description: Description{Plain: "First payment today"},
				Type:        string(PaymentScheduleTypeImmediate),
			},
			{
				ID:          "s2",
				Amount:      2500,
				Description: Description{Plain: "Second payment in 2 weeks"},
				DueAt:       &due,
				Type:        "deferred",
			},
		},
	}
	keys := jsonKeys(t, term)
	assertHasKeys(t, keys, "id", "title", "schedules")
	// Optional description omitted when empty.
	assertNoKeys(t, keys, "description")

	raw, err := json.Marshal(term)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var back PaymentTerm
	if err := json.Unmarshal(raw, &back); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(term, back) {
		t.Fatalf("round-trip mismatch:\n got %+v\nwant %+v", back, term)
	}
}

func TestInstrumentGroupOmitempty(t *testing.T) {
	// Min/Max default to 0/1 in the schema; with zero values omitempty keeps the
	// wire form minimal, leaving defaults to the consumer.
	keys := jsonKeys(t, InstrumentGroup{Types: []string{"card"}})
	assertHasKeys(t, keys, "types")
	assertNoKeys(t, keys, "min", "max")

	keys = jsonKeys(t, InstrumentGroup{Types: []string{"gift_card"}, Min: 1, Max: 2})
	assertHasKeys(t, keys, "types", "min", "max")
}

func TestActionsRoundTrip(t *testing.T) {
	actions := Actions{
		ActionTypeThreeDSChallenge: {
			{ID: "act1", Config: map[string]interface{}{
				"payment_instrument_id": "pi_1",
				"url":                   "https://example.test/3ds",
			}},
		},
	}
	raw, err := json.Marshal(actions)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var back Actions
	if err := json.Unmarshal(raw, &back); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	got := back[ActionTypeThreeDSChallenge]
	if len(got) != 1 || got[0].ID != "act1" {
		t.Fatalf("actions round-trip mismatch: %+v", back)
	}
	if got[0].Config["url"] != "https://example.test/3ds" {
		t.Fatalf("action config not preserved: %+v", got[0].Config)
	}
}

func TestBusinessSplitPaymentsConfigJSON(t *testing.T) {
	cfg := BusinessSplitPaymentsConfig{
		AllowedCombinations: [][]InstrumentGroup{
			{
				{Types: []string{"card"}, Min: 1, Max: 1},
				{Types: []string{"gift_card"}, Max: 3},
			},
		},
	}
	keys := jsonKeys(t, cfg)
	assertHasKeys(t, keys, "allowed_combinations")

	raw, _ := json.Marshal(cfg)
	var back BusinessSplitPaymentsConfig
	if err := json.Unmarshal(raw, &back); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(back.AllowedCombinations) != 1 || len(back.AllowedCombinations[0]) != 2 {
		t.Fatalf("split payments config round-trip mismatch: %+v", back)
	}
}

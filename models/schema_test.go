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

package models_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dhananjay2021/ucp-go-sdk/models"
)

// schemaDir is the path to the UCP spec schemas directory.
// Assumes the ucp repo is a sibling to go-sdk.
const schemaDir = "../../ucp/spec/schemas/shopping"

// TestModelsMatchSchemaFields verifies that our model field names match the JSON schema.
// This catches mismatches like "AddressLines" vs "street_address".
func TestModelsMatchSchemaFields(t *testing.T) {
	tests := []struct {
		name       string
		schemaFile string
		model      interface{}
	}{
		{
			name:       "PostalAddress",
			schemaFile: "types/postal_address.json",
			model: models.PostalAddress{
				StreetAddress:   "123 Main St",
				ExtendedAddress: "Apt 4",
				AddressLocality: "San Francisco",
				AddressRegion:   "CA",
				AddressCountry:  "US",
				PostalCode:      "94102",
				FirstName:       "John",
				LastName:        "Doe",
				FullName:        "John Doe",
				PhoneNumber:     "+14155551234",
			},
		},
		{
			name:       "Buyer",
			schemaFile: "types/buyer.json",
			model: models.Buyer{
				FirstName:   "John",
				LastName:    "Doe",
				FullName:    "John Doe",
				Email:       "john@example.com",
				PhoneNumber: "+14155551234",
			},
		},
		{
			name:       "LineItemCreateRequest",
			schemaFile: "types/line_item.create_req.json",
			model: models.LineItemCreateRequest{
				Item: models.ItemCreateRequest{
					ID: "product-123",
				},
				Quantity: 2,
			},
		},
		{
			name:       "ItemResponse",
			schemaFile: "types/item_resp.json",
			model: models.ItemResponse{
				ID:       "product-123",
				Title:    "Test Product",
				Price:    1999,
				ImageURL: "https://example.com/image.jpg",
			},
		},
		{
			name:       "TotalResponse",
			schemaFile: "types/total_resp.json",
			model: models.TotalResponse{
				Type:        models.TotalTypeSubtotal,
				Amount:      1999,
				DisplayText: "Subtotal",
			},
		},
		{
			name:       "Link",
			schemaFile: "types/link.json",
			model: models.Link{
				Type:  "privacy_policy",
				Title: "Privacy Policy",
				URL:   "https://example.com/privacy",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal model to JSON
			jsonBytes, err := json.Marshal(tt.model)
			if err != nil {
				t.Fatalf("Failed to marshal %s: %v", tt.name, err)
			}

			// Parse JSON to get field names
			var jsonMap map[string]interface{}
			if err := json.Unmarshal(jsonBytes, &jsonMap); err != nil {
				t.Fatalf("Failed to unmarshal %s JSON: %v", tt.name, err)
			}

			// Load schema if available
			schemaPath := filepath.Join(schemaDir, tt.schemaFile)
			schemaBytes, err := os.ReadFile(schemaPath)
			if err != nil {
				t.Skipf("Schema file not found (run from go-sdk directory): %s", schemaPath)
			}

			// Parse schema
			var schema map[string]interface{}
			if err := json.Unmarshal(schemaBytes, &schema); err != nil {
				t.Fatalf("Failed to parse schema %s: %v", tt.schemaFile, err)
			}

			// Get schema properties
			properties, ok := schema["properties"].(map[string]interface{})
			if !ok {
				t.Skipf("Schema %s has no properties field", tt.schemaFile)
			}

			// Verify each JSON field exists in schema
			for field := range jsonMap {
				if _, exists := properties[field]; !exists {
					t.Errorf("Field %q in model %s not found in schema %s", field, tt.name, tt.schemaFile)
				}
			}

			t.Logf("%s: %d fields validated against schema", tt.name, len(jsonMap))
		})
	}
}

// TestModelJSONRoundTrip verifies models can be serialized and deserialized correctly.
func TestModelJSONRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		model interface{}
		check func(t *testing.T, decoded map[string]interface{})
	}{
		{
			name: "PostalAddress uses snake_case",
			model: models.PostalAddress{
				StreetAddress:   "123 Main St",
				AddressLocality: "San Francisco",
				AddressRegion:   "CA",
				PostalCode:      "94102",
			},
			check: func(t *testing.T, m map[string]interface{}) {
				assertFieldExists(t, m, "street_address")
				assertFieldExists(t, m, "address_locality")
				assertFieldExists(t, m, "address_region")
				assertFieldExists(t, m, "postal_code")
				// Verify wrong names are NOT present
				assertFieldNotExists(t, m, "AddressLines")
				assertFieldNotExists(t, m, "City")
				assertFieldNotExists(t, m, "State")
			},
		},
		{
			name: "ItemResponse uses integer price",
			model: models.ItemResponse{
				ID:    "prod-1",
				Title: "Widget",
				Price: 1999,
			},
			check: func(t *testing.T, m map[string]interface{}) {
				price, ok := m["price"].(float64) // JSON numbers are float64
				if !ok {
					t.Error("price field should be a number")
					return
				}
				if price != 1999 {
					t.Errorf("price should be 1999, got %v", price)
				}
			},
		},
		{
			name: "TotalResponse uses integer amounts",
			model: models.TotalResponse{
				Type:        models.TotalTypeSubtotal,
				Amount:      5000,
				DisplayText: "Subtotal",
			},
			check: func(t *testing.T, m map[string]interface{}) {
				// Verify amount is an integer (float64 in JSON)
				v, ok := m["amount"].(float64)
				if !ok {
					t.Error("amount should be a number")
					return
				}
				if v != float64(int(v)) {
					t.Errorf("amount should be an integer, got %v", v)
				}
				// Verify type field
				assertFieldExists(t, m, "type")
				if m["type"] != "subtotal" {
					t.Errorf("type should be 'subtotal', got %v", m["type"])
				}
			},
		},
		{
			name: "Link uses type field",
			model: models.Link{
				Type:  "terms_of_service",
				Title: "Terms of Service",
				URL:   "https://example.com/tos",
			},
			check: func(t *testing.T, m map[string]interface{}) {
				assertFieldExists(t, m, "type")
				assertFieldNotExists(t, m, "rel")
			},
		},
		{
			name: "CheckoutCreateRequest structure",
			model: models.CheckoutCreateRequest{
				Currency: "USD",
				LineItems: []models.LineItemCreateRequest{
					{
						Item:     models.ItemCreateRequest{ID: "prod-1"},
						Quantity: 1,
					},
				},
				Payment: models.PaymentCreateRequest{
					Instruments: []models.PaymentInstrument{
						{
							ID:        "card-1",
							HandlerID: "handler-1",
							Type:      "card",
						},
					},
					SelectedInstrumentID: "card-1",
				},
			},
			check: func(t *testing.T, m map[string]interface{}) {
				assertFieldExists(t, m, "currency")
				assertFieldExists(t, m, "line_items")
				assertFieldExists(t, m, "payment")

				lineItems, ok := m["line_items"].([]interface{})
				if !ok || len(lineItems) == 0 {
					t.Error("line_items should be a non-empty array")
					return
				}

				item := lineItems[0].(map[string]interface{})
				assertFieldExists(t, item, "item")
				assertFieldExists(t, item, "quantity")

				itemData := item["item"].(map[string]interface{})
				assertFieldExists(t, itemData, "id")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			jsonBytes, err := json.Marshal(tt.model)
			if err != nil {
				t.Fatalf("Failed to marshal: %v", err)
			}

			// Unmarshal to map for field inspection
			var m map[string]interface{}
			if err := json.Unmarshal(jsonBytes, &m); err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}

			tt.check(t, m)
		})
	}
}

// TestEnumValues verifies enum types have expected values.
func TestEnumValues(t *testing.T) {
	tests := []struct {
		name     string
		enumVal  string
		expected string
	}{
		{"CheckoutStatusIncomplete", string(models.CheckoutStatusIncomplete), "incomplete"},
		{"CheckoutStatusCompleted", string(models.CheckoutStatusCompleted), "completed"},
		{"CheckoutStatusCanceled", string(models.CheckoutStatusCanceled), "canceled"},
		{"CheckoutStatusRequiresEscalation", string(models.CheckoutStatusRequiresEscalation), "requires_escalation"},
		{"CheckoutStatusReadyForComplete", string(models.CheckoutStatusReadyForComplete), "ready_for_complete"},
		{"FulfillmentMethodTypeShipping", string(models.FulfillmentMethodTypeShipping), "shipping"},
		{"FulfillmentMethodTypePickup", string(models.FulfillmentMethodTypePickup), "pickup"},
		{"AllocationMethodAcross", string(models.AllocationMethodAcross), "across"},
		{"AllocationMethodEach", string(models.AllocationMethodEach), "each"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.enumVal != tt.expected {
				t.Errorf("%s = %q, want %q", tt.name, tt.enumVal, tt.expected)
			}
		})
	}
}

// TestSchemaRequiredFields verifies required fields in request types.
func TestSchemaRequiredFields(t *testing.T) {
	tests := []struct {
		name           string
		schemaFile     string
		requiredFields []string
	}{
		{
			name:           "CheckoutCreateRequest",
			schemaFile:     "checkout.create_req.json",
			requiredFields: []string{"currency", "line_items", "payment"},
		},
		{
			name:           "LineItemCreateRequest",
			schemaFile:     "types/line_item.create_req.json",
			requiredFields: []string{"item", "quantity"},
		},
		{
			name:           "ItemCreateRequest",
			schemaFile:     "types/item.create_req.json",
			requiredFields: []string{"id"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schemaPath := filepath.Join(schemaDir, tt.schemaFile)
			schemaBytes, err := os.ReadFile(schemaPath)
			if err != nil {
				t.Skipf("Schema file not found: %s", schemaPath)
			}

			var schema map[string]interface{}
			if err := json.Unmarshal(schemaBytes, &schema); err != nil {
				t.Fatalf("Failed to parse schema: %v", err)
			}

			required, ok := schema["required"].([]interface{})
			if !ok {
				t.Skipf("Schema has no required field")
			}

			requiredSet := make(map[string]bool)
			for _, r := range required {
				if s, ok := r.(string); ok {
					requiredSet[s] = true
				}
			}

			for _, field := range tt.requiredFields {
				if !requiredSet[field] {
					t.Errorf("Expected %q to be required in %s", field, tt.schemaFile)
				}
			}
		})
	}
}

// TestFieldNamingConventions verifies JSON field names use correct conventions.
func TestFieldNamingConventions(t *testing.T) {
	// Create a model with all addressable fields
	addr := models.PostalAddress{
		StreetAddress:   "123 Main St",
		ExtendedAddress: "Apt 4",
		AddressLocality: "San Francisco",
		AddressRegion:   "CA",
		AddressCountry:  "US",
		PostalCode:      "94102",
		FirstName:       "John",
		LastName:        "Doe",
		PhoneNumber:     "+1234567890",
	}

	jsonBytes, err := json.Marshal(addr)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	jsonStr := string(jsonBytes)

	// Verify Schema.org naming (snake_case)
	expectedFields := []string{
		"street_address",
		"extended_address",
		"address_locality",
		"address_region",
		"address_country",
		"postal_code",
		"first_name",
		"last_name",
		"phone_number",
	}

	for _, field := range expectedFields {
		if !strings.Contains(jsonStr, `"`+field+`"`) {
			t.Errorf("Expected field %q not found in JSON output", field)
		}
	}

	// Verify old/wrong naming is NOT present
	wrongFields := []string{
		"AddressLines",
		"City",
		"State",
		"ZipCode",
		"Country",
		"FirstName",  // Should be snake_case
		"LastName",   // Should be snake_case
		"PhoneNumber", // Should be snake_case
	}

	for _, field := range wrongFields {
		if strings.Contains(jsonStr, `"`+field+`"`) {
			t.Errorf("Wrong field naming %q found in JSON output (should be snake_case)", field)
		}
	}
}

// Helper functions
func assertFieldExists(t *testing.T, m map[string]interface{}, field string) {
	t.Helper()
	if _, ok := m[field]; !ok {
		t.Errorf("Expected field %q not found in JSON", field)
	}
}

func assertFieldNotExists(t *testing.T, m map[string]interface{}, field string) {
	t.Helper()
	if _, ok := m[field]; ok {
		t.Errorf("Unexpected field %q found in JSON", field)
	}
}

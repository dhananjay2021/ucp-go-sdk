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

// Package validation provides JSON Schema validation and capability negotiation.
package validation

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

// SchemaValidator validates JSON data against UCP schemas.
type SchemaValidator struct {
	schemaCache map[string][]byte
	mu          sync.RWMutex
	httpClient  *http.Client
}

// NewSchemaValidator creates a new schema validator.
func NewSchemaValidator() *SchemaValidator {
	return &SchemaValidator{
		schemaCache: make(map[string][]byte),
		httpClient:  &http.Client{},
	}
}

// ValidationError represents a schema validation error.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s", e.Field, e.Message)
	}
	return e.Message
}

// ValidationResult contains the result of schema validation.
type ValidationResult struct {
	Valid  bool              `json:"valid"`
	Errors []ValidationError `json:"errors,omitempty"`
}

// LoadSchema loads a schema from a URL and caches it.
func (v *SchemaValidator) LoadSchema(url string) ([]byte, error) {
	v.mu.RLock()
	if schema, ok := v.schemaCache[url]; ok {
		v.mu.RUnlock()
		return schema, nil
	}
	v.mu.RUnlock()

	// Fetch schema
	resp, err := v.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch schema from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch schema from %s: status %d", url, resp.StatusCode)
	}

	schema, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema from %s: %w", url, err)
	}

	// Cache the schema
	v.mu.Lock()
	v.schemaCache[url] = schema
	v.mu.Unlock()

	return schema, nil
}

// LoadSchemaFromBytes loads a schema from bytes and caches it under a key.
func (v *SchemaValidator) LoadSchemaFromBytes(key string, schema []byte) {
	v.mu.Lock()
	v.schemaCache[key] = schema
	v.mu.Unlock()
}

// ValidateJSON performs basic JSON validation.
// Note: For full JSON Schema validation, use a library like github.com/santhosh-tekuri/jsonschema/v5
func (v *SchemaValidator) ValidateJSON(data []byte) *ValidationResult {
	var parsed interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return &ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{Message: fmt.Sprintf("invalid JSON: %s", err.Error())},
			},
		}
	}
	return &ValidationResult{Valid: true}
}

// ValidateRequired checks that required fields are present in a JSON object.
func ValidateRequired(data map[string]interface{}, required []string) *ValidationResult {
	result := &ValidationResult{Valid: true}

	for _, field := range required {
		if _, ok := data[field]; !ok {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Field:   field,
				Message: "required field is missing",
			})
		}
	}

	return result
}

// ValidateCheckoutRequest validates a checkout create request.
func ValidateCheckoutRequest(data map[string]interface{}) *ValidationResult {
	required := []string{"line_items", "currency"}
	result := ValidateRequired(data, required)
	if !result.Valid {
		return result
	}

	// Validate line_items is an array
	lineItems, ok := data["line_items"].([]interface{})
	if !ok {
		return &ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{Field: "line_items", Message: "must be an array"},
			},
		}
	}

	// Validate each line item
	for i, item := range lineItems {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Field:   fmt.Sprintf("line_items[%d]", i),
				Message: "must be an object",
			})
			continue
		}

		// Validate line item required fields
		itemRequired := []string{"item", "quantity"}
		for _, field := range itemRequired {
			if _, ok := itemMap[field]; !ok {
				result.Valid = false
				result.Errors = append(result.Errors, ValidationError{
					Field:   fmt.Sprintf("line_items[%d].%s", i, field),
					Message: "required field is missing",
				})
			}
		}
	}

	// Validate currency
	currency, ok := data["currency"].(string)
	if !ok || len(currency) != 3 {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "currency",
			Message: "must be a 3-letter ISO 4217 currency code",
		})
	}

	return result
}

//go:build ignore

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

// Post-processor for generated UCP models.
// Adds typed enums, cleans up code, and improves Go idioms.
//
// Usage: go run postprocess.go <generated_file.go>

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"regexp"
	"strings"
)

// EnumDefinition represents an enum type to be added
type EnumDefinition struct {
	TypeName string
	BaseType string
	Values   []EnumValue
	Comment  string
}

type EnumValue struct {
	Name    string
	Value   string
	Comment string
}

// Enums to add based on UCP specification
var enumDefinitions = []EnumDefinition{
	{
		TypeName: "CheckoutStatus",
		BaseType: "string",
		Comment:  "CheckoutStatus represents the status of a checkout session.",
		Values: []EnumValue{
			{Name: "CheckoutStatusOpen", Value: "open", Comment: "Checkout is open and accepting updates"},
			{Name: "CheckoutStatusCompleted", Value: "completed", Comment: "Checkout has been completed successfully"},
			{Name: "CheckoutStatusExpired", Value: "expired", Comment: "Checkout session has expired"},
			{Name: "CheckoutStatusRequiresEscalation", Value: "requires_escalation", Comment: "Checkout requires user action via continue_url"},
		},
	},
	{
		TypeName: "FulfillmentMethodType",
		BaseType: "string",
		Comment:  "FulfillmentMethodType represents the type of fulfillment method.",
		Values: []EnumValue{
			{Name: "FulfillmentMethodTypeShipping", Value: "shipping", Comment: "Items shipped to buyer's address"},
			{Name: "FulfillmentMethodTypePickup", Value: "pickup", Comment: "Buyer picks up items at location"},
			{Name: "FulfillmentMethodTypeDigital", Value: "digital", Comment: "Digital/electronic delivery"},
		},
	},
	{
		TypeName: "OrderStatus",
		BaseType: "string",
		Comment:  "OrderStatus represents the status of an order.",
		Values: []EnumValue{
			{Name: "OrderStatusPending", Value: "pending", Comment: "Order is pending processing"},
			{Name: "OrderStatusConfirmed", Value: "confirmed", Comment: "Order has been confirmed"},
			{Name: "OrderStatusProcessing", Value: "processing", Comment: "Order is being processed"},
			{Name: "OrderStatusShipped", Value: "shipped", Comment: "Order has been shipped"},
			{Name: "OrderStatusDelivered", Value: "delivered", Comment: "Order has been delivered"},
			{Name: "OrderStatusCancelled", Value: "cancelled", Comment: "Order has been cancelled"},
		},
	},
	{
		TypeName: "OrderLineItemStatus",
		BaseType: "string",
		Comment:  "OrderLineItemStatus represents the status of a line item in an order.",
		Values: []EnumValue{
			{Name: "OrderLineItemStatusPending", Value: "pending", Comment: "Line item is pending"},
			{Name: "OrderLineItemStatusConfirmed", Value: "confirmed", Comment: "Line item is confirmed"},
			{Name: "OrderLineItemStatusCancelled", Value: "cancelled", Comment: "Line item is cancelled"},
		},
	},
	{
		TypeName: "MessageType",
		BaseType: "string",
		Comment:  "MessageType represents the type/severity of a message.",
		Values: []EnumValue{
			{Name: "MessageTypeError", Value: "error", Comment: "Error message - action required"},
			{Name: "MessageTypeWarning", Value: "warning", Comment: "Warning message - may need attention"},
			{Name: "MessageTypeInfo", Value: "info", Comment: "Informational message"},
		},
	},
	{
		TypeName: "LinkType",
		BaseType: "string",
		Comment:  "LinkType represents the type of a link.",
		Values: []EnumValue{
			{Name: "LinkTypeTermsOfService", Value: "terms_of_service", Comment: "Terms of service link"},
			{Name: "LinkTypePrivacyPolicy", Value: "privacy_policy", Comment: "Privacy policy link"},
			{Name: "LinkTypeReturnPolicy", Value: "return_policy", Comment: "Return policy link"},
			{Name: "LinkTypeShippingPolicy", Value: "shipping_policy", Comment: "Shipping policy link"},
		},
	},
	{
		TypeName: "AllocationMethod",
		BaseType: "string",
		Comment:  "AllocationMethod represents how a discount is allocated across items.",
		Values: []EnumValue{
			{Name: "AllocationMethodAcross", Value: "across", Comment: "Distributed across all items"},
			{Name: "AllocationMethodEach", Value: "each", Comment: "Applied to each item individually"},
		},
	},
	{
		TypeName: "PaymentStatus",
		BaseType: "string",
		Comment:  "PaymentStatus represents the status of a payment.",
		Values: []EnumValue{
			{Name: "PaymentStatusPending", Value: "pending", Comment: "Payment is pending"},
			{Name: "PaymentStatusAuthorized", Value: "authorized", Comment: "Payment has been authorized"},
			{Name: "PaymentStatusCaptured", Value: "captured", Comment: "Payment has been captured"},
			{Name: "PaymentStatusFailed", Value: "failed", Comment: "Payment has failed"},
			{Name: "PaymentStatusRefunded", Value: "refunded", Comment: "Payment has been refunded"},
		},
	},
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run postprocess.go <generated_file.go>")
		os.Exit(1)
	}

	filename := os.Args[1]

	// Read the generated file
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Process the content
	processed := processContent(string(content))

	// Format the output
	formatted, err := format.Source([]byte(processed))
	if err != nil {
		// If formatting fails, write unformatted
		fmt.Printf("Warning: gofmt failed: %v\n", err)
		formatted = []byte(processed)
	}

	// Write back
	if err := os.WriteFile(filename, formatted, 0644); err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Post-processed: %s\n", filename)
}

func processContent(content string) string {
	var buf bytes.Buffer

	// Find the end of the package declaration and imports
	packageEnd := strings.Index(content, "\n\ntype ")
	if packageEnd == -1 {
		packageEnd = strings.Index(content, "\ntype ")
	}
	if packageEnd == -1 {
		// No types found, just add enums at end
		buf.WriteString(content)
		buf.WriteString("\n")
		buf.WriteString(generateEnums())
		return buf.String()
	}

	// Write everything up to the first type
	buf.WriteString(content[:packageEnd])
	buf.WriteString("\n\n")

	// Add enum definitions
	buf.WriteString(generateEnums())
	buf.WriteString("\n")

	// Write the rest of the content
	buf.WriteString(content[packageEnd:])

	result := buf.String()

	// Clean up: remove mapstructure and yaml tags if present
	result = cleanupTags(result)

	// Clean up: fix common naming issues
	result = fixNaming(result)

	return result
}

func generateEnums() string {
	var buf bytes.Buffer

	buf.WriteString("// ============================================\n")
	buf.WriteString("// Enum Types\n")
	buf.WriteString("// ============================================\n\n")

	for _, enum := range enumDefinitions {
		// Type definition
		buf.WriteString(fmt.Sprintf("// %s\n", enum.Comment))
		buf.WriteString(fmt.Sprintf("type %s %s\n\n", enum.TypeName, enum.BaseType))

		// Constants
		buf.WriteString(fmt.Sprintf("// %s values\n", enum.TypeName))
		buf.WriteString("const (\n")
		for _, v := range enum.Values {
			buf.WriteString(fmt.Sprintf("\t// %s\n", v.Comment))
			buf.WriteString(fmt.Sprintf("\t%s %s = %q\n", v.Name, enum.TypeName, v.Value))
		}
		buf.WriteString(")\n\n")
	}

	return buf.String()
}

func cleanupTags(content string) string {
	// Remove yaml and mapstructure tags, keep only json
	re := regexp.MustCompile(` yaml:"[^"]*"`)
	content = re.ReplaceAllString(content, "")

	re = regexp.MustCompile(` mapstructure:"[^"]*"`)
	content = re.ReplaceAllString(content, "")

	// Clean up extra spaces in struct tags
	re = regexp.MustCompile(`json:"([^"]+)"\s+\x60`)
	content = re.ReplaceAllString(content, `json:"$1"` + "`")

	return content
}

func fixNaming(content string) string {
	// Fix common capitalization issues
	replacements := map[string]string{
		"Url ":  "URL ",
		"Url\n": "URL\n",
		"Uri ":  "URI ",
		"Uri\n": "URI\n",
		"Api":   "API",
		" Id ":  " ID ",
		" Id\n": " ID\n",
		"(Id ":  "(ID ",
	}

	for old, new := range replacements {
		content = strings.ReplaceAll(content, old, new)
	}

	return content
}

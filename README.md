<!--
   Copyright 2026 UCP Authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
-->

<p align="center">
  <h1 align="center">UCP Go SDK</h1>
</p>

<p align="center">
  <b>Official Go library for the Universal Commerce Protocol (UCP).</b>
</p>

## Overview

This repository contains the Go SDK for the
[Universal Commerce Protocol (UCP)](https://ucp.dev). It provides Go types,
validation, and client/server utilities for building UCP-compliant applications.

## Features

- **Models**: Go structs for all UCP schemas (checkout, order, payment, fulfillment, etc.)
- **REST Client**: Typed HTTP client for consuming UCP APIs from platforms/agents
- **Server Helpers**: HTTP handlers and middleware for implementing UCP endpoints
- **Validation**: JSON Schema validation and capability negotiation
- **Extensions**: Extended types for UCP extensions (fulfillment, discounts, buyer consent)

## Installation

```bash
go get github.com/dhananjay2021/ucp-go-sdk
```

## Quick Start

### Platform Client (Consuming UCP APIs)

Use the client package to interact with UCP-compliant merchants:

```go
package main

import (
    "context"
    "log"

    "github.com/dhananjay2021/ucp-go-sdk/client"
    "github.com/dhananjay2021/ucp-go-sdk/extensions"
    "github.com/dhananjay2021/ucp-go-sdk/models"
)

func main() {
    // Create a UCP client
    ucpClient := client.NewClient("https://merchant.example.com",
        client.WithAPIKey("your-api-key"),
    )

    ctx := context.Background()

    // Discover merchant capabilities
    profile, err := ucpClient.FetchProfile(ctx)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Merchant UCP version: %s", profile.UCP.Version)

    // Check for required capabilities
    if !client.HasCapability(profile, client.CapabilityCheckout) {
        log.Fatal("Merchant does not support checkout")
    }

    // Create a checkout session
    checkout, err := ucpClient.CreateCheckout(ctx, &extensions.ExtendedCheckoutCreateRequest{
        LineItems: []models.LineItemCreateRequest{
            {
                Item: models.ItemCreateRequest{
                    Name:  "Product Name",
                    Price: "29.99",
                },
                Quantity: 1,
            },
        },
        Currency: "USD",
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Created checkout: %s (status: %s)", checkout.ID, checkout.Status)
}
```

### Business Server (Implementing UCP APIs)

Use the server package to implement UCP-compliant endpoints:

```go
package main

import (
    "log"
    "net/http"

    "github.com/dhananjay2021/ucp-go-sdk/client"
    "github.com/dhananjay2021/ucp-go-sdk/extensions"
    "github.com/dhananjay2021/ucp-go-sdk/models"
    "github.com/dhananjay2021/ucp-go-sdk/server"
)

func main() {
    // Configure the server
    config := server.Config{
        Version: "2026-01-11",
        Capabilities: []models.CapabilityDiscovery{
            {
                CapabilityBase: models.CapabilityBase{
                    Name:    client.CapabilityCheckout,
                    Version: "2026-01-11",
                    Spec:    "https://ucp.dev/specification/checkout",
                    Schema:  "https://ucp.dev/schemas/shopping/checkout.json",
                },
            },
        },
        Services: models.Services{
            client.ServiceShopping: models.UCPService{
                Version: "2026-01-11",
                Spec:    "https://ucp.dev/specification/shopping",
                Rest: &models.RestTransport{
                    Schema:   "https://ucp.dev/schemas/services/shopping/rest.openapi.json",
                    Endpoint: "http://localhost:8080",
                },
            },
        },
    }

    // Create server and register handlers
    srv := server.NewServer(config)
    srv.HandleCreateCheckout(handleCreateCheckout)
    srv.HandleUpdateCheckout(handleUpdateCheckout)
    srv.HandleCompleteCheckout(handleCompleteCheckout)

    // Apply middleware
    handler := server.Chain(srv,
        server.LoggingMiddleware,
        server.CORSMiddleware([]string{"*"}),
    )

    log.Println("Starting UCP server on :8080")
    http.ListenAndServe(":8080", handler)
}

func handleCreateCheckout(r *http.Request, req *extensions.ExtendedCheckoutCreateRequest) (*extensions.ExtendedCheckoutResponse, error) {
    // Implement your checkout logic here
    return &extensions.ExtendedCheckoutResponse{
        ID:     "chk-123",
        Status: models.CheckoutStatusIncomplete,
        // ...
    }, nil
}
```

## Package Structure

```
go-sdk/
├── models/          # Go types for all UCP schemas
├── client/          # REST client for consuming UCP APIs
├── server/          # HTTP handlers for implementing UCP endpoints
├── validation/      # JSON Schema validation and capability negotiation
├── extensions/      # Extended types for UCP extensions
├── internal/        # Internal utilities
└── examples/        # Example implementations
    ├── business_server/   # Example merchant server
    └── platform_client/   # Example platform client
```

## Models

The `models` package contains Go structs for all UCP schemas:

- **UCP Core**: `Version`, `CapabilityBase`, `DiscoveryProfile`, `UCPProfile`
- **Checkout**: `CheckoutCreateRequest`, `CheckoutUpdateRequest`, `CheckoutResponse`
- **Payment**: `PaymentResponse`, `PaymentHandlerResponse`, `CardCredential`
- **Fulfillment**: `FulfillmentRequest`, `FulfillmentResponse`, `ShippingDestination`
- **Order**: `Order`, `OrderLineItem`, `Adjustment`
- **Discount**: `DiscountCreateRequest`, `DiscountResponse`
- **Buyer Consent**: `BuyerConsentCreateRequest`, `BuyerConsentResponse`

## Client Package

The `client` package provides a REST client for platforms and agents:

```go
// Create client with options
c := client.NewClient(baseURL,
    client.WithAPIKey("key"),
    client.WithAccessToken("token"),
    client.WithTimeout(30*time.Second),
)

// Discovery
profile, _ := c.FetchProfile(ctx)

// Checkout operations
checkout, _ := c.CreateCheckout(ctx, req)
checkout, _ := c.GetCheckout(ctx, id)
checkout, _ := c.UpdateCheckout(ctx, id, updateReq)
checkout, _ := c.CompleteCheckout(ctx, id)
checkout, _ := c.CancelCheckout(ctx, id)

// Order operations
order, _ := c.GetOrder(ctx, id)
```

## Server Package

The `server` package helps implement UCP-compliant endpoints:

```go
// Create server
srv := server.NewServer(config)

// Register handlers
srv.HandleCreateCheckout(handler)
srv.HandleGetCheckout(handler)
srv.HandleUpdateCheckout(handler)
srv.HandleCompleteCheckout(handler)
srv.HandleCancelCheckout(handler)
srv.HandleGetOrder(handler)

// Available middleware
server.LoggingMiddleware
server.CORSMiddleware(allowedOrigins)
server.APIKeyMiddleware(validKeys)
server.BearerTokenMiddleware(validator)
server.RequestIDMiddleware

// Response helpers
server.WriteJSON(w, statusCode, data)
server.WriteError(w, statusCode, code, message)
```

## Validation Package

The `validation` package provides capability negotiation:

```go
// Create negotiator with platform capabilities
negotiator := validation.NewCapabilityNegotiator(platformCaps)

// Negotiate with a business profile
result := negotiator.Negotiate(businessProfile, requiredCaps)

if result.Success {
    // Use result.CommonCapabilities
    // Use result.NegotiatedVersion
}
```

## Extensions Package

The `extensions` package provides extended types that combine base schemas with extensions:

```go
// Extended checkout with fulfillment, discounts, and buyer consent
type ExtendedCheckoutResponse struct {
    // Base checkout fields...
    Fulfillment *FulfillmentResponse
    Discounts   []DiscountResponse
    Buyer       *BuyerConsentResponse
}
```

## Running Examples

### Business Server

```bash
cd examples/business_server
go run main.go
```

The server starts on port 8080 with a discovery endpoint at `/.well-known/ucp`.

### Platform Client

```bash
# Start the business server first, then:
cd examples/platform_client
go run main.go
```

## Development

### Prerequisites

- Go 1.22 or later

### Building

```bash
go build ./...
```

### Testing

```bash
go test ./...
```

### Generating Models

The models can be regenerated from UCP JSON schemas:

```bash
# Ensure the UCP spec repository is available
./generate_models.sh
```

## Contributing

We welcome community contributions. See our [Contribution Guide](https://github.com/Universal-Commerce-Protocol/ucp/blob/main/CONTRIBUTING.md) for details.

## License

UCP is an open-source project under the [Apache License 2.0](LICENSE).

# Getting Started with UCP in Go

**Build interoperable commerce experiences with the Universal Commerce Protocol**

## TL;DR

UCP (Universal Commerce Protocol) is an open standard for e-commerce APIs. This guide shows you how to use the Go SDK to either consume merchant APIs as a platform, or expose UCP-compliant endpoints as a business. Install with `go get github.com/dhananjay2021/ucp-go-sdk@latest` and follow along.

• • •

The e-commerce landscape is fragmented. Every platform, every merchant, and every payment provider speaks a different language. If you've ever integrated with multiple shopping APIs, you know the pain: inconsistent data models, varying authentication schemes, and endless edge cases.

Enter the **Universal Commerce Protocol (UCP)** — an open standard that brings order to this chaos.

In this guide, I'll show you how to get started with UCP in Go, whether you're building a platform that consumes merchant APIs or a business that needs to expose UCP-compliant endpoints.

## What is UCP?

[UCP](https://ucp.dev) is an open protocol that standardizes how commerce systems communicate. Co-developed by **Google** and **Shopify**, and supported by major retailers like Target, Walmart, Etsy, and Wayfair, UCP is quickly becoming the standard for agentic commerce.

Think of it as "HTTP for shopping" — a common language for:

- **Checkout flows** — Cart creation, buyer info, payment, and order confirmation
- **Fulfillment** — Shipping options, delivery tracking
- **Payments** — Tokenization, payment handlers, refunds
- **Discovery** — Capability negotiation between systems

The protocol is designed for the AI agent era, where autonomous systems need to discover and interact with merchants programmatically.

## Installing the Go SDK

```bash
go get github.com/dhananjay2021/ucp-go-sdk@latest
```

The SDK provides:
- **Type-safe models** for all UCP schemas
- **HTTP client** for consuming UCP APIs
- **Server helpers** for implementing UCP endpoints
- **Middleware** for logging, CORS, authentication

## Part 1: Consuming a UCP API (Platform Client)

Let's build a client that discovers a merchant's capabilities, creates a checkout, and completes a purchase.

### Step 1: Create a Client and Discover Capabilities

```go
package main

import (
    "context"
    "log"
    "github.com/dhananjay2021/ucp-go-sdk/client"
)

func main() {
    ucpClient := client.NewClient("https://merchant.example.com",
        client.WithAPIKey("your-api-key"),
    )

    ctx := context.Background()

    // Discover what the merchant supports
    profile, err := ucpClient.FetchProfile(ctx)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Merchant UCP Version: %s", profile.UCP.Version)
    
    for _, cap := range profile.UCP.Capabilities {
        log.Printf("  ✓ %s (v%s)", cap.Name, cap.Version)
    }

    // Check for required capability before proceeding
    if !client.HasCapability(profile, client.CapabilityCheckout) {
        log.Fatal("Merchant doesn't support checkout")
    }
}
```

UCP's discovery mechanism means your code can adapt to different merchants dynamically — no hardcoded assumptions.

### Step 2: Create a Checkout Session

In UCP, the platform sends only **item IDs** — the merchant looks up product details (title, price, etc.) and returns them in the response. This keeps the platform lightweight and ensures merchants control their catalog.

```go
import (
    "github.com/dhananjay2021/ucp-go-sdk/extensions"
    "github.com/dhananjay2021/ucp-go-sdk/models"
)

checkout, err := ucpClient.CreateCheckout(ctx, &extensions.ExtendedCheckoutCreateRequest{
    LineItems: []models.LineItemCreateRequest{
        {
            Item:     models.ItemCreateRequest{ID: "SKU-001"},
            Quantity: 1,
        },
        {
            Item:     models.ItemCreateRequest{ID: "SKU-002"},
            Quantity: 2,
        },
    },
    Currency: "USD",
    Payment:  models.PaymentCreateRequest{},
})
if err != nil {
    log.Fatal(err)
}

log.Printf("Checkout created: %s (Status: %s)", checkout.ID, checkout.Status)

// Amounts are in cents (minor currency units)
for _, total := range checkout.Totals {
    log.Printf("  %s: %d cents", total.Type, total.Amount)
}
```

### Step 3: Update with Buyer Information

```go
checkout, err = ucpClient.UpdateCheckout(ctx, checkout.ID, 
    &extensions.ExtendedCheckoutUpdateRequest{
        ID:       checkout.ID,
        Currency: "USD",
        LineItems: []models.LineItemUpdateRequest{
            {Item: models.ItemUpdateRequest{ID: "SKU-001"}, Quantity: 1},
            {Item: models.ItemUpdateRequest{ID: "SKU-002"}, Quantity: 2},
        },
        Payment: models.PaymentUpdateRequest{},
        Buyer: &models.BuyerWithConsentUpdateRequest{
            Email:       "jane@example.com",
            FirstName:   "Jane",
            LastName:    "Doe",
            PhoneNumber: "+1-555-123-4567",
        },
    })
```

### Step 4: Add Payment and Complete

```go
// Add payment instrument
checkout, err = ucpClient.UpdateCheckout(ctx, checkout.ID, 
    &extensions.ExtendedCheckoutUpdateRequest{
        ID:       checkout.ID,
        Currency: "USD",
        LineItems: []models.LineItemUpdateRequest{
            {Item: models.ItemUpdateRequest{ID: "SKU-001"}, Quantity: 1},
        },
        Payment: models.PaymentUpdateRequest{
            SelectedInstrumentID: "pi-001",
            Instruments: []models.PaymentInstrument{
                {
                    ID:         "pi-001",
                    HandlerID:  "default",
                    Type:       models.PaymentInstrumentTypeCard,
                    Brand:      "visa",
                    LastDigits: "4242",
                },
            },
        },
    })

// Complete the checkout
if checkout.Status == models.CheckoutStatusReadyForComplete {
    checkout, err = ucpClient.CompleteCheckout(ctx, checkout.ID)
    log.Printf("Order confirmed: %s", checkout.Order.ID)
}
```

## Part 2: Implementing a UCP Server (Business/Merchant)

Now let's flip the script. If you're a merchant, here's how to expose UCP-compliant endpoints.

### Step 1: Configure Your Server

```go
package main

import (
    "log"
    "net/http"
    "github.com/dhananjay2021/ucp-go-sdk/client"
    "github.com/dhananjay2021/ucp-go-sdk/models"
    "github.com/dhananjay2021/ucp-go-sdk/server"
)

func main() {
    config := server.Config{
        Version: "2026-01-11",
        Capabilities: []models.CapabilityDiscovery{
            {
                CapabilityBase: models.CapabilityBase{
                    Name:    client.CapabilityCheckout,
                    Version: "2026-01-11",
                },
            },
        },
        Services: models.Services{
            client.ServiceShopping: models.UCPService{
                Version: "2026-01-11",
                Rest:    &models.RestTransport{Endpoint: "https://your-store.com"},
            },
        },
    }

    srv := server.NewServer(config)
    srv.HandleCreateCheckout(handleCreateCheckout)
    srv.HandleUpdateCheckout(handleUpdateCheckout)
    srv.HandleCompleteCheckout(handleCompleteCheckout)

    handler := server.Chain(srv,
        server.LoggingMiddleware,
        server.CORSMiddleware([]string{"*"}),
    )

    log.Println("UCP server running on :8080")
    http.ListenAndServe(":8080", handler)
}
```

### Step 2: Implement Handlers

The merchant looks up products by ID and returns full item details with pricing:

```go
import "github.com/dhananjay2021/ucp-go-sdk/extensions"

// Product catalog (in real app, this would be a database)
var catalog = map[string]struct {
    Title string
    Price int // cents
}{
    "SKU-001": {Title: "Wireless Headphones", Price: 14999},
    "SKU-002": {Title: "Phone Case", Price: 2999},
}

func handleCreateCheckout(r *http.Request, req *extensions.ExtendedCheckoutCreateRequest) (*extensions.ExtendedCheckoutResponse, error) {
    var subtotal int
    lineItems := make([]models.LineItemResponse, len(req.LineItems))

    for i, li := range req.LineItems {
        product, ok := catalog[li.Item.ID]
        if !ok {
            return nil, server.BadRequestError("unknown product: " + li.Item.ID)
        }
        
        itemTotal := product.Price * li.Quantity
        subtotal += itemTotal
        
        lineItems[i] = models.LineItemResponse{
            ID: generateUniqueID(),
            Item: models.ItemResponse{
                ID:    li.Item.ID,
                Title: product.Title,
                Price: product.Price,
            },
            Quantity: li.Quantity,
            Totals: []models.TotalResponse{
                {Type: models.TotalTypeSubtotal, Amount: itemTotal},
            },
        }
    }

    tax := subtotal * 875 / 10000  // 8.75% tax

    return &extensions.ExtendedCheckoutResponse{
        ID:        generateUniqueID(),
        Status:    models.CheckoutStatusIncomplete,
        Currency:  req.Currency,
        LineItems: lineItems,
        Totals: []models.TotalResponse{
            {Type: models.TotalTypeSubtotal, Amount: subtotal},
            {Type: models.TotalTypeTax, Amount: tax},
            {Type: models.TotalTypeTotal, Amount: subtotal + tax},
        },
        Links: []models.Link{
            {Type: "terms_of_service", URL: "https://example.com/terms"},
            {Type: "privacy_policy", URL: "https://example.com/privacy"},
        },
        Payment: models.PaymentResponse{
            Handlers: []models.PaymentHandlerResponse{
                {
                    ID:      "default",
                    Name:    "dev.ucp.tokenization",
                    Version: "2026-01-11",
                    Spec:    "https://ucp.dev/handlers/tokenization",
                    ConfigSchema:      "https://ucp.dev/handlers/tokenization/config.json",
                    InstrumentSchemas: []string{"https://ucp.dev/schemas/card_payment_instrument.json"},
                    Config:            map[string]interface{}{},
                },
            },
        },
    }, nil
}
```

### Step 3: Handle Errors Gracefully

The SDK provides typed error helpers:

```go
func handleGetCheckout(r *http.Request, id string) (*extensions.ExtendedCheckoutResponse, error) {
    checkout, found := store.Get(id)
    if !found {
        return nil, server.NotFoundError("checkout not found")
    }
    return checkout, nil
}
```

## Testing Your Integration

Run the example server:

```bash
cd examples/business_server && go run main.go
```

In another terminal, test the client:

```bash
cd examples/platform_client && go run main.go
```

Or hit the discovery endpoint:

```bash
curl http://localhost:8080/.well-known/ucp | jq
```

## Why UCP Matters

- **Interoperability** — Write once, integrate with any UCP-compliant merchant
- **AI-Ready** — Designed for autonomous agents to discover and transact
- **Type Safety** — The Go SDK gives you compile-time guarantees
- **Extensible** — Capabilities like fulfillment and discounts are opt-in

## What's Next?

- 📖 Read the [UCP Specification](https://ucp.dev/specification/overview)
- 💻 Explore the [SDK on GitHub](https://github.com/dhananjay2021/ucp-go-sdk)
- 📦 Check it out on [pkg.go.dev](https://pkg.go.dev/github.com/dhananjay2021/ucp-go-sdk)

## Further Reading

Want to dive deeper into UCP? Check out these official resources:

- [Building the Universal Commerce Protocol](https://shopify.engineering/ucp) — Shopify Engineering's deep dive into UCP's architecture, capability negotiation, and why they built it
- [Under the Hood: Universal Commerce Protocol](https://developers.googleblog.com/en/under-the-hood-universal-commerce-protocol-ucp/) — Google's perspective on UCP and how it powers agentic commerce
- [Agentic Commerce: AI Tools and Protocol for Retailers](https://blog.google/products/ads-commerce/agentic-commerce-ai-tools-protocol-retailers-platforms/) — Google's announcement on how UCP enables AI-powered shopping experiences

• • •

*Have questions or feedback? [Open an issue on GitHub](https://github.com/dhananjay2021/ucp-go-sdk/issues) or connect with me on [LinkedIn](https://linkedin.com/in/YOUR-LINKEDIN-USERNAME).*

**Happy building!**

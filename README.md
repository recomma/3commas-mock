# 3Commas Mock Server

A mock HTTP server for testing 3Commas API integrations, following the same patterns as [`github.com/recomma/hyperliquid-mock`](https://github.com/recomma/hyperliquid-mock).

## Features

- **3 Core Endpoints**: ListBots, ListDeals, GetDeal
- **Programmatic State Management**: Add/update/remove bots and deals
- **Rich Bot Events**: Full event structure with action, coin, type, status, price, size, etc.
- **Error Simulation**: Rate limiting, 404s, and custom errors
- **Thread-Safe**: Safe for concurrent access
- **Easy Testing**: Simple httptest-based server for integration tests

## Installation

```bash
go get github.com/recomma/3commas-mock
```

## Quick Start

```go
package mytest

import (
    "testing"
    "github.com/recomma/3commas-mock/server"
)

func TestMyIntegration(t *testing.T) {
    // Create mock server
    mockServer := server.NewTestServer(t)
    defer mockServer.Close()

    // Add a bot
    mockServer.AddBot(server.Bot{
        ID:        1,
        Name:      "Test Bot",
        Enabled:   true,
        AccountID: 12345,
    })

    // Add a deal with events
    mockServer.AddDeal(1, server.Deal{
        ID:           101,
        BotID:        1,
        Pair:         "USDT_BTC",
        Status:       "active",
        ToCurrency:   "BTC",
        FromCurrency: "USDT",
        Events: []server.BotEvent{
            {
                CreatedAt:     "2024-01-15T10:30:00.000Z",
                Action:        "place",
                Coin:          "BTC",
                Type:          "buy",
                Status:        "active",
                Price:         "50000.0",
                Size:          "0.0002",
                OrderType:     "base",
                OrderSize:     1,
                OrderPosition: 1,
                IsMarket:      false,
            },
        },
    })

    // Use mockServer.URL() in your client
    // client := NewMyClient(mockServer.URL())
}
```

## API Endpoints

The mock server implements 3 endpoints from the 3Commas API:

### GET /ver1/bots

Lists bots with optional filtering.

**Query Parameters:**
- `scope` (optional): `enabled` or `disabled`

**Example:**
```bash
curl http://localhost/ver1/bots?scope=enabled
```

### GET /ver1/deals

Lists deals with optional filtering.

**Query Parameters:**
- `bot_id` (optional): Filter by bot ID
- `scope` (optional): Filter by status (e.g., `active`, `finished`)

**Example:**
```bash
curl http://localhost/ver1/deals?bot_id=1&scope=active
```

### GET /ver1/deals/{deal_id}/show

Get a specific deal by ID, including full bot_events history.

**Example:**
```bash
curl http://localhost/ver1/deals/101/show
```

## State Management

### Bots

```go
// Add a bot
mockServer.AddBot(server.Bot{
    ID:        1,
    Name:      "My Bot",
    Enabled:   true,
    AccountID: 123,
})

// Get a bot
bot, ok := mockServer.GetBot(1)

// Update a bot
enabled := false
mockServer.UpdateBot(1, server.BotUpdate{
    Enabled: &enabled,
})

// Remove a bot (also removes all its deals)
mockServer.RemoveBot(1)

// Get all bots
bots := mockServer.GetAllBots()
```

### Deals

```go
// Add a deal
mockServer.AddDeal(1, server.Deal{
    ID:     101,
    BotID:  1,
    Status: "active",
    Events: []server.BotEvent{},
})

// Get a deal
deal, ok := mockServer.GetDealByID(101)

// Update a deal
status := "completed"
mockServer.UpdateDeal(101, server.DealUpdate{
    Status: &status,
})

// Add an event to a deal
mockServer.AddBotEvent(101, server.BotEvent{
    CreatedAt:     "2024-01-15T10:32:00.000Z",
    Action:        "place",
    Coin:          "BTC",
    Type:          "buy",
    Status:        "active",
    Price:         "48750.0",
    Size:          "0.0004",
    OrderType:     "safety",
    OrderSize:     2,
    OrderPosition: 1,
    IsMarket:      false,
})

// Remove a deal
mockServer.RemoveDeal(101)

// Get all deals for a bot
deals := mockServer.GetBotDeals(1)

// Get all deals
allDeals := mockServer.GetAllDeals()
```

## Error Simulation

```go
// Enable rate limiting (429 responses)
mockServer.SetRateLimitError(true, 60) // 60 second retry-after

// Set bot-specific error
mockServer.SetBotError(1, fmt.Errorf("bot error"))

// Set deal-specific error
mockServer.SetDealError(101, fmt.Errorf("deal error"))

// Clear all errors
mockServer.ClearErrors()

// Reset everything (clears state and errors)
mockServer.Reset()
```

## Bot Event Structure

The mock server uses a rich bot event structure for detailed testing:

```go
type BotEvent struct {
    CreatedAt     string `json:"created_at"`      // ISO 8601 timestamp
    Action        string `json:"action"`          // "place", "cancel", "cancelled", "modify"
    Coin          string `json:"coin"`            // "BTC", "ETH", etc.
    Type          string `json:"type"`            // "buy", "sell"
    Status        string `json:"status"`          // "active", "filled", "cancelled"
    Price         string `json:"price"`           // Decimal string
    Size          string `json:"size"`            // Decimal string
    OrderType     string `json:"order_type"`      // "base", "safety", "take_profit"
    OrderSize     int    `json:"order_size"`      // Order size category
    OrderPosition int    `json:"order_position"`  // Safety order position
    IsMarket      bool   `json:"is_market"`       // Market vs limit order
}
```

This structure allows for precise simulation of order lifecycle events in tests.

## Architecture

The mock server is built using:

- **oapi-codegen**: Generates server interface from OpenAPI 3.0 spec
- **httptest**: Provides lightweight test server
- **Custom types**: Rich event structures for testing flexibility

The implementation filters the full 3Commas OpenAPI spec to only generate the 3 required endpoints, keeping the codebase minimal while maintaining type safety.

## Development

### Generate Code

Code is generated from the [3commas-openapi](https://github.com/recomma/3commas-openapi) submodule:

```bash
go generate ./...
```

### Run Tests

```bash
go test ./server -v
```

## License

Apache 2.0

## See Also

- [3commas-sdk-go](https://github.com/recomma/3commas-sdk-go) - Go SDK for 3Commas API
- [3commas-openapi](https://github.com/recomma/3commas-openapi) - OpenAPI spec for 3Commas
- [hyperliquid-mock](https://github.com/recomma/hyperliquid-mock) - Similar mock pattern for Hyperliquid

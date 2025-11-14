# 3Commas Mock Server Specification

## Overview

This document specifies the requirements for a mock HTTP server that simulates the 3Commas API for testing purposes. The mock server should be implemented as a standalone Go package (`github.com/recomma/3commas-mock`) following the same patterns as `github.com/recomma/hyperliquid-mock`.

## Purpose

The mock server enables end-to-end testing of the Recomma application by:
- Simulating 3Commas API responses without making real API calls
- Allowing programmatic control of bot and deal states
- Supporting error simulation and edge case testing
- Enabling integration tests to run in CI/CD without external dependencies

## Required API Endpoints

### 1. List Bots
**Endpoint:** `GET /public/api/ver1/bots`

**Query Parameters:**
- `scope` (string, optional): Filter by bot status
  - `enabled` - Only return enabled bots (most commonly used)
  - `disabled` - Only return disabled bots
  - If omitted, return all bots

**Response:** `200 OK`
```json
[
  {
    "id": 1,
    "account_id": 12345,
    "account_name": "Test Account",
    "is_enabled": true,
    "max_safety_orders": 5,
    "active_safety_orders_count": 2,
    "pairs": ["USDT_BTC", "USDT_ETH"],
    "strategy_list": [{"strategy": "long"}],
    "max_active_deals": 3,
    "active_deals_count": 1,
    "deletable": true,
    "created_at": "2024-01-01T00:00:00.000Z",
    "updated_at": "2024-01-01T00:00:00.000Z",
    "base_order_volume": 10.0,
    "safety_order_volume": 10.0,
    "safety_order_step_percentage": 2.5,
    "take_profit": 1.5,
    "martingale_volume_coefficient": 1.0,
    "martingale_step_coefficient": 1.0,
    "max_deal_funds": 100.0,
    "name": "Test Bot"
  }
]
```

**Error Response:** `429 Too Many Requests`
```json
{
  "error": "rate limit exceeded",
  "error_description": "You have exceeded the rate limit. Please try again later."
}
```

### 2. List Deals
**Endpoint:** `GET /public/api/ver1/deals`

**Query Parameters:**
- `bot_id` (integer, required): Bot ID to list deals for
- `scope` (string, optional): Filter by deal status (commonly `active`)
- `limit` (integer, optional): Number of deals per page (default: 100)
- `offset` (integer, optional): Pagination offset (default: 0)

**Response:** `200 OK`
```json
[
  {
    "id": 101,
    "type": "Deal::Base",
    "bot_id": 1,
    "bot_name": "Test Bot",
    "pair": "USDT_BTC",
    "status": "active",
    "created_at": "2024-01-15T10:30:00.000Z",
    "updated_at": "2024-01-15T10:35:00.000Z",
    "closed_at": null,
    "finished?": false,
    "current_active_safety_orders_count": 2,
    "completed_safety_orders_count": 0,
    "completed_manual_safety_orders_count": 0,
    "cancellable?": true,
    "panic_sellable?": true,
    "trailing_enabled": false,
    "tsl_enabled": false,
    "stop_loss_timeout_enabled": false,
    "active_manual_safety_orders": 0,
    "bought_amount": "0.001",
    "bought_volume": "50.0",
    "bought_average_price": "50000.0",
    "base_order_volume": "10.0",
    "safety_order_volume": "20.0",
    "deal_has_error": false,
    "from_currency": "USDT",
    "to_currency": "BTC",
    "take_profit": "1.5",
    "final_profit": "0.0",
    "martingale_coefficient": "1.0",
    "martingale_volume_coefficient": "1.0",
    "martingale_step_coefficient": "1.0",
    "max_safety_orders": 5,
    "active_safety_orders": 2,
    "max_deal_funds": "100.0",
    "actual_profit_percentage": "0.0",
    "actual_profit": "0.0"
  }
]
```

**Error Response:** `404 Not Found`
```json
{
  "error": "bot not found"
}
```

### 3. Get Deal by ID
**Endpoint:** `GET /public/api/ver1/deals/{deal_id}/show`

**Path Parameters:**
- `deal_id` (integer, required): Deal ID

**Response:** `200 OK`
```json
{
  "id": 101,
  "type": "Deal::Base",
  "bot_id": 1,
  "bot_name": "Test Bot",
  "pair": "USDT_BTC",
  "status": "active",
  "created_at": "2024-01-15T10:30:00.000Z",
  "updated_at": "2024-01-15T10:35:00.000Z",
  "from_currency": "USDT",
  "to_currency": "BTC",
  "base_order_volume": "10.0",
  "bought_amount": "0.001",
  "bought_volume": "50.0",
  "bought_average_price": "50000.0",
  "take_profit": "1.5",
  "bot_events": [
    {
      "created_at": "2024-01-15T10:30:00.000Z",
      "action": "place",
      "coin": "BTC",
      "type": "buy",
      "status": "active",
      "price": "50000.0",
      "size": "0.0002",
      "order_type": "base",
      "order_size": 1,
      "order_position": 1,
      "is_market": false
    },
    {
      "created_at": "2024-01-15T10:32:00.000Z",
      "action": "place",
      "coin": "BTC",
      "type": "buy",
      "status": "active",
      "price": "48750.0",
      "size": "0.0004",
      "order_type": "safety",
      "order_size": 2,
      "order_position": 1,
      "is_market": false
    },
    {
      "created_at": "2024-01-15T10:35:00.000Z",
      "action": "place",
      "coin": "BTC",
      "type": "sell",
      "status": "active",
      "price": "50750.0",
      "size": "0.0006",
      "order_type": "take_profit",
      "order_size": 3,
      "order_position": 1,
      "is_market": false
    }
  ]
}
```

**Error Response:** `404 Not Found`
```json
{
  "error": "deal not found"
}
```

## Bot Event Structure

Bot events represent order lifecycle actions. The `bot_events` array in a deal contains the history of all order actions.

### Event Fields

| Field | Type | Description | Example Values |
|-------|------|-------------|----------------|
| `created_at` | string (ISO 8601) | When the event was created | `"2024-01-15T10:30:00.000Z"` |
| `action` | string | Event action type | `"place"`, `"cancel"`, `"cancelled"`, `"modify"` |
| `coin` | string | Cryptocurrency symbol | `"BTC"`, `"ETH"`, `"SOL"` |
| `type` | string | Order side | `"buy"`, `"sell"` |
| `status` | string | Order status | `"active"`, `"filled"`, `"cancelled"` |
| `price` | string | Order price (as decimal string) | `"50000.0"`, `"3000.5"` |
| `size` | string | Order size/quantity (as decimal string) | `"0.001"`, `"1.5"` |
| `order_type` | string | Deal order type | `"base"`, `"safety"`, `"take_profit"` |
| `order_size` | integer | Order size category | `1`, `2`, `3` |
| `order_position` | integer | Safety order position | `1`, `2`, `3` |
| `is_market` | boolean | Whether order is market order | `true`, `false` |

### Event Actions

- `"place"` - New order placed
- `"modify"` - Existing order modified (price or size changed)
- `"cancel"` - Cancel action initiated
- `"cancelled"` - Order was cancelled

### Order Types

- `"base"` - Base order (initial entry order)
- `"safety"` - Safety order (averaging down/up)
- `"take_profit"` - Take profit order (exit order)

### Event Lifecycle Examples

**Base order placed:**
```json
{
  "created_at": "2024-01-15T10:30:00.000Z",
  "action": "place",
  "coin": "BTC",
  "type": "buy",
  "status": "active",
  "price": "50000.0",
  "size": "0.0002",
  "order_type": "base",
  "order_size": 1,
  "order_position": 1,
  "is_market": false
}
```

**Order modification (price changed):**
```json
{
  "created_at": "2024-01-15T10:31:00.000Z",
  "action": "modify",
  "coin": "BTC",
  "type": "buy",
  "status": "active",
  "price": "49800.0",
  "size": "0.0002",
  "order_type": "base",
  "order_size": 1,
  "order_position": 1,
  "is_market": false
}
```

**Order cancellation:**
```json
{
  "created_at": "2024-01-15T10:32:00.000Z",
  "action": "cancel",
  "coin": "BTC",
  "type": "sell",
  "status": "cancelled",
  "price": "50750.0",
  "size": "0.0006",
  "order_type": "take_profit",
  "order_size": 3,
  "order_position": 1,
  "is_market": false
}
```

## Mock Server API

### Test Server Creation

```go
package threecommasmock

import (
    "net/http/httptest"
    "testing"
)

// NewTestServer creates a new 3Commas mock server for testing
func NewTestServer(t *testing.T) *TestServer {
    // Implementation creates httptest.Server with configured handlers
}

type TestServer struct {
    server *httptest.Server
    // Internal state management
}

// URL returns the mock server's base URL
func (ts *TestServer) URL() string {
    return ts.server.URL
}

// Close shuts down the mock server
func (ts *TestServer) Close() {
    ts.server.Close()
}
```

### Bot Management

```go
// AddBot adds a bot to the mock server's state
func (ts *TestServer) AddBot(bot Bot) {
    // Bot is a simplified structure for test configuration
}

type Bot struct {
    ID        int
    Name      string
    Enabled   bool
    AccountID int
    // Optional: additional fields as needed for testing
}

// GetBot retrieves a bot by ID
func (ts *TestServer) GetBot(botID int) (Bot, bool) {
    // Returns bot and existence flag
}

// UpdateBot updates a bot's state
func (ts *TestServer) UpdateBot(botID int, updates BotUpdate) error {
    // Allows changing enabled state, etc.
}

// RemoveBot removes a bot from the mock
func (ts *TestServer) RemoveBot(botID int) {
    // Removes bot and all its deals
}
```

### Deal Management

```go
// AddDeal adds a deal to the mock server's state
func (ts *TestServer) AddDeal(botID int, deal Deal) error {
    // Returns error if bot doesn't exist
}

type Deal struct {
    ID        int
    BotID     int
    Pair      string
    Status    string
    ToCurrency string  // e.g., "BTC"
    FromCurrency string // e.g., "USDT"
    CreatedAt string   // ISO 8601
    UpdatedAt string   // ISO 8601
    Events    []BotEvent
}

type BotEvent struct {
    CreatedAt     string  // ISO 8601
    Action        string  // "place", "cancel", "cancelled", "modify"
    Coin          string  // "BTC", "ETH", etc.
    Type          string  // "buy", "sell"
    Status        string  // "active", "filled", "cancelled"
    Price         string  // Decimal string
    Size          string  // Decimal string
    OrderType     string  // "base", "safety", "take_profit"
    OrderSize     int
    OrderPosition int
    IsMarket      bool
}

// GetDeal retrieves a deal by ID
func (ts *TestServer) GetDeal(dealID int) (Deal, bool) {
    // Returns deal and existence flag
}

// UpdateDeal updates a deal's state
func (ts *TestServer) UpdateDeal(dealID int, updates DealUpdate) error {
    // Allows changing status, adding events, etc.
}

// AddBotEvent adds a new bot event to an existing deal
func (ts *TestServer) AddBotEvent(dealID int, event BotEvent) error {
    // Appends event to deal's event list
}

// RemoveDeal removes a deal from the mock
func (ts *TestServer) RemoveDeal(dealID int) {
    // Removes deal from state
}
```

### Error Simulation

```go
// SetRateLimitError configures the mock to return rate limit errors
func (ts *TestServer) SetRateLimitError(enabled bool, retryAfter int) {
    // When enabled, return 429 responses
    // retryAfter specifies seconds in Retry-After header
}

// SetBotError configures errors for specific bot operations
func (ts *TestServer) SetBotError(botID int, err error) {
    // Returns 404 or 500 for operations on this bot
}

// SetDealError configures errors for specific deal operations
func (ts *TestServer) SetDealError(dealID int, err error) {
    // Returns 404 or 500 for operations on this deal
}

// ClearErrors removes all configured errors
func (ts *TestServer) ClearErrors() {
    // Reset to normal operation
}
```

### State Inspection

```go
// GetAllBots returns all bots in the mock
func (ts *TestServer) GetAllBots() []Bot {
    // Returns all configured bots
}

// GetAllDeals returns all deals in the mock
func (ts *TestServer) GetAllDeals() []Deal {
    // Returns all configured deals
}

// GetBotDeals returns deals for a specific bot
func (ts *TestServer) GetBotDeals(botID int) []Deal {
    // Returns deals filtered by bot ID
}

// Reset clears all state
func (ts *TestServer) Reset() {
    // Removes all bots, deals, and errors
}
```

## VCR Integration (Future Enhancement)

The SDK uses `go-vcr` to record real API interactions. The mock server could support loading these recordings:

```go
// LoadFromVCRCassette loads bots and deals from a VCR cassette file
func (ts *TestServer) LoadFromVCRCassette(path string) error {
    // Parses VCR YAML file
    // Extracts API responses
    // Populates mock state with bots and deals
}
```

This would allow:
1. Recording real 3Commas API interactions during development
2. Using those recordings as test fixtures
3. Ensuring mock responses match real API behavior

**VCR Cassette Format Example:**
```yaml
---
version: 2
interactions:
- request:
    method: GET
    url: https://api.3commas.io/public/api/ver1/bots?scope=enabled
  response:
    status: 200
    headers:
      Content-Type: application/json
    body: '[{"id":1,"name":"Test Bot","is_enabled":true}]'
- request:
    method: GET
    url: https://api.3commas.io/public/api/ver1/deals?bot_id=1
  response:
    status: 200
    body: '[{"id":101,"bot_id":1,"status":"active"}]'
```

## Usage Example

```go
package main_test

import (
    "context"
    "testing"

    threecommasmock "github.com/recomma/3commas-mock/server"
    tc "github.com/recomma/3commas-sdk-go/threecommas"
    "github.com/stretchr/testify/require"
)

func TestE2E_DealProcessing(t *testing.T) {
    t.Parallel()

    // Create mock server
    mockServer := threecommasmock.NewTestServer(t)
    defer mockServer.Close()

    // Configure mock state
    mockServer.AddBot(threecommasmock.Bot{
        ID:      1,
        Name:    "Test Bot",
        Enabled: true,
    })

    mockServer.AddDeal(1, threecommasmock.Deal{
        ID:           101,
        BotID:        1,
        Pair:         "USDT_BTC",
        Status:       "active",
        ToCurrency:   "BTC",
        FromCurrency: "USDT",
        CreatedAt:    "2024-01-15T10:30:00.000Z",
        UpdatedAt:    "2024-01-15T10:30:00.000Z",
        Events: []threecommasmock.BotEvent{
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

    // Create SDK client pointing to mock
    client, err := tc.New3CommasClient(
        tc.WithBaseURL(mockServer.URL()),
        tc.WithAPIKey("test-key"),
        tc.WithPrivatePEM([]byte("test-secret")),
    )
    require.NoError(t, err)

    // Test ListBots
    ctx := context.Background()
    bots, err := client.ListBots(ctx, tc.WithScopeForListBots(tc.Enabled))
    require.NoError(t, err)
    require.Len(t, bots, 1)
    require.Equal(t, 1, bots[0].Id)

    // Test GetListOfDeals
    deals, err := client.GetListOfDeals(ctx, tc.WithBotIdForListDeals(1))
    require.NoError(t, err)
    require.Len(t, deals, 1)
    require.Equal(t, 101, deals[0].Id)

    // Test GetDealForID
    deal, err := client.GetDealForID(ctx, tc.DealPathId(101))
    require.NoError(t, err)
    require.Equal(t, 101, deal.Id)
    require.Len(t, deal.Events(), 1)

    // Simulate event addition (e.g., safety order triggered)
    mockServer.AddBotEvent(101, threecommasmock.BotEvent{
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

    // Fetch again, should see new event
    deal, err = client.GetDealForID(ctx, tc.DealPathId(101))
    require.NoError(t, err)
    require.Len(t, deal.Events(), 2)
}

func TestErrorHandling(t *testing.T) {
    t.Parallel()

    mockServer := threecommasmock.NewTestServer(t)
    defer mockServer.Close()

    // Enable rate limiting
    mockServer.SetRateLimitError(true, 60)

    client, _ := tc.New3CommasClient(
        tc.WithBaseURL(mockServer.URL()),
        tc.WithAPIKey("test-key"),
    )

    ctx := context.Background()
    _, err := client.ListBots(ctx)
    require.Error(t, err)

    // Check that it's a rate limit error
    var apiErr *tc.APIError
    require.ErrorAs(t, err, &apiErr)
    require.Equal(t, 429, apiErr.StatusCode)
}
```

## Implementation Guidelines

### Architecture

1. **State Management**
   - Use in-memory maps for bots and deals
   - Thread-safe access with mutexes
   - Fast lookups by ID

2. **HTTP Handlers**
   - Implement standard library `http.Handler` interface
   - Use `httptest.Server` for test server lifecycle
   - Return proper HTTP status codes and headers

3. **Response Marshaling**
   - Use JSON encoding/decoding
   - Match 3Commas API response format exactly
   - Support optional fields

4. **Error Handling**
   - Configurable error injection
   - Realistic error responses (matching 3Commas API errors)
   - Rate limit simulation with proper headers

### Testing the Mock Itself

The mock server package should include its own tests:
- Test that handlers return correct responses
- Test state management (add/update/remove)
- Test error simulation
- Test concurrent access safety

### Documentation

Include comprehensive documentation:
- Package-level GoDoc
- Examples in GoDoc
- README with quick start guide
- Reference to this specification

### Compatibility

- Go version: Same as Recomma project (Go 1.25.0)
- Dependencies: Minimal (standard library + testing)
- Similar to hyperliquid-mock patterns

## Future Enhancements

### Phase 2 Features (Optional)

1. **Webhook Simulation**
   - Simulate 3Commas webhooks for deal updates
   - Push-based event notifications

2. **Advanced Filtering**
   - Support more query parameters
   - Pagination simulation

3. **Performance Testing**
   - Configurable response delays
   - Load testing support

4. **Recording Mode**
   - Act as proxy to real 3Commas API
   - Record interactions for later playback
   - Generate VCR cassettes automatically

## References

- 3Commas API Documentation: https://github.com/3commas-io/3commas-official-api-docs
- Recomma 3Commas SDK: https://github.com/recomma/3commas-sdk-go
- Hyperliquid Mock (pattern reference): https://github.com/recomma/hyperliquid-mock
- go-vcr: https://github.com/dnaeon/go-vcr

## Questions / Clarifications

When implementing the mock server, clarify:

1. **VCR Integration Priority**: Should VCR cassette loading be part of the initial implementation or Phase 2?
2. **Authentication**: Should the mock verify API keys/signatures, or just accept any credentials in test mode?
3. **Webhook Support**: Are webhooks needed for E2E tests, or is polling sufficient?
4. **Deal Status Transitions**: Should the mock support automatic state transitions (e.g., deal → completed), or is manual control sufficient?

## Acceptance Criteria

The mock server implementation is complete when:

- ✅ All three API endpoints are implemented
- ✅ Test server can be created and used in tests
- ✅ Bots and deals can be added/updated/removed
- ✅ Error simulation is supported
- ✅ State inspection methods are available
- ✅ Package includes comprehensive tests
- ✅ Package documentation is complete
- ✅ Integration tests pass using the mock
- ✅ Usage examples are provided

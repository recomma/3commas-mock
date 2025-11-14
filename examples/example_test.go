package examples

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/recomma/3commas-mock/server"
)

// TestExampleBasicUsage demonstrates basic usage of the mock server
func TestExampleBasicUsage(t *testing.T) {
	// Create mock server
	mockServer := server.NewTestServer(t)
	defer mockServer.Close()

	// Configure mock state
	mockServer.AddBot(server.Bot{
		ID:        1,
		Name:      "Test Bot",
		Enabled:   true,
		AccountID: 12345,
	})

	mockServer.AddDeal(1, server.Deal{
		ID:           101,
		BotID:        1,
		Pair:         "USDT_BTC",
		Status:       "active",
		ToCurrency:   "BTC",
		FromCurrency: "USDT",
		CreatedAt:    "2024-01-15T10:30:00.000Z",
		UpdatedAt:    "2024-01-15T10:30:00.000Z",
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

	// Test ListBots
	resp, err := http.Get(mockServer.URL() + "/ver1/bots?scope=enabled")
	if err != nil {
		t.Fatalf("failed to get bots: %v", err)
	}
	defer resp.Body.Close()

	var bots []server.Bot
	if err := json.NewDecoder(resp.Body).Decode(&bots); err != nil {
		t.Fatalf("failed to decode bots: %v", err)
	}

	if len(bots) != 1 {
		t.Fatalf("expected 1 bot, got %d", len(bots))
	}
	if bots[0].ID != 1 {
		t.Fatalf("expected bot ID 1, got %d", bots[0].ID)
	}

	// Test GetDeal
	resp, err = http.Get(mockServer.URL() + "/ver1/deals/101/show")
	if err != nil {
		t.Fatalf("failed to get deal: %v", err)
	}
	defer resp.Body.Close()

	var deal server.Deal
	if err := json.NewDecoder(resp.Body).Decode(&deal); err != nil {
		t.Fatalf("failed to decode deal: %v", err)
	}

	if deal.ID != 101 {
		t.Fatalf("expected deal ID 101, got %d", deal.ID)
	}
	if len(deal.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(deal.Events))
	}
}

// TestExampleEventAddition demonstrates adding events to deals
func TestExampleEventAddition(t *testing.T) {
	mockServer := server.NewTestServer(t)
	defer mockServer.Close()

	mockServer.AddBot(server.Bot{ID: 1, Name: "Test Bot", Enabled: true})
	mockServer.AddDeal(1, server.Deal{
		ID:     101,
		BotID:  1,
		Status: "active",
		Events: []server.BotEvent{},
	})

	// Simulate safety order being placed
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

	// Fetch deal and verify event was added
	deal, ok := mockServer.GetDealByID(101)
	if !ok {
		t.Fatal("deal not found")
	}
	if len(deal.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(deal.Events))
	}
}

// TestExampleErrorSimulation demonstrates error simulation
func TestExampleErrorSimulation(t *testing.T) {
	mockServer := server.NewTestServer(t)
	defer mockServer.Close()

	// Enable rate limiting
	mockServer.SetRateLimitError(true, 60)

	resp, err := http.Get(mockServer.URL() + "/ver1/bots")
	if err != nil {
		t.Fatalf("failed to get bots: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("expected status 429, got %d", resp.StatusCode)
	}

	retryAfter := resp.Header.Get("Retry-After")
	if retryAfter != "60" {
		t.Fatalf("expected Retry-After: 60, got %s", retryAfter)
	}

	// Clear errors
	mockServer.ClearErrors()

	resp, err = http.Get(mockServer.URL() + "/ver1/bots")
	if err != nil {
		t.Fatalf("failed to get bots: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
}

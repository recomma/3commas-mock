package server

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestNewTestServer(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Close()

	if ts.URL() == "" {
		t.Fatal("expected non-empty URL")
	}
}

func TestListBots_Empty(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Close()

	resp, err := http.Get(ts.URL() + "/ver1/bots")
	if err != nil {
		t.Fatalf("failed to GET /ver1/bots: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var bots []Bot
	if err := json.NewDecoder(resp.Body).Decode(&bots); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(bots) != 0 {
		t.Fatalf("expected 0 bots, got %d", len(bots))
	}
}

func TestListBots_WithBots(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Close()

	// Add some bots
	ts.AddBot(Bot{
		ID:        1,
		Name:      "Test Bot 1",
		Enabled:   true,
		AccountID: 123,
	})
	ts.AddBot(Bot{
		ID:        2,
		Name:      "Test Bot 2",
		Enabled:   false,
		AccountID: 123,
	})

	resp, err := http.Get(ts.URL() + "/ver1/bots")
	if err != nil {
		t.Fatalf("failed to GET /ver1/bots: %v", err)
	}
	defer resp.Body.Close()

	var bots []Bot
	if err := json.NewDecoder(resp.Body).Decode(&bots); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(bots) != 2 {
		t.Fatalf("expected 2 bots, got %d", len(bots))
	}
}

func TestListBots_FilterByScope(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Close()

	ts.AddBot(Bot{ID: 1, Name: "Enabled Bot", Enabled: true, AccountID: 123})
	ts.AddBot(Bot{ID: 2, Name: "Disabled Bot", Enabled: false, AccountID: 123})

	// Filter for enabled bots
	resp, err := http.Get(ts.URL() + "/ver1/bots?scope=enabled")
	if err != nil {
		t.Fatalf("failed to GET /ver1/bots?scope=enabled: %v", err)
	}
	defer resp.Body.Close()

	var bots []Bot
	if err := json.NewDecoder(resp.Body).Decode(&bots); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(bots) != 1 {
		t.Fatalf("expected 1 enabled bot, got %d", len(bots))
	}
	if !bots[0].Enabled {
		t.Fatal("expected bot to be enabled")
	}
}

func TestGetDeal_NotFound(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Close()

	resp, err := http.Get(ts.URL() + "/ver1/deals/999/show")
	if err != nil {
		t.Fatalf("failed to GET /ver1/deals/999/show: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestGetDeal_Success(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Close()

	// Add bot and deal
	ts.AddBot(Bot{ID: 1, Name: "Test Bot", Enabled: true, AccountID: 123})
	err := ts.AddDeal(1, Deal{
		ID:           101,
		BotID:        1,
		Pair:         "USDT_BTC",
		Status:       "active",
		ToCurrency:   "BTC",
		FromCurrency: "USDT",
		CreatedAt:    "2024-01-15T10:30:00.000Z",
		UpdatedAt:    "2024-01-15T10:30:00.000Z",
		Events: []BotEvent{
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
	if err != nil {
		t.Fatalf("failed to add deal: %v", err)
	}

	resp, err := http.Get(ts.URL() + "/ver1/deals/101/show")
	if err != nil {
		t.Fatalf("failed to GET /ver1/deals/101/show: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var deal Deal
	if err := json.NewDecoder(resp.Body).Decode(&deal); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if deal.ID != 101 {
		t.Fatalf("expected deal ID 101, got %d", deal.ID)
	}
	if len(deal.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(deal.Events))
	}
	if deal.Events[0].Action != "place" {
		t.Fatalf("expected action 'place', got '%s'", deal.Events[0].Action)
	}
}

func TestAddBotEvent(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Close()

	ts.AddBot(Bot{ID: 1, Name: "Test Bot", Enabled: true, AccountID: 123})
	ts.AddDeal(1, Deal{
		ID:        101,
		BotID:     1,
		Pair:      "USDT_BTC",
		Status:    "active",
		CreatedAt: "2024-01-15T10:30:00.000Z",
		UpdatedAt: "2024-01-15T10:30:00.000Z",
		Events:    []BotEvent{},
	})

	// Add an event
	err := ts.AddBotEvent(101, BotEvent{
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
	if err != nil {
		t.Fatalf("failed to add bot event: %v", err)
	}

	// Verify event was added
	deal, ok := ts.GetDealByID(101)
	if !ok {
		t.Fatal("deal not found")
	}
	if len(deal.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(deal.Events))
	}
}

func TestRateLimitError(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Close()

	ts.SetRateLimitError(true, 60)

	resp, err := http.Get(ts.URL() + "/ver1/bots")
	if err != nil {
		t.Fatalf("failed to GET /ver1/bots: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("expected status 429, got %d", resp.StatusCode)
	}

	retryAfter := resp.Header.Get("Retry-After")
	if retryAfter != "60" {
		t.Fatalf("expected Retry-After: 60, got %s", retryAfter)
	}
}

func TestListDeals_FilterByBot(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Close()

	ts.AddBot(Bot{ID: 1, Name: "Bot 1", Enabled: true, AccountID: 123})
	ts.AddBot(Bot{ID: 2, Name: "Bot 2", Enabled: true, AccountID: 123})

	ts.AddDeal(1, Deal{ID: 101, BotID: 1, Status: "active"})
	ts.AddDeal(1, Deal{ID: 102, BotID: 1, Status: "active"})
	ts.AddDeal(2, Deal{ID: 103, BotID: 2, Status: "active"})

	resp, err := http.Get(ts.URL() + "/ver1/deals?bot_id=1")
	if err != nil {
		t.Fatalf("failed to GET /ver1/deals?bot_id=1: %v", err)
	}
	defer resp.Body.Close()

	var deals []Deal
	if err := json.NewDecoder(resp.Body).Decode(&deals); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(deals) != 2 {
		t.Fatalf("expected 2 deals for bot 1, got %d", len(deals))
	}
}

func TestReset(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Close()

	ts.AddBot(Bot{ID: 1, Name: "Test Bot", Enabled: true, AccountID: 123})
	ts.AddDeal(1, Deal{ID: 101, BotID: 1, Status: "active"})
	ts.SetRateLimitError(true, 60)

	ts.Reset()

	// Verify state was cleared
	bots := ts.GetAllBots()
	if len(bots) != 0 {
		t.Fatalf("expected 0 bots after reset, got %d", len(bots))
	}

	deals := ts.GetAllDeals()
	if len(deals) != 0 {
		t.Fatalf("expected 0 deals after reset, got %d", len(deals))
	}

	// Verify rate limit was cleared
	resp, err := http.Get(ts.URL() + "/ver1/bots")
	if err != nil {
		t.Fatalf("failed to GET /ver1/bots: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200 after reset, got %d", resp.StatusCode)
	}
}

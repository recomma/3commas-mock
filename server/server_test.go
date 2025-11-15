package server

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/recomma/3commas-mock/tcmock"
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

	var bots []tcmock.Bot
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
	ts.AddBot(NewBot(1, "Test Bot 1", 123, true))
	ts.AddBot(NewBot(2, "Test Bot 2", 123, false))

	resp, err := http.Get(ts.URL() + "/ver1/bots")
	if err != nil {
		t.Fatalf("failed to GET /ver1/bots: %v", err)
	}
	defer resp.Body.Close()

	var bots []tcmock.Bot
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

	ts.AddBot(NewBot(1, "Enabled Bot", 123, true))
	ts.AddBot(NewBot(2, "Disabled Bot", 123, false))

	// Filter for enabled bots
	resp, err := http.Get(ts.URL() + "/ver1/bots?scope=enabled")
	if err != nil {
		t.Fatalf("failed to GET /ver1/bots?scope=enabled: %v", err)
	}
	defer resp.Body.Close()

	var bots []tcmock.Bot
	if err := json.NewDecoder(resp.Body).Decode(&bots); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(bots) != 1 {
		t.Fatalf("expected 1 enabled bot, got %d", len(bots))
	}
	if !bots[0].IsEnabled {
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
	ts.AddBot(NewBot(1, "Test Bot", 123, true))

	deal := NewDeal(101, 1, "USDT_BTC", "active")
	AddBotEvent(&deal, "Placing base order. Price: 50000.0 USDT Size: 0.0002 BTC")

	err := ts.AddDeal(deal)
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

	var dealResp tcmock.Deal
	if err := json.NewDecoder(resp.Body).Decode(&dealResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if dealResp.Id != 101 {
		t.Fatalf("expected deal ID 101, got %d", dealResp.Id)
	}
	if len(dealResp.BotEvents) != 1 {
		t.Fatalf("expected 1 event, got %d", len(dealResp.BotEvents))
	}
	if dealResp.BotEvents[0].Message == nil {
		t.Fatal("expected bot event to have a message")
	}
	if *dealResp.BotEvents[0].Message != "Placing base order. Price: 50000.0 USDT Size: 0.0002 BTC" {
		t.Fatalf("expected specific message, got '%s'", *dealResp.BotEvents[0].Message)
	}
}

func TestAddBotEvent(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Close()

	ts.AddBot(NewBot(1, "Test Bot", 123, true))
	deal := NewDeal(101, 1, "USDT_BTC", "active")
	ts.AddDeal(deal)

	// Add an event
	err := ts.AddBotEventToDeal(101, "Placing safety order. Price: 48750.0 USDT Size: 0.0004 BTC")
	if err != nil {
		t.Fatalf("failed to add bot event: %v", err)
	}

	// Verify event was added
	dealResp, ok := ts.GetDealByID(101)
	if !ok {
		t.Fatal("deal not found")
	}
	if len(dealResp.BotEvents) != 1 {
		t.Fatalf("expected 1 event, got %d", len(dealResp.BotEvents))
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

	ts.AddBot(NewBot(1, "Bot 1", 123, true))
	ts.AddBot(NewBot(2, "Bot 2", 123, true))

	ts.AddDeal(NewDeal(101, 1, "USDT_BTC", "active"))
	ts.AddDeal(NewDeal(102, 1, "USDT_ETH", "active"))
	ts.AddDeal(NewDeal(103, 2, "USDT_BTC", "active"))

	resp, err := http.Get(ts.URL() + "/ver1/deals?bot_id=1")
	if err != nil {
		t.Fatalf("failed to GET /ver1/deals?bot_id=1: %v", err)
	}
	defer resp.Body.Close()

	var deals []tcmock.Deal
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

	ts.AddBot(NewBot(1, "Test Bot", 123, true))
	ts.AddDeal(NewDeal(101, 1, "USDT_BTC", "active"))
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

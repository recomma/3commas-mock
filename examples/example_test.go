package examples

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/recomma/3commas-mock/server"
	"github.com/recomma/3commas-mock/tcmock"
)

// TestExampleBasicUsage demonstrates basic usage of the mock server
func TestExampleBasicUsage(t *testing.T) {
	// Create mock 3comma server
	mockServer := server.NewTestServer(t)
	defer mockServer.Close()

	// Configure mock state
	mockServer.AddBot(server.NewBot(1, "Test Bot", 12345, true))

	deal := server.NewDeal(101, 1, "USDT_BTC", "active")
	server.AddBotEvent(&deal, "Placing base order. Price: 50000.0 USDT Size: 0.0002 BTC")
	mockServer.AddDeal(deal)

	// Test ListBots
	resp, err := http.Get(mockServer.URL() + "/ver1/bots?scope=enabled")
	if err != nil {
		t.Fatalf("failed to get bots: %v", err)
	}
	defer resp.Body.Close()

	var bots []tcmock.Bot
	if err := json.NewDecoder(resp.Body).Decode(&bots); err != nil {
		t.Fatalf("failed to decode bots: %v", err)
	}

	if len(bots) != 1 {
		t.Fatalf("expected 1 bot, got %d", len(bots))
	}
	if bots[0].Id != 1 {
		t.Fatalf("expected bot ID 1, got %d", bots[0].Id)
	}

	// Test GetDeal
	resp, err = http.Get(mockServer.URL() + "/ver1/deals/101/show")
	if err != nil {
		t.Fatalf("failed to get deal: %v", err)
	}
	defer resp.Body.Close()

	var dealResp tcmock.Deal
	if err := json.NewDecoder(resp.Body).Decode(&dealResp); err != nil {
		t.Fatalf("failed to decode deal: %v", err)
	}

	if dealResp.Id != 101 {
		t.Fatalf("expected deal ID 101, got %d", dealResp.Id)
	}
	if len(dealResp.BotEvents) != 1 {
		t.Fatalf("expected 1 event, got %d", len(dealResp.BotEvents))
	}
}

// TestExampleEventAddition demonstrates adding events to deals
func TestExampleEventAddition(t *testing.T) {
	mockServer := server.NewTestServer(t)
	defer mockServer.Close()

	mockServer.AddBot(server.NewBot(1, "Test Bot", 12345, true))
	deal := server.NewDeal(101, 1, "USDT_BTC", "active")
	mockServer.AddDeal(deal)

	// Simulate safety order being placed
	mockServer.AddBotEventToDeal(101, "Placing safety order. Price: 48750.0 USDT Size: 0.0004 BTC")

	// Fetch deal and verify event was added
	dealResp, ok := mockServer.GetDealByID(101)
	if !ok {
		t.Fatal("deal not found")
	}
	if len(dealResp.BotEvents) != 1 {
		t.Fatalf("expected 1 event, got %d", len(dealResp.BotEvents))
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

// TestExampleVCRLoading demonstrates loading real API responses from VCR cassettes
func TestExampleVCRLoading(t *testing.T) {
	mockServer := server.NewTestServer(t)
	defer mockServer.Close()

	// Load a VCR cassette with a real recorded deal from 3commas API
	// Note: go-vcr appends .yaml automatically, so omit the extension
	err := mockServer.LoadVCRCassette("../testdata/fixtures/deal_2376446537")
	if err != nil {
		t.Fatalf("failed to load VCR cassette: %v", err)
	}

	// Now the deal is available via the API with ALL real data
	resp, err := http.Get(mockServer.URL() + "/ver1/deals/2376446537/show")
	if err != nil {
		t.Fatalf("failed to get deal: %v", err)
	}
	defer resp.Body.Close()

	var deal tcmock.Deal
	if err := json.NewDecoder(resp.Body).Decode(&deal); err != nil {
		t.Fatalf("failed to decode deal: %v", err)
	}

	// Verify we got real data from the VCR recording
	if deal.Pair != "USDT_DOGE" {
		t.Fatalf("expected pair USDT_DOGE, got %s", deal.Pair)
	}

	// Most importantly: bot_events are preserved from the real API!
	if len(deal.BotEvents) != 3 {
		t.Fatalf("expected 3 bot events from VCR, got %d", len(deal.BotEvents))
	}

	// The messages are real ones from 3commas
	if deal.BotEvents[0].Message == nil ||
		*deal.BotEvents[0].Message != "Placing averaging order (9 out of 9). Price: market Size: 25.0008 USDT (110.0 DOGE)" {
		t.Fatal("bot_events not preserved correctly from VCR")
	}

	// Bot was auto-created from the deal data
	bot, ok := mockServer.GetBot(16511317)
	if !ok {
		t.Fatal("bot should have been auto-created from deal")
	}
	if bot.AccountName != "Demo Account 2080398" {
		t.Fatalf("expected account name from VCR, got %s", bot.AccountName)
	}
}

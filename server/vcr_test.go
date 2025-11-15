package server

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/recomma/3commas-mock/tcmock"
)

func TestLoadVCRCassette_SingleDeal(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Close()

	// Load VCR cassette with a real recorded deal
	// Note: go-vcr appends .yaml automatically, so omit extension
	err := ts.LoadVCRCassette("../testdata/fixtures/deal_2376446537")
	if err != nil {
		t.Fatalf("failed to load VCR cassette: %v", err)
	}

	// Verify the deal was loaded
	deal, ok := ts.GetDealByID(2376446537)
	if !ok {
		t.Fatal("deal 2376446537 not found after loading VCR cassette")
	}

	// Verify deal fields from real API
	if deal.Pair != "USDT_DOGE" {
		t.Errorf("expected pair USDT_DOGE, got %s", deal.Pair)
	}
	if deal.Status != "bought" {
		t.Errorf("expected status bought, got %s", deal.Status)
	}
	if deal.BotId != 16511317 {
		t.Errorf("expected bot_id 16511317, got %d", deal.BotId)
	}

	// MOST IMPORTANT: Verify bot_events were preserved from real API
	if len(deal.BotEvents) != 3 {
		t.Fatalf("expected 3 bot events, got %d", len(deal.BotEvents))
	}

	// Check first event
	if deal.BotEvents[0].Message == nil {
		t.Fatal("expected first event to have a message")
	}
	expectedMsg := "Placing averaging order (9 out of 9). Price: market Size: 25.0008 USDT (110.0 DOGE)"
	if *deal.BotEvents[0].Message != expectedMsg {
		t.Errorf("expected first event message: %s, got: %s", expectedMsg, *deal.BotEvents[0].Message)
	}

	// Verify bot was auto-created
	bot, ok := ts.GetBot(16511317)
	if !ok {
		t.Fatal("bot 16511317 should have been auto-created")
	}
	if bot.AccountName != "Demo Account 2080398" {
		t.Errorf("expected account name from deal, got %s", bot.AccountName)
	}

	// Test that the deal is accessible via HTTP API
	resp, err := http.Get(ts.URL() + "/ver1/deals/2376446537/show")
	if err != nil {
		t.Fatalf("failed to GET deal: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var dealResp tcmock.Deal
	if err := json.NewDecoder(resp.Body).Decode(&dealResp); err != nil {
		t.Fatalf("failed to decode deal: %v", err)
	}

	// Verify bot_events are returned via API
	if len(dealResp.BotEvents) != 3 {
		t.Fatalf("expected 3 bot events via API, got %d", len(dealResp.BotEvents))
	}
}

func TestLoadVCRCassette_DuplicateError(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Close()

	// Load cassette first time - should succeed
	err := ts.LoadVCRCassette("../testdata/fixtures/deal_2376446537")
	if err != nil {
		t.Fatalf("first load failed: %v", err)
	}

	// Load same cassette again - should fail with duplicate error
	err = ts.LoadVCRCassette("../testdata/fixtures/deal_2376446537")
	if err == nil {
		t.Fatal("expected error for duplicate deal ID, got nil")
	}

	if err.Error() != "failed to process interaction 0 from ../testdata/fixtures/deal_2376446537: duplicate deal ID 2376446537 found in VCR cassette" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestLoadVCRCassettes_Multiple(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Close()

	// Load multiple cassettes at once
	err := ts.LoadVCRCassettes(
		"../testdata/fixtures/deal_2376446537",
		// Add more cassette paths here as you create them
	)
	if err != nil {
		t.Fatalf("failed to load cassettes: %v", err)
	}

	// Verify deal from first cassette
	deal, ok := ts.GetDealByID(2376446537)
	if !ok {
		t.Fatal("deal from first cassette not found")
	}
	if len(deal.BotEvents) != 3 {
		t.Fatalf("expected 3 bot events, got %d", len(deal.BotEvents))
	}
}

func TestExtractDealID(t *testing.T) {
	tests := []struct {
		url      string
		expected int
		wantErr  bool
	}{
		{
			url:      "https://api.3commas.io/public/api/ver1/deals/2376446537/show",
			expected: 2376446537,
			wantErr:  false,
		},
		{
			url:      "/ver1/deals/123/show",
			expected: 123,
			wantErr:  false,
		},
		{
			url:      "/ver1/deals",
			expected: 0,
			wantErr:  true,
		},
		{
			url:      "/ver1/bots",
			expected: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			id, err := ExtractDealID(tt.url)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error for URL %s, got nil", tt.url)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for URL %s: %v", tt.url, err)
				}
				if id != tt.expected {
					t.Errorf("expected ID %d, got %d", tt.expected, id)
				}
			}
		})
	}
}

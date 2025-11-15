package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/recomma/3commas-mock/tcmock"
	"gopkg.in/dnaeon/go-vcr.v3/cassette"
)

// URL pattern matchers
var (
	dealShowPattern  = regexp.MustCompile(`/ver1/deals/(\d+)/show$`)
	dealsListPattern = regexp.MustCompile(`/ver1/deals(\?.*)?$`)
	botsListPattern  = regexp.MustCompile(`/ver1/bots(\?.*)?$`)
)

// LoadVCRCassette loads a VCR cassette and populates mock server state
// - Bots and Deals are loaded with ALL their data from real API responses
// - bot_events are PRESERVED exactly as recorded (this is the valuable part!)
// - Duplicate IDs will return an error
// - Non-2xx responses are skipped
func (ts *TestServer) LoadVCRCassette(cassettePath string) error {
	// Load cassette from file
	c, err := cassette.Load(cassettePath)
	if err != nil {
		return fmt.Errorf("failed to load VCR cassette %s: %w", cassettePath, err)
	}

	// Process each interaction
	for i, interaction := range c.Interactions {
		// Skip non-successful responses
		if interaction.Response.Code < 200 || interaction.Response.Code >= 300 {
			continue
		}

		// Only process GET requests
		if interaction.Request.Method != http.MethodGet {
			continue
		}

		// Try to match against known patterns
		if err := ts.processInteraction(interaction); err != nil {
			return fmt.Errorf("failed to process interaction %d from %s: %w", i, cassettePath, err)
		}
	}

	return nil
}

// LoadVCRCassettes loads multiple VCR cassettes in sequence
func (ts *TestServer) LoadVCRCassettes(cassettePaths ...string) error {
	for _, path := range cassettePaths {
		if err := ts.LoadVCRCassette(path); err != nil {
			return err
		}
	}
	return nil
}

// processInteraction processes a single VCR interaction and adds entities to state
func (ts *TestServer) processInteraction(interaction *cassette.Interaction) error {
	url := interaction.Request.URL
	body := interaction.Response.Body

	// Match against deal show endpoint: /ver1/deals/{id}/show
	if matches := dealShowPattern.FindStringSubmatch(url); matches != nil {
		return ts.loadDealFromJSON(body)
	}

	// Match against deals list endpoint: /ver1/deals
	if dealsListPattern.MatchString(url) {
		return ts.loadDealsListFromJSON(body)
	}

	// Match against bots list endpoint: /ver1/bots
	if botsListPattern.MatchString(url) {
		return ts.loadBotsListFromJSON(body)
	}

	// Unknown endpoint - skip silently
	return nil
}

// loadDealFromJSON unmarshals a single deal and adds it to state
func (ts *TestServer) loadDealFromJSON(jsonBody string) error {
	var deal tcmock.Deal
	if err := json.Unmarshal([]byte(jsonBody), &deal); err != nil {
		return fmt.Errorf("failed to unmarshal deal: %w", err)
	}

	// Check for duplicate
	ts.mu.RLock()
	if _, exists := ts.deals[deal.Id]; exists {
		ts.mu.RUnlock()
		return fmt.Errorf("duplicate deal ID %d found in VCR cassette", deal.Id)
	}
	ts.mu.RUnlock()

	// Check if bot exists, create a minimal one if not
	ts.mu.RLock()
	botExists := ts.bots[deal.BotId] != nil
	ts.mu.RUnlock()

	if !botExists {
		// Create a minimal bot based on deal data
		bot := NewBot(deal.BotId, deal.BotName, deal.AccountId, true)
		bot.AccountName = deal.AccountName
		ts.AddBot(bot)
	}

	// Add deal with all its data (including bot_events preserved)
	ts.mu.Lock()
	ts.deals[deal.Id] = &deal
	ts.mu.Unlock()

	return nil
}

// loadDealsListFromJSON unmarshals a list of deals and adds them to state
func (ts *TestServer) loadDealsListFromJSON(jsonBody string) error {
	var deals []tcmock.Deal
	if err := json.Unmarshal([]byte(jsonBody), &deals); err != nil {
		return fmt.Errorf("failed to unmarshal deals list: %w", err)
	}

	for _, deal := range deals {
		// Check for duplicate
		ts.mu.RLock()
		if _, exists := ts.deals[deal.Id]; exists {
			ts.mu.RUnlock()
			return fmt.Errorf("duplicate deal ID %d found in VCR cassette", deal.Id)
		}
		ts.mu.RUnlock()

		// Check if bot exists, create a minimal one if not
		ts.mu.RLock()
		botExists := ts.bots[deal.BotId] != nil
		ts.mu.RUnlock()

		if !botExists {
			// Create a minimal bot based on deal data
			bot := NewBot(deal.BotId, deal.BotName, deal.AccountId, true)
			bot.AccountName = deal.AccountName
			ts.AddBot(bot)
		}

		// Add deal with all its data (including bot_events preserved)
		ts.mu.Lock()
		dealCopy := deal // Create a copy to get the correct pointer
		ts.deals[dealCopy.Id] = &dealCopy
		ts.mu.Unlock()
	}

	return nil
}

// loadBotsListFromJSON unmarshals a list of bots and adds them to state
func (ts *TestServer) loadBotsListFromJSON(jsonBody string) error {
	var bots []tcmock.Bot
	if err := json.Unmarshal([]byte(jsonBody), &bots); err != nil {
		return fmt.Errorf("failed to unmarshal bots list: %w", err)
	}

	for _, bot := range bots {
		// Check for duplicate
		ts.mu.RLock()
		if _, exists := ts.bots[bot.Id]; exists {
			ts.mu.RUnlock()
			return fmt.Errorf("duplicate bot ID %d found in VCR cassette", bot.Id)
		}
		ts.mu.RUnlock()

		// Add bot with all its data
		ts.mu.Lock()
		botCopy := bot // Create a copy to get the correct pointer
		ts.bots[botCopy.Id] = &botCopy
		ts.mu.Unlock()
	}

	return nil
}

// ExtractDealID is a helper to extract deal ID from a deal show URL
func ExtractDealID(url string) (int, error) {
	matches := dealShowPattern.FindStringSubmatch(url)
	if matches == nil {
		return 0, fmt.Errorf("URL does not match deal show pattern: %s", url)
	}
	return strconv.Atoi(matches[1])
}

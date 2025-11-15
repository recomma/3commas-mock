package server

import (
	"fmt"

	"github.com/recomma/3commas-mock/tcmock"
)

// Bot Management

// AddBot adds a bot to the mock server's state
func (ts *TestServer) AddBot(bot tcmock.Bot) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.bots[bot.Id] = &bot
}

// GetBot retrieves a bot by ID
func (ts *TestServer) GetBot(botID int) (tcmock.Bot, bool) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	bot, ok := ts.bots[botID]
	if !ok {
		return tcmock.Bot{}, false
	}
	return *bot, true
}

// UpdateBotEnabled updates a bot's enabled state
func (ts *TestServer) UpdateBotEnabled(botID int, enabled bool) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	bot, ok := ts.bots[botID]
	if !ok {
		return fmt.Errorf("bot %d not found", botID)
	}

	bot.IsEnabled = enabled

	return nil
}

// UpdateBotName updates a bot's name
func (ts *TestServer) UpdateBotName(botID int, name string) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	bot, ok := ts.bots[botID]
	if !ok {
		return fmt.Errorf("bot %d not found", botID)
	}

	bot.Name = &name

	return nil
}

// RemoveBot removes a bot from the mock
func (ts *TestServer) RemoveBot(botID int) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	delete(ts.bots, botID)

	// Remove all deals for this bot
	for dealID, deal := range ts.deals {
		if deal.BotId == botID {
			delete(ts.deals, dealID)
		}
	}
}

// GetAllBots returns all bots in the mock
func (ts *TestServer) GetAllBots() []tcmock.Bot {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	result := make([]tcmock.Bot, 0, len(ts.bots))
	for _, bot := range ts.bots {
		result = append(result, *bot)
	}
	return result
}

// Deal Management

// AddDeal adds a deal to the mock server's state
func (ts *TestServer) AddDeal(deal tcmock.Deal) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	// Check if bot exists
	if _, ok := ts.bots[deal.BotId]; !ok {
		return fmt.Errorf("bot %d not found", deal.BotId)
	}

	ts.deals[deal.Id] = &deal
	return nil
}

// GetDealByID retrieves a deal by ID (state management method)
func (ts *TestServer) GetDealByID(dealID int) (tcmock.Deal, bool) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	deal, ok := ts.deals[dealID]
	if !ok {
		return tcmock.Deal{}, false
	}
	return *deal, true
}

// UpdateDealStatus updates a deal's status
func (ts *TestServer) UpdateDealStatus(dealID int, status string) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	deal, ok := ts.deals[dealID]
	if !ok {
		return fmt.Errorf("deal %d not found", dealID)
	}

	deal.Status = tcmock.DealStatus(status)

	return nil
}

// AddBotEventToDeal adds a new bot event to an existing deal
// message: Human-readable event description
func (ts *TestServer) AddBotEventToDeal(dealID int, message string) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	deal, ok := ts.deals[dealID]
	if !ok {
		return fmt.Errorf("deal %d not found", dealID)
	}

	AddBotEvent(deal, message)
	return nil
}

// RemoveDeal removes a deal from the mock
func (ts *TestServer) RemoveDeal(dealID int) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	delete(ts.deals, dealID)
}

// GetAllDeals returns all deals in the mock
func (ts *TestServer) GetAllDeals() []tcmock.Deal {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	result := make([]tcmock.Deal, 0, len(ts.deals))
	for _, deal := range ts.deals {
		result = append(result, *deal)
	}
	return result
}

// GetBotDeals returns deals for a specific bot
func (ts *TestServer) GetBotDeals(botID int) []tcmock.Deal {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	result := make([]tcmock.Deal, 0)
	for _, deal := range ts.deals {
		if deal.BotId == botID {
			result = append(result, *deal)
		}
	}
	return result
}

// Error Simulation

// SetRateLimitError configures the mock to return rate limit errors
func (ts *TestServer) SetRateLimitError(enabled bool, retryAfter int) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.rateLimitEnabled = enabled
	ts.rateLimitRetry = retryAfter
}

// SetBotError configures errors for specific bot operations
func (ts *TestServer) SetBotError(botID int, err error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.botErrors[botID] = err
}

// SetDealError configures errors for specific deal operations
func (ts *TestServer) SetDealError(dealID int, err error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.dealErrors[dealID] = err
}

// ClearErrors removes all configured errors
func (ts *TestServer) ClearErrors() {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.botErrors = make(map[int]error)
	ts.dealErrors = make(map[int]error)
	ts.rateLimitEnabled = false
	ts.rateLimitRetry = 0
}

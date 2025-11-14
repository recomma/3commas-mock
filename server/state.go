package server

import "fmt"

// Bot Management

// AddBot adds a bot to the mock server's state
func (ts *TestServer) AddBot(bot Bot) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.bots[bot.ID] = &bot
}

// GetBot retrieves a bot by ID
func (ts *TestServer) GetBot(botID int) (Bot, bool) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	bot, ok := ts.bots[botID]
	if !ok {
		return Bot{}, false
	}
	return *bot, true
}

// UpdateBot updates a bot's state
func (ts *TestServer) UpdateBot(botID int, updates BotUpdate) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	bot, ok := ts.bots[botID]
	if !ok {
		return fmt.Errorf("bot %d not found", botID)
	}

	if updates.Enabled != nil {
		bot.Enabled = *updates.Enabled
	}
	if updates.Name != nil {
		bot.Name = *updates.Name
	}

	return nil
}

// RemoveBot removes a bot from the mock
func (ts *TestServer) RemoveBot(botID int) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	delete(ts.bots, botID)

	// Remove all deals for this bot
	for dealID, deal := range ts.deals {
		if deal.BotID == botID {
			delete(ts.deals, dealID)
		}
	}
}

// GetAllBots returns all bots in the mock
func (ts *TestServer) GetAllBots() []Bot {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	result := make([]Bot, 0, len(ts.bots))
	for _, bot := range ts.bots {
		result = append(result, *bot)
	}
	return result
}

// Deal Management

// AddDeal adds a deal to the mock server's state
func (ts *TestServer) AddDeal(botID int, deal Deal) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	// Check if bot exists
	if _, ok := ts.bots[botID]; !ok {
		return fmt.Errorf("bot %d not found", botID)
	}

	deal.BotID = botID
	ts.deals[deal.ID] = &deal
	return nil
}

// GetDealByID retrieves a deal by ID (state management method)
func (ts *TestServer) GetDealByID(dealID int) (Deal, bool) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	deal, ok := ts.deals[dealID]
	if !ok {
		return Deal{}, false
	}
	return *deal, true
}

// UpdateDeal updates a deal's state
func (ts *TestServer) UpdateDeal(dealID int, updates DealUpdate) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	deal, ok := ts.deals[dealID]
	if !ok {
		return fmt.Errorf("deal %d not found", dealID)
	}

	if updates.Status != nil {
		deal.Status = *updates.Status
	}

	return nil
}

// AddBotEvent adds a new bot event to an existing deal
func (ts *TestServer) AddBotEvent(dealID int, event BotEvent) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	deal, ok := ts.deals[dealID]
	if !ok {
		return fmt.Errorf("deal %d not found", dealID)
	}

	deal.Events = append(deal.Events, event)
	return nil
}

// RemoveDeal removes a deal from the mock
func (ts *TestServer) RemoveDeal(dealID int) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	delete(ts.deals, dealID)
}

// GetAllDeals returns all deals in the mock
func (ts *TestServer) GetAllDeals() []Deal {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	result := make([]Deal, 0, len(ts.deals))
	for _, deal := range ts.deals {
		result = append(result, *deal)
	}
	return result
}

// GetBotDeals returns deals for a specific bot
func (ts *TestServer) GetBotDeals(botID int) []Deal {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	result := make([]Deal, 0)
	for _, deal := range ts.deals {
		if deal.BotID == botID {
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

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/recomma/3commas-mock/tcmock"
)

// TestServer wraps the mock 3Commas server for testing
type TestServer struct {
	server *httptest.Server
	mu     sync.RWMutex

	// State
	bots  map[int]*tcmock.Bot
	deals map[int]*tcmock.Deal

	// Error simulation
	rateLimitEnabled bool
	rateLimitRetry   int
	botErrors        map[int]error
	dealErrors       map[int]error
}

// NewTestServer creates a new mock 3Commas server for testing
func NewTestServer(t *testing.T) *TestServer {
	ts := &TestServer{
		bots:       make(map[int]*tcmock.Bot),
		deals:      make(map[int]*tcmock.Deal),
		botErrors:  make(map[int]error),
		dealErrors: make(map[int]error),
	}

	// Create HTTP handler using the generated HandlerFromMux
	handler := tcmock.HandlerFromMux(ts, http.NewServeMux())
	ts.server = httptest.NewServer(handler)

	return ts
}

// URL returns the base URL of the mock server
func (ts *TestServer) URL() string {
	return ts.server.URL
}

// Close shuts down the mock server
func (ts *TestServer) Close() {
	ts.server.Close()
}

// Reset clears all state
func (ts *TestServer) Reset() {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.bots = make(map[int]*tcmock.Bot)
	ts.deals = make(map[int]*tcmock.Deal)
	ts.botErrors = make(map[int]error)
	ts.dealErrors = make(map[int]error)
	ts.rateLimitEnabled = false
}

// ListBots implements the ServerInterface method for GET /ver1/bots
func (ts *TestServer) ListBots(w http.ResponseWriter, r *http.Request, params tcmock.ListBotsParams) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	// Check for rate limit simulation
	if ts.rateLimitEnabled {
		if ts.rateLimitRetry > 0 {
			w.Header().Set("Retry-After", fmt.Sprintf("%d", ts.rateLimitRetry))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]string{
			"error":             "rate limit exceeded",
			"error_description": "You have exceeded the rate limit. Please try again later.",
		})
		return
	}

	// Filter bots based on scope parameter
	var result []tcmock.Bot
	for _, bot := range ts.bots {
		// Apply scope filter if provided
		if params.Scope != nil {
			if *params.Scope == tcmock.Enabled && !bot.IsEnabled {
				continue
			}
			if *params.Scope == tcmock.Disabled && bot.IsEnabled {
				continue
			}
		}
		result = append(result, *bot)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// ListDeals implements the ServerInterface method for GET /ver1/deals
func (ts *TestServer) ListDeals(w http.ResponseWriter, r *http.Request, params tcmock.ListDealsParams) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	// Filter deals based on parameters
	var result []tcmock.Deal
	for _, deal := range ts.deals {
		// Apply bot_id filter if provided
		if params.BotId != nil && deal.BotId != *params.BotId {
			continue
		}

		// Apply scope filter if provided
		if params.Scope != nil && tcmock.DealStatus(*params.Scope) != deal.Status {
			continue
		}

		result = append(result, *deal)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// GetDeal implements the ServerInterface method for GET /ver1/deals/{deal_id}/show
func (ts *TestServer) GetDeal(w http.ResponseWriter, r *http.Request, dealID tcmock.DealPathId) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	// Check for deal-specific error
	if err := ts.dealErrors[dealID]; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	deal, ok := ts.deals[dealID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "deal not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deal)
}

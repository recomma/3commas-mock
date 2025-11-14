package server

// BotEvent represents a bot event with full structured data for testing
// This overrides the simple generated type to provide rich event details
type BotEvent struct {
	CreatedAt     string `json:"created_at"`
	Action        string `json:"action"`         // "place", "cancel", "cancelled", "modify"
	Coin          string `json:"coin"`           // "BTC", "ETH", etc.
	Type          string `json:"type"`           // "buy", "sell"
	Status        string `json:"status"`         // "active", "filled", "cancelled"
	Price         string `json:"price"`          // Decimal string
	Size          string `json:"size"`           // Decimal string
	OrderType     string `json:"order_type"`     // "base", "safety", "take_profit"
	OrderSize     int    `json:"order_size"`     // Order size category
	OrderPosition int    `json:"order_position"` // Safety order position
	IsMarket      bool   `json:"is_market"`      // Whether order is market order
}

// Bot represents a simplified bot for testing
type Bot struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Enabled   bool   `json:"is_enabled"`
	AccountID int    `json:"account_id"`
	// Add other fields as needed
}

// Deal represents a simplified deal for testing
type Deal struct {
	ID           int        `json:"id"`
	BotID        int        `json:"bot_id"`
	Pair         string     `json:"pair"`
	Status       string     `json:"status"`
	ToCurrency   string     `json:"to_currency"`
	FromCurrency string     `json:"from_currency"`
	CreatedAt    string     `json:"created_at"`
	UpdatedAt    string     `json:"updated_at"`
	Events       []BotEvent `json:"bot_events"`
}

// BotUpdate represents fields that can be updated on a bot
type BotUpdate struct {
	Enabled *bool
	Name    *string
}

// DealUpdate represents fields that can be updated on a deal
type DealUpdate struct {
	Status *string
}

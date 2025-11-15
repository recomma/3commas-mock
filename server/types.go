package server

import (
	"time"

	"github.com/oapi-codegen/nullable"
	"github.com/recomma/3commas-mock/tcmock"
)

// Helper types and functions for creating test data with the full generated types

// NewBot creates a minimal Bot with required fields populated
// This is a helper to make it easier to create test bots
func NewBot(id int, name string, accountID int, enabled bool) tcmock.Bot {
	now := time.Now()
	namePtr := &name
	strategy := tcmock.BotStrategyLong
	return tcmock.Bot{
		Id:                          id,
		Name:                        namePtr,
		AccountId:                   accountID,
		AccountName:                 "Test Account",
		IsEnabled:                   enabled,
		CreatedAt:                   now,
		ActiveDealsCount:            0,
		FinishedDealsCount:          "0",
		ActiveDeals:                 []tcmock.Deal{},
		ActiveDealsBtcProfit:        "0",
		ActiveDealsUsdProfit:        "0",
		BtcFundsLockedInActiveDeals: "0",
		FundsLockedInActiveDeals:    "0",
		FinishedDealsProfitUsd:      "0",
		Deletable:                   true,
		Strategy:                    &strategy,
	}
}

// NewDeal creates a minimal Deal with required fields populated
// This is a helper to make it easier to create test deals
func NewDeal(id int, botID int, pair string, status string) tcmock.Deal {
	now := time.Now()
	// Extract currency from pair (assumes format like "USDT_BTC")
	toCurrency := "BTC"
	if len(pair) > 3 {
		toCurrency = pair[len(pair)-3:]
	}
	return tcmock.Deal{
		Id:                        id,
		BotId:                     botID,
		Pair:                      pair,
		Status:                    tcmock.DealStatus(status),
		CreatedAt:                 now,
		UpdatedAt:                 now,
		BotEvents:                 []struct {
			CreatedAt *time.Time `json:"created_at,omitempty"`
			Message   *string    `json:"message,omitempty"`
		}{},
		AccountId:                        1,
		AccountName:                      "Test Account",
		BotName:                          "Test Bot",
		FromCurrency:                     "USDT",
		ToCurrency:                       toCurrency,
		BaseOrderVolume:                  "10",
		BaseOrderVolumeType:              "quote_currency",
		BoughtAmount:                     "0",
		BoughtAveragePrice:               "0",
		BoughtVolume:                     "0",
		SoldAmount:                       "0",
		SoldAveragePrice:                 "0",
		SoldVolume:                       "0",
		ActualProfitPercentage:           "0",
		ActualProfit:                     nullable.NewNullableWithValue("0"),
		ActualUsdProfit:                  nullable.NewNullableWithValue("0"),
		FinalProfit:                      "0",
		FinalProfitPercentage:            "0",
		Cancellable:                      false,
		CompletedManualSafetyOrdersCount: 0,
		CompletedSafetyOrdersCount:       0,
		CurrentActiveSafetyOrders:        0,
		CurrentActiveSafetyOrdersCount:   0,
		CurrentPrice:                     "0",
		TakeProfitPrice:                  "0",
		MaxSafetyOrders:                  0,
		ActiveManualSafetyOrders:         0,
		TrailingEnabled:                  false,
		TslEnabled:                       false,
		StopLossPercentage:               "0",
		ErrorMessage:                     nullable.NewNullableWithValue(""),
		ProfitCurrency:                   "quote_currency",
		StopLossType:                     "stop_loss",
		SafetyOrderStepPercentage:        "0",
		TakeProfitType:                   "total",
		StopLossTimeoutEnabled:           false,
		StopLossTimeoutInSeconds:         0,
		AddFundable:                      false,
		SmartTradeConvertable:            false,
		PanicSellable:                    false,
		MarketType:                       "spot",
		OrderbookPriceCurrency:           "USDT",
		SafetyOrderVolume:                "10",
		SafetyOrderVolumeType:            "quote_currency",
		MartingaleStepCoefficient:        "1.0",
		MartingaleVolumeCoefficient:      "1.0",
		MinProfitPercentage:              "0",
		SafetyStrategyList:               []map[string]interface{}{},
		SlToBreakevenEnabled:             false,
		CloseStrategyList:                []map[string]interface{}{},
		TakeProfitSteps:                  []struct {
			AmountPercentage   *float32                      `json:"amount_percentage,omitempty"`
			Editable           *bool                         `json:"editable,omitempty"`
			ExecutionTimestamp nullable.Nullable[time.Time] `json:"execution_timestamp,omitempty"`
			Id                 *int                          `json:"id,omitempty"`
			InitialAmount      *string                       `json:"initial_amount,omitempty"`
			PanicSellable      *bool                         `json:"panic_sellable,omitempty"`
			Price              *string                       `json:"price,omitempty"`
			ProfitPercentage   *float32                      `json:"profit_percentage,omitempty"`
			Status             *string                       `json:"status,omitempty"`
			TradeId            *int                          `json:"trade_id,omitempty"`
		}{},
		Type:                             "simple",
	}
}

// AddBotEvent adds a bot event to a deal
// message: Human-readable event description
func AddBotEvent(deal *tcmock.Deal, message string) {
	now := time.Now()
	msg := message
	deal.BotEvents = append(deal.BotEvents, struct {
		CreatedAt *time.Time `json:"created_at,omitempty"`
		Message   *string    `json:"message,omitempty"`
	}{
		CreatedAt: &now,
		Message:   &msg,
	})
}

package model

import "time"

type BacktestListItem struct {
	ID          uint      `json:"id"`
	Symbol      string    `json:"symbol"`
	Timeframe   string    `json:"timeframe"`
	TotalReturn float64   `json:"total_return"`
	TradeCount  int       `json:"trade_count"`
	CreatedAt   time.Time `json:"created_at"`
}

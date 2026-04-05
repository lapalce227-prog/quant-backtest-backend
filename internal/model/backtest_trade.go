package model

import "time"

type BacktestTrade struct {
	BaseModel
	BacktestRunID uint      `gorm:"not null;index" json:"backtest_run_id"`
	Symbol        string    `gorm:"size:64;not null" json:"symbol"`
	Side          string    `gorm:"size:16;not null" json:"side"`
	EntryTime     time.Time `gorm:"not null" json:"entry_time"`
	EntryPrice    float64   `gorm:"type:decimal(20,8);not null" json:"entry_price"`
	ExitTime      time.Time `gorm:"not null" json:"exit_time"`
	ExitPrice     float64   `gorm:"type:decimal(20,8);not null" json:"exit_price"`
	Quantity      float64   `gorm:"type:decimal(20,8);not null" json:"quantity"`
	PnL           float64   `gorm:"column:pnl;type:decimal(20,8);not null" json:"pnl"`
	ReturnRate    float64   `gorm:"type:decimal(12,6);not null" json:"return_rate"`
	Commission    float64   `gorm:"type:decimal(20,8);not null" json:"commission"`
}

func (BacktestTrade) TableName() string {
	return "backtest_trades"
}

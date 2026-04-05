package model

import "time"

type BacktestRun struct {
	BaseModel
	StrategyID      uint       `gorm:"not null;index" json:"strategy_id"`
	MarketDatasetID uint       `gorm:"not null;index" json:"market_dataset_id"`
	Symbol          string     `gorm:"size:64;not null;comment:backtest config snapshot" json:"symbol"`
	Timeframe       string     `gorm:"size:32;not null;comment:backtest config snapshot" json:"timeframe"`
	StartTime       time.Time  `gorm:"not null" json:"start_time"`
	EndTime         time.Time  `gorm:"not null" json:"end_time"`
	InitialCapital  float64    `gorm:"type:decimal(20,8);not null" json:"initial_capital"`
	CommissionRate  float64    `gorm:"type:decimal(10,8);not null" json:"commission_rate"`
	SlippageRate    float64    `gorm:"type:decimal(10,8);not null" json:"slippage_rate"`
	RunStatus       string     `gorm:"size:32;not null;index" json:"run_status"`
	TotalReturn     float64    `gorm:"type:decimal(12,6)" json:"total_return"`
	MaxDrawdown     float64    `gorm:"type:decimal(12,6)" json:"max_drawdown"`
	WinRate         float64    `gorm:"type:decimal(12,6)" json:"win_rate"`
	SharpeRatio     float64    `gorm:"type:decimal(12,6)" json:"sharpe_ratio"`
	TradeCount      int        `gorm:"not null;default:0" json:"trade_count"`
	EquityCurve     string     `gorm:"type:longtext" json:"equity_curve"`
	Summary         string     `gorm:"type:json" json:"summary"`
	StartedAt       *time.Time `json:"started_at"`
	FinishedAt      *time.Time `json:"finished_at"`
}

func (BacktestRun) TableName() string {
	return "backtest_runs"
}

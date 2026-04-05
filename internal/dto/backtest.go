package dto

import "time"

type CreateBacktestRequest struct {
	StrategyID      uint      `json:"strategy_id" binding:"required"`
	MarketDatasetID uint      `json:"market_dataset_id" binding:"required"`
	StartTime       time.Time `json:"start_time" binding:"required"`
	EndTime         time.Time `json:"end_time" binding:"required"`
	InitialCapital  float64   `json:"initial_capital" binding:"required"`
	CommissionRate  float64   `json:"commission_rate" binding:"required"`
	SlippageRate    float64   `json:"slippage_rate" binding:"required"`
}

type ListBacktestsRequest struct {
	StrategyID      uint `form:"strategy_id"`
	MarketDatasetID uint `form:"market_dataset_id"`
}

type BacktestListItem struct {
	ID              uint      `json:"id"`
	StrategyID      uint      `json:"strategy_id"`
	MarketDatasetID uint      `json:"market_dataset_id"`
	Symbol          string    `json:"symbol"`
	Timeframe       string    `json:"timeframe"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	InitialCapital  float64   `json:"initial_capital"`
	RunStatus       string    `json:"run_status"`
	TotalReturn     float64   `json:"total_return"`
	MaxDrawdown     float64   `json:"max_drawdown"`
	WinRate         float64   `json:"win_rate"`
	SharpeRatio     float64   `json:"sharpe_ratio"`
	TradeCount      int       `json:"trade_count"`
	CreatedAt       time.Time `json:"created_at"`
}

type BacktestMetrics struct {
	TotalReturn float64 `json:"total_return"`
	MaxDrawdown float64 `json:"max_drawdown"`
	WinRate     float64 `json:"win_rate"`
	SharpeRatio float64 `json:"sharpe_ratio"`
	TradeCount  int     `json:"trade_count"`
	FinalEquity float64 `json:"final_equity"`
}

type BacktestConfigSnapshot struct {
	Symbol         string    `json:"symbol"`
	Timeframe      string    `json:"timeframe"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	InitialCapital float64   `json:"initial_capital"`
	CommissionRate float64   `json:"commission_rate"`
	SlippageRate   float64   `json:"slippage_rate"`
}

type EquityPointResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Equity    float64   `json:"equity"`
}

type BacktestDetailResponse struct {
	ID             uint                   `json:"id"`
	Strategy       StrategyResponse       `json:"strategy"`
	MarketDataset  MarketDatasetResponse  `json:"market_dataset"`
	ConfigSnapshot BacktestConfigSnapshot `json:"config_snapshot"`
	Metrics        BacktestMetrics        `json:"metrics"`
	EquityCurve    []EquityPointResponse  `json:"equity_curve"`
	StartedAt      *time.Time             `json:"started_at"`
	FinishedAt     *time.Time             `json:"finished_at"`
	CreatedAt      time.Time              `json:"created_at"`
}

type BacktestTradeResponse struct {
	ID         uint      `json:"id"`
	Symbol     string    `json:"symbol"`
	Side       string    `json:"side"`
	EntryTime  time.Time `json:"entry_time"`
	EntryPrice float64   `json:"entry_price"`
	ExitTime   time.Time `json:"exit_time"`
	ExitPrice  float64   `json:"exit_price"`
	Quantity   float64   `json:"quantity"`
	PnL        float64   `json:"pnl"`
	ReturnRate float64   `json:"return_rate"`
	Commission float64   `json:"commission"`
	CreatedAt  time.Time `json:"created_at"`
}

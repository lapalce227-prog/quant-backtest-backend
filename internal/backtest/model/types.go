package model

import "time"

type Candle struct {
	Timestamp time.Time `json:"timestamp"`
	Open      float64   `json:"open"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Close     float64   `json:"close"`
	Volume    float64   `json:"volume"`
}

type Signal string

const (
	SignalHold Signal = "hold"
	SignalBuy  Signal = "buy"
	SignalSell Signal = "sell"
)

type Position struct {
	Quantity   float64   `json:"quantity"`
	EntryPrice float64   `json:"entry_price"`
	EntryTime  time.Time `json:"entry_time"`
	Commission float64   `json:"commission"`
	IsOpen     bool      `json:"is_open"`
}

type Trade struct {
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
}

type EquityPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Equity    float64   `json:"equity"`
}

type MovingAverageCrossParams struct {
	ShortWindow int `json:"short_window"`
	LongWindow  int `json:"long_window"`
}

type Config struct {
	Symbol         string                   `json:"symbol"`
	Timeframe      string                   `json:"timeframe"`
	InitialCapital float64                  `json:"initial_capital"`
	CommissionRate float64                  `json:"commission_rate"`
	SlippageRate   float64                  `json:"slippage_rate"`
	StrategyParams MovingAverageCrossParams `json:"strategy_params"`
}

type RunSummary struct {
	InitialCapital float64 `json:"initial_capital"`
	FinalEquity    float64 `json:"final_equity"`
	TotalReturn    float64 `json:"total_return"`
	MaxDrawdown    float64 `json:"max_drawdown"`
	WinRate        float64 `json:"win_rate"`
	SharpeRatio    float64 `json:"sharpe_ratio"`
	TradeCount     int     `json:"trade_count"`
}

type Result struct {
	Summary     RunSummary    `json:"summary"`
	Trades      []Trade       `json:"trades"`
	EquityCurve []EquityPoint `json:"equity_curve"`
}

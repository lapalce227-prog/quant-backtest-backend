package repository

import "gorm.io/gorm"

type Repositories struct {
	Strategy      *StrategyRepository
	MarketDataset *MarketDatasetRepository
	BacktestRun   *BacktestRunRepository
	BacktestTrade *BacktestTradeRepository
}

func New(db *gorm.DB) *Repositories {
	return &Repositories{
		Strategy:      NewStrategyRepository(db),
		MarketDataset: NewMarketDatasetRepository(db),
		BacktestRun:   NewBacktestRunRepository(db),
		BacktestTrade: NewBacktestTradeRepository(db),
	}
}

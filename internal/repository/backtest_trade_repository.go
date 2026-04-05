package repository

import (
	"gorm.io/gorm"

	"tradingsystem/internal/model"
)

type BacktestTradeRepository struct {
	db *gorm.DB
}

func NewBacktestTradeRepository(db *gorm.DB) *BacktestTradeRepository {
	return &BacktestTradeRepository{db: db}
}

func (r *BacktestTradeRepository) DB() *gorm.DB {
	return r.db
}

func (r *BacktestTradeRepository) CreateBatch(trades []model.BacktestTrade) error {
	if len(trades) == 0 {
		return nil
	}

	return r.db.Create(&trades).Error
}

func (r *BacktestTradeRepository) ListByBacktestRunID(backtestRunID uint) ([]model.BacktestTrade, error) {
	var trades []model.BacktestTrade
	if err := r.db.
		Where("backtest_run_id = ?", backtestRunID).
		Order("entry_time asc").
		Find(&trades).Error; err != nil {
		return nil, err
	}

	return trades, nil
}

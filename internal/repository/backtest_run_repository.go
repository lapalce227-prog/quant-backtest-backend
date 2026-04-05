package repository

import (
	"gorm.io/gorm"

	"tradingsystem/internal/model"
)

type BacktestRunRepository struct {
	db *gorm.DB
}

func NewBacktestRunRepository(db *gorm.DB) *BacktestRunRepository {
	return &BacktestRunRepository{db: db}
}

func (r *BacktestRunRepository) DB() *gorm.DB {
	return r.db
}

func (r *BacktestRunRepository) Create(run *model.BacktestRun) error {
	return r.db.Create(run).Error
}

func (r *BacktestRunRepository) Save(run *model.BacktestRun) error {
	return r.db.Save(run).Error
}

func (r *BacktestRunRepository) GetByID(id uint) (*model.BacktestRun, error) {
	var run model.BacktestRun
	if err := r.db.First(&run, id).Error; err != nil {
		return nil, err
	}

	return &run, nil
}

func (r *BacktestRunRepository) List(strategyID uint, marketDatasetID uint) ([]model.BacktestRun, error) {
	query := r.db.Model(&model.BacktestRun{}).Order("id desc")

	if strategyID > 0 {
		query = query.Where("strategy_id = ?", strategyID)
	}
	if marketDatasetID > 0 {
		query = query.Where("market_dataset_id = ?", marketDatasetID)
	}

	var runs []model.BacktestRun
	if err := query.Find(&runs).Error; err != nil {
		return nil, err
	}

	return runs, nil
}

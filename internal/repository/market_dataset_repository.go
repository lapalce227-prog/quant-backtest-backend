package repository

import (
	"gorm.io/gorm"

	"tradingsystem/internal/model"
)

type MarketDatasetRepository struct {
	db *gorm.DB
}

func NewMarketDatasetRepository(db *gorm.DB) *MarketDatasetRepository {
	return &MarketDatasetRepository{db: db}
}

func (r *MarketDatasetRepository) DB() *gorm.DB {
	return r.db
}

func (r *MarketDatasetRepository) Create(dataset *model.MarketDataset) error {
	return r.db.Create(dataset).Error
}

func (r *MarketDatasetRepository) DeleteByID(id uint) error {
	return r.db.Delete(&model.MarketDataset{}, id).Error
}

func (r *MarketDatasetRepository) GetByID(id uint) (*model.MarketDataset, error) {
	var dataset model.MarketDataset
	if err := r.db.First(&dataset, id).Error; err != nil {
		return nil, err
	}

	return &dataset, nil
}

func (r *MarketDatasetRepository) List(symbol string, timeframe string) ([]model.MarketDataset, error) {
	query := r.db.Model(&model.MarketDataset{}).Order("id desc")

	if symbol != "" {
		query = query.Where("symbol = ?", symbol)
	}
	if timeframe != "" {
		query = query.Where("timeframe = ?", timeframe)
	}

	var datasets []model.MarketDataset
	if err := query.Find(&datasets).Error; err != nil {
		return nil, err
	}

	return datasets, nil
}

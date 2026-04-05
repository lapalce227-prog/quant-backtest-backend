package repository

import (
	"gorm.io/gorm"

	"tradingsystem/internal/model"
)

type StrategyRepository struct {
	db *gorm.DB
}

func NewStrategyRepository(db *gorm.DB) *StrategyRepository {
	return &StrategyRepository{db: db}
}

func (r *StrategyRepository) DB() *gorm.DB {
	return r.db
}

func (r *StrategyRepository) Create(strategy *model.Strategy) error {
	return r.db.Create(strategy).Error
}

func (r *StrategyRepository) Save(strategy *model.Strategy) error {
	return r.db.Save(strategy).Error
}

func (r *StrategyRepository) DeleteByID(id uint) error {
	return r.db.Delete(&model.Strategy{}, id).Error
}

func (r *StrategyRepository) GetByID(id uint) (*model.Strategy, error) {
	var strategy model.Strategy
	if err := r.db.First(&strategy, id).Error; err != nil {
		return nil, err
	}

	return &strategy, nil
}

func (r *StrategyRepository) List() ([]model.Strategy, error) {
	var strategies []model.Strategy
	if err := r.db.Order("id desc").Find(&strategies).Error; err != nil {
		return nil, err
	}

	return strategies, nil
}

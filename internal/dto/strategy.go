package dto

import "time"

type CreateStrategyRequest struct {
	Name         string                 `json:"name" binding:"required"`
	Description  string                 `json:"description"`
	StrategyType string                 `json:"strategy_type" binding:"required"`
	Parameters   map[string]interface{} `json:"parameters" binding:"required"`
}

type UpdateStrategyRequest struct {
	Name         string                 `json:"name" binding:"required"`
	Description  string                 `json:"description"`
	StrategyType string                 `json:"strategy_type" binding:"required"`
	Parameters   map[string]interface{} `json:"parameters" binding:"required"`
}

type StrategyResponse struct {
	ID           uint                   `json:"id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	StrategyType string                 `json:"strategy_type"`
	Parameters   map[string]interface{} `json:"parameters"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

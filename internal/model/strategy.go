package model

type Strategy struct {
	BaseModel
	Name         string `gorm:"size:128;not null;uniqueIndex" json:"name"`
	Description  string `gorm:"type:text" json:"description"`
	StrategyType string `gorm:"size:64;not null;index" json:"strategy_type"`
	Parameters   string `gorm:"type:json;not null" json:"parameters"`
}

func (Strategy) TableName() string {
	return "strategies"
}

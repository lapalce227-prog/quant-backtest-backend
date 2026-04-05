package model

import "time"

type MarketDataset struct {
	BaseModel
	Symbol      string    `gorm:"size:64;not null;index:idx_symbol_timeframe" json:"symbol"`
	Timeframe   string    `gorm:"size:32;not null;index:idx_symbol_timeframe" json:"timeframe"`
	SourceName  string    `gorm:"size:255;not null" json:"source_name"`
	StoragePath string    `gorm:"column:storage_path;size:512;not null" json:"storage_path"`
	StartTime   time.Time `gorm:"not null" json:"start_time"`
	EndTime     time.Time `gorm:"not null" json:"end_time"`
	RecordCount int64     `gorm:"not null" json:"record_count"`
}

func (MarketDataset) TableName() string {
	return "market_datasets"
}

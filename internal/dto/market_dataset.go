package dto

import "time"

type ListMarketDatasetsRequest struct {
	Symbol    string `form:"symbol"`
	Timeframe string `form:"timeframe"`
}

type MarketDatasetResponse struct {
	ID          uint      `json:"id"`
	Symbol      string    `json:"symbol"`
	Timeframe   string    `json:"timeframe"`
	SourceName  string    `json:"source_name"`
	StoragePath string    `json:"storage_path"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	RecordCount int64     `json:"record_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

package service

import (
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"tradingsystem/internal/dto"
	"tradingsystem/internal/model"
	"tradingsystem/internal/repository"
)

const datasetStorageDir = "data/market_datasets"

type MarketDatasetService struct {
	repo *repository.MarketDatasetRepository
}

func NewMarketDatasetService(repo *repository.MarketDatasetRepository) *MarketDatasetService {
	return &MarketDatasetService{repo: repo}
}

func (s *MarketDatasetService) Import(symbol string, timeframe string, fileHeader *multipart.FileHeader) (*dto.MarketDatasetResponse, error) {
	if err := os.MkdirAll(datasetStorageDir, 0o755); err != nil {
		return nil, err
	}

	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(fileHeader.Filename))
	storagePath := filepath.Join(datasetStorageDir, filename)

	dst, err := os.Create(storagePath)
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(dst, src); err != nil {
		dst.Close()
		return nil, err
	}
	if err := dst.Close(); err != nil {
		return nil, err
	}

	startTime, endTime, recordCount, err := inspectCSV(storagePath)
	if err != nil {
		_ = os.Remove(storagePath)
		return nil, err
	}

	dataset := &model.MarketDataset{
		Symbol:      symbol,
		Timeframe:   timeframe,
		SourceName:  fileHeader.Filename,
		StoragePath: storagePath,
		StartTime:   startTime,
		EndTime:     endTime,
		RecordCount: recordCount,
	}
	if err := s.repo.Create(dataset); err != nil {
		_ = os.Remove(storagePath)
		return nil, err
	}

	return toMarketDatasetResponse(dataset), nil
}

func (s *MarketDatasetService) List(symbol string, timeframe string) ([]dto.MarketDatasetResponse, error) {
	datasets, err := s.repo.List(symbol, timeframe)
	if err != nil {
		return nil, err
	}

	items := make([]dto.MarketDatasetResponse, 0, len(datasets))
	for _, dataset := range datasets {
		items = append(items, *toMarketDatasetResponse(&dataset))
	}

	return items, nil
}

func (s *MarketDatasetService) GetByID(id uint) (*dto.MarketDatasetResponse, error) {
	dataset, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return toMarketDatasetResponse(dataset), nil
}

func (s *MarketDatasetService) Delete(id uint) error {
	dataset, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if err := s.repo.DeleteByID(id); err != nil {
		return err
	}

	if dataset.StoragePath != "" {
		_ = os.Remove(dataset.StoragePath)
	}

	return nil
}

func toMarketDatasetResponse(dataset *model.MarketDataset) *dto.MarketDatasetResponse {
	return &dto.MarketDatasetResponse{
		ID:          dataset.ID,
		Symbol:      dataset.Symbol,
		Timeframe:   dataset.Timeframe,
		SourceName:  dataset.SourceName,
		StoragePath: dataset.StoragePath,
		StartTime:   dataset.StartTime,
		EndTime:     dataset.EndTime,
		RecordCount: dataset.RecordCount,
		CreatedAt:   dataset.CreatedAt,
		UpdatedAt:   dataset.UpdatedAt,
	}
}

func inspectCSV(path string) (time.Time, time.Time, int64, error) {
	file, err := os.Open(path)
	if err != nil {
		return time.Time{}, time.Time{}, 0, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	header, err := reader.Read()
	if err != nil {
		return time.Time{}, time.Time{}, 0, err
	}

	indexes, err := validateCSVHeader(header)
	if err != nil {
		return time.Time{}, time.Time{}, 0, err
	}

	var startTime time.Time
	var endTime time.Time
	var recordCount int64

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return time.Time{}, time.Time{}, 0, err
		}

		ts, err := parseTimestamp(record[indexes["timestamp"]])
		if err != nil {
			return time.Time{}, time.Time{}, 0, err
		}

		if _, err := strconv.ParseFloat(record[indexes["open"]], 64); err != nil {
			return time.Time{}, time.Time{}, 0, err
		}
		if _, err := strconv.ParseFloat(record[indexes["high"]], 64); err != nil {
			return time.Time{}, time.Time{}, 0, err
		}
		if _, err := strconv.ParseFloat(record[indexes["low"]], 64); err != nil {
			return time.Time{}, time.Time{}, 0, err
		}
		if _, err := strconv.ParseFloat(record[indexes["close"]], 64); err != nil {
			return time.Time{}, time.Time{}, 0, err
		}
		if _, err := strconv.ParseFloat(record[indexes["volume"]], 64); err != nil {
			return time.Time{}, time.Time{}, 0, err
		}

		if recordCount == 0 {
			startTime = ts
		}
		endTime = ts
		recordCount++
	}

	if recordCount == 0 {
		return time.Time{}, time.Time{}, 0, fmt.Errorf("csv has no data rows")
	}

	return startTime, endTime, recordCount, nil
}

func validateCSVHeader(header []string) (map[string]int, error) {
	required := []string{"timestamp", "open", "high", "low", "close", "volume"}
	indexes := make(map[string]int, len(required))

	for i, column := range header {
		indexes[column] = i
	}

	for _, column := range required {
		if _, ok := indexes[column]; !ok {
			return nil, fmt.Errorf("missing required column: %s", column)
		}
	}

	return indexes, nil
}

func parseTimestamp(value string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, layout := range layouts {
		if ts, err := time.Parse(layout, value); err == nil {
			return ts, nil
		}
	}

	return time.Time{}, fmt.Errorf("unsupported timestamp format: %s", value)
}

func ParseTimestamp(value string) (time.Time, error) {
	return parseTimestamp(value)
}

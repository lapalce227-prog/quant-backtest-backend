package service

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"

	btengine "tradingsystem/internal/backtest/engine"
	btmodel "tradingsystem/internal/backtest/model"
	"tradingsystem/internal/dto"
	"tradingsystem/internal/model"
	"tradingsystem/internal/repository"
)

type BacktestService struct {
	strategyRepo      *repository.StrategyRepository
	marketDatasetRepo *repository.MarketDatasetRepository
	backtestRunRepo   *repository.BacktestRunRepository
	backtestTradeRepo *repository.BacktestTradeRepository
}

func NewBacktestService(repos *repository.Repositories) *BacktestService {
	return &BacktestService{
		strategyRepo:      repos.Strategy,
		marketDatasetRepo: repos.MarketDataset,
		backtestRunRepo:   repos.BacktestRun,
		backtestTradeRepo: repos.BacktestTrade,
	}
}

func (s *BacktestService) Create(req dto.CreateBacktestRequest) (*dto.BacktestListItem, error) {
	if req.InitialCapital <= 0 {
		return nil, fmt.Errorf("initial_capital must be greater than zero")
	}
	if req.EndTime.Before(req.StartTime) {
		return nil, fmt.Errorf("end_time must be greater than or equal to start_time")
	}

	strategyEntity, err := s.strategyRepo.GetByID(req.StrategyID)
	if err != nil {
		return nil, err
	}
	marketDatasetEntity, err := s.marketDatasetRepo.GetByID(req.MarketDatasetID)
	if err != nil {
		return nil, err
	}

	params, err := parseMovingAverageCrossParams(strategyEntity)
	if err != nil {
		return nil, err
	}

	candles, err := btengine.LoadCandlesFromCSV(marketDatasetEntity.StoragePath)
	if err != nil {
		return nil, err
	}
	filteredCandles := filterCandlesByTimeRange(candles, req.StartTime, req.EndTime)
	if len(filteredCandles) == 0 {
		return nil, fmt.Errorf("no candles found in selected time range")
	}

	engine := btengine.NewMovingAverageCrossEngine(params)
	result, err := engine.Run(filteredCandles, btmodel.Config{
		Symbol:         marketDatasetEntity.Symbol,
		Timeframe:      marketDatasetEntity.Timeframe,
		InitialCapital: req.InitialCapital,
		CommissionRate: req.CommissionRate,
		SlippageRate:   req.SlippageRate,
		StrategyParams: params,
	})
	if err != nil {
		return nil, err
	}

	now := time.Now()
	summaryJSON, err := json.Marshal(result.Summary)
	if err != nil {
		return nil, err
	}
	equityCurveJSON, err := json.Marshal(toEquityPointResponses(result.EquityCurve))
	if err != nil {
		return nil, err
	}

	run := &model.BacktestRun{
		StrategyID:      strategyEntity.ID,
		MarketDatasetID: marketDatasetEntity.ID,
		Symbol:          marketDatasetEntity.Symbol,
		Timeframe:       marketDatasetEntity.Timeframe,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		InitialCapital:  req.InitialCapital,
		CommissionRate:  req.CommissionRate,
		SlippageRate:    req.SlippageRate,
		RunStatus:       "completed",
		TotalReturn:     result.Summary.TotalReturn,
		MaxDrawdown:     result.Summary.MaxDrawdown,
		WinRate:         result.Summary.WinRate,
		SharpeRatio:     result.Summary.SharpeRatio,
		TradeCount:      result.Summary.TradeCount,
		EquityCurve:     string(equityCurveJSON),
		Summary:         string(summaryJSON),
		StartedAt:       &now,
		FinishedAt:      &now,
	}

	db := s.backtestRunRepo.DB()
	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(run).Error; err != nil {
			return err
		}

		trades := toBacktestTradeModels(run.ID, result.Trades)
		if len(trades) > 0 {
			if err := tx.Create(&trades).Error; err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &dto.BacktestListItem{
		ID:              run.ID,
		StrategyID:      run.StrategyID,
		MarketDatasetID: run.MarketDatasetID,
		Symbol:          run.Symbol,
		Timeframe:       run.Timeframe,
		StartTime:       run.StartTime,
		EndTime:         run.EndTime,
		InitialCapital:  run.InitialCapital,
		RunStatus:       run.RunStatus,
		TotalReturn:     run.TotalReturn,
		MaxDrawdown:     run.MaxDrawdown,
		WinRate:         run.WinRate,
		SharpeRatio:     run.SharpeRatio,
		TradeCount:      run.TradeCount,
		CreatedAt:       run.CreatedAt,
	}, nil
}

func (s *BacktestService) List() ([]model.BacktestListItem, error) {
	items := make([]model.BacktestListItem, 0)

	err := s.backtestRunRepo.DB().
		Model(&model.BacktestRun{}).
		Select("id, symbol, timeframe, total_return, trade_count, created_at").
		Order("created_at desc").
		Find(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s *BacktestService) GetByID(id uint) (*dto.BacktestDetailResponse, error) {
	run, err := s.backtestRunRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	strategyEntity, err := s.strategyRepo.GetByID(run.StrategyID)
	if err != nil {
		return nil, err
	}
	marketDatasetEntity, err := s.marketDatasetRepo.GetByID(run.MarketDatasetID)
	if err != nil {
		return nil, err
	}

	strategyDTO, err := toStrategyResponse(strategyEntity)
	if err != nil {
		return nil, err
	}

	var equityCurve []dto.EquityPointResponse
	if run.EquityCurve != "" {
		if err := json.Unmarshal([]byte(run.EquityCurve), &equityCurve); err != nil {
			return nil, err
		}
	}

	return &dto.BacktestDetailResponse{
		ID:            run.ID,
		Strategy:      *strategyDTO,
		MarketDataset: *toMarketDatasetResponse(marketDatasetEntity),
		ConfigSnapshot: dto.BacktestConfigSnapshot{
			Symbol:         run.Symbol,
			Timeframe:      run.Timeframe,
			StartTime:      run.StartTime,
			EndTime:        run.EndTime,
			InitialCapital: run.InitialCapital,
			CommissionRate: run.CommissionRate,
			SlippageRate:   run.SlippageRate,
		},
		Metrics: dto.BacktestMetrics{
			TotalReturn: run.TotalReturn,
			MaxDrawdown: run.MaxDrawdown,
			WinRate:     run.WinRate,
			SharpeRatio: run.SharpeRatio,
			TradeCount:  run.TradeCount,
			FinalEquity: extractFinalEquity(run.Summary),
		},
		EquityCurve: equityCurve,
		StartedAt:   run.StartedAt,
		FinishedAt:  run.FinishedAt,
		CreatedAt:   run.CreatedAt,
	}, nil
}

func (s *BacktestService) ListTrades(backtestRunID uint) ([]dto.BacktestTradeResponse, error) {
	if _, err := s.backtestRunRepo.GetByID(backtestRunID); err != nil {
		return nil, err
	}

	trades, err := s.backtestTradeRepo.ListByBacktestRunID(backtestRunID)
	if err != nil {
		return nil, err
	}

	items := make([]dto.BacktestTradeResponse, 0, len(trades))
	for _, trade := range trades {
		items = append(items, dto.BacktestTradeResponse{
			ID:         trade.ID,
			Symbol:     trade.Symbol,
			Side:       trade.Side,
			EntryTime:  trade.EntryTime,
			EntryPrice: trade.EntryPrice,
			ExitTime:   trade.ExitTime,
			ExitPrice:  trade.ExitPrice,
			Quantity:   trade.Quantity,
			PnL:        trade.PnL,
			ReturnRate: trade.ReturnRate,
			Commission: trade.Commission,
			CreatedAt:  trade.CreatedAt,
		})
	}

	return items, nil
}

func parseMovingAverageCrossParams(strategyEntity *model.Strategy) (btmodel.MovingAverageCrossParams, error) {
	if strategyEntity.StrategyType != "moving_average_cross" {
		return btmodel.MovingAverageCrossParams{}, fmt.Errorf("unsupported strategy_type: %s", strategyEntity.StrategyType)
	}

	var params btmodel.MovingAverageCrossParams
	if err := json.Unmarshal([]byte(strategyEntity.Parameters), &params); err != nil {
		return btmodel.MovingAverageCrossParams{}, err
	}
	if params.ShortWindow <= 0 || params.LongWindow <= 0 || params.ShortWindow >= params.LongWindow {
		return btmodel.MovingAverageCrossParams{}, fmt.Errorf("invalid moving average parameters")
	}

	return params, nil
}

func filterCandlesByTimeRange(candles []btmodel.Candle, start time.Time, end time.Time) []btmodel.Candle {
	filtered := make([]btmodel.Candle, 0, len(candles))
	for _, candle := range candles {
		if candle.Timestamp.Before(start) || candle.Timestamp.After(end) {
			continue
		}
		filtered = append(filtered, candle)
	}

	return filtered
}

func toEquityPointResponses(points []btmodel.EquityPoint) []dto.EquityPointResponse {
	items := make([]dto.EquityPointResponse, 0, len(points))
	for _, point := range points {
		items = append(items, dto.EquityPointResponse{
			Timestamp: point.Timestamp,
			Equity:    point.Equity,
		})
	}

	return items
}

func extractFinalEquity(summaryJSON string) float64 {
	if summaryJSON == "" {
		return 0
	}

	var summary btmodel.RunSummary
	if err := json.Unmarshal([]byte(summaryJSON), &summary); err != nil {
		return 0
	}

	return summary.FinalEquity
}

func toBacktestTradeModels(backtestRunID uint, trades []btmodel.Trade) []model.BacktestTrade {
	items := make([]model.BacktestTrade, 0, len(trades))
	for _, trade := range trades {
		items = append(items, model.BacktestTrade{
			BacktestRunID: backtestRunID,
			Symbol:        trade.Symbol,
			Side:          trade.Side,
			EntryTime:     trade.EntryTime,
			EntryPrice:    trade.EntryPrice,
			ExitTime:      trade.ExitTime,
			ExitPrice:     trade.ExitPrice,
			Quantity:      trade.Quantity,
			PnL:           trade.PnL,
			ReturnRate:    trade.ReturnRate,
			Commission:    trade.Commission,
		})
	}

	return items
}

package engine

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	btmetrics "tradingsystem/internal/backtest/metrics"
	btmodel "tradingsystem/internal/backtest/model"
	btstrategy "tradingsystem/internal/backtest/strategy"
)

type Engine struct {
	strategy btstrategy.Strategy
}

func NewMovingAverageCrossEngine(params btmodel.MovingAverageCrossParams) *Engine {
	return &Engine{
		strategy: btstrategy.NewMovingAverageCrossStrategy(params),
	}
}

func (e *Engine) RunCSV(path string, config btmodel.Config) (*btmodel.Result, error) {
	candles, err := LoadCandlesFromCSV(path)
	if err != nil {
		return nil, err
	}

	return e.Run(candles, config)
}

func (e *Engine) Run(candles []btmodel.Candle, config btmodel.Config) (*btmodel.Result, error) {
	if e.strategy == nil {
		return nil, fmt.Errorf("strategy is required")
	}
	if len(candles) == 0 {
		return nil, fmt.Errorf("candles are required")
	}
	if config.InitialCapital <= 0 {
		return nil, fmt.Errorf("initial capital must be greater than zero")
	}

	cash := config.InitialCapital
	position := btmodel.Position{}
	trades := make([]btmodel.Trade, 0)
	equityCurve := make([]btmodel.EquityPoint, 0, len(candles))

	for i, candle := range candles {
		signal := e.strategy.Signal(candles, i)

		if signal == btmodel.SignalBuy && !position.IsOpen {
			position = openLongPosition(candle, cash, config.CommissionRate, config.SlippageRate)
			if position.IsOpen {
				cash = 0
			}
		}

		if signal == btmodel.SignalSell && position.IsOpen {
			trade, newCash := closeLongPosition(config.Symbol, candle, position, config.CommissionRate, config.SlippageRate)
			trades = append(trades, trade)
			cash = newCash
			position = btmodel.Position{}
		}

		equityCurve = append(equityCurve, btmodel.EquityPoint{
			Timestamp: candle.Timestamp,
			Equity:    markToMarket(cash, position, candle.Close),
		})
	}

	if position.IsOpen {
		lastCandle := candles[len(candles)-1]
		trade, newCash := closeLongPosition(config.Symbol, lastCandle, position, config.CommissionRate, config.SlippageRate)
		trades = append(trades, trade)
		cash = newCash
		position = btmodel.Position{}
		equityCurve[len(equityCurve)-1].Equity = cash
	}

	summary := btmetrics.BuildSummary(config.InitialCapital, trades, equityCurve)

	return &btmodel.Result{
		Summary:     summary,
		Trades:      trades,
		EquityCurve: equityCurve,
	}, nil
}

func LoadCandlesFromCSV(path string) ([]btmodel.Candle, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	header, err := reader.Read()
	if err != nil {
		return nil, err
	}

	indexes, err := validateHeader(header)
	if err != nil {
		return nil, err
	}

	candles := make([]btmodel.Candle, 0)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		candle, err := parseCandle(record, indexes)
		if err != nil {
			return nil, err
		}
		candles = append(candles, candle)
	}

	if len(candles) == 0 {
		return nil, fmt.Errorf("csv has no data rows")
	}

	return candles, nil
}

func openLongPosition(candle btmodel.Candle, cash float64, commissionRate float64, slippageRate float64) btmodel.Position {
	entryPrice := candle.Close * (1 + slippageRate)
	if entryPrice <= 0 || cash <= 0 {
		return btmodel.Position{}
	}

	quantity := cash / (entryPrice * (1 + commissionRate))
	if quantity <= 0 {
		return btmodel.Position{}
	}

	grossAmount := quantity * entryPrice
	commission := grossAmount * commissionRate

	return btmodel.Position{
		Quantity:   quantity,
		EntryPrice: entryPrice,
		EntryTime:  candle.Timestamp,
		Commission: commission,
		IsOpen:     true,
	}
}

func closeLongPosition(symbol string, candle btmodel.Candle, position btmodel.Position, commissionRate float64, slippageRate float64) (btmodel.Trade, float64) {
	exitPrice := candle.Close * (1 - slippageRate)
	grossProceeds := position.Quantity * exitPrice
	exitCommission := grossProceeds * commissionRate
	netProceeds := grossProceeds - exitCommission

	entryCost := position.Quantity * position.EntryPrice
	totalCommission := position.Commission + exitCommission
	pnl := netProceeds - entryCost - position.Commission
	returnRate := 0.0
	if entryCost > 0 {
		returnRate = pnl / entryCost
	}

	return btmodel.Trade{
		Symbol:     symbol,
		Side:       "long",
		EntryTime:  position.EntryTime,
		EntryPrice: position.EntryPrice,
		ExitTime:   candle.Timestamp,
		ExitPrice:  exitPrice,
		Quantity:   position.Quantity,
		PnL:        pnl,
		ReturnRate: returnRate,
		Commission: totalCommission,
	}, netProceeds
}

func markToMarket(cash float64, position btmodel.Position, currentClose float64) float64 {
	if !position.IsOpen {
		return cash
	}

	return position.Quantity * currentClose
}

func validateHeader(header []string) (map[string]int, error) {
	required := []string{"timestamp", "open", "high", "low", "close", "volume"}
	indexes := make(map[string]int, len(header))
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

func parseCandle(record []string, indexes map[string]int) (btmodel.Candle, error) {
	timestamp, err := parseTimestamp(record[indexes["timestamp"]])
	if err != nil {
		return btmodel.Candle{}, err
	}

	openPrice, err := strconv.ParseFloat(record[indexes["open"]], 64)
	if err != nil {
		return btmodel.Candle{}, err
	}
	highPrice, err := strconv.ParseFloat(record[indexes["high"]], 64)
	if err != nil {
		return btmodel.Candle{}, err
	}
	lowPrice, err := strconv.ParseFloat(record[indexes["low"]], 64)
	if err != nil {
		return btmodel.Candle{}, err
	}
	closePrice, err := strconv.ParseFloat(record[indexes["close"]], 64)
	if err != nil {
		return btmodel.Candle{}, err
	}
	volume, err := strconv.ParseFloat(record[indexes["volume"]], 64)
	if err != nil {
		return btmodel.Candle{}, err
	}

	return btmodel.Candle{
		Timestamp: timestamp,
		Open:      openPrice,
		High:      highPrice,
		Low:       lowPrice,
		Close:     closePrice,
		Volume:    volume,
	}, nil
}

func parseTimestamp(value string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf("unsupported timestamp format: %s", value)
}

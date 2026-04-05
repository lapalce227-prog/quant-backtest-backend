package strategy

import btmodel "tradingsystem/internal/backtest/model"

type Strategy interface {
	Signal(candles []btmodel.Candle, index int) btmodel.Signal
}

type MovingAverageCrossStrategy struct {
	shortWindow int
	longWindow  int
}

func NewMovingAverageCrossStrategy(params btmodel.MovingAverageCrossParams) *MovingAverageCrossStrategy {
	return &MovingAverageCrossStrategy{
		shortWindow: params.ShortWindow,
		longWindow:  params.LongWindow,
	}
}

func (s *MovingAverageCrossStrategy) Signal(candles []btmodel.Candle, index int) btmodel.Signal {
	if index <= 0 || index >= len(candles) {
		return btmodel.SignalHold
	}
	if s.shortWindow <= 0 || s.longWindow <= 0 || s.shortWindow >= s.longWindow {
		return btmodel.SignalHold
	}
	if index < s.longWindow {
		return btmodel.SignalHold
	}

	prevShort := simpleMovingAverage(candles, index-1, s.shortWindow)
	prevLong := simpleMovingAverage(candles, index-1, s.longWindow)
	currShort := simpleMovingAverage(candles, index, s.shortWindow)
	currLong := simpleMovingAverage(candles, index, s.longWindow)

	if prevShort <= prevLong && currShort > currLong {
		return btmodel.SignalBuy
	}
	if prevShort >= prevLong && currShort < currLong {
		return btmodel.SignalSell
	}

	return btmodel.SignalHold
}

func simpleMovingAverage(candles []btmodel.Candle, endIndex int, window int) float64 {
	start := endIndex - window + 1
	if start < 0 {
		return 0
	}

	sum := 0.0
	for i := start; i <= endIndex; i++ {
		sum += candles[i].Close
	}

	return sum / float64(window)
}

package metrics

import (
	"math"

	btmodel "tradingsystem/internal/backtest/model"
)

func BuildSummary(initialCapital float64, trades []btmodel.Trade, equityCurve []btmodel.EquityPoint) btmodel.RunSummary {
	finalEquity := initialCapital
	if len(equityCurve) > 0 {
		finalEquity = equityCurve[len(equityCurve)-1].Equity
	}

	return btmodel.RunSummary{
		InitialCapital: initialCapital,
		FinalEquity:    finalEquity,
		TotalReturn:    calcTotalReturn(initialCapital, finalEquity),
		MaxDrawdown:    calcMaxDrawdown(equityCurve),
		WinRate:        calcWinRate(trades),
		SharpeRatio:    calcSharpeRatio(equityCurve),
		TradeCount:     len(trades),
	}
}

func calcTotalReturn(initialCapital float64, finalEquity float64) float64 {
	if initialCapital == 0 {
		return 0
	}
	return (finalEquity - initialCapital) / initialCapital
}

func calcMaxDrawdown(equityCurve []btmodel.EquityPoint) float64 {
	if len(equityCurve) == 0 {
		return 0
	}

	peak := equityCurve[0].Equity
	maxDrawdown := 0.0

	for _, point := range equityCurve {
		if point.Equity > peak {
			peak = point.Equity
		}
		if peak == 0 {
			continue
		}

		drawdown := (peak - point.Equity) / peak
		if drawdown > maxDrawdown {
			maxDrawdown = drawdown
		}
	}

	return maxDrawdown
}

func calcWinRate(trades []btmodel.Trade) float64 {
	if len(trades) == 0 {
		return 0
	}

	wins := 0
	for _, trade := range trades {
		if trade.PnL > 0 {
			wins++
		}
	}

	return float64(wins) / float64(len(trades))
}

func calcSharpeRatio(equityCurve []btmodel.EquityPoint) float64 {
	if len(equityCurve) < 2 {
		return 0
	}

	returns := make([]float64, 0, len(equityCurve)-1)
	for i := 1; i < len(equityCurve); i++ {
		prev := equityCurve[i-1].Equity
		curr := equityCurve[i].Equity
		if prev == 0 {
			continue
		}
		returns = append(returns, (curr-prev)/prev)
	}

	if len(returns) == 0 {
		return 0
	}

	mean := mean(returns)
	stddev := standardDeviation(returns, mean)
	if stddev == 0 {
		return 0
	}

	return mean / stddev * math.Sqrt(252)
}

func mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	sum := 0.0
	for _, value := range values {
		sum += value
	}

	return sum / float64(len(values))
}

func standardDeviation(values []float64, avg float64) float64 {
	if len(values) < 2 {
		return 0
	}

	sum := 0.0
	for _, value := range values {
		diff := value - avg
		sum += diff * diff
	}

	variance := sum / float64(len(values)-1)
	return math.Sqrt(variance)
}

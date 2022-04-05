package perf

import (
	"math"

	"github.com/colngroup/zero2algo/broker"
	"github.com/gonum/floats"
)

type TradeReport struct {
	TradeCount     float64
	TotalNetProfit float64
	AvgNetProfit   float64
	GrossProfit    float64
	GrossLoss      float64
	ProfitFactor   float64

	PercentProfitable  float64
	MaxProfit, MaxLoss float64

	AvgProfit float64
	AvgLoss   float64

	MaxLossStreak int

	Kelly    float64
	OptimalF float64

	TotalTimeInMarketSec float64
	AvgHoldSec           float64

	winningCount, winningPct float64
	losingCount, losingPct   float64
}

func NewTradeReport(trades []broker.Trade) *TradeReport {
	if len(trades) == 0 {
		return nil
	}

	var report TradeReport
	var lossStreak int

	profits := make([]float64, len(trades))

	for i := range trades {
		t := trades[i]
		report.TotalTimeInMarketSec += t.HoldPeriod.Seconds()

		profit := t.Profit.InexactFloat64()
		profits[i] = profit
		switch {
		case profit > 0:
			report.winningCount++
			report.GrossProfit += profit
			if lossStreak > report.MaxLossStreak {
				report.MaxLossStreak = lossStreak
			}
			lossStreak = 0
		case profit < 0:
			report.losingCount++
			report.GrossLoss += math.Abs(profit)
			lossStreak++
		}
	}
	report.MaxProfit = floats.Max(profits)
	report.MaxLoss = math.Abs(floats.Min(profits))

	report.TradeCount = report.winningCount + report.losingCount

	report.TotalNetProfit = report.GrossProfit - report.GrossLoss
	report.AvgNetProfit = report.TotalNetProfit / report.TradeCount
	report.ProfitFactor = report.GrossProfit / NNZ(report.GrossLoss, 1)

	report.AvgProfit = report.GrossProfit / NNZ(report.winningCount, 1)
	report.AvgLoss = report.GrossLoss / NNZ(report.losingCount, 1)

	report.winningPct = report.winningCount / report.TradeCount
	report.losingPct = 1 - report.winningPct
	report.PercentProfitable = report.winningPct

	report.AvgHoldSec = report.TotalTimeInMarketSec / report.TradeCount

	report.Kelly = KellyCriterion(report.ProfitFactor, report.winningPct)
	report.OptimalF = OptimalF(profits)

	return &report
}

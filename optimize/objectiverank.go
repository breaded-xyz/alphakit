package optimize

import (
	"github.com/colngroup/zero2algo/perf"
)

type ObjectiveRanker func(a, b perf.PerformanceReport) bool

func SharpeRanker(a, b perf.PerformanceReport) bool {
	return a.Portfolio.Sharpe < b.Portfolio.Sharpe
}

func PRRRanker(a, b perf.PerformanceReport) bool {
	return a.Trade.PRR < b.Trade.PRR
}

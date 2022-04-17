package optimize

import (
	"github.com/colngroup/zero2algo/perf"
	"golang.org/x/exp/slices"
)

func SharpeSort(reports []perf.PerformanceReport) {
	slices.SortFunc(reports, func(a, b perf.PerformanceReport) bool {
		return a.Portfolio.Sharpe < b.Portfolio.Sharpe
	})
}

func PRRSort(reports []perf.PerformanceReport) {
	slices.SortFunc(reports, func(a, b perf.PerformanceReport) bool {
		return a.Trade.PRR < b.Trade.PRR
	})
}

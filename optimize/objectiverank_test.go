package optimize

import (
	"testing"

	"github.com/colngroup/zero2algo/perf"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

func TestSharpeRanker(t *testing.T) {
	give := []perf.PerformanceReport{
		{Portfolio: &perf.PortfolioReport{Sharpe: 2}},
		{Portfolio: &perf.PortfolioReport{Sharpe: 0.9}},
		{Portfolio: &perf.PortfolioReport{Sharpe: 2.5}},
	}

	want := []perf.PerformanceReport{give[1], give[0], give[2]}

	slices.SortFunc(give, SharpeRanker)

	assert.Equal(t, want, give)
}

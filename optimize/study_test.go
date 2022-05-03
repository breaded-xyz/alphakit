package optimize

import (
	"testing"

	"github.com/colngroup/zero2algo/perf"
	"github.com/stretchr/testify/assert"
)

func TestSummarize(t *testing.T) {

	give := Report{
		Backtests: []perf.PerformanceReport{
			{
				Trade:     &perf.TradeReport{PRR: 2.0},
				Portfolio: &perf.PortfolioReport{CAGR: 0.8, Sharpe: 1.0, Calmar: 2.0},
			},
			{
				Trade:     &perf.TradeReport{PRR: 4.0},
				Portfolio: &perf.PortfolioReport{CAGR: 1.5, Sharpe: 2.0, Calmar: 2.0},
			},
		},
	}
	want := Report{
		PRR:     3.0,
		CAGR:    1.15,
		Sharpe:  1.5,
		Calmar:  2.0,
		Backtests: []perf.PerformanceReport{
			{
				Trade:     &perf.TradeReport{PRR: 2.0},
				Portfolio: &perf.PortfolioReport{CAGR: 0.8, Sharpe: 1.0, Calmar: 2.0},
			},
			{
				Trade:     &perf.TradeReport{PRR: 4.0},
				Portfolio: &perf.PortfolioReport{CAGR: 1.5, Sharpe: 2.0, Calmar: 2.0},
			},
		},
	}

	act := Summarize(give)
	assert.Equal(t, want, act)
}

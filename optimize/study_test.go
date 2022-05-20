package optimize

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thecolngroup/alphakit/perf"
)

func TestSummarize(t *testing.T) {

	give := Report{
		Backtests: []perf.PerformanceReport{
			{
				Trade:     &perf.TradeReport{PRR: 2.0, TradeCount: 5},
				Portfolio: &perf.PortfolioReport{MaxDrawdown: 0.3, CAGR: 0.8, Sharpe: 1.0, Calmar: 2.0},
			},
			{
				Trade:     &perf.TradeReport{PRR: 4.0, TradeCount: 10},
				Portfolio: &perf.PortfolioReport{MaxDrawdown: 0.2, CAGR: 1.5, Sharpe: 2.0, Calmar: 2.0},
			},
		},
	}
	want := Report{
		PRR:         3.0,
		MDD:         0.25,
		CAGR:        1.15,
		Sharpe:      1.5,
		Calmar:      2.0,
		SampleCount: 2,
		TradeCount:  15,
		Backtests: []perf.PerformanceReport{
			{
				Trade:     &perf.TradeReport{PRR: 2.0, TradeCount: 5},
				Portfolio: &perf.PortfolioReport{MaxDrawdown: 0.3, CAGR: 0.8, Sharpe: 1.0, Calmar: 2.0},
			},
			{
				Trade:     &perf.TradeReport{PRR: 4.0, TradeCount: 10},
				Portfolio: &perf.PortfolioReport{MaxDrawdown: 0.2, CAGR: 1.5, Sharpe: 2.0, Calmar: 2.0},
			},
		},
	}

	act := Summarize(give)
	assert.Equal(t, want, act)
}

package optimize

import (
	"github.com/colngroup/zero2algo/internal/util"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/perf"
)

// Study is an optimization experiment, prepared and executed by an Optimizer.
// The subject is a trading algo and its parameter space.
// Hypothesis: the optimized algo params will generate positive market returns in live trading
// Null hypothesis: algo has zero positive expectancy of returns in the tested param space
// Independent variable (aka predictor / feature): algo parameter space defined by []ParamSet
// Dependent variable: algo performance (Sharpe, CAGR et al) measured by Report
// Control variables: price data sample, backtest simulator settings et al
type Study struct {
	TrainingPSets   []ParamSet
	TrainingSamples [][]market.Kline
	TrainingResults map[ParamSetID]Report

	ValidationPSets   []ParamSet
	ValidationSamples [][]market.Kline
	ValidationResults map[ParamSetID]Report
}

type ParamSetID string

type ParamSet struct {
	ID     ParamSetID
	Params map[string]any
}

func NewParamSet() ParamSet {
	return ParamSet{ID: ParamSetID(util.NewID())}
}

type Report struct {
	Subject ParamSet

	PRR    float64
	CAGR   float64
	Sharpe float64
	Calmar float64

	Backtests []perf.PerformanceReport
}

func (r *Report) AddResult(backtest ...perf.PerformanceReport) {
	r.Backtests = append(r.Backtests, backtest...)

	for i := range r.Backtests {
		r.PRR += r.Backtests[i].Trade.PRR
		r.CAGR += r.Backtests[i].Portfolio.CAGR
		r.Sharpe += r.Backtests[i].Portfolio.Sharpe
		r.Calmar += r.Backtests[i].Portfolio.Calmar
	}

	count := float64(len(r.Backtests))
	r.PRR /= count
	r.CAGR /= count
	r.Sharpe /= count
	r.Calmar /= count
}

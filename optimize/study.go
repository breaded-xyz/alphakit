package optimize

import (
	"github.com/colngroup/zero2algo/internal/util"
	"github.com/colngroup/zero2algo/market"
)

// Study is an optimization experiment. The subject is a trading algo and its parameter space.
// Hypothesis: the optimized algo params will generate positive market returns in live trading
// Null hypothesis: algo has zero positive expectancy of returns in the tested param space
// Independent variable (aka predictor / feature): algo parameter space
// Dependent variable: algo performance (Sharpe, CAGR et al)
// Control variables: price data sample, backtest simulator settings etc
// Various optimization methods exist in litreature.
// Included is an example method for a 'brute-force peak objective' approach which
// tests all given parameter combinations and selects the highest ranked (peak) param set.
// Method:
//
// Setup:
// - Generate 1 or more price data samples (inc. over/under sampling)
// - Split sample into in-sample (training) and out-of-sample (validation) datasets
// - Generate 1 or more param sets using the cartesian product of given ranges that define the param space
//
// Training:
// - Execute each algo param set over the in-sample price data
// - Average the performance for each param set over the in-sample data
// - Rank the param sets based on the performance objective (Profit Factor, T-Score etc)
// Optimization:
// - Execute the highest ranked ("trained") algo param over the out-of-sample price data
// - Accept or reject the hypothesis based on statistical significance of the performance report
type Study struct {
	Samples   [][]market.Kline
	ParamSets map[ParamSetID]ParamSet
	InSample  map[ParamSetID]ParamSetReport
	OutSample map[ParamSetID]ParamSetReport
}

type ParamSet struct {
	ID     ParamSetID
	Params map[string]any
}

type ParamSetID string

type ParamSetReport struct {
	CAGR   float64
	PRR    float64
	Sharpe float64
	Calmar float64
}

func NewParamSet() ParamSet {
	return ParamSet{ID: ParamSetID(util.NewID())}
}

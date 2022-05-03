package optimize

import (
	"github.com/colngroup/zero2algo/internal/util"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/perf"
)

// Study is an optimization experiment, prepared and executed by an Optimizer.
// First, a training (in-sample) phase is conducted, followed by a validation (out-of-sample) phase.
// The validation phase should report the out-of-sample (OOS) performance of the optimum param set.
//
// The experiment can be summarised as:
//
// - Hypothesis: the optimized algo params will generate positive market returns in live trading.
//
// - Null hypothesis: algo has zero positive expectancy of returns in the tested param space.
//
// - Independent variable (aka predictor / feature): algo parameter space defined by []ParamSet.
//
// - Dependent variable: algo backtest performance (Sharpe, CAGR et al) measured by Report.
//
// - Control variables: price data samples, backtest simulator settings etc.
//
// - Method: as defined by the Optimizer implementation (e.g. brute force, genetic et al) and its ObjectiveRanker func.
type Study struct {
	Training        []ParamSet
	TrainingSamples [][]market.Kline
	TrainingResults map[ParamSetID]Report

	Validation        []ParamSet
	ValidationSamples [][]market.Kline
	ValidationResults map[ParamSetID]Report
}

func NewStudy() Study {
	return Study{
		TrainingResults:   make(map[ParamSetID]Report),
		ValidationResults: make(map[ParamSetID]Report),
	}
}

type ParamSet struct {
	ID     ParamSetID
	Params ParamMap
}

type ParamSetID string

type ParamMap map[string]any

func NewParamSet() ParamSet {
	return ParamSet{
		ID:     ParamSetID(util.NewID()),
		Params: make(map[string]any),
	}
}

type Report struct {
	ID string `csv:"id"`

	Subject ParamSet `csv:",inline"`

	PRR    float64 `csv:"prr"`
	MDD    float64 `csv:"mdd"`
	CAGR   float64 `csv:"cagr"`
	Sharpe float64 `csv:"sharpe"`
	Calmar float64 `csv:"calmar"`

	SampleCount int `csv:"sample_count"`
	TradeCount  int `csv:"trade_count"`

	Backtests []perf.PerformanceReport `csv:"-"`
}

func Summarize(report Report) Report {

	for i := range report.Backtests {
		backtest := report.Backtests[i]

		report.SampleCount++
		report.TradeCount += int(backtest.Trade.TradeCount)

		report.PRR += backtest.Trade.PRR
		report.MDD += backtest.Portfolio.MaxDrawdown
		report.CAGR += backtest.Portfolio.CAGR
		report.Sharpe += backtest.Portfolio.Sharpe
		report.Calmar += backtest.Portfolio.Calmar
	}

	count := float64(report.SampleCount)
	report.PRR /= count
	report.MDD /= count
	report.CAGR /= count
	report.Sharpe /= count
	report.Calmar /= count

	return report
}

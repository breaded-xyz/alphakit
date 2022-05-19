package optimize

import (
	"github.com/thecolngroup/zerotoalgo/internal/util"
	"github.com/thecolngroup/zerotoalgo/market"
	"github.com/thecolngroup/zerotoalgo/perf"
)

// Study is an optimization experiment, prepared and executed by an Optimizer.
// First, a training (in-sample) phase is conducted, followed by a validation (out-of-sample) phase.
// The validation phase reports the out-of-sample (OOS) performance of the optimum param set.
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
	ID string

	Training        []ParamSet
	TrainingSamples [][]market.Kline
	TrainingResults map[ParamSetID]Report

	Validation        []ParamSet
	ValidationSamples [][]market.Kline
	ValidationResults map[ParamSetID]Report
}

func NewStudy() Study {
	return Study{
		ID:                string(util.NewID()),
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

//func (pm ParamMap) MarshalCSV() ([]byte, error) {
//	return []byte(fmt.Sprint(pm)), nil
//}

func NewParamSet() ParamSet {
	return ParamSet{
		ID:     ParamSetID(util.NewID()),
		Params: make(map[string]any),
	}
}

type Report struct {
	ID      string   `csv:"id"`
	Phase   Phase    `csv:"phase"`
	Subject ParamSet `csv:",inline"`

	PRR    float64 `csv:"prr"`
	MDD    float64 `csv:"mdd"`
	CAGR   float64 `csv:"cagr"`
	Sharpe float64 `csv:"sharpe"`
	Calmar float64 `csv:"calmar"`
	WinPct float64 `csv:"win_pct"`

	Kelly    float64 `csv:"kelly"`
	OptimalF float64 `csv:"optimalf"`

	SampleCount int `csv:"sample_count"`
	TradeCount  int `csv:"trade_count"`

	Backtests []perf.PerformanceReport `csv:"-"`
}

func NewReport() Report {
	return Report{
		ID: string(util.NewID()),
	}
}

func Summarize(report Report) Report {

	for i := range report.Backtests {
		backtest := report.Backtests[i]

		if backtest.Trade == nil || backtest.Portfolio == nil {
			continue
		}

		report.SampleCount++
		report.TradeCount += int(backtest.Trade.TradeCount)

		report.PRR += backtest.Trade.PRR
		report.MDD += backtest.Portfolio.MaxDrawdown
		report.CAGR += backtest.Portfolio.CAGR
		report.Sharpe += backtest.Portfolio.Sharpe
		report.Calmar += backtest.Portfolio.Calmar
		report.WinPct += backtest.Trade.PercentProfitable

		report.Kelly += backtest.Trade.Kelly
		report.OptimalF += backtest.Trade.OptimalF
	}

	count := float64(report.SampleCount)
	report.PRR /= count
	report.MDD /= count
	report.CAGR /= count
	report.Sharpe /= count
	report.Calmar /= count
	report.WinPct /= count
	report.Kelly /= count
	report.OptimalF /= count

	return report
}

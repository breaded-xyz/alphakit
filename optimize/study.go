// Copyright 2022 The Coln Group Ltd
// SPDX-License-Identifier: MIT

package optimize

import (
	"github.com/thecolngroup/alphakit/market"
	"github.com/thecolngroup/alphakit/perf"
	"github.com/thecolngroup/gou/id"
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
	TrainingSamples map[AssetID][]market.Kline
	TrainingResults map[ParamSetID]PhaseReport

	Validation        []ParamSet
	ValidationSamples map[AssetID][]market.Kline
	ValidationResults map[ParamSetID]PhaseReport
}

// NewStudy returns a new study.
// Use an Optimizer implementation to prepare and execute the study.
func NewStudy() *Study {
	return &Study{
		ID:                string(id.New()),
		TrainingSamples:   make(map[AssetID][]market.Kline),
		TrainingResults:   make(map[ParamSetID]PhaseReport),
		ValidationSamples: make(map[AssetID][]market.Kline),
		ValidationResults: make(map[ParamSetID]PhaseReport),
	}
}

// AssetID is a string identifer for the asset associated with a sample.
// Typically the symbol of the asset, e.g. btcusdt.
type AssetID string

// PhaseReport is the aggregated performance of a ParamSet across one or more price samples (trials)
// The summary method is owned by the Optimizer implementation, but will typically be the mean (avg) of the individual trials.
type PhaseReport struct {
	ID      string   `csv:"id"`
	Phase   Phase    `csv:"phase"`
	Subject ParamSet `csv:"paramset_,inline"`

	PRR        float64 `csv:"prr"`
	MDD        float64 `csv:"mdd"`
	CAGR       float64 `csv:"cagr"`
	HistVolAnn float64 `csv:"vol"`
	Sharpe     float64 `csv:"sharpe"`
	Calmar     float64 `csv:"calmar"`
	WinPct     float64 `csv:"win_pct"`

	Kelly    float64 `csv:"kelly"`
	OptimalF float64 `csv:"optimalf"`

	SampleCount    int `csv:"sample_count"`
	RoundTurnCount int `csv:"roundturn_count"`

	Trials []perf.PerformanceReport `csv:"-"`
}

// NewReport returns a new empty report with an initialized ID.
func NewReport() PhaseReport {
	return PhaseReport{
		ID: string(id.New()),
	}
}

// Summarize inspects the individual backtests in the report and updates the summary fields.
// Summary takes the average of the metric across all the trials in a phase.
func Summarize(report PhaseReport) PhaseReport {

	if len(report.Trials) == 0 {
		return report
	}

	for i := range report.Trials {
		backtest := report.Trials[i]

		if backtest.PortfolioReport == nil || backtest.TradeReport == nil {
			continue
		}

		if backtest.TradeReport.TradeCount == 0 {
			continue
		}

		report.SampleCount++
		report.RoundTurnCount += int(backtest.TradeReport.TradeCount)

		report.PRR += backtest.TradeReport.PRR
		report.MDD += backtest.PortfolioReport.MaxDrawdown
		report.HistVolAnn += backtest.PortfolioReport.HistVolAnn
		report.CAGR += backtest.PortfolioReport.CAGR
		report.Sharpe += backtest.PortfolioReport.Sharpe
		report.Calmar += backtest.PortfolioReport.Calmar
		report.WinPct += backtest.TradeReport.PercentProfitable

		report.Kelly += backtest.TradeReport.Kelly
		report.OptimalF += backtest.TradeReport.OptimalF
	}

	if report.RoundTurnCount == 0 {
		return report
	}

	count := float64(report.SampleCount)
	report.PRR /= count
	report.MDD /= count
	report.HistVolAnn /= count
	report.CAGR /= count
	report.Sharpe /= count
	report.Calmar /= count
	report.WinPct /= count
	report.Kelly /= count
	report.OptimalF /= count

	return report
}

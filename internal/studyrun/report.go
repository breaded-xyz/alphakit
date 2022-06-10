package studyrun

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/thecolngroup/alphakit/optimize"
	"github.com/thecolngroup/alphakit/perf"
	"github.com/thecolngroup/util"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

var _summaryReportHeader = []string{
	"PRR",
	"MDD",
	"WinPct",
	"CAGR",
	"Sharpe",
	"Calmar",
	"Kelly",
	"OptimalF",
	"Samples",
	"Trades",
}

// phaseReport is a wrapper on optimize.PhaseReport that adds a PK for saving to CSV.
type phaseReport struct {
	StudyID string               `csv:"study_id"`
	Report  optimize.PhaseReport `csv:"phasereport_,inline"`
}

// trialReport is a wrapper on perf.PerformanceReport that adds a compound key for saving to CSV.
type trialReport struct {
	StudyID       string                 `csv:"study_id"`
	PhaseReportID string                 `csv:"phasereport_id"`
	Backtest      perf.PerformanceReport `csv:"backtest_,inline"`
}

// prepareStudyForCSV returns data that is ready for saving to CSV.
func prepareStudyForCSV(study optimize.Study) ([]phaseReport, []trialReport) {

	var summaries []phaseReport
	var backtests []trialReport

	flattenResults := func(results map[optimize.ParamSetID]optimize.PhaseReport) {
		for k := range results {
			report := results[k]
			summaries = append(summaries, phaseReport{
				StudyID: study.ID,
				Report:  report,
			})
			for i := range report.Trials {
				backtests = append(backtests, trialReport{
					StudyID:       study.ID,
					PhaseReportID: report.ID,
					Backtest:      report.Trials[i],
				})
			}
		}
	}

	flattenResults(study.TrainingResults)
	flattenResults(study.ValidationResults)

	return summaries, backtests
}

// printSummary prints a summary report to stdout.
func printSummary(report optimize.PhaseReport) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(_summaryReportHeader)
	table.Append([]string{
		fmt.Sprintf("%.2f", report.PRR),
		fmt.Sprintf("%.2f", report.MDD*100),
		fmt.Sprintf("%.2f", report.WinPct*100),
		fmt.Sprintf("%.2f", report.CAGR*100),
		fmt.Sprintf("%.2f", report.Sharpe),
		fmt.Sprintf("%.2f", report.Calmar),
		fmt.Sprintf("%.2f", report.Kelly),
		fmt.Sprintf("%.2f", report.OptimalF),
		fmt.Sprintf("%d", report.SampleCount),
		fmt.Sprintf("%d", report.TradeCount),
	})
	table.Render()

}

// printParams pretty prints a map.
func printParams(params map[string]any) {
	keys := maps.Keys(params)
	slices.Sort(keys)
	for _, k := range keys {
		fmt.Printf("%s: %s\n", k, util.ToString(params[k]))
	}
}

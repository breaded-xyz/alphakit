package studyrun

import (
	"fmt"
	"os"

	"github.com/colngroup/zero2algo/internal/util"
	"github.com/colngroup/zero2algo/optimize"
	"github.com/colngroup/zero2algo/perf"
	"github.com/olekukonko/tablewriter"
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

type SummaryReport struct {
	StudyID string          `csv:"study_id"`
	Summary optimize.Report `csv:",inline"`
}

type BacktestReport struct {
	StudyID   string                 `csv:"study_id"`
	SummaryID string                 `csv:"summary_id"`
	Backtest  perf.PerformanceReport `csv:",inline"`
}

func PrepareStudyForCSV(study optimize.Study) ([]SummaryReport, []BacktestReport) {

	var summaries []SummaryReport
	var backtests []BacktestReport

	flattenResults := func(results map[optimize.ParamSetID]optimize.Report) {
		for k := range results {
			report := results[k]
			summaries = append(summaries, SummaryReport{
				StudyID: study.ID,
				Summary: report,
			})
			for i := range report.Backtests {
				backtests = append(backtests, BacktestReport{
					StudyID:   study.ID,
					SummaryID: report.ID,
					Backtest:  report.Backtests[i],
				})
			}
		}
	}

	flattenResults(study.TrainingResults)
	flattenResults(study.ValidationResults)

	return summaries, backtests
}

func PrintSummaryReport(report optimize.Report) {
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

func PrintParams(params map[string]any) {
	keys := maps.Keys(params)
	slices.Sort(keys)
	for _, k := range keys {
		fmt.Printf("%s: %s\n", k, util.ToString(params[k]))
	}
}

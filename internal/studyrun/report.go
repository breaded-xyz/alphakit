package studyrun

import (
	"fmt"
	"os"

	"github.com/colngroup/zero2algo/internal/util"
	"github.com/colngroup/zero2algo/optimize"
	"github.com/colngroup/zero2algo/perf"
	"github.com/olekukonko/tablewriter"
)

var _summaryReportHeader = []string{
	"PRR",
	"MDD",
	"WinPct",
	"CAGR",
	"Sharpe",
	"Calmar",
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
	return nil, nil
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
		fmt.Sprintf("%d", report.SampleCount),
		fmt.Sprintf("%d", report.TradeCount),
	})
	table.Render()

}

func PrintParams(params map[string]any) {
	for k, v := range params {
		fmt.Printf("%s: %s\n", k, util.ToString(v))
	}
}

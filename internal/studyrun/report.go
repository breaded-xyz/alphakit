package studyrun

import (
	"github.com/colngroup/zero2algo/optimize"
	"github.com/colngroup/zero2algo/perf"
)

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

func PrintSummary(report optimize.Report) {

}

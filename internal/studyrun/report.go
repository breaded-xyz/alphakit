package studyrun

import (
	"time"

	"github.com/thecolngroup/alphakit/broker"
	"github.com/thecolngroup/alphakit/optimize"
	"github.com/thecolngroup/alphakit/perf"
)

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

type tradeDetailRow struct {
	StudyID       string       `csv:"study_id"`
	PhaseReportID string       `csv:"phasereport_id"`
	BacktestID    string       `csv:"backtest_id"`
	Trade         broker.Trade `csv:"trade_,inline"`
}

type curveDetailRow struct {
	StudyID       string    `csv:"study_id"`
	PhaseReportID string    `csv:"phasereport_id"`
	BacktestID    string    `csv:"backtest_id"`
	Time          time.Time `csv:"time"`
	Amount        float64   `csv:"amount"`
}

// prepareStudyForCSV returns data that is ready for saving to CSV.
func prepareStudyForCSV(study optimize.Study) ([]phaseReport, []trialReport, []tradeDetailRow, []curveDetailRow) {

	var phaseReports []phaseReport
	var trialReports []trialReport

	flattenResults := func(results map[optimize.ParamSetID]optimize.PhaseReport) {
		for k := range results {
			report := results[k]
			phaseReports = append(phaseReports, phaseReport{
				StudyID: study.ID,
				Report:  report,
			})
			for i := range report.Trials {
				trialReports = append(trialReports, trialReport{
					StudyID:       study.ID,
					PhaseReportID: report.ID,
					Backtest:      report.Trials[i],
				})
			}
		}
	}

	flattenResults(study.TrainingResults)
	flattenResults(study.ValidationResults)

	if len(study.Validation) == 0 || len(study.ValidationResults) == 0 {
		return phaseReports, trialReports, nil, nil
	}

	optimaReport := study.ValidationResults[study.Validation[0].ID]

	var tradeRows []tradeDetailRow
	var curveRows []curveDetailRow
	for _, trial := range optimaReport.Trials {
		for _, trade := range trial.TradeReport.Trades {
			tradeRows = append(tradeRows, tradeDetailRow{
				StudyID:       study.ID,
				PhaseReportID: optimaReport.ID,
				BacktestID:    trial.ID,
				Trade:         trade,
			})
		}

		curve := trial.PortfolioReport.EquityCurve
		sortedKeys := curve.SortKeys()
		for _, k := range sortedKeys {
			curveRows = append(curveRows, curveDetailRow{
				StudyID:       study.ID,
				PhaseReportID: optimaReport.ID,
				BacktestID:    trial.ID,
				Time:          k.Time(),
				Amount:        curve[k].InexactFloat64(),
			})
		}

	}

	return phaseReports, trialReports, tradeRows, curveRows
}

package studyrun

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/thecolngroup/alphakit/optimize"
	"github.com/thecolngroup/util"
)

const _filenameFriendlyTimeFormat = "20060102T150405"

// WriteStudy writes a study to CSV.
func WriteStudy(path string, study optimize.Study) error {

	summaries, backtests := PrepareStudyForCSV(study)

	if err := writeSummaryReports(path, summaries); err != nil {
		return err
	}

	if err := writeBacktestReports(path, backtests); err != nil {
		return err
	}

	return nil
}

func writeSummaryReports(path string, reports []SummaryReport) error {
	prefix := time.Now().UTC().Format(_filenameFriendlyTimeFormat)
	out := filepath.Join(path, fmt.Sprintf("%s-summaryreports.csv", prefix))
	if err := saveDataToCSV(out, reports); err != nil {
		return err
	}

	return nil
}

func writeBacktestReports(path string, reports []BacktestReport) error {
	prefix := time.Now().UTC().Format(_filenameFriendlyTimeFormat)
	out := filepath.Join(path, fmt.Sprintf("%s-backtestreports.csv", prefix))
	if err := saveDataToCSV(out, reports); err != nil {
		return err
	}

	return nil
}

func saveDataToCSV(filename string, data interface{}) error {

	encMap := func(m map[string]any) ([]byte, error) {
		return []byte(fmt.Sprint(m)), nil
	}

	encParamMap := func(m optimize.ParamMap) ([]byte, error) {
		return []byte(fmt.Sprint(m)), nil
	}

	return util.WriteToCSV(filename, data, encMap, encParamMap)
}

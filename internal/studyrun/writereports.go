package studyrun

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/thecolngroup/alphakit/optimize"
	"github.com/thecolngroup/util"
)

const _filenameFriendlyTimeFormat = "20060102T150405"

// writeStudy writes a study to CSV.
func writeStudy(path string, study optimize.Study) error {

	summaries, backtests := prepareStudyForCSV(study)

	if err := writePhaseReports(path, summaries); err != nil {
		return err
	}

	if err := writeTrialReports(path, backtests); err != nil {
		return err
	}

	return nil
}

func writePhaseReports(path string, reports []phaseReport) error {
	prefix := time.Now().UTC().Format(_filenameFriendlyTimeFormat)
	out := filepath.Join(path, fmt.Sprintf("%s-phasereports.csv", prefix))
	if err := saveDataToCSV(out, reports); err != nil {
		return err
	}

	return nil
}

func writeTrialReports(path string, reports []trialReport) error {
	prefix := time.Now().UTC().Format(_filenameFriendlyTimeFormat)
	out := filepath.Join(path, fmt.Sprintf("%s-trialreports.csv", prefix))
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

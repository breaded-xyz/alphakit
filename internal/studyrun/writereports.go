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

	phaseReports, trialReports, optimaTrades, optimaCurves := prepareStudyForCSV(study)

	prefix := time.Now().UTC().Format(_filenameFriendlyTimeFormat)

	if err := writePhaseReports(path, prefix, phaseReports); err != nil {
		return err
	}

	if err := writeTrialReports(path, prefix, trialReports); err != nil {
		return err
	}

	if err := writeTrades(path, prefix, optimaTrades); err != nil {
		return err
	}

	if err := writeCurves(path, prefix, optimaCurves); err != nil {
		return err
	}

	return nil
}

func writePhaseReports(path string, prefix string, reports []phaseReport) error {
	out := filepath.Join(path, fmt.Sprintf("%s-phasereports.csv", prefix))
	if err := saveDataToCSV(out, reports); err != nil {
		return err
	}
	return nil
}

func writeTrialReports(path string, prefix string, reports []trialReport) error {
	out := filepath.Join(path, fmt.Sprintf("%s-trialreports.csv", prefix))
	if err := saveDataToCSV(out, reports); err != nil {
		return err
	}
	return nil
}

func writeTrades(path string, prefix string, trades []tradeDetailRow) error {
	out := filepath.Join(path, fmt.Sprintf("%s-optima-trades.csv", prefix))
	if err := saveDataToCSV(out, trades); err != nil {
		return err
	}
	return nil
}

func writeCurves(path string, prefix string, curves []curveDetailRow) error {
	out := filepath.Join(path, fmt.Sprintf("%s-optima-equitycurves.csv", prefix))
	if err := saveDataToCSV(out, curves); err != nil {
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

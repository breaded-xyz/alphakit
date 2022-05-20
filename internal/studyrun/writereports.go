package studyrun

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jszwec/csvutil"
	"github.com/thecolngroup/alphakit/optimize"
)

const _filenameFriendlyTimeFormat = "20060102T150405"

func WriteStudy(path string, study optimize.Study) error {

	summaries, backtests := PrepareStudyForCSV(study)

	if err := WriteSummaryReports(path, summaries); err != nil {
		return err
	}

	if err := WriteBacktestReports(path, backtests); err != nil {
		return err
	}

	return nil
}

func WriteSummaryReports(path string, reports []SummaryReport) error {
	prefix := time.Now().UTC().Format(_filenameFriendlyTimeFormat)
	out := filepath.Join(path, fmt.Sprintf("%s-summaryreports.csv", prefix))
	if err := SaveStructToCSV(out, reports); err != nil {
		return err
	}

	return nil
}

func WriteBacktestReports(path string, reports []BacktestReport) error {
	prefix := time.Now().UTC().Format(_filenameFriendlyTimeFormat)
	out := filepath.Join(path, fmt.Sprintf("%s-backtestreports.csv", prefix))
	if err := SaveStructToCSV(out, reports); err != nil {
		return err
	}

	return nil
}

func SaveStructToCSV(filename string, data interface{}) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Wrap file in CSV struct encoder
	w := csv.NewWriter(f)
	enc := csvutil.NewEncoder(w)
	enc.Tag = "csv"
	enc.Register(func(m map[string]any) ([]byte, error) {
		return []byte(fmt.Sprint(m)), nil
	})
	enc.Register(func(m optimize.ParamMap) ([]byte, error) {
		return []byte(fmt.Sprint(m)), nil
	})

	// Only write header if new file
	info, err := f.Stat()
	if err != nil {
		return err
	}
	if info.Size() > 0 {
		enc.AutoHeader = false
	}

	// Write report
	err = enc.Encode(data)
	if err != nil {
		return err
	}
	w.Flush()

	return nil
}

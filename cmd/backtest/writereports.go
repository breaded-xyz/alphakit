package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/colngroup/zero2algo/perf"
	"github.com/jszwec/csvutil"
)

const _filenameFriendlyTimeFormat = "20060102T150405"

func writeReports(path string, reports []perf.PerformanceReport) error {
	prefix := time.Now().UTC().Format(_filenameFriendlyTimeFormat)
	out := filepath.Join(path, fmt.Sprintf("%s-performancereport.csv", prefix))
	if err := saveStructToCSV(out, reports); err != nil {
		return err
	}

	return nil
}

func saveStructToCSV(filename string, data interface{}) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Wrap file in CSV struct encoder
	w := csv.NewWriter(f)
	enc := csvutil.NewEncoder(w)
	enc.Tag = "json"

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

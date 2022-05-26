package market

import (
	"encoding/csv"
	"io/fs"
	"os"
	"path/filepath"
)

// ReadKlinesFromCSV reads all the .csv files in a given directory or a single file into a slice of Klines.
// Wraps a CSVKlineReader for convenience. For finer grained memory management use the base kline reader.
func ReadKlinesFromCSV(path string) ([]Kline, error) {
	var prices []Kline

	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".csv" {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		//nolint:errcheck // Read ops only so safe to ignore err return
		defer file.Close()
		klines, err := NewCSVKlineReader(csv.NewReader(file)).ReadAll()
		if err != nil {
			return err
		}
		prices = append(prices, klines...)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return prices, nil
}
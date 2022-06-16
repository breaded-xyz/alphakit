// Copyright 2022 The Coln Group Ltd
// SPDX-License-Identifier: MIT

package market

import (
	"encoding/csv"
	"io/fs"
	"os"
	"path/filepath"
)

// ReadKlinesFromCSV reads all the .csv files in a given directory or a single file into a slice of Klines.
// Wraps a default CSVKlineReader with Binance decoder for convenience.
// For finer grained memory management use the base kline reader.
func ReadKlinesFromCSV(path string) ([]Kline, error) {
	return ReadKlinesFromCSVWithDecoder(path, MakeCSVKlineReader(NewBinanceCSVKlineReader))
}

// ReadKlinesFromCSVWithDecoder permits using a custom CSVKlineReader.
func ReadKlinesFromCSVWithDecoder(path string, maker MakeCSVKlineReader) ([]Kline, error) {
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
		reader := maker(csv.NewReader(file))
		klines, err := reader.ReadAll()
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

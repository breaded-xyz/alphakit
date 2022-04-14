package main

import (
	"encoding/csv"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/colngroup/zero2algo/market"
)

func readPrices(path string) ([]market.Kline, error) {
	var prices []market.Kline

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
		defer file.Close()
		klines, err := market.NewCSVKlineReader(csv.NewReader(file)).ReadAll()
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

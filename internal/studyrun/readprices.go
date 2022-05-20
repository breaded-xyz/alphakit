package studyrun

import (
	"encoding/csv"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/thecolngroup/alphakit/market"
)

func ReadPricesFromConfig(config map[string]any) ([][]market.Kline, error) {

	root, ok := config["samples"]
	if !ok {
		return nil, errors.New("'samples' key not found")
	}
	paths := root.(map[string]any)

	var prices [][]market.Kline
	for _, v := range paths {
		path := v.(string)
		series, err := ReadPriceSeries(path)
		if err != nil {
			return nil, err
		}
		prices = append(prices, series)
	}
	return prices, nil
}

func ReadPriceSeries(path string) ([]market.Kline, error) {
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

package studyrun

import (
	"errors"

	"github.com/thecolngroup/alphakit/market"
)

// ReadPricesFromConfig reads the price samples from a config file params.
func ReadPricesFromConfig(config map[string]any) ([][]market.Kline, error) {

	root, ok := config["samples"]
	if !ok {
		return nil, errors.New("'samples' key not found")
	}
	paths := root.(map[string]any)

	var prices [][]market.Kline
	for _, v := range paths {
		path := v.(string)
		series, err := market.CSVKlineReadAllDir(path)
		if err != nil {
			return nil, err
		}
		prices = append(prices, series)
	}
	return prices, nil
}

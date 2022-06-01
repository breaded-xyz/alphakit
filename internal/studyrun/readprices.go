package studyrun

import (
	"errors"
	"fmt"

	"github.com/thecolngroup/alphakit/market"
	"github.com/thecolngroup/util"
)

// readPricesFromConfig reads the price samples from a config file params.
func readPricesFromConfig(config map[string]any, typeRegistry map[string]any) ([][]market.Kline, error) {

	if _, ok := config["samples"]; !ok {
		return nil, errors.New("'samples' key not found")
	}
	root := config["samples"].(map[string]any)

	// Load decoder from type registry
	if _, ok := root["decoder"]; !ok {
		return nil, errors.New("'decoder' key not found")
	}
	decoder := util.ToString(root["decoder"])
	if _, ok := typeRegistry[decoder]; !ok {
		return nil, fmt.Errorf("'%s' key not found in type registry", decoder)
	}
	maker := typeRegistry[decoder].(market.MakeCSVKlineReader)

	var prices [][]market.Kline
	for _, v := range root["paths"].([]any) {
		path := v.(string)
		series, err := market.ReadKlinesFromCSVWithDecoder(path, maker)
		if err != nil {
			return nil, err
		}
		prices = append(prices, series)
	}
	return prices, nil
}

package studyrun

import (
	"errors"
	"fmt"

	"github.com/thecolngroup/alphakit/market"
	"github.com/thecolngroup/alphakit/optimize"
	"github.com/thecolngroup/gou/conv"
)

// readPricesFromConfig reads the price samples from a config file params.
func readPricesFromConfig(config map[string]any, typeRegistry map[string]any) (map[optimize.AssetID][]market.Kline, error) {

	if _, ok := config["samples"]; !ok {
		return nil, errors.New("'samples' key not found")
	}
	root := config["samples"].([]any)
	samples := make(map[optimize.AssetID][]market.Kline)

	for _, sub := range root {

		cfg := sub.(map[string]any)

		// Load decoder from type registry
		if _, ok := cfg["decoder"]; !ok {
			return nil, errors.New("'decoder' key not found")
		}
		decoder := conv.ToString(cfg["decoder"])
		if _, ok := typeRegistry[decoder]; !ok {
			return nil, fmt.Errorf("'%s' key not found in type registry", decoder)
		}
		maker := typeRegistry[decoder].(market.MakeCSVKlineReader)

		// Load path to price files from config
		path := cfg["path"].(string)
		series, err := market.ReadKlinesFromCSVWithDecoder(path, maker)
		if err != nil {
			return nil, err
		}

		// Load asset key from config
		assetID := optimize.AssetID(cfg["asset"].(string))
		samples[assetID] = series
	}

	return samples, nil
}

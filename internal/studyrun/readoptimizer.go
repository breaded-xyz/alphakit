package studyrun

import (
	"errors"
	"fmt"

	"github.com/thecolngroup/alphakit/optimize"
	"github.com/thecolngroup/alphakit/trader"
	"github.com/thecolngroup/util"
)

// readBruteOptimizerFromConfig creates a new brute optimizer from a config file params.
func readBruteOptimizerFromConfig(config map[string]any, botRegistry map[string]trader.MakeFromConfig) (*optimize.BruteOptimizer, error) {

	var optimizer optimize.BruteOptimizer

	// Load root config
	if _, ok := config["optimizer"]; !ok {
		return nil, errors.New("'optimizer' key not found")
	}
	root := config["optimizer"].(map[string]any)

	// Load bot from type registry
	if _, ok := root["bot"]; !ok {
		return nil, errors.New("'bot' key not found")
	}
	bot := util.ToString(root["bot"])
	if _, ok := botRegistry[bot]; !ok {
		return nil, fmt.Errorf("'%s' key not found in type registry", bot)
	}

	optimizer = optimize.NewBruteOptimizer()
	optimizer.MakeBot = botRegistry[bot]
	optimizer.SampleSplitPct = util.ToFloat(root["samplesplitpct"])
	optimizer.WarmupBarCount = util.ToInt(root["warmupbarcount"])
	optimizer.Ranker = optimize.SharpeRanker

	return &optimizer, nil
}

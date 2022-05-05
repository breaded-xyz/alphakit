package studyrun

import (
	"errors"
	"fmt"

	"github.com/colngroup/zero2algo/internal/util"
	"github.com/colngroup/zero2algo/optimize"
	"github.com/colngroup/zero2algo/trader"
)

func ReadBruteOptimizerFromConfig(config map[string]any) (*optimize.BruteOptimizer, error) {

	var optimizer optimize.BruteOptimizer

	// Load root config
	if _, ok := config["optimizer"]; !ok {
		return nil, errors.New("'optimizer' key not found")
	}
	root := config["optimizer"].(map[string]any)

	// Load bot from type registry
	if _, ok := config["bot"]; !ok {
		return nil, errors.New("'bot' key not found")
	}
	bot := util.ToString(config["bot"])
	if _, ok := _typeRegistry[bot]; !ok {
		return nil, fmt.Errorf("'%s' key not found in type registry", bot)
	}

	optimizer = optimize.NewBruteOptimizer()
	optimizer.MakeBot = _typeRegistry[bot].(trader.MakeFromConfig)
	optimizer.SampleSplitPct = util.ToFloat(root["samplesplitpct"])
	optimizer.WarmupBarCount = util.ToInt(root["warmupbarcount"])
	optimizer.Ranker = optimize.SharpeRanker

	return &optimizer, nil
}

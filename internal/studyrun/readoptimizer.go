package studyrun

import (
	"errors"
	"fmt"

	"github.com/colngroup/zero2algo/broker/backtest"
	"github.com/colngroup/zero2algo/internal/util"
	"github.com/colngroup/zero2algo/optimize"
	"github.com/colngroup/zero2algo/trader"
)

func ReadBruteOptimizerFromConfig(config map[string]any) (optimize.Optimizer, error) {

	var optimizer optimize.BruteOptimizer

	// Load root config
	if _, ok := config["optimizer"]; !ok {
		return nil, errors.New("'optimizer' key not found")
	}
	root := config["optimizer"].(map[string]any)

	// Load bot from type registry
	bot, ok := root["bot"].(string)
	if !ok {
		return nil, errors.New("'bot' key not found")
	}
	if _, ok := _typeRegistry[bot]; !ok {
		return nil, fmt.Errorf("'%s' key not found in type registry", bot)
	}
	optimizer.MakeBot = _typeRegistry[bot].(trader.MakeBot)
	
	optimizer.MakeDealer = backtest.MakeDealer
	optimizer.SampleSplitPct = util.ToFloat(root["sampleSplitPct"])
	optimizer.WarmupBarCount = util.ToInt(root["warmupBarCount"])
	optimizer.Ranker = optimize.SharpeRanker

	return &optimizer, nil
}

// Copyright 2022 The Coln Group Ltd
// SPDX-License-Identifier: MIT

package studyrun

import (
	"errors"
	"fmt"

	"github.com/thecolngroup/alphakit/optimize"
	"github.com/thecolngroup/alphakit/trader"
	"github.com/thecolngroup/gou/conv"
)

// readBruteOptimizerFromConfig creates a new brute optimizer from a config file params.
func readBruteOptimizerFromConfig(config map[string]any, typeRegistry map[string]any) (*optimize.BruteOptimizer, error) {

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
	bot := conv.ToString(root["bot"])
	if _, ok := typeRegistry[bot]; !ok {
		return nil, fmt.Errorf("'%s' key not found in type registry", bot)
	}

	optimizer = optimize.NewBruteOptimizer()
	optimizer.MakeBot = typeRegistry[bot].(trader.MakeFromConfig)
	optimizer.SampleSplitPct = conv.ToFloat(root["samplesplitpct"])
	optimizer.WarmupBarCount = conv.ToInt(root["warmupbarcount"])
	optimizer.Ranker = optimize.SharpeRanker

	return &optimizer, nil
}

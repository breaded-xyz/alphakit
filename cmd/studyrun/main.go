// Copyright 2022 The Coln Group Ltd
// SPDX-License-Identifier: MIT

package main

import (
	"log"
	"os"

	"github.com/thecolngroup/alphakit/cmd/studyrun/app"
	"github.com/thecolngroup/alphakit/market"
	"github.com/thecolngroup/alphakit/trader"
	"github.com/thecolngroup/alphakit/trader/hodl"
	"github.com/thecolngroup/alphakit/trader/trend"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	return app.Run(
		args,
		map[string]any{
			"hodl":        trader.MakeFromConfig(hodl.MakeBotFromConfig),
			"trend.cross": trader.MakeFromConfig(trend.MakeCrossBotFromConfig),
			"trend.apex":  trader.MakeFromConfig(trend.MakeApexBotFromConfig),
			"binance":     market.MakeCSVKlineReader(market.NewBinanceCSVKlineReader),
			"metatrader":  market.MakeCSVKlineReader(market.NewMetaTraderCSVKlineReader),
		},
		app.BuildVersion{
			GitTag:    buildGitTag,
			GitCommit: buildGitCommit,
			Time:      buildTime,
			User:      buildUser,
		},
	)
}

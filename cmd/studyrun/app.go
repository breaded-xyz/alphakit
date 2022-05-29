package main

import (
	"github.com/thecolngroup/alphakit/internal/studyrun"
	"github.com/thecolngroup/alphakit/trader"
	"github.com/thecolngroup/alphakit/trader/hodl"
	"github.com/thecolngroup/alphakit/trader/trend"
)

type BuildVersion struct {
	GitTag    string
	GitCommit string
	Time      string
	User      string
}

func RunApp(args []string, botRegistry map[string]trader.MakeFromConfig, build BuildVersion) error {
	app := studyrun.App{
		Args:      args,
		GitCommit: build.GitCommit,
		GitTag:    build.GitTag,
		BuildTime: build.Time,
		BuildUser: build.User,
		BotRegistry: map[string]trader.MakeFromConfig{
			"hodl":        trader.MakeFromConfig(hodl.MakeBotFromConfig),
			"trend.cross": trader.MakeFromConfig(trend.MakeCrossBotFromConfig),
			"trend.apex":  trader.MakeFromConfig(trend.MakeApexBotFromConfig),
		},
	}
	return app.Run()
}

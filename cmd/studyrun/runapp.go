package main

import (
	"github.com/thecolngroup/alphakit/internal/studyrun"
	"github.com/thecolngroup/alphakit/trader"
	"github.com/thecolngroup/alphakit/trader/hodl"
	"github.com/thecolngroup/alphakit/trader/trend"
)

// BuildVersion describes key info to identify the build of the app.
type BuildVersion struct {
	GitTag    string
	GitCommit string
	Time      string
	User      string
}

// RunApp is the entrypoint for executing the studyrun cmd outside of the alphakit project.
// Param args are the cmd line args and excludes the app name.
// Param botRegistry enables injection of bots to be loaded by string name from config.
// Param build is optional and will render key build version info to the console during execution.
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

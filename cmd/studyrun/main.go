package main

import (
	"log"
	"os"

	"github.com/thecolngroup/alphakit/cmd/studyrun/app"
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
		map[string]trader.MakeFromConfig{
			"hodl":        trader.MakeFromConfig(hodl.MakeBotFromConfig),
			"trend.cross": trader.MakeFromConfig(trend.MakeCrossBotFromConfig),
			"trend.apex":  trader.MakeFromConfig(trend.MakeApexBotFromConfig),
		},
		app.BuildVersion{
			GitTag:    buildGitTag,
			GitCommit: buildGitCommit,
			Time:      buildTime,
			User:      buildUser,
		},
	)
}

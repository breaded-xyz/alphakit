package app

import (
	"github.com/thecolngroup/alphakit/internal/studyrun"
)

// BuildVersion describes key info to identify the build of the app.
type BuildVersion struct {
	GitTag    string
	GitCommit string
	Time      string
	User      string
}

// Run is the entrypoint for executing the studyrun cmd outside of the alphakit project.
// Param args are the cmd line args (excluding the cmd name).
// Param typeRegistry maps a string name used in a config file to a type such as a trading bot.
// Param build will render strings to the console during execution.
func Run(args []string, typeRegistry map[string]any, build BuildVersion) error {
	app := studyrun.App{
		Args:         args,
		GitCommit:    build.GitCommit,
		GitTag:       build.GitTag,
		BuildTime:    build.Time,
		BuildUser:    build.User,
		TypeRegistry: typeRegistry,
	}
	return app.Run()
}

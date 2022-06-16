// Copyright 2022 The Coln Group Ltd
// SPDX-License-Identifier: MIT

package studyrun

import (
	"context"
	"fmt"
	"strings"

	"github.com/schollz/progressbar/v3"
	"github.com/thecolngroup/alphakit/optimize"
	"golang.org/x/exp/maps"
)

// App exposes an entry point to execute the studyrun application.
type App struct {
	Args         []string
	GitTag       string
	GitCommit    string
	BuildTime    string
	BuildUser    string
	TypeRegistry map[string]any
}

// Run executes the studyrun application. Should be called from main.
func (app *App) Run() error {
	print(_titleArt)

	print("\n ----- Build Info -----\n")
	print("App: studyrun\n")
	fmt.Printf("Tag: %s, Commit: %s\n", app.GitTag, app.GitCommit)
	fmt.Printf("Time: %s, User: %s\n", app.BuildTime, app.BuildUser)

	if len(app.Args) < 2 {
		print("Expect args: [config filename] [output path]\n")
	}
	configFilename, outputPath := app.Args[0], app.Args[1]

	print("\n----- Study Configuration -----\n")
	fmt.Printf("Reading config '%s' ... ", configFilename)
	config, err := readConfig(configFilename)
	if err != nil {
		return err
	}
	print("done\n")

	print("Reading price samples... ")
	samples, err := readPricesFromConfig(config, app.TypeRegistry)
	if err != nil {
		return err
	}
	print("done\n")

	print("Reading param space... ")
	psets, err := readParamSpaceFromConfig(config)
	if err != nil {
		return err
	}
	print("done\n")

	print("Reading dealer... ")
	makeDealer, err := readDealerFromConfig(config)
	if err != nil {
		return err
	}
	print("done\n")

	print("Reading optimizer... ")
	optimizer, err := readBruteOptimizerFromConfig(config, app.TypeRegistry)
	if err != nil {
		return err
	}
	optimizer.MakeDealer = makeDealer
	print("done\n")

	print("\n----- Study Execution -----\n")

	print("Preparing study... ")
	stepCount, err := optimizer.Prepare(psets, samples)
	if err != nil {
		return err
	}
	print("done\n")
	fmt.Printf("Estimated backtests # required: %d\n", stepCount)

	print("Running study... ")
	bar := progressbar.Default(int64(stepCount), "Running backtests... ")
	stepCh, err := optimizer.Start(context.Background())
	if err != nil {
		return err
	}
	var errs []string
	for step := range stepCh {
		_ = bar.Add(1)
		if step.Err != nil {
			errs = append(errs, step.Err.Error())
		}
	}
	_ = bar.Finish()
	print("Study complete\n")

	if len(errs) > 0 {
		print("Errors encountered during optimization:\n")
		print(strings.Join(errs, "\n"))
	}

	fmt.Printf("Writing study results to output directory '%s'... ", outputPath)
	if err := optimize.WriteStudyResultToCSV(outputPath, optimizer.Study()); err != nil {
		return err
	}
	print("done\n")

	print("\n----- Optima Training Result -----\n")
	validationResult := maps.Values(optimizer.Study().ValidationResults)[0]
	trainingResult := optimizer.Study().TrainingResults[validationResult.Subject.ID]
	printSummary(trainingResult)

	print("\n----- Optima Validation Result -----\n")
	printSummary(validationResult)

	print("\n----- Optima Params -----\n")
	printParams(validationResult.Subject.Params)

	return nil
}

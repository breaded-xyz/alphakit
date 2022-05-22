package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/schollz/progressbar/v3"
	"github.com/thecolngroup/alphakit/internal/studyrun"
	"golang.org/x/exp/maps"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	print(_titleArt)

	print("\n ----- Build Info -----\n")
	print("App: studyrun\n")
	fmt.Printf("Tag: %s, Commit: %s\n", buildGitTag, buildGitCommit)
	fmt.Printf("Time: %s, User: %s\n", buildTime, buildUser)

	if len(args) < 2 {
		print("Expect args: [config filename] [output path]\n")
	}
	configFilename, outputPath := args[0], args[1]

	print("\n----- Study Configuration -----\n")
	fmt.Printf("Reading config '%s' ... ", configFilename)
	config, err := studyrun.ReadConfig(configFilename)
	if err != nil {
		return err
	}
	print("done\n")

	print("Reading price samples... ")
	samples, err := studyrun.ReadPricesFromConfig(config)
	if err != nil {
		return err
	}
	print("done\n")

	print("Reading param space... ")
	psets, err := studyrun.ReadParamSpaceFromConfig(config)
	if err != nil {
		return err
	}
	print("done\n")

	print("Reading dealer... ")
	makeDealer, err := studyrun.ReadDealerFromConfig(config)
	if err != nil {
		return err
	}
	print("done\n")

	print("Reading optimizer... ")
	optimizer, err := studyrun.ReadBruteOptimizerFromConfig(config)
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
	if err := studyrun.WriteStudy(outputPath, optimizer.Study()); err != nil {
		return err
	}
	print("done\n")

	print("\n----- Optima Training Result -----\n")
	validationResult := maps.Values(optimizer.Study().ValidationResults)[0]
	trainingResult := optimizer.Study().TrainingResults[validationResult.Subject.ID]
	studyrun.PrintSummaryReport(trainingResult)

	print("\n----- Optima Validation Result -----\n")
	studyrun.PrintSummaryReport(validationResult)

	print("\n----- Optima Params -----\n")
	studyrun.PrintParams(validationResult.Subject.Params)

	return nil
}

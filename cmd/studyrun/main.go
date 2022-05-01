package main

import (
	"context"
	"log"
	"os"

	"github.com/colngroup/zero2algo/internal/studyrun"
	"github.com/schollz/progressbar/v3"
)

const _outputDir = ".out"

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {

	print("Reading study config... ")
	config, err := studyrun.ReadConfig(args[0])
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

	print("Reading param set... ")
	psets, err := studyrun.ReadParamSetFromConfig(config)
	if err != nil {
		return err
	}
	print("done\n")

	print("Reading optimizer... ")
	optimizer, err := studyrun.ReadBruteOptimizerFromConfig(config)
	if err != nil {
		return err
	}
	print("done\n")

	print("Preparing optimizer... ")
	stepCount, err := optimizer.Prepare(psets, samples)
	if err != nil {
		return err
	}
	print("done\n")

	print("Executing optimizer... ")
	bar := progressbar.Default(int64(stepCount), "Running backtests")
	stepCh, err := optimizer.Start(context.Background())
	if err != nil {
		return err
	}
	for range stepCh {
		bar.Add(1)
	}
	bar.Finish()

	return nil
}

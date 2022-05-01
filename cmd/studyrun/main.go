package main

import (
	"context"
	"fmt"
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
	print("Executing studyrun...\n")

	fmt.Printf("Reading config '%s' ... ", args[0])
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

	print("Preparing study... ")
	stepCount, err := optimizer.Prepare(psets, samples)
	if err != nil {
		return err
	}
	print("done\n")

	print("Running study... ")
	bar := progressbar.Default(int64(stepCount), "Running backtests... ")
	stepCh, err := optimizer.Start(context.Background())
	if err != nil {
		return err
	}
	for range stepCh {
		bar.Add(1)
	}
	bar.Finish()
	print("Study complete\n")

	//spew.Dump(maps.Values(optimizer.Study().ValidationResults)[0].Backtests)

	print("studyrun execution complete\n")

	return nil
}

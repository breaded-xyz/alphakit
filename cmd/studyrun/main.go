package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/colngroup/zero2algo/internal/studyrun"
	"github.com/davecgh/go-spew/spew"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/exp/maps"
)

const _outputDir = ".out"

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	print("studyrun\n")
	fmt.Printf("Tag %s, Commit %s\n", buildGitTag, buildGitCommit)
	fmt.Printf("Time: %s, User: %s\n", buildTime, buildUser)

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

	print("Raw read out:\n\n")
	optimaResult := maps.Values(optimizer.Study().ValidationResults)[0]
	spew.Dump(optimaResult.Subject.Params)
	spew.Config.MaxDepth = 1
	spew.Dump(optimaResult)

	return nil
}

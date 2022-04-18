package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/colngroup/zero2algo/optimize"
	"github.com/colngroup/zero2algo/perf"
	"github.com/schollz/progressbar/v3"
)

const _outputDir = ".out"

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	print("Reading prices... ")
	_, err := readPrices(args[0])
	if err != nil {
		return err
	}
	print("done\n")

	print("Reading config... ")
	config, err := readConfig(args[1])
	if err != nil {
		return err
	}
	print("done\n")

	print("Preparing optimizer...\n")
	optimizer := optimize.NewBruteOptimizer()
	err = optimizer.Configure(config)
	if err != nil {
		return err
	}
	cycles, err := optimizer.Prepare(config)
	if err != nil {
		return err
	}
	print("done\n")

	bar := progressbar.Default(int64(cycles), "Running backtests")
	results := make([]perf.PerformanceReport, 0, cycles)
	_, err = optimizer.Start(context.Background())
	if err != nil {
		return err
	}

	bar.Finish()

	top := results[len(results)-1]
	perf.PrintSummary(top)
	fmt.Println(top.Description)

	if err := writeReports(_outputDir, results); err != nil {
		return err
	}

	return nil
}

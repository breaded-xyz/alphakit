package main

import (
	"log"
	"os"

	"github.com/colngroup/zero2algo/internal/studyrun"
	"github.com/davecgh/go-spew/spew"
)

const _outputDir = ".out"

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {

	print("Reading study config... ")
	config, err := studyrun.ReadConfig(args[1])
	if err != nil {
		return err
	}
	print("done\n")

	spew.Dump(config)

	/*print("Reading prices... ")
	//_, err := studyrun.readPrices(args[0])
	if err != nil {
		return err
	}
	print("done\n")



	print("Preparing optimizer...\n")
	//optimizer := optimize.NewBruteOptimizer()
	//err = optimizer.Configure(config)
	//if err != nil {
	//	return err
	//}
	///cycles, err := optimizer.Prepare(config)
	if err != nil {
		return err
	}
	print("done\n")

	//bar := progressbar.Default(int64(cycles), "Running backtests")
	//results := make([]perf.PerformanceReport, 0, cycles)
	//_, err = optimizer.Start(context.Background())
	if err != nil {
		return err
	}

	//bar.Finish()

	//if err := writeReports(_outputDir, results); err != nil {
	//	return err
	//}

	return nil*/
	return nil
}

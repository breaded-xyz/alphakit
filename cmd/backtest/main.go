package main

import (
	"log"
	"os"

	"github.com/colngroup/zero2algo/optimize"
	"github.com/colngroup/zero2algo/perf"
)

var ()

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	prices, err := readPrices(_priceDir)
	if err != nil {
		return err
	}

	testCases := optimize.BuildBacktestCases(_params)
	results := make([]perf.PerformanceReport, 0, len(testCases))
	for i := range testCases {
		dealer := _dealer
		bot := _bot
		if err := bot.Configure(testCases[i]); err != nil {
			return err
		}
		result, err := execBacktest(&bot, &dealer, prices)
		if err != nil {
			return err
		}
		results = append(results, result)
	}

	optimize.SharpeSort(results)
	perf.PrintSummary(results[len(results)-1])

	return nil
}

package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/colngroup/zero2algo/optimize"
	"github.com/colngroup/zero2algo/perf"
	"github.com/colngroup/zero2algo/trader"
	"github.com/gammazero/workerpool"
	"github.com/schollz/progressbar/v3"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	print("Reading prices... ")
	prices, err := readPrices(_priceDir)
	if err != nil {
		return err
	}
	print("done\n")

	print("Generating backtest cases...\n")
	testCases := optimize.BuildBacktestCases(_params)
	bar := progressbar.Default(int64(len(testCases)), "Running backtests")

	results := make([]perf.PerformanceReport, 0, len(testCases))
	wp := workerpool.New(16)
	for i := range testCases {
		i := i
		wp.Submit(func() {
			tCase := testCases[i]
			dealer := _dealerMakeFunc()
			if err := dealer.Configure(tCase); err != nil {
				if errors.Is(err, trader.ErrInvalidConfig) {
					return
				}
				panic(err)
			}

			bot := _botMakeFunc()
			bot.SetDealer(dealer)
			if err := bot.Configure(tCase); err != nil {
				if errors.Is(err, trader.ErrInvalidConfig) {
					return
				}
				panic(err)
			}

			result, err := execBacktest(bot, dealer, prices)
			if err != nil {
				panic(err)
			}
			result.Description = fmt.Sprintf("%+v", tCase)
			results = append(results, result)

			bar.Add(1)
		})
	}

	wp.StopWait()
	bar.Finish()

	optimize.SharpeSort(results)
	print("Top strategy by Sharpe:\n")
	top := results[len(results)-1]
	perf.PrintSummary(top)
	fmt.Println(top.Description)

	if err := writeReports(_outputDir, results); err != nil {
		return err
	}

	return nil
}

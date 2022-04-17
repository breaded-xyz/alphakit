package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/colngroup/zero2algo/broker/backtest"
	"github.com/colngroup/zero2algo/optimize"
	"github.com/colngroup/zero2algo/perf"
	"github.com/colngroup/zero2algo/trader"
	"github.com/gammazero/workerpool"
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
	prices, err := readPrices(args[0])
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

	print("Generating backtest cases...\n")
	testCases := optimize.BruteCaseBuilder(config)
	bar := progressbar.Default(int64(len(testCases)), "Running backtests")

	results := make([]perf.PerformanceReport, 0, len(testCases))
	wp := workerpool.New(16)
	var mu sync.Mutex
	for i := range testCases {
		i := i
		wp.Submit(func() {
			tCase := testCases[i]
			dealer := backtest.NewDealer()
			if err := dealer.Configure(tCase); err != nil {
				if errors.Is(err, trader.ErrInvalidConfig) {
					return
				}
				panic(err)
			}

			bot := _typeRegistry[config["bot"].(string)].(botMakerFunc)()
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

			mu.Lock()
			results = append(results, result)
			mu.Unlock()

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

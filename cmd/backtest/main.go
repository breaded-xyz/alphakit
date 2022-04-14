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
)

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

	wp := workerpool.New(16)

	for _, tCase := range testCases {

		tCase := tCase

		wp.Submit(func() {

			dealer := _dealerMakeFunc()
			if err := dealer.Configure(tCase); err != nil {
				//ch <- err
			}

			bot := _botMakeFunc()
			bot.SetDealer(dealer)
			if err := bot.Configure(tCase); err != nil {
				if errors.Is(err, trader.ErrInvalidConfig) {
					return
				}
				//ch <- err
			}

			result, err := execBacktest(bot, dealer, prices)
			if err != nil {
				//ch <- err
			}
			result.Strategy = fmt.Sprintf("%+v", tCase)
			results = append(results, result)

			fmt.Printf("Done: %s\n", result.Strategy)
		})

	}

	wp.StopWait()

	optimize.SharpeSort(results)
	top := results[len(results)-1]
	perf.PrintSummary(top)
	fmt.Println(top.Strategy)

	return nil
}

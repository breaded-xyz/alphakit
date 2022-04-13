package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/colngroup/zero2algo/optimize"
	"github.com/colngroup/zero2algo/perf"
	"github.com/colngroup/zero2algo/trader"
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
	testCaseCount := len(testCases)

	results := make([]perf.PerformanceReport, 0, len(testCases))
	for i := range testCases {
		tCase := testCases[i]

		dealer := _dealerMakeFunc()
		if err := dealer.Configure(tCase); err != nil {
			return err
		}

		bot := _botMakeFunc()
		bot.SetDealer(dealer)
		if err := bot.Configure(tCase); err != nil {
			if errors.Is(err, trader.ErrInvalidConfig) {
				continue
			}
			return err
		}

		result, err := execBacktest(bot, dealer, prices)
		if err != nil {
			return err
		}
		result.Strategy = fmt.Sprintf("%+v", tCase)
		results = append(results, result)

		fmt.Printf("%d/%d complete: %s\n", i+1, testCaseCount, result.Strategy)
	}

	optimize.SharpeSort(results)
	top := results[len(results)-1]
	perf.PrintSummary(top)
	fmt.Println(top.Strategy)

	return nil
}

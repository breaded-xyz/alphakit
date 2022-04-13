package main

import (
	"context"
	"encoding/csv"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/colngroup/zero2algo/broker/backtest"
	"github.com/colngroup/zero2algo/dec"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/optimize"
	"github.com/colngroup/zero2algo/perf"
	"github.com/colngroup/zero2algo/trader/hodl"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {

	pricePath := args[0]

	params := defineParamRange()
	testCases := optimize.BuildBacktestCases(params)

	results := make([]perf.PerformanceReport, 0, len(testCases))

	prices, err := readPrices(pricePath)
	if err != nil {
		return err
	}

	wg := new(sync.WaitGroup)
	//resultCh := make(chan perf.PerformanceReport)

	for i := range testCases {
		wg.Add(1)

		go func(prices []market.Kline, tCase optimize.TestCase) {
			defer wg.Done()

			// Create a special simulated dealer for each test case run
			dealer := backtest.NewDealer()
			dealer.SetInitialCapital(dec.New(1000))

			// Create a new bot initialized with our dealer
			// Hodl Bot implements a basic buy and hold algo
			bot := hodl.New(market.NewAsset("BTCUSD"), dealer)
			// The bot is configured with the params in the test case
			bot.Configure(tCase)

			// Iterate prices sending each price interval to the dealer and then to the bot
			for _, price := range prices {
				dealer.ReceivePrice(context.Background(), price)
				bot.ReceivePrice(context.Background(), price)
			}
			// Close the bot which will liquidate any open position resulting in a final trade
			bot.Close(context.Background())

			// Generate a performance report for the test case and add it to the result set
			trades, _, _ := dealer.ListTrades(context.Background(), nil)
			equity := dealer.EquityHistory()
			report := perf.NewPerformanceReport(trades, equity)
			perf.PrintSummary(report)
		}(prices, testCases[i])
	}
	//results = append(results, <-resultCh)
	wg.Wait()

	// Rank results based on the test case with the highest sharpe ratio
	optimize.SharpeSort(results)
	//perf.PrintSummary(results[len(results)-1])

	return nil
}

func defineParamRange() optimize.ParamRange {
	return map[string][]any{
		hodl.BuyBarIndex:  {0, 1, 1000},
		hodl.SellBarIndex: {0, 1000, 2000},
	}
}

func readPrices(path string) ([]market.Kline, error) {
	var prices []market.Kline

	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		klines, err := market.NewCSVKlineReader(csv.NewReader(file)).ReadAll()
		if err != nil {
			return err
		}
		prices = append(prices, klines...)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return prices, nil
}

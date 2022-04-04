package zero2algo

import (
	"context"
	"encoding/csv"
	"io"
	"os"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/broker/backtest"
	"github.com/colngroup/zero2algo/dec"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/perf"
)

func ExampleBacktest() {
	// Verbose error handling ommitted for brevity

	// Create a special simulated dealer for backtesting
	dealer := backtest.NewDealer()

	// Identify the asset to trade
	asset := market.NewAsset("BTCUSD")

	// Read a .csv file of historical prices (aka candlestick data)
	file, _ := os.Open("prices.csv")
	defer file.Close()
	reader := market.NewCSVKlineReader(csv.NewReader(file))

	// Iterate prices sending each price interval to the backtest dealer
	// When connected to a live exchange we would not be required to supply the price to the dealer!
	var i int
	for {
		price, err := reader.Read()
		if err == io.EOF {
			break
		}
		dealer.ReceivePrice(context.Background(), price)
		if i == 0 {
			dealer.PlaceOrder(context.Background(), broker.NewOrder(asset, broker.Buy, dec.New(1)))
		}
		i++
	}
	// Close the position and create the trade
	dealer.PlaceOrder(context.Background(), broker.NewOrder(asset, broker.Sell, dec.New(1)))

	// Generate a performance report from the dealer execution history
	trades, _, _ := dealer.ListTrades(context.Background(), nil)
	report := perf.NewPerformanceReport(trades, dealer.Equity())
	perf.PrintPerformanceReportSummary(report)
}

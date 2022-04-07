package backtest

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/dec"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/perf"
)

func Example() {
	// Verbose error handling ommitted for brevity

	// Create a special simulated dealer for backtesting with initial capital of 1000
	dealer := NewDealer()
	dealer.SetAccountBalance(dec.New(1000))

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

		// Place an order for 1 BTC at start of the price series
		if i == 0 {
			dealer.PlaceOrder(context.Background(), broker.NewOrder(asset, broker.Buy, dec.New(1)))
		}
		i++
	}
	// Close the position and create the trade
	dealer.PlaceOrder(context.Background(), broker.NewOrder(asset, broker.Sell, dec.New(1)))

	// Generate a performance report from the dealer execution history
	trades, _, _ := dealer.ListTrades(context.Background(), nil)
	equity := dealer.Equity()
	report := perf.NewPerformanceReport(trades, equity)
	fmt.Printf("Your backtest return is %.2f%%", report.Portfolio.EquityReturn*100)
}

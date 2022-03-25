package zero2algo

import (
	"context"
	"encoding/csv"
	"io"
	"os"

	"github.com/colngroup/zero2algo/broker/backtest"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/perf"
	"github.com/colngroup/zero2algo/tradebot"
)

func ExampleBacktest() {
	// Verbose error handling ommitted for brevity

	// Create a special simulated dealer for backtesting
	dealer := backtest.NewDealer()

	// Identify the asset to trade
	asset := market.NewAsset("BTCUSD")

	// Create a new bot initialized with our dealer
	// HodlBot implements a basic buy and hold algo
	bot := tradebot.NewHodlBot(asset, dealer)

	// Read a .csv file of historical prices (aka candlestick data)
	file, _ := os.Open("prices.csv")
	defer file.Close()
	reader := market.NewCSVKlineReader(csv.NewReader(file))

	// Iterate prices sending each price interval to the backtest dealer and then to the bot
	// When connected to a live exchange we would not be required to supply the price to the dealer!
	for {
		price, err := reader.Read()
		if err == io.EOF {
			break
		}
		_ = dealer.ReceivePrice(context.Background(), price)
		_ = bot.ReceivePrice(context.Background(), price)
	}
	// Close the bot which will liquidate the held position resulting in a trade
	_ = bot.Close(context.Background())

	// Generate a performance report from the dealer execution history
	trades, _, _ := dealer.ListTrades(context.Background(), nil)
	report := perf.NewReport(trades, dealer.ListEquityHistory())
	perf.PrintReportSummary(report)
}

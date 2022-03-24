package zero2algo

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/colngroup/zero2algo/broker/backtest"
	"github.com/colngroup/zero2algo/perf"
	"github.com/colngroup/zero2algo/price"
	"github.com/colngroup/zero2algo/tradebot"
)

func ExampleBacktest() {
	// Verbose error handling ommitted for brevity

	// Create a special simulated dealer for backtesting
	dealer := backtest.NewDealer()

	// Create a new bot initialized with our dealer
	// HodlBot implements a basic buy and hold algo
	bot := tradebot.NewHodlBot(dealer)

	// Read a .csv file of historical prices (aka candlestick data)
	file, _ := os.Open("prices.csv")
	defer file.Close()
	reader := price.NewCSVKlineReader(csv.NewReader(file))

	// Iterate prices sending each price interval to the backtest dealer and then to the bot
	// When connected to a live exchange we would not be required to supply the price to the dealer!
	for {
		price, err := reader.Read()
		if err == io.EOF {
			break
		}
		_ = dealer.ReceivePrice(price)
		_ = bot.ReceivePrice(price)
	}
	// Close the bot which will liquidate the held position resulting in a trade
	bot.Close()

	// Generate a performance report from the dealer execution history
	report := perf.NewReport(dealer.ListTrades(), dealer.EquityCurve())
	perf.PrintReportSummary(report)
}

package zero2algo

func Example() {
	// Create a special simulated dealer for backtesting
	dealer := NewBacktestDealer()

	// Create a new bot initialized with our dealer
	// HodlBot implements a basic buy and hold algo
	bot := NewHodlBot(dealer)

	// Read a .csv file of historical prices into a slice of klines (aka candlestick data)
	prices, _ := ReadKlinesFromCSV("prices.csv")

	// Iterate prices sending each price interval to the backtest dealer and then to the bot
	for _, price := range prices {
		_ = dealer.Receive(price)
		_ = bot.Receive(price)
	}

	// Generate a performance report once all price data has been iterated
	report := NewReport(dealer.ListTradeHistory(), dealer.ListEquityHistory())
	PrintReportSummary(report)
}

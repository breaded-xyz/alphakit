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
	// In the real world with the dealer connected to an exchange we would not be required to supply the price!
	for _, price := range prices {
		_ = dealer.ReceivePrice(price)
		_ = bot.ReceivePrice(price)
	}

	// Generate a performance report once all price data has been iterated
	report := NewReport(dealer.ListTrades(), dealer.EquityCurve())
	PrintReportSummary(report)
}

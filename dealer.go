package zero2algo

type Dealer interface {
	ListTrades() []Trade
	ListEquityHistory() []Equity
}

type SimulatedDealer interface {
	Dealer
	PriceReceiver
}

package zero2algo

type Dealer interface {
	ListTrades() []Trade
}

type SimulatedDealer interface {
	Dealer
	PriceReceiver
	EquityCurve() []Equity
}

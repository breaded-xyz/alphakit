package zero2algo

type Dealer interface {
	PriceReceiver
	ListTradeHistory() []Trade
	ListEquityHistory() []Equity
}

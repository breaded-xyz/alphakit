package broker

import "github.com/colngroup/zero2algo/price"

type Dealer interface {
	ListTrades() []Trade
}

type SimulatedDealer interface {
	Dealer
	price.Receiver
	EquityCurve() []Equity
}

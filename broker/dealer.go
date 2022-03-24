package broker

import "github.com/colngroup/zero2algo/pricing"

type Dealer interface {
	ListTrades() []Trade
}

type SimulatedDealer interface {
	Dealer
	pricing.Receiver
	EquityCurve() []Equity
}

package broker

import (
	"context"

	"github.com/colngroup/zero2algo/net"
	"github.com/colngroup/zero2algo/pricing"
)

type Dealer interface {
	ListTrades(context.Context, *net.ListOpts) ([]Trade, *net.Response, error)
}

type SimulatedDealer interface {
	Dealer
	pricing.Receiver
	EquityCurve() []Equity
}

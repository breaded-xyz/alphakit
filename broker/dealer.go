package broker

import (
	"context"

	"github.com/colngroup/zero2algo/netapi"
	"github.com/colngroup/zero2algo/pricing"
)

type Dealer interface {
	ListTrades(context.Context, *netapi.ListOpts) ([]Trade, *netapi.Response, error)
}

type SimulatedDealer interface {
	Dealer
	pricing.Receiver
	ListEquityHistory() []Equity
}

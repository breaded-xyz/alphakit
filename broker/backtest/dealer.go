package backtest

import (
	"context"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/netapi"
)

// Enforce at compile time that the type implements the interface
var _ broker.SimulatedDealer = (*Dealer)(nil)

type Dealer struct {
}

func NewDealer() *Dealer {
	return nil
}

func (d *Dealer) ListTrades(ctx context.Context, opts *netapi.ListOpts) ([]broker.Trade, *netapi.Response, error) {
	return nil, nil, nil
}

func (d *Dealer) ListEquityHistory() []broker.Equity {
	return nil
}

func (d *Dealer) ReceivePrice(ctx context.Context, price market.Kline) error {
	return nil
}

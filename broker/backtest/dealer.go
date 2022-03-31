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
	simulator *Simulator
}

func NewDealer() *Dealer {
	return &Dealer{
		simulator: NewSimulator(),
	}
}

func (d *Dealer) PlaceOrder(ctx context.Context, order broker.Order) (*broker.Order, *netapi.Response, error) {
	order, err := d.simulator.AddOrder(order)
	return &order, nil, err
}

func (d *Dealer) ListPositions(ctx context.Context, opts *netapi.ListOpts) ([]broker.Position, *netapi.Response, error) {
	return d.simulator.Positions(), nil, nil
}

func (d *Dealer) ListTrades(ctx context.Context, opts *netapi.ListOpts) ([]broker.Trade, *netapi.Response, error) {
	return d.simulator.Trades(), nil, nil
}

func (d *Dealer) Equity() broker.EquitySeries {
	return d.simulator.Equity()
}

func (d *Dealer) ReceivePrice(ctx context.Context, price market.Kline) error {
	return d.simulator.Next(price)
}

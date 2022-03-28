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

func (d *Dealer) PlaceOrder(ctx context.Context, order broker.Order) (*broker.Order, *netapi.Response, error) {
	return nil, nil, nil
}

func (d *Dealer) ListPositions(ctx context.Context, opts *netapi.ListOpts) ([]broker.Position, *netapi.Response, error) {
	return nil, nil, nil
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

func (d *Dealer) processOrder(order broker.Order) broker.Order {
	switch order.State() {
	case broker.Pending:
		// Open order and add to set of working orders
	case broker.Open:
		// Evaluate price match
	case broker.Filled:
		// Update positions
	case broker.Closed:
		// Move to set of closed orders
	}
	return order
}

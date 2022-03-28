package backtest

import (
	"context"
	"errors"
	"time"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/dec"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/netapi"
)

var ErrInvalidOrderState = errors.New("order is not valid for processing")

// Enforce at compile time that the type implements the interface
var _ broker.SimulatedDealer = (*Dealer)(nil)

type Dealer struct {
	simulationTime time.Time
	marketPrice    market.Kline

	openOrders map[broker.DealID]broker.Order
}

func NewDealer() *Dealer {
	return &Dealer{
		openOrders: make(map[broker.DealID]broker.Order),
	}
}

func (d *Dealer) PlaceOrder(ctx context.Context, order broker.Order) (*broker.Order, *netapi.Response, error) {
	if order.Side == 0 || order.Type == 0 || order.State() != broker.Pending || !order.Size.IsPositive() {
		return nil, nil, ErrInvalidOrderState
	}

	order = d.processOrder(order, d.simulationTime, d.marketPrice)
	d.updateOpenOrders(order)

	return &order, nil, nil
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

	// Iterate open orders in OpenedAt order and process

	return nil
}

func (d *Dealer) processOrder(order broker.Order, t time.Time, price market.Kline) broker.Order {

	switch order.State() {
	case broker.Pending:
		order.ID = broker.NewID()
		order.OpenedAt = t
		order = d.processOrder(order, t, price)
	case broker.Open:
		switch order.Type {
		case broker.Limit:
			if !dec.Between(order.LimitPrice, price.L, price.H) {
				return order
			}
			order.FilledPrice = order.LimitPrice
		case broker.Market:
			order.FilledPrice = price.C
		}
		order.FilledAt = t
		order.FilledSize = order.Size
		order = d.processOrder(order, t, price)
	case broker.Filled:
		order.ClosedAt = t
	}

	return order
}

func (d *Dealer) updateOpenOrders(order broker.Order) {
	switch order.State() {
	case broker.Open:
		d.openOrders[order.ID] = order
	case broker.Closed:
		delete(d.openOrders, order.ID)
	}
}

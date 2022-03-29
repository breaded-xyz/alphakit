package backtest

import (
	"context"
	"errors"
	"time"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/dec"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/netapi"
	"github.com/shopspring/decimal"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

var ErrInvalidOrderState = errors.New("order is not valid for processing")

// Enforce at compile time that the type implements the interface
var _ broker.SimulatedDealer = (*Dealer)(nil)

type Dealer struct {
	simulationTime time.Time
	marketPrice    market.Kline

	orders map[broker.DealID]broker.Order
}

func NewDealer() *Dealer {
	return &Dealer{
		simulationTime: time.Now().UTC(),
		orders:         make(map[broker.DealID]broker.Order),
	}
}

func (d *Dealer) PlaceOrder(ctx context.Context, order broker.Order) (*broker.Order, *netapi.Response, error) {
	if order.Side == 0 || order.Type == 0 || order.State() != broker.Pending || !order.Size.IsPositive() {
		return nil, nil, ErrInvalidOrderState
	}

	order = d.processOrder(order)

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

	d.simulationTime = closeTime(d.marketPrice.Start, price.Start)
	d.marketPrice = price

	// Iterate open orders in the order they were placed with the dealer
	// Go maps do not maintain insertion order so we must sort the keys in a slice first
	// The key is a ULID seeded from a time and supports lexicographic sorting
	ks := maps.Keys(d.orders)
	slices.Sort(ks)
	for _, k := range ks {
		order := d.orders[k]
		if order.State() == broker.Open {
			d.processOrder(order)
		}
	}

	// Calculate equity history point

	return nil
}

func (d *Dealer) processOrder(order broker.Order) broker.Order {

	switch order.State() {
	case broker.Pending:
		order = d.processOrder(d.openOrder(order))
	case broker.Open:
		if matchedPrice := matchOrder(order, d.marketPrice); matchedPrice.IsPositive() {
			order = d.processOrder(d.fillOrder(order, matchedPrice))
		}
	case broker.Filled:
		order = d.closeOrder(order)
	}

	return order
}

func (d *Dealer) openOrder(order broker.Order) broker.Order {
	order.ID = broker.NewID()
	order.OpenedAt = d.simulationTime
	d.orders[order.ID] = order
	return order
}

func (d *Dealer) fillOrder(order broker.Order, matchedPrice decimal.Decimal) broker.Order {
	order.FilledAt = d.simulationTime
	order.FilledPrice = matchedPrice
	order.FilledSize = order.Size
	d.orders[order.ID] = order
	return order
}

func (d *Dealer) closeOrder(order broker.Order) broker.Order {
	order.ClosedAt = d.simulationTime
	d.orders[order.ID] = order
	return order
}

func matchOrder(order broker.Order, quote market.Kline) decimal.Decimal {
	var matchedPrice decimal.Decimal

	switch order.Type {
	case broker.Limit:
		if dec.Between(order.LimitPrice, quote.L, quote.H) {
			matchedPrice = order.LimitPrice
		}
	case broker.Market:
		matchedPrice = quote.C
	}

	return matchedPrice
}

func closeTime(prevStart, currStart time.Time) time.Time {
	if prevStart.IsZero() || currStart.Before(prevStart) {
		return currStart
	}
	interval := currStart.UTC().Sub(prevStart.UTC())
	return currStart.UTC().Add(interval)
}

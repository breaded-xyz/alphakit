package backtest

import (
	"context"
	"log"
	"time"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/dec"
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
		log.Default().Println("Process Pending Order")
		order.OpenedAt = d.simulationTime()
		d.processOrder(order)
	case broker.Open:
		log.Default().Println("Process Open Order")
		if order.Type == broker.Limit {
			if !dec.Between(order.LimitPrice, d.marketPrice().L, d.marketPrice().H) {
				break
			}
		}
		order.FilledAt = d.simulationTime()
		order.FilledPrice = d.marketPrice().C
		order.FilledSize = order.Size
		d.processOrder(order)
	case broker.Filled:
		log.Default().Println("Process Filled Order")
		order.ClosedAt = d.simulationTime()
		d.processOrder(order)
	case broker.Closed:
		// Move to set of closed orders
		log.Default().Println("Process Closed Order")
	}
	return order
}

func (d *Dealer) simulationTime() time.Time {
	return time.Now().UTC()
}

func (d *Dealer) marketPrice() market.Kline {
	return market.Kline{}
}

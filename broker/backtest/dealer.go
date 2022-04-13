package backtest

import (
	"context"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/netapi"
	"github.com/shopspring/decimal"
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

func NewDealerWithCost(cost Coster) *Dealer {
	return &Dealer{
		simulator: NewSimulatorWithCost(cost),
	}
}

func (d *Dealer) Configure(config map[string]any) error {
	return d.simulator.Configure(config)
}

func (d *Dealer) SetInitialCapital(amount decimal.Decimal) {
	d.simulator.SetInitialCapital(amount)
}

func (d *Dealer) GetBalance(ctx context.Context) (*broker.AccountBalance, *netapi.Response, error) {
	acc := d.simulator.Balance()
	return &acc, nil, nil
}

func (d *Dealer) PlaceOrder(ctx context.Context, order broker.Order) (*broker.Order, *netapi.Response, error) {
	order, err := d.simulator.AddOrder(order)
	return &order, nil, err
}

func (d *Dealer) CancelOrders(ctx context.Context) (*netapi.Response, error) {
	d.simulator.CancelOrders()
	return nil, nil
}

func (d *Dealer) ListPositions(ctx context.Context, opts *netapi.ListOpts) ([]broker.Position, *netapi.Response, error) {
	return d.simulator.Positions(), nil, nil
}

func (d *Dealer) ListTrades(ctx context.Context, opts *netapi.ListOpts) ([]broker.Trade, *netapi.Response, error) {
	return d.simulator.Trades(), nil, nil
}

func (d *Dealer) EquityHistory() broker.EquitySeries {
	return d.simulator.EquityHistory()
}

func (d *Dealer) ReceivePrice(ctx context.Context, price market.Kline) error {
	return d.simulator.Next(price)
}

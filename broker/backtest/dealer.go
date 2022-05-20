package backtest

import (
	"context"

	"github.com/shopspring/decimal"
	"github.com/thecolngroup/alphakit/broker"
	"github.com/thecolngroup/alphakit/market"
	"github.com/thecolngroup/alphakit/web"
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

func (d *Dealer) SetInitialCapital(amount decimal.Decimal) {
	d.simulator.SetInitialCapital(amount)
}

func (d *Dealer) GetBalance(ctx context.Context) (*broker.AccountBalance, *web.Response, error) {
	acc := d.simulator.Balance()
	return &acc, nil, nil
}

func (d *Dealer) PlaceOrder(ctx context.Context, order broker.Order) (*broker.Order, *web.Response, error) {
	order, err := d.simulator.AddOrder(order)
	return &order, nil, err
}

func (d *Dealer) CancelOrders(ctx context.Context) (*web.Response, error) {
	d.simulator.CancelOrders()
	return nil, nil
}

func (d *Dealer) ListPositions(ctx context.Context, opts *web.ListOpts) ([]broker.Position, *web.Response, error) {
	return d.simulator.Positions(), nil, nil
}

func (d *Dealer) ListTrades(ctx context.Context, opts *web.ListOpts) ([]broker.Trade, *web.Response, error) {
	return d.simulator.Trades(), nil, nil
}

func (d *Dealer) EquityHistory() broker.EquitySeries {
	return d.simulator.EquityHistory()
}

func (d *Dealer) ReceivePrice(ctx context.Context, price market.Kline) error {
	return d.simulator.Next(price)
}

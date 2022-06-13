// Package backtest provides a simultated dealer implementation for running backtests.
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

// Dealer is a SimulatedDealer implementation for backtesting.
type Dealer struct {
	simulator *Simulator
}

// NewDealer creates a new Dealer with a default simulator state.
func NewDealer() *Dealer {
	return &Dealer{
		simulator: NewSimulator(),
	}
}

// NewDealerWithCost creates a new Dealer with the given cost model.
func NewDealerWithCost(cost Coster) *Dealer {
	return &Dealer{
		simulator: NewSimulatorWithCost(cost),
	}
}

// SetInitialCapital sets the initial trading balance for the dealer.
func (d *Dealer) SetInitialCapital(amount decimal.Decimal) {
	d.simulator.SetInitialCapital(amount)
}

// GetBalance returns the current balance of the dealer.
func (d *Dealer) GetBalance(ctx context.Context) (*broker.AccountBalance, *web.Response, error) {
	acc := d.simulator.Balance()
	return &acc, nil, nil
}

// PlaceOrder places an order on the dealer.
func (d *Dealer) PlaceOrder(ctx context.Context, order broker.Order) (*broker.Order, *web.Response, error) {
	order, err := d.simulator.AddOrder(order)
	return &order, nil, err
}

// CancelOrders cancels all open (resting) orders on the dealer.
func (d *Dealer) CancelOrders(ctx context.Context) (*web.Response, error) {
	d.simulator.CancelOrders()
	return nil, nil
}

// ListPositions returns all historical (closed) and open positions.
func (d *Dealer) ListPositions(ctx context.Context, opts *web.ListOpts) ([]broker.Position, *web.Response, error) {
	return d.simulator.Positions(), nil, nil
}

// ListRoundTurns returns all historical round-turns.
func (d *Dealer) ListRoundTurns(ctx context.Context, opts *web.ListOpts) ([]broker.RoundTurn, *web.Response, error) {
	return d.simulator.RoundTurns(), nil, nil
}

// EquityHistory returns the equity history (equity curve) of the dealer.
func (d *Dealer) EquityHistory() broker.EquitySeries {
	return d.simulator.EquityHistory()
}

// ReceivePrice initiates processing by supplying the next market price to the simulator.
func (d *Dealer) ReceivePrice(ctx context.Context, price market.Kline) error {
	return d.simulator.Next(price)
}

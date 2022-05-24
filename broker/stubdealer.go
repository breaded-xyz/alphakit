package broker

import (
	"context"

	"github.com/shopspring/decimal"
	"github.com/thecolngroup/alphakit/market"
	"github.com/thecolngroup/alphakit/web"
)

var _ SimulatedDealer = (*StubDealer)(nil)

// StubDealer is a test double for a simulated dealer that does nothing.
type StubDealer struct {
}

// ReceivePrice not implemented.
func (d *StubDealer) ReceivePrice(ctx context.Context, price market.Kline) error {
	return nil
}

// SetInitialCapital not implemented.
func (d *StubDealer) SetInitialCapital(amount decimal.Decimal) {}

// GetBalance not implemented.
func (d *StubDealer) GetBalance(ctx context.Context) (*AccountBalance, *web.Response, error) {
	return nil, nil, nil
}

// PlaceOrder not implemented.
func (d *StubDealer) PlaceOrder(ctx context.Context, order Order) (*Order, *web.Response, error) {
	return nil, nil, nil
}

// CancelOrders not implemented.
func (d *StubDealer) CancelOrders(ctx context.Context) (*web.Response, error) {
	return nil, nil
}

// ListPositions not implemented.
func (d *StubDealer) ListPositions(ctx context.Context, opts *web.ListOpts) ([]Position, *web.Response, error) {
	return nil, nil, nil
}

// ListTrades not implemented.
func (d *StubDealer) ListTrades(ctx context.Context, opts *web.ListOpts) ([]Trade, *web.Response, error) {
	return nil, nil, nil
}

// EquityHistory not implemented.
func (d *StubDealer) EquityHistory() EquitySeries {
	return nil
}

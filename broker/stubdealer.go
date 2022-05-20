package broker

import (
	"context"

	"github.com/shopspring/decimal"
	"github.com/thecolngroup/alphakit/market"
	"github.com/thecolngroup/alphakit/web"
)

var _ SimulatedDealer = (*StubDealer)(nil)

type StubDealer struct {
}

func (d *StubDealer) ReceivePrice(ctx context.Context, price market.Kline) error {
	return nil
}

func (d *StubDealer) Configure(map[string]any) error { return nil }

func (d *StubDealer) SetInitialCapital(amount decimal.Decimal) {}

func (d *StubDealer) GetBalance(ctx context.Context) (*AccountBalance, *web.Response, error) {
	return nil, nil, nil
}

func (d *StubDealer) PlaceOrder(ctx context.Context, order Order) (*Order, *web.Response, error) {
	return nil, nil, nil
}

func (d *StubDealer) CancelOrders(ctx context.Context) (*web.Response, error) {
	return nil, nil
}

func (d *StubDealer) ListPositions(ctx context.Context, opts *web.ListOpts) ([]Position, *web.Response, error) {
	return nil, nil, nil
}

func (d *StubDealer) ListTrades(ctx context.Context, opts *web.ListOpts) ([]Trade, *web.Response, error) {
	return nil, nil, nil
}

func (d *StubDealer) EquityHistory() EquitySeries {
	return nil
}

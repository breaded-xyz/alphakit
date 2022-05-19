package broker

import (
	"context"

	"github.com/shopspring/decimal"
	"github.com/thecolngroup/zerotoalgo/market"
	"github.com/thecolngroup/zerotoalgo/netapi"
)

var _ SimulatedDealer = (*StubDealer)(nil)

type StubDealer struct {
}

func (d *StubDealer) ReceivePrice(ctx context.Context, price market.Kline) error {
	return nil
}

func (d *StubDealer) Configure(map[string]any) error { return nil }

func (d *StubDealer) SetInitialCapital(amount decimal.Decimal) {}

func (d *StubDealer) GetBalance(ctx context.Context) (*AccountBalance, *netapi.Response, error) {
	return nil, nil, nil
}

func (d *StubDealer) PlaceOrder(ctx context.Context, order Order) (*Order, *netapi.Response, error) {
	return nil, nil, nil
}

func (d *StubDealer) CancelOrders(ctx context.Context) (*netapi.Response, error) {
	return nil, nil
}

func (d *StubDealer) ListPositions(ctx context.Context, opts *netapi.ListOpts) ([]Position, *netapi.Response, error) {
	return nil, nil, nil
}

func (d *StubDealer) ListTrades(ctx context.Context, opts *netapi.ListOpts) ([]Trade, *netapi.Response, error) {
	return nil, nil, nil
}

func (d *StubDealer) EquityHistory() EquitySeries {
	return nil
}

package broker

import (
	"context"

	"github.com/colngroup/zero2algo/netapi"
	"github.com/stretchr/testify/mock"
)

var _ Dealer = (*MockDealer)(nil)

type MockDealer struct {
	mock.Mock
}

func (d *MockDealer) PlaceOrder(ctx context.Context, order Order) (*Order, *netapi.Response, error) {
	args := d.Called(ctx, order)

	if len(args) == 0 {
		return nil, nil, nil
	}

	return args.Get(0).(*Order), args.Get(1).(*netapi.Response), args.Error(2)
}

func (d *MockDealer) ListPositions(ctx context.Context, opts *netapi.ListOpts) ([]Position, *netapi.Response, error) {
	args := d.Called(ctx, opts)
	return args.Get(0).([]Position), args.Get(1).(*netapi.Response), args.Error(2)
}

func (d *MockDealer) ListTrades(ctx context.Context, opts *netapi.ListOpts) ([]Trade, *netapi.Response, error) {
	args := d.Called(ctx, opts)
	return args.Get(0).([]Trade), args.Get(1).(*netapi.Response), args.Error(2)
}

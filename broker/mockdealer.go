package broker

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/thecolngroup/alphakit/web"
)

var _ Dealer = (*MockDealer)(nil)

type MockDealer struct {
	mock.Mock
}

func (d *MockDealer) GetBalance(ctx context.Context) (*AccountBalance, *web.Response, error) {
	args := d.Called(ctx)
	return args.Get(0).(*AccountBalance), args.Get(1).(*web.Response), args.Error(2)
}

func (d *MockDealer) PlaceOrder(ctx context.Context, order Order) (*Order, *web.Response, error) {
	args := d.Called(ctx, order)

	if len(args) == 0 {
		return nil, nil, nil
	}

	return args.Get(0).(*Order), args.Get(1).(*web.Response), args.Error(2)
}

func (d *MockDealer) CancelOrders(ctx context.Context) (*web.Response, error) {
	args := d.Called(ctx)
	return args.Get(0).(*web.Response), args.Error(1)
}

func (d *MockDealer) ListPositions(ctx context.Context, opts *web.ListOpts) ([]Position, *web.Response, error) {
	args := d.Called(ctx, opts)
	return args.Get(0).([]Position), args.Get(1).(*web.Response), args.Error(2)
}

func (d *MockDealer) ListTrades(ctx context.Context, opts *web.ListOpts) ([]Trade, *web.Response, error) {
	args := d.Called(ctx, opts)
	return args.Get(0).([]Trade), args.Get(1).(*web.Response), args.Error(2)
}

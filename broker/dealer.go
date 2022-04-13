package broker

import (
	"context"

	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/netapi"
	"github.com/shopspring/decimal"
)

type Dealer interface {
	GetBalance(context.Context) (*AccountBalance, *netapi.Response, error)
	PlaceOrder(context.Context, Order) (*Order, *netapi.Response, error)
	CancelOrders(context.Context) (*netapi.Response, error)
	ListPositions(context.Context, *netapi.ListOpts) ([]Position, *netapi.Response, error)
	ListTrades(context.Context, *netapi.ListOpts) ([]Trade, *netapi.Response, error)
}

type SimulatedDealer interface {
	Dealer
	market.Receiver
	Configure(map[string]any) error
	EquityHistory() EquitySeries
	SetInitialCapital(amount decimal.Decimal)
}

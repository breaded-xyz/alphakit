package broker

import (
	"context"

	"github.com/shopspring/decimal"
	"github.com/thecolngroup/zerotoalgo/market"
	"github.com/thecolngroup/zerotoalgo/netapi"
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
	EquityHistory() EquitySeries
	SetInitialCapital(amount decimal.Decimal)
}

type MakeSimulatedDealer func() (SimulatedDealer, error)

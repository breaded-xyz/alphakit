// Package broker provides an API for interacting with 3rd party exchanges,
// and a simulated dealer for backtesting in the child package backtest.
package broker

import (
	"context"

	"github.com/shopspring/decimal"
	"github.com/thecolngroup/alphakit/market"
	"github.com/thecolngroup/alphakit/web"
)

// Dealer is an interface for interacting with a 3rd party exchange and placing orders in the market.
type Dealer interface {
	GetBalance(context.Context) (*AccountBalance, *web.Response, error)
	PlaceOrder(context.Context, Order) (*Order, *web.Response, error)
	CancelOrders(context.Context) (*web.Response, error)
	ListPositions(context.Context, *web.ListOpts) ([]Position, *web.Response, error)
	ListRoundTurns(context.Context, *web.ListOpts) ([]RoundTurn, *web.Response, error)
}

// SimulatedDealer is a Dealer that can be used for backtesting.
type SimulatedDealer interface {
	Dealer
	market.Receiver
	EquityHistory() EquitySeries
	SetInitialCapital(amount decimal.Decimal)
}

// MakeSimulatedDealer returns a new SimulatedDealer
type MakeSimulatedDealer func() (SimulatedDealer, error)

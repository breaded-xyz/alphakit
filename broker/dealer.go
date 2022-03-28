package broker

import (
	"context"

	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/netapi"
	"github.com/google/uuid"
)

type DealID string

func NewID() DealID {
	return DealID(uuid.NewString())
}

type Dealer interface {
	PlaceOrder(context.Context, Order) (*Order, *netapi.Response, error)
	ListPositions(context.Context, *netapi.ListOpts) ([]Position, *netapi.Response, error)
	ListTrades(context.Context, *netapi.ListOpts) ([]Trade, *netapi.Response, error)
}

type SimulatedDealer interface {
	Dealer
	market.Receiver
	ListEquityHistory() []Equity
}

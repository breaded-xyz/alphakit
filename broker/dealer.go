package broker

import (
	"context"
	"math/rand"
	"time"

	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/netapi"
	"github.com/oklog/ulid/v2"
)

type DealID string

func NewID() DealID {
	t := time.Now().UTC()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	return DealID(ulid.MustNew(ulid.Timestamp(t), entropy).String())
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

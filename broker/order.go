package broker

import (
	"time"

	"github.com/colngroup/zero2algo/market"
	"github.com/shopspring/decimal"
)

type OrderSide int

const (
	Buy OrderSide = iota + 1
	Sell
)

func (s OrderSide) String() string {
	return [...]string{"Buy", "Sell"}[s]
}

type OrderType int

const (
	Market OrderType = iota + 1
	Limit
)

func (s OrderType) String() string {
	return [...]string{"Market", "Limit"}[s]
}

type OrderState int

const (
	Pending = iota
	Open
	Filled
	Closed
)

func (s OrderState) String() string {
	return [...]string{"Pending", "Open", "Filled", "Closed"}[s]
}

type Order struct {
	ID       DealID
	OpenedAt time.Time
	FilledAt time.Time
	ClosedAt time.Time

	Asset      market.Asset
	Side       OrderSide
	Type       OrderType
	LimitPrice decimal.Decimal
	Size       decimal.Decimal
	ReduceOnly bool

	FilledPrice decimal.Decimal
	FilledSize  decimal.Decimal
}

func NewOrder(asset market.Asset, side OrderSide, size decimal.Decimal) Order {
	return Order{
		Asset: asset,
		Side:  side,
		Size:  size,
	}
}

func (o *Order) State() OrderState {
	switch {
	case !o.ClosedAt.IsZero():
		return Closed
	case !o.FilledAt.IsZero():
		return Filled
	case !o.OpenedAt.IsZero():
		return Open
	}
	return Pending
}

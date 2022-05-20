package broker

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/thecolngroup/alphakit/market"
)

type OrderSide int

const (
	Buy OrderSide = iota + 1
	Sell
)

func (s OrderSide) String() string {
	return [...]string{"None", "Buy", "Sell"}[s]
}

func (s OrderSide) Opposite() OrderSide {
	switch s {
	case Buy:
		return Sell
	case Sell:
		return Buy
	}
	return 0
}

type OrderType int

const (
	Market OrderType = iota + 1
	Limit
)

func (s OrderType) String() string {
	return [...]string{"None", "Market", "Limit"}[s]
}

type OrderState int

const (
	OrderPending = iota
	OrderOpen
	OrderFilled
	OrderClosed
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
		Type:  Market,
	}
}

func (o *Order) State() OrderState {
	switch {
	case !o.ClosedAt.IsZero():
		return OrderClosed
	case !o.FilledAt.IsZero():
		return OrderFilled
	case !o.OpenedAt.IsZero():
		return OrderOpen
	}
	return OrderPending
}

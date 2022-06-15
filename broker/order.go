package broker

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/thecolngroup/alphakit/market"
)

// OrderSide represents the side of an order: Buy (long) or Sell (short).
type OrderSide int

const (
	// Buy (long)
	Buy OrderSide = iota + 1

	// Sell (short)
	Sell
)

func (s OrderSide) String() string {
	return [...]string{"None", "Buy", "Sell"}[s]
}

// MarshalText is used to output as a string for CSV rendering.
func (s OrderSide) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

// Opposite returns the opposite side of the order.
func (s OrderSide) Opposite() OrderSide {
	switch s {
	case Buy:
		return Sell
	case Sell:
		return Buy
	}
	return 0
}

// OrderType represents an order type
type OrderType int

const (
	// Market order type is executed at market price (taker).
	Market OrderType = iota + 1

	// Limit order type is executed at a specified price (maker).
	Limit
)

func (t OrderType) String() string {
	return [...]string{"None", "Market", "Limit"}[t]
}

// MarshalText is used to output as a string for CSV rendering.
func (t OrderType) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

// OrderState represents the state of an order as it is processed by a dealer.
type OrderState int

const (
	// OrderPending represents an order that has not been processed by a dealer.
	OrderPending = iota

	// OrderOpen represents an order that has been opened by a dealer but not yet filled.
	OrderOpen

	// OrderFilled represents an order that has been filled by a dealer at a price level.
	OrderFilled

	// OrderClosed represents an order that has been closed by a dealer.
	OrderClosed
)

func (s OrderState) String() string {
	return [...]string{"Pending", "Open", "Filled", "Closed"}[s]
}

// Order represents an order to be placed using a dealer.
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

	Fee decimal.Decimal
}

// NewOrder creates a new order with the minimum required fields to be valid.
func NewOrder(asset market.Asset, side OrderSide, size decimal.Decimal) Order {
	return Order{
		Asset: asset,
		Side:  side,
		Size:  size,
		Type:  Market,
	}
}

// State returns the state of the order based on the order timestamps.
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

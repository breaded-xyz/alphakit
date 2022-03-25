package broker

import (
	"github.com/colngroup/zero2algo/market"
	"github.com/shopspring/decimal"
)

type OrderSide int

const (
	Buy OrderSide = iota + 1
	Sell
)

type Order struct {
	Asset      market.Asset
	Side       OrderSide
	Size       decimal.Decimal
	ReduceOnly bool
}

func NewOrder(asset market.Asset, side OrderSide, size decimal.Decimal) Order {
	return Order{
		Asset: asset,
		Side:  side,
		Size:  size,
	}
}

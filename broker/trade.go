package broker

import (
	"time"

	"github.com/colngroup/zero2algo/market"
	"github.com/shopspring/decimal"
)

type Trade struct {
	ID        DealID
	CreatedAt time.Time
	Asset     market.Asset
	Side      OrderSide
	Size      decimal.Decimal
	Profit    decimal.Decimal
}

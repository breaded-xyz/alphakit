package broker

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/thecolngroup/zerotoalgo/market"
)

type Trade struct {
	ID         DealID
	CreatedAt  time.Time
	Asset      market.Asset
	Side       OrderSide
	Size       decimal.Decimal
	Profit     decimal.Decimal
	HoldPeriod time.Duration
}

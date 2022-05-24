package broker

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/thecolngroup/alphakit/market"
)

// Trade is the result of opening and closing a position i.e. a roundtrip / roundturn.
type Trade struct {
	ID         DealID
	CreatedAt  time.Time
	Asset      market.Asset
	Side       OrderSide
	Size       decimal.Decimal
	Profit     decimal.Decimal
	HoldPeriod time.Duration
}

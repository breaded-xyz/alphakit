package broker

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/thecolngroup/alphakit/market"
)

// Trade is the result of opening and closing a position i.e. a roundtrip / roundturn.
type Trade struct {
	ID         DealID          `csv:"id"`
	CreatedAt  time.Time       `csv:"created_at"`
	Asset      market.Asset    `csv:",inline"`
	Side       OrderSide       `csv:"side"`
	Size       decimal.Decimal `csv:"size"`
	Profit     decimal.Decimal `csv:"profit"`
	HoldPeriod time.Duration   `csv:"hold_period"`
}

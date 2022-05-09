package risk

import (
	"context"

	"github.com/colngroup/zero2algo/market"
	"github.com/shopspring/decimal"
)

// FullRisker assumes that the current price is the maximum downside risk.
// This is equivalent to assumming a position size of 1 and no stopp loss.
type FullRisker struct {
	price decimal.Decimal
}

// NewFullRisker returns a new FullRisker.
func NewFullRisker() *FullRisker {
	return &FullRisker{}
}

// ReceivePrice updates the FullRisker with the next price.
func (r *FullRisker) ReceivePrice(ctx context.Context, price market.Kline) error {
	r.price = price.C
	return nil
}

// Risk returns a unitary measure of risk based on the current price.
// Will always be the current price.
func (r *FullRisker) Risk() decimal.Decimal {
	return r.price
}

// Valid returns true if the risker has enough data to be calculated.
// Will always return true.
func (r *FullRisker) Valid() bool {
	return true
}

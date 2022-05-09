package risk

import (
	"context"

	"github.com/colngroup/zero2algo/internal/dec"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/ta"
	"github.com/shopspring/decimal"
)

var _ Risker = (*SDRisker)(nil)

// SDRisker is a Risker that uses the standard deviation of a moving window.
type SDRisker struct {
	SD *ta.SD
}

// NewSDRisker returns a new SDRisker.
func NewSDRisker(length int, factor float64) *SDRisker {
	return &SDRisker{
		SD: ta.NewSDWithFactor(length, factor),
	}
}

// ReceivePrice updates the SDRisker with the next price.
func (r *SDRisker) ReceivePrice(ctx context.Context, price market.Kline) error {
	return r.SD.Update(price.C.InexactFloat64())
}

// Risk returns a unitary measure of risk based on the current price.
func (r *SDRisker) Risk() decimal.Decimal {
	return dec.New(r.SD.Value())
}

// Valid returns true if the risker has enough data to be calculated.
func (r *SDRisker) Valid() bool {
	return r.SD.Valid()
}

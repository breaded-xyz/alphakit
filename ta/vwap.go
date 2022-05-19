package ta

import (
	"github.com/thecolngroup/zerotoalgo/market"
)

var _ Indicator[market.Kline] = (*VWAP)(nil)

// VWAP is a volume weighted average price.
type VWAP struct {
	cumPV  float64
	cumVol float64
	series []float64
}

// NewVWAP creates a new VWAP indicator with default parameters.
func NewVWAP() *VWAP {
	return &VWAP{}
}

// Update updates the indicator with the next value(s).
func (ind *VWAP) Update(prices ...market.Kline) error {

	for i := range prices {
		avgPrice := HLC3(prices[i])
		vol := prices[i].Volume

		if avgPrice == 0 || vol == 0 {
			continue
		}

		ind.cumPV += avgPrice * vol
		ind.cumVol += vol
		vwap := ind.cumPV / ind.cumVol
		ind.series = append(ind.series, vwap)
	}

	return nil
}

// Valid returns true if the indicator is valid.
// An indicator is invalid if it hasn't received enough values yet.
func (ind *VWAP) Valid() bool {
	return len(ind.series) > 0
}

// Value returns the current value of the indicator.
func (ind *VWAP) Value() float64 {
	return Lookback(ind.series, 0)
}

// History returns the historical values of the indicator.
func (ind *VWAP) History() []float64 {
	return ind.series
}

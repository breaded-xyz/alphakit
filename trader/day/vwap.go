package day

import (
	"log"
	"math"

	"github.com/colngroup/zero2algo/internal/util"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/ta"
)

//var _ ta.Indicator = (*VWAP)(nil)

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
		ind.cumPV += ta.HLC3(prices[i]) * util.NNZ(prices[i].Volume, 1)
		ind.cumVol += prices[i].Volume
		vwap := ind.cumPV / ind.cumVol
		if math.IsNaN(vwap) {
			log.Fatalf("price: %+v cumPV: %f cumVol: %f", prices[i], ind.cumPV, ind.cumVol)

		}
		ind.series = append(ind.series, vwap)
	}

	return nil
}

// Valid returns true if the indicator is valid.
// An indicator is invalid if it hasn't received enough values yet.
func (ind *VWAP) Valid() bool {
	return true
}

// Value returns the current value of the indicator.
func (ind *VWAP) Value() float64 {
	return ta.Lookback(ind.series, 0)
}

// History returns the historical values of the indicator.
func (ind *VWAP) History() []float64 {
	return ind.series
}

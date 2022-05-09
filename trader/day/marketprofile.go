package day

import (
	"sort"

	"github.com/gonum/floats"
	"github.com/gonum/stat"
)

// MarketProfile is a histogram of market price and volume for a session.
// Intent is to show the price points with most volume during the session.
// The profile gives key features such as:
//
// Point of control (POC)
//
// Value area high (VAH)
//
// Value area low (VAL)
//
// Session High
//
// Session Low
type MarketProfile struct {
	Bins []float64
	Hist []float64

	POC  float64
	VAH  float64
	VAL  float64
	High float64
	Low  float64
}

// NewMarketProfile creates a new MarketProfile for the given price series.
func NewMarketProfile(nBins int, prices, volumes []float64) *MarketProfile {

	var mp MarketProfile
	mp.Bins = make([]float64, nBins+1)

	sort.Float64s(prices)
	mp.High = floats.Max(prices)
	mp.Low = floats.Min(prices)

	mp.Bins = floats.Span(mp.Bins, mp.Low, mp.High+1)

	mp.Hist = stat.Histogram(nil, mp.Bins, prices, volumes)

	return &mp
}

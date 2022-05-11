package day

import (
	"github.com/colngroup/zero2algo/internal/util"
	"github.com/davecgh/go-spew/spew"
	"github.com/gonum/floats"
	"github.com/gonum/stat"
)

// DefaultValueAreaPercentage is the percentage of the total volume used to calculate the value area.
const DefaultValueAreaPercentage = 0.7

// VolumeProfile is a histogram of market price and volume.
// Intent is to show the price points with most volume during a period.
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
type VolumeProfile struct {
	Bins []float64
	Hist []float64

	POC  float64
	VAH  float64
	VAL  float64
	High float64
	Low  float64
}

type Level struct {
	Price  float64
	Volume float64
}

// NewVolumeProfile creates a new profile for the given price and volume series.
// nBins is the number of bins to use, higher numbers for greater accurracy.
// Prices and volumes must be of the same length.
func NewVolumeProfile(nBins int, levels []Level) *VolumeProfile {

	var vp VolumeProfile

	var sortedPrices, volumes []float64
	for _, level := range levels {
		sortedPrices = append(sortedPrices, level.Price)
		volumes = append(volumes, level.Volume)
	}

	vp.High = floats.Max(sortedPrices)
	vp.Low = floats.Min(sortedPrices)
	vp.Bins = make([]float64, nBins)
	vp.Bins = floats.Span(vp.Bins, vp.Low, vp.High+1)

	vp.Hist = stat.Histogram(nil, vp.Bins, sortedPrices, volumes)

	pocIdx := floats.MaxIdx(vp.Hist)
	vp.POC = stat.Mean([]float64{vp.Bins[pocIdx], vp.Bins[pocIdx+1]}, nil)

	vaTotalVol := util.RoundTo(floats.Sum(volumes)*DefaultValueAreaPercentage, 1)
	vaCumVol := vp.Hist[pocIdx]

	spew.Dump(vaTotalVol, vaCumVol)

	var vaIdx int

	for i := 1; vaCumVol <= vaTotalVol; i++ {

		if pocIdx+i == 0 {
		}

		vaCumVol += vp.Hist[pocIdx+i] + vp.Hist[pocIdx-i]
		vaIdx = i
		spew.Dump(vaCumVol)
	}

	spew.Dump(vaIdx)

	return &vp
}

package day

import (
	"sort"

	"github.com/gonum/floats"
	"github.com/gonum/stat"
	"golang.org/x/exp/slices"
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

type level struct {
	Price  float64
	Volume float64
}

// NewVolumeProfile creates a new profile for the given price and volume series.
// nBins is the number of bins to use, higher numbers for greater accurracy.
// Prices and volumes must be of the same length.
func NewVolumeProfile(nBins int, prices, volumes []float64) *VolumeProfile {

	var vp VolumeProfile

	levels := make([]level, len(prices))
	for i := range prices {
		levels[i] = level{prices[i], volumes[i]}
	}
	slices.SortFunc(levels, func(i, j level) bool {
		return i.Price < j.Price
	})

	sortedPrices := prices
	sort.Float64s(sortedPrices)
	vp.High = floats.Max(sortedPrices)
	vp.Low = floats.Min(sortedPrices)
	vp.Bins = make([]float64, nBins+1)
	vp.Bins = floats.Span(vp.Bins, vp.Low, vp.High+1)

	sortedVolumes := make([]float64, len(volumes))
	for i := range levels {
		sortedVolumes[i] = levels[i].Volume
	}

	vp.Hist = stat.Histogram(nil, vp.Bins, sortedPrices, sortedVolumes)

	pocIdx := floats.MaxIdx(vp.Hist)
	vp.POC = vp.Bins[pocIdx]

	vaTotalVol := floats.Sum(sortedVolumes) * DefaultValueAreaPercentage

	vaCumVol := vp.Hist[pocIdx]
	var vaPOCOffsetIdx int

	for i := 1; vaCumVol <= vaTotalVol; i++ {

		var vahVol, valVol float64
		if (pocIdx + i) <= len(vp.Hist)-1 {
			vahVol = vp.Hist[pocIdx+i]
		}
		if (pocIdx - i) >= 0 {
			valVol = vp.Hist[pocIdx-i]
		}

		vaCumVol += (vahVol + valVol)
		vaPOCOffsetIdx = i
	}

	vahIdx := pocIdx + vaPOCOffsetIdx
	if vahIdx > len(vp.Bins)-1 {
		vahIdx = len(vp.Bins) - 1
	}

	valIdx := pocIdx - vaPOCOffsetIdx
	if valIdx < 0 {
		valIdx = 0
	}

	vp.VAH = vp.Bins[vahIdx]
	vp.VAL = vp.Bins[valIdx]

	return &vp
}

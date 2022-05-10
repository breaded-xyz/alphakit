package day

import (
	"sort"

	"github.com/gonum/floats"
	"github.com/gonum/stat"
	"golang.org/x/exp/slices"
)

const DefaultValueAreaPercentage = 0.7

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

type pricevol struct {
	Price  float64
	Volume float64
}

// NewMarketProfile creates a new MarketProfile for the given price series.
func NewMarketProfile(nBins int, prices, volumes []float64) *MarketProfile {

	var mp MarketProfile
	mp.Bins = make([]float64, nBins+1)

	pvs := make([]pricevol, len(prices))
	for i := range prices {
		pvs[i] = pricevol{prices[i], volumes[i]}
	}
	slices.SortFunc(pvs, func(i, j pricevol) bool {
		return i.Price < j.Price
	})

	sort.Float64s(prices)
	mp.High = floats.Max(prices)
	mp.Low = floats.Min(prices)
	mp.Bins = floats.Span(mp.Bins, mp.Low, mp.High+1)

	sortedVolumes := make([]float64, len(volumes))
	for i := range pvs {
		sortedVolumes[i] = pvs[i].Volume
	}

	//spew.Dump(prices[floats.MaxIdx(sortedVolumes)])

	mp.Hist = stat.Histogram(nil, mp.Bins, prices, sortedVolumes)

	pocIdx := floats.MaxIdx(mp.Hist)
	mp.POC = stat.Mean([]float64{mp.Bins[pocIdx]}, nil)

	vaTotalVol := floats.Sum(mp.Hist) * DefaultValueAreaPercentage

	vaCumVol := mp.Hist[pocIdx]
	var vaPOCOffset int

	for i := 1; vaCumVol <= vaTotalVol; i++ {
		vahVol := mp.Hist[pocIdx+i] + mp.Hist[pocIdx+(i+1)]
		valVol := mp.Hist[pocIdx-i] + mp.Hist[pocIdx-(i+1)]

		/*	//hVolIdx, lVolIdx := 1, 1
			//
					if pocIdx+hVolIdx < len(mp.Hist) {
			/			vahVol = mp.Hist[pocIdx+hVolIdx]
					}
					if pocIdx-lVolIdx >= 0 {
						valVol = mp.Hist[pocIdx-lVolIdx]
					}*/

		switch {
		case vahVol > valVol:
			//vahIdx++
			//valIdx++
			vaCumVol += vahVol
		case vahVol <= valVol:
			//valIdx--
			//vahIdx++
			vaCumVol += valVol
			//case vahVol == valVol:
			//	valIdx++
			//	vaCumVol += valVol
		}

		vaPOCOffset = i
	}

	mp.VAH = stat.Mean([]float64{mp.Bins[pocIdx+vaPOCOffset]}, nil)
	mp.VAL = stat.Mean([]float64{mp.Bins[pocIdx-vaPOCOffset]}, nil)

	return &mp
}

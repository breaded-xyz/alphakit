package ta

import (
	"math"
)

const (
	DefaultALMAOffset = 0.85
	DefaultALMASigma  = 6
)

var _ Indicator = (*ALMA)(nil)

// ALMA is a modern low lag moving average.
// Ported from https://www.tradingview.com/pine-script-reference/#fun_alma
type ALMA struct {
	sample []float64
	series []float64
	length int
	offset float64
	sigma  float64
}

func NewALMA(length int) *ALMA {
	return NewALMAWithSigma(length, DefaultALMAOffset, DefaultALMASigma)
}

func NewALMAWithSigma(length int, offset, sigma float64) *ALMA {
	return &ALMA{
		length: length,
		offset: offset,
		sigma:  sigma,
	}
}

func (ind *ALMA) Update(v ...float64) error {
	for i := range v {
		ind.sample = WindowAppend(ind.sample, ind.length-1, v[i])

		if ind.length < 1 {
			ind.series = append(ind.series, v[i])
			continue
		}

		m := math.Floor(ind.offset * (float64(ind.length) - 1))
		s := float64(ind.length) / ind.sigma
		var norm, sum float64
		for i := 0; i < len(ind.sample); i++ {
			weight := math.Exp(-1 * math.Pow(float64(i)-m, 2) / (2 * math.Pow(s, 2)))
			norm += weight
			sum += ind.sample[i] * weight
		}
		ma := sum / norm
		ind.series = append(ind.series, ma)
	}

	return nil
}
func (ind *ALMA) Valid() bool {
	return len(ind.sample) >= ind.length
}

func (ind *ALMA) Value() float64 {
	return Lookback(ind.series, 0)
}

func (ind *ALMA) History() []float64 {
	return ind.series
}

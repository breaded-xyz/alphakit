package ta

import (
	"math"
)

const (
	// DefaultALMAOffset is the default offset for the ALMA indicator.
	DefaultALMAOffset = 0.85

	// DefaultALMASigma is the default sigma for the ALMA indicator.
	DefaultALMASigma = 6
)

var _ Indicator[float64] = (*ALMA)(nil)

// ALMA is a modern low lag moving average.
// Ported from https://www.tradingview.com/pine-script-reference/#fun_alma
type ALMA struct {
	Length int
	Offset float64
	Sigma  float64

	sample []float64
	series []float64
}

// NewALMA creates a new ALMA indicator with default parameters.
func NewALMA(length int) *ALMA {
	return NewALMAWithSigma(length, DefaultALMAOffset, DefaultALMASigma)
}

// NewALMAWithSigma creates a new ALMA indicator with the given offset and sigma.
func NewALMAWithSigma(length int, offset, sigma float64) *ALMA {
	return &ALMA{
		Length: length,
		Offset: offset,
		Sigma:  sigma,
	}
}

// Update updates the indicator with the next value(s).
func (ind *ALMA) Update(v ...float64) error {
	for i := range v {
		ind.sample = WindowAppend(ind.sample, ind.Length-1, v[i])

		if ind.Length < 1 {
			ind.series = append(ind.series, v[i])
			continue
		}

		m := math.Floor(ind.Offset * (float64(ind.Length) - 1))
		s := float64(ind.Length) / ind.Sigma
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

// Valid returns true if the indicator is valid.
// An indicator is invalid if it hasn't received enough values yet.
func (ind *ALMA) Valid() bool {
	return len(ind.sample) >= ind.Length
}

// Value returns the current value of the indicator.
func (ind *ALMA) Value() float64 {
	return Lookback(ind.series, 0)
}

// History returns the historical values of the indicator.
func (ind *ALMA) History() []float64 {
	return ind.series
}

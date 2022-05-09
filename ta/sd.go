package ta

import "github.com/gonum/stat"

// SD is a sample standard deviation indicator.
type SD struct {
	// Length is the number of values to use in the calculation.
	Length int

	// Factor is the factor to multiply the standard deviation by.
	Factor float64

	sample []float64
	series []float64
}

// NewSD returns a new SD indicator with default factor of 1.
func NewSD(length int) *SD {
	return NewSDWithFactor(length, 1)
}

// NewSDWithFactor returns a new SD indicator with the given factor.
func NewSDWithFactor(length int, factor float64) *SD {
	return &SD{
		Length: length,
		Factor: factor,
	}
}

// Update updates the indicator with the next value(s).
func (ind *SD) Update(v ...float64) error {

	for i := range v {
		ind.sample = WindowAppend(ind.sample, ind.Length-1, v[i])
		sd := stat.StdDev(ind.sample, nil)
		sd *= ind.Factor
		ind.series = append(ind.series, sd)
	}

	return nil
}

// Valid returns true if the indicator has enough data to be calculated.
func (ind *SD) Valid() bool {
	return len(ind.sample) >= ind.Length
}

// Value returns the current value of the indicator.
func (ind *SD) Value() float64 {
	return Lookback(ind.series, 0)
}

// History returns the history of the indicator.
func (ind *SD) History() []float64 {
	return ind.series
}

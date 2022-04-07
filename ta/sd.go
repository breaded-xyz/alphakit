package ta

import "github.com/gonum/stat"

type SD struct {
	sample []float64
	series []float64
	length int
	factor float64
}

func NewSD(length int) *SD {
	return NewSDWithFactor(length, 1)
}

func NewSDWithFactor(length int, factor float64) *SD {
	return &SD{
		length: length,
		factor: factor,
	}
}

func (ind *SD) Update(v ...float64) error {

	for i := range v {
		ind.sample = WindowAppend(ind.sample, ind.length-1, v[i])
		sd := stat.StdDev(ind.sample, nil)
		sd *= ind.factor
		ind.series = append(ind.series, sd)
	}

	return nil
}
func (ind *SD) Valid() bool {
	return len(ind.sample) >= ind.length
}

func (ind *SD) Value() float64 {
	return Lookback(ind.series, 0)
}

func (ind *SD) History() []float64 {
	return ind.series
}

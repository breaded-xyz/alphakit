package ta

import (
	"github.com/iamjinlei/go-tart"
)

var _ Indicator = (*KAMA)(nil)

type KAMA struct {
	ma     *tart.Ma
	series []float64
	length int
}

func NewKAMA(length int) Indicator {
	kama := &KAMA{
		length: length,
		ma:     tart.NewMa(tart.KAMA, int64(length)),
	}
	return kama
}

func (ind *KAMA) Update(v ...float64) error {
	for i := range v {
		ind.series = append(ind.series, ind.ma.Update(v[i]))
	}
	return nil
}

func (ind *KAMA) Value() float64 {
	return Lookback(ind.series, 0)
}

func (ind *KAMA) History() []float64 {
	return ind.series
}

func (ind *KAMA) Valid() bool {
	return len(ind.series) >= ind.length
}

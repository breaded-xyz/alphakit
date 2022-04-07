package ta

import (
	"sort"

	"github.com/gonum/stat"
)

// MMI (Market Meaness Index) is a statistical measure between 0 - 100
// that indicates if the series exhibits serial correlation (trendiness).
// Reference: https://financial-hacker.com/the-market-meanness-index/
type MMI struct {
	sample []float64
	length int

	smoother Indicator
}

func NewMMI(length int) *MMI {
	return &MMI{
		length:   length,
		smoother: NewALMA(length),
	}
}

func NewMMIWithSmoother(length int, smoother Indicator) *MMI {
	return &MMI{
		length:   length,
		smoother: smoother,
	}
}

func (ind *MMI) Update(v ...float64) error {

	for i := range v {
		ind.sample = WindowAppend(ind.sample, ind.length-1, v[i])

		m := Median(ind.sample)
		var nh, nl float64
		for i := 1; i < len(ind.sample); i++ {
			p1, p0 := Lookback(ind.sample, i), Lookback(ind.sample, i-1)
			if p1 > m && p1 > p0 {
				nl++
			} else if p1 < m && p1 < p0 {
				nh++
			}
		}
		mmi := (nl + nh) / float64(len(ind.sample)-1)
		if err := ind.smoother.Update(mmi); err != nil {
			return err
		}
	}

	return nil
}
func (ind *MMI) Valid() bool {
	return len(ind.sample) >= ind.length
}

func (ind *MMI) Value() float64 {
	return ind.smoother.Value()
}

func (ind *MMI) History() []float64 {
	return ind.smoother.History()
}

func Median(v []float64) float64 {
	x := make([]float64, len(v))
	copy(x, v)
	sort.Float64s(x)
	return stat.Quantile(0.5, stat.Empirical, x, nil)
}

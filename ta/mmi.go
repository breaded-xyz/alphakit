package ta

var _ Indicator[float64] = (*MMI)(nil)

// MMI (Market Meaness Index) is a statistical measure between 0 - 100
// that indicates if the series exhibits serial correlation (trendiness).
// Reference: https://financial-hacker.com/the-market-meanness-index/.
type MMI struct {
	// Length is the number of values to use for the calculation.
	Length int

	// Smoother is the indicator used to smooth the MMI.
	Smoother Indicator[float64]

	sample []float64
}

// NewMMI returns a new MMI indicator with a default ALMA smoother.
// The smoothing length is the same as the given MMI length.
func NewMMI(length int) *MMI {
	return &MMI{
		Length:   length,
		Smoother: NewALMA(length),
	}
}

// NewMMIWithSmoother returns a new MMI indicator with the given smoother.
func NewMMIWithSmoother(length int, smoother Indicator[float64]) *MMI {
	return &MMI{
		Length:   length,
		Smoother: smoother,
	}
}

// Update updates the indicator with the next value(s).
func (ind *MMI) Update(v ...float64) error {

	for i := range v {
		ind.sample = WindowAppend(ind.sample, ind.Length-1, v[i])

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
		if err := ind.Smoother.Update(mmi); err != nil {
			return err
		}
	}

	return nil
}

// Valid returns true if the indicator has enough data to be calculated.
func (ind *MMI) Valid() bool {
	return len(ind.sample) >= ind.Length
}

// Value returns the current value of the indicator.
func (ind *MMI) Value() float64 {
	return ind.Smoother.Value()
}

// History returns the historical data of the indicator.
func (ind *MMI) History() []float64 {
	return ind.Smoother.History()
}

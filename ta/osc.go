package ta

var _ Indicator = (*ALMA)(nil)

// Osc is a composite of a fast and slow moving average indicator.
// Osc value = fast value - slow value.
// Osc is not normalized and has an unbounded range.
type Osc struct {
	series []float64
	fast   Indicator
	slow   Indicator
}

func NewOsc(fast, slow Indicator) *Osc {
	return &Osc{
		fast: fast,
		slow: slow,
	}
}

func (ind *Osc) Update(v ...float64) error {
	for i := range v {
		if err := ind.fast.Update(v[i]); err != nil {
			return err
		}
		if err := ind.slow.Update(v[i]); err != nil {
			return err
		}
		ind.series = append(ind.series, ind.fast.Value()-ind.slow.Value())
	}

	return nil
}
func (ind *Osc) Valid() bool {
	return ind.fast.Valid() && ind.slow.Valid()
}

func (ind *Osc) Value() float64 {
	return Lookback(ind.series, 0)
}

func (ind *Osc) History() []float64 {
	return ind.series
}

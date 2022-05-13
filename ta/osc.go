package ta

var _ Indicator[float64] = (*Osc)(nil)

// Osc is a composite of a fast and slow moving average indicator.
// Osc value = fast value minus slow value.
// Osc is not normalized and has an unbounded range.
type Osc struct {
	// Fast is the fast moving average indicator.
	Fast Indicator[float64]

	// Slow is the slow moving average indicator.
	Slow Indicator[float64]

	series []float64
}

// NewOsc returns a new oscillator with the given fast and slow moving averages.
func NewOsc(fast, slow Indicator[float64]) *Osc {
	return &Osc{
		Fast: fast,
		Slow: slow,
	}
}

// Update updates the indicator with the next value(s).
func (ind *Osc) Update(v ...float64) error {
	for i := range v {
		if err := ind.Fast.Update(v[i]); err != nil {
			return err
		}
		if err := ind.Slow.Update(v[i]); err != nil {
			return err
		}
		ind.series = append(ind.series, ind.Fast.Value()-ind.Slow.Value())
	}

	return nil
}

// Valid returns true if the indicator has enough data to be calculated.
func (ind *Osc) Valid() bool {
	return ind.Fast.Valid() && ind.Slow.Valid()
}

// Value returns the current value of the indicator.
func (ind *Osc) Value() float64 {
	return Lookback(ind.series, 0)
}

// History returns the history of the indicator.
func (ind *Osc) History() []float64 {
	return ind.series
}

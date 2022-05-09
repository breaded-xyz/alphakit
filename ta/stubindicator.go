package ta

// StubIndicator is a test double for an indicator.
type StubIndicator struct {
	// Values is the history of the indicator.
	Values []float64

	// IsValid is the validity of the indicator.
	IsValid bool
}

// Update is not implemented.
func (ind *StubIndicator) Update(v ...float64) error {
	return nil
}

// Valid returns IsValid.
func (ind *StubIndicator) Valid() bool {
	return ind.IsValid
}

// Value returns the latest value in Values
func (ind *StubIndicator) Value() float64 {
	return Lookback(ind.Values, 0)
}

// History returns the history of the indicator.
func (ind *StubIndicator) History() []float64 {
	return ind.Values
}

package ta

type StubIndicator struct {
	Values  []float64
	IsValid bool
}

func (ind *StubIndicator) Update(v ...float64) error {
	return nil
}

func (ind *StubIndicator) Valid() bool {
	return ind.IsValid
}

func (ind *StubIndicator) Value() float64 {
	return Lookback(ind.Values, 0)
}

func (ind *StubIndicator) History() []float64 {
	return ind.Values
}

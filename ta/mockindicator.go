package ta

import (
	"github.com/stretchr/testify/mock"
)

// MockIndicator is a mock implementation of the Indicator interface.
type MockIndicator struct {
	mock.Mock
}

// Update updates the indicator with the next value(s).
func (ind *MockIndicator) Update(v ...float64) error {
	args := ind.Called(v)
	return args.Error(0)
}

// Valid returns true if the indicator has enough data to be calculated.
func (ind *MockIndicator) Valid() bool {
	args := ind.Called()
	return args.Bool(0)
}

// Value returns the current value of the indicator.
func (ind *MockIndicator) Value() float64 {
	args := ind.Called()
	return args.Get(0).(float64)
}

// History returns the history of the indicator.
func (ind *MockIndicator) History() []float64 {
	args := ind.Called()
	return args.Get(0).([]float64)
}

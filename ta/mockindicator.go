package ta

import (
	"github.com/stretchr/testify/mock"
)

type MockIndicator struct {
	mock.Mock
}

func (ind *MockIndicator) Update(v ...float64) error {
	args := ind.Called(v)
	return args.Error(0)
}

func (ind *MockIndicator) Valid() bool {
	args := ind.Called()
	return args.Bool(0)
}

func (ind *MockIndicator) Value() float64 {
	args := ind.Called()
	return args.Get(0).(float64)
}

func (ind *MockIndicator) History() []float64 {
	args := ind.Called()
	return args.Get(0).([]float64)
}

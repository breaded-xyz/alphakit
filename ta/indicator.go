package ta

import "errors"

var ErrInvalidIndicatorState = errors.New("indicator state is invalid")

type Indicator interface {
	Update(v ...float64) error
	Valid() bool
	Value() float64
	History() []float64
}

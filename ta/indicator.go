package ta

type Indicator interface {
	Update(v ...float64) error
	Valid() bool
	Value() float64
	History() []float64
}

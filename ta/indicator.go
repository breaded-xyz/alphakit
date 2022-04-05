package ta

type Indicator interface {
	Update(v ...float64) error
	Value() float64
}

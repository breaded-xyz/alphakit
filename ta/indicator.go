// Package ta offers technical analysis functions for price series data.
package ta

// Indicator is the interface for all technical analysis functions.
type Indicator[T any] interface {

	// Update the indicator with new inputs (typically a price series).
	Update(v ...T) error

	// Value returns the latest value of the indicator.
	Value() float64

	// History returns all the historical indicator values (including the latest value).
	History() []float64

	// Valid returns true if the indicator is valid.
	Valid() bool
}

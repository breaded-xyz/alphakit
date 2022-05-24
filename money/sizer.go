// Package money provides money management and position sizing
package money

import (
	"github.com/shopspring/decimal"
	"github.com/thecolngroup/dec"
	"github.com/thecolngroup/util"
)

// DefaultStepSize is the default step size (rounding) for position sizes.
const DefaultStepSize = 0.01

// Sizer is the interface for position sizing.
type Sizer interface {
	Size(price, capital, risk decimal.Decimal) decimal.Decimal
}

// FixedSizer is a fixed position sizing strategy.
// Size is fixed capital divided by price (rounded to StepSize).
type FixedSizer struct {
	FixedCapital decimal.Decimal
	StepSize     float64
}

// NewFixedSizer returns a new FixedSizer with the given fixed capital and a default step size.
func NewFixedSizer(capital decimal.Decimal) *FixedSizer {
	return &FixedSizer{
		FixedCapital: capital,
		StepSize:     DefaultStepSize,
	}
}

// Size returns the fixed capital divided by price (rounded to StepSize).
// Returns 0 if size is NaN/Inf.
func (s *FixedSizer) Size(price, capital, risk decimal.Decimal) decimal.Decimal {
	size := s.FixedCapital.Div(price).InexactFloat64()
	size = util.RoundTo(size, s.StepSize)
	return dec.New(util.NN(size, 0))
}

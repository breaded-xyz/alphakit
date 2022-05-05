package money

import (
	"github.com/colngroup/zero2algo/internal/dec"
	"github.com/colngroup/zero2algo/internal/util"
	"github.com/shopspring/decimal"
)

var _defaultStepSize = 0.001

type Sizer interface {
	Size(price, capital, risk decimal.Decimal) decimal.Decimal
}

type FixedSizer struct {
	FixedCapital decimal.Decimal
	StepSize     float64
}

func NewFixedSizer(capital decimal.Decimal) *FixedSizer {
	return &FixedSizer{
		FixedCapital: capital,
		StepSize:     _defaultStepSize,
	}
}

func (s *FixedSizer) Size(price, capital, risk decimal.Decimal) decimal.Decimal {
	size := s.FixedCapital.Div(price).InexactFloat64()
	size = util.RoundTo(size, s.StepSize)
	return dec.New(size)
}

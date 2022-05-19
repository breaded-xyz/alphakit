package money

import (
	"github.com/shopspring/decimal"
	"github.com/thecolngroup/zerotoalgo/internal/dec"
	"github.com/thecolngroup/zerotoalgo/internal/util"
)

var _defaultStepSize = 0.01

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
	return dec.New(util.NN(size, 0))
}

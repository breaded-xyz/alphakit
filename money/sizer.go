package money

import (
	"github.com/shopspring/decimal"
)

type Sizer interface {
	Size(price, capital, risk decimal.Decimal) decimal.Decimal
}

type FixedSizer struct {
	FixedCapital decimal.Decimal
}

func NewFixedSizer(capital decimal.Decimal) *FixedSizer {
	return &FixedSizer{
		FixedCapital: capital,
	}
}

func (s *FixedSizer) Size(price, capital, risk decimal.Decimal) decimal.Decimal {
	return s.FixedCapital.Div(price)
}

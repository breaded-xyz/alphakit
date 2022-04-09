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

func (s *FixedSizer) Size(price, capital, risk decimal.Decimal) decimal.Decimal {
	return s.FixedCapital.Div(price)
}

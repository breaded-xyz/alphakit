package money

import (
	"math"

	"github.com/colngroup/zero2algo/dec"
	"github.com/shopspring/decimal"
)

type Sizer interface {
	Size(price, capital, risk decimal.Decimal) decimal.Decimal
}

type SafeFSizer struct {
	InitialCapital decimal.Decimal
	F, ScaleF      float64
}

func (s *SafeFSizer) Size(price, capital, risk decimal.Decimal) decimal.Decimal {

	sqrtGrowthFactor := 1.0
	profit := capital.Sub(s.InitialCapital)
	if profit.IsPositive() {
		capitalGrowthFactor := 1 + profit.Div(capital).InexactFloat64()
		sqrtGrowthFactor = math.Sqrt(capitalGrowthFactor)
	}
	safeF := s.F * s.ScaleF * sqrtGrowthFactor
	margin := capital.InexactFloat64() * safeF

	size := margin / risk.InexactFloat64()

	return dec.New(size)
}

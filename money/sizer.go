package money

import (
	"math"

	"github.com/colngroup/zero2algo/dec"
	"github.com/shopspring/decimal"
)

type Sizer interface {
	Size(price, capital decimal.Decimal, risk float64) decimal.Decimal
}

type SafeFSizer struct {
	InitialCapital decimal.Decimal
	F, ScaleF      float64
}

func (s *SafeFSizer) Size(price, capital decimal.Decimal, risk float64) decimal.Decimal {

	sqrtGrowthFactor := 1.0
	profit := capital.Sub(s.InitialCapital)
	if profit.IsPositive() {
		capitalGrowthFactor := 1 + profit.Div(capital).InexactFloat64()
		sqrtGrowthFactor = math.Sqrt(capitalGrowthFactor)
	}
	safeF := s.F * s.ScaleF * sqrtGrowthFactor
	margin := capital.InexactFloat64() * safeF

	unitRiskAmount := price.InexactFloat64() * risk
	size := margin / unitRiskAmount

	return dec.New(size)
}

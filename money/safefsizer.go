// Copyright 2022 The Coln Group Ltd
// SPDX-License-Identifier: MIT

package money

import (
	"math"

	"github.com/shopspring/decimal"
	"github.com/thecolngroup/gou/dec"
	"github.com/thecolngroup/gou/num"
)

var _ Sizer = (*SafeFSizer)(nil)

// SafeFSizer is a Sizer that uses a fixed fraction method (e.g. Kelly / OptimalF) with a safety margin.
type SafeFSizer struct {
	InitialCapital decimal.Decimal
	F              float64
	ScaleF         float64
	StepSize       float64
}

// NewSafeFSizer returns a new SafeFSizer with the given initial capital, fixed fraction, scale factor, and step size.
func NewSafeFSizer(initialCapital decimal.Decimal, f, scaleF float64) *SafeFSizer {
	return &SafeFSizer{
		InitialCapital: initialCapital,
		F:              f,
		ScaleF:         scaleF,
		StepSize:       DefaultStepSize,
	}
}

// Size moderates a fixed fraction position size by:
//
// - scaling f using the ScaleF fraction (typically 0.5)
//
// - only using the square root of the profits as trading capital.
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

	size = num.RoundTo(size, s.StepSize)

	return dec.New(num.NN(size, 0))
}

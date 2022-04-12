package perf

import (
	"math"
	"time"

	"github.com/colngroup/zero2algo/broker"
	"github.com/gonum/stat"
)

const (
	SharpeAnnualRiskFreeRate  = 0.0
	SharpeDailyToAnnualFactor = 252
	SharpeDailyRiskFreeRate   = SharpeAnnualRiskFreeRate / SharpeDailyToAnnualFactor
)

func DiffReturns(curve broker.EquitySeries) []float64 {
	diffs := make([]float64, len(curve)-1)
	vs := curve.SortValuesByTime()
	for i := range vs {
		if i == 0 {
			continue
		}
		diffs[i-1] = (vs[i].Sub(vs[i-1]).Div(vs[i-1])).InexactFloat64()
	}
	return diffs
}

func ReduceEOD(curve broker.EquitySeries) broker.EquitySeries {
	reduced := make(broker.EquitySeries, 0)
	eodH, eodM := 0, 0 // End of day = midnight
	for k, v := range curve {
		h, m, _ := time.UnixMilli(int64(k)).Clock()
		if h == eodH && m == eodM {
			reduced[k] = v
		}
	}
	return reduced
}

// CAGR Compound Annual Growth Rate
func CAGR(initial, final float64, days int) float64 {
	growthRate := (final - initial) / initial
	x := 1 + growthRate
	y := 365.0 / float64(days)
	return math.Pow(x, y) - 1
}

// SharpeRatio is the annualised value using daily risk free rate and daily returns
func SharpeRatio(daily []float64, riskFreeRate float64) float64 {
	xr := make([]float64, len(daily)) // Excess returns
	for i := range daily {
		xr[i] = daily[i] - riskFreeRate
	}

	mxr := stat.Mean(xr, nil)                         // Mean excess returns
	sd := stat.StdDev(xr, nil)                        // SD excess returns
	dsr := mxr / sd                                   // Daily Sharpe
	return dsr * math.Sqrt(SharpeDailyToAnnualFactor) // Scale daily to annual
}

func CalmarRatio(cagr, mdd float64) float64 {
	return cagr / mdd
}

func KellyCriterion(profitFactor, winP float64) float64 {
	return (profitFactor*winP - (1 - winP)) / profitFactor
}

// NN (Not Number) returns y if x is NaN or Inf.
func NN(x, y float64) float64 {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return y
	}
	return x
}

// NNZ (Not Number or Zero) returns y if x is NaN or Inf or Zero.
func NNZ(x, y float64) float64 {
	if NN(x, y) == y || x == 0 {
		return y
	}
	return x
}

package perf

import (
	"math"

	"github.com/gonum/stat"
)

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

// SE (Standard Error)
func SE(xs []float64) float64 {
	sd := stat.StdDev(xs, nil)
	se := stat.StdErr(sd, float64(len(xs)))
	return se
}

// PRR (Pessimistic Return Ratio)
func PRR(profit, loss, winningN, losingN float64) float64 {
	winF := 1 / math.Sqrt(1+winningN)
	loseF := 1 / math.Sqrt(1+losingN)
	return (1 - winF) / (1 + loseF) * (1 + profit) / (1 + loss)
}

// StatN returns the statistically significant number of samples required based on the distribution of a series
// From: https://www.elitetrader.com/et/threads/minimum-number-of-trades-required-for-backtesting-results-to-be-trusted.356588/page-2
func StatN(xs []float64) float64 {
	sd := stat.StdDev(xs, nil)
	m := stat.Mean(xs, nil)
	n := math.Pow(4*(sd/m), 2)
	return n
}

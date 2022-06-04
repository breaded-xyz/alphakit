package perf

import (
	"math"

	"gonum.org/v1/gonum/stat"
)

const (
	// SharpeDefaultAnnualRiskFreeRate is the default risk free rate for Sharpe Ratio.
	SharpeDefaultAnnualRiskFreeRate = 0.0

	// SharpeDailyToAnnualFactor is the factor to convert daily Sharpe to annual Sharpe.
	SharpeDailyToAnnualFactor = 252

	// SharpeDefaultDailyRiskFreeRate is the daily rate based on the default annual rate
	SharpeDefaultDailyRiskFreeRate = SharpeDefaultAnnualRiskFreeRate / SharpeDailyToAnnualFactor
)

// SharpeRatio is the annualised value using a daily risk free rate and daily returns.
// Param daily is the percentage daily returns from the portfolio.
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

// CAGR Compound Annual Growth Rate
func CAGR(initial, final float64, days int) float64 {
	growthRate := (final - initial) / initial
	x := 1 + growthRate
	y := 365.0 / float64(days)
	return math.Pow(x, y) - 1
}

// CalmarRatio relates the capaital growth rate to the maximum drawdown.
func CalmarRatio(cagr, mdd float64) float64 {
	return cagr / mdd
}

// KellyCriterion is the famous method for trade sizing.
func KellyCriterion(profitFactor, winP float64) float64 {
	return (profitFactor*winP - (1 - winP)) / profitFactor
}

// SE (Standard Error)
func SE(xs []float64) float64 {
	sd := stat.StdDev(xs, nil)
	se := stat.StdErr(sd, float64(len(xs)))
	return se
}

// PRR (Pessimistic Return Ratio) is the profit factor with a penalty for a lower number of trades.
func PRR(profit, loss, winningN, losingN float64) float64 {
	winF := 1 / math.Sqrt(1+winningN)
	loseF := 1 / math.Sqrt(1+losingN)
	return (1 - winF) / (1 + loseF) * (1 + profit) / (1 + loss)
}

// StatN returns the statistically significant number of samples required based on the distribution of a series.
// From: https://www.elitetrader.com/et/threads/minimum-number-of-trades-required-for-backtesting-results-to-be-trusted.356588/page-2
func StatN(xs []float64) float64 {
	sd := stat.StdDev(xs, nil)
	m := stat.Mean(xs, nil)
	n := math.Pow(4*(sd/m), 2)
	return n
}

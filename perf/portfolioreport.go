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

type PortfolioReport struct {
	PeriodStart time.Time
	PeriodEnd   time.Time
	Period      time.Duration

	StartEquity  float64
	EndEquity    float64
	EquityReturn float64

	CAGR float64

	MaxDrawdown float64
	MDDRecovery time.Duration

	Sharpe float64
	Calmar float64

	drawdowns []Drawdown
	mdd       Drawdown
}

func NewPortfolioReport(curve broker.EquitySeries) *PortfolioReport {
	if len(curve) == 0 {
		return nil
	}

	t := curve.SortKeys()
	tStart, tEnd := t[0], t[len(t)-1]

	var report PortfolioReport
	report.PeriodStart, report.StartEquity = tStart.Time(), curve[tStart].InexactFloat64()
	report.PeriodEnd, report.EndEquity = tEnd.Time(), curve[tEnd].InexactFloat64()
	report.Period = report.PeriodEnd.Sub(report.PeriodStart)

	report.EquityReturn = (report.EndEquity - report.StartEquity) / NNZ(report.StartEquity, 1)
	report.CAGR = CAGR(report.StartEquity, report.EndEquity, int(report.Period.Hours())/24)

	report.drawdowns = Drawdowns(curve)
	report.mdd = MaxDrawdown(report.drawdowns)
	report.MaxDrawdown = report.mdd.Pct
	report.MDDRecovery = report.mdd.Recovery

	returns := DiffOfReturns(ReduceToEOD(curve))
	report.Sharpe = SharpeRatio(returns, SharpeAnnualRiskFreeRate)
	report.Calmar = CalmarRatio(report.CAGR, report.MaxDrawdown)

	return &report
}

func DiffOfReturns(curve broker.EquitySeries) []float64 {
	returns := make([]float64, len(curve)-1)
	vs := curve.SortValuesByTime()
	for i := range vs {
		if i == 0 {
			continue
		}
		returns[i-1] = (vs[i].Sub(vs[i-1]).Div(vs[i-1])).InexactFloat64()
	}
	return returns
}

func ReduceToEOD(curve broker.EquitySeries) broker.EquitySeries {
	reduced := make(broker.EquitySeries, 0)
	eodH, eodM := 0, 0 // End of day = midnight
	for k, v := range curve {
		h, m, _ := time.Unix(int64(k), 0).Clock()
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

package perf

import (
	"time"

	"github.com/colngroup/zero2algo/broker"
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
	report.CAGR = NN(CAGR(report.StartEquity, report.EndEquity, int(report.Period.Hours())/24), 0)

	daily := ReduceEOD(curve)
	if len(daily) == 0 {
		return &report
	}

	report.drawdowns = Drawdowns(daily)
	report.mdd = MaxDrawdown(report.drawdowns)
	report.MaxDrawdown = report.mdd.Pct
	report.MDDRecovery = report.mdd.Recovery

	returns := DiffReturns(daily)
	report.Sharpe = SharpeRatio(returns, SharpeAnnualRiskFreeRate)
	report.Calmar = CalmarRatio(report.CAGR, report.MaxDrawdown)

	return &report
}

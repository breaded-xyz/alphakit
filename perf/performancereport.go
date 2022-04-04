package perf

import "github.com/colngroup/zero2algo/broker"

type PerformanceReport struct {
	Trade     *TradeReport
	Portfolio *PortfolioReport
}

func NewPerformanceReport(trades []broker.Trade, equity broker.EquitySeries) PerformanceReport {
	return PerformanceReport{
		Trade:     NewTradeReport(trades),
		Portfolio: NewPortfolioReport(equity),
	}
}

func PrintPerformanceReportSummary(report PerformanceReport) {

}

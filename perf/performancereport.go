package perf

import "github.com/colngroup/zero2algo/broker"

type PerformanceReport struct {
	Portfolio PortfolioReport
	Trade     TradeReport
}

func NewPerformanceReport(trades []broker.Trade, equity broker.EquitySeries) PerformanceReport {
	return PerformanceReport{}
}

func PrintPerformanceReportSummary(report PerformanceReport) {

}

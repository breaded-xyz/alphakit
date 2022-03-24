package perf

import "github.com/colngroup/zero2algo/broker"

type Report struct {
}

func NewReport(trades []broker.Trade, curve []broker.Equity) Report {
	return Report{}
}

func PrintReportSummary(report Report) {

}

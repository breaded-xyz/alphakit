package perf

import "github.com/colngroup/zero2algo/broker"

type Report struct {
}

func NewReport(trades []broker.Trade, equity broker.EquitySeries) Report {
	return Report{}
}

func PrintReportSummary(report Report) {

}

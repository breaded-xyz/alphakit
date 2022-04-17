package perf

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/colngroup/zero2algo/broker"
	"github.com/olekukonko/tablewriter"
)

const _friendlyReportTimeFormat = time.RFC822

var _summaryReportHeader = []string{
	"start",
	"end",
	"trades",
	"win",
	"cagr",
	"prr",
	"mdd",
	"sharpe",
	"calmar",
}

type PerformanceReport struct {
	Description string           `json:"description"`
	Properties  map[string]any   `json:",inline"`
	Trade       *TradeReport     `json:",inline"`
	Portfolio   *PortfolioReport `json:",inline"`
}

func NewPerformanceReport(trades []broker.Trade, equity broker.EquitySeries) PerformanceReport {
	return PerformanceReport{
		Trade:     NewTradeReport(trades),
		Portfolio: NewPortfolioReport(equity),
	}
}

func PrintSummary(r PerformanceReport) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(_summaryReportHeader)
	table.Append([]string{
		r.Portfolio.PeriodStart.UTC().Format(_friendlyReportTimeFormat),
		r.Portfolio.PeriodEnd.UTC().Format(_friendlyReportTimeFormat),
		strconv.Itoa(int(r.Trade.TradeCount)),
		fmt.Sprintf("%.2f%%", r.Trade.PercentProfitable*100),
		fmt.Sprintf("%.2f%%", r.Portfolio.CAGR*100),
		fmt.Sprintf("%.2f", r.Trade.PRR),
		fmt.Sprintf("%.2f%%", r.Portfolio.MaxDrawdown*100),
		fmt.Sprintf("%.2f", r.Portfolio.Sharpe),
		fmt.Sprintf("%.2f", r.Portfolio.Calmar),
	})
	table.Render()
}

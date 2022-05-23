// Package perf provides performance metrics for a trading algo,
// using a record of trades and the equity curve from a Dealer in package broker.
package perf

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/thecolngroup/alphakit/broker"
	"github.com/thecolngroup/alphakit/internal/util"
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

// PerformanceReport is a report on the performance of a trading algo.
// It contains a TradeReport and a PortfolioReport.
//
// - TradeReport reports metrics related to the discrete trades (aka roundrtrip / roundturn).
//
//-  PorfolioReport reports metrics related to the portfolio equity curve.
type PerformanceReport struct {
	ID         string           `csv:"id"`
	Trade      *TradeReport     `csv:",inline"`
	Portfolio  *PortfolioReport `csv:",inline"`
	Properties map[string]any   `csv:"properties"`
}

// NewPerformanceReport creates a new PerformanceReport.
func NewPerformanceReport(trades []broker.Trade, equity broker.EquitySeries) PerformanceReport {
	return PerformanceReport{
		ID:         string(util.NewID()),
		Trade:      NewTradeReport(trades),
		Portfolio:  NewPortfolioReport(equity),
		Properties: make(map[string]any),
	}
}

// PrintSummary prints a summary of the performance report to stdout.
func PrintSummary(r PerformanceReport) {
	if r.Trade == nil || r.Portfolio == nil {
		println("No trades and/or equity data")
		return
	}

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

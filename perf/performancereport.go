// Package perf provides performance metrics for a trading algo,
// using a record of roundturns and the equity curve from a Dealer in package broker.
package perf

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/thecolngroup/alphakit/broker"
	"github.com/thecolngroup/alphakit/market"
	"github.com/thecolngroup/util"
)

const _friendlyReportTimeFormat = time.RFC822

var _summaryReportHeader = []string{
	"start",
	"end",
	"roundturns",
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
// - TradeReport reports metrics related to the discrete roundturns (aka roundrtrip / roundturn).
//
//-  PorfolioReport reports metrics related to the portfolio equity curve.
type PerformanceReport struct {
	ID              string           `csv:"id"`
	Asset           market.Asset     `csv:"asset_,inline"`
	TradeReport     *TradeReport     `csv:",inline"`
	PortfolioReport *PortfolioReport `csv:",inline"`
	Properties      map[string]any   `csv:"properties"`
}

// NewPerformanceReport creates a new PerformanceReport.
func NewPerformanceReport(roundturns []broker.RoundTurn, equity broker.EquitySeries) PerformanceReport {
	return PerformanceReport{
		ID:              string(util.NewID()),
		TradeReport:     NewTradeReport(roundturns),
		PortfolioReport: NewPortfolioReport(equity),
		Properties:      make(map[string]any),
	}
}

// PrintSummary prints a summary of the performance report to stdout.
func PrintSummary(r PerformanceReport) {
	if r.TradeReport == nil || r.PortfolioReport == nil {
		println("No roundturns and/or equity data")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(_summaryReportHeader)
	table.Append([]string{
		r.PortfolioReport.PeriodStart.UTC().Format(_friendlyReportTimeFormat),
		r.PortfolioReport.PeriodEnd.UTC().Format(_friendlyReportTimeFormat),
		strconv.Itoa(int(r.TradeReport.TradeCount)),
		fmt.Sprintf("%.2f%%", r.TradeReport.PercentProfitable*100),
		fmt.Sprintf("%.2f%%", r.PortfolioReport.CAGR*100),
		fmt.Sprintf("%.2f", r.TradeReport.PRR),
		fmt.Sprintf("%.2f%%", r.PortfolioReport.MaxDrawdown*100),
		fmt.Sprintf("%.2f", r.PortfolioReport.Sharpe),
		fmt.Sprintf("%.2f", r.PortfolioReport.Calmar),
	})
	table.Render()
}

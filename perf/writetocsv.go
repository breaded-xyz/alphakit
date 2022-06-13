package perf

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"github.com/thecolngroup/alphakit/broker"
	"github.com/thecolngroup/util"
)

// WritePerformanceReportToCSV writes a performance report to a CSV file.
func WritePerformanceReportToCSV(filename string, report *PerformanceReport) error {
	encMap := func(m map[string]any) ([]byte, error) {
		return []byte(fmt.Sprint(m)), nil
	}
	return util.WriteToCSV(filename, report, encMap)
}

// WriteroundturnsToCSV writes a slice of roundturns to a CSV file.
func WriteroundturnsToCSV(filename string, roundturns []broker.RoundTurn) error {
	return util.WriteToCSV(filename, roundturns)
}

type equitySeriesRow struct {
	Time   time.Time       `csv:"time"`
	Amount decimal.Decimal `csv:"amount"`
}

// WriteEquitySeriesToCSV writes an equity curve to a CSV file.
func WriteEquitySeriesToCSV(filename string, series broker.EquitySeries) error {
	rows := make([]equitySeriesRow, len(series))
	ks := series.SortKeys()
	for i := 0; i < len(ks); i++ {
		rows[i] = equitySeriesRow{
			Time:   ks[i].Time(),
			Amount: series[ks[i]],
		}
	}
	return util.WriteToCSV(filename, rows)
}

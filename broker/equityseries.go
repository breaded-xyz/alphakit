package broker

import (
	"github.com/shopspring/decimal"
)

// EquitySeries is a time series of equity values.
type EquitySeries TimeSeries[decimal.Decimal]

// SortKeys returns a sorted slice of keys in ascending chronological order.
func (es EquitySeries) SortKeys() []Timestamp {
	return TimeSeries[decimal.Decimal](es).SortKeys()
}

// SortValuesByTime returns a sorted slice of values in ascending chronological order.
func (es EquitySeries) SortValuesByTime() []decimal.Decimal {
	return TimeSeries[decimal.Decimal](es).SortValuesByTime()
}

package broker

import (
	"github.com/shopspring/decimal"
)

type EquitySeries TimeSeries[decimal.Decimal]

func (es EquitySeries) SortKeys() []Timestamp {
	return TimeSeries[decimal.Decimal](es).SortKeys()
}

func (es EquitySeries) SortValuesByTime() []decimal.Decimal {
	return TimeSeries[decimal.Decimal](es).SortValuesByTime()
}

package broker

import (
	"github.com/shopspring/decimal"
)

type Timestamp int64

type TimeSeries[V any] map[Timestamp]V

type EquitySeries TimeSeries[decimal.Decimal]

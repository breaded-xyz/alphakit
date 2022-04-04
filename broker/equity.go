package broker

import (
	"time"

	"github.com/shopspring/decimal"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type Timestamp int64

func (t Timestamp) Time() time.Time {
	return time.Unix(int64(t), 0)
}

type TimeSeries[V any] map[Timestamp]V

func (ts TimeSeries[V]) SortKeys() []Timestamp {
	ks := maps.Keys(ts)
	slices.Sort(ks)
	return ks
}

func (ts TimeSeries[V]) SortValuesByTime() []V {
	ks := ts.SortKeys()
	sorted := make([]V, len(ks))
	for i := range ks {
		sorted[i] = ts[ks[i]]
	}
	return sorted
}

type EquitySeries TimeSeries[decimal.Decimal]

func (es EquitySeries) SortKeys() []Timestamp {
	return TimeSeries[decimal.Decimal](es).SortKeys()
}

func (es EquitySeries) SortValuesByTime() []decimal.Decimal {
	return TimeSeries[decimal.Decimal](es).SortValuesByTime()
}

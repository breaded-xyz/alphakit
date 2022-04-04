package broker

import (
	"time"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type Timestamp int64

func (t Timestamp) Time() time.Time {
	return time.UnixMilli(int64(t))
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

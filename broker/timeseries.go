package broker

import (
	"time"

	"golang.org/x/exp/constraints"
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
	return SortedMapValues(ts)
}

func SortedMapValues[K constraints.Ordered, V any](m map[K]V) []V {
	sorted := make([]V, len(m))
	ks := maps.Keys(m)
	slices.Sort(ks)
	for i := range ks {
		sorted[i] = m[ks[i]]
	}
	return sorted
}

package broker

import (
	"time"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// Timestamp represents a unix timestamp in milliseconds.
type Timestamp int64

// Time returns the time.Time representation of the Timestamp.
func (t Timestamp) Time() time.Time {
	return time.UnixMilli(int64(t))
}

// TimeSeries represents a time series of values.
type TimeSeries[V any] map[Timestamp]V

// SortKeys sorts the keys of the time series by time in ascending order.
func (ts TimeSeries[V]) SortKeys() []Timestamp {
	ks := maps.Keys(ts)
	slices.Sort(ks)
	return ks
}

// SortValuesByTime sorts the values of the time series by time in ascending order.
func (ts TimeSeries[V]) SortValuesByTime() []V {
	return SortMapValues(ts)
}

// SortMapValues sorts the values of the map.
func SortMapValues[K constraints.Ordered, V any](m map[K]V) []V {
	sorted := make([]V, len(m))
	ks := maps.Keys(m)
	slices.Sort(ks)
	for i := range ks {
		sorted[i] = m[ks[i]]
	}
	return sorted
}

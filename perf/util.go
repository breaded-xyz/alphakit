package perf

import (
	"time"

	"github.com/thecolngroup/alphakit/broker"
)

// DiffReturns converts an equity curve of absolute amounts
// into a series of percentage differences.
func DiffReturns(curve broker.EquitySeries) []float64 {
	diffs := make([]float64, len(curve)-1)
	vs := curve.SortValuesByTime()
	for i := range vs {
		if i == 0 {
			continue
		}
		diffs[i-1] = (vs[i].Sub(vs[i-1]).Div(vs[i-1])).InexactFloat64()
	}
	return diffs
}

// ReduceEOD filters the equity curve to the end of day values.
// End of day is defined as equity point with hour and minute 0.
func ReduceEOD(curve broker.EquitySeries) broker.EquitySeries {
	reduced := make(broker.EquitySeries, 0)
	eodH, eodM := 0, 0 // End of day = midnight
	for k, v := range curve {
		h, m, _ := time.UnixMilli(int64(k)).Clock()
		if h == eodH && m == eodM {
			reduced[k] = v
		}
	}
	return reduced
}

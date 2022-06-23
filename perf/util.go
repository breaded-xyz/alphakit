// Copyright 2022 The Coln Group Ltd
// SPDX-License-Identifier: MIT

package perf

import (
	"time"

	"github.com/thecolngroup/alphakit/broker"
	"github.com/thecolngroup/gou/num"
)

// SimpleReturns converts an equity curve of absolute amounts into a series of percentage differences.
func SimpleReturns(curve broker.EquitySeries) []float64 {
	diffs := make([]float64, len(curve)-1)
	vs := curve.SortValuesByTime()
	for i := range vs {
		if i == 0 {
			continue
		}
		t0, t1 := vs[i].InexactFloat64(), num.NNZ(vs[i-1].InexactFloat64(), 1)
		diffs[i-1] = (t0 - t1) / t1
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

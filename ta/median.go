// Copyright 2022 The Coln Group Ltd
// SPDX-License-Identifier: MIT

package ta

import (
	"sort"

	"gonum.org/v1/gonum/stat"
)

// Median returns the median of the given values.
func Median(v []float64) float64 {
	x := make([]float64, len(v))
	copy(x, v)
	sort.Float64s(x)
	return stat.Quantile(0.5, stat.Empirical, x, nil)
}

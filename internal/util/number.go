package util

import (
	"math"
	"sort"

	"github.com/gonum/stat"
	"golang.org/x/exp/constraints"
)

func Round2DP(x float64) float64 {
	return math.Round(x*100) / 100
}

func RoundTo(x, y float64) float64 {
	return y * math.Round(x/y)
}

func Median(v []float64) float64 {
	x := make([]float64, len(v))
	copy(x, v)
	sort.Float64s(x)
	return stat.Quantile(0.5, stat.Empirical, x, nil)
}

// NN (Not Number) returns y if x is NaN or Inf.
func NN(x, y float64) float64 {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return y
	}
	return x
}

// NNZ (Not Number or Zero) returns y if x is NaN or Inf or Zero.
func NNZ(x, y float64) float64 {
	if NN(x, y) == y || x == 0 {
		return y
	}
	return x
}

// Between returns true if x >= lower and x <= upper.
func Between[T constraints.Ordered](v, lower, upper T) bool {
	return v >= lower && v <= upper
}

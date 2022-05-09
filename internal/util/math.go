package util

import "math"

func Round2DP(x float64) float64 {
	return math.Round(x*100) / 100
}

func RoundTo(x, y float64) float64 {
	return y * math.Round(x/y)
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

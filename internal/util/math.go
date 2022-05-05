package util

import "math"

func Round2DP(x float64) float64 {
	return math.Round(x*100) / 100
}

func RoundTo(x, y float64) float64 {
	return y * math.Round(x/y)
}

package perf

import (
	"github.com/gonum/floats"
)

func OptimalF(trades []float64) float64 {

	maxLoss := floats.Min(trades)
	var maxTWR, optimalF float64

	for i := 1.0; i <= 100.0; i++ {
		twr := 1.0
		f := i / 100
		for j := range trades {
			if trades[j] == 0 {
				continue
			}
			hpr := 1 + f*(-trades[j]/maxLoss)
			twr *= hpr
		}
		if twr > maxTWR {
			maxTWR = twr
			optimalF = f
		}
	}

	return optimalF
}

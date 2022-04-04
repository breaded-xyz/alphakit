package perf

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSharpeRatio(t *testing.T) {
	give := []float64{0.1, 0.2, -0.15, 0.1, 0.8, -0.3, 0.2}
	exp := 6.20
	act := SharpeRatio(give, SharpeAnnualRiskFreeRate)
	assert.Equal(t, exp, round2DP(act))
}

func round2DP(x float64) float64 {
	return math.Round(x*100) / 100
}

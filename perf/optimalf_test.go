package perf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptimalF(t *testing.T) {
	roundturns := []float64{10, 20, 50, -10, 40, -40}
	f := OptimalF(roundturns)
	assert.EqualValues(t, 0.45, f)
}

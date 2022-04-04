package perf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptimalF(t *testing.T) {
	trades := []float64{10, 20, 50, -10, 40, -40}
	f := OptimalF(trades)
	assert.EqualValues(t, 0.45, f)
}

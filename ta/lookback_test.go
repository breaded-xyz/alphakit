package ta

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLookback(t *testing.T) {
	arr := []float64{0, 1, 2, 3, 4, 5}
	assert.EqualValues(t, 5, Lookback(arr, 0))
	assert.EqualValues(t, 0, Lookback(arr, 5))
	assert.EqualValues(t, 4, Lookback(arr, 1))

	assert.EqualValues(t, 0, Lookback(arr, 6))
}

func TestLookbackWindow(t *testing.T) {
	arr := []float64{0, 1, 2, 3, 4, 5}
	assert.EqualValues(t, []float64{5}, Window(arr, 0))
	assert.EqualValues(t, []float64{4, 5}, Window(arr, 1))
	assert.EqualValues(t, []float64{0, 1, 2, 3, 4, 5}, Window(arr, 5))

	// Return max available window
	assert.EqualValues(t, []float64{0, 1, 2, 3, 4, 5}, Window(arr, 10))
}

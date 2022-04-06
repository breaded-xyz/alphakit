package ta

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrossUp(t *testing.T) {
	// Case -ve to zero (no cross)
	assert.False(t, CrossUp([]float64{-1, 0}, 0))

	// Case zero to zero (no cross)
	assert.False(t, CrossUp([]float64{0, 0}, 0))

	// Case +ve to -ve (cross down)
	assert.False(t, CrossUp([]float64{10, -1}, 0))

	// Case -ve to +ve (cross up)
	assert.True(t, CrossUp([]float64{-1, 1}, 0))

	// Case x non-Zero
	assert.True(t, CrossUp([]float64{8, 10}, 9))
}

func TestCrossDown(t *testing.T) {
	// Case +ve to zero (no cross)
	assert.False(t, CrossDown([]float64{1, 0}, 0))

	// Case zero to zero (no cross)
	assert.False(t, CrossDown([]float64{0, 0}, 0))

	// Case -ve to zero (cross up)
	assert.False(t, CrossDown([]float64{-10, 1}, 0))

	// Case +ve to -ve (cross down)
	assert.True(t, CrossDown([]float64{1, -1}, 0))

	// Case x non-Zero
	assert.True(t, CrossDown([]float64{10, 8}, 9))
}

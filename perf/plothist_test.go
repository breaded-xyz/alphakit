package perf

import (
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlotHistStdNormDist(t *testing.T) {

	// Draw some random values from the standard normal distribution.
	rand.Seed(int64(0))
	series := make([]float64, 10000)
	for i := range series {
		series[i] = rand.NormFloat64()
	}

	writer, err := os.Create("./testdata/out/hist-stdnorm.png")
	assert.NoError(t, err)
	byteCount, err := PlotHistStdNormDist(series, writer, 100, "std normal hist", "png")
	assert.NoError(t, err)
	assert.NotZero(t, byteCount)
	assert.NoError(t, writer.Close())
}

func TestPlotHist(t *testing.T) {

	// Draw some random values from the standard normal distribution.
	rand.Seed(int64(0))
	series := make([]float64, 10000)
	for i := range series {
		series[i] = rand.NormFloat64()
	}

	writer, err := os.Create("./testdata/out/hist.png")
	assert.NoError(t, err)
	byteCount, err := PlotHist(series, writer, 50, "hist", "png")
	assert.NoError(t, err)
	assert.NotZero(t, byteCount)
	assert.NoError(t, writer.Close())
}

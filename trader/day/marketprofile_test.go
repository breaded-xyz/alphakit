package day

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestMarketProfile(t *testing.T) {

	givePrices := []float64{10.1, 10.3, 11, 12.1, 3.2, 15}
	giveVolumes := []float64{10, 8, 22, 19, 20, 5}
	giveBins := 10

	wantHist := []float64{20, 0, 0, 0, 0, 18, 41, 0, 0, 5}

	act := NewMarketProfile(giveBins, givePrices, giveVolumes)

	assert.Equal(t, wantHist, act.Hist)

	spew.Dump(act)
}

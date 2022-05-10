package day

import (
	"encoding/csv"
	"os"
	"path"
	"testing"

	"github.com/colngroup/zero2algo/internal/util"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/ta"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

const testdataPath string = "../../internal/testdata/"

func TestNewMarketProfile(t *testing.T) {

	givePrices := []float64{10.1, 10.3, 11, 12.1, 3.2, 15}
	giveVolumes := []float64{10, 8, 22, 19, 20, 5}
	giveBins := 10

	wantHist := []float64{20, 0, 0, 0, 0, 18, 41, 0, 0, 5}

	act := NewVolumeProfile(giveBins, givePrices, giveVolumes)

	assert.Equal(t, wantHist, act.Hist)

	spew.Dump(act)
}

func TestMarketProfileWithPriceFile(t *testing.T) {

	file, _ := os.Open(path.Join(testdataPath, "btcusdt-1m-2022-05-06.csv"))
	defer func() {
		assert.NoError(t, file.Close())
	}()

	prices, err := market.NewCSVKlineReader(csv.NewReader(file)).ReadAll()
	assert.NoError(t, err)
	priceLevels := make([]float64, len(prices))
	vols := make([]float64, len(prices))
	for i := range prices {
		priceLevels[i] = util.RoundTo(ta.OHLC4(prices[i]), 0.1)
		vols[i] = util.RoundTo(prices[i].Volume, 1.0)
	}

	spew.Dump(prices[0].Start, prices[len(prices)-1].Start)
	mp := NewVolumeProfile(1440, priceLevels, vols)

	spew.Dump(mp.POC, mp.VAH, mp.VAL)
}

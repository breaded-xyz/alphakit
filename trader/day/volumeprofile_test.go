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
	"golang.org/x/exp/slices"
)

const testdataPath string = "../../internal/testdata/"

func TestMarketProfileWithPriceFile(t *testing.T) {

	file, _ := os.Open(path.Join(testdataPath, "btcusdt-1m-2022-05-06.csv"))
	defer func() {
		assert.NoError(t, file.Close())
	}()

	prices, err := market.NewCSVKlineReader(csv.NewReader(file)).ReadAll()
	assert.NoError(t, err)
	var levels []VolumeLevel
	for i := range prices {
		levels = append(levels, VolumeLevel{
			Price:  util.RoundTo(ta.HLC3(prices[i]), 1.0),
			Volume: util.RoundTo(prices[i].Volume, 1.0),
		})
	}
	slices.SortStableFunc(levels, func(i, j VolumeLevel) bool {
		return i.Price < j.Price
	})

	spew.Dump(prices[0].Start, prices[len(prices)-1].Start)
	vp := NewVolumeProfile(100, levels)

	assert.Equal(t, 35359.0, vp.Low)
	assert.Equal(t, 35753.232323232325, vp.VAL)
	assert.Equal(t, 36033.0101010101, vp.POC)
	assert.Equal(t, 36325.50505050505, vp.VAH)
	assert.Equal(t, 36617.0, vp.High)
}

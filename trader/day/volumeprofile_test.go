package day

import (
	"encoding/csv"
	"os"
	"path"
	"testing"

	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/ta"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

const testdataPath string = "./testdata/"

func TestMarketProfileWithPriceFile(t *testing.T) {

	file, _ := os.Open(path.Join(testdataPath, "BTCUSDT-1m-2022-04-08.csv"))
	defer func() {
		assert.NoError(t, file.Close())
	}()

	prices, err := market.NewCSVKlineReader(csv.NewReader(file)).ReadAll()
	assert.NoError(t, err)
	var levels []VolumeLevel
	for i := range prices {

		hlc3 := ta.HLC3(prices[i])
		vol := prices[i].Volume

		if hlc3 == 0 || vol == 0 {
			continue
		}

		levels = append(levels, VolumeLevel{
			Price:  hlc3,
			Volume: vol,
		})
	}
	slices.SortStableFunc(levels, func(i, j VolumeLevel) bool {
		return i.Price < j.Price
	})

	spew.Dump(prices[0].Start, prices[len(prices)-1].Start)
	vp := NewVolumeProfile(10, levels)

	spew.Dump(vp)

	//assert.Equal(t, 35359.0, vp.Low)
	//assert.Equal(t, 35753.232323232325, vp.VAL)
	//assert.Equal(t, 36033.0101010101, vp.POC)
	//assert.Equal(t, 36325.50505050505, vp.VAH)
	//assert.Equal(t, 36617.0, vp.High)
}

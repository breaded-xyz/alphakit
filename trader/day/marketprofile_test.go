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
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

const testdataPath string = "../../internal/testdata/"

func TestNewMarketProfile(t *testing.T) {

	givePrices := []float64{10.1, 10.3, 11, 12.1, 3.2, 15}
	giveVolumes := []float64{10, 8, 22, 19, 20, 5}
	giveBins := 10

	wantHist := []float64{20, 0, 0, 0, 0, 18, 41, 0, 0, 5}

	act := NewMarketProfile(giveBins, givePrices, giveVolumes)

	assert.Equal(t, wantHist, act.Hist)

	spew.Dump(act)
}

func TestMarketProfileWithPriceFile(t *testing.T) {

	file, _ := os.Open(path.Join(testdataPath, "btcusdt-1m-2022-05-04.csv"))
	defer func() {
		assert.NoError(t, file.Close())
	}()

	prices, err := market.NewCSVKlineReader(csv.NewReader(file)).ReadAll()
	assert.NoError(t, err)

	/*dayStart := time.Date(2022, time.April, 30, 0, 0, 0, 0, time.UTC)
	var dayStartIdx int
	for i := range prices {
		if prices[i].Start.After(dayStart) {
			dayStartIdx = i - 1
			break
		}
	}

	//prices = prices[dayStartIdx:]*/
	priceLevels := make([]float64, len(prices))
	vols := make([]float64, len(prices))
	for i := range prices {
		priceLevels[i] = util.RoundTo(ta.OHLC4(prices[i]), 1.0)
		vols[i] = util.RoundTo(prices[i].Volume, 1.0)
	}

	//spew.Dump(hlc3s, vols)

	spew.Dump(prices[0].Start, prices[len(prices)-1].Start)
	mp := NewMarketProfile(100, priceLevels, vols)

	spew.Dump(mp.POC, mp.VAH, mp.VAL, mp.High, mp.Low)

	// Make a plot and set its title.
	p := plot.New()

	p.Title.Text = "Histogram"

	// Create a histogram of our values drawn
	// from the standard normal.

	//pts := make([]plotter.XYer, len(mp.Hist))
	//for i := range pts {
	//	xys := make(plotter.XYs, len(mp.Hist))
	//	pts[i] = plotter.XYValues
	//

	var xys plotter.XYs
	for i := range mp.Hist {
		xys = append(xys, plotter.XY{X: mp.Bins[i], Y: mp.Hist[i]})
	}

	h, err := plotter.NewHistogram(xys, len(mp.Bins))
	if err != nil {
		panic(err)
	}

	p.Add(h)

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "hist.png"); err != nil {
		panic(err)
	}
}

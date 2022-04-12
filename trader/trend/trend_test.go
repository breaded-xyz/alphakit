package trend

import (
	"context"
	"encoding/csv"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/colngroup/zero2algo/broker/backtest"
	"github.com/colngroup/zero2algo/dec"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/money"
	"github.com/colngroup/zero2algo/perf"
	"github.com/colngroup/zero2algo/ta"
	"github.com/davecgh/go-spew/spew"
)

const testdataPath string = "../../testdata/"

func TestTrendBot(t *testing.T) {
	asset := market.NewAsset("BTCUSDT")

	dealer := backtest.NewDealer()
	dealer.SetInitialCapital(dec.New(1000))

	predicter := *NewPredicter(
		ta.NewOsc(ta.NewALMA(1), ta.NewALMA(256)),
		ta.NewSDWithFactor(512, 1.5),
		ta.NewMMIWithSmoother(200, ta.NewALMA(200)))

	bot := Bot{
		EnterLong:  1,
		ExitLong:   -0.9,
		EnterShort: -1,
		ExitShort:  0.6,
		asset:      asset,
		dealer:     dealer,
		predicter:  predicter,
		//risker:     NewSDRisk(512, 1.5),
		//sizer: &money.SafeFSizer{
		//	InitialCapital: dec.New(1000),
		//	F:              0.5,
		//	ScaleF:         0.5,
		//},
		risker: NewMaxRisk(),
		sizer:  &money.FixedSizer{FixedCapital: dec.New(1000)},
	}

	filepath.WalkDir(path.Join(testdataPath, "btcusdt-1h-2021-Q1.csv"),
		func(path string, d fs.DirEntry, err error) error {

			file, _ := os.Open(path)
			defer file.Close()
			prices, _ := market.NewCSVKlineReader(csv.NewReader(file)).ReadAll()

			for _, price := range prices {
				if err := dealer.ReceivePrice(context.Background(), price); err != nil {
					t.Fatal(err)
				}
				if err := bot.ReceivePrice(context.Background(), price); err != nil {
					t.Fatal(err)
				}
			}

			return nil
		})
	bot.Close(context.Background())

	trades, _, _ := dealer.ListTrades(context.Background(), nil)
	equity := dealer.EquityHistory()
	report := perf.NewPerformanceReport(trades, equity)

	spew.Dump(report)
}

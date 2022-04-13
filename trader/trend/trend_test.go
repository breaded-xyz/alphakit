package trend

import (
	"context"
	"encoding/csv"
	"os"
	"path"
	"testing"

	"github.com/colngroup/zero2algo/broker/backtest"
	"github.com/colngroup/zero2algo/dec"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/money"
	"github.com/colngroup/zero2algo/perf"
	"github.com/colngroup/zero2algo/risk"
	"github.com/colngroup/zero2algo/ta"
)

const testdataPath string = "../../testdata/"

func TestTrendBot(t *testing.T) {
	dealer := backtest.NewDealer()
	dealer.SetInitialCapital(dec.New(1000))

	predicter := NewPredicter(
		ta.NewOsc(ta.NewALMA(1), ta.NewALMA(256)),
		ta.NewSDWithFactor(512, 1.5),
		ta.NewMMIWithSmoother(200, ta.NewALMA(200)))

	bot := Bot{
		EnterLong:  1,
		ExitLong:   -0.9,
		EnterShort: -1,
		ExitShort:  0.6,
		asset:      market.NewAsset("BTCUSDT"),
		dealer:     dealer,
		Predicter:  predicter,
		Risker:     risk.NewFullRisk(),
		Sizer:      &money.FixedSizer{FixedCapital: dec.New(1000)},
	}

	file, _ := os.Open(path.Join(testdataPath, "btcusdt-1h-2021-Q1.csv"))
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

	bot.Close(context.Background())

	trades, _, _ := dealer.ListTrades(context.Background(), nil)
	equity := dealer.EquityHistory()
	report := perf.NewPerformanceReport(trades, equity)
	perf.PrintSummary(report)
}

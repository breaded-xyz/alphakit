package trend

import (
	"context"
	"encoding/csv"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thecolngroup/alphakit/broker/backtest"
	"github.com/thecolngroup/alphakit/internal/dec"
	"github.com/thecolngroup/alphakit/market"
	"github.com/thecolngroup/alphakit/money"
	"github.com/thecolngroup/alphakit/perf"
	"github.com/thecolngroup/alphakit/risk"
	"github.com/thecolngroup/alphakit/ta"
)

const testdataPath string = "../../internal/testdata/"

func TestBotWithCrossPredicter(t *testing.T) {
	dealer := backtest.NewDealer()
	dealer.SetInitialCapital(dec.New(1000))

	predicter := NewCrossPredicter(
		ta.NewOsc(ta.NewALMA(32), ta.NewALMA(64)),
		ta.NewMMIWithSmoother(200, ta.NewALMA(200)))

	bot := Bot{
		EnterLong:  1,
		ExitLong:   -0.9,
		EnterShort: -1,
		ExitShort:  0.9,
		Asset:      market.NewAsset("BTCUSDT"),
		dealer:     dealer,
		Predicter:  predicter,
		Risker:     risk.NewFullRisker(),
		Sizer:      money.NewFixedSizer(dec.New(1000)),
	}

	file, _ := os.Open(path.Join(testdataPath, "btcusdt-1h-2021-Q1.csv"))
	defer func() {
		assert.NoError(t, file.Close())
	}()

	prices, _ := market.NewCSVKlineReader(csv.NewReader(file)).ReadAll()
	for _, price := range prices {
		if err := dealer.ReceivePrice(context.Background(), price); err != nil {
			t.Fatal(err)
		}
		if err := bot.ReceivePrice(context.Background(), price); err != nil {
			t.Fatal(err)
		}
	}

	assert.NoError(t, bot.Close(context.Background()))

	trades, _, _ := dealer.ListTrades(context.Background(), nil)
	equity := dealer.EquityHistory()
	report := perf.NewPerformanceReport(trades, equity)
	perf.PrintSummary(report)
}

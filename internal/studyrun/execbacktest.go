package studyrun

import (
	"context"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/perf"
	"github.com/colngroup/zero2algo/trader"
)

func execBacktest(bot trader.Bot, dealer broker.SimulatedDealer, prices []market.Kline) (perf.PerformanceReport, error) {
	var empty perf.PerformanceReport

	for _, price := range prices {
		if err := dealer.ReceivePrice(context.Background(), price); err != nil {
			return empty, err
		}
		if err := bot.ReceivePrice(context.Background(), price); err != nil {
			return empty, err
		}
	}
	bot.Close(context.Background())
	trades, _, _ := dealer.ListTrades(context.Background(), nil)
	equity := dealer.EquityHistory()
	report := perf.NewPerformanceReport(trades, equity)

	return report, nil
}

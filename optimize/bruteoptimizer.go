package optimize

import (
	"context"
	"errors"
	"math"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/perf"
	"github.com/colngroup/zero2algo/trader"
)

type CaseResult struct {
	CAGR   float64
	PRR    float64
	Sharpe float64
	Calmar float64
}

type Step struct {
	Result perf.PerformanceReport
}

type BruteOptimizer struct {
	Params         map[string]any
	Samples        [][]market.Kline
	SampleSplitPct float64
	Warmup         int
	MakeBot        func() trader.ConfigurableBot
	MakeDealer     func() broker.SimulatedDealer

	testCases map[string]map[string]any
	results   map[string]CaseResult
}

func (o *BruteOptimizer) Prepare() {}

func (o *BruteOptimizer) Start(ctx context.Context) (chan<- Step, error) {

	resultCh := make(chan Step)

	// Run training phase
	// Select top rank params based on average metrics
	// Run OOS backtest

	for k := range o.testCases {
		testCase := o.testCases[k]
		sampleReports := make([]perf.PerformanceReport, 0, len(o.Samples))
	testCase:
		for i := range o.Samples {
			sample := o.Samples[i]

			dealer := o.MakeDealer()
			bot := o.MakeBot()
			if err := bot.Configure(testCase); err != nil {
				if errors.Is(err, trader.ErrInvalidConfig) {
					continue testCase
				}
				panic(err)
			}
			for i := range sample {
				price := sample[i]
				if err := dealer.ReceivePrice(ctx, price); err != nil {
					panic(err)
				}
				if err := bot.ReceivePrice(ctx, price); err != nil {
					panic(err)
				}
			}
			bot.Close(context.Background())
			trades, _, _ := dealer.ListTrades(context.Background(), nil)
			equity := dealer.EquityHistory()
			result := perf.NewPerformanceReport(trades, equity)
			sampleReports = append(sampleReports, result)
			resultCh <- Step{
				Result: result,
			}
		}
		// Create testCase Report
		o.results[k] = newCaseResult(sampleReports)
	}

	return nil, nil
}

func newCaseResult(reports []perf.PerformanceReport) CaseResult {
	return CaseResult{}
}

func splitSample(sample []market.Kline, splitPct float64) (a, b []market.Kline) {
	splitIndex := float64(len(sample)) * splitPct
	splitIndex = math.Floor(splitIndex)
	return sample[:int(splitIndex)], sample[int(splitIndex):]
}

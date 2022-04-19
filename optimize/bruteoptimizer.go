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

type ParamRange map[string]any

type OptimizerStep struct {
	Params ParamSet
	Result ParamSetReport
}

type BruteOptimizer struct {
	SampleSplitPct float64
	WarmupBarCount int
	MakeBot        func() trader.ConfigurableBot
	MakeDealer     func() broker.SimulatedDealer

	study *Study
}

func (o *BruteOptimizer) Prepare(in ParamRange, samples [][]market.Kline) (int, error) {

	products := CartesianBuilder(in)
	for i := range products {
		pSet := NewParamSet()
		pSet.Params = products[i]
		o.study.ParamSets[pSet.ID] = pSet
	}

	o.study.Samples = samples

	steps := len(o.study.ParamSets) * len(samples) // Training phase
	steps += len(samples)                          // Validation phase for single optima

	return steps, nil
}

func (o *BruteOptimizer) Start(ctx context.Context) (chan<- OptimizerStep, error) {

	resultCh := make(chan OptimizerStep)

	// Run training phase
	// Select top rank params based on average metrics
	// Run OOS backtest

	for k := range o.study.ParamSets {
		pSet := o.study.ParamSets[k]
		perSampleReports := make([]perf.PerformanceReport, 0, len(o.study.Samples))
	pSetBacktest:
		for i := range o.study.Samples {
			sample := o.study.Samples[i]

			dealer := o.MakeDealer()
			bot := o.MakeBot()
			if err := bot.Configure(pSet.Params); err != nil {
				if errors.Is(err, trader.ErrInvalidConfig) {
					continue pSetBacktest
				}
				panic(err)
			}

			if err := bot.Warmup(ctx, sample[:o.WarmupBarCount]); err != nil {
				panic(err)
			}

			for i := range sample[o.WarmupBarCount:] {
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
			perSampleReports = append(perSampleReports, result)
			resultCh <- OptimizerStep{}
		}
		// Create testCase Report
		o.study.InSample[pSet.ID] = newParamSetReport(perSampleReports)
	}

	return nil, nil
}

func newParamSetReport(reports []perf.PerformanceReport) ParamSetReport {
	return ParamSetReport{}
}

func splitSample(sample []market.Kline, splitPct float64) (a, b []market.Kline) {
	splitIndex := float64(len(sample)) * splitPct
	splitIndex = math.Floor(splitIndex)
	return sample[:int(splitIndex)], sample[int(splitIndex):]
}

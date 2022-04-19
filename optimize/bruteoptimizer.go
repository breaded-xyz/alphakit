package optimize

import (
	"context"
	"errors"
	"math"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/perf"
	"github.com/colngroup/zero2algo/trader"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type Phase int

const (
	Training Phase = iota
	Validation
)

type ParamRange map[string]any

type OptimizerStep struct {
	Params ParamSet
	Result ParamSetReport
	Err    error
}

type BruteOptimizer struct {
	SampleSplitPct float64
	WarmupBarCount int
	MakeBot        func() trader.ConfigurableBot
	MakeDealer     func() broker.SimulatedDealer
	RankFunc       ObjectiveRanker

	study *Study
}

func (o *BruteOptimizer) Prepare(in ParamRange, samples [][]market.Kline) (int, error) {

	products := CartesianBuilder(in)
	for i := range products {
		pSet := NewParamSet()
		pSet.Params = products[i]
		o.study.Training[pSet.ID] = pSet
	}

	o.study.Samples = samples

	steps := len(o.study.Training) * len(samples) // Training phase
	steps += len(samples)                         // Validation phase for optimum

	return steps, nil
}

func (o *BruteOptimizer) Start(ctx context.Context) (chan<- OptimizerStep, error) {

	resultCh := make(chan OptimizerStep)

	go func() {
		trainingCh, err := o.startPhase(ctx, Training)
		if err != nil {
			resultCh <- OptimizerStep{Err: err}
		}

		select {
		case <-trainingCh:
		}

		results := maps.Values(o.study.TrainingResults)
		slices.SortFunc(results, o.RankFunc)

		validationCh, err := o.startPhase(ctx, Validation)
		if err != nil {
			resultCh <- OptimizerStep{Err: err}
		}

		select {
		case <-validationCh:
		}
	}()

	return resultCh, nil
}

func (o *BruteOptimizer) startPhase(ctx context.Context, phase Phase) (chan OptimizerStep, error) {

	resultCh := make(chan OptimizerStep)

	var phasePSets map[ParamSetID]ParamSet

	switch phase {
	case Training:
		phasePSets = o.study.Training
	case Validation:
		phasePSets = o.study.Optima
	default:
		return nil, errors.New("invalid phase")
	}

	for k := range phasePSets {
		pSet := o.study.Training[k]
		perSampleReports := make([]perf.PerformanceReport, 0, len(o.study.Samples))
	pSetBacktest:
		for i := range o.study.Samples {
			sample := splitSample(phase, o.study.Samples[i], o.SampleSplitPct)

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

		o.study.TrainingResults[pSet.ID] = newParamSetReport(perSampleReports)
	}

	return resultCh, nil
}

func newParamSetReport(reports []perf.PerformanceReport) ParamSetReport {
	return ParamSetReport{}
}

func splitSample(phase Phase, sample []market.Kline, splitPct float64) []market.Kline {

	splitIndex := float64(len(sample)) * splitPct
	splitIndex = math.Floor(splitIndex)

	var split []market.Kline

	switch phase {
	case Training:
		split = sample[:int(splitIndex)]
	case Validation:
		split = sample[int(splitIndex):]
	}

	return split
}

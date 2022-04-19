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

// BruteOptimizer implements a 'brute-force peak objective' optimization approach which
// tests all given parameter combinations and selects the highest ranked (peak) param set.
// Prepare:
// - Accept 1 or more price data samples (inc. over/under sampling)
// - Split sample into in-sample (training) and out-of-sample (validation) datasets
// - Generate 1 or more param sets using the cartesian product of given ranges that define the param space
// Train:
// - Execute each algo param set over the in-sample price data
// - Average the performance for each param set over the in-sample data
// - Rank the param sets based on the performance objective (Profit Factor, T-Score etc)
// Validate:
// - Execute the highest ranked ("trained") algo param over the out-of-sample price data
// - Accept or reject the hypothesis based on statistical significance of the performance report
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

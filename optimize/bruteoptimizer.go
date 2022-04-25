package optimize

import (
	"context"
	"errors"
	"math"
	"sync"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/perf"
	"github.com/colngroup/zero2algo/trader"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type Phase int

const (
	Training Phase = iota + 1
	Validation
)

type ParamRange map[string]any

type OptimizerStep struct {
	Phase    Phase
	ParamSet ParamSet
	Result   perf.PerformanceReport
	Err      error
}

// BruteOptimizer implements a 'brute-force peak objective' optimization approach which
// tests all given parameter combinations and selects the highest ranked (peak) param set.
// Optima is selected by the given ObjectiveRanker func.
// Optimization method in 3 stages:
// Prepare:
// - Accept 1 or more price data samples
// - Split sample price data into in-sample (training) and out-of-sample (validation) datasets
// - Generate 1 or more param sets using the cartesian product of given ranges that define the param space
// Train:
// - Execute each algo param set over the in-sample price data
// - Average the performance for each param set over the in-sample data
// - Rank the param sets based on the performance objective (Profit Factor, T-Score etc)
// Validate:
// - Execute the highest ranked ("trained") algo param over the out-of-sample price data
// - Accept or reject the hypothesis based on statistical significance of the study report
type BruteOptimizer struct {
	SampleSplitPct float64
	WarmupBarCount int
	MakeBot        trader.MakeBot
	MakeDealer     broker.MakeSimulatedDealer
	Ranker         ObjectiveRanker

	study Study
}

type bruteOptimizerJob struct {
	ParamSet       ParamSet
	Sample         []market.Kline
	WarmupBarCount int
	MakeBot        trader.MakeBot
	MakeDealer     broker.MakeSimulatedDealer
}

func (o *BruteOptimizer) Prepare(in ParamRange, samples [][]market.Kline) (int, error) {

	products := CartesianBuilder(in)
	for i := range products {
		pSet := NewParamSet()
		pSet.Params = products[i]
		o.study.TrainingPSets = append(o.study.TrainingPSets, pSet)
	}

	for i := range samples {
		training, validation := splitSample(samples[i], o.SampleSplitPct)
		o.study.TrainingSamples = append(o.study.TrainingSamples, training)
		o.study.OptimaSamples = append(o.study.OptimaSamples, validation)
	}

	steps := len(o.study.TrainingPSets) * len(samples) // Training phase
	steps += len(samples)                              // Validation phase for optimum

	return steps, nil
}

func (o *BruteOptimizer) Start(ctx context.Context) (<-chan OptimizerStep, error) {

	outCh := make(chan OptimizerStep)

	go func() {
		defer close(outCh)

		doneCh := make(chan struct{})
		defer close(doneCh)

		// Training phase
		trainigJobCh := o.enqueueJobs(o.study.TrainingPSets, o.study.TrainingSamples)
		trainingOutCh := processBruteJobs(ctx, doneCh, trainigJobCh)
		for step := range trainingOutCh {
			step.Phase = Training
			outCh <- step
			if step.Err != nil {
				if errors.Is(step.Err, trader.ErrInvalidConfig) {
					continue
				} else {
					return
				}
			}
			report := o.study.TrainingResults[step.ParamSet.ID]
			report.AddResult(step.Result)
			o.study.TrainingResults[step.ParamSet.ID] = report
		}

		// Evaluate traning results and select top ranked pset for validation phase
		slices.SortFunc(maps.Values(o.study.TrainingResults), o.Ranker)
		o.study.OptimaPSets = append(o.study.OptimaPSets, o.study.TrainingPSets[len(o.study.TrainingPSets)-1])

		// Validation phase
		validationJobCh := o.enqueueJobs(o.study.OptimaPSets, o.study.OptimaSamples)
		validationOutCh := processBruteJobs(ctx, doneCh, validationJobCh)
		for step := range validationOutCh {
			step.Phase = Validation
			outCh <- step
			if step.Err != nil {
				return
			}
			report := o.study.OptimaResults[step.ParamSet.ID]
			report.AddResult(step.Result)
			o.study.OptimaResults[step.ParamSet.ID] = report
		}
	}()

	return outCh, nil
}

func (o *BruteOptimizer) enqueueJobs(pSets []ParamSet, samples [][]market.Kline) <-chan bruteOptimizerJob {

	// A buffered channel enables us to enqueue jobs and close the channel in a single function to simplify the call flow
	// Without a buffer the loop would block awaiting a ready receiver for the jobs
	jobCh := make(chan bruteOptimizerJob, len(pSets)*len(samples))
	defer close(jobCh)

	// Enqueue a job for each pset and price series combination
	for i := range pSets {
		for j := range samples {
			jobCh <- bruteOptimizerJob{
				ParamSet:       pSets[i],
				Sample:         samples[j],
				WarmupBarCount: o.WarmupBarCount,
				MakeBot:        o.MakeBot,
				MakeDealer:     o.MakeDealer,
			}
		}
	}
	return jobCh
}

func processBruteJobs(ctx context.Context, doneCh <-chan struct{}, jobCh <-chan bruteOptimizerJob) <-chan OptimizerStep {

	outCh := make(chan OptimizerStep)

	go func() {
		defer close(outCh)
		var wg sync.WaitGroup
		next := true

		for next {
			select {
			case <-ctx.Done():
				next = false
			case <-doneCh:
				next = false
			case job, ok := <-jobCh:
				if !ok {
					next = false
					break
				}
				wg.Add(1)
				go func() {
					defer wg.Done()
					dealer, err := job.MakeDealer(job.ParamSet.Params)
					if err != nil {
						outCh <- OptimizerStep{ParamSet: job.ParamSet, Err: err}
					}
					bot, err := job.MakeBot(job.ParamSet.Params)
					if err != nil {
						outCh <- OptimizerStep{ParamSet: job.ParamSet, Err: err}
					}
					if err := bot.Warmup(ctx, job.Sample[:job.WarmupBarCount]); err != nil {
						outCh <- OptimizerStep{ParamSet: job.ParamSet, Err: err}
					}

					perf, err := runBacktest(ctx, bot, dealer, job.Sample[job.WarmupBarCount:])
					outCh <- OptimizerStep{ParamSet: job.ParamSet, Result: perf, Err: err}
				}()
			}
		}
		wg.Wait()
	}()

	return outCh
}

func runBacktest(ctx context.Context, bot trader.Bot, dealer broker.SimulatedDealer, prices []market.Kline) (perf.PerformanceReport, error) {
	var empty perf.PerformanceReport

	for i := range prices {
		price := prices[i]
		if err := dealer.ReceivePrice(ctx, price); err != nil {
			return empty, err
		}
		if err := bot.ReceivePrice(ctx, price); err != nil {
			return empty, err
		}
	}

	bot.Close(ctx)
	trades, _, err := dealer.ListTrades(context.Background(), nil)
	if err != nil {
		return empty, err
	}
	equity := dealer.EquityHistory()
	report := perf.NewPerformanceReport(trades, equity)

	return report, nil
}

func splitSample(sample []market.Kline, splitPct float64) (a, b []market.Kline) {
	splitIndex := float64(len(sample)) * splitPct
	splitIndex = math.Floor(splitIndex)
	return sample[:int(splitIndex)], sample[int(splitIndex):]
}

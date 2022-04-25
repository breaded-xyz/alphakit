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

type optimizerJob struct {
	ParamSet       ParamSet
	Sample         []market.Kline
	WarmupBarCount int
	MakeBot        func() trader.ConfigurableBot
	MakeDealer     func() broker.SimulatedDealer
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
// - Accept or reject the hypothesis based on statistical significance of the study report
type BruteOptimizer struct {
	SampleSplitPct float64
	WarmupBarCount int
	MakeBot        func() trader.ConfigurableBot
	MakeDealer     func() broker.SimulatedDealer
	RankFunc       ObjectiveRanker

	study Study
}

func (o *BruteOptimizer) Prepare(in ParamRange, samples [][]market.Kline) (int, error) {

	products := CartesianBuilder(in)
	for i := range products {
		pSet := NewParamSet()
		pSet.Params = products[i]
		o.study.Training = append(o.study.Training, pSet)
	}

	o.study.Samples = samples

	steps := len(o.study.Training) * len(samples) // Training phase
	steps += len(samples)                         // Validation phase for optimum

	return steps, nil
}

func (o *BruteOptimizer) Start(ctx context.Context) (<-chan OptimizerStep, error) {

	outChan := make(chan OptimizerStep)

	go func() {
		defer close(outChan)

		doneChan := make(chan struct{})
		defer close(doneChan)

		// Training phase
		pSets, samples := prepareTraining(o.study, o.SampleSplitPct)
		trainingOutChan := o.startPhase(ctx, doneChan, pSets, samples)
		for step := range trainingOutChan {
			outChan <- step
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
		slices.SortFunc(maps.Values(o.study.TrainingResults), o.RankFunc)
		o.study.Optima = append(o.study.Optima, o.study.Training[len(o.study.Training)-1])

		// Validation phase

	}()

	return outChan, nil
}

func (o *BruteOptimizer) startPhase(ctx context.Context, doneCh <-chan struct{}, pSets []ParamSet, samples [][]market.Kline) <-chan OptimizerStep {
	outChan := make(chan OptimizerStep)

	go func() {
		defer close(outChan)

		jobChan := make(chan optimizerJob)
		jobOutChan := processJobs(ctx, doneCh, jobChan)

		for i := range pSets {
			for j := range samples {
				jobChan <- optimizerJob{
					ParamSet:       pSets[i],
					Sample:         samples[j],
					WarmupBarCount: o.WarmupBarCount,
					MakeBot:        o.MakeBot,
					MakeDealer:     o.MakeDealer,
				}
			}
		}
		close(jobChan)

		for step := range jobOutChan {
			outChan <- step
		}
	}()

	return outChan
}

func processJobs(ctx context.Context, doneCh <-chan struct{}, jobCh <-chan optimizerJob) <-chan OptimizerStep {

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
					dealer := job.MakeDealer()
					bot := job.MakeBot()

					if err := bot.Configure(job.ParamSet.Params); err != nil {
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

func prepareTraining(study Study, splitPct float64) ([]ParamSet, [][]market.Kline) {
	var trainingSamples [][]market.Kline
	for i := range study.Samples {
		training, _ := splitSample(study.Samples[i], splitPct)
		trainingSamples = append(trainingSamples, training)
	}
	return study.Training, trainingSamples
}

func prepareValidation(study Study, splitPct float64) ([]ParamSet, [][]market.Kline) {
	var validationSamples [][]market.Kline
	for i := range study.Samples {
		_, validation := splitSample(study.Samples[i], splitPct)
		validationSamples = append(validationSamples, validation)
	}
	return study.Optima, validationSamples
}

func splitSample(sample []market.Kline, splitPct float64) (a, b []market.Kline) {
	splitIndex := float64(len(sample)) * splitPct
	splitIndex = math.Floor(splitIndex)
	return sample[:int(splitIndex)], sample[int(splitIndex):]
}

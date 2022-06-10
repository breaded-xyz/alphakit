package optimize

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thecolngroup/alphakit/broker"
	"github.com/thecolngroup/alphakit/market"
	"github.com/thecolngroup/alphakit/perf"
	"github.com/thecolngroup/alphakit/trader"
	"github.com/thecolngroup/dec"
	"golang.org/x/exp/maps"
)

func TestBruteOptimizer_Prepare(t *testing.T) {
	tests := []struct {
		name               string
		giveParamRange     ParamMap
		giveSamples        map[AssetID][]market.Kline
		giveSampleSplitPct float64
		wantSteps          int
		wantStudy          Study
	}{
		{
			name:           "ok",
			giveParamRange: map[string]any{"A": []any{1, 2}, "B": []any{10}},
			giveSamples: map[AssetID][]market.Kline{
				AssetID("asset_x"): {{C: dec.New(10)}, {C: dec.New(20)}, {C: dec.New(30)}, {C: dec.New(40)}},
				AssetID("asset_y"): {{C: dec.New(50)}, {C: dec.New(60)}, {C: dec.New(70)}},
			},
			giveSampleSplitPct: 0.5,
			wantSteps:          6,
			wantStudy: Study{
				Training: []ParamSet{
					{Params: map[string]any{"A": 1, "B": 10}},
					{Params: map[string]any{"A": 2, "B": 10}},
				},
				TrainingSamples: map[AssetID][]market.Kline{
					AssetID("asset_x"): {{C: dec.New(10)}, {C: dec.New(20)}},
					AssetID("asset_y"): {{C: dec.New(50)}, {C: dec.New(60)}},
				},
				ValidationSamples: map[AssetID][]market.Kline{
					AssetID("asset_x"): {{C: dec.New(30)}},
					AssetID("asset_y"): {{C: dec.New(60)}},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			optimizer := NewBruteOptimizer()
			optimizer.SampleSplitPct = tt.giveSampleSplitPct
			actSteps, err := optimizer.Prepare(tt.giveParamRange, tt.giveSamples)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantSteps, actSteps)
			assert.Len(t, optimizer.study.Training, len(tt.wantStudy.Training))
			assert.ElementsMatch(t, maps.Values(optimizer.study.TrainingSamples), maps.Values(tt.wantStudy.TrainingSamples))
		})
	}
}

func TestBruteOptimizer_Start(t *testing.T) {
	tests := []struct {
		name      string
		giveStudy *Study
		wantStudy *Study
	}{
		{
			name: "select top ranked pset for validation",
			giveStudy: &Study{
				Training: []ParamSet{
					{ID: "1", Params: map[string]any{"A": 1, "B": 10}},
					{ID: "2", Params: map[string]any{"A": 2, "B": 10}},
				},
				TrainingResults: map[ParamSetID]PhaseReport{
					"1": {PRR: 2, TradeCount: 2, Subject: ParamSet{ID: "1", Params: map[string]any{"A": 1, "B": 10}}},
					"2": {PRR: 4, TradeCount: 2, Subject: ParamSet{ID: "2", Params: map[string]any{"A": 2, "B": 10}}},
				},
				ValidationResults: make(map[ParamSetID]PhaseReport),
			},
			wantStudy: &Study{
				Training: []ParamSet{
					{ID: "1", Params: map[string]any{"A": 1, "B": 10}},
					{ID: "2", Params: map[string]any{"A": 2, "B": 10}},
				},
				TrainingResults: map[ParamSetID]PhaseReport{
					"1": {PRR: 2, TradeCount: 2, Subject: ParamSet{ID: "1", Params: map[string]any{"A": 1, "B": 10}}},
					"2": {PRR: 4, TradeCount: 2, Subject: ParamSet{ID: "2", Params: map[string]any{"A": 2, "B": 10}}},
				},
				Validation: []ParamSet{
					{ID: "2", Params: map[string]any{"A": 2, "B": 10}},
				},
				ValidationResults: map[ParamSetID]PhaseReport{
					"2": {},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			optimizer := BruteOptimizer{
				MakeBot:    func(config map[string]any) (trader.Bot, error) { return &trader.StubBot{}, nil },
				MakeDealer: func() (broker.SimulatedDealer, error) { return &broker.StubDealer{}, nil },
				Ranker:     PRRRanker,
				study:      tt.giveStudy,
			}

			stepCh, err := optimizer.Start(context.Background())
			assert.NoError(t, err)
			for range stepCh {
			}
			act := optimizer.Study()
			assert.Equal(t, tt.wantStudy, act)
		})
	}
}

func TestBruteOptimizer_enqueueJobs(t *testing.T) {

	givePSets := []ParamSet{
		{ID: "0", Params: map[string]any{"A": 0, "B": 1}},
		{ID: "1", Params: map[string]any{"Y": 25, "Z": 26}},
	}
	giveSamples := map[AssetID][]market.Kline{
		"asset_x": {{C: dec.New(10)}, {C: dec.New(20)}},
		"asset_y": {{C: dec.New(30)}, {C: dec.New(40)}, {C: dec.New(50)}},
	}
	want := 4 // Expect 4 enqueued jobs in buffered channel

	optimizer := BruteOptimizer{}
	ch := optimizer.enqueueJobs(givePSets, giveSamples)
	act := len(ch)
	assert.Equal(t, want, act)
}

func TestProcessBruteJobs(t *testing.T) {

	giveSample := []market.Kline{{C: dec.New(10)}, {C: dec.New(20)}, {C: dec.New(30)}, {C: dec.New(40)}, {C: dec.New(50)}}
	giveMakeBot := func(map[string]any) (trader.Bot, error) { return &trader.StubBot{}, nil }
	giveMakeDealer := func() (broker.SimulatedDealer, error) { return &broker.StubDealer{}, nil }
	giveJobCh := make(chan bruteOptimizerJob)
	giveDoneCh := make(chan struct{})
	defer close(giveDoneCh)

	outChan := processBruteJobs(context.Background(), giveDoneCh, giveJobCh, 8)

	giveJobCh <- bruteOptimizerJob{
		ParamSet: ParamSet{ID: "0", Params: map[string]any{"A": 0, "B": 1}},
		Sample:   giveSample, WarmupBarCount: 3, MakeBot: giveMakeBot, MakeDealer: giveMakeDealer,
	}
	giveJobCh <- bruteOptimizerJob{
		ParamSet: ParamSet{ID: "1", Params: map[string]any{"Y": 25, "Z": 26}},
		Sample:   giveSample, WarmupBarCount: 3, MakeBot: giveMakeBot, MakeDealer: giveMakeDealer,
	}
	close(giveJobCh)

	want := []OptimizerTrial{
		{PSet: ParamSet{ID: "0", Params: map[string]any{"A": 0, "B": 1}}},
		{PSet: ParamSet{ID: "1", Params: map[string]any{"Y": 25, "Z": 26}}},
	}

	var act []OptimizerTrial
	for step := range outChan {
		step.Result = perf.PerformanceReport{} // Set to empty for test equality
		act = append(act, step)
	}

	assert.ElementsMatch(t, want, act)
}

func TestSplitSample(t *testing.T) {
	tests := []struct {
		name         string
		giveSample   []market.Kline
		giveSplitPct float64
		wantASample  []market.Kline
		wantBSample  []market.Kline
	}{
		{
			name:         "ok split",
			giveSample:   []market.Kline{{C: dec.New(10)}, {C: dec.New(20)}, {C: dec.New(30)}},
			giveSplitPct: 0.3,
			wantASample:  []market.Kline{{C: dec.New(10)}},
			wantBSample:  []market.Kline{{C: dec.New(20)}, {C: dec.New(30)}},
		},
		{
			name:         "50/50",
			giveSample:   []market.Kline{{C: dec.New(10)}, {C: dec.New(20)}, {C: dec.New(30)}, {C: dec.New(40)}},
			giveSplitPct: 0.5,
			wantASample:  []market.Kline{{C: dec.New(10)}, {C: dec.New(20)}},
			wantBSample:  []market.Kline{{C: dec.New(30)}, {C: dec.New(40)}},
		},
		{
			name:         "no split = same samples in A & B",
			giveSample:   []market.Kline{{C: dec.New(10)}, {C: dec.New(20)}, {C: dec.New(30)}},
			giveSplitPct: 0,
			wantASample:  []market.Kline{{C: dec.New(10)}, {C: dec.New(20)}, {C: dec.New(30)}},
			wantBSample:  []market.Kline{{C: dec.New(10)}, {C: dec.New(20)}, {C: dec.New(30)}},
		},
		{
			name:         "100 pct split",
			giveSample:   []market.Kline{{C: dec.New(10)}, {C: dec.New(20)}, {C: dec.New(30)}},
			giveSplitPct: 1,
			wantASample:  []market.Kline{{C: dec.New(10)}, {C: dec.New(20)}, {C: dec.New(30)}},
			wantBSample:  []market.Kline{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actASample, actBSample := splitSample(tt.giveSample, tt.giveSplitPct)
			assert.Equal(t, tt.wantASample, actASample)
			assert.Equal(t, tt.wantBSample, actBSample)
		})
	}
}

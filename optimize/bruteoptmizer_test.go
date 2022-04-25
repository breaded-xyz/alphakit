package optimize

import (
	"context"
	"testing"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/dec"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/trader"
	"github.com/stretchr/testify/assert"
)

func TestBruteOptimizer_EnqueueJobs(t *testing.T) {

	givePSets := []ParamSet{
		{ID: "0", Params: map[string]any{"A": 0, "B": 1}},
		{ID: "1", Params: map[string]any{"Y": 25, "Z": 26}},
	}
	giveSamples := [][]market.Kline{
		{{C: dec.New(10)}, {C: dec.New(20)}},
		{{C: dec.New(30)}, {C: dec.New(40)}, {C: dec.New(50)}},
	}
	want := 4 // Expect 4 enqueued jobs

	optimizer := BruteOptimizer{}
	ch := optimizer.enqueueJobs(givePSets, giveSamples)
	act := len(ch)
	assert.Equal(t, want, act)
}

func TestProcessBruteJobs(t *testing.T) {

	giveSample := []market.Kline{{C: dec.New(10)}, {C: dec.New(20)}, {C: dec.New(30)}, {C: dec.New(40)}, {C: dec.New(50)}}
	giveMakeBot := func(map[string]any) (trader.Bot, error) { return &trader.StubBot{}, nil }
	giveMakeDealer := func(map[string]any) (broker.SimulatedDealer, error) { return &broker.StubDealer{}, nil }
	giveJobCh := make(chan bruteOptimizerJob)
	giveDoneCh := make(chan struct{})
	defer close(giveDoneCh)

	outChan := processBruteJobs(context.Background(), giveDoneCh, giveJobCh)

	giveJobCh <- bruteOptimizerJob{
		ParamSet: ParamSet{ID: "0", Params: map[string]any{"A": 0, "B": 1}},
		Sample:   giveSample, WarmupBarCount: 3, MakeBot: giveMakeBot, MakeDealer: giveMakeDealer,
	}
	giveJobCh <- bruteOptimizerJob{
		ParamSet: ParamSet{ID: "1", Params: map[string]any{"Y": 25, "Z": 26}},
		Sample:   giveSample, WarmupBarCount: 3, MakeBot: giveMakeBot, MakeDealer: giveMakeDealer,
	}
	close(giveJobCh)

	want := []OptimizerStep{
		{ParamSet: ParamSet{ID: "0", Params: map[string]any{"A": 0, "B": 1}}},
		{ParamSet: ParamSet{ID: "1", Params: map[string]any{"Y": 25, "Z": 26}}},
	}

	var act []OptimizerStep
	for step := range outChan {
		act = append(act, step)
	}

	assert.ElementsMatch(t, want, act)
}

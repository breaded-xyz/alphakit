package optimize

import (
	"context"
	"testing"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/dec"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/trader"
	"github.com/davecgh/go-spew/spew"
)

func TestProcessJobs(t *testing.T) {

	t.Parallel()

	giveSample := []market.Kline{
		{C: dec.New(10)},
		{C: dec.New(20)},
		{C: dec.New(30)},
		{C: dec.New(40)},
		{C: dec.New(50)},
	}
	giveMakeBot := func() trader.ConfigurableBot { return &trader.StubBot{} }
	giveMakeDealer := func() broker.SimulatedDealer { return &broker.StubDealer{} }

	jobCh := make(chan optimizerJob)
	doneCh := make(chan struct{})
	defer close(doneCh)

	outChan := processJobs(context.Background(), doneCh, jobCh)

	jobCh <- optimizerJob{
		ParamSet:       ParamSet{ID: "0", Params: map[string]any{"A": 0, "B": 1}},
		Sample:         giveSample,
		WarmupBarCount: 3,
		MakeBot:        giveMakeBot,
		MakeDealer:     giveMakeDealer,
	}
	jobCh <- optimizerJob{
		ParamSet:       ParamSet{ID: "1", Params: map[string]any{"Y": 25, "Z": 26}},
		Sample:         giveSample,
		WarmupBarCount: 3,
		MakeBot:        giveMakeBot,
		MakeDealer:     giveMakeDealer,
	}
	close(jobCh)

	for step := range outChan {
		spew.Dump(step)
	}

}

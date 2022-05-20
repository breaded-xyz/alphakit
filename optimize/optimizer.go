package optimize

import (
	"context"

	"github.com/thecolngroup/alphakit/market"
	"github.com/thecolngroup/alphakit/perf"
)

type Phase int

const (
	Training Phase = iota + 1
	Validation
)

func (p Phase) String() string {
	return [...]string{"None", "Training", "Validation"}[p]
}

type OptimizerStep struct {
	Phase  Phase
	PSet   ParamSet
	Result perf.PerformanceReport
	Err    error
}

type Optimizer interface {
	Prepare(ParamMap, [][]market.Kline) (int, error)
	Start(context.Context) (<-chan OptimizerStep, error)
	Study() Study
}

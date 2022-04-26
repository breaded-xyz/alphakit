package optimize

import (
	"context"

	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/perf"
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

type Optimizer interface {
	Prepare(ParamRange, [][]market.Kline) (int, error)
	Start(context.Context) (<-chan OptimizerStep, error)
	Study() Study
}

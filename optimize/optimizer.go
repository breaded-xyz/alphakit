package optimize

import (
	"context"

	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/perf"
)

type StepResult struct {
	Report perf.PerformanceReport
}

type Optimizer interface {
	Configure(map[string]any) error
	Prepare(params map[string]any, reader market.KlineReader) (int, error)
	Start(context.Context) (chan<- StepResult, error)
}

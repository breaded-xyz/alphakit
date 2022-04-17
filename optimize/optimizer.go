package optimize

import (
	"context"

	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/perf"
	"github.com/colngroup/zero2algo/trader"
)

type StepResult struct {
	Report perf.PerformanceReport
}

type Optimizer interface {
	Configure(map[string]any) error
	Prepare(trader.BotMakerFunc, map[string]any) (int, error)
	Start(context.Context, []market.Kline, float64) (chan<- StepResult, error)
}

// Package optimize provides a set of services for optimizing algo parameters.
// Parameter optimization is a process of systematically searching for the optimal set of parameters given a target objective.
// Multiple methods are available, each implementing the Optimizer interface.
package optimize

import (
	"context"

	"github.com/thecolngroup/alphakit/market"
	"github.com/thecolngroup/alphakit/perf"
)

// Phase is the phase of the optimization method.
type Phase int

const (
	// Training is the in-sample phase used to select the optimal parameters.
	Training Phase = iota + 1

	// Validation is the out-of-sample phase used to evaluate the performance of the optimal parameters.
	Validation
)

// String returns the string representation of the phase.
func (p Phase) String() string {
	return [...]string{"None", "Training", "Validation"}[p]
}

// MarshalText is used to output the phase as a string for CSV rendering.
func (p Phase) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

// Optimizer is the interface for an optimization method.
type Optimizer interface {
	Prepare(ParamMap, map[AssetID][]market.Kline) (int, error)
	Start(context.Context) (<-chan OptimizerTrial, error)
	Study() Study
}

// OptimizerTrial is a discrete trial conducted by an Optimizer on a single ParamSet.
type OptimizerTrial struct {
	Phase  Phase
	PSet   ParamSet
	Result perf.PerformanceReport
	Err    error
}

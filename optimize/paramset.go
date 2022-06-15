package optimize

import "github.com/thecolngroup/gou/id"

// ParamSet is a set of algo parameters to trial.
type ParamSet struct {
	ID     ParamSetID `csv:"id"`
	Params ParamMap   `csv:"params"`
}

// ParamSetID is a unique identifier for a ParamSet.
type ParamSetID string

// ParamMap is a map of algo parameters.
type ParamMap map[string]any

// NewParamSet returns a new param set with initialized ID
func NewParamSet() ParamSet {
	return ParamSet{
		ID:     ParamSetID(id.New()),
		Params: make(map[string]any),
	}
}

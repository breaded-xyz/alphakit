package optimize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildBacktestCases(t *testing.T) {
	give := ParamRange{
		"fast": []any{0, 10, 20},
		"slow": []any{0, 100},
	}

	want := []TestCase{
		{"fast": 0, "slow": 0},
		{"fast": 10, "slow": 0},
		{"fast": 20, "slow": 0},
		{"fast": 0, "slow": 100},
		{"fast": 10, "slow": 100},
		{"fast": 20, "slow": 100},
	}

	act := BuildBacktestCases(give)
	assert.ElementsMatch(t, want, act)
}

package optimize

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

func TestSharpeRanker(t *testing.T) {
	give := []Report{
		{Sharpe: 2},
		{Sharpe: 0.9},
		{Sharpe: 2.5},
	}

	want := []Report{give[1], give[0], give[2]}

	slices.SortFunc(give, SharpeRanker)

	assert.Equal(t, want, give)
}

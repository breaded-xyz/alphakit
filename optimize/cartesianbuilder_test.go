// Copyright 2022 The Coln Group Ltd
// SPDX-License-Identifier: MIT

package optimize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCartesianBuilder(t *testing.T) {
	give := map[string]any{
		"capital": 1000,
		"fast":    []any{0, 10, 20},
		"slow":    []any{0, 100},
	}

	want := []CartesianProduct{
		{"capital": 1000, "fast": 0, "slow": 0},
		{"capital": 1000, "fast": 10, "slow": 0},
		{"capital": 1000, "fast": 20, "slow": 0},
		{"capital": 1000, "fast": 0, "slow": 100},
		{"capital": 1000, "fast": 10, "slow": 100},
		{"capital": 1000, "fast": 20, "slow": 100},
	}

	act := CartesianBuilder(give)
	assert.ElementsMatch(t, want, act)
}

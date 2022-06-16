// Copyright 2022 The Coln Group Ltd
// SPDX-License-Identifier: MIT

package ta

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSD(t *testing.T) {

	giveV := []float64{10, 89, 20, 43, 44, 10}
	giveN := 3
	want := 19.347695814575268

	ind := NewSD(giveN)
	err := ind.Update(giveV...)
	act := ind.Value()

	assert.NoError(t, err)
	assert.True(t, ind.Valid())
	assert.Equal(t, want, act)
}

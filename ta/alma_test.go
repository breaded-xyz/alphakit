// Copyright 2022 The Coln Group Ltd
// SPDX-License-Identifier: MIT

package ta

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestALMA(t *testing.T) {

	tests := []struct {
		name       string
		giveV      []float64
		giveLength int
		want       float64
	}{
		{
			name:       "Valid sample",
			giveV:      []float64{10, 89, 20, 43, 44, 33, 19},
			giveLength: 3,
			want:       32.68047906324239,
		},
		{
			name:       "0 length window",
			giveV:      []float64{10, 89, 20, 43, 44, 33, 19},
			giveLength: 0,
			want:       19,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ind := NewALMA(tt.giveLength)
			err := ind.Update(tt.giveV...)
			act := ind.Value()

			assert.NoError(t, err)
			assert.True(t, ind.Valid())
			assert.Equal(t, tt.want, act)
		})
	}

}

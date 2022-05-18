package util

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundTo(t *testing.T) {
	giveX, giveY := 0.036, 0.01
	want := 0.04
	act := RoundTo(giveX, giveY)
	assert.Equal(t, want, act)
}

func TestNNZ(t *testing.T) {
	tests := []struct {
		name  string
		giveX float64
		giveY float64
		want  float64
	}{
		{
			name:  "inf",
			giveX: math.Inf(0),
			giveY: 99,
			want:  99,
		},
		{
			name:  "NaN",
			giveX: math.NaN(),
			giveY: 99,
			want:  99,
		},
		{
			name:  "zero",
			giveX: 0,
			giveY: 99,
			want:  99,
		},
		{
			name:  "no effect",
			giveX: 5,
			giveY: 99,
			want:  5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act := NNZ(tt.giveX, tt.giveY)
			assert.Equal(t, tt.want, act)
		})
	}
}

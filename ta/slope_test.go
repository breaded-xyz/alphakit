package ta

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlope(t *testing.T) {
	tests := []struct {
		name   string
		giveT1 float64
		giveT2 float64
		want   int
	}{
		{
			name:   "slope up",
			giveT1: 10,
			giveT2: 15,
			want:   1,
		},
		{
			name:   "slope down",
			giveT1: 15,
			giveT2: 10,
			want:   -1,
		},
		{
			name:   "flat",
			giveT1: 15,
			giveT2: 15,
			want:   0,
		},
		{
			name:   "negative numbers",
			giveT1: -10,
			giveT2: -5,
			want:   1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act := Slope(tt.giveT1, tt.giveT2)
			assert.Equal(t, tt.want, act)
		})
	}
}

package ta

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrossUp(t *testing.T) {
	tests := []struct {
		name       string
		giveSeries []float64
		giveX      float64
		want       bool
	}{
		{
			name:       "cross up",
			giveSeries: []float64{10, 20},
			giveX:      15,
			want:       true,
		},
		{
			name:       "cross up start at x",
			giveSeries: []float64{15, 20},
			giveX:      15,
			want:       true,
		},
		{
			name:       "cross up end at x",
			giveSeries: []float64{15, 20},
			giveX:      20,
			want:       false,
		},
		{
			name:       "cross down",
			giveSeries: []float64{1, -1},
			giveX:      0,
			want:       false,
		},
		{
			name:       "flat",
			giveSeries: []float64{15, 15},
			giveX:      15,
			want:       false,
		},
		{
			name:       "only last 2 series values are evaluated",
			giveSeries: []float64{4, 7, 0, 1},
			giveX:      0.5,
			want:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act := CrossUp(tt.giveSeries, tt.giveX)
			assert.Equal(t, tt.want, act)
		})
	}
}

func TestCrossDown(t *testing.T) {
	tests := []struct {
		name       string
		giveSeries []float64
		giveX      float64
		want       bool
	}{
		{
			name:       "cross down",
			giveSeries: []float64{20, 10},
			giveX:      15,
			want:       true,
		},
		{
			name:       "cross down start at x",
			giveSeries: []float64{20, 15},
			giveX:      20,
			want:       true,
		},
		{
			name:       "cross down end at x",
			giveSeries: []float64{25, 20},
			giveX:      20,
			want:       false,
		},
		{
			name:       "cross up",
			giveSeries: []float64{-1, 1},
			giveX:      0,
			want:       false,
		},
		{
			name:       "flat",
			giveSeries: []float64{15, 15},
			giveX:      15,
			want:       false,
		},
		{
			name:       "only last 2 series values are evaluated",
			giveSeries: []float64{4, 7, 1, -1},
			giveX:      0,
			want:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act := CrossDown(tt.giveSeries, tt.giveX)
			assert.Equal(t, tt.want, act)
		})
	}
}

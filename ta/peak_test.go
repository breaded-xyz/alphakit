package ta

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPeak(t *testing.T) {
	tests := []struct {
		name       string
		giveSeries []float64
		want       bool
	}{
		{
			name:       "peak",
			giveSeries: []float64{5, 10, 5},
			want:       true,
		},
		{
			name:       "no peak: flat top",
			giveSeries: []float64{5, 10, 10, 5},
			want:       false,
		},
		{
			name:       "no peak: trend down with flat",
			giveSeries: []float64{20, 10, 10, 5},
			want:       false,
		},
		{
			name:       "no peak: trend up with flat",
			giveSeries: []float64{5, 10, 10, 20},
			want:       false,
		},
		{
			name:       "no peak: valley",
			giveSeries: []float64{10, 5, 10},
			want:       false,
		},
		{
			name:       "no peak: trend up",
			giveSeries: []float64{10, 20, 30},
			want:       false,
		},
		{
			name:       "no peak: trend down",
			giveSeries: []float64{30, 20, 10},
			want:       false,
		},
		{
			name:       "no peak: flat",
			giveSeries: []float64{10, 10, 10},
			want:       false,
		},
		{
			name:       "no peak: empty",
			giveSeries: []float64{},
			want:       false,
		},
		{
			name:       "no peak: missing data",
			giveSeries: []float64{10, 5},
			want:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act := Peak(tt.giveSeries)
			assert.Equal(t, tt.want, act)
		})
	}
}

func TestValley(t *testing.T) {
	tests := []struct {
		name       string
		giveSeries []float64
		want       bool
	}{
		{
			name:       "valley",
			giveSeries: []float64{10, 5, 10},
			want:       true,
		},
		{
			name:       "no valley: flat bottom",
			giveSeries: []float64{10, 5, 5, 10},
			want:       false,
		},
		{
			name:       "no valley: trend down with flat",
			giveSeries: []float64{20, 10, 10, 5},
			want:       false,
		},
		{
			name:       "no valley: trend up with flat",
			giveSeries: []float64{5, 10, 10, 20},
			want:       false,
		},
		{
			name:       "no valley: peak",
			giveSeries: []float64{5, 10, 5},
			want:       false,
		},
		{
			name:       "no valley: trend up",
			giveSeries: []float64{10, 20, 30},
			want:       false,
		},
		{
			name:       "no valley: trend down",
			giveSeries: []float64{30, 20, 10},
			want:       false,
		},
		{
			name:       "no valley: flat",
			giveSeries: []float64{10, 10, 10},
			want:       false,
		},
		{
			name:       "no valley: empty",
			giveSeries: []float64{},
			want:       false,
		},
		{
			name:       "no valley: missing data",
			giveSeries: []float64{10, 5},
			want:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act := Valley(tt.giveSeries)
			assert.Equal(t, tt.want, act)
		})
	}
}

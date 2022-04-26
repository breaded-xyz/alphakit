package perf

import (
	"math"
	"testing"
	"time"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/internal/dec"
	"github.com/stretchr/testify/assert"
)

func TestDiffReturns(t *testing.T) {
	give := broker.EquitySeries{
		1: dec.New(10),
		2: dec.New(20), // 20 - 10 = 10 / 10 = 1
		3: dec.New(30), // 30 - 20 = 10 / 20 = 0.5
		4: dec.New(5),  // 5 - 30 = -25 / 30 = -0.8333333333333333
	}
	want := []float64{1, 0.5, -0.8333333333333333}
	act := DiffReturns(give)
	assert.Equal(t, want, act)
}

func TestReduceEoD(t *testing.T) {

	datum := time.Date(2022, time.January, 1, 0, 0, 0, 0, time.Local)
	give := broker.EquitySeries{
		broker.Timestamp(datum.UnixMilli()):                     dec.New(10),
		broker.Timestamp(datum.Add(25 * time.Hour).UnixMilli()): dec.New(20),
		broker.Timestamp(datum.Add(48 * time.Hour).UnixMilli()): dec.New(30),
	}
	want := broker.EquitySeries{
		broker.Timestamp(datum.UnixMilli()):                     dec.New(10),
		broker.Timestamp(datum.Add(48 * time.Hour).UnixMilli()): dec.New(30),
	}

	act := ReduceEOD(give)
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

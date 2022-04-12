package perf

import (
	"math"
	"testing"
	"time"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/dec"
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

func TestCAGR(t *testing.T) {
	giveInitial := 1000.0
	giveFinal := 2500.0
	giveDays := 190
	want := 4.81
	act := CAGR(giveInitial, giveFinal, giveDays)
	assert.Equal(t, want, round2DP(act))
}

func TestKellyCriterion(t *testing.T) {
	giveProfitFactor := 1.6
	giveWinP := 0.7
	want := 0.51
	act := KellyCriterion(giveProfitFactor, giveWinP)
	assert.Equal(t, want, round2DP(act))
}

func TestSharpeRatio(t *testing.T) {
	give := []float64{0.1, 0.2, -0.15, 0.1, 0.8, -0.3, 0.2}
	exp := 6.20
	act := SharpeRatio(give, SharpeAnnualRiskFreeRate)
	assert.Equal(t, exp, round2DP(act))
}

func round2DP(x float64) float64 {
	return math.Round(x*100) / 100
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

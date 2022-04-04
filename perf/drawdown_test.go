package perf

import (
	"testing"
	"time"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/dec"
	"github.com/stretchr/testify/assert"
)

func TestDrawdowns(t *testing.T) {

	give := broker.EquitySeries{
		1: dec.New(10),
		2: dec.New(20),
		3: dec.New(15),
		4: dec.New(10),
		5: dec.New(18),
		6: dec.New(19),
		7: dec.New(30),
	}

	exp := []Drawdown{
		{
			HighAt:   time.UnixMilli(2),
			LowAt:    time.UnixMilli(4),
			StartAt:  time.UnixMilli(3),
			EndAt:    time.UnixMilli(7),
			High:     dec.New(20),
			Low:      dec.New(10),
			Recovery: 4 * time.Millisecond,
			Amount:   dec.New(10),
			Pct:      0.5,
		},
	}

	act := Drawdowns(give)
	assert.Equal(t, exp, act)

	/*
		// Case: initial equity is first peak
		curve = []alpha.Equity{
			{At: start.Add(barSize * 0), C: dec.Int64(15)}, // Peak #1
			{At: start.Add(barSize * 1), C: dec.Int64(0)},  // Trough #1
			{At: start.Add(barSize * 2), C: dec.Int64(15)}, // Drawdown #1 = 10
		}
		ds = Drawdowns(curve)
		assert.Len(t, ds, 1)
		dec.AssertEqual(t, dec.Int64(15), ds[0].Amount)

		// Case: sequence of drawdowns
		curve = []alpha.Equity{
			{At: start.Add(barSize * 0), C: dec.Int64(15)}, // Peak #1
			{At: start.Add(barSize * 1), C: dec.Int64(0)},  // Trough #1
			{At: start.Add(barSize * 2), C: dec.Int64(15)}, // Drawdown #1 = 15
			{At: start.Add(barSize * 3), C: dec.Int64(20)}, // Peak #2
			{At: start.Add(barSize * 4), C: dec.Int64(19)},
			{At: start.Add(barSize * 5), C: dec.Int64(10)}, // Trough #2
			{At: start.Add(barSize * 6), C: dec.Int64(11)},
			{At: start.Add(barSize * 7), C: dec.Int64(17)},
			{At: start.Add(barSize * 8), C: dec.Int64(30)},  // Drawdown #2 = 10 | Peak #3
			{At: start.Add(barSize * 9), C: dec.Int64(25)},  // Trough #3
			{At: start.Add(barSize * 10), C: dec.Int64(30)}, // Drawdown #3 = 5 | Peak #4
			{At: start.Add(barSize * 11), C: dec.Int64(10)},
			{At: start.Add(barSize * 12), C: dec.Int64(5)},  // Trough #4
			{At: start.Add(barSize * 13), C: dec.Int64(25)}, // Drawdown #4 because end of curve (status = open)
		}
		ds = Drawdowns(curve)
		assert.Len(t, ds, 4)
		dec.AssertEqual(t, dec.Int64(5), ds[2].Amount)

		// Case: max drawdown
		mdd := MaxDrawdown(ds)
		dec.AssertEqual(t, dec.Int64(25), mdd.Amount)*/
}

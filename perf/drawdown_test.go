package perf

import (
	"testing"
)

func TestDrawdowns(t *testing.T) {

	/*start := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
	barSize := time.Hour

	// Case: base
	curve := []alpha.Equity{
		{At: start.Add(barSize * 0), C: dec.Int64(10)},
		{At: start.Add(barSize * 1), C: dec.Int64(20)}, // Peak #1
		{At: start.Add(barSize * 2), C: dec.Int64(15)},
		{At: start.Add(barSize * 3), C: dec.Int64(10)}, // Trough #1
		{At: start.Add(barSize * 4), C: dec.Int64(18)},
		{At: start.Add(barSize * 5), C: dec.Int64(19)},
		{At: start.Add(barSize * 6), C: dec.Int64(30)}, // Drawdown #1 = 10
	}
	ds := Drawdowns(curve)
	assert.Len(t, ds, 1)
	assert.Equal(t, start.Add(barSize*1), ds[0].HighAt)
	assert.Equal(t, start.Add(barSize*3), ds[0].LowAt)
	dec.AssertEqual(t, dec.Int64(20), ds[0].High)
	dec.AssertEqual(t, dec.Int64(10), ds[0].Low)
	dec.AssertEqual(t, dec.Int64(10), ds[0].Amount)
	assert.Equal(t, 0.5, ds[0].Pct)
	assert.Equal(t, 4*time.Hour, ds[0].Recovery)

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

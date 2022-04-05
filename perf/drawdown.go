package perf

import (
	"time"

	"github.com/colngroup/zero2algo/broker"
	"github.com/shopspring/decimal"
)

// Drawdown is a pattern in the equity curve and represents a peak to valley and recovery.
type Drawdown struct {
	HighAt  time.Time
	LowAt   time.Time
	StartAt time.Time
	EndAt   time.Time

	High decimal.Decimal
	Low  decimal.Decimal

	Recovery time.Duration

	Amount decimal.Decimal
	Pct    float64

	IsOpen bool
}

// Drawdowns extracts all the drawdowns from the equity curve.
func Drawdowns(curve broker.EquitySeries) []Drawdown {
	if len(curve) == 0 {
		return nil
	}

	var dds []Drawdown

	// Iterate the series in chronological order
	for i, k := range curve.SortKeys() {
		t, v := k.Time(), curve[k]

		// Init an empty DD to begin tracking changes as we walk the equity curve
		if i == 0 {
			dds = append(dds, Drawdown{HighAt: t, High: v, LowAt: t, Low: v})
			continue
		}

		// Get pointer to latest DD
		dd := &dds[len(dds)-1]

		// Case: end of curve is reached
		// If a drawdown is open close it based on last equity point
		// IsOpen field is set to flag drawdown is not a complete recovery
		if i == len(curve)-1 && !dd.StartAt.IsZero() {
			dd.EndAt = t
			dd.Recovery = t.Sub(dd.StartAt)
			dd.Amount = dd.High.Sub(dd.Low)
			dd.Pct = dd.Amount.Div(dd.High).InexactFloat64()
			dd.IsOpen = true
			continue
		}

		// Case: new lower low
		if v.LessThanOrEqual(dd.Low) {
			// Open a new DD if not already started
			if dd.StartAt.IsZero() {
				dd.StartAt = t
			}
			// Update the DD low
			dd.LowAt, dd.Low = t, v
			continue
		}

		// Case: new higher high
		if v.GreaterThanOrEqual(dd.High) {

			// If DD not open then continue to mark high and low to curve
			if dd.StartAt.IsZero() {
				dd.HighAt, dd.High = t, v
				dd.LowAt, dd.Low = t, v
				continue
			}

			// Else the DD was open and has recovered from a low so close it
			dd.EndAt = t
			dd.Recovery = t.Sub(dd.StartAt)
			dd.Amount = dd.High.Sub(dd.Low)
			dd.Pct = dd.Amount.Div(dd.High).InexactFloat64()

			// Open new empty DD ready for next iteration
			dds = append(dds, Drawdown{HighAt: t, High: v, LowAt: t, Low: v})
			continue
		}

	}

	// If final DD was empty then strip from slice
	if dds[len(dds)-1].StartAt.IsZero() {
		dds = dds[:len(dds)-1]
	}

	return dds
}

// MaxDrawdown finds the largest drawdown based on the percentage amount.
func MaxDrawdown(dds []Drawdown) (max Drawdown) {
	for i := range dds {
		d := dds[i]
		if d.Pct >= max.Pct {
			max = d
		}
	}
	return max
}

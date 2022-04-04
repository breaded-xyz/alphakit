package perf

import (
	"time"

	"github.com/colngroup/zero2algo/broker"
	"github.com/shopspring/decimal"
)

// Drawdown is a pattern in the equity curve: peak to valley to recovery.
type Drawdown struct {
	High   decimal.Decimal
	HighAt time.Time

	Low   decimal.Decimal
	LowAt time.Time

	StartAt time.Time
	EndAt   time.Time

	Recovery time.Duration

	Amount decimal.Decimal
	Pct    float64

	IsOpen bool
}

// Drawdowns extracts all the drawdowns from the equity curve using close value.
func Drawdowns(curve broker.EquitySeries) []Drawdown {
	if len(curve) == 0 {
		return nil
	}

	var dds []Drawdown

	/*for i := range curve {
		curr := curve[i]

		// Init a starting (empty) DD to begin tracking changes as we walk the equity curve
		if i == 0 {
			dds = append(dds, Drawdown{HighAt: curr.At, High: curr.C, LowAt: curr.At, Low: curr.C})
			continue
		}

		// Get pointer to latest DD
		dd := &dds[len(dds)-1]

		// End of curve is reached so calc DD values based on last equity point
		if i == len(curve)-1 && !dd.StartAt.IsZero() {
			dd.EndAt = curr.At
			dd.Recovery = curr.At.Sub(dd.StartAt)
			dd.Amount = dd.High.Sub(dd.Low)
			dd.Pct = dd.Amount.Div(dd.High).InexactFloat64()
			dd.IsOpen = true
			continue
		}

		// Track lower low for current open DD
		if curr.C.LessThanOrEqual(dd.Low) {
			// Detect if this is the start of the drawdown
			if dd.StartAt.IsZero() {
				dd.StartAt = curr.At
			}
			dd.LowAt, dd.Low = curr.At, curr.C
			continue
		}

		// Track higher high for current open DD
		if curr.C.GreaterThanOrEqual(dd.High) {

			// If current DD in initial empty state then move high and low up together
			if dd.StartAt.IsZero() {
				dd.HighAt, dd.High = curr.At, curr.C
				dd.LowAt, dd.Low = curr.At, curr.C
				continue
			}

			// Else a DD has recovered from a low so close current drawdown
			dd.EndAt = curr.At
			dd.Recovery = curr.At.Sub(dd.StartAt)
			dd.Amount = dd.High.Sub(dd.Low)
			dd.Pct = dd.Amount.Div(dd.High).InexactFloat64()

			// Open new DD to continue tracking
			dds = append(dds, Drawdown{HighAt: curr.At, High: curr.C, LowAt: curr.At, Low: curr.C})
			continue
		}

	}

	// If final DD was empty then strip from sequence
	if dds[len(dds)-1].StartAt.IsZero() {
		dds = dds[:len(dds)-1]
	}*/

	return dds
}

// MaxDrawdown finds the biggest drawdown based on the currency amount.
func MaxDrawdown(dds []Drawdown) (max Drawdown) {
	for i := range dds {
		d := dds[i]
		if d.Amount.GreaterThanOrEqual(max.Amount) {
			max = d
		}
	}
	return max
}

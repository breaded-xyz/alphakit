package backtest

import (
	"testing"
	"time"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/dec"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestPerpCostFunding2(t *testing.T) {
	// Test evaluates how state mutates over time in response to elapsed duration input
	// so the same PerpCost instance is used in all sub tests
	cost := PerpCost{
		FundingHourPct: dec.New(0.1),
	}
	givePosition := broker.Position{
		OpenedAt: time.Now(),
		Size:     dec.New(1),
	}
	givePrice := dec.New(10)
	var giveElapsed time.Duration

	tests := []struct {
		name string
		give time.Duration
		want decimal.Decimal
	}{
		{
			name: "initial funding",
			give: 10 * time.Hour,
			want: dec.New(10),
		},
		{
			name: "next epoch",
			give: 2 * time.Hour,
			want: dec.New(2),
		},
		{
			name: "no funding as interval less than hour",
			give: 50 * time.Minute,
			want: decimal.Zero,
		},
		{
			name: "now fund as partial interval completes",
			give: 10 * time.Minute,
			want: dec.New(1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			giveElapsed += tt.give
			act := cost.Funding(givePosition, givePrice, giveElapsed)
			// Assert IntPart because decimal.Decimal stores the same equivalent value
			// differently depending on its source (float vs int vs string)
			assert.Equal(t, tt.want.IntPart(), act.IntPart())
		})
	}
}

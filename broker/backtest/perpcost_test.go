package backtest

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/thecolngroup/alphakit/broker"
	"github.com/thecolngroup/gou/dec"
)

func TestPerpCostFunding(t *testing.T) {
	// Test evaluates how state mutates over time in response to elapsed duration input
	// Same PerpCost instance is used in all sub tests
	cost := PerpCoster{
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

func TestPerpCostSlippage(t *testing.T) {
	cost := PerpCoster{
		SlippagePct: dec.New(0.1),
	}
	exp := dec.New(1)
	act := cost.Slippage(dec.New(10))
	assert.True(t, act.Equal(exp))
}

func TestPerpCostSpread(t *testing.T) {
	cost := PerpCoster{
		SpreadPct: dec.New(0.4),
	}
	exp := dec.New(2)
	act := cost.Spread(dec.New(10))
	assert.True(t, act.Equal(exp))
}

func TestPerpCostTransaction(t *testing.T) {
	cost := PerpCoster{
		TransactionPct: dec.New(0.5),
	}
	order := broker.Order{
		FilledPrice: dec.New(2),
		FilledSize:  dec.New(5),
	}
	exp := dec.New(5)
	act := cost.Transaction(order)
	assert.True(t, act.Equal(exp))
}

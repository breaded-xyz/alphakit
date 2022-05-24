package money

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/thecolngroup/dec"
)

func TestSafeFSizer_Size(t *testing.T) {
	tests := []struct {
		name        string
		giveSizer   Sizer
		givePrice   decimal.Decimal
		giveCapital decimal.Decimal
		giveRisk    decimal.Decimal
		want        decimal.Decimal
	}{
		{
			name: "ok",
			giveSizer: &SafeFSizer{
				InitialCapital: dec.New(500),
				F:              0.1,
				ScaleF:         0.5,
				StepSize:       DefaultStepSize,
			},
			givePrice:   dec.New(100),
			giveCapital: dec.New(1000),
			giveRisk:    dec.New(10),
			want:        dec.New(6.12),
		},
		{
			name: "zero or neg profit = fixed capital growth factor at 1.0",
			giveSizer: &SafeFSizer{
				InitialCapital: dec.New(5000),
				F:              0.1,
				ScaleF:         0.5,
				StepSize:       DefaultStepSize,
			},
			givePrice:   dec.New(100),
			giveCapital: dec.New(1000),
			giveRisk:    dec.New(10),
			want:        dec.New(5),
		},
		{
			name: "negative capital",
			giveSizer: &SafeFSizer{
				InitialCapital: dec.New(5000),
				F:              0.1,
				ScaleF:         0.5,
				StepSize:       DefaultStepSize,
			},
			givePrice:   dec.New(100),
			giveCapital: dec.New(-1000),
			giveRisk:    dec.New(10),
			want:        dec.New(-5),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act := tt.giveSizer.Size(tt.givePrice, tt.giveCapital, tt.giveRisk)
			spew.Dump(act)
			assert.Equal(t, tt.want, act)
		})
	}

}

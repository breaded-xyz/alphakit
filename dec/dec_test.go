package dec

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	exp := decimal.NewFromFloat(10)
	act := New(10)
	assert.Equal(t, exp, act)
}

func TestBetween(t *testing.T) {

	tests := []struct {
		name      string
		giveValue decimal.Decimal
		giveLower decimal.Decimal
		giveUpper decimal.Decimal
		want      bool
	}{
		{
			name:      "inside bounds",
			giveValue: New(10),
			giveLower: New(1),
			giveUpper: New(20),
			want:      true,
		},
		{
			name:      "inside at lower bound",
			giveValue: New(1),
			giveLower: New(1),
			giveUpper: New(20),
			want:      true,
		},
		{
			name:      "inside at upper bound",
			giveValue: New(20),
			giveLower: New(1),
			giveUpper: New(20),
			want:      true,
		},
		{
			name:      "outside lower bound",
			giveValue: New(0),
			giveLower: New(1),
			giveUpper: New(20),
			want:      false,
		},
		{
			name:      "outside upper bound",
			giveValue: New(21),
			giveLower: New(1),
			giveUpper: New(20),
			want:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act := Between(tt.giveValue, tt.giveLower, tt.giveUpper)
			assert.Equal(t, tt.want, act)
		})
	}
}

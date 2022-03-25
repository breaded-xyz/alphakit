package tradebot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHoldBot_Configure(t *testing.T) {
	tests := []struct {
		name string
		give map[string]any
		want HodlBot
		err  error
	}{
		{
			name: "Success",
			give: map[string]any{"buybarindex": 1, "sellbarindex": 1000},
			want: HodlBot{
				BuyBarIndex:  1,
				SellBarIndex: 1000,
			},
			err: nil,
		},
		{
			name: "Default config state is zero",
			give: map[string]any{"buybarindex": 0, "sellbarindex": 0},
			want: HodlBot{
				BuyBarIndex:  0,
				SellBarIndex: 0,
			},
			err: nil,
		},
		{
			name: "Buy index >= sell index",
			give: map[string]any{"buybarindex": 10, "sellbarindex": 5},
			want: HodlBot{
				BuyBarIndex:  0,
				SellBarIndex: 0,
			},
			err: ErrInvalidConfig,
		},
		{
			name: "Not int",
			give: map[string]any{"buybarindex": 10.5, "sellbarindex": 5},
			want: HodlBot{
				BuyBarIndex:  0,
				SellBarIndex: 0,
			},
			err: ErrInvalidConfig,
		},
		{
			name: "Key not found",
			give: map[string]any{"notakey": 10, "sellbarindex": 5},
			want: HodlBot{
				BuyBarIndex:  0,
				SellBarIndex: 0,
			},
			err: ErrInvalidConfig,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bot HodlBot
			err := bot.Configure(tt.give)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, bot)
		})
	}
}

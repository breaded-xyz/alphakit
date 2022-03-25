package tradebot

import (
	"context"
	"testing"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/broker/backtest"
	"github.com/colngroup/zero2algo/market"
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
			name: "ok: buy index < sell index",
			give: map[string]any{"buybarindex": 1, "sellbarindex": 1000},
			want: HodlBot{
				BuyBarIndex:  1,
				SellBarIndex: 1000,
			},
			err: nil,
		},
		{
			name: "ok: no sell",
			give: map[string]any{"buybarindex": 10, "sellbarindex": 0},
			want: HodlBot{
				BuyBarIndex:  10,
				SellBarIndex: 0,
			},
			err: nil,
		},
		{
			name: "ok: default",
			give: map[string]any{"buybarindex": 0, "sellbarindex": 0},
			want: HodlBot{
				BuyBarIndex:  0,
				SellBarIndex: 0,
			},
			err: nil,
		},
		{
			name: "err: buy index >= sell index",
			give: map[string]any{"buybarindex": 10, "sellbarindex": 5},
			want: HodlBot{
				BuyBarIndex:  0,
				SellBarIndex: 0,
			},
			err: ErrInvalidConfig,
		},
		{
			name: "err: not int",
			give: map[string]any{"buybarindex": 10.5, "sellbarindex": 5},
			want: HodlBot{
				BuyBarIndex:  0,
				SellBarIndex: 0,
			},
			err: ErrInvalidConfig,
		},
		{
			name: "err: neg int",
			give: map[string]any{"buybarindex": -1, "sellbarindex": 5},
			want: HodlBot{
				BuyBarIndex:  0,
				SellBarIndex: 0,
			},
			err: ErrInvalidConfig,
		},
		{
			name: "err: key not found",
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

func TestHodlBot_evalAlgo(t *testing.T) {
	tests := []struct {
		name string
		give []int
		want broker.OrderSide
	}{
		{
			name: "ok: default state",
			give: []int{0, 0, 0},
			want: broker.Buy,
		},
		{
			name: "ok: buy",
			give: []int{10, 10, 20},
			want: broker.Buy,
		},
		{
			name: "ok: sell",
			give: []int{20, 10, 20},
			want: broker.Sell,
		},
		{
			name: "ok: no sell",
			give: []int{0, 10, 0},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bot HodlBot
			actual := bot.evalAlgo(tt.give[0], tt.give[1], tt.give[2])
			assert.Equal(t, tt.want, actual)
		})
	}
}

func TestHodlBot_ReceivePrice(t *testing.T) {

	bot := HodlBot{
		dealer: backtest.NewDealer(),
	}
	err := bot.ReceivePrice(context.Background(), market.Kline{})
	assert.NoError(t, err)
	assert.Equal(t, 1, bot.barIndex)
}

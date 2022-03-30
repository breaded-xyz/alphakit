package bot

import (
	"context"
	"testing"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/dec"
	"github.com/colngroup/zero2algo/market"
	"github.com/stretchr/testify/assert"
)

func TestHoldBotConfigure(t *testing.T) {
	tests := []struct {
		name string
		give map[string]any
		want HodlBot
		err  error
	}{
		{
			name: "buy index < sell index",
			give: map[string]any{"buybarindex": 1, "sellbarindex": 1000},
			want: HodlBot{
				BuyBarIndex:  1,
				SellBarIndex: 1000,
			},
			err: nil,
		},
		{
			name: "no sell",
			give: map[string]any{"buybarindex": 10, "sellbarindex": 0},
			want: HodlBot{
				BuyBarIndex:  10,
				SellBarIndex: 0,
			},
			err: nil,
		},
		{
			name: "default",
			give: map[string]any{"buybarindex": 0, "sellbarindex": 0},
			want: HodlBot{
				BuyBarIndex:  0,
				SellBarIndex: 0,
			},
			err: nil,
		},
		{
			name: "buy index >= sell index",
			give: map[string]any{"buybarindex": 10, "sellbarindex": 5},
			want: HodlBot{
				BuyBarIndex:  0,
				SellBarIndex: 0,
			},
			err: ErrInvalidConfig,
		},
		{
			name: "not int",
			give: map[string]any{"buybarindex": 10.5, "sellbarindex": 5},
			want: HodlBot{
				BuyBarIndex:  0,
				SellBarIndex: 0,
			},
			err: ErrInvalidConfig,
		},
		{
			name: "neg int",
			give: map[string]any{"buybarindex": -1, "sellbarindex": 5},
			want: HodlBot{
				BuyBarIndex:  0,
				SellBarIndex: 0,
			},
			err: ErrInvalidConfig,
		},
		{
			name: "key not found",
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

func TestHodlBotEvalAlgo(t *testing.T) {
	tests := []struct {
		name string
		give []int // barIndex, buyIndex, sellIndex
		want broker.OrderSide
	}{
		{
			name: "default state",
			give: []int{0, 0, 0},
			want: broker.Buy,
		},
		{
			name: "buy",
			give: []int{10, 10, 20},
			want: broker.Buy,
		},
		{
			name: "sell",
			give: []int{20, 10, 20},
			want: broker.Sell,
		},
		{
			name: "no sell",
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

func TestHodlBotReceivePrice(t *testing.T) {
	expOrder := broker.Order{Type: broker.Market, Side: broker.Buy, Size: dec.New(1)}
	mock := &broker.MockDealer{}
	mock.On("PlaceOrder", context.Background(), expOrder)

	bot := HodlBot{dealer: mock}
	err := bot.ReceivePrice(context.Background(), market.Kline{})
	assert.NoError(t, err)
	assert.Equal(t, 1, bot.barIndex)
	mock.AssertExpectations(t)
}

func TestHodlBotClose(t *testing.T) {
	expOrder := broker.Order{Type: broker.Market, Side: broker.Sell, Size: dec.New(1), ReduceOnly: true}
	mock := &broker.MockDealer{}
	mock.On("PlaceOrder", context.Background(), expOrder)

	bot := HodlBot{dealer: mock}
	err := bot.Close(context.Background())
	assert.NoError(t, err)
	mock.AssertExpectations(t)
}

package hodl

import (
	"context"
	"testing"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/internal/dec"
	"github.com/colngroup/zero2algo/market"
	"github.com/stretchr/testify/assert"
)

func TestBot_evalAlgo(t *testing.T) {
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
			var bot Bot
			act := bot.evalAlgo(tt.give[0], tt.give[1], tt.give[2])
			assert.Equal(t, tt.want, act)
		})
	}
}

func TestBot_ReceivePrice(t *testing.T) {
	expOrder := broker.Order{Type: broker.Market, Side: broker.Buy, Size: dec.New(1)}
	mock := &broker.MockDealer{}
	mock.On("PlaceOrder", context.Background(), expOrder)

	bot := Bot{dealer: mock}
	err := bot.ReceivePrice(context.Background(), market.Kline{})
	assert.NoError(t, err)
	assert.Equal(t, 1, bot.barIndex)
	mock.AssertExpectations(t)
}

func TestBot_Close(t *testing.T) {
	expOrder := broker.Order{Type: broker.Market, Side: broker.Sell, Size: dec.New(1), ReduceOnly: true}
	mock := &broker.MockDealer{}
	mock.On("PlaceOrder", context.Background(), expOrder)

	bot := Bot{dealer: mock}
	err := bot.Close(context.Background())
	assert.NoError(t, err)
	mock.AssertExpectations(t)
}

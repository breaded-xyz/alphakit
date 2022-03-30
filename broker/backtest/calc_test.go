package backtest

import (
	"testing"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/dec"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestProfit(t *testing.T) {
	tests := []struct {
		name string
		give broker.Position
		want decimal.Decimal
	}{
		{
			name: "buy side profit",
			give: broker.Position{
				Side:             broker.Buy,
				Price:            dec.New(10),
				Size:             dec.New(2),
				LiquidationPrice: dec.New(20),
			},
			want: dec.New(20),
		},
		{
			name: "sell side profit",
			give: broker.Position{
				Side:             broker.Sell,
				Price:            dec.New(100),
				Size:             dec.New(2),
				LiquidationPrice: dec.New(50),
			},
			want: dec.New(100),
		},
		{
			name: "buy side loss",
			give: broker.Position{
				Side:             broker.Buy,
				Price:            dec.New(10),
				Size:             dec.New(2),
				LiquidationPrice: dec.New(5),
			},
			want: dec.New(-10),
		},
		{
			name: "sell side loss",
			give: broker.Position{
				Side:             broker.Sell,
				Price:            dec.New(10),
				Size:             dec.New(2),
				LiquidationPrice: dec.New(20),
			},
			want: dec.New(-20),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act := Profit(tt.give, tt.give.LiquidationPrice)
			assert.Equal(t, tt.want, act)
		})
	}
}

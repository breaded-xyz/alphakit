package backtest

import (
	"context"
	"testing"
	"time"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/dec"
	"github.com/colngroup/zero2algo/market"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestDealerProcessOrder(t *testing.T) {
	tests := []struct {
		name      string
		give      broker.Order
		wantOrder broker.Order
		wantState broker.OrderState
	}{
		{
			name: "ok: market order filled",
			give: broker.Order{
				Type: broker.Market,
				Size: dec.New(1),
			},
			wantOrder: broker.Order{
				FilledPrice: dec.New(10),
				FilledSize:  dec.New(1),
			},
			wantState: broker.Closed,
		},
		{
			name: "ok: limit order filled",
			give: broker.Order{
				Type:       broker.Limit,
				LimitPrice: dec.New(8),
				Size:       dec.New(1),
			},
			wantOrder: broker.Order{
				FilledPrice: dec.New(8),
				FilledSize:  dec.New(1),
			},
			wantState: broker.Closed,
		},
		{
			name: "ok: limit order open but not filled",
			give: broker.Order{
				Type:       broker.Limit,
				LimitPrice: dec.New(100),
				Size:       dec.New(1),
			},
			wantOrder: broker.Order{},
			wantState: broker.Open,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dealer := NewDealer()
			dealer.marketPrice = market.Kline{O: dec.New(8), H: dec.New(15), L: dec.New(5), C: dec.New(10)}
			act := dealer.processOrder(tt.give)
			assert.Equal(t, tt.wantOrder.FilledSize, act.FilledSize)
			assert.Equal(t, tt.wantOrder.FilledPrice, act.FilledPrice)
			assert.Equal(t, tt.wantState, act.State())
		})
	}
}

func TestDealerPlaceOrder_InvalidArgs(t *testing.T) {
	tests := []struct {
		name string
		give broker.Order
		want error
	}{
		{
			name: "err: invalid side",
			give: broker.Order{
				Side: 0,
				Type: broker.Market,
				Size: dec.New(1),
			},
			want: ErrInvalidOrderState,
		},
		{
			name: "err: invalid type",
			give: broker.Order{
				Side: broker.Buy,
				Type: 0,
				Size: dec.New(1),
			},
			want: ErrInvalidOrderState,
		},
		{
			name: "err: invalid size",
			give: broker.Order{
				Side: broker.Buy,
				Type: broker.Market,
				Size: dec.New(0),
			},
			want: ErrInvalidOrderState,
		},
		{
			name: "err: no pending state",
			give: broker.Order{
				OpenedAt: time.Now(),
			},
			want: ErrInvalidOrderState,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dealer Dealer
			act, _, err := dealer.PlaceOrder(context.Background(), tt.give)
			assert.ErrorIs(t, err, tt.want)
			assert.Nil(t, act)
		})
	}
}

func TestMatchOrder(t *testing.T) {
	tests := []struct {
		name      string
		giveOrder broker.Order
		giveQuote market.Kline
		want      decimal.Decimal
	}{
		{
			name: "ok: match limit",
			giveOrder: broker.Order{
				Type:       broker.Limit,
				LimitPrice: dec.New(12),
			},
			giveQuote: market.Kline{O: dec.New(8), H: dec.New(15), L: dec.New(5), C: dec.New(10)},
			want:      dec.New(12),
		},
		{
			name: "ok: match limit lower bound inclusive",
			giveOrder: broker.Order{
				Type:       broker.Limit,
				LimitPrice: dec.New(5),
			},
			giveQuote: market.Kline{O: dec.New(8), H: dec.New(15), L: dec.New(5), C: dec.New(10)},
			want:      dec.New(5),
		},
		{
			name: "ok: match limit upper bound inclusive",
			giveOrder: broker.Order{
				Type:       broker.Limit,
				LimitPrice: dec.New(15),
			},
			giveQuote: market.Kline{O: dec.New(8), H: dec.New(15), L: dec.New(5), C: dec.New(10)},
			want:      dec.New(15),
		},
		{
			name: "ok: no match limit below lower bound",
			giveOrder: broker.Order{
				Type:       broker.Limit,
				LimitPrice: dec.New(2),
			},
			giveQuote: market.Kline{O: dec.New(8), H: dec.New(15), L: dec.New(5), C: dec.New(10)},
			want:      decimal.Decimal{},
		},
		{
			name: "ok: no match limit above upper bound",
			giveOrder: broker.Order{
				Type:       broker.Limit,
				LimitPrice: dec.New(100),
			},
			giveQuote: market.Kline{O: dec.New(8), H: dec.New(15), L: dec.New(5), C: dec.New(10)},
			want:      decimal.Decimal{},
		},
		{
			name: "ok: always match market on close price",
			giveOrder: broker.Order{
				Type: broker.Market,
			},
			giveQuote: market.Kline{O: dec.New(8), H: dec.New(15), L: dec.New(5), C: dec.New(10)},
			want:      dec.New(10),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act := matchOrder(tt.giveOrder, tt.giveQuote)
			assert.Equal(t, tt.want, act)
		})
	}
}

func TestDealerOpenOrder(t *testing.T) {
	dealer := NewDealer()
	order := dealer.openOrder(broker.Order{})
	assert.EqualValues(t, broker.Open, order.State())
	assert.Contains(t, dealer.orders, order.ID)
}

func TestDealerFillOrder(t *testing.T) {
	dealer := NewDealer()
	order := dealer.fillOrder(broker.Order{}, dec.New(100))
	assert.EqualValues(t, broker.Filled, order.State())
}

func TestDealerCloseOrder(t *testing.T) {
	dealer := NewDealer()
	order := dealer.closeOrder(broker.Order{})
	assert.EqualValues(t, broker.Closed, order.State())
}

package backtest

import (
	"context"
	"testing"
	"time"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/dec"
	"github.com/colngroup/zero2algo/market"
	"github.com/stretchr/testify/assert"
)

func TestDealer_processOrder(t *testing.T) {

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
			var dealer Dealer
			act := dealer.processOrder(tt.give,
				time.Now().UTC(),
				market.Kline{O: dec.New(8), H: dec.New(15), L: dec.New(5), C: dec.New(10)})
			assert.Equal(t, tt.wantOrder.FilledSize, act.FilledSize)
			assert.Equal(t, tt.wantOrder.FilledPrice, act.FilledPrice)
			assert.Equal(t, tt.wantState, act.State())
		})
	}
}

func TestDealer_PlaceOrder(t *testing.T) {

	tests := []struct {
		name    string
		give    broker.Order
		wantErr error
	}{
		{
			name: "err: invalid side",
			give: broker.Order{
				Side: 0,
				Type: broker.Market,
				Size: dec.New(1),
			},
			wantErr: ErrInvalidOrderState,
		},
		{
			name: "err: invalid type",
			give: broker.Order{
				Side: broker.Buy,
				Type: 0,
				Size: dec.New(1),
			},
			wantErr: ErrInvalidOrderState,
		},
		{
			name: "err: invalid size",
			give: broker.Order{
				Side: broker.Buy,
				Type: broker.Market,
				Size: dec.New(0),
			},
			wantErr: ErrInvalidOrderState,
		},
		{
			name: "err: no pending state",
			give: broker.Order{
				OpenedAt: time.Now(),
			},
			wantErr: ErrInvalidOrderState,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dealer Dealer
			act, _, err := dealer.PlaceOrder(context.Background(), tt.give)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Nil(t, act)
		})
	}
}

func TestDealer_updateOpenOrders(t *testing.T) {
	dealer := NewDealer()

	// Add/update open orders if in open state
	dealer.updateOpenOrders(broker.Order{ID: "1", OpenedAt: time.Now()})
	assert.Contains(t, dealer.openOrders, broker.DealID("1"))

	// Delete open order if in closed state
	dealer.updateOpenOrders(broker.Order{ID: "1", ClosedAt: time.Now()})
	assert.NotContains(t, dealer.openOrders, broker.DealID("1"))
}

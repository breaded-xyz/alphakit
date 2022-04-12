package trend

import (
	"context"
	"testing"
	"time"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/dec"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/netapi"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestBot_executeExitOrder(t *testing.T) {
	tests := []struct {
		name      string
		giveSide  broker.OrderSide
		givePrice decimal.Decimal
		giveSize  decimal.Decimal
		want      broker.Order
	}{
		{
			name:      "exit ok",
			giveSide:  broker.Buy,
			givePrice: dec.New(10),
			giveSize:  dec.New(2),
			want: broker.Order{
				Side:       broker.Sell,
				Type:       broker.Market,
				Size:       dec.New(2),
				ReduceOnly: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dealer broker.MockDealer
			dealer.On("CancelOrders", context.Background()).Return((*netapi.Response)(nil), nil)
			dealer.On("PlaceOrder", context.Background(), tt.want).Return(&tt.want, (*netapi.Response)(nil), nil)
			bot := Bot{dealer: &dealer}
			act, err := bot.executeExitOrder(context.Background(), tt.giveSide, tt.giveSize)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, act)
			dealer.AssertExpectations(t)
		})
	}
}

func TestBot_executeEnterOrder(t *testing.T) {
	tests := []struct {
		name      string
		giveSide  broker.OrderSide
		givePrice decimal.Decimal
		giveSize  decimal.Decimal
		giveRisk  decimal.Decimal
		want      broker.BracketOrder
	}{
		{
			name:      "enter short",
			giveSide:  broker.Sell,
			givePrice: dec.New(10),
			giveSize:  dec.New(2),
			giveRisk:  dec.New(1),
			want: broker.BracketOrder{
				Enter: broker.Order{
					Side: broker.Sell,
					Type: broker.Market,
					Size: dec.New(2),
				},
				Stop: broker.Order{
					Side:       broker.Buy,
					Type:       broker.Limit,
					LimitPrice: dec.New(11),
					Size:       dec.New(2),
					ReduceOnly: true,
				},
			},
		},
		{
			name:      "enter long",
			giveSide:  broker.Buy,
			givePrice: dec.New(10),
			giveSize:  dec.New(2),
			giveRisk:  dec.New(1),
			want: broker.BracketOrder{
				Enter: broker.Order{
					Side: broker.Buy,
					Type: broker.Market,
					Size: dec.New(2),
				},
				Stop: broker.Order{
					Side:       broker.Sell,
					Type:       broker.Limit,
					LimitPrice: dec.New(9),
					Size:       dec.New(2),
					ReduceOnly: true,
				},
			},
		},
		{
			name:      "enter max risk no stop",
			giveSide:  broker.Buy,
			givePrice: dec.New(10),
			giveSize:  dec.New(2),
			giveRisk:  dec.New(10),
			want: broker.BracketOrder{
				Enter: broker.Order{
					Side: broker.Buy,
					Type: broker.Market,
					Size: dec.New(2),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dealer broker.MockDealer
			dealer.On("CancelOrders", context.Background()).Return((*netapi.Response)(nil), nil)
			dealer.On("PlaceOrder", context.Background(), tt.want.Enter).Return(&tt.want.Enter, (*netapi.Response)(nil), nil).Once()
			dealer.On("PlaceOrder", context.Background(), tt.want.Stop).Return(&tt.want.Stop, (*netapi.Response)(nil), nil).Maybe()
			bot := Bot{dealer: &dealer}
			act, err := bot.executeEnterOrder(context.Background(), tt.giveSide, tt.givePrice, tt.giveSize, tt.giveRisk)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, act)
			dealer.AssertExpectations(t)
		})
	}
}

func TestBot_filterPositions(t *testing.T) {
	fixed := time.Date(2022, time.January, 01, 0, 0, 0, 0, time.Local)

	tests := []struct {
		name          string
		givePositions []broker.Position
		giveAsset     market.Asset
		giveSide      broker.OrderSide
		giveState     broker.PositionState
		want          []broker.Position
	}{
		{
			name: "filtered",
			givePositions: []broker.Position{
				{
					OpenedAt: fixed,
					Asset:    market.Asset{Symbol: "BTCUSD"},
					Side:     broker.Buy,
				},
				{
					OpenedAt: fixed,
					ClosedAt: time.Now(),
					Asset:    market.Asset{Symbol: "BTCUSD"},
					Side:     broker.Buy,
				},
				{
					OpenedAt: fixed,
					Asset:    market.Asset{Symbol: "ETHUSD"},
					Side:     broker.Buy,
				},
				{
					OpenedAt: fixed,
					Asset:    market.Asset{Symbol: "BTCUSD"},
					Side:     broker.Sell,
				},
			},
			giveAsset: market.Asset{Symbol: "BTCUSD"},
			giveSide:  broker.Buy,
			giveState: broker.OrderOpen,
			want: []broker.Position{
				{
					OpenedAt: fixed,
					Asset:    market.Asset{Symbol: "BTCUSD"},
					Side:     broker.Buy,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act := filterPositions(tt.givePositions, tt.giveAsset, tt.giveSide, tt.giveState)
			assert.Equal(t, tt.want, act)
		})
	}
}

func TestBot_signal(t *testing.T) {
	tests := []struct {
		name           string
		giveEnterLong  float64
		giveEnterShort float64
		giveExitLong   float64
		giveExitShort  float64
		givePrediction float64
		wantEnter      broker.OrderSide
		wantExit       broker.OrderSide
	}{
		{
			name:           "default zero state",
			giveEnterLong:  1.0,
			giveEnterShort: -1.0,
			giveExitLong:   -0.7,
			giveExitShort:  0.7,
			givePrediction: 0,
			wantEnter:      0,
			wantExit:       0,
		},
		{
			name:           "flat",
			giveEnterLong:  1.0,
			giveEnterShort: -1.0,
			giveExitLong:   -0.7,
			giveExitShort:  0.7,
			givePrediction: 0.5,
			wantEnter:      0,
			wantExit:       0,
		},
		{
			name:           "go long",
			giveEnterLong:  1.0,
			giveEnterShort: -1.0,
			giveExitLong:   -0.7,
			giveExitShort:  0.7,
			givePrediction: 1.0,
			wantEnter:      broker.Buy,
			wantExit:       broker.Sell,
		},
		{
			name:           "go short",
			giveEnterLong:  1.0,
			giveEnterShort: -1.0,
			giveExitLong:   -0.7,
			giveExitShort:  0.7,
			givePrediction: -1.0,
			wantEnter:      broker.Sell,
			wantExit:       broker.Buy,
		},
		{
			name:           "exit long only",
			giveEnterLong:  1.0,
			giveEnterShort: -1.0,
			giveExitLong:   -0.7,
			giveExitShort:  0.7,
			givePrediction: -0.7,
			wantEnter:      0,
			wantExit:       broker.Buy,
		},
		{
			name:           "exit short only",
			giveEnterLong:  1.0,
			giveEnterShort: -1.0,
			giveExitLong:   -0.7,
			giveExitShort:  0.7,
			givePrediction: 0.7,
			wantEnter:      0,
			wantExit:       broker.Sell,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bot := Bot{
				EnterLong:  tt.giveEnterLong,
				ExitLong:   tt.giveExitLong,
				EnterShort: tt.giveEnterShort,
				ExitShort:  tt.giveExitShort,
			}
			actEnter, actExit := bot.signal(tt.givePrediction)
			assert.Equal(t, tt.wantEnter, actEnter)
			assert.Equal(t, tt.wantExit, actExit)
		})
	}
}

func TestBot_getOpenPosition(t *testing.T) {
	fixed := time.Date(2022, time.January, 01, 0, 0, 0, 0, time.Local)
	givePositions := []broker.Position{
		{ID: "1", OpenedAt: fixed, Side: broker.Buy},
		{ID: "2", OpenedAt: fixed, Side: broker.Sell},
		{ID: "3", OpenedAt: fixed, ClosedAt: time.Now(), Side: broker.Buy},
	}
	var dealer broker.MockDealer
	dealer.On("ListPositions", context.Background(), (*netapi.ListOpts)(nil)).Return(givePositions, (*netapi.Response)(nil), nil)

	giveSide := broker.Sell
	want := broker.Position{ID: "2", OpenedAt: fixed, Side: broker.Sell}

	bot := Bot{dealer: &dealer}
	act, _ := bot.getOpenPosition(context.Background(), giveSide)
	assert.Equal(t, act, want)
	dealer.AssertExpectations(t)
}

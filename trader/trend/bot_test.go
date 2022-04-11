package trend

import (
	"testing"
	"time"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/market"
	"github.com/stretchr/testify/assert"
)

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

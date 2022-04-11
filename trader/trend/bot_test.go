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
		name string
		give any
		want any
	}{
		{
			name: "",
			give: nil,
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var act any
			assert.Equal(t, tt.want, act)
		})
	}
}

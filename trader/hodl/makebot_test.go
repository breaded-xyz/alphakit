// Copyright 2022 The Coln Group Ltd
// SPDX-License-Identifier: MIT

package hodl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thecolngroup/alphakit/trader"
)

func TestMakeBotFromConfig(t *testing.T) {
	tests := []struct {
		name    string
		give    map[string]any
		wantBot trader.Bot
		wantErr error
	}{
		{
			name: "buy index < sell index",
			give: map[string]any{"buybarindex": 1, "sellbarindex": 1000},
			wantBot: &Bot{
				BuyBarIndex:  1,
				SellBarIndex: 1000,
			},
			wantErr: nil,
		},
		{
			name: "no sell",
			give: map[string]any{"buybarindex": 10, "sellbarindex": 0},
			wantBot: &Bot{
				BuyBarIndex:  10,
				SellBarIndex: 0,
			},
			wantErr: nil,
		},
		{
			name: "default",
			give: map[string]any{"buybarindex": 0, "sellbarindex": 0},
			wantBot: &Bot{
				BuyBarIndex:  0,
				SellBarIndex: 0,
			},
			wantErr: nil,
		},
		{
			name:    "buy index >= sell index",
			give:    map[string]any{"buybarindex": 10, "sellbarindex": 5},
			wantBot: nil,
			wantErr: trader.ErrInvalidConfig,
		},
		{
			name:    "not int",
			give:    map[string]any{"buybarindex": 10.5, "sellbarindex": 5},
			wantBot: nil,
			wantErr: trader.ErrInvalidConfig,
		},
		{
			name:    "neg int",
			give:    map[string]any{"buybarindex": -1, "sellbarindex": 5},
			wantBot: nil,
			wantErr: trader.ErrInvalidConfig,
		},
		{
			name:    "key not found",
			give:    map[string]any{"notakey": 10, "sellbarindex": 5},
			wantBot: nil,
			wantErr: trader.ErrInvalidConfig,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act, err := MakeBotFromConfig(tt.give)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantBot, act)
		})
	}
}

package perf

import (
	"testing"
	"time"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/internal/dec"
	"github.com/stretchr/testify/assert"
)

func TestTradeReport(t *testing.T) {

	give := []broker.Trade{
		{
			Side:       broker.Buy,
			Profit:     dec.New(-10),
			HoldPeriod: time.Hour * 96,
		},
		{
			Side:       broker.Buy,
			Profit:     dec.New(-20),
			HoldPeriod: time.Hour * 24,
		},
		{
			Side:       broker.Sell,
			Profit:     dec.New(100),
			HoldPeriod: time.Hour * 192,
		},
		{
			Side:       broker.Sell,
			Profit:     dec.New(10),
			HoldPeriod: time.Hour * 48,
		},
	}

	want := &TradeReport{
		TradeCount:           4,
		TotalNetProfit:       80,
		AvgNetProfit:         20,
		GrossProfit:          110,
		GrossLoss:            30,
		ProfitFactor:         3.6666666666666665,
		PRR:                  0.9594309793501488,
		PercentProfitable:    0.5,
		MaxProfit:            100,
		MaxLoss:              20,
		AvgProfit:            55,
		AvgLoss:              15,
		MaxLossStreak:        2,
		Kelly:                0.36363636363636365,
		OptimalF:             0.37,
		TotalTimeInMarketSec: 1296000,
		AvgHoldSec:           324000,
		StatN:                120,
		winningCount:         2,
		winningPct:           0.5,
		losingCount:          2,
		losingPct:            0.5,
	}

	act := NewTradeReport(give)
	assert.Equal(t, want, act)
}

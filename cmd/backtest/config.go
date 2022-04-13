package main

import (
	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/broker/backtest"
	"github.com/colngroup/zero2algo/optimize"
	"github.com/colngroup/zero2algo/trader"
	"github.com/colngroup/zero2algo/trader/trend"
)

var (
	_priceDir = "/Users/richklee/Dropbox/dev-share/github.com/thecolngroup/go-alpha/prices/binance/spot/BTCUSDT-H1/2021"

	_dealerMakeFunc = func() broker.SimulatedDealer {
		return backtest.NewDealer()
	}

	_botMakeFunc = func() trader.ConfigurableBot {
		return trend.NewBot()
	}

	_params = optimize.ParamRange{
		"asset":             {"btcusdt"},
		"initialCapital":    {1000.0},
		"spreadPct":         {0.0},
		"slippagePct":       {0.0},
		"transactionPct":    {0.0},
		"fundingHourPct":    {0.0},
		"enterLong":         {1.0},
		"enterShort":        {-1.0},
		"exitLong":          {-1.0, -0.9, -0.6},
		"exitShort":         {1.0, 0.9, 0.6},
		"maFastLength":      {1, 8, 16, 32, 64, 128},
		"maSlowLength":      {32, 64, 128, 256, 512},
		"maSDFilterLength":  {128, 256, 512},
		"maSDFilterFactor":  {1.0, 1.25, 1.5},
		"mmiLength":         {200, 300},
		"mmiSmootherLength": {100, 200, 300},
		"riskerSDLength":    {0},
		"riskerSDFactor":    {1.0},
		"sizerF":            {0.0},
		"sizerScaleF":       {1.0},
	}
)

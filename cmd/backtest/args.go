package main

import (
	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/broker/backtest"
	"github.com/colngroup/zero2algo/optimize"
	"github.com/colngroup/zero2algo/trader"
	"github.com/colngroup/zero2algo/trader/trend"
)

var (
	_outputDir = ".out"

	//_priceDir = "/Users/richklee/Dropbox/dev-share/github.com/thecolngroup/go-alpha/prices/binance/spot/LTCUSDT-H1/Y18192021"
	//_priceDir = "/Users/richklee/Dropbox/dev-share/github.com/thecolngroup/go-alpha/prices/binance/spot/SOLUSDT-H1/2021"
	//_priceDir = "/Users/richklee/Dropbox/dev-share/github.com/thecolngroup/go-alpha/prices/binance/spot/BTCUSDT-H1/2021"
	_priceDir = "/Users/richklee/Dropbox/dev-share/github.com/thecolngroup/go-alpha/prices/binance/spot/BTCUSDT-H1/Y17181920"
	//_priceDir = "/Users/richklee/Dropbox/dev-share/github.com/thecolngroup/go-alpha/prices/binance/spot/ETHUSDT-H1/Y18192021"

	_dealerMakeFunc = func() broker.SimulatedDealer {
		return backtest.NewDealer()
	}

	_botMakeFunc = func() trader.ConfigurableBot {
		return trend.NewBot()
	}

	/*_params = optimize.ParamRange{
		"inSample":          {false},
		"asset":             {"btcusdt"},
		"barSize":           {"H1"},
		"initialCapital":    {1000.0},
		"spreadPct":         {0.0005},
		"slippagePct":       {0.001},
		"transactionPct":    {0.0026},
		"fundingHourPct":    {0.000025},
		"enterLong":         {1.0},
		"enterShort":        {-1.0},
		"exitLong":          {-1},
		"exitShort":         {1.0},
		"maFastLength":      {16},
		"maSlowLength":      {32},
		"maSDFilterLength":  {128},
		"maSDFilterFactor":  {1.25},
		"mmiLength":         {300},
		"mmiSmootherLength": {300},
		"riskerSDLength":    {0.0},
		"riskerSDFactor":    {1.0},
		"sizerF":            {0.0},
		"sizerScaleF":       {1.0},
	}*/

	/*_params = optimize.ParamRange{
		"inSample":          {false},
		"asset":             {"solusdt"},
		"barSize":           {"H1"},
		"initialCapital":    {1000.0},
		"spreadPct":         {0.0005},
		"slippagePct":       {0.0005},
		"transactionPct":    {0.0005},
		"fundingHourPct":    {0.000025},
		"enterLong":         {1.0},
		"enterShort":        {-1.0},
		"exitLong":          {-0.9},
		"exitShort":         {0.6},
		"maFastLength":      {1},
		"maSlowLength":      {128, 256, 386, 512},
		"maSDFilterLength":  {512},
		"maSDFilterFactor":  {1.5},
		"mmiLength":         {200},
		"mmiSmootherLength": {200},
		"riskerSDLength":    {512},
		"riskerSDFactor":    {1.5},
		"sizerF":            {0.25},
		"sizerScaleF":       {0.5},
	}*/

	_params = optimize.ParamRange{
		"inSample":          {true},
		"asset":             {"btcusdt"},
		"barSize":           {"H1"},
		"initialCapital":    {1000.0},
		"spreadPct":         {0.0},
		"slippagePct":       {0.0},
		"transactionPct":    {0.0},
		"fundingHourPct":    {0.0},
		"enterLong":         {1.0},
		"enterShort":        {-1.0},
		"exitLong":          {-0.9},
		"exitShort":         {0.6},
		"maFastLength":      {1},
		"maSlowLength":      {128, 256, 512},
		"maSDFilterLength":  {128, 256, 512, 1024},
		"maSDFilterFactor":  {1.0, 1.25, 1.5, 1.75},
		"mmiLength":         {200, 300},
		"mmiSmootherLength": {200, 300},
		"riskerSDLength":    {0},
		"riskerSDFactor":    {1.0},
		"sizerF":            {0.0},
		"sizerScaleF":       {1.0},
	}
)

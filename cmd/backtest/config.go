package main

import (
	"github.com/colngroup/zero2algo/broker/backtest"
	"github.com/colngroup/zero2algo/optimize"
	"github.com/colngroup/zero2algo/trader/trend"
)

var _priceDir = ".prices/btcusdt-h1/"

var _params = optimize.ParamRange{
	"initialCapital": {1000},
	"spread":         {0.005},
}

var _dealer = *backtest.NewDealer()

var _bot = trend.Bot{}

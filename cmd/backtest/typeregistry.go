package main

import (
	"github.com/colngroup/zero2algo/trader"
	"github.com/colngroup/zero2algo/trader/trend"
)

type botMakerFunc func() trader.ConfigurableBot

var _typeRegistry = map[string]any{
	"trend.bot": botMakerFunc(func() trader.ConfigurableBot { return trend.NewBot() }),
}

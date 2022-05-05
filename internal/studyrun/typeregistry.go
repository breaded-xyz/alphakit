package studyrun

import (
	"github.com/colngroup/zero2algo/trader"
	"github.com/colngroup/zero2algo/trader/hodl"
	"github.com/colngroup/zero2algo/trader/trend"
)

var _typeRegistry = map[string]any{
	"hodl":        trader.MakeFromConfig(hodl.MakeBotFromConfig),
	"trend.cross": trader.MakeFromConfig(trend.MakeCrossBotFromConfig),
	"trend.apex":  trader.MakeFromConfig(trend.MakeApexBotFromConfig),
}

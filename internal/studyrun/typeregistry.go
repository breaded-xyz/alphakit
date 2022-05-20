package studyrun

import (
	"github.com/thecolngroup/alphakit/trader"
	"github.com/thecolngroup/alphakit/trader/hodl"
	"github.com/thecolngroup/alphakit/trader/trend"
)

var _typeRegistry = map[string]any{
	"hodl":        trader.MakeFromConfig(hodl.MakeBotFromConfig),
	"trend.cross": trader.MakeFromConfig(trend.MakeCrossBotFromConfig),
	"trend.apex":  trader.MakeFromConfig(trend.MakeApexBotFromConfig),
}

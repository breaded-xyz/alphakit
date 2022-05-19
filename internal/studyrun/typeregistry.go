package studyrun

import (
	"github.com/thecolngroup/zerotoalgo/trader"
	"github.com/thecolngroup/zerotoalgo/trader/hodl"
	"github.com/thecolngroup/zerotoalgo/trader/trend"
)

var _typeRegistry = map[string]any{
	"hodl":        trader.MakeFromConfig(hodl.MakeBotFromConfig),
	"trend.cross": trader.MakeFromConfig(trend.MakeCrossBotFromConfig),
	"trend.apex":  trader.MakeFromConfig(trend.MakeApexBotFromConfig),
}

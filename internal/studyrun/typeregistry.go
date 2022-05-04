package studyrun

import (
	"github.com/colngroup/zero2algo/trader"
	"github.com/colngroup/zero2algo/trader/trend"
)

var _typeRegistry = map[string]any{
	"trend.cross": trader.MakeFromConfig(trend.MakeCrossBotFromConfig),
	"trend.apex":  trader.MakeFromConfig(trend.MakeApexBotFromConfig),
}

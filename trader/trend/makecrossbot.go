package trend

import (
	"github.com/thecolngroup/alphakit/market"
	"github.com/thecolngroup/alphakit/money"
	"github.com/thecolngroup/alphakit/risk"
	"github.com/thecolngroup/alphakit/ta"
	"github.com/thecolngroup/alphakit/trader"
	"github.com/thecolngroup/dec"
	"github.com/thecolngroup/util"
)

var _ trader.MakeFromConfig = MakeCrossBotFromConfig

// MakeCrossBotFromConfig returns a bot configured with a CrossPredicter and sensible defaults if config is missing.
func MakeCrossBotFromConfig(config map[string]any) (trader.Bot, error) {

	var bot Bot

	bot.Asset = market.NewAsset(util.ToString(config["asset"]))

	bot.EnterLong = util.NNZ(util.ToFloat(config["enterlong"]), 1.0)
	bot.EnterShort = util.NNZ(util.ToFloat(config["entershort"]), -1.0)
	bot.ExitLong = util.NNZ(util.ToFloat(config["exitlong"]), -1.0)
	bot.ExitShort = util.NNZ(util.ToFloat(config["exitshort"]), 1.0)

	maFastLength := util.ToInt(config["mafastlength"])
	maSlowLength := util.ToInt(config["maslowlength"])
	if maFastLength >= maSlowLength {
		return nil, trader.ErrInvalidConfig
	}
	maOsc := ta.NewOsc(ta.NewALMA(maFastLength), ta.NewALMA(maSlowLength))
	mmi := ta.NewMMI(util.ToInt(config["mmilength"]))
	bot.Predicter = NewCrossPredicter(maOsc, mmi)

	riskSDLength := util.ToInt(config["riskersdlength"])
	if riskSDLength > 0 {
		bot.Risker = risk.NewSDRisker(riskSDLength, util.ToFloat(config["riskersdfactor"]))
	} else {
		bot.Risker = risk.NewFullRisker()
	}

	initialCapital := dec.New(util.NNZ(util.ToFloat(config["initialcapital"]), 1000))
	sizerF := util.ToFloat(config["sizerf"])
	if sizerF > 0 {
		bot.Sizer = money.NewSafeFSizer(initialCapital, sizerF, util.ToFloat(config["sizerscalef"]))
	} else {
		bot.Sizer = money.NewFixedSizer(initialCapital)
	}

	return &bot, nil
}

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

var _ trader.MakeFromConfig = MakeApexBotFromConfig

// MakeApexBotFromConfig returns a bot configured with an ApexPredicter.
func MakeApexBotFromConfig(config map[string]any) (trader.Bot, error) {

	var bot Bot

	bot.Asset = market.NewAsset(util.ToString(config["asset"]))

	bot.EnterLong = util.NNZ(util.ToFloat(config["enterlong"]), 1.0)
	bot.EnterShort = util.NNZ(util.ToFloat(config["entershort"]), -1.0)
	bot.ExitLong = util.NNZ(util.ToFloat(config["exitlong"]), -1.0)
	bot.ExitShort = util.NNZ(util.ToFloat(config["exitshort"]), 1.0)

	maLength := util.ToInt(config["malength"])
	ma := ta.NewALMA(maLength)
	mmi := ta.NewMMI(util.ToInt(config["mmilength"]))
	predicter := NewApexPredicter(ma, mmi)
	predicter.ApexDelta = util.NNZ(util.ToFloat(config["apexdelta"]), 0.5)
	bot.Predicter = predicter

	riskSDLength := util.ToInt(config["riskersdlength"])
	if riskSDLength > 0 {
		bot.Risker = risk.NewSDRisker(riskSDLength, util.ToFloat(config["riskersdfactor"]))
	} else {
		bot.Risker = risk.NewFullRisker()
	}

	initialCapital := dec.New(util.ToFloat(config["initialcapital"]))
	sizerF := util.ToFloat(config["sizerf"])
	if sizerF > 0 {
		bot.Sizer = money.NewSafeFSizer(initialCapital, sizerF, util.ToFloat(config["sizerscalef"]))
	} else {
		bot.Sizer = money.NewFixedSizer(initialCapital)
	}

	return &bot, nil
}

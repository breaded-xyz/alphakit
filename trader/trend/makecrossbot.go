package trend

import (
	"github.com/thecolngroup/alphakit/market"
	"github.com/thecolngroup/alphakit/money"
	"github.com/thecolngroup/alphakit/risk"
	"github.com/thecolngroup/alphakit/ta"
	"github.com/thecolngroup/alphakit/trader"
	"github.com/thecolngroup/gou/conv"
	"github.com/thecolngroup/gou/dec"
	"github.com/thecolngroup/gou/num"
)

var _ trader.MakeFromConfig = MakeCrossBotFromConfig

// MakeCrossBotFromConfig returns a bot configured with a CrossPredicter and sensible defaults if config is missing.
func MakeCrossBotFromConfig(config map[string]any) (trader.Bot, error) {

	var bot Bot

	bot.Asset = market.NewAsset(conv.ToString(config["asset"]))

	bot.EnterLong = num.NNZ(conv.ToFloat(config["enterlong"]), 1.0)
	bot.EnterShort = num.NNZ(conv.ToFloat(config["entershort"]), -1.0)
	bot.ExitLong = num.NNZ(conv.ToFloat(config["exitlong"]), -1.0)
	bot.ExitShort = num.NNZ(conv.ToFloat(config["exitshort"]), 1.0)

	maFastLength := conv.ToInt(config["mafastlength"])
	maSlowLength := conv.ToInt(config["maslowlength"])
	if maFastLength >= maSlowLength {
		return nil, trader.ErrInvalidConfig
	}
	maOsc := ta.NewOsc(ta.NewALMA(maFastLength), ta.NewALMA(maSlowLength))
	mmi := ta.NewMMI(conv.ToInt(config["mmilength"]))
	bot.Predicter = NewCrossPredicter(maOsc, mmi)

	riskSDLength := conv.ToInt(config["riskersdlength"])
	if riskSDLength > 0 {
		bot.Risker = risk.NewSDRisker(riskSDLength, conv.ToFloat(config["riskersdfactor"]))
	} else {
		bot.Risker = risk.NewFullRisker()
	}

	initialCapital := dec.New(num.NNZ(conv.ToFloat(config["initialcapital"]), 1000))
	sizerF := conv.ToFloat(config["sizerf"])
	if sizerF > 0 {
		bot.Sizer = money.NewSafeFSizer(initialCapital, sizerF, conv.ToFloat(config["sizerscalef"]))
	} else {
		bot.Sizer = money.NewFixedSizer(initialCapital)
	}

	return &bot, nil
}

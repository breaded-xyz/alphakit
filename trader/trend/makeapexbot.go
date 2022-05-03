package trend

import (
	"github.com/colngroup/zero2algo/internal/dec"
	"github.com/colngroup/zero2algo/internal/util"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/money"
	"github.com/colngroup/zero2algo/risk"
	"github.com/colngroup/zero2algo/ta"
	"github.com/colngroup/zero2algo/trader"
)

var _ trader.MakeFromConfig = MakeApexBotFromConfig

func MakeApexBotFromConfig(config map[string]any) (trader.Bot, error) {

	var bot Bot

	bot.asset = market.NewAsset(util.ToString(config["asset"]))

	bot.EnterLong = config["enterlong"].(float64)
	bot.EnterShort = config["entershort"].(float64)
	bot.ExitLong = config["exitlong"].(float64)
	bot.ExitShort = config["exitshort"].(float64)

	maFastLength := util.ToInt(config["mafastlength"])
	maSlowLength := util.ToInt(config["maslowlength"])
	if maFastLength >= maSlowLength {
		return nil, trader.ErrInvalidConfig
	}
	maOsc := ta.NewOsc(ta.NewALMA(maFastLength), ta.NewALMA(maSlowLength))
	maSDFilter := ta.NewSD(util.ToInt(config["masdfilterlength"]))
	mmi := ta.NewMMIWithSmoother(util.ToInt(config["mmilength"]), ta.NewALMA(util.ToInt(config["mmismootherlength"])))
	bot.Predicter = NewApexPredicter(maOsc, maSDFilter, mmi)

	riskSDLength := util.ToInt(config["riskersdlength"])
	if riskSDLength > 0 {
		bot.Risker = risk.NewSDRisker(util.ToInt(config["riskersdlength"]), config["riskersdfactor"].(float64))
	} else {
		bot.Risker = risk.NewFullRisker()
	}

	initialCapital := dec.New(config["initialcapital"].(float64))
	sizerF := config["sizerf"].(float64)
	if sizerF > 0 {
		bot.Sizer = money.NewSafeFSizer(initialCapital, sizerF, config["sizerscalef"].(float64))
	} else {
		bot.Sizer = money.NewFixedSizer(initialCapital)
	}

	return &bot, nil
}

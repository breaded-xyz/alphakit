// Copyright 2022 The Coln Group Ltd
// SPDX-License-Identifier: MIT

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

var _ trader.MakeFromConfig = MakeApexBotFromConfig

// MakeApexBotFromConfig returns a bot configured with an ApexPredicter.
func MakeApexBotFromConfig(config map[string]any) (trader.Bot, error) {

	var bot Bot

	bot.Asset = market.NewAsset(conv.ToString(config["asset"]))

	bot.EnterLong = num.NNZ(conv.ToFloat(config["enterlong"]), 1.0)
	bot.EnterShort = num.NNZ(conv.ToFloat(config["entershort"]), -1.0)
	bot.ExitLong = num.NNZ(conv.ToFloat(config["exitlong"]), -1.0)
	bot.ExitShort = num.NNZ(conv.ToFloat(config["exitshort"]), 1.0)

	maLength := conv.ToInt(config["malength"])
	ma := ta.NewALMA(maLength)
	mmi := ta.NewMMI(conv.ToInt(config["mmilength"]))
	predicter := NewApexPredicter(ma, mmi)
	predicter.ApexDelta = num.NNZ(conv.ToFloat(config["apexdelta"]), 0.5)
	bot.Predicter = predicter

	riskSDLength := conv.ToInt(config["riskersdlength"])
	if riskSDLength > 0 {
		bot.Risker = risk.NewSDRisker(riskSDLength, conv.ToFloat(config["riskersdfactor"]))
	} else {
		bot.Risker = risk.NewFullRisker()
	}

	initialCapital := dec.New(conv.ToFloat(config["initialcapital"]))
	sizerF := conv.ToFloat(config["sizerf"])
	if sizerF > 0 {
		bot.Sizer = money.NewSafeFSizer(initialCapital, sizerF, conv.ToFloat(config["sizerscalef"]))
	} else {
		bot.Sizer = money.NewFixedSizer(initialCapital)
	}

	return &bot, nil
}

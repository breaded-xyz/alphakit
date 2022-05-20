package hodl

import (
	"github.com/thecolngroup/alphakit/internal/util"
	"github.com/thecolngroup/alphakit/trader"
)

// MakeBotFromConfig builds a valid Bot from a given set of config params.
func MakeBotFromConfig(config map[string]any) (trader.Bot, error) {
	var hodl Bot

	if _, ok := config["buybarindex"]; !ok {
		return nil, trader.ErrInvalidConfig
	}
	buyBarIndex := util.ToInt(config["buybarindex"])

	if _, ok := config["sellbarindex"]; !ok {
		return nil, trader.ErrInvalidConfig
	}
	sellBarIndex := util.ToInt(config["sellbarindex"])

	switch {
	case buyBarIndex == 0 && sellBarIndex == 0:
		break
	case buyBarIndex >= 0 && sellBarIndex == 0:
		break
	case buyBarIndex < 0 || sellBarIndex < 0:
		return nil, trader.ErrInvalidConfig
	case buyBarIndex >= sellBarIndex:
		return nil, trader.ErrInvalidConfig
	}

	hodl.BuyBarIndex = buyBarIndex
	hodl.SellBarIndex = sellBarIndex

	return &hodl, nil
}

package hodl

import (
	"github.com/colngroup/zero2algo/internal/util"
	"github.com/colngroup/zero2algo/trader"
)

const (
	BuyBarIndexKey  = "buybarindex"
	SellBarIndexKey = "sellbarindex"
)

// MakeBotFromConfig builds a valid Bot from a given set of config params.
func MakeBotFromConfig(config map[string]any) (trader.Bot, error) {
	var hodl Bot

	if _, ok := config[BuyBarIndexKey]; !ok {
		return nil, trader.ErrInvalidConfig
	}
	buyBarIndex := util.ToInt(config[BuyBarIndexKey])

	if _, ok := config[SellBarIndexKey]; !ok {
		return nil, trader.ErrInvalidConfig
	}
	sellBarIndex := util.ToInt(config[SellBarIndexKey])

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

package hodl

import "github.com/colngroup/zero2algo/trader"

const (
	BuyBarIndex  = "buybarindex"
	SellBarIndex = "sellbarindex"
)

func MakeBot(config map[string]any) (trader.Bot, error) {
	var hodl Bot

	buyBarIndex, ok := config[BuyBarIndex].(int)
	if !ok {
		return nil, trader.ErrInvalidConfig
	}
	sellBarIndex, ok := config[SellBarIndex].(int)
	if !ok {
		return nil, trader.ErrInvalidConfig
	}

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

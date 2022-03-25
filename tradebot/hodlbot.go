package tradebot

import (
	"context"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/pricing"
)

var _ ConfigurableBot = (*HodlBot)(nil)

type HodlBot struct {
	// Algo parameters
	BuyBarIndex  int
	SellBarIndex int

	dealer broker.Dealer
}

func NewHodlBot(dealer broker.Dealer) *HodlBot {
	return &HodlBot{
		dealer: dealer,
	}
}

func (b *HodlBot) Configure(config map[string]any) error {
	buyBarIndex, ok := config["buybarindex"].(int)
	if !ok {
		return ErrInvalidConfig
	}
	sellBarIndex, ok := config["sellbarindex"].(int)
	if !ok {
		return ErrInvalidConfig
	}

	switch {
	case buyBarIndex == 0 && sellBarIndex == 0:
		break
	case buyBarIndex >= sellBarIndex:
		return ErrInvalidConfig
	}

	b.BuyBarIndex = buyBarIndex
	b.SellBarIndex = sellBarIndex

	return nil
}

func (b *HodlBot) ReceivePrice(ctx context.Context, price pricing.Kline) error {
	return nil
}

func (b *HodlBot) Close() error {
	return nil
}

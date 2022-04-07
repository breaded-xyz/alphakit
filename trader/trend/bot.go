package trend

import (
	"context"

	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/trader"
)

var _ trader.ConfigurableBot = (*Bot)(nil)

type Bot struct {
	predicter Predicter
}

func (b *Bot) Configure(config map[string]any) error {
	return nil
}

func (b *Bot) ReceivePrice(ctx context.Context, price market.Kline) error {
	return nil
}

func (b *Bot) Close(ctx context.Context) error {
	return nil
}

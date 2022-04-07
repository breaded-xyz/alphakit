package trend

import (
	"context"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/trader"
)

var _ trader.ConfigurableBot = (*Bot)(nil)

type Bot struct {
	EnterLong float64
	ExitLong  float64

	EnterShort float64
	ExitShort  float64

	dealer    broker.Dealer
	predicter Predicter
}

func (b *Bot) Configure(config map[string]any) error {
	return nil
}

func (b *Bot) ReceivePrice(ctx context.Context, price market.Kline) error {

	if err := b.predicter.ReceivePrice(ctx, price); err != nil {
		return err
	}

	if !b.predicter.Valid() {
		return nil
	}

	enterSide, exitSide := b.evalAlgo(b.predicter.Predict())

	if exitSide != 0 {
		if err := b.closePosition(ctx, price, exitSide); err != nil {
			return err
		}
	}

	if enterSide != 0 {
		if err := b.openPosition(ctx, price, enterSide); err != nil {
			return err
		}
	}

	return nil
}

func (b *Bot) evalAlgo(prediction float64) (enter, exit broker.OrderSide) {

	switch {
	case prediction == 0:
		return
	case prediction >= b.EnterLong:
		return broker.Buy, broker.Sell
	case prediction >= b.ExitShort:
		return 0, broker.Sell
	case prediction <= b.ExitLong:
		return 0, broker.Buy
	case prediction <= b.EnterShort:
		return broker.Sell, broker.Buy
	}

	return
}

func (b *Bot) closePosition(ctx context.Context, price market.Kline, side broker.OrderSide) error {
	return nil
}

func (b *Bot) openPosition(ctx context.Context, price market.Kline, side broker.OrderSide) error {
	return nil
}

func (b *Bot) Close(ctx context.Context) error {
	return nil
}

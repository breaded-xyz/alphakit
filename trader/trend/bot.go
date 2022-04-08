package trend

import (
	"context"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/dec"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/trader"
)

var _ trader.ConfigurableBot = (*Bot)(nil)

type Bot struct {
	EnterLong float64
	ExitLong  float64

	EnterShort float64
	ExitShort  float64

	asset      market.Asset
	positioner PositionManager
	predicter  Predicter
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

	enter, exit := b.evalAlgo(b.predicter.Predict())

	//size, b.sizePosition(price)

	if err := b.executeSignal(ctx, enter, exit); err != nil {
		return err
	}

	return nil
}

func (b *Bot) Close(ctx context.Context) error {
	return b.positioner.LiquidateAll()
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

func (b *Bot) executeSignal(ctx context.Context, enter, exit broker.OrderSide) error {

	switch {
	case enter == 0 && exit == 0:
		return nil
	case enter == broker.Buy:
		return b.positioner.EnterLong(ctx, dec.New(0))
	case enter == broker.Sell:
		return b.positioner.EnterShort()
	case exit == broker.Buy:
		return b.positioner.ExitLong()
	case exit == broker.Sell:
		return b.positioner.ExitShort()
	}

	return nil
}

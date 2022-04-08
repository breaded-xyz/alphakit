package trend

import (
	"context"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/dec"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/money"
	"github.com/colngroup/zero2algo/ta"
	"github.com/colngroup/zero2algo/trader"
	"github.com/shopspring/decimal"
)

var _ trader.ConfigurableBot = (*Bot)(nil)

type Bot struct {
	EnterLong float64
	ExitLong  float64

	EnterShort float64
	ExitShort  float64

	asset      market.Asset
	dealer     broker.Dealer
	predicter  Predicter
	sizer      money.Sizer
	risk       ta.Indicator
	openedSize decimal.Decimal
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
	if enter == 0 && exit == 0 {
		return nil
	}

	balance, _, err := b.dealer.GetBalance(ctx)
	if err != nil {
		return err
	}
	capital := balance.Trade
	risk := dec.New(b.risk.Value())

	if err := b.exit(ctx, exit, price.C, b.openedSize, risk); err != nil {
		return err
	}

	size := b.sizer.Size(price.C, capital, risk)
	if err := b.enter(ctx, enter, price.C, size, risk); err != nil {
		return err
	}
	b.openedSize = size

	return nil
}

func (b *Bot) Close(ctx context.Context) error {
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

func (b *Bot) enter(ctx context.Context, side broker.OrderSide, price, size, risk decimal.Decimal) error {
	if _, err := b.dealer.CancelOrders(ctx); err != nil {
		return err
	}

	order := broker.Order{
		Asset: b.asset,
		Type:  broker.Market,
		Side:  side,
		Size:  size,
	}
	if _, _, err := b.dealer.PlaceOrder(ctx, order); err != nil {
		return err
	}

	stop := broker.Order{
		Asset:      b.asset,
		Type:       broker.Limit,
		Side:       side.Opposite(),
		Size:       size,
		LimitPrice: price.Sub(risk),
		ReduceOnly: true,
	}
	if _, _, err := b.dealer.PlaceOrder(ctx, stop); err != nil {
		return err
	}

	return nil
}

func (b *Bot) exit(ctx context.Context, side broker.OrderSide, price, size, risk decimal.Decimal) error {
	if _, err := b.dealer.CancelOrders(ctx); err != nil {
		return err
	}

	order := broker.Order{
		Asset:      b.asset,
		Type:       broker.Market,
		Side:       side.Opposite(),
		Size:       size,
		ReduceOnly: true,
	}
	if _, _, err := b.dealer.PlaceOrder(ctx, order); err != nil {
		return err
	}
	return nil
}

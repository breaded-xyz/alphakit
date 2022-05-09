// Hodl package offers a buy and hold trading algo.
package hodl

import (
	"context"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/internal/dec"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/trader"
)

var _ trader.Bot = (*Bot)(nil)

// Bot implements a buy and hold algo.
// Should only be used for backtesting purposes.
// At the set bar index a long position of size of 1 is opened.
type Bot struct {
	// BuyBarIndex is the index in the price sequence to open the position.
	BuyBarIndex int

	// SellBarIndex is the index in the price sequence to close the position.
	SellBarIndex int

	asset    market.Asset
	dealer   broker.Dealer
	barIndex int
}

// New returns a new Bot.
func New(asset market.Asset, dealer broker.Dealer) *Bot {
	return &Bot{
		asset:  asset,
		dealer: dealer,
	}
}

// SetDealer is the dealer to use for order execution.
// Should only be given a simulated dealer for backtesting.
func (b *Bot) SetDealer(dealer broker.Dealer) {
	b.dealer = dealer
}

// Warmup is not used.
func (b *Bot) Warmup(ctx context.Context, prices []market.Kline) error {
	return nil
}

// ReceivePrice updates the algo with the next market price.
func (b *Bot) ReceivePrice(ctx context.Context, price market.Kline) error {
	defer func() { b.barIndex++ }()

	signal := b.evalAlgo(b.barIndex, b.BuyBarIndex, b.SellBarIndex)
	if signal == 0 {
		return nil
	}

	order := broker.NewOrder(b.asset, signal, dec.New(1))
	if _, _, err := b.dealer.PlaceOrder(ctx, order); err != nil {
		return err
	}

	return nil
}

func (b *Bot) evalAlgo(index, buybar, sellbar int) broker.OrderSide {
	var signal broker.OrderSide

	switch {
	case index == buybar:
		signal = broker.Buy
	case sellbar == 0:
		break
	case index == sellbar:
		signal = broker.Sell
	}

	return signal
}

// Close closes any open position.
func (b *Bot) Close(ctx context.Context) error {
	order := broker.NewOrder(b.asset, broker.Sell, dec.New(1))
	order.ReduceOnly = true
	if _, _, err := b.dealer.PlaceOrder(ctx, order); err != nil {
		return err
	}
	return nil
}

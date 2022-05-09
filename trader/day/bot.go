package day

import (
	"context"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/trader"
)

var _ trader.Bot = (*Bot)(nil)

// Bot is a trader.Bot implementation for day trading.
type Bot struct {
}

// NewBot creates a new Bot instance.
func NewBot() *Bot {
	return &Bot{}
}

// SetDealer sets the broker used for placing orders.
func (b *Bot) SetDealer(dealer broker.Dealer) {

}

// Warmup seeds the Predicter and Risker with historical price data.
func (b *Bot) Warmup(ctx context.Context, prices []market.Kline) error {

	return nil
}

// ReceivePrice updates the algo with latest market price potentially triggering buy and/or sell orders.
func (b *Bot) ReceivePrice(ctx context.Context, price market.Kline) error {

	return nil
}

// Close exits all open positions at current market price.
func (b *Bot) Close(ctx context.Context) error {

	return nil
}

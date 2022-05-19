package trader

import (
	"context"

	"github.com/thecolngroup/zerotoalgo/broker"
	"github.com/thecolngroup/zerotoalgo/market"
)

var _ Bot = (*StubBot)(nil)

// StubBot is a testing double.
type StubBot struct {
}

// SetDealer not implemented.
func (b *StubBot) SetDealer(dealer broker.Dealer) {
}

// Warmup not implemented.
func (b *StubBot) Warmup(ctx context.Context, prices []market.Kline) error {
	return nil
}

// ReceivePrice not implemented.
func (b *StubBot) ReceivePrice(ctx context.Context, price market.Kline) error {
	return nil
}

// Close not implemented.
func (b *StubBot) Close(ctx context.Context) error {
	return nil
}

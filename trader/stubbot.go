package trader

import (
	"context"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/market"
)

type StubBot struct {
}

func (b *StubBot) SetDealer(dealer broker.Dealer) {

}

func (b *StubBot) Warmup(ctx context.Context, prices []market.Kline) error {
	return nil
}

func (b *StubBot) Configure(config map[string]any) error {
	return nil
}

func (b *StubBot) ReceivePrice(ctx context.Context, price market.Kline) error {
	return nil
}

func (b *StubBot) Close(ctx context.Context) error {
	return nil
}

package tradebot

import (
	"context"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/dec"
	"github.com/colngroup/zero2algo/market"
)

var _ ConfigurableBot = (*HodlBot)(nil)

type HodlBot struct {
	BuyBarIndex  int
	SellBarIndex int

	asset    market.Asset
	dealer   broker.Dealer
	barIndex int
}

func NewHodlBot(asset market.Asset, dealer broker.Dealer) *HodlBot {
	return &HodlBot{
		asset:  asset,
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
	case buyBarIndex >= 0 && sellBarIndex == 0:
		break
	case buyBarIndex < 0 || sellBarIndex < 0:
		return ErrInvalidConfig
	case buyBarIndex >= sellBarIndex:
		return ErrInvalidConfig
	}

	b.BuyBarIndex = buyBarIndex
	b.SellBarIndex = sellBarIndex

	return nil
}

func (b *HodlBot) ReceivePrice(ctx context.Context, price market.Kline) error {
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

func (b *HodlBot) evalAlgo(index, buybar, sellbar int) broker.OrderSide {
	var signal broker.OrderSide

	switch {
	case index == buybar:
		signal = broker.Buy
		break
	case sellbar == 0:
		break
	case index == sellbar:
		signal = broker.Sell
		break
	}

	return signal
}

func (b *HodlBot) Close(ctx context.Context) error {
	order := broker.NewOrder(b.asset, broker.Sell, dec.New(1))
	order.ReduceOnly = true
	if _, _, err := b.dealer.PlaceOrder(ctx, order); err != nil {
		return err
	}
	return nil
}

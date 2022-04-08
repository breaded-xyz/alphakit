package trend

import (
	"context"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/market"
	"github.com/shopspring/decimal"
)

// Positioner offers a facade over the Dealer interface and manages the position
// (market exposure) for a single asset on behalf of a Bot.
type Positioner struct {
	dealer   broker.Dealer
	asset    market.Asset
	openSize decimal.Decimal
}

// EnterLong opens a new long position and closes any opened short positions.
// Resting orders will cancelled pior to opening the new position.
func (p *Positioner) EnterLong(ctx context.Context, size decimal.Decimal) error {
	if err := p.enter(ctx, broker.Buy, size); err != nil {
		return err
	}
	p.openSize = p.openSize.Add(size)

	return nil
}

// EnterShort opens a new short position and closes any opened long positions.
// Resting orders will cancelled pior to opening the new position.
func (p *Positioner) EnterShort() error {
	return nil
}

// ExitLong closes any open long positions.
func (p *Positioner) ExitLong() error {
	return nil
}

// ExitShort closes any open short positions.
func (p *Positioner) ExitShort() error {
	return nil
}

// LiquidateAll closes any and all open positions.
func (p *Positioner) LiquidateAll() error {
	return nil
}

func (p *Positioner) enter(ctx context.Context, side broker.OrderSide, size decimal.Decimal) error {
	if _, err := p.dealer.CancelOrders(ctx); err != nil {
		return err
	}

	order := broker.Order{
		Asset: p.asset,
		Side:  side,
		Type:  broker.Market,
		Size:  size,
	}
	_, _, err := p.dealer.PlaceOrder(ctx, order)
	if err != nil {
		return err
	}

	return nil
}

func (p *Positioner) exit(ctx context.Context, side broker.OrderSide, size decimal.Decimal) error {
	if _, err := p.dealer.CancelOrders(ctx); err != nil {
		return err
	}

	order := broker.Order{
		Asset:      p.asset,
		Side:       side.Opposite(),
		Type:       broker.Market,
		Size:       size,
		ReduceOnly: true,
	}
	_, _, err := p.dealer.PlaceOrder(ctx, order)
	if err != nil {
		return err
	}

	return nil
}

package trend

import "github.com/colngroup/zero2algo/broker"

// Positioner offers a facade over the Dealer interface and manages the position
// (market exposure) for a single asset on behalf of a Bot.
type Positioner struct {
	dealer broker.Dealer
}

// EnterLong opens a new long position and closes any opened short positions.
// Resting orders will cancelled pior to opening the new position.
func (p *Positioner) EnterLong() error {
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

/*
func (b *Bot) closePosition(ctx context.Context, price market.Kline, side broker.OrderSide) error {

	position, err := b.getOpenedPosition(ctx, side)
	if err != nil {
		return err
	}
	if position == nil {
		return nil
	}

	b.dealer.CancelAllOrders()

	order := broker.Order{
		Asset:      b.asset,
		Side:       side.Opposite(),
		Type:       broker.Market,
		Size:       position.Size,
		ReduceOnly: true,
	}
	placedOrder, res, err := b.dealer.PlaceOrder(ctx, order)
	if err != nil {
		return err
	}
	spew.Dump(placedOrder, res, err)

	return nil
}

func (b *Bot) openPosition(ctx context.Context, price market.Kline, side broker.OrderSide) error {
	return nil
}

func (b *Bot) getOpenedPosition(ctx context.Context, side broker.OrderSide) (*broker.Position, error) {
	positions, _, err := b.dealer.ListPositions(ctx, nil)
	if err != nil {
		return nil, err
	}
	if positions = filter(positions, b.asset, side); len(positions) == 0 {
		return nil, nil
	}
	return &positions[0], nil
}

func filter(positions []broker.Position, asset market.Asset, side broker.OrderSide) []broker.Position {
	return nil
}
*/

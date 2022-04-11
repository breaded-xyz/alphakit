package broker

import (
	"github.com/colngroup/zero2algo/market"
)

func FilterPositions(positions []Position, asset market.Asset, side OrderSide, state PositionState) []Position {

	filtered := make([]Position, 0, len(positions))
	for i := range positions {
		p := positions[i]
		if p.Asset.Equal(asset) && p.Side == side && p.State() == state {
			filtered = append(filtered, p)
		}
	}

	return filtered
}

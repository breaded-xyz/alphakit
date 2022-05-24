package broker

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/thecolngroup/alphakit/market"
)

// PositionState represents the state of a position as it is processed by a dealer.
type PositionState int

const (
	// PositionPending represents a position that has not been processed by a dealer.
	PositionPending = iota

	// PositionOpen represents a position that has been opened by a dealer.
	PositionOpen

	// PositionClosed represents a position that has been closed by a dealer.
	PositionClosed
)

func (s PositionState) String() string {
	return [...]string{"Pending", "Open", "Closed"}[s]
}

// Position represents a position in a market for a given asset.
type Position struct {
	ID               DealID
	OpenedAt         time.Time
	ClosedAt         time.Time
	Asset            market.Asset
	Side             OrderSide
	Price            decimal.Decimal
	Size             decimal.Decimal
	LiquidationPrice decimal.Decimal
}

// State returns the state of the position based on the position timestamps.
func (p *Position) State() PositionState {
	switch {
	case !p.ClosedAt.IsZero():
		return PositionClosed
	case !p.OpenedAt.IsZero():
		return PositionOpen
	}
	return PositionPending
}

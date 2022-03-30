package broker

import (
	"time"

	"github.com/colngroup/zero2algo/market"
	"github.com/shopspring/decimal"
)

type PositionState int

const (
	PositionPending = iota
	PositionOpen
	PositionClosed
)

func (s PositionState) String() string {
	return [...]string{"Pending", "Open", "Closed"}[s]
}

type Position struct {
	OpenedAt time.Time
	ClosedAt time.Time
	Asset    market.Asset
	Side     OrderSide
	Price    decimal.Decimal
	Size     decimal.Decimal
}

func (p *Position) State() PositionState {
	switch {
	case !p.ClosedAt.IsZero():
		return PositionClosed
	case !p.OpenedAt.IsZero():
		return PositionOpen
	}
	return PositionPending
}

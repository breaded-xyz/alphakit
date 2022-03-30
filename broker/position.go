package broker

import (
	"time"

	"github.com/colngroup/zero2algo/market"
	"github.com/shopspring/decimal"
)

type PositionState int

const (
	PositionOpen = iota + 1
	PositionClosed
)

func (s PositionState) String() string {
	return [...]string{"Open", "Closed"}[s]
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
		return PositionClosed
	}
	return 0
}

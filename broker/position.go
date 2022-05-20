package broker

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/thecolngroup/alphakit/market"
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
	ID               DealID
	OpenedAt         time.Time
	ClosedAt         time.Time
	Asset            market.Asset
	Side             OrderSide
	Price            decimal.Decimal
	Size             decimal.Decimal
	LiquidationPrice decimal.Decimal
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

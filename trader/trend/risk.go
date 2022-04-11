package trend

import (
	"github.com/colngroup/zero2algo/market"
	"github.com/shopspring/decimal"
)

type Risker interface {
	market.Receiver
	Risk() decimal.Decimal
	Valid() bool
}

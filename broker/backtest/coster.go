package backtest

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/thecolngroup/alphakit/broker"
)

// Coster is a cost model used by a dealer to apply trading charges and fees.
type Coster interface {
	Slippage(price decimal.Decimal) decimal.Decimal
	Spread(price decimal.Decimal) decimal.Decimal
	Transaction(broker.Order) decimal.Decimal
	Funding(position broker.Position, price decimal.Decimal, elapsed time.Duration) decimal.Decimal
}

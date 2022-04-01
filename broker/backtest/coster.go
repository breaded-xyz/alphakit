package backtest

import (
	"time"

	"github.com/colngroup/zero2algo/broker"
	"github.com/shopspring/decimal"
)

type Coster interface {
	Slippage(price decimal.Decimal) decimal.Decimal
	Spread(price decimal.Decimal) decimal.Decimal
	Transaction(order broker.Order) decimal.Decimal
	Funding(position broker.Position, elapsed time.Duration) decimal.Decimal
}

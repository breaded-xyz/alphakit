package backtest

import (
	"github.com/colngroup/zero2algo/broker"
	"github.com/shopspring/decimal"
)

func Profit(position broker.Position, price decimal.Decimal) decimal.Decimal {
	profit := price.Sub(position.Price).Mul(position.Size)
	if position.Side == broker.Sell {
		profit = profit.Neg()
	}
	return profit
}

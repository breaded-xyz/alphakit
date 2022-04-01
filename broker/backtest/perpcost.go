package backtest

import (
	"time"

	"github.com/colngroup/zero2algo/broker"
	"github.com/shopspring/decimal"
)

var _ Coster = (*PerpCost)(nil)

type PerpCost struct {
	SpreadPct      decimal.Decimal
	SlippagePct    decimal.Decimal
	TransactionPct decimal.Decimal
	FundingHourPct decimal.Decimal

	lastestFundingHour int
}

func NewPerpCost() *PerpCost {
	return &PerpCost{}
}

func (c *PerpCost) Slippage(price decimal.Decimal) decimal.Decimal {
	return price.Mul(c.SlippagePct)
}

func (c *PerpCost) Spread(price decimal.Decimal) decimal.Decimal {
	return price.Mul(c.SpreadPct)
}

func (c *PerpCost) Transaction(order broker.Order) decimal.Decimal {
	return order.FilledPrice.Mul(order.FilledSize).Mul(c.TransactionPct)
}

func (c *PerpCost) Funding(position broker.Position, elapsed time.Duration) decimal.Decimal {
	return decimal.Zero
}

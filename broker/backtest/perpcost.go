package backtest

import (
	"math"
	"time"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/internal/dec"
	"github.com/shopspring/decimal"
)

var _ Coster = (*PerpCost)(nil)

type PerpCost struct {
	SpreadPct      decimal.Decimal
	SlippagePct    decimal.Decimal
	TransactionPct decimal.Decimal
	FundingHourPct decimal.Decimal

	lastFundingHour float64
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

func (c *PerpCost) Funding(position broker.Position, price decimal.Decimal, elapsed time.Duration) decimal.Decimal {

	if position.State() != broker.OrderOpen {
		return decimal.Zero
	}

	hours := math.Trunc(elapsed.Hours())
	excess := hours - c.lastFundingHour

	if excess == 0 {
		return decimal.Zero
	}

	c.lastFundingHour = hours
	perHourCost := position.Size.Mul(price).Mul(c.FundingHourPct)
	totalCost := perHourCost.Mul(dec.New(excess))

	return totalCost
}

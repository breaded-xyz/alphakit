// Copyright 2022 The Coln Group Ltd
// SPDX-License-Identifier: MIT

package backtest

import (
	"math"
	"time"

	"github.com/shopspring/decimal"
	"github.com/thecolngroup/alphakit/broker"
	"github.com/thecolngroup/gou/dec"
)

var _ Coster = (*PerpCoster)(nil)

// PerpCoster implements the Coster interface for Perpetual Future style assets.
type PerpCoster struct {
	SpreadPct      decimal.Decimal
	SlippagePct    decimal.Decimal
	TransactionPct decimal.Decimal
	FundingHourPct decimal.Decimal

	lastFundingHour float64
}

// NewPerpCoster creates a new PerpCoster.
func NewPerpCoster() *PerpCoster {
	return &PerpCoster{}
}

// Slippage returns the cost of slippage incurred by an order.
// Slippage is a fraction of the order price.
func (c *PerpCoster) Slippage(price decimal.Decimal) decimal.Decimal {
	return price.Mul(c.SlippagePct)
}

// Spread returns the cost of the spread incurred by an order.
// Half the SpreadPct field is used, representing the difference from the mid-price to your quote.
func (c *PerpCoster) Spread(price decimal.Decimal) decimal.Decimal {
	if !c.SpreadPct.IsPositive() {
		return decimal.Zero
	}
	return price.Mul(c.SpreadPct.Div(dec.New(2)))
}

// Transaction returns the cost of a transaction, calculated as a fraction of the order price and size.
func (c *PerpCoster) Transaction(order broker.Order) decimal.Decimal {
	return order.FilledPrice.Mul(order.FilledSize).Mul(c.TransactionPct)
}

// Funding returns the funding fee for a position, calculated on an hourly basis.
func (c *PerpCoster) Funding(position broker.Position, price decimal.Decimal, elapsed time.Duration) decimal.Decimal {

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

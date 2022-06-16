// Copyright 2022 The Coln Group Ltd
// SPDX-License-Identifier: MIT

package broker

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/thecolngroup/alphakit/market"
)

// RoundTurn is the result of opening and closing a position aka round-trip.
type RoundTurn struct {
	ID         DealID          `csv:"id"`
	CreatedAt  time.Time       `csv:"created_at"`
	Asset      market.Asset    `csv:",inline"`
	Side       OrderSide       `csv:"side"`
	Profit     decimal.Decimal `csv:"profit"`
	HoldPeriod time.Duration   `csv:"hold_period"`
	TradeCount int             `csv:"trade_count"`
}

// Copyright 2022 The Coln Group Ltd
// SPDX-License-Identifier: MIT

// Package trader provides an API for building trading bots.
// A bot receives prices and execute orders with a broker.
// Child packages offer specific bot implementations.
package trader

import (
	"context"
	"errors"

	"github.com/thecolngroup/alphakit/broker"
	"github.com/thecolngroup/alphakit/market"
)

// Bot is the primary interface for a trading algo.
type Bot interface {
	// Warmup the indicators used by the bot with historical data prior to active trading.
	// The amount of price data required is typically equivalent to the longest lookback.
	Warmup(context.Context, []market.Kline) error

	// Sets the dealer to be used for order execution.
	SetDealer(broker.Dealer)

	// Sets the asset to trade.
	SetAsset(market.Asset)

	// Receive gives the bot the next market price and evaluates the algo,
	// potentially generating new broker orders.
	market.Receiver

	// Clean-up the bot before close down, e.g. close open positions.
	Close(context.Context) error
}

// ErrInvalidConfig is returned by MakeFromConfig.
var ErrInvalidConfig = errors.New("invalid bot config")

// MakeFromConfig is a factory for building a tailored bot from a given config.
// Used by the optimize package to mint new bots for backtesting.
type MakeFromConfig func(config map[string]any) (Bot, error)

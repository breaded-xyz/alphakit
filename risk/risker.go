// Copyright 2022 The Coln Group Ltd
// SPDX-License-Identifier: MIT

// Package risk offers an API to evaluate trade risk.
// Used in conjunction with the money package to size positions.
package risk

import (
	"github.com/shopspring/decimal"
	"github.com/thecolngroup/alphakit/market"
)

// Risker is an interface that defines the methods needed to evaluate trade risk.
type Risker interface {
	// ReceivePrice updates the risker with the next price.
	market.Receiver

	// Risk returns a unitary measure of risk based on the current price.
	Risk() decimal.Decimal

	// Valid returns true if the risker has enough data to be calculated.
	Valid() bool
}

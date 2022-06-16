// Copyright 2022 The Coln Group Ltd
// SPDX-License-Identifier: MIT

package market

// Asset represents a tradeable asset identified by its symbol, e.g. BTCUSD.
type Asset struct {
	Symbol string `csv:"symbol"`
}

// NewAsset creates a new Asset with the given symbol.
func NewAsset(symbol string) Asset {
	return Asset{
		Symbol: symbol,
	}
}

// Equal asserts equality based on the symbol.
func (a *Asset) Equal(other Asset) bool {
	return a.Symbol == other.Symbol
}

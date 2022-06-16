// Copyright 2022 The Coln Group Ltd
// SPDX-License-Identifier: MIT

package market

// KlineReader is an interface for reading candlesticks.
type KlineReader interface {
	Read() (Kline, error)
	ReadAll() ([]Kline, error)
}

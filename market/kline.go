// Package market provides an API to read and process market price data
package market

import (
	"time"

	"github.com/shopspring/decimal"
)

// Kline represents a single candlestick.
type Kline struct {
	Start  time.Time
	O      decimal.Decimal
	H      decimal.Decimal
	L      decimal.Decimal
	C      decimal.Decimal
	Volume float64
}

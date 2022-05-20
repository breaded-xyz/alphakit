package ta

import (
	"github.com/shopspring/decimal"
	"github.com/thecolngroup/alphakit/market"
)

// PriceSelector is a selector that returns a price value for the given kline.
type PriceSelector func(price market.Kline) float64

// HL2 returns the average of the high and low prices.
func HL2(price market.Kline) float64 {
	return decimal.Avg(price.H, price.L).InexactFloat64()
}

// HLC3 returns the average of the high, low and close prices.
func HLC3(price market.Kline) float64 {
	return decimal.Avg(price.H, price.L, price.C).InexactFloat64()
}

// OHLC4 returns the average of the open, high, low and close prices.
func OHLC4(price market.Kline) float64 {
	return decimal.Avg(price.O, price.H, price.L, price.C).InexactFloat64()
}

// Close returns the close price.
func Close(price market.Kline) float64 {
	return price.C.InexactFloat64()
}

package ta

import (
	"github.com/colngroup/zero2algo/market"
	"github.com/shopspring/decimal"
)

// PriceSelector is a selector that returns a price value for the given kline.
type PriceSelector func(price market.Kline) float64

// HL2 returns the average of the high and low prices.
func HL2(price market.Kline) float64 {
	return decimal.Avg(price.H, price.L).InexactFloat64()
}

// Close returns the close price.
func Close(price market.Kline) float64 {
	return price.C.InexactFloat64()
}

package ta

import (
	"github.com/colngroup/zero2algo/market"
	"github.com/shopspring/decimal"
)

type PriceSelector func(price market.Kline) float64

func HL2(price market.Kline) float64 {
	return decimal.Avg(price.H, price.L).InexactFloat64()
}

func Close(price market.Kline) float64 {
	return price.C.InexactFloat64()
}

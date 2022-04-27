package trader

import "github.com/colngroup/zero2algo/market"

type Predicter interface {
	market.Receiver
	Predict() float64
	Valid() bool
}

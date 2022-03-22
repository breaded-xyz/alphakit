package zero2algo

import "time"

type Kline struct {
	Start  time.Time
	O      Price
	H      Price
	L      Price
	C      Price
	Volume float64
}
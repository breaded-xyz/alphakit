package zero2algo

import "time"

type Kline struct {
	Start  time.Time
	O      Money
	H      Money
	L      Money
	C      Money
	Volume float64
}

func ReadKlinesFromCSV(file string) ([]Kline, error) {
	return nil, nil
}

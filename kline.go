package zero2algo

import (
	"encoding/csv"
	"time"

	"github.com/shopspring/decimal"
)

type Kline struct {
	Start  time.Time
	O      decimal.Decimal
	H      decimal.Decimal
	L      decimal.Decimal
	C      decimal.Decimal
	Volume float64
}

type KlineReader interface {
	Read() (Kline, error)
}

func NewCSVKlineReader(csv *csv.Reader) KlineReader {
	return nil
}

/*func (r *KlineReader) Read() (*Kline, error) {
	return nil, nil
}

func (r *KlineReader) ReadAll() ([]Kline, error) {
	return nil, nil
}

func ReadKlinesFromCSV(file string) ([]Kline, error) {
	return nil, nil
}*/

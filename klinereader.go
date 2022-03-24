package zero2algo

import "encoding/csv"

var _ KlineReader = (*CSVKlineReader)(nil)

type KlineReader interface {
	Read() (Kline, error)
}

type CSVKlineReader struct {
	csv *csv.Reader
}

func NewCSVKlineReader(csv *csv.Reader) *CSVKlineReader {
	return &CSVKlineReader{
		csv: csv,
	}
}

func (r *CSVKlineReader) Read() (Kline, error) {
	return Kline{}, nil
}

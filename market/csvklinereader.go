package market

import (
	"encoding/csv"
	"errors"
	"io"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

var (
	ErrNotEnoughColumns    = errors.New("require 5 cols to parse kline")
	ErrInvalidTimeFormat   = errors.New("require col[1] start time in unix format")
	ErrInvalidPriceFormat  = errors.New("require col[2..5] OHLC prices in valid decimal format")
	ErrInvalidVolumeFormat = errors.New("require col[6] volume in valid float format")
)

var _ KlineReader = (*CSVKlineReader)(nil)

type CSVKlineReader struct {
	csv *csv.Reader
}

func NewCSVKlineReader(csv *csv.Reader) *CSVKlineReader {
	return &CSVKlineReader{
		csv: csv,
	}
}

func (r *CSVKlineReader) Read() (Kline, error) {
	var k Kline

	rec, err := r.csv.Read()
	if err != nil {
		return k, err
	}

	return newKlineFromCSVRecord(rec)
}

func (r *CSVKlineReader) ReadAll() ([]Kline, error) {
	var ks []Kline
	for {
		k, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		ks = append(ks, k)
	}

	return ks, nil
}

func newKlineFromCSVRecord(record []string) (Kline, error) {
	var k, empty Kline
	var err error

	if len(record) < 5 {
		return k, ErrNotEnoughColumns
	}

	msec, err := strconv.ParseInt(record[0], 10, 64)
	if err != nil {
		return empty, ErrInvalidTimeFormat
	}
	k.Start = time.UnixMilli(msec).UTC()

	if k.O, err = decimal.NewFromString(record[1]); err != nil {
		return empty, ErrInvalidPriceFormat
	}
	if k.H, err = decimal.NewFromString(record[2]); err != nil {
		return empty, ErrInvalidPriceFormat
	}
	if k.L, err = decimal.NewFromString(record[3]); err != nil {
		return empty, ErrInvalidPriceFormat
	}
	if k.C, err = decimal.NewFromString(record[4]); err != nil {
		return empty, ErrInvalidPriceFormat
	}

	if len(record) > 5 {
		if k.Volume, err = strconv.ParseFloat(record[5], 64); err != nil {
			return empty, ErrInvalidVolumeFormat
		}
	}

	return k, nil
}

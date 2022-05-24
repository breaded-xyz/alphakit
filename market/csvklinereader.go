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
	// ErrNotEnoughColumns is returned when the CSV price record does not have enough columns.
	ErrNotEnoughColumns = errors.New("require 5 cols to parse kline")

	// ErrInvalidTimeFormat is returned when the CSV price record does not have a valid time unix milli format.
	ErrInvalidTimeFormat = errors.New("require col[0] start time in unix millisecond format")

	// ErrInvalidPriceFormat is returned when the CSV price record does not prices in expected format.
	ErrInvalidPriceFormat = errors.New("require col[1..4] OHLC prices in valid decimal format")

	// ErrInvalidVolumeFormat is returned when the CSV price record does not have a valid volume format.
	ErrInvalidVolumeFormat = errors.New("require col[5] volume in valid float format")
)

var _ KlineReader = (*CSVKlineReader)(nil)

// CSVKlineReader is a KlineReader that reads from a CSV file.
type CSVKlineReader struct {
	csv *csv.Reader
}

// NewCSVKlineReader creates a new CSVKlineReader.
func NewCSVKlineReader(csv *csv.Reader) *CSVKlineReader {
	return &CSVKlineReader{
		csv: csv,
	}
}

// Read reads the next Kline from the underlying CSV data.
func (r *CSVKlineReader) Read() (Kline, error) {
	var k Kline

	rec, err := r.csv.Read()
	if err != nil {
		return k, err
	}

	return newKlineFromCSVRecord(rec)
}

// ReadAll reads all the Klines from the underlying CSV data.
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

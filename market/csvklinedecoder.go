package market

import (
	"encoding/csv"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

// MetaTraderTimeFormat is the time format expected by the MetaTrader decoder when cols [0] and [1] are used.
const MetaTraderTimeFormat = "02/01/2006 15:04"

var (
	// ErrNotEnoughColumns is returned when the CSV price record does not have enough columns.
	ErrNotEnoughColumns = errors.New("not enough columns")

	// ErrInvalidTimeFormat is returned when the CSV price record does not have a valid time unix milli format.
	ErrInvalidTimeFormat = errors.New("cannot parse time string")

	// ErrInvalidPriceFormat is returned when the CSV price record does not prices in expected format.
	ErrInvalidPriceFormat = errors.New("OHLC prices must be in valid decimal format")

	// ErrInvalidVolumeFormat is returned when the CSV price record does not have a valid volume format.
	ErrInvalidVolumeFormat = errors.New("volume must be in valid float format")
)

// CSVKlineDecoder is an extension point for CSVKlineReader to support custom file formats.
type CSVKlineDecoder func(record []string) (Kline, error)

// NewBinanceCSVKlineReader creates a new CSVKlineReader for Binance CSV files.
func NewBinanceCSVKlineReader(csv *csv.Reader) *CSVKlineReader {
	return &CSVKlineReader{
		csv:     csv,
		decoder: BinanceCSVKlineDecoder,
	}
}

// BinanceCSVKlineDecoder decodes a CSV record from Binance into a Kline.
func BinanceCSVKlineDecoder(record []string) (Kline, error) {
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

// NewMetaTraderCSVKlineReader creates a new CSVKlineReader for MetaTrader CSV files.
func NewMetaTraderCSVKlineReader(csv *csv.Reader) *CSVKlineReader {
	csv.Comma = ';'
	return &CSVKlineReader{
		csv:     csv,
		decoder: MetaTraderCSVKlineDecoder,
	}
}

// MetaTraderCSVKlineDecoder decodes a CSV record from MetaTrader into a Kline.
func MetaTraderCSVKlineDecoder(record []string) (Kline, error) {
	var k, empty Kline
	var err error

	if len(record) < 6 {
		return k, ErrNotEnoughColumns
	}

	tStr := fmt.Sprintf("%s %s", record[0], record[1])
	t, err := time.Parse(MetaTraderTimeFormat, tStr)
	if err != nil {
		return empty, ErrInvalidTimeFormat
	}
	k.Start = t.UTC()

	if k.O, err = decimal.NewFromString(record[2]); err != nil {
		return empty, ErrInvalidPriceFormat
	}
	if k.H, err = decimal.NewFromString(record[3]); err != nil {
		return empty, ErrInvalidPriceFormat
	}
	if k.L, err = decimal.NewFromString(record[4]); err != nil {
		return empty, ErrInvalidPriceFormat
	}
	if k.C, err = decimal.NewFromString(record[5]); err != nil {
		return empty, ErrInvalidPriceFormat
	}

	if len(record) > 5 {
		if k.Volume, err = strconv.ParseFloat(record[6], 64); err != nil {
			return empty, ErrInvalidVolumeFormat
		}
	}

	return k, nil
}

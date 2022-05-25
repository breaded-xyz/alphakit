package market

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadKlinesFromCSV(t *testing.T) {

	prices, err := ReadKlinesFromCSV("testdata/BTCUSDT-1h-2021-Q1.csv")
	assert.NoError(t, err)
	assert.Len(t, prices, 2158)
}

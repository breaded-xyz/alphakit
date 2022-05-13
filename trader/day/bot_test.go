package day

import (
	"context"
	"testing"

	"github.com/colngroup/zero2algo/internal/studyrun"
	"github.com/stretchr/testify/assert"
)

func TestBot(t *testing.T) {
	bot := NewBot()

	testdataPath := "/Users/richklee/Dropbox/dev-share/github.com/colngroup/zero2algo/prices/spot/btcusdt-m1/all"

	prices, err := studyrun.ReadPriceSeries(testdataPath)
	assert.NoError(t, err)

	for i := range prices {
		err := bot.ReceivePrice(context.Background(), prices[i])
		assert.NoError(t, err)
	}

	studyrun.SaveStructToCSV("./testdata/results-9.csv", bot.Results)
}

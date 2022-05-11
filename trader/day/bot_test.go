package day

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/colngroup/zero2algo/internal/studyrun"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestBot(t *testing.T) {
	bot := NewBot()

	prices, err := studyrun.ReadPriceSeries(filepath.Join(testdataPath, "day"))
	assert.NoError(t, err)

	for i := range prices {
		err := bot.ReceivePrice(context.Background(), prices[i])
		assert.NoError(t, err)

	}

	spew.Dump(bot.Results)
}

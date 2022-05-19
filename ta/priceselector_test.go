package ta

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thecolngroup/zerotoalgo/internal/dec"
	"github.com/thecolngroup/zerotoalgo/market"
)

func TestHL2(t *testing.T) {
	give := market.Kline{H: dec.New(10), L: dec.New(5)}
	want := 7.5
	act := HL2(give)
	assert.Equal(t, want, act)
}

func TestClose(t *testing.T) {
	give := market.Kline{C: dec.New(75)}
	want := 75.0
	act := Close(give)
	assert.Equal(t, want, act)
}

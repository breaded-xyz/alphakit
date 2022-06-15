package ta

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thecolngroup/alphakit/market"
	"github.com/thecolngroup/gou/dec"
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

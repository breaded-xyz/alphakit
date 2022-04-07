package ta

import (
	"testing"

	"github.com/colngroup/zero2algo/dec"
	"github.com/colngroup/zero2algo/market"
	"github.com/stretchr/testify/assert"
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

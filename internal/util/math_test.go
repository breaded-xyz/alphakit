package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundTo(t *testing.T) {
	giveX, giveY := 0.036, 0.01
	want := 0.3
	act := RoundTo(giveX, giveY)
	assert.Equal(t, want, act)
}

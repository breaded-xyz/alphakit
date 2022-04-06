package ta

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestALMA(t *testing.T) {
	giveV := []float64{10, 89, 20, 43, 44, 33, 19}
	giveN := 3
	want := 32.68047906324239

	ind := NewALMA(giveN)
	err := ind.Update(giveV...)
	act := ind.Value()

	assert.NoError(t, err)
	assert.True(t, ind.Valid())
	assert.Equal(t, want, act)
}

package ta

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMMI(t *testing.T) {

	giveV := make([]float64, 1000)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < len(giveV); i++ {
		giveV[i] = r.Float64()
	}
	giveLength := 300
	wantLowerLimit := 0.65
	wantUpperLimit := 0.85

	ind := NewMMI(giveLength)
	err := ind.Update(giveV...)
	act := ind.Value()

	assert.NoError(t, err)
	assert.True(t, ind.Valid())
	assert.Greater(t, act, wantLowerLimit)
	assert.Less(t, act, wantUpperLimit)

	t.Log(act)
}

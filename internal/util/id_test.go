package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

func TestID(t *testing.T) {

	ids := []ID{
		NewID(),
		NewIDWithTime(time.Now().Add(time.Hour * 2)),
		NewIDWithTime(time.Now().Add(time.Hour * 1)),
	}
	assert.False(t, slices.IsSorted(ids))

	slices.Sort(ids)
	assert.True(t, slices.IsSorted(ids))
}

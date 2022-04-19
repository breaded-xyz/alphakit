package util

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

type ID string

func NewID() ID {
	return NewIDWithTime(time.Now())
}

func NewIDWithTime(t time.Time) ID {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	return ID(ulid.MustNew(ulid.Timestamp(t), entropy).String())
}

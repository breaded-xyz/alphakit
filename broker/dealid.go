package broker

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

type DealID string

func NewID() DealID {
	return NewIDWithTime(time.Now().UTC())
}

func NewIDWithTime(t time.Time) DealID {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	return DealID(ulid.MustNew(ulid.Timestamp(t), entropy).String())
}

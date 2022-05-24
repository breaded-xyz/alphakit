package broker

import (
	"time"

	"github.com/thecolngroup/util"
)

// DealID is a unique identifier for a dealer data entity.
type DealID string

// NewID returns a new DealID seeded with the current time.
func NewID() DealID {
	return DealID(util.NewID())
}

// NewIDWithTime returns a new DealID seeded with the given time.
func NewIDWithTime(t time.Time) DealID {
	return DealID(util.NewIDWithTime(t))
}

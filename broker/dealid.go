// Copyright 2022 The Coln Group Ltd
// SPDX-License-Identifier: MIT

package broker

import (
	"time"

	"github.com/thecolngroup/gou/id"
)

// DealID is a unique identifier for a dealer data entity.
type DealID string

// NewID returns a new DealID seeded with the current time.
func NewID() DealID {
	return DealID(id.New())
}

// NewIDWithTime returns a new DealID seeded with the given time.
func NewIDWithTime(t time.Time) DealID {
	return DealID(id.NewWithTime(t))
}

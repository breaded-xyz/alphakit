// Copyright 2022 The Coln Group Ltd
// SPDX-License-Identifier: MIT

package perf

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thecolngroup/alphakit/broker"
	"github.com/thecolngroup/gou/dec"
)

func TestDiffPctReturns(t *testing.T) {
	give := broker.EquitySeries{
		1: dec.New(10),
		2: dec.New(20), // 20 - 10 = 10 / 10 = 1
		3: dec.New(30), // 30 - 20 = 10 / 20 = 0.5
		4: dec.New(5),  // 5 - 30 = -25 / 30 = -0.8333333333333334
	}
	want := []float64{1, 0.5, -0.8333333333333334}
	act := DiffPctReturns(give)
	assert.Equal(t, want, act)
}

func TestReduceEoD(t *testing.T) {

	datum := time.Date(2022, time.January, 1, 0, 0, 0, 0, time.Local)
	give := broker.EquitySeries{
		broker.Timestamp(datum.UnixMilli()):                     dec.New(10),
		broker.Timestamp(datum.Add(25 * time.Hour).UnixMilli()): dec.New(20),
		broker.Timestamp(datum.Add(48 * time.Hour).UnixMilli()): dec.New(30),
	}
	want := broker.EquitySeries{
		broker.Timestamp(datum.UnixMilli()):                     dec.New(10),
		broker.Timestamp(datum.Add(48 * time.Hour).UnixMilli()): dec.New(30),
	}

	act := ReduceEOD(give)
	assert.Equal(t, want, act)
}

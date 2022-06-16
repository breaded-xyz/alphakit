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

func TestDrawdowns(t *testing.T) {
	give := broker.EquitySeries{
		1:  dec.New(15), // Peak 1
		2:  dec.New(0),  // Valley 1
		3:  dec.New(15), // Drawdown 1 = 15
		4:  dec.New(20), // Peak 2
		5:  dec.New(19),
		6:  dec.New(10), // Valley 2
		7:  dec.New(11),
		8:  dec.New(17),
		9:  dec.New(30), // Drawdown 2 = 10 & Peak 3
		10: dec.New(25), // Valley 3
		11: dec.New(30), // Drawdown 3 = 5 & Peak 4
		12: dec.New(10),
		13: dec.New(5),  // Valley 4
		14: dec.New(25), // Drawdown 4 (status = open)
	}

	exp := []Drawdown{
		{
			HighAt:   time.UnixMilli(1),
			LowAt:    time.UnixMilli(2),
			StartAt:  time.UnixMilli(2),
			EndAt:    time.UnixMilli(3),
			High:     dec.New(15),
			Low:      dec.New(0),
			Recovery: 1 * time.Millisecond,
			Amount:   dec.New(15),
			Pct:      1,
		},
		{
			HighAt:   time.UnixMilli(4),
			LowAt:    time.UnixMilli(6),
			StartAt:  time.UnixMilli(5),
			EndAt:    time.UnixMilli(9),
			High:     dec.New(20),
			Low:      dec.New(10),
			Recovery: 4 * time.Millisecond,
			Amount:   dec.New(10),
			Pct:      0.5,
		},
		{
			HighAt:   time.UnixMilli(9),
			LowAt:    time.UnixMilli(10),
			StartAt:  time.UnixMilli(10),
			EndAt:    time.UnixMilli(11),
			High:     dec.New(30),
			Low:      dec.New(25),
			Recovery: 1 * time.Millisecond,
			Amount:   dec.New(5),
			Pct:      0.1666666666666667,
		},
		{
			HighAt:   time.UnixMilli(11),
			LowAt:    time.UnixMilli(13),
			StartAt:  time.UnixMilli(12),
			EndAt:    time.UnixMilli(14),
			High:     dec.New(30),
			Low:      dec.New(5),
			Recovery: 2 * time.Millisecond,
			Amount:   dec.New(25),
			Pct:      0.8333333333333333,
			IsOpen:   true,
		},
	}

	act := Drawdowns(give)
	assert.Equal(t, exp, act)
}

package perf

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thecolngroup/zerotoalgo/broker"
	"github.com/thecolngroup/zerotoalgo/internal/dec"
)

func TestPortfolioReport(t *testing.T) {
	datum := time.Date(2022, time.January, 1, 0, 0, 0, 0, time.Local)
	give := broker.EquitySeries{
		broker.Timestamp(datum.UnixMilli()):                     dec.New(10),
		broker.Timestamp(datum.Add(24 * time.Hour).UnixMilli()): dec.New(20),
		broker.Timestamp(datum.Add(96 * time.Hour).UnixMilli()): dec.New(40),
		broker.Timestamp(datum.Add(48 * time.Hour).UnixMilli()): dec.New(5),
		broker.Timestamp(datum.Add(72 * time.Hour).UnixMilli()): dec.New(10),
	}
	want := PortfolioReport{
		PeriodStart:  datum,
		PeriodEnd:    datum.Add(96 * time.Hour),
		Period:       96 * time.Hour,
		StartEquity:  10,
		EndEquity:    40,
		EquityReturn: 3,
	}
	act := NewPortfolioReport(give)
	assert.Equal(t, want.PeriodStart, act.PeriodStart)
	assert.Equal(t, want.PeriodEnd, act.PeriodEnd)
	assert.Equal(t, want.Period, act.Period)
	assert.Equal(t, want.StartEquity, act.StartEquity)
	assert.Equal(t, want.EndEquity, act.EndEquity)
	assert.Equal(t, want.EquityReturn, act.EquityReturn)
}

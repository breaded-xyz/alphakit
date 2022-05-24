package broker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimestamp_Time(t *testing.T) {

	var tstamp Timestamp
	exp := time.Date(2022, time.January, 1, 12, 45, 30, 0, time.Local)
	tstamp = Timestamp(exp.UnixMilli())

	act := tstamp.Time()
	assert.Equal(t, exp, act)
}

func TestTimeSeries_SortKeys(t *testing.T) {

	exp := []Timestamp{1, 2, 3, 4, 5}
	give := TimeSeries[int]{
		2: 2,
		5: 5,
		4: 4,
		1: 1,
		3: 3,
	}
	act := give.SortKeys()
	assert.Equal(t, exp, act)
}

func TestTimeSeries_SortValuesByTime(t *testing.T) {

	exp := []int{1, 2, 3, 4, 5}
	give := TimeSeries[int]{
		2: 2,
		5: 5,
		4: 4,
		1: 1,
		3: 3,
	}
	act := give.SortValuesByTime()
	assert.Equal(t, exp, act)
}

func TestSortedMapKeys(t *testing.T) {
	m := map[int]string{
		2: "2",
		0: "0",
		1: "1",
	}
	want := []string{"0", "1", "2"}
	act := SortMapValues(m)
	assert.Equal(t, want, act)
}

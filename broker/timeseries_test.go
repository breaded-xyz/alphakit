package broker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimestampTime(t *testing.T) {

	var tstamp Timestamp
	exp := time.Date(2022, time.January, 1, 12, 45, 30, 0, time.Local)
	tstamp = Timestamp(exp.UnixMilli())

	act := tstamp.Time()
	assert.Equal(t, exp, act)
}

func TestTimeSeriesSortKeys(t *testing.T) {

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

func TestTimeSeriesSortValuesByTime(t *testing.T) {

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

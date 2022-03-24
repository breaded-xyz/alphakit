package price

import (
	"encoding/csv"
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCSVKlineReader_Read(t *testing.T) {

	var assertKlineEq = func(t *testing.T, exp, act Kline) {
		assert.Equal(t, exp.Start, act.Start)
		assert.True(t, exp.O.Equal(act.O))
		assert.True(t, exp.H.Equal(act.H))
		assert.True(t, exp.L.Equal(act.L))
		assert.True(t, exp.C.Equal(act.C))
		assert.Equal(t, exp.Volume, act.Volume)
	}

	tests := []struct {
		name string
		give string
		want Kline
		err  error
	}{
		{
			name: "Read DOHLCV",
			give: "1609459200000,28923.63000000,29031.34000000,28690.17000000,28995.13000000,2311.81144500",
			want: Kline{
				Start:  time.UnixMilli(1609459200000).UTC(),
				O:      decimal.NewFromFloat(28923.63),
				H:      decimal.NewFromFloat(29031.34),
				L:      decimal.NewFromFloat(28690.17),
				C:      decimal.NewFromFloat(28995.13),
				Volume: 2311.81144500},
			err: nil,
		},
		{
			name: "Read DOHLC",
			give: "1609459200000,28923.63000000,29031.34000000,28690.17000000,28995.13000000",
			want: Kline{
				Start:  time.UnixMilli(1609459200000).UTC(),
				O:      decimal.NewFromFloat(28923.63),
				H:      decimal.NewFromFloat(29031.34),
				L:      decimal.NewFromFloat(28690.17),
				C:      decimal.NewFromFloat(28995.13),
				Volume: 0},
			err: nil,
		},
		{
			name: "Not enough columns",
			give: "1609459200000,28923.63000000,29031.34000000",
			want: Kline{},
			err:  ErrNotEnoughColumns,
		},
		{
			name: "Invalid time format",
			give: "23/12/2021,28923.63000000,29031.34000000,28690.17000000,28995.13000000",
			want: Kline{},
			err:  ErrInvalidTimeFormat,
		},
		{
			name: "Invalid price format",
			give: "1609459200000,sixty,29031.34000000,28690.17000000,28995.13000000",
			want: Kline{},
			err:  ErrInvalidPriceFormat,
		},
		{
			name: "Invalid volume format",
			give: "1609459200000,28923.63000000,29031.34000000,28690.17000000,28995.13000000,vol",
			want: Kline{},
			err:  ErrInvalidVolumeFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewCSVKlineReader(csv.NewReader(strings.NewReader(tt.give)))
			kline, err := reader.Read()
			assert.Equal(t, tt.err, err)
			assertKlineEq(t, tt.want, kline)
		})
	}
}

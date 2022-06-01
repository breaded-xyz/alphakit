package market

import (
	"encoding/csv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thecolngroup/dec"
)

var assertKlineEq = func(t *testing.T, exp, act Kline) {
	assert.Equal(t, exp.Start, act.Start)
	assert.True(t, exp.O.Equal(act.O))
	assert.True(t, exp.H.Equal(act.H))
	assert.True(t, exp.L.Equal(act.L))
	assert.True(t, exp.C.Equal(act.C))
	assert.Equal(t, exp.Volume, act.Volume)
}

func TestCSVKlineReader_ReadWithBinanceDecoder(t *testing.T) {
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
				O:      dec.New(28923.63),
				H:      dec.New(29031.34),
				L:      dec.New(28690.17),
				C:      dec.New(28995.13),
				Volume: 2311.81144500},
			err: nil,
		},
		{
			name: "Read DOHLC",
			give: "1609459200000,28923.63000000,29031.34000000,28690.17000000,28995.13000000",
			want: Kline{
				Start:  time.UnixMilli(1609459200000).UTC(),
				O:      dec.New(28923.63),
				H:      dec.New(29031.34),
				L:      dec.New(28690.17),
				C:      dec.New(28995.13),
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
			reader := NewBinanceCSVKlineReader(csv.NewReader(strings.NewReader(tt.give)))
			kline, err := reader.Read()
			assert.Equal(t, tt.err, err)
			assertKlineEq(t, tt.want, kline)
		})
	}
}

func TestCSVKlineReader_ReadAllWithDefaultDecoder(t *testing.T) {
	records := []string{
		"1609459200000,28923.63000000,29031.34000000,28690.17000000,28995.13000000,2311.81144500",
		"1609459300000,28928.63000000,30031.34000000,22690.17000000,28495.13000000,3000.00",
	}
	reader := NewCSVKlineReader(csv.NewReader(strings.NewReader(strings.Join(records, "\n"))))
	klines, err := reader.ReadAll()
	assert.NoError(t, err)
	assert.Len(t, klines, 2)
}

func TestCSVKlineReader_ReadWithMetaTraderDecoder(t *testing.T) {

	tests := []struct {
		name string
		give string
		want Kline
		err  error
	}{
		{
			name: "Read DOHLCV",
			give: "11/12/2008;16:00;779.527679;780.964756;777.527679;779.964756;5",
			want: Kline{
				Start:  time.Date(2008, 12, 11, 16, 0, 0, 0, time.UTC),
				O:      dec.New(779.527679),
				H:      dec.New(780.964756),
				L:      dec.New(777.527679),
				C:      dec.New(779.964756),
				Volume: 5},
			err: nil,
		},
		{
			name: "Not enough columns",
			give: "1609459200000;28923.63000000;29031.34000000",
			want: Kline{},
			err:  ErrNotEnoughColumns,
		},
		{
			name: "Invalid time format",
			give: "23/12/2021;t;28923.63000000;29031.34000000;28690.17000000;28995.13000000",
			want: Kline{},
			err:  ErrInvalidTimeFormat,
		},
		{
			name: "Invalid price format",
			give: "11/12/2008;00:00;sixty;29031.34000000;28690.17000000;28995.13000000",
			want: Kline{},
			err:  ErrInvalidPriceFormat,
		},
		{
			name: "Invalid volume format",
			give: "11/12/2008;00:00;779.527679;780.964756;777.527679;779.964756;vol",
			want: Kline{},
			err:  ErrInvalidVolumeFormat,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewMetaTraderCSVKlineReader(csv.NewReader(strings.NewReader(tt.give)))
			kline, err := reader.Read()
			assert.Equal(t, tt.err, err)
			assertKlineEq(t, tt.want, kline)
		})
	}
}

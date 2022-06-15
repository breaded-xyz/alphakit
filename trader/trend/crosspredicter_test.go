package trend

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thecolngroup/alphakit/market"
	"github.com/thecolngroup/alphakit/ta"
	"github.com/thecolngroup/gou/dec"
	"github.com/thecolngroup/gou/num"
)

func TestCrossPredicter_ReceivePrice(t *testing.T) {
	var giveOsc, giveMMI ta.MockIndicator
	giveOsc.On("Update", []float64{10}).Return(error(nil))
	giveMMI.On("Update", []float64{7}).Return(error(nil))
	givePrice := market.Kline{C: dec.New(10)}
	givePrev := 3.0

	predicter := NewCrossPredicter(&giveOsc, &giveMMI)
	predicter.PriceSelector = ta.Close
	predicter.prev = givePrev
	err := predicter.ReceivePrice(context.Background(), givePrice)

	assert.NoError(t, err)
	giveOsc.AssertExpectations(t)
	giveMMI.AssertExpectations(t)
}

func TestCrossPredicter_Predict(t *testing.T) {
	tests := []struct {
		name          string
		giveOscValues []float64
		giveMMIValues []float64
		want          float64
	}{
		{
			name:          "flat @ 0",
			giveOscValues: []float64{10, 10},
			giveMMIValues: []float64{0, 0},
			want:          0,
		},
		{
			name:          "flat @ 0.1, MMI down-trend only",
			giveOscValues: []float64{10, 10},
			giveMMIValues: []float64{75, 70},
			want:          0.1,
		},
		{
			name:          "long @ 1.0, cross up zero w/ MMI",
			giveOscValues: []float64{-20, 20},
			giveMMIValues: []float64{75, 70},
			want:          1.0,
		},
		{
			name:          "long @ 0.9, cross up w/no MMI",
			giveOscValues: []float64{-20, 20},
			giveMMIValues: []float64{70, 75},
			want:          0.9,
		},
		{
			name:          "short @ -1.0",
			giveOscValues: []float64{20, -20},
			giveMMIValues: []float64{75, 70},
			want:          -1.0,
		},
		{
			name:          "short @ -0.9",
			giveOscValues: []float64{20, -20},
			giveMMIValues: []float64{70, 75},
			want:          -0.9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			predicter := NewCrossPredicter(
				&ta.StubIndicator{Values: tt.giveOscValues},
				&ta.StubIndicator{Values: tt.giveMMIValues},
			)
			act := predicter.Predict()
			assert.Equal(t, tt.want, num.Round2DP(act))
		})
	}

}

func TestCrossPredicter_Valid(t *testing.T) {
	predicter := NewCrossPredicter(
		&ta.StubIndicator{IsValid: true},
		&ta.StubIndicator{IsValid: true},
	)

	want := true
	act := predicter.Valid()
	assert.Equal(t, want, act)
}

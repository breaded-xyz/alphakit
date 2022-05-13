package trend

import (
	"context"

	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/ta"
	"github.com/colngroup/zero2algo/trader"
)

var _ trader.Predicter = (*CrossPredicter)(nil)

// CrossPredicter predicts price direction based on a moving average cross
// with a market meaness index filter.
type CrossPredicter struct {
	// PriceSelector is the kline component to use for price. Close by default.
	PriceSelector ta.PriceSelector

	// Osc is the moving average oscillator to use for entery and exit signals.
	Osc ta.Indicator[float64]

	// MMI is the signal filter.
	MMI ta.Indicator[float64]

	prev float64
}

// NewCrossPredicter creates a new predicter with Close quote price selector.
func NewCrossPredicter(osc, mmi ta.Indicator[float64]) *CrossPredicter {
	return &CrossPredicter{
		PriceSelector: ta.Close,
		Osc:           osc,
		MMI:           mmi,
	}
}

// ReceivePrice updates the prediction algo with the next market price.
// Call Predict() to get the resulting score.
func (p *CrossPredicter) ReceivePrice(ctx context.Context, price market.Kline) error {

	v := p.PriceSelector(price)
	vDiff := v - p.prev
	p.prev = v

	if err := p.Osc.Update(v); err != nil {
		return err
	}
	if err := p.MMI.Update(vDiff); err != nil {
		return err
	}

	return nil
}

// Predict returns a score to indicate confidence of price direction.
//
// 1.0 = Long trend cross over with MMI in confluence.
//
// 0.9 = Long trend cross over (no MMI confluence).
//
// -0.9 = Short trend cross over (no MMI confluence).
//
// -1.0 = Short trend with MMI confluence.
//
// [0.0, 0.1] = Flat trend.
func (p *CrossPredicter) Predict() float64 {

	var score float64

	if mmiSlope := ta.Slope(ta.Lookback(p.MMI.History(), 1), ta.Lookback(p.MMI.History(), 0)); mmiSlope < 0 {
		score += 0.1
	}

	switch {
	case ta.CrossUp(p.Osc.History(), 0):
		score += 0.9
	case ta.CrossDown(p.Osc.History(), 0):
		score = -(score + 0.9)
	}

	return score
}

// Valid returns true if Osc and MMI indicators are valid.
func (p *CrossPredicter) Valid() bool {
	return p.Osc.Valid() && p.MMI.Valid()
}

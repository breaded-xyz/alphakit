package day

import (
	"context"
	"time"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/ta"
	"github.com/colngroup/zero2algo/trader"
	"golang.org/x/exp/slices"
)

var _ trader.Bot = (*Bot)(nil)

type sessionRow struct {
	Start             time.Time
	YLow              float64
	YVAL              float64
	YPOC              float64
	YVAH              float64
	YHigh             float64
	SessionOpenPrice  float64
	FirstHourPrice    float64
	FirstHourVWAP     float64
	LinRegAlpha       float64
	LinRegBeta        float64
	LinRegR2          float64
	YNakedLowDistPct  float64
	YNakedVALDistPct  float64
	YNakedPOCDistPct  float64
	YNakedVAHDistPct  float64
	YNakedHighDistPct float64
	CrossYLow         bool
	CrossYVAL         bool
	CrossYPOC         bool
	CrossYVAH         bool
	CrossYHigh        bool
}

// Bot is a trader.Bot implementation for day trading.
type Bot struct {
	Prices   []market.Kline
	Levels   []VolumeLevel
	Profiles []*VolumeProfile
	Results  []sessionRow

	vwap VWAP
}

// NewBot creates a new Bot instance.
func NewBot() *Bot {
	return &Bot{}
}

// SetDealer sets the broker used for placing orders.
func (b *Bot) SetDealer(dealer broker.Dealer) {

}

// Warmup seeds the Predicter and Risker with historical price data.
func (b *Bot) Warmup(ctx context.Context, prices []market.Kline) error {

	return nil
}

// ReceivePrice updates the algo with latest market price potentially triggering buy and/or sell orders.
func (b *Bot) ReceivePrice(ctx context.Context, price market.Kline) error {

	b.Prices = append(b.Prices, price)
	_ = b.vwap.Update(price)

	// Add new Level for kline using HL2
	levelsNow := VolumeLevel{
		Price:  ta.HLC3(price),
		Volume: price.Volume,
	}
	b.Levels = append(b.Levels, levelsNow)

	// Initialize new day
	if h, m, s := price.Start.UTC().Clock(); h == 0 && m == 0 && s == 0 {
		b.vwap = VWAP{}

		if len(b.Levels) < 1440 {
			return nil
		}

		ySession := b.Levels[len(b.Levels)-1440:]
		slices.SortStableFunc(ySession, func(i, j VolumeLevel) bool {
			return i.Price < j.Price
		})

		vp := NewVolumeProfile(10, ySession)
		b.Profiles = append(b.Profiles, vp)
		b.Results = append(b.Results, sessionRow{
			Start:            price.Start.UTC(),
			SessionOpenPrice: price.O.InexactFloat64(),
			YLow:             vp.Low,
			YVAL:             vp.VAL,
			YPOC:             vp.POC,
			YVAH:             vp.VAH,
			YHigh:            vp.High,
		})
		return nil
	}

	if h, m, s := price.Start.UTC().Clock(); h > 0 && (h < 23 && m < 59 && s < 59) {

		if len(b.Results) == 0 {
			return nil
		}

		low, high := price.L.InexactFloat64(), price.H.InexactFloat64()

		session := b.Results[len(b.Results)-1]
		if crossed(session.YLow, low, high) {
			session.CrossYLow = true
		}
		if crossed(session.YVAL, low, high) {
			session.CrossYVAL = true
		}
		if crossed(session.YPOC, low, high) {
			session.CrossYPOC = true
		}
		if crossed(session.YVAH, low, high) {
			session.CrossYVAH = true
		}
		if crossed(session.YHigh, low, high) {
			session.CrossYHigh = true
		}
		b.Results[len(b.Results)-1] = session
	}

	// Initialize distance to key levels after first hour
	if h, m, s := price.Start.UTC().Clock(); h == 1 && m == 0 && s == 0 {
		if len(b.Results) == 0 {
			return nil
		}
		session := b.Results[len(b.Results)-1]
		session.FirstHourPrice = price.O.InexactFloat64()
		session.FirstHourVWAP = b.vwap.Value()

		if !session.CrossYLow {
			session.YNakedLowDistPct = (session.YLow - session.FirstHourPrice) / session.FirstHourPrice
		}
		if !session.CrossYVAL {
			session.YNakedVALDistPct = (session.YVAL - session.FirstHourPrice) / session.FirstHourPrice
		}
		if !session.CrossYPOC {
			session.YNakedPOCDistPct = (session.YPOC - session.FirstHourPrice) / session.FirstHourPrice
		}
		if !session.CrossYVAH {
			session.YNakedVAHDistPct = (session.YVAH - session.FirstHourPrice) / session.FirstHourPrice
		}
		if !session.CrossYHigh {
			session.YNakedHighDistPct = (session.YHigh - session.FirstHourPrice) / session.FirstHourPrice
		}

		/*xs := make([]float64, 60)
		ys := make([]float64, 60)
		vwapSeries := ta.Window(b.vwap.History(), 60)
		for i := range vwapSeries {
			xs[i] = float64(i)
			ys[i] = vwapSeries[i]
		}
		alpha, beta := stat.LinearRegression(xs, ys, nil, false)
		r2 := stat.RSquared(xs, ys, nil, alpha, beta)
		session.LinRegAlpha = alpha
		session.LinRegBeta = beta
		session.LinRegR2 = r2*/

		b.Results[len(b.Results)-1] = session
	}

	return nil
}

func crossed(v, lower, upper float64) bool {
	return v >= lower && v <= upper
}

// Close exits all open positions at current market price.
func (b *Bot) Close(ctx context.Context) error {

	return nil
}

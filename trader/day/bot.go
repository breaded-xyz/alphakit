package day

import (
	"context"
	"time"

	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/internal/util"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/ta"
	"github.com/colngroup/zero2algo/trader"
	"github.com/gonum/stat"
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
	SessionOpen       float64
	HourClose         float64
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

	// Add new Level for kline using HL2
	levelsNow := VolumeLevel{
		Price:  util.RoundTo(ta.HL2(price), 1.0),
		Volume: util.RoundTo(price.Volume, 1.0),
	}
	b.Levels = append(b.Levels, levelsNow)

	// Initialize new day
	if h, m, s := price.Start.UTC().Clock(); h == 0 && m == 0 && s == 0 {
		if len(b.Levels) < 1440 {
			return nil
		}

		ySession := b.Levels[len(b.Levels)-1440:]
		slices.SortStableFunc(ySession, func(i, j VolumeLevel) bool {
			return i.Price < j.Price
		})

		vp := NewVolumeProfile(100, ySession)
		b.Profiles = append(b.Profiles, vp)
		b.Results = append(b.Results, sessionRow{
			Start:       price.Start.UTC(),
			SessionOpen: price.O.InexactFloat64(),
			YLow:        vp.Low,
			YVAL:        vp.VAL,
			YPOC:        vp.POC,
			YVAH:        vp.VAH,
			YHigh:       vp.High,
		})
		return nil
	}

	if h, m, s := price.Start.UTC().Clock(); h > 0 && (h < 23 && m < 59 && s < 59) {

		if len(b.Results) == 0 || len(b.Levels) == 0 {
			return nil
		}

		session := b.Results[len(b.Results)-1]
		if crossed(b.Levels[len(b.Levels)-2].Price, levelsNow.Price, session.YLow) {
			session.CrossYLow = true
		}
		if crossed(b.Levels[len(b.Levels)-2].Price, levelsNow.Price, session.YVAL) {
			session.CrossYVAL = true
		}
		if crossed(b.Levels[len(b.Levels)-2].Price, levelsNow.Price, session.YPOC) {
			session.CrossYPOC = true
		}
		if crossed(b.Levels[len(b.Levels)-2].Price, levelsNow.Price, session.YVAH) {
			session.CrossYVAH = true
		}
		if crossed(b.Levels[len(b.Levels)-2].Price, levelsNow.Price, session.YHigh) {
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
		session.HourClose = price.C.InexactFloat64()

		if !session.CrossYLow {
			session.YNakedLowDistPct = (session.YLow - session.HourClose) / session.HourClose
		}
		if !session.CrossYVAL {
			session.YNakedVALDistPct = (session.YVAL - session.HourClose) / session.HourClose
		}
		if !session.CrossYPOC {
			session.YNakedPOCDistPct = (session.YPOC - session.HourClose) / session.HourClose
		}
		if !session.CrossYVAH {
			session.YNakedVAHDistPct = (session.YVAH - session.HourClose) / session.HourClose
		}
		if !session.CrossYHigh {
			session.YNakedHighDistPct = (session.YHigh - session.HourClose) / session.HourClose
		}

		xs := make([]float64, 60)
		ys := make([]float64, 60)
		sessionPrices := b.Prices[len(b.Prices)-60:]
		for i := range sessionPrices {
			xs[i] = float64(i)
			ys[i] = sessionPrices[i].C.InexactFloat64()
		}
		alpha, beta := stat.LinearRegression(xs, ys, nil, false)
		r2 := stat.RSquared(xs, ys, nil, alpha, beta)
		session.LinRegAlpha = alpha
		session.LinRegBeta = beta
		session.LinRegR2 = r2

		b.Results[len(b.Results)-1] = session
		return nil
	}

	if h, m, s := price.Start.UTC().Clock(); h > 0 && (h < 23 && m < 59 && s < 59) {

		if len(b.Results) == 0 || len(b.Levels) == 0 {
			return nil
		}

		session := b.Results[len(b.Results)-1]
		if crossed(b.Levels[len(b.Levels)-2].Price, levelsNow.Price, session.YLow) {
			session.CrossYLow = true
		}
		if crossed(b.Levels[len(b.Levels)-2].Price, levelsNow.Price, session.YVAL) {
			session.CrossYVAL = true
		}
		if crossed(b.Levels[len(b.Levels)-2].Price, levelsNow.Price, session.YPOC) {
			session.CrossYPOC = true
		}
		if crossed(b.Levels[len(b.Levels)-2].Price, levelsNow.Price, session.YVAH) {
			session.CrossYVAH = true
		}
		if crossed(b.Levels[len(b.Levels)-2].Price, levelsNow.Price, session.YHigh) {
			session.CrossYHigh = true
		}
		b.Results[len(b.Results)-1] = session
	}

	return nil
}

func crossed(a, b, c float64) bool {
	return ta.CrossUp([]float64{a, b}, c) || ta.CrossDown([]float64{a, b}, c)
}

// Close exits all open positions at current market price.
func (b *Bot) Close(ctx context.Context) error {

	return nil
}

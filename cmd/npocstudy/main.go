package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/colngroup/zero2algo/internal/studyrun"
	"github.com/colngroup/zero2algo/market"
	"github.com/colngroup/zero2algo/ta"
	"golang.org/x/exp/slices"
)

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

var (
	prices []market.Kline
	levels []ta.VolumeLevel
	//profiles []*ta.VolumeProfile
	results []sessionRow
	vwap    ta.VWAP
)

func main() {

	var err error

	pricePath := "/Users/richklee/Dropbox/dev-share/github.com/colngroup/zero2algo/prices/spot/btcusdt-m1/all"

	prices, err = studyrun.ReadPriceSeries(pricePath)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	for i := range prices {
		fmt.Printf("%s\n", prices[i].Start)
		//spew.Dump(prices[i])
		//spew.Dump(i)
		_ = receivePrice(prices[i])
	}

	_ = studyrun.SaveStructToCSV("results.csv", results)
}

func receivePrice(price market.Kline) error {

	if price.Start.UTC().Equal(time.Date(2018, 06, 19, 00, 00, 00, 0, time.UTC)) {
		println("break")
	}

	_ = vwap.Update(price)

	// Add new Level for kline using HL2
	levelsNow := ta.VolumeLevel{
		Price:  ta.HLC3(price),
		Volume: price.Volume,
	}
	levels = append(levels, levelsNow)

	// Initialize new day
	if h, m, s := price.Start.UTC().Clock(); h == 0 && m == 0 && s == 0 {
		vwap = ta.VWAP{}

		if len(levels) < 1440 {
			return nil
		}

		ySession := levels[len(levels)-1440:]
		slices.SortStableFunc(ySession, func(i, j ta.VolumeLevel) bool {
			return i.Price < j.Price
		})

		vp := ta.NewVolumeProfile(10, ySession)
		//profiles = append(profiles, vp)
		results = append(results, sessionRow{
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

		if len(results) == 0 {
			return nil
		}

		low, high := price.L.InexactFloat64(), price.H.InexactFloat64()

		session := results[len(results)-1]
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
		results[len(results)-1] = session
	}

	// Initialize distance to key levels after first hour
	if h, m, s := price.Start.UTC().Clock(); h == 1 && m == 0 && s == 0 {
		if len(results) == 0 {
			return nil
		}
		session := results[len(results)-1]
		session.FirstHourPrice = price.O.InexactFloat64()
		session.FirstHourVWAP = vwap.Value()

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

		results[len(results)-1] = session
	}

	return nil
}

func crossed(v, lower, upper float64) bool {
	return v >= lower && v <= upper
}
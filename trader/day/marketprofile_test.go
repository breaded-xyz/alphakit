package day

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestMarketProfile(t *testing.T) {

	prices := []float64{10, 10, 12, 12, 1}
	volumes := []float64{1, 1, 2, 3, 100}

	mp := NewMarketProfile(4, prices, volumes)

	spew.Dump(mp.Bins, mp.Hist)

}

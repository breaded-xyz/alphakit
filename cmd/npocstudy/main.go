package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/colngroup/zero2algo/internal/studyrun"
	"github.com/colngroup/zero2algo/trader/day"
)

func main() {
	bot := day.NewBot()

	testdataPath := "/Users/richklee/Dropbox/dev-share/github.com/colngroup/zero2algo/prices/spot/btcusdt-m1/all"

	prices, err := studyrun.ReadPriceSeries(testdataPath)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	for i := range prices {
		fmt.Printf("%s\n", prices[i].Start)
		err := bot.ReceivePrice(context.Background(), prices[i])
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}
	}

	studyrun.SaveStructToCSV("./testdata/results-9.csv", bot.Results)
}

package optimize

import (
	"context"
	"fmt"

	"github.com/thecolngroup/alphakit/market"
	"github.com/thecolngroup/alphakit/trader/trend"
)

func Example() {
	// Verbose error handling ommitted for brevity

	// Identify the bot (algo) to optimize by supplying a factory function
	// Here we're using the classic MA cross variant of trend bot
	bot := trend.MakeCrossBotFromConfig

	// Define the parameter space to optimize
	// Param names must match those expected by the MakeBot function passed to optimizer
	paramSpace := ParamMap{
		"mafastlength": []float64{1, 10, 20, 30},
		"maslowlength": []float64{30, 40, 50, 60},
		"mmilength":    []float64{200, 300},
	}

	// Read price samples to use for optimization
	btcPriceSample, _ := market.ReadKlinesFromCSV("testdata/")
	ethPriceSample, _ := market.ReadKlinesFromCSV("testdata/")
	priceSamples := [][]market.Kline{btcPriceSample, ethPriceSample}

	// Create a new brute style optimizer
	optimizer := NewBruteOptimizer()
	optimizer.SampleSplitPct = 0.5
	optimizer.WarmupBarCount = 300
	optimizer.MakeBot = bot

	// Prepare the optimizer and get an estimate on the number of trials (backtests) required
	trialCount, _ := optimizer.Prepare(paramSpace, priceSamples)
	fmt.Printf("%d trials to run during optimization\n", trialCount)

	// Start the optimization process and monitor with a receive channel
	trials, _ := optimizer.Start(context.Background())
	for range trials {
	}

	// Inspect the study to get the optimized param set and results
	study := optimizer.Study()
	optimaPSet := study.Validation[0]
	optimaResult := study.ValidationResults[optimaPSet.ID]

	// Output: Optima sharpe ratio is 0.0
	fmt.Printf("Optima sharpe ratio is %.2f", optimaResult.Sharpe)
}

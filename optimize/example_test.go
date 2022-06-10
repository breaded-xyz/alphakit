package optimize

import (
	"context"
	"fmt"
	"log"

	"github.com/thecolngroup/alphakit/market"
	"github.com/thecolngroup/alphakit/trader/trend"
)

func Example() {
	// Verbose error handling omitted for brevity

	// Identify the bot (algo) to optimize by supplying a factory function
	// Here we're using the classic moving average (MA) cross variant of trend bot
	bot := trend.MakeCrossBotFromConfig

	// Define the parameter space to optimize
	// Param names must match those expected by the MakeBot function passed to optimizer
	// Here we're optimizing the lookback period of a fast and slow MA
	// and the Market Meanness Index (MMI) filter
	paramSpace := ParamMap{
		"mafastlength": []any{30, 90, 180},
		"maslowlength": []any{90, 180, 360},
		"mmilength":    []any{200, 300},
	}

	// Read price samples to use for optimization
	btc, err := market.ReadKlinesFromCSV("./testdata/BTCUSDT-1H/")
	if err != nil {
		log.Fatal(err)
	}
	eth, err := market.ReadKlinesFromCSV("./testdata/ETHUSDT-1H/")
	if err != nil {
		log.Fatal(err)
	}
	priceSamples := map[AssetID][]market.Kline{"btc": btc, "eth": eth}

	// Create a new brute style optimizer with a default simulated dealer (no broker costs)
	// The default optimization objective is the param set with the highest sharpe ratio
	optimizer := NewBruteOptimizer()
	optimizer.SampleSplitPct = 0   // Do not split samples due to small sample size
	optimizer.WarmupBarCount = 360 // Set as maximum lookback of your param space
	optimizer.MakeBot = bot        // Tell the optimizer which bot to use

	// Prepare the optimizer and get an estimate on the number of trials (backtests) required
	trialCount, _ := optimizer.Prepare(paramSpace, priceSamples)
	fmt.Printf("%d backtest trials to run during optimization\n", trialCount)

	// Start the optimization process and monitor with a receive-only channel
	// Trials will execute concurrently with a default worker pool matching num of CPUs
	trials, _ := optimizer.Start(context.Background())
	for range trials {
	}

	// Inspect the study results following optimization
	study := optimizer.Study()
	if len(study.ValidationResults) == 0 {
		fmt.Println("Optima not found because highest ranked param set made no trades during optimization trials.")
		return
	}

	// Read out the optimal param set and results
	optimaPSet := study.Validation[0]
	fmt.Printf("Optima params: fast: %d slow: %d MMI: %d\n",
		optimaPSet.Params["mafastlength"], optimaPSet.Params["maslowlength"], optimaPSet.Params["mmilength"])
	optimaResult := study.ValidationResults[optimaPSet.ID]
	fmt.Printf("Optima sharpe ratio is %.1f", optimaResult.Sharpe)

	// Output:
	// 38 backtest trials to run during optimization
	// Optima params: fast: 30 slow: 90 MMI: 200
	// Optima sharpe ratio is 2.5
}

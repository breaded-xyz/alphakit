# Alpha Kit

A framework for algorithmic trading in Go and serverless cloud

_"Master the latest features of Go and learn how to design, validate and deploy sound algorithmic trading strategies."_

Companion code repository for the forthcoming book __Zero to Algo__.

<<Insert gif of studyrun execution>>

## Inspiration

The majority of open source algo trading frameworks, especially in Go, focus purely on trade execution - that's a sure fire way to get rekt. The most important activity is researching and validating practical market edges - which is the focus of alphakit. Furthermore, I wanted a clean architecture that could easily be composed and executed serverless in the cloud, rather than a monolithic black box running on a carpet server under a desk!

## What's included?

A complete starter kit for developing algorithmic trading strategies in the Go language:

- Example buy-and-hold and trend-following algos
- Backtest and walk-forward engine to evaluate algos
- Performance reports and metrics
- Brute force parameter optimization method
- Command app to execute research studies from a config file
- Path to serverless production deployment in the cloud (coming soon)
- Uses latest Go language (1.18) features including generics
- Idiomatic Go style using community accepted best practises
- Pragmatic use of concurrency, go routines and channels

## Install

go get "github.com/thecolngroup/alphakit"

## Getting started

The canonical example that brings together many of the framework components is in the optimize package. A further well documented example in the backtest package demonstrates how to use the simulated dealer to study algos without an optimizer.

```

func Example() {
	// Verbose error handling ommitted for brevity

	// Identify the bot (algo) to optimize by supplying a factory function
	// Here we're using the classic moving average (MA) cross variant of trend bot
	bot := trend.MakeCrossBotFromConfig

	// Define the parameter space to optimize
	// Param names must match those expected by the MakeBot function passed to optimizer
	// Here we're optimizing the lookback period of a fast and slow MA
	// and the Market Meaness Index (MMI) filter
	paramSpace := ParamMap{
		"mafastlength": []any{30, 90, 180},
		"maslowlength": []any{90, 180, 360},
		"mmilength":    []any{200, 300},
	}

	// Read price samples to use for optimization
	btc, _ := market.ReadKlinesFromCSV("testdata/btcusdt-1h/")
	eth, _ := market.ReadKlinesFromCSV("testdata/ethusdt-1h/")
	priceSamples := [][]market.Kline{btc, eth}

	// Create a new brute style optimizer with a default simulated dealer (no broker costs)
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
	fmt.Printf("Optima sharpe ratio is %.2f", optimaResult.Sharpe)

	// Output:
	// 38 backtest trials to run during optimization
	// Optima params: fast: 30 slow: 90 MMI: 200
	// Optima sharpe ratio is 2.46
}

```

## Working with price data

## Building a trading bot

In the trader package you will find a couple of example bots: hodl and trend. The hodl bot is useful for benchmarking an asset, and the trend bot serves as a template for developing your own algo.

### Prediction

### Risk Management

### Money Management

## Command app: studyrun

## Connecting to a live trading venue

## Further reading

## Contributing

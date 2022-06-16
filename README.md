<!--
 Copyright 2022 The Coln Group Ltd
 SPDX-License-Identifier: MIT
-->

# Alphakit

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/gomods/athens)
[![Go Report Card](https://goreportcard.com/badge/github.com/thecolngroup/alphakit)](https://goreportcard.com/report/github.com/thecolngroup/alphakit)
[![Go](https://github.com/thecolngroup/alphakit/actions/workflows/go.yml/badge.svg)](https://github.com/thecolngroup/alphakit/actions/workflows/go.yml)

Introducing a framework for algorithmic trading in Go and serverless cloud

```
           /$$           /$$                 /$$       /$$   /$$    
          | $$          | $$                | $$      |__/  | $$    
  /$$$$$$ | $$  /$$$$$$ | $$$$$$$   /$$$$$$ | $$   /$$ /$$ /$$$$$$  
 |____  $$| $$ /$$__  $$| $$__  $$ |____  $$| $$  /$$/| $$|_  $$_/  
  /$$$$$$$| $$| $$  \ $$| $$  \ $$  /$$$$$$$| $$$$$$/ | $$  | $$    
 /$$__  $$| $$| $$  | $$| $$  | $$ /$$__  $$| $$_  $$ | $$  | $$ /$$
|  $$$$$$$| $$| $$$$$$$/| $$  | $$|  $$$$$$$| $$ \  $$| $$  |  $$$$/
 \_______/|__/| $$____/ |__/  |__/ \_______/|__/  \__/|__/   \___/  
              | $$                                                  
              | $$                                                  
              |__/                                                  
```

Companion code repository for the forthcoming book __Zero to Algo__

> "Master the latest features of Go and learn how to design, validate and deploy sound algorithmic trading strategies."

## Inspiration

The majority of open source algo trading frameworks, especially in Go, focus purely on trade execution - that's a sure fire way to get rekt. The most important precursor to a successful trading system is researching and validating practical market edges - which is the focus of alphakit. Furthermore, I wanted a composable architecture that could easily be executed serverless in the cloud using features such as cloud functions and messaging queues.

## What's included?

A complete starter kit for developing algorithmic trading strategies in the Go language:

- Example buy-and-hold and trend-following algos
- Backtest and walk-forward engine to evaluate algos
- Performance reports and metrics
- Brute force parameter optimization method
- Command app to execute research studies from a config file
- Scaffold for serverless production deployment in the cloud (coming soon)
- Uses latest Go language (1.18) features including generics
- Idiomatic Go style using community accepted best practices
- Pragmatic use of concurrency, go routines and channels
- Extensive test coverage where it matters

## Install

`go get "github.com/thecolngroup/alphakit"`

## Getting started

⚠️  API is pre v1 and is not stable

The canonical example that brings together many of the framework components is in the `optimize` package and reproduced below. A further well documented example in the `backtest` package demonstrates how to use a simulated dealer to study algos without an optimizer.

```go

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
 btc, _ := market.ReadKlinesFromCSV("testdata/btcusdt-1h/")
 eth, _ := market.ReadKlinesFromCSV("testdata/ethusdt-1h/")
 priceSamples := map[AssetID][]market.Kline{"btc": btc, "eth": eth}

 // Create a new brute style optimizer with a default simulated dealer (no broker costs)
 // The default optimization objective is the param set with the highest sharpe ratio
 optimizer := NewBruteOptimizer()
 optimizer.SampleSplitPct = 0.5   // Use first 50% as in-sample training data, and remainder for out-of-sample validation
 optimizer.WarmupBarCount = 360 // Set as maximum lookback of your param space - 360 is the longest lookback for slow MA
 optimizer.MakeBot = bot        // Tell the optimizer which bot to use

 // Prepare the optimizer and get an estimate on the number of trials (backtests) required
 trialCount, _ := optimizer.Prepare(paramSpace, priceSamples)
 fmt.Printf("%d backtest trials to run during optimization\n", trialCount)

 // Start the optimization process and monitor with a receive-only channel
 // Trials will execute concurrently with a default worker pool matching the num of CPUs
 trials, _ := optimizer.Start(context.Background())
 for range trials {
   // Monitor for errors and progress
 }

 // Inspect the study results following optimization
 study := optimizer.Study()

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

## Beware testing bias

Alphakit's aim is to demonstrate best practice when developing algo strategies. This requires a sound approach to backtesting and performance analysis in order to mimimize bias, the most important of which is overfitting whereby you mistake random chance for edge. This is a good summary of the challenges: <https://robotwealth.com/backtesting-bias-feels-good-until-you-blow-up/>.

To this end, the `BruteOptimizer` implementation enables you to split the given price samples into in-sample and out-of-sample buckets.

```go

// Use first 50% of data as in-sample training data, and remainder for out-of-sample validation
optimizer.SampleSplitPct = 0.5 

```

Training is conducted on the in-sample and performance validation on the out-of-sample. However, this is not a pancea and can still result in an overfitted algo if you attempt to optimize too many parameters at the same time.

There are a number of useful articles on <https://financial-hacker.com/> that explore the pitfalls of backtesting in more detail.

## Fundamental architecture patterns

The core assumption underlying the framework is that price data enters the system at a defined interval. Each time a new kline arrives it triggers an evaluation process owned by a bot that may result in 1 or more new orders being issued to a dealer.

Every component that participates in this processing implements the `market.Receiver` interface and accepts a kline (and a context to control long running operations).

The `broker` package offers an API to mediate the interaction between bot and trading venue. A bot creates market positions by placing orders through an implementation of `Dealer`. A simulated dealer in the `backtest` package (also a price receiver) allows you study and validate algos outsde of an optimizer.

In future releases new `Dealer` implementations will enable you to connect to specific trading venues.

## Working with price data

The price data used in the unit tests and examples is sourced from Binance. It's a good source of clean crypto data going back to late 2017. See <https://github.com/binance/binance-public-data/>.

Alphakit offers an API for price data in the `market` package. The primary representation is in the form of a candlestick (OHLC) - also known as a kline. `CSVKlineReader` reads klines from a .csv file, it can be extended to decode data from various sources with a `CSVKlineDecoder`. The default decoder supports the Binance data format which uses a unix millisecond format. A further decoder is also provided for MetaTrader data files.

Convenience functions for reading individual CSV files or walking a directory are also included.

## Performance reports

Package `perf` provides comprehensive performance reporting for your algo, enabling you to track industry standard metrics such as CAGR, return rate, sharpe ratio, and drawdowns.

To create a new report use the equity history and trade history data from a dealer.

## Trading costs

Many algos appear to be viable until you correctly factor in trading costs! Package `backtest` offers a `PerpCoster` implementation that simulates typical costs you might expect when trading crypto perpetual futures, including an hourly funding rate fee. See the tests in package `backtest` to understand how costs are applied during backtesting.

## Building a trading bot

In the `trader` package you will find a couple of example bots: hodl and trend. The hodl bot is useful for benchmarking an asset, and the trend bot serves as a template for developing your own algo.

The following notes refer to how the bot in the `trend` package operates.

### Prediction

`trader.Predicter` is a simple interface that returns a value between -1 and 1. A value of 1 signals maximum confidence in opening a long position, whilst -1 maximum confidence in opening a short position. 0 indicates no directional bias. Other values between -1 and 1 indicate varying confidence in direction.

`CrossPredicter` uses a fast and slow moving average cross with a Market Meanness Index (MMI) filter to determine the prediction.

`ApexPredicter` uses peak and valley detection in a smoothed price series with an MMI filter.

To understand more about trend following and MMI this is a great starting point: <https://financial-hacker.com/trend-and-exploiting-it/>

The trend bot interprets the prediction value according to a set of threshold values for opening and closing positions, namely:

- `EnterLong`
- `EnterShort`
- `ExitLong`
- `ExitShort`

By varying these threshold values you can create asymmetric entry and exit conditions.

### Risk Management

Package `risk` provides methods to calculate unit risk, used as an input to position sizing and stop loss specification.

By default 'full-risk' will be used which assumes no stop-loss. As an alternative a standard deviation method is also provided.

### Money Management

Package `money` provides methods to size a position. By default (and recommended for initial backtesting) a position based on a fixed capital amount is used.

A more sophisticated option is to a use a fixed fraction method given by the `SafeFSizer` type. You can determine the optimal value of 'f' by using the OptimalF or Kelly value from a performance report.

## Command app: studyrun

The command app `studyrun` enables you to execute optimization studies by specifying a `.toml` config file.

See the test in `cmd/studyrun` to understand the syntax and play with a working example.

If you wish to use your own custom bots or price data decoders with the command app you'll need to update the contents of the type registry passed to the app from `main`.

The command app will execute an optimization study using `BruteOptimizer` and dump out the results in .csv format.

If you wish to execute the `studyrun` process from outside Alphakit, an entrypoint function is availble in the package `github.com/thecolngroup/alphakit/cmd/studyrun/app`.

## Connecting to a live trading venue

Future releases will provide implementations of `broker.Dealer` for specific trading venues. Contributions welcome!

## Further reading

- <https://financial-hacker.com/>
- <https://robotwealth.com/>
- <https://quantocracy.com/>
- <https://zorro-project.com/manual/>

## Contributing

Please fork and raise a PR or submit an issue.

Contributions should comply with:

- Default golangci linters (see config file in root)
- Uber style guide: <https://github.com/uber-go/guide/blob/master/style.md>

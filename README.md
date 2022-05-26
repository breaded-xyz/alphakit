# Alpha Kit

Framework for algorithmic trading in Go and serverless cloud

_"Master the latest features of Go and learn how to design, validate and deploy sound algorithmic trading strategies."_

Companion code repository for the forthcoming book __Zero to Algo__.

<<Insert gif of studyrun execution>>

## Inspiration

Most open source algo trading frameworks, especially in Go, focus purely on trade execution - thats a sure fire way to get rekt! The most important activity is researching and validating practical market edges, which is the focus of alphakit. Furthermore, I wanted a clean architecture that could easily be composed and executed serverless in the cloud, rather than a monolithic black box running in Docker.

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

The canonical example that brings together many of the framework components is in the optimize package. Enclosed below:

<<>>

A further well documented example is in the backtest package that demonstrates how to use the simulated dealer to study algos without using the optimizer.

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

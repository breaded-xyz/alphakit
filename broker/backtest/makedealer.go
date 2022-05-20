package backtest

import (
	"github.com/thecolngroup/alphakit/broker"
	"github.com/thecolngroup/alphakit/internal/dec"
)

func MakeDealerFromConfig(config map[string]any) (broker.SimulatedDealer, error) {

	dealer := NewDealer()

	dealer.simulator.SetInitialCapital(dec.New(config["initialcapital"].(float64)))
	dealer.simulator.cost = &PerpCost{
		SpreadPct:      dec.New(config["spreadpct"].(float64)),
		SlippagePct:    dec.New(config["slippagepct"].(float64)),
		TransactionPct: dec.New(config["transactionpct"].(float64)),
		FundingHourPct: dec.New(config["fundinghourpct"].(float64)),
	}

	return dealer, nil
}

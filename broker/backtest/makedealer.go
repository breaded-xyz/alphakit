package backtest

import (
	"github.com/colngroup/zero2algo/broker"
	"github.com/colngroup/zero2algo/internal/dec"
)

func MakeDealer(config map[string]any) (broker.SimulatedDealer, error) {

	var dealer Dealer

	dealer.simulator.SetInitialCapital(dec.New(config["initialcapital"].(float64)))
	dealer.simulator.cost = &PerpCost{
		SpreadPct:      dec.New(config["spreadpct"].(float64)),
		SlippagePct:    dec.New(config["slippagepct"].(float64)),
		TransactionPct: dec.New(config["transactionpct"].(float64)),
		FundingHourPct: dec.New(config["fundinghourpct"].(float64)),
	}

	return &dealer, nil
}

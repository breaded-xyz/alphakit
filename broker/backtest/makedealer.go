package backtest

import (
	"github.com/thecolngroup/alphakit/broker"
	"github.com/thecolngroup/dec"
	"github.com/thecolngroup/util"
)

// MakeDealerFromConfig mints a new dealer from a config source.
func MakeDealerFromConfig(config map[string]any) (broker.SimulatedDealer, error) {
	dealer := NewDealer()

	dealer.simulator.SetInitialCapital(dec.New(config["initialcapital"].(float64)))
	dealer.simulator.cost = &PerpCoster{
		SpreadPct:      dec.New(util.ToFloat(config["spreadpct"])),
		SlippagePct:    dec.New(util.ToFloat(config["slippagepct"])),
		TransactionPct: dec.New(util.ToFloat(config["transactionpct"])),
		FundingHourPct: dec.New(util.ToFloat(config["fundinghourpct"])),
	}

	return dealer, nil
}

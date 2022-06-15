package backtest

import (
	"github.com/thecolngroup/alphakit/broker"
	"github.com/thecolngroup/gou/conv"
	"github.com/thecolngroup/gou/dec"
)

// MakeDealerFromConfig mints a new dealer from a config source.
func MakeDealerFromConfig(config map[string]any) (broker.SimulatedDealer, error) {
	dealer := NewDealer()

	dealer.simulator.SetInitialCapital(dec.New(config["initialcapital"].(float64)))
	dealer.simulator.cost = &PerpCoster{
		SpreadPct:      dec.New(conv.ToFloat(config["spreadpct"])),
		SlippagePct:    dec.New(conv.ToFloat(config["slippagepct"])),
		TransactionPct: dec.New(conv.ToFloat(config["transactionpct"])),
		FundingHourPct: dec.New(conv.ToFloat(config["fundinghourpct"])),
	}

	return dealer, nil
}

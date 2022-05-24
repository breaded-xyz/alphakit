// Package studyrun is internal and not part of exported API.
package studyrun

import (
	"errors"

	"github.com/thecolngroup/alphakit/broker"
	"github.com/thecolngroup/alphakit/broker/backtest"
)

func ReadDealerFromConfig(config map[string]any) (broker.MakeSimulatedDealer, error) {

	var makeDealer broker.MakeSimulatedDealer

	if _, ok := config["dealer"]; !ok {
		return nil, errors.New("'dealer' key not found")
	}
	root := config["dealer"].(map[string]any)
	makeDealer = func() (broker.SimulatedDealer, error) {
		return backtest.MakeDealerFromConfig(root)
	}

	return makeDealer, nil
}

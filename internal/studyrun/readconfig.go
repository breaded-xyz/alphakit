// Copyright 2022 The Coln Group Ltd
// SPDX-License-Identifier: MIT

package studyrun

import (
	"os"

	"github.com/spf13/viper"
)

const _configFormat = "toml"

// readConfig reads the toml configuration file and returns a string map of params.
func readConfig(filename string) (map[string]any, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	//nolint:errcheck
	defer file.Close()
	viper.SetConfigType(_configFormat)
	if err := viper.ReadConfig(file); err != nil {
		return nil, err
	}
	return viper.AllSettings(), nil
}

// Copyright 2022 The Coln Group Ltd
// SPDX-License-Identifier: MIT

package studyrun

import "errors"

// readParamSpaceFromConfig creates a new param space from a config file params.
func readParamSpaceFromConfig(config map[string]any) (map[string]any, error) {

	pset, ok := config["paramspace"].(map[string]any)
	if !ok {
		return nil, errors.New("'paramspace' key not found")
	}

	return pset, nil
}

package studyrun

import "errors"

func ReadParamSetFromConfig(config map[string]any) (map[string]any, error) {

	pset, ok := config["paramset"].(map[string]any)
	if !ok {
		return nil, errors.New("'paramset' key not found")
	}

	return pset, nil
}

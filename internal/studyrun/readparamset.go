package studyrun

import "errors"

func ReadParamSpaceFromConfig(config map[string]any) (map[string]any, error) {

	pset, ok := config["paramspace"].(map[string]any)
	if !ok {
		return nil, errors.New("'paramspace' key not found")
	}

	return pset, nil
}

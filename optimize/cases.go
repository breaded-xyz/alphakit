package optimize

import (
	"github.com/schwarmco/go-cartesian-product"
)

type keyValuePair struct {
	k string
	v any
}

func BuildTestCases(paramRanges map[string][]any) []map[string]any {

	testCases := make([]map[string]any, 0)

	// Prepare a slice of sets to pass to the cartesian func
	// Each element in a set is a key-value pair (param name, param value)
	inputSets := make([][]any, 0)
	for k, vArr := range paramRanges {
		set := make([]any, 0)
		for _, v := range vArr {
			kv := keyValuePair{k, v}
			set = append(set, kv)
		}
		inputSets = append(inputSets, set)
	}

	// Produce the cartesian products passing in the input sets
	// Marshal each product (a slice of key-value pairs) back to a test case map
	productCh := cartesian.Iter(inputSets...)
	for product := range productCh {
		tCase := make(map[string]any, len(product))
		for _, kv := range product {
			tCase[kv.(keyValuePair).k] = kv.(keyValuePair).v
		}
		testCases = append(testCases, tCase)
	}

	return testCases
}

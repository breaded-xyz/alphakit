package optimize

import (
	"github.com/schwarmco/go-cartesian-product"
)

type TestCase map[string]any

type ParamRange map[string][]any

type keyValuePair struct {
	k string
	v any
}

func BuildBacktestCases(in ParamRange) []TestCase {

	testCases := make([]TestCase, 0)

	// Prepare a slice of sets to pass to the cartesian func
	// Each element in a set is a key-value pair (param name, param value)
	inputSets := make([][]any, 0)
	for paramName, paramValues := range in {
		set := make([]any, 0)
		for i := range paramValues {
			kv := keyValuePair{paramName, paramValues[i]}
			set = append(set, kv)
		}
		inputSets = append(inputSets, set)
	}

	// Produce the cartesian products passing in the input sets
	// Marshal each product (a slice of key-value pairs) back to a test case map
	productCh := cartesian.Iter(inputSets...)
	for product := range productCh {
		tCase := make(map[string]any, len(product))
		for i := range product {
			kv := product[i].(keyValuePair)
			tCase[kv.k] = kv.v
		}
		testCases = append(testCases, tCase)
	}

	return testCases
}

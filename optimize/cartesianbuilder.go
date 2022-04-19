package optimize

import (
	"reflect"

	"github.com/schwarmco/go-cartesian-product"
)

type CartesianProduct map[string]any

type keyValuePair struct {
	k string
	v any
}

func CartesianBuilder(in map[string]any) []CartesianProduct {
	products := make([]CartesianProduct, 0)

	// Prepare a slice of sets to pass to the cartesian func
	// Each element in a set is a key-value pair (param name, param value)
	// All elements in a set share the same key name
	cartesianInputSets := make([][]any, 0)
	for paramName, paramValue := range in {
		set := make([]any, 0)

		// Marshal scalar values into a slice
		if reflect.ValueOf(paramValue).Kind() != reflect.Slice {
			paramValue = []any{paramValue}
		}
		paramValue := paramValue.([]any)

		// Make a key-value pair for each param in the range
		for i := range paramValue {
			kv := keyValuePair{paramName, paramValue[i]}
			set = append(set, kv)
		}
		cartesianInputSets = append(cartesianInputSets, set)
	}

	// Produce the cartesian products passing in the input sets
	// Marshal each product (a slice of key-value pairs) back to a param set
	productCh := cartesian.Iter(cartesianInputSets...)
	for product := range productCh {
		pSet := make(map[string]any, len(product))
		for i := range product {
			kv := product[i].(keyValuePair)
			pSet[kv.k] = kv.v
		}
		products = append(products, pSet)
	}

	return products
}

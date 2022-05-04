package util

import "fmt"

func ToString(v any) string {
	return fmt.Sprint(v)
}

func ToInt(v any) int {
	var i int
	switch t := v.(type) {
	case int:
		i = t
	case int64:
		i = int(t)
	case float64:
		i = int(t)
	}
	return i
}

func ToFloat(v any) float64 {
	var f float64
	switch t := v.(type) {
	case int:
		f = float64(t)
	case int64:
		f = float64(t)
	case float64:
		f = t
	}
	return f
}

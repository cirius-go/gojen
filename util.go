package gojen

import (
	"encoding/json"
	"fmt"
)

func recordFrom[K comparable, V comparable](d []K, predicates ...func(d int, v K) V) map[K]V {
	res := make(map[K]V)
	for i, v := range d {
		if len(predicates) > 0 {
			predicate := predicates[0]
			if predicate != nil {
				res[v] = predicate(i, v)
			}
		} else {
			var t V
			res[v] = t
		}

	}

	return res
}

func contains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func toIntMap[T any](s []T) map[int]T {
	res := make(map[int]T)
	for i, v := range s {
		res[i] = v
	}

	return res
}

func mapVals[K comparable, V any](m map[K]V) []V {
	res := make([]V, 0, len(m))
	for _, v := range m {
		res = append(res, v)
	}
	return res
}

func findMapKey[K, V comparable](m map[K]V, v V) (K, bool) {
	var k K
	for k, val := range m {
		if val == v {
			return k, true
		}
	}
	return k, false
}

func mergeMaps[K comparable, V any](maps ...map[K]V) map[K]V {
	result := make(map[K]V)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

func filterMap[K comparable, V any](m map[K]V, keys []K) map[K]V {
	result := make(map[K]V)
	for k, v := range m {
		exists := false
		for _, key := range keys {
			if k == key {
				exists = true
				break
			}
		}

		if exists {
			result[k] = v
		}
	}

	return result
}

func printJSON(v any) {
	b, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		fmt.Println(err)
		fmt.Println(v)
		return
	}

	fmt.Println(string(b))
}

// Utility function to check confirmation input
func isConfirmed(input string) bool {
	switch input {
	case "y", "Y", "true", "1":
		return true
	default:
		return false
	}
}

func cloneSlice[T any](s []T) []T {
	if s == nil {
		return nil
	}
	var result = make([]T, len(s))
	copy(result, s)
	return result
}

func cloneMap[K comparable, V any](m map[K]V) map[K]V {
	result := make(map[K]V)
	for k, v := range m {
		result[k] = v
	}
	return result
}

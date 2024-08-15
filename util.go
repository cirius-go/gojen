package gojen

import (
	"encoding/json"
	"fmt"
)

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

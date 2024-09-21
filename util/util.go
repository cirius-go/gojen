package util

import (
	"encoding/json"
	"slices"
	"sort"
	"strings"
)

// RecordFrom creates a map from a slice of keys and a slice of predicates.
func RecordFrom[K comparable, V comparable](d []K, predicates ...func(i int, k K) V) map[K]V {
	res := make(map[K]V)

	for i, k := range d {
		var predicate func(d int, v K) V
		for i := range predicates {
			if predicates[i] != nil {
				predicate = predicates[i]
				break
			}
		}
		if predicate != nil {
			res[k] = predicate(i, k)
		} else {
			var t V
			res[k] = t
		}
	}

	return res
}

// SliceToMapExisting creates a map from a slice of keys.
func SliceToMapExisting[K comparable](d []K) MapExisting[K] {
	res := make(map[K]struct{})
	for _, v := range d {
		res[v] = struct{}{}
	}
	return res
}

// MapExisting is a map type that contains only the keys.
type MapExisting[K comparable] map[K]struct{}

func (m MapExisting[K]) Contains(k K) bool {
	_, ok := m[k]
	return ok
}

func (m MapExisting[K]) Keys() []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func (m MapExisting[K]) Add(ks ...K) {
	for _, k := range ks {
		m[k] = struct{}{}
	}
}

func (m MapExisting[K]) String() string {
	if m == nil {
		return ""
	}
	keys := m.Keys()
	b, _ := json.Marshal(&keys)
	return string(b)
}

// PFunc is a function type that takes a parameter.
type PFunc[P any] func(P) error

// PRFunc is a function type that takes a parameter and returns a value.
type PRFunc[P, R any] func(P) R

// MkSpace creates a string of spaces.
func MkSpace(n int) string {
	return strings.Repeat(" ", n)
}

// LoopStrMap loops through sorted keys of a string map.
func LoopStrMap[V any](m map[string]V, h func(string, V)) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		h(k, m[k])
	}
}

func SortSlice[K, V comparable](orders []K, v []V, fn func(V) K) []V {
	type pair struct {
		order K
		value V
		index int
	}

	pairs := make([]pair, len(v))
	for i, val := range v {
		pairs[i] = pair{order: fn(val), value: val, index: i}
	}

	sort.Slice(pairs, func(i, j int) bool {
		iIndex := slices.Index(orders, pairs[i].order)
		jIndex := slices.Index(orders, pairs[j].order)

		if iIndex == -1 && jIndex == -1 {
			// If both are not in orders, maintain original order
			return pairs[i].index < pairs[j].index
		}
		if iIndex == -1 {
			return false // i goes to the end
		}
		if jIndex == -1 {
			return true // j goes to the end
		}
		return iIndex < jIndex
	})

	result := make([]V, len(v))
	for i, p := range pairs {
		result[i] = p.value
	}

	return result
}

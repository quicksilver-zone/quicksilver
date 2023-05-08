package utils

import (
	"fmt"
	"sort"

	"golang.org/x/exp/constraints"
)

func Keys[V interface{}](in map[string]V) []string {
	out := make([]string, 0)

	for k := range in {
		out = append(out, k)
	}

	sort.Strings(out)

	return out
}

// SortSlice sorts a slice of type T elements that implement constraints.Ordered.
// Mutates input slice s.
func SortSlice[T constraints.Ordered](s []T) {
	sort.Slice(s, func(i, j int) bool {
		return s[i] < s[j]
	})
}

func Unique[V interface{}](in []V) []V {
	keys := make(map[string]struct{}, len(in))
	list := []V{}
	for _, entry := range in {
		key := fmt.Sprintf("%v", entry)
		if _, ok := keys[key]; !ok {
			keys[key] = struct{}{}
			list = append(list, entry)
		}
	}
	return list
}

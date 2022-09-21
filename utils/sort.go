package utils

import "sort"

func Keys[V interface{}](in map[string]V) []string {
	out := make([]string, 0)

	for k := range in {
		out = append(out, k)
	}

	sort.Strings(out)

	return out
}

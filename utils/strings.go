package utils

// create a map from a slice of strings for efficient lookup
func StringSliceToMap(in []string) map[string]bool {
	out := make(map[string]bool, len(in))
	for _, i := range in {
		out[i] = true
	}
	return out
}

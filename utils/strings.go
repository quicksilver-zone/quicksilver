package utils

// LengthPrefixString returns length-prefixed bytes representation of a string.
func LengthPrefixString(s string) []byte {
	bz := []byte(s)
	bzLen := len(bz)
	return append([]byte{byte(bzLen)}, bz...)
}

// create a map from a slice of strings for efficient lookup
func StringSliceToMap(in []string) map[string]bool {
	out := make(map[string]bool, len(in))
	for _, i := range in {
		out[i] = true
	}
	return out
}

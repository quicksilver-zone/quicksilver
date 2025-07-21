package utils

import (
	"fmt"
)

// helper function to convert error map to slice for multierr
func ErrorMapToSlice(errs map[string]error) []error {
	errList := make([]error, 0, len(errs)) // pre-allocate memory for the slice
	for _, err := range Keys(errs) {
		errList = append(errList, fmt.Errorf("%s: %w", err, errs[err]))
	}
	return errList
}

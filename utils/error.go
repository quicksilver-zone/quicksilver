package utils

import (
	"fmt"
)

// helper function to convert error map to slice for multierr
func ErrorMapToSlice(errs map[string]error) []error {
	var errList []error
	for _, err := range Keys(errs) {
		errList = append(errList, fmt.Errorf("%s: %w", err, errs[err]))
	}
	return errList
}

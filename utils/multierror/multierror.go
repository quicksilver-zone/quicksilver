package multierror

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// New returns an error aggregate using the given map.
func New(errs map[string]error) MultiError {
	return MultiError{errs}
}

// MultiError represents aggregated errors, contained in a map.
type MultiError struct {
	Errors map[string]error
}

func (e MultiError) Error() string {
	return e.details(0)
}

func (e MultiError) details(d int) string {
	str := "{"
	d++

	keys := make([]string, 0, len(e.Errors))
	for k := range e.Errors {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		str += indent(key, e.Errors[key], d)
	}
	d--
	str += fmt.Sprintf("\n%v}", indentString("  ", d))

	return str
}

func indent(k string, v error, d int) string {
	istr := indentString("  ", d)

	var typeErrors *MultiError
	if errors.As(v, &typeErrors) {
		return fmt.Sprintf("\n%v\"%v\": %v", istr, k, typeErrors.details(d))
	}

	return fmt.Sprintf("\n%v\"%v\": \"%v\"", istr, k, v)
}

func indentString(indent string, n int) string {
	return strings.Repeat(indent, n)
}

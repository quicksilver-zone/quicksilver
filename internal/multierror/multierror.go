package multierror

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// New returns an error aggregate using the given map.
func New(errors map[string]error) Errors {
	return Errors(errors)
}

// Error represents aggregated errors, contained in a map.
type Errors map[string]error

// Error implements the error interface.
func (e Errors) Error() string {
	return e.details(0)
}

// MarshalJSON implements json.Marshaler interface.
func (e Errors) MarshalJSON() ([]byte, error) {
	jsonErrors := make(map[string]interface{})
	for key, v := range e {
		switch err := v.(type) {
		case Errors:
			j, jerr := err.MarshalJSON()
			if jerr != nil {
				jsonErrors[key] = err.Error()
				continue
			}
			jsonErrors[key] = json.RawMessage(j)
		default:
			jsonErrors[key] = err.Error()
		}
	}

	return json.Marshal(jsonErrors)
}

func (e Errors) details(d int) string {
	str := "{"
	d++

	keys := make([]string, 0, len(e))
	for k := range e {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		str += indent(key, e[key], d)
	}
	d--
	str += fmt.Sprintf("\n%v}", indentString("  ", d))

	return str
}

func indent(k string, v error, d int) string {
	istr := indentString("  ", d)

	switch err := v.(type) {
	case Errors:
		return fmt.Sprintf("\n%v\"%v\": %v", istr, k, err.details(d))
	default:
		return fmt.Sprintf("\n%v\"%v\": \"%v\"", istr, k, v)
	}
}

func indentString(indent string, n int) string {
	return strings.Repeat(indent, n)
}

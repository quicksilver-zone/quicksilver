package types

import (
	fmt "fmt"
	"sort"
	"strings"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/airdrop module sentinel errors
var (
	ErrZoneDropNotFound    = sdkerrors.Register(ModuleName, 1, "zone airdrop not found")
	ErrClaimRecordNotFound = sdkerrors.Register(ModuleName, 2, "claim record not found")
	ErrUnknownStatus       = sdkerrors.Register(ModuleName, 3, "unknown status")
)

// NewMultiError returns an error aggregate using the given map.
func NewMultiError(errors map[string]error) Errors {
	return Errors{errors}
}

// Error represents aggregated errors, contained in a map.
type Errors struct {
	Errors map[string]error
}

func (e Errors) Error() string {
	return e.details(0)
}

func (e Errors) details(d int) string {
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

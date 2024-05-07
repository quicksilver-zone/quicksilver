package osmoutils

import (
	"fmt"

	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
)

// ValidateAddressList validates a slice of addresses.
//
// Parameters:
// - i: The parameter to validate.
//
// Returns:
// - An error if any of the strings are not addresses
func ValidateAddressList(i interface{}) error {
	whitelist, ok := i.([]string)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for _, a := range whitelist {
		if _, err := addressutils.AccAddressFromBech32(a, ""); err != nil {
			return fmt.Errorf("invalid address")
		}
	}

	return nil
}

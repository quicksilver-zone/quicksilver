package types

import (
	"fmt"
	"strings"
)

func ValidateEpochIdentifierInterface(i any) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	err := ValidateEpochIdentifierString(v)

	return err
}

func ValidateEpochIdentifierString(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return fmt.Errorf("blank epoch identifier: %s", s)
	}
	return nil
}

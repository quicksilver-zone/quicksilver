package types

import fmt "fmt"

type (
	ErrInvalidParameter struct {
		Type interface{}
	}
)

func (e ErrInvalidParameter) Error() string {
	return fmt.Sprintf("invalid parameter type: %T", e.Type)
}

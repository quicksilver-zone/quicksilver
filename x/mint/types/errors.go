package types

import "fmt"

type (
	ErrInvalidParameter struct {
		Type interface{}
	}
)

func (e ErrInvalidParameter) Error() string {
	return fmt.Sprintf("invalid parameter type: %T", e.Type)
}

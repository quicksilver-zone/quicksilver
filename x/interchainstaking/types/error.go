package types

import (
	"errors"
	"fmt"
)

type (
	ErrNoRegisteredZoneForChainId struct {
		Id string
	}
)

func (e ErrNoRegisteredZoneForChainId) Error() string {
	return fmt.Sprintf("no registered zone for chain id: %s", e.Id)
}

var (
	ErrCoinAmountNil              = errors.New("coin amount is nil")
	ErrValidatorAlreadyInDenyList = errors.New("validator already in deny list")
	ErrEmptyChainID               = errors.New("chain-id cannot be empty")
	ErrErrorHappened              = errors.New("an error happened")
	ErrUnexpectedNegative         = errors.New("unexpected negative value")
)

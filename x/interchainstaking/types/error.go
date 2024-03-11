package types

import (
	"errors"
)

var ErrCoinAmountNil = errors.New("coin amount is nil")
var ErrValidatorAlreadyInDenyList = errors.New("validator already in deny list")

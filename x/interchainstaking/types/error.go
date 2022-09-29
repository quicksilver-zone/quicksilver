package types

import (
	"errors"
)

var (
	ErrInvalidVersion = errors.New("invalid version")
	ErrMaxChannels    = errors.New("max channels exceeded")
	ErrCoinAmountNil  = errors.New("coin amount is nil")
)

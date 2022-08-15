package types

import (
	"errors"
)

var (
	ErrInvalidVersion = errors.New("invalid version")
	ErrMaxChannels    = errors.New("max channels exceeded")
)

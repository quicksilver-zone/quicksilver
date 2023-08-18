package types

import (
	"errors"

	sdkerrors "cosmossdk.io/errors"
)

var (
	ErrCoinAmountNil           = errors.New("coin amount is nil")
	ErrInvalidSubzoneAuthority = sdkerrors.Register(ModuleName, 1, "invalid authority for subzone")
)

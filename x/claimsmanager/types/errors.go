package types

import (
	fmt "fmt"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/claimsmanager module sentinel errors
var (
	ErrUndefinedAttribute   = sdkerrors.Register(ModuleName, 1, "expected attribute not defined")
	ErrNegativeAttribute    = sdkerrors.Register(ModuleName, 2, "expected attribute must not be negative")
	ErrNotPositive          = sdkerrors.Register(ModuleName, 3, "expected attribute must be positive")
	ErrClaimTypeOutOfBounds = sdkerrors.Register(ModuleName, 4, fmt.Sprintf("invalid claim type, expects range [1-%d]", len(ClaimType_value)-1))
)

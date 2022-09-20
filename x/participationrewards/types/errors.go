package types

import (
	fmt "fmt"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/participationrewards module sentinel errors
var (
	ErrUndefinedAttribute            = sdkerrors.Register(ModuleName, 1, "expected attribute not defined")
	ErrNegativeAttribute             = sdkerrors.Register(ModuleName, 2, "expected attribute must not be negative")
	ErrNegativeDistributionRatio     = sdkerrors.Register(ModuleName, 3, "distribution ratio must not be negative")
	ErrInvalidTotalProportions       = sdkerrors.Register(ModuleName, 4, "total distribution proportions must be 1.0")
	ErrNotPositive                   = sdkerrors.Register(ModuleName, 5, "expected attribute must be positive")
	ErrUnknownProtocolDataType       = sdkerrors.Register(ModuleName, 6, "unknown protocol data type")
	ErrUnimplementedProtocolDataType = sdkerrors.Register(ModuleName, 7, "protocol data type not implemented")
	ErrClaimTypeOutOfBounds          = sdkerrors.Register(ModuleName, 8, fmt.Sprintf("invalid claim type, expects range [1-%d]", len(ClaimType_value)-1))
	ErrNothingToAllocate             = sdkerrors.Register(ModuleName, 9, "balance is zero, nothing to allocate")
)

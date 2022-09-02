package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/participationrewards module sentinel errors
var (
	ErrUndefinedAttribute        = sdkerrors.Register(ModuleName, 1, "expected attribute not defined")
	ErrNegativeDistributionRatio = sdkerrors.Register(ModuleName, 2, "distribution ratio must not be negative")
	ErrInvalidTotalProportions   = sdkerrors.Register(ModuleName, 3, "total distribution proportions must be 1.0")
	ErrSliceLengthMismatch       = sdkerrors.Register(ModuleName, 4, "expected data / key / proofops fields to be equal length")
)

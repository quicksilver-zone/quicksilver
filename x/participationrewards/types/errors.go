package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/participationrewards module sentinel errors
var (
	ErrUndefinedAttribute        = sdkerrors.Register(ModuleName, 1, "expected attribute not defined")
	ErrNegativeAttribute         = sdkerrors.Register(ModuleName, 2, "expected attribute must not be negative")
	ErrNegativeDistributionRatio = sdkerrors.Register(ModuleName, 3, "distribution ratio must not be negative")
	ErrInvalidTotalProportions   = sdkerrors.Register(ModuleName, 4, "total distribution proportions must be 1.0")
)

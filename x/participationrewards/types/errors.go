package types

import sdkioerrors "cosmossdk.io/errors"

// x/participationrewards module sentinel errors.
var (
	ErrUndefinedAttribute            = sdkioerrors.Register(ModuleName, 1, "expected attribute not defined")
	ErrNegativeAttribute             = sdkioerrors.Register(ModuleName, 2, "expected attribute must not be negative")
	ErrNegativeDistributionRatio     = sdkioerrors.Register(ModuleName, 3, "distribution ratio must not be negative")
	ErrInvalidTotalProportions       = sdkioerrors.Register(ModuleName, 4, "total distribution proportions must be 1.0")
	ErrNotPositive                   = sdkioerrors.Register(ModuleName, 5, "expected attribute must be positive")
	ErrUnknownProtocolDataType       = sdkioerrors.Register(ModuleName, 6, "unknown protocol data type")
	ErrUnimplementedProtocolDataType = sdkioerrors.Register(ModuleName, 7, "protocol data type not implemented")
	ErrNothingToAllocate             = sdkioerrors.Register(ModuleName, 9, "balance is zero, nothing to allocate")
	ErrInvalidAssetName              = sdkioerrors.Register(ModuleName, 10, "invalid ibc asset name")
	ErrInvalidChainID                = sdkioerrors.Register(ModuleName, 11, "invalid chain id")
	ErrInvalidDenom                  = sdkioerrors.Register(ModuleName, 12, "invalid denom")
	ErrCannotUnmarshal               = sdkioerrors.Register(ModuleName, 13, "unable to unmarshal")
)

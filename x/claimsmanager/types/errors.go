package types

import (
	"fmt"

	sdkioerrors "cosmossdk.io/errors"
)

// x/claimsmanager module sentinel errors.
var (
	ErrUndefinedAttribute   = sdkioerrors.Register(ModuleName, 1, "expected attribute not defined")
	ErrNegativeAttribute    = sdkioerrors.Register(ModuleName, 2, "expected attribute must not be negative")
	ErrNotPositive          = sdkioerrors.Register(ModuleName, 3, "expected attribute must be positive")
	ErrClaimTypeOutOfBounds = sdkioerrors.Register(ModuleName, 4, fmt.Sprintf("invalid claim type, expects range [1-%d]", len(ClaimType_value)-1))
)

package types

import (
	"fmt"

	sdkioerrors "cosmossdk.io/errors"
)

// x/airdrop module sentinel errors.
var (
	ErrZoneDropNotFound     = sdkioerrors.Register(ModuleName, 1, "zone airdrop not found")
	ErrClaimRecordNotFound  = sdkioerrors.Register(ModuleName, 2, "claim record not found")
	ErrUnknownStatus        = sdkioerrors.Register(ModuleName, 3, "unknown status")
	ErrUndefinedAttribute   = sdkioerrors.Register(ModuleName, 4, "expected attribute not defined")
	ErrInvalidDuration      = sdkioerrors.Register(ModuleName, 5, "invalid duration")
	ErrActionOutOfBounds    = sdkioerrors.Register(ModuleName, 6, fmt.Sprintf("invalid action, expects range [1-%d]", len(Action_value)-1))
	ErrActionWeights        = sdkioerrors.Register(ModuleName, 7, "sum of action weights must be 1.0")
	ErrDuplicateZoneDrop    = sdkioerrors.Register(ModuleName, 8, "duplicate zone drop")
	ErrDuplicateClaimRecord = sdkioerrors.Register(ModuleName, 9, "duplicate claim record")
	ErrAllocationExceeded   = sdkioerrors.Register(ModuleName, 10, "claim records allocations exceed zone drop allocation")
	ErrNoClaimRecords       = sdkioerrors.Register(ModuleName, 11, "no claim records for zone drop")
	ErrZoneDropExpired      = sdkioerrors.Register(ModuleName, 12, "nothing to claim, this zone drop has expired")
	ErrActionCompleted      = sdkioerrors.Register(ModuleName, 13, "nothing to claim, action already completed")
	ErrNegativeAttribute    = sdkioerrors.Register(ModuleName, 14, "expected attribute must not be negative")
)

package types

import (
	"fmt"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/airdrop module sentinel errors
var (
	ErrZoneDropNotFound     = sdkerrors.Register(ModuleName, 1, "zone airdrop not found")
	ErrClaimRecordNotFound  = sdkerrors.Register(ModuleName, 2, "claim record not found")
	ErrUnknownStatus        = sdkerrors.Register(ModuleName, 3, "unknown status")
	ErrUndefinedAttribute   = sdkerrors.Register(ModuleName, 4, "expected attribute not defined")
	ErrInvalidDuration      = sdkerrors.Register(ModuleName, 5, "invalid duration")
	ErrActionOutOfBounds    = sdkerrors.Register(ModuleName, 6, fmt.Sprintf("invalid action, expects range [1-%d]", len(Action_value)-1))
	ErrActionWeights        = sdkerrors.Register(ModuleName, 7, "sum of action weights must be 1.0")
	ErrDuplicateZoneDrop    = sdkerrors.Register(ModuleName, 8, "duplicate zone drop")
	ErrDuplicateClaimRecord = sdkerrors.Register(ModuleName, 9, "duplicate claim record")
	ErrAllocationExceeded   = sdkerrors.Register(ModuleName, 10, "claim records allocations exceed zone drop allocation")
	ErrNoClaimRecords       = sdkerrors.Register(ModuleName, 11, "no claim records for zone drop")
	ErrZoneDropExpired      = sdkerrors.Register(ModuleName, 12, "nothing to claim, this zone drop has expired")
	ErrActionCompleted      = sdkerrors.Register(ModuleName, 13, "nothing to claim, action already completed")
	ErrNegativeAttribute    = sdkerrors.Register(ModuleName, 14, "expected attribute must not be negative")
)

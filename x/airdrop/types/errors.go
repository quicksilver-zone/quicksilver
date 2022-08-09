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
	ErrActionOutOfBounds    = sdkerrors.Register(ModuleName, 6, fmt.Sprintf("invalid action, expects range [0-%d]", len(Action_value)-1))
	ErrDuplicateZoneDrop    = sdkerrors.Register(ModuleName, 7, "duplicate zone drop")
	ErrDuplicateClaimRecord = sdkerrors.Register(ModuleName, 8, "duplicate claim record")
	ErrAllocationExceeded   = sdkerrors.Register(ModuleName, 9, "claim records allocations exceed zone drop allocation")
)

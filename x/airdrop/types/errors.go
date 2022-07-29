package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/airdrop module sentinel errors
var (
	ErrZoneDropNotFound    = sdkerrors.Register(ModuleName, 1, "zone airdrop not found")
	ErrClaimRecordNotFound = sdkerrors.Register(ModuleName, 2, "claim record not found")
)

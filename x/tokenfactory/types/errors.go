package types

// DONTCOVER

import (
	fmt "fmt"

	sdkioerrors "cosmossdk.io/errors"
)

// x/tokenfactory module sentinel errors.
var (
	ErrDenomExists              = sdkioerrors.Register(ModuleName, 2, "attempting to create a denom that already exists (has bank metadata)")
	ErrUnauthorized             = sdkioerrors.Register(ModuleName, 3, "unauthorized account")
	ErrInvalidDenom             = sdkioerrors.Register(ModuleName, 4, "invalid denom")
	ErrInvalidCreator           = sdkioerrors.Register(ModuleName, 5, "invalid creator")
	ErrInvalidAuthorityMetadata = sdkioerrors.Register(ModuleName, 6, "invalid authority metadata")
	ErrInvalidGenesis           = sdkioerrors.Register(ModuleName, 7, "invalid genesis")
	ErrSubdenomTooLong          = sdkioerrors.Register(ModuleName, 8, fmt.Sprintf("subdenom too long, max length is %d bytes", MaxSubdenomLength))
	ErrCreatorTooLong           = sdkioerrors.Register(ModuleName, 9, fmt.Sprintf("creator too long, max length is %d bytes", MaxCreatorLength))
	ErrDenomDoesNotExist        = sdkioerrors.Register(ModuleName, 10, "denom does not exist")
)

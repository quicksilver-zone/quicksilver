package claimsmanager

import (
	sdkioerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/ingenuity-build/quicksilver/x/claimsmanager/keeper"
)

// NewHandler returns a handler for claimsmanager module messages
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		_ = ctx.WithEventManager(sdk.NewEventManager())

		return nil, sdkioerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized claimsmanager message type: %T", msg)
	}
}

package participationrewards

import (
	sdkioerrors "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/quicksilver-zone/quicksilver/v7/x/participationrewards/keeper"
	"github.com/quicksilver-zone/quicksilver/v7/x/participationrewards/types"
)

func NewProposalHandler(k *keeper.Keeper) govv1beta1.Handler {
	return func(ctx sdk.Context, content govv1beta1.Content) error {
		switch c := content.(type) {
		case *types.AddProtocolDataProposal:
			return keeper.HandleAddProtocolDataProposal(ctx, k, c)

		default:
			return sdkioerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized participationrewards proposal content type: %T", c)
		}
	}
}

package interchainstaking

import (
	sdkioerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// NewHandler returns a handler for interchainstaking module messages
func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgRequestRedemption:
			res, err := msgServer.RequestRedemption(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgSignalIntent:
			res, err := msgServer.SignalIntent(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			return nil, sdkioerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized interchainstaking message type: %T", msg)
		}
	}
}

func NewProposalHandler(k keeper.Keeper) govv1beta1.Handler {
	return func(ctx sdk.Context, content govv1beta1.Content) error {
		switch c := content.(type) {
		case *types.RegisterZoneProposal:
			return k.HandleRegisterZoneProposal(ctx, c)
		case *types.UpdateZoneProposal:
			return k.HandleUpdateZoneProposal(ctx, c)

		default:
			return sdkioerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized interchainstaking proposal content type: %T", c)
		}
	}
}

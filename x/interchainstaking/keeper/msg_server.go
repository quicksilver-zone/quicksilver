package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) RegisterZone(goCtx context.Context, msg *types.MsgRegisterZone) (*types.MsgRegisterZoneResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	zone := types.RegisteredZone{Identifier: msg.Identifier, ChainId: msg.ChainId, LocalDenom: msg.LocalDenom, RemoteDenom: msg.RemoteDenom}

	// generate new deposit address

	// generate delegate addresses

	k.SetRegisteredZone(ctx, zone)

	// if err := k.IsSendEnabledCoins(ctx, msg.Amount...); err != nil {
	// 	return nil, err
	// }

	// from, err := sdk.AccAddressFromBech32(msg.FromAddress)
	// if err != nil {
	// 	return nil, err
	// }
	// to, err := sdk.AccAddressFromBech32(msg.ToAddress)
	// if err != nil {
	// 	return nil, err
	// }

	// if k.BlockedAddr(to) {
	// 	return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive funds", msg.ToAddress)
	// }

	// err = k.SendCoins(ctx, from, to, msg.Amount)
	// if err != nil {
	// 	return nil, err
	// }

	// defer func() {
	// 	for _, a := range msg.Amount {
	// 		if a.Amount.IsInt64() {
	// 			telemetry.SetGaugeWithLabels(
	// 				[]string{"tx", "msg", "send"},
	// 				float32(a.Amount.Int64()),
	// 				[]metrics.Label{telemetry.NewLabel("denom", a.Denom)},
	// 			)
	// 		}
	// 	}
	// }()

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeRegisterZone,
			sdk.NewAttribute(types.AttributeKeyChainId, msg.ChainId),
		),
	})

	return &types.MsgRegisterZoneResponse{}, nil
}

package keeper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	sdkioerrors "cosmossdk.io/errors"
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/quicksilver-zone/quicksilver/x/airdrop/types"
)

type msgServer struct {
	*Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper *Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) Claim(goCtx context.Context, msg *types.MsgClaim) (*types.MsgClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	action := types.Action(msg.Action)

	amount, err := k.Keeper.Claim(ctx, msg.ChainId, action, msg.Address, msg.Proofs)
	if err != nil {
		return nil, err
	}

	return &types.MsgClaimResponse{Amount: amount}, nil
}

func (k msgServer) IncentivePoolSpend(goCtx context.Context, msg *types.MsgIncentivePoolSpend) (*types.MsgIncentivePoolSpendResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.GetAuthority() != msg.Authority {
		return nil, sdkioerrors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.GetAuthority(), msg.Authority)
	}

	to, err := sdk.AccAddressFromBech32(msg.ToAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", err)
	}

	if err := k.bankKeeper.IsSendEnabledCoins(ctx, msg.Amount...); err != nil {
		return nil, err
	}

	if k.bankKeeper.BlockedAddr(to) {
		return nil, sdkioerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive funds", msg.ToAddress)
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, to, msg.Amount)
	if err != nil {
		return nil, err
	}

	defer func() {
		for _, a := range msg.Amount {
			if a.Amount.IsInt64() {
				telemetry.SetGaugeWithLabels(
					[]string{"tx", "msg", "send"},
					float32(a.Amount.Int64()),
					[]metrics.Label{telemetry.NewLabel("denom", a.Denom)},
				)
			}
		}
	}()

	return &types.MsgIncentivePoolSpendResponse{}, nil
}

func (k msgServer) RegisterZoneDrop(goCtx context.Context, msg *types.MsgRegisterZoneDrop) (*types.MsgRegisterZoneDropResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if k.GetAuthority() != msg.Authority {
		return &types.MsgRegisterZoneDropResponse{}, sdkioerrors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.GetAuthority(), msg.Authority)
	}

	_, found := k.icsKeeper.GetZone(ctx, msg.ZoneDrop.ChainId)
	if !found {
		return &types.MsgRegisterZoneDropResponse{}, fmt.Errorf("zone not found, %q", msg.ZoneDrop.ChainId)
	}

	if msg.ZoneDrop.StartTime.Before(ctx.BlockTime()) {
		return &types.MsgRegisterZoneDropResponse{}, errors.New("zone airdrop already started")
	}

	// decompress claim records
	crsb, err := types.Decompress(msg.ClaimRecords)
	if err != nil {
		return &types.MsgRegisterZoneDropResponse{}, err
	}

	// unmarshal json
	var crs types.ClaimRecords
	if err := json.Unmarshal(crsb, &crs); err != nil {
		return &types.MsgRegisterZoneDropResponse{}, err
	}

	// process ZoneDrop
	k.SetZoneDrop(ctx, *msg.ZoneDrop)
	for i, cr := range crs {
		if err := k.SetClaimRecord(ctx, cr); err != nil {
			return &types.MsgRegisterZoneDropResponse{}, fmt.Errorf("invalid zonedrop proposal claim record [%d]: %w", i, err)
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
		sdk.NewEvent(
			types.EventTypeRegisterZoneDrop,
			sdk.NewAttribute(types.AttributeKeyZoneID, msg.ZoneDrop.ChainId),
		),
	})

	return &types.MsgRegisterZoneDropResponse{}, nil
}

func (k msgServer) UpdateParams(goCtx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if k.authority != req.Authority {
		return nil, sdkioerrors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, req.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.SetParams(ctx, req.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}

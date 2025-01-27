package keeper

import (
	"context"

	"github.com/armon/go-metrics"

	sdkioerrors "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/quicksilver-zone/quicksilver/x/supply/types"
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

func (k msgServer) IncentivePoolSpend(goCtx context.Context, msg *types.MsgIncentivePoolSpend) (*types.MsgIncentivePoolSpendResponse, error) {
	if k.govAuthority != msg.Authority {
		return nil, sdkioerrors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.govAuthority, msg.Authority)
	}

	to, err := sdk.AccAddressFromBech32(msg.ToAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", err)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.bankKeeper.IsSendEnabledCoins(ctx, msg.Amount...); err != nil {
		return nil, err
	}

	if k.bankKeeper.BlockedAddr(to) {
		return nil, sdkioerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive funds", msg.ToAddress)
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.AirdropAccount, to, msg.Amount)
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

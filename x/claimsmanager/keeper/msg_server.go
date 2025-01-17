package keeper

import (
	"context"
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
)

type msgServer struct {
	*Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: &keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) SubmitClaim(goCtx context.Context, msg *types.MsgSubmitClaim) (*types.MsgSubmitClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.GetClaimsEnabled(ctx) {
		return nil, errors.New("claims currently disabled")
	}
	// fetch zone
	zone, ok := k.icsKeeper.GetZone(ctx, msg.Zone)
	if !ok {
		return nil, fmt.Errorf("invalid zone, chain id %q not found", msg.Zone)
	}
	var pd types.ProtocolData
	pd, ok = k.GetProtocolData(ctx, types.ProtocolDataTypeConnection, msg.SrcZone)
	if !ok {
		return nil, fmt.Errorf("unable to obtain connection protocol data for %q", msg.SrcZone)
	}

	// protocol data
	iConnectionData, err := types.UnmarshalProtocolData(types.ProtocolDataTypeConnection, pd.Data)
	if err != nil {
		k.Logger(ctx).Error("SubmitClaim: error unmarshalling protocol data")
		return nil, fmt.Errorf("unable to unmarshal connection protocol data for %q", msg.SrcZone)
	}
	connectionData, ok := iConnectionData.(*types.ConnectionProtocolData)
	if !ok {
		return nil, fmt.Errorf("unable to cast connection protocol data for %q", msg.SrcZone)
	}

	for i, proof := range msg.Proofs {
		pl := fmt.Sprintf("Proof [%d]", i)

		if proof.Height != connectionData.LastEpoch {
			return nil, fmt.Errorf(
				"invalid claim for last epoch, %s expected height %d, got %d",
				pl,
				connectionData.LastEpoch,
				proof.Height,
			)
		}

		// if we are claiming against Quicksilver, use the SelfProofOpsFn.
		if msg.SrcZone == ctx.ChainID() {
			if err := k.ValidateSelfProofOps(
				ctx,
				k,
				"epoch",
				proof.ProofType,
				proof.Key,
				proof.Data,
				proof.ProofOps,
			); err != nil {
				return nil, fmt.Errorf("%s: %w", pl, err)
			}
		} else {
			if err := k.ValidateProofOps(
				ctx,
				k.IBCKeeper,
				connectionData.ConnectionID,
				connectionData.ChainID,
				proof.Height,
				proof.ProofType,
				proof.Key,
				proof.Data,
				proof.ProofOps,
			); err != nil {
				return nil, fmt.Errorf("%s: %w", pl, err)
			}
		}
	}

	// if we get here all data was validated; verifyClaim will write the claim to the correct store.
	if mod, ok := k.PrSubmodules[msg.ClaimType]; ok {
		// vertifyClaim needs to return the amount!
		amount, err := mod.ValidateClaim(ctx, k.Keeper, msg)
		if err != nil {
			return nil, fmt.Errorf("claim validation failed: %w", err)
		}
		claim := types.NewClaim(msg.UserAddress, zone.ChainId, msg.ClaimType, msg.SrcZone, amount)
		k.SetClaim(ctx, &claim)
	}

	return &types.MsgSubmitClaimResponse{}, nil
}

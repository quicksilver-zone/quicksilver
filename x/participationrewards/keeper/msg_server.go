package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/utils"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
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

// SubmitClaim is used to verify, by proof, that the given user address has a
// claim against the given asset type for the given zone.
func (k msgServer) SubmitClaim(goCtx context.Context, msg *types.MsgSubmitClaim) (*types.MsgSubmitClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// fetch zone
	zone, ok := k.icsKeeper.GetZone(ctx, msg.Zone)
	// this is fine for launch, but we cannot always guarantee the zone we are querying is a registered zone.
	// add last_epoch_height to connection protocol, so we can fetch this epochly.

	if !ok {
		return nil, fmt.Errorf("invalid zone, chain id \"%s\" not found", msg.Zone)
	}

	for i, proof := range msg.Proofs {
		pl := fmt.Sprintf("Proof [%d]", i)

		if proof.Height != zone.LastEpochHeight {
			return nil, fmt.Errorf(
				"invalid claim for last epoch, %s expected height %d, got %d",
				pl,
				zone.LastEpochHeight,
				proof.Height,
			)
		}

		if err := utils.ValidateProofOps(
			ctx,
			&k.icsKeeper.IBCKeeper,
			zone.ConnectionId,
			zone.ChainId,
			proof.Height,
			proof.ProofType,
			proof.Key,
			proof.Data,
			proof.ProofOps,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", pl, err)
		}
	}

	// if we get here all data was validated; verifyClaim will write the claim to the correct store.
	if mod, ok := k.prSubmodules[msg.ClaimType]; ok {
		// vertifyClaim needs to return the amount!
		amount, err := mod.ValidateClaim(ctx, k.Keeper, msg)
		if err != nil {
			return nil, fmt.Errorf("claim validation failed: %v", err)
		}
		claim := k.NewClaim(ctx, msg.UserAddress, zone.ChainId, msg.ClaimType, msg.SrcZone, amount)
		k.SetClaim(ctx, &claim)
	}

	return &types.MsgSubmitClaimResponse{}, nil
}

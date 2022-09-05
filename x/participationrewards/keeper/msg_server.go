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
			"lockup",
			proof.Key,
			proof.Data,
			proof.ProofOps,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", pl, err)
		}
	}

	// if we get here all data was validated; verifyClaim will write the claim to the correct store.
	if mod, ok := k.prSubmodules[msg.ProofType]; ok {
		if err := mod.VerifyClaim(ctx, k.Keeper, msg); err != nil {
			return nil, fmt.Errorf("claim verification failed: %v", err)
		}
	}

	return &types.MsgSubmitClaimResponse{}, nil
}

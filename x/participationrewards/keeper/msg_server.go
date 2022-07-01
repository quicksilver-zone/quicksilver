package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/utils"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

// SubmitClaim is used to verify, by proof, that the given user address has a
// claim against the given asset type for the given zone.
func (k *Keeper) SubmitClaim(goCtx context.Context, msg *types.MsgSubmitClaim) (*types.MsgSubmitClaimResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	// fetch zone
	zone, ok := k.icsKeeper.GetRegisteredZoneInfo(ctx, msg.Zone)
	if !ok {
		return nil, fmt.Errorf("invalid zone, chain id \"%s\" not found", msg.Zone)
	}

	if msg.Height != zone.LastEpochHeight {
		return nil, fmt.Errorf("invalid claim for last epoch, expected height %d, got %d", zone.LastEpochHeight, msg.Height)
	}

	height := msg.Height
	claimCount := len(msg.Data) // validate basic should check that the length of these fields is equal.
	for claimIdx := 0; claimIdx < claimCount; claimIdx++ {
		if err := utils.ValidateProofOps(ctx, &k.icsKeeper.IBCKeeper, zone.ConnectionId, zone.ChainId, height, "lockup", msg.Key[claimIdx], msg.Data[claimIdx], msg.ProofOps[claimIdx]); err != nil {
			return nil, err
		}
	}
	// if we get here all data was validated; verifyClaim will write the claim to the correct store.
	if mod, ok := k.prSubmodules[msg.ProofType]; ok {
		if err := mod.VerifyClaim(ctx, k, msg); err != nil {
			return nil, fmt.Errorf("claim verification failed", err)
		}
	}

	return &types.MsgSubmitClaimResponse{}, nil
}

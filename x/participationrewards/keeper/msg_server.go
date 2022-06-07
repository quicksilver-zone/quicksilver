package keeper

import (
	"context"
	"fmt"

	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

// SubmitClaim is used to verify, by proof, that the given user address has a
// claim against the given asset type for the given zone.
func (k *Keeper) SubmitClaim(goCtx context.Context, msg *types.MsgSubmitClaim) (*types.MsgSubmitClaimResponse, error) {
	// TODO: implement

	/*ctx := sdk.UnwrapSDKContext(goCtx)

	// get zone
	zone, ok := k.icsKeeper.GetRegisteredZoneInfo(ctx, msg.Zone)
	if !ok {
		return nil, fmt.Errorf("invalid zone, chain id \"%s\" not found", msg.Zone)
	}*/

	return nil, fmt.Errorf("not implemented")
}

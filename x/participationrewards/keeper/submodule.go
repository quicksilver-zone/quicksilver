package keeper

import (
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

// Submodule defines the interface for for tracking off-chain qAssets with
// regards to participation rewards claims.
type Submodule interface {
	Hooks(ctx sdk.Context, keeper *Keeper)
	ValidateClaim(ctx sdk.Context, k *Keeper, msg *types.MsgSubmitClaim) (math.Int, error)
}

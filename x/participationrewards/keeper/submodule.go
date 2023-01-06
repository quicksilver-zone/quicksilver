package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

// Submodule defines the interface for for tracking off-chain qAssets with
// regards to participation rewards claims.
type Submodule interface {
	Hooks(ctx sdk.Context, keeper Keeper)
	IsActive() bool
	IsReady() bool
	ValidateClaim(ctx sdk.Context, k *Keeper, msg *types.MsgSubmitClaim) (uint64, error)
}

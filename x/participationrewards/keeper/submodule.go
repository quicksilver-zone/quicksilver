package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

type Submodule interface {
	Hooks(ctx sdk.Context, keeper Keeper)
	IsActive() bool
	IsReady() bool
	VerifyClaim(ctx sdk.Context, k *Keeper, msg *types.MsgSubmitClaim) error
}

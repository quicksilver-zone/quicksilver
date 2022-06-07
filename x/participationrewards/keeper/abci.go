package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

// BeginBlocker of participationrewards module
func (k Keeper) BeginBlocker(ctx sdk.Context) {
	// hartbeat logger (for dev & debugging)
	if ctx.BlockHeight()%int64(10) == 0 {
		k.Logger(ctx).Info("up and running")
		k.Logger(ctx).Info(
			"module account",
			"account", k.accountKeeper.GetModuleAccount(ctx, types.ModuleName),
			"address", k.accountKeeper.GetModuleAddress(types.ModuleName),
		)
	}
}

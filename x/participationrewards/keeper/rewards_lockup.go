package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

func (k Keeper) allocateLockupRewards(ctx sdk.Context, allocation math.Int) error {
	k.Logger(ctx).Info("allocateLockupRewards", "allocation", allocation)

	// allocate staking incentives into fee collector account to be moved to on next begin blocker by staking module
	coins := sdk.NewCoins(
		sdk.NewCoin(
			k.stakingKeeper.BondDenom(ctx),
			allocation,
		),
	)

	return k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, coins)
}

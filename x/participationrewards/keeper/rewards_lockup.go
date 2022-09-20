package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

func (k Keeper) allocateLockupRewards(ctx sdk.Context, allocation sdk.Int) error {
	k.Logger(ctx).Info("allocateLockupRewards", "allocation", allocation)

	// allocate staking incentives into fee collector account to be moved to on next begin blocker by staking module
	coins := sdk.NewCoins(
		sdk.NewCoin(
			k.stakingKeeper.BondDenom(ctx),
			allocation,
		),
	)
	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, coins); err != nil {
		return err
	}

	return nil
}

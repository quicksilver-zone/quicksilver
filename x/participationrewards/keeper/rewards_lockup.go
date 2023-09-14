package keeper

import (
	"errors"

	"github.com/quicksilver-zone/quicksilver/x/participationrewards/types"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) AllocateLockupRewards(ctx sdk.Context, allocation math.Int) error {
	if allocation.IsNegative() {
		k.Logger(ctx).Error("invalid allocation requested", "allocation", allocation)
		return errors.New("invalid allocation")
	}
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

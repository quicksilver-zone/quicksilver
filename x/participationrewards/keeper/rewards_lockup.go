package keeper

import (
	"errors"
	"fmt"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/v7/x/participationrewards/types"
)

func (k Keeper) AllocateLockupRewards(ctx sdk.Context, allocation math.Int) error {
	if allocation.IsNegative() {
		k.Logger(ctx).Error("invalid allocation requested", "allocation", allocation)
		return errors.New("invalid allocation")
	}
	k.Logger(ctx).Info("allocateLockupRewards", "allocation", allocation)

	bondDenom, err := k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return fmt.Errorf("failed to get bond denom: %w", err)
	}
	// allocate staking incentives into fee collector account to be moved to on next begin blocker by staking module
	coins := sdk.NewCoins(
		sdk.NewCoin(
			bondDenom,
			allocation,
		),
	)
	return k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, coins)
}

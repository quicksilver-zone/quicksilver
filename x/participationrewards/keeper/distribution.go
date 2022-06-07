package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

type rewardsAllocation struct {
	ValidatorSelection sdk.Coins
	Holdings           sdk.Coins
	Lockup             sdk.Coins
}

func (k Keeper) getRewardsAllocations(ctx sdk.Context) rewardsAllocation {
	var allocation rewardsAllocation

	moduleAddress := k.accountKeeper.GetModuleAddress(types.ModuleName)
	moduleBalances := k.bankKeeper.GetAllBalances(ctx, moduleAddress)

	k.Logger(ctx).Info("module account", "address", moduleAddress, "balances", moduleBalances)

	if moduleBalances.Empty() {
		k.Logger(ctx).Info("nothing to distribute...")

		return allocation
	}

	// get distribution proportions (params)
	params := k.GetParams(ctx)
	k.Logger(ctx).Info("module parameters", "params", params)

	// split participation rewards allocations
	allocation.ValidatorSelection = sdk.NewCoins(
		k.GetAllocation(
			ctx,
			moduleBalances[0],
			params.DistributionProportions.ValidatorSelectionAllocation,
		),
	)
	allocation.Holdings = sdk.NewCoins(
		k.GetAllocation(
			ctx,
			moduleBalances[0],
			params.DistributionProportions.HoldingsAllocation,
		),
	)
	allocation.Lockup = sdk.NewCoins(
		k.GetAllocation(
			ctx,
			moduleBalances[0],
			params.DistributionProportions.LockupAllocation,
		),
	)

	// use sum to check total distribution to collect and allocate dust
	total := moduleBalances[0]
	sum := allocation.Lockup.Add(allocation.ValidatorSelection...).Add(allocation.Holdings...)
	dust := total.Sub(sum[0])
	k.Logger(ctx).Info(
		"rewards distribution",
		"total", total,
		"validatorSelectionAllocation", allocation.ValidatorSelection,
		"holdingsAllocation", allocation.Holdings,
		"lockupAllocation", allocation.Lockup,
		"sum", sum,
		"dust", dust,
	)

	// Add dust to validator choice allocation (favors decentralization)
	k.Logger(ctx).Info("add dust to validatorSelectionAllocation...")
	allocation.ValidatorSelection = allocation.ValidatorSelection.Add(dust)

	return allocation
}

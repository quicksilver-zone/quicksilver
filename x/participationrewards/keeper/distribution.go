package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

type rewardsAllocation struct {
	ValidatorSelection sdk.Coins
	Holdings           sdk.Coins
	Lockup             sdk.Coins
}

type tokenValues struct {
	Tokens map[string]tokenValue
}

type tokenValue struct {
	Symbol     string
	Multiplier int64
	Value      sdk.Dec
}

// getRewardsAllocations returns an instance of rewardsAllocation with values
// set according to the module balance and set DistributionProportions
// parameters.
func (k Keeper) getRewardsAllocations(ctx sdk.Context) rewardsAllocation {
	var allocation rewardsAllocation

	denom := k.stakingKeeper.BondDenom(ctx)
	moduleAddress := k.accountKeeper.GetModuleAddress(types.ModuleName)
	moduleBalance := k.bankKeeper.GetBalance(ctx, moduleAddress, denom)

	k.Logger(ctx).Info("module account", "address", moduleAddress, "balance", moduleBalance)

	if moduleBalance.IsZero() {
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
			moduleBalance,
			params.DistributionProportions.ValidatorSelectionAllocation,
		),
	)
	allocation.Holdings = sdk.NewCoins(
		k.GetAllocation(
			ctx,
			moduleBalance,
			params.DistributionProportions.HoldingsAllocation,
		),
	)
	allocation.Lockup = sdk.NewCoins(
		k.GetAllocation(
			ctx,
			moduleBalance,
			params.DistributionProportions.LockupAllocation,
		),
	)

	// use sum to check total distribution to collect and allocate dust
	total := moduleBalance
	sum := allocation.Lockup.Add(allocation.ValidatorSelection...).Add(allocation.Holdings...)
	dust := total.SubAmount(sum.AmountOf(denom))
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

// allocateZoneRewards executes zone based rewards allocation. This entails
// rewards that are proportionally distributed to zones based on the tvl for
// each zone relative to the tvl of the QS protocol.
func (k Keeper) allocateZoneRewards(ctx sdk.Context, tvs tokenValues, allocation rewardsAllocation) error {
	k.Logger(ctx).Info("allocateZoneRewards", "token values", tvs, "allocation", allocation)

	if err := k.setZoneAllocations(ctx, tvs, allocation); err != nil {
		return err
	}

	k.allocateValidatorSelectionRewards(ctx)

	if err := k.allocateHoldingsRewards(ctx); err != nil {
		k.Logger(ctx).Error(err.Error())
		// TODO: remove once allocateHoldingsRewards is implemented: >>>
		if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, allocation.Holdings); err != nil {
			k.Logger(ctx).Error(err.Error())
		}
		// <<<
	}

	return nil
}

// setZoneAllocations returns the proportional zone rewards allocations as a
// map indexed by the zone id.
func (k Keeper) setZoneAllocations(ctx sdk.Context, tvs tokenValues, allocation rewardsAllocation) error {
	k.Logger(ctx).Info("setZoneAllocations", "allocation", allocation)

	zoneProps := make(map[string]sdk.Dec)

	otvl := sdk.NewDec(0)
	// pass 1: iterate zones - set tvl & calc overall tvl
	for _, zone := range k.icsKeeper.AllZones(ctx) {
		zone := zone
		tv, exists := tvs.Tokens[zone.BaseDenom]
		if !exists {
			err := fmt.Errorf("unable to obtain token value for zone %s", zone.ChainId)
			return err
		}
		ztvl := k.icsKeeper.GetDelegatedAmount(ctx, &zone).Amount.ToDec().
			Quo(sdk.NewDec(tv.Multiplier)).
			Mul(tv.Value)
		// set the zone tvl here, we will overwrite it with the correct
		// proportion once we have the overall protocol tvl;
		zoneProps[zone.ChainId] = ztvl
		k.Logger(ctx).Info("zone tvl", "zone", zone.ChainId, "tvl", ztvl)

		otvl = otvl.Add(ztvl)
	}

	// check overall protocol tvl
	if otvl.IsZero() {
		err := fmt.Errorf("protocol tvl is zero")
		return err
	}

	// pass 2: iterate zones - calc zone tvl proportion & set allocations
	for _, zone := range k.icsKeeper.AllZones(ctx) {
		// explicit memory referencing
		zone := zone

		ztvl, exists := zoneProps[zone.ChainId]
		if !exists {
			// this panic should never occur as the first iteration inserts a
			// a value into the map (see above).
			panic("unable to obtain zone proportion on second zone pass")
		}

		zp := ztvl.Quo(otvl)
		k.Logger(ctx).Info("zone proportion", "zone", zone.ChainId, "proportion", zp)

		zone.ValidatorSelectionAllocation = sdk.NewCoins(
			sdk.NewCoin(
				k.stakingKeeper.BondDenom(ctx),
				allocation.ValidatorSelection.AmountOfNoDenomValidation(k.stakingKeeper.BondDenom(ctx)).ToDec().
					Mul(zp).TruncateInt(),
			),
		)

		zone.HoldingsAllocation = sdk.NewCoins(
			sdk.NewCoin(
				k.stakingKeeper.BondDenom(ctx),
				allocation.Holdings.AmountOfNoDenomValidation(k.stakingKeeper.BondDenom(ctx)).ToDec().
					Mul(zp).TruncateInt(),
			),
		)

		k.icsKeeper.SetZone(ctx, &zone)
	}

	return nil
}

// distributeToUsers sends the allocated user rewards to the user address.
func (k Keeper) distributeToUsers(ctx sdk.Context, userAllocations []userAllocation) error {
	k.Logger(ctx).Info("distributeToUsers", "allocations", userAllocations)
	hasError := false

	for _, ua := range userAllocations {
		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(ua.Address), ua.Coins)
		if err != nil {
			k.Logger(ctx).Error("distribute to user", "address", ua.Address, "coins", ua.Coins)
			hasError = true
		}
	}

	if hasError {
		return fmt.Errorf("errors occurred while distributing rewards, review logs")
	}

	return nil
}

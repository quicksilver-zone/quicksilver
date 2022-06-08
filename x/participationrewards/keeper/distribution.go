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

// allocateZoneRewards allocates both validator selection rewards and holdings
// rewards across zones, based on the TVL for each zone in proportion to the
// overall protocol TVL (sum of zone TVLs).
/*func (k Keeper) allocateZoneRewards(ctx sdk.Context, allocation rewardsAllocation) error {
	var valuescb Callback = func(k Keeper, ctx sdk.Context, response []byte, query icqtypes.Query) error {
		tvResponse := QueryTokenValuesResponse{}
		err := k.cdc.Unmarshal(response, &tvResponse)
		if err != nil {
			return err
		}

		zoneProportions, err := k.getZoneProportions(ctx, tvs)
		if err != nil {
			return err
		}

		if err := k.allocateValidatorSelectionRewards(ctx, allocation.ValidatorSelection, zoneProportions); err != nil {
			k.Logger(ctx).Error(err.Error())
		}

		if err := k.allocateHoldingsRewards(ctx, allocation.Holdings, zoneProportions); err != nil {
			k.Logger(ctx).Error(err.Error())
		}

		return nil
	}

	// obtain zones token values
	valuesQuery := QueryTokenValuesRequest{}
	bz := k.cdc.MustMarshal(&valuesQuery)

	// Request to obtain a comparable value for tokens across all zones
	k.icqKeeper.MakeRequest(
		ctx,
		ibc.ConnectionId,
		ibc.ChainId,
		"cosmos.distribution.v1beta1.Query/DelegationTotalRewards",
		bz,
		sdk.NewInt(-1),
		types.ModuleName,
		valuescb,
	)

	return fmt.Errorf("not implemented (stub)")
}*/

// TODO: remove when above is properly implemented
func (k Keeper) allocateZoneRewards(ctx sdk.Context, tvs tokenValues, allocation rewardsAllocation) error {
	k.Logger(ctx).Info("allocateZoneRewards", "token values", tvs, "allocation", allocation)

	zoneProportions, err := k.getZoneProportions(ctx, tvs)
	if err != nil {
		return err
	}

	if err := k.allocateValidatorSelectionRewards(ctx, allocation.ValidatorSelection, zoneProportions); err != nil {
		k.Logger(ctx).Error(err.Error())
	}

	if err := k.allocateHoldingsRewards(ctx, allocation.Holdings, zoneProportions); err != nil {
		k.Logger(ctx).Error(err.Error())
	}

	return nil
}

func (k Keeper) getZoneProportions(ctx sdk.Context, tvs tokenValues) (map[string]sdk.Dec, error) {
	k.Logger(ctx).Info("getZoneProportions", "token values", tvs)

	zoneProps := make(map[string]sdk.Dec)

	otvl := sdk.NewDec(0)
	for _, zone := range k.icsKeeper.AllRegisteredZones(ctx) {
		tv, exists := tvs.Tokens[zone.BaseDenom]
		if !exists {
			err := fmt.Errorf("unable to obtain token value for zone %s", zone.ChainId)
			return nil, err
		}

		ztvl := zone.GetDelegatedAmount().Amount.ToDec().
			Quo(sdk.NewDec(tv.Multiplier)).
			Mul(tv.Value)
		// set the zone tvl here, we will overwrite it with the correct
		// proportion once we have the overall tvl;
		zoneProps[zone.ChainId] = ztvl
		k.Logger(ctx).Info("zone tvl", "zone", zone.ChainId, "tvl", ztvl)

		otvl = otvl.Add(ztvl)
	}

	for zid, ztvl := range zoneProps {
		zoneProps[zid] = ztvl.Quo(otvl)
		k.Logger(ctx).Info("zone proportion", "zone", zid, "proportion", zoneProps[zid])
	}

	return zoneProps, nil
}

func (k Keeper) getZoneAllocations(ctx sdk.Context, zoneProps map[string]sdk.Dec, allocation sdk.Coins) map[string]sdk.Coins {
	k.Logger(ctx).Info("getZoneAllocations", "proportions", zoneProps, "allocation", allocation)

	zoneAllocations := make(map[string]sdk.Coins)

	for zid, zp := range zoneProps {
		zoneAllocations[zid] = sdk.NewCoins(
			sdk.NewCoin(
				allocation.GetDenomByIndex(0),
				allocation.AmountOf(allocation.GetDenomByIndex(0)).ToDec().
					Mul(zp).TruncateInt(),
			),
		)
	}

	return zoneAllocations
}

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
		return fmt.Errorf("errors occured while distributing rewards, review logs")
	}

	return nil
}

package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/internal/multierror"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

type RewardsAllocation struct {
	ValidatorSelection sdk.Int
	Holdings           sdk.Int
	Lockup             sdk.Int
}

type tokenValues map[string]sdk.Dec

// GetRewardsAllocations returns an instance of rewardsAllocation with values
// set according to the given moduleBalance and distribution proportions.
func GetRewardsAllocations(moduleBalance sdk.Int, proportions types.DistributionProportions) (*RewardsAllocation, error) {
	if moduleBalance.IsNil() || moduleBalance.IsZero() {
		return nil, types.ErrNothingToAllocate
	}

	if sum := proportions.Total(); !sum.Equal(sdk.OneDec()) {
		return nil, fmt.Errorf("%w: got %v", types.ErrInvalidTotalProportions, sum)
	}

	var allocation RewardsAllocation

	// split participation rewards allocations
	allocation.ValidatorSelection = sdk.NewDecFromInt(moduleBalance).Mul(proportions.ValidatorSelectionAllocation).TruncateInt()
	allocation.Holdings = sdk.NewDecFromInt(moduleBalance).Mul(proportions.HoldingsAllocation).TruncateInt()
	allocation.Lockup = sdk.NewDecFromInt(moduleBalance).Mul(proportions.LockupAllocation).TruncateInt()

	// use sum to check total distribution to collect and allocate dust
	sum := allocation.Lockup.Add(allocation.ValidatorSelection).Add(allocation.Holdings)
	dust := moduleBalance.Sub(sum)

	// Add dust to validator choice allocation (favors decentralization)
	allocation.ValidatorSelection = allocation.ValidatorSelection.Add(dust)

	return &allocation, nil
}

func (k Keeper) calcTokenValues(ctx sdk.Context) (tokenValues, error) {
	k.Logger(ctx).Info("calcTokenValues")

	tvs := make(map[string]sdk.Dec)

	// get base zone (Cosmos)
	var cosmosZone *icstypes.Zone
	k.icsKeeper.IterateZones(ctx, func(_ int64, zone icstypes.Zone) (stop bool) {
		if zone.AccountPrefix == "cosmos" {
			cosmosZone = &zone
			return true
		}
		return false
	})
	if cosmosZone == nil {
		return nil, fmt.Errorf("unable to find Cosmos zone")
	}

	// add base value
	tvs[cosmosZone.BaseDenom] = sdk.OneDec()

	// capture errors from iterator
	errors := make(map[string]error)
	k.IteratePrefixedProtocolDatas(ctx, "osmosis/pools", func(idx int64, data types.ProtocolData) bool {
		idxLabel := fmt.Sprintf("index[%d]", idx)
		ipool, err := UnmarshalProtocolData(types.ProtocolDataOsmosisPool, data.Data)
		if err != nil {
			errors[idxLabel] = err
			return true
		}
		pool, _ := ipool.(types.OsmosisPoolProtocolData)

		// pool must be a cosmos pair
		if len(pool.Zones) != 2 {
			// not a pair: skip
			return false
		}

		// values to be captured and used
		//  - baseIBCDenom -> the cosmos IBC denom in this pair
		//  - queryIBCDenom -> the target IBC denom in this pair
		//  - valueDenom -> the target zone.BaseDenom
		var baseIBCDenom, queryIBCDenom, valueDenom string
		isCosmosPair := false

		for chainID, denom := range pool.Zones {
			zone, ok := k.icsKeeper.GetZone(ctx, chainID)
			if !ok {
				errors[idxLabel] = fmt.Errorf("zone not found, %s", chainID)
				return true
			}

			if zone.AccountPrefix == "cosmos" {
				isCosmosPair = true
				baseIBCDenom = denom
				continue
			}

			queryIBCDenom = denom
			valueDenom = zone.BaseDenom
		}

		if isCosmosPair {
			value, err := pool.PoolData.SpotPrice(ctx, baseIBCDenom, queryIBCDenom)
			if err != nil {
				errors[idxLabel] = err
				return true
			}

			tvs[valueDenom] = value
		}

		return false
	})

	if len(errors) > 0 {
		return nil, multierror.New(errors)
	}

	return tvs, nil
}

// allocateZoneRewards executes zone based rewards allocation. This entails
// rewards that are proportionally distributed to zones based on the tvl for
// each zone relative to the tvl of the QS protocol.
func (k Keeper) allocateZoneRewards(ctx sdk.Context, tvs tokenValues, allocation RewardsAllocation) error {
	k.Logger(ctx).Info("allocateZoneRewards", "token values", tvs, "allocation", allocation)

	if err := k.setZoneAllocations(ctx, tvs, allocation); err != nil {
		return err
	}

	k.allocateValidatorSelectionRewards(ctx)

	if err := k.allocateHoldingsRewards(ctx); err != nil {
		return err
	}

	return nil
}

// setZoneAllocations returns the proportional zone rewards allocations as a
// map indexed by the zone id.
func (k Keeper) setZoneAllocations(ctx sdk.Context, tvs tokenValues, allocation RewardsAllocation) error {
	k.Logger(ctx).Info("setZoneAllocations", "allocation", allocation)

	otvl := sdk.NewDec(0)
	// pass 1: iterate zones - set tvl & calc overall tvl
	for _, zone := range k.icsKeeper.AllZones(ctx) {
		// explicit memory referencing
		zone := zone

		tv, exists := tvs[zone.BaseDenom]
		if !exists {
			err := fmt.Errorf("unable to obtain token value for zone %s", zone.ChainId)
			return err
		}
		ztvl := sdk.NewDecFromInt(k.icsKeeper.GetDelegatedAmount(ctx, &zone).Amount).Mul(tv)

		zone.Tvl = ztvl
		k.icsKeeper.SetZone(ctx, &zone)

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

		if zone.Tvl.IsNil() {
			zone.Tvl = sdk.ZeroDec()
		}

		zp := zone.Tvl.Quo(otvl)
		k.Logger(ctx).Info("zone proportion", "zone", zone.ChainId, "proportion", zp)

		zone.ValidatorSelectionAllocation = sdk.NewDecFromInt(allocation.ValidatorSelection).Mul(zp).TruncateInt().Uint64()
		zone.HoldingsAllocation = sdk.NewDecFromInt(allocation.Holdings).Mul(zp).TruncateInt().Uint64()

		k.icsKeeper.SetZone(ctx, &zone)
	}

	return nil
}

// distributeToUsers sends the allocated user rewards to the user address.
func (k Keeper) distributeToUsers(ctx sdk.Context, userAllocations []userAllocation) error {
	k.Logger(ctx).Info("distributeToUsers", "allocations", userAllocations)
	hasError := false

	for _, ua := range userAllocations {
		coins := sdk.NewCoins(
			sdk.NewCoin(
				k.stakingKeeper.BondDenom(ctx),
				ua.Amount,
			),
		)
		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(ua.Address), coins)
		if err != nil {
			k.Logger(ctx).Error("distribute to user", "address", ua.Address, "coins", coins)
			hasError = true
		}
	}

	if hasError {
		return fmt.Errorf("errors occurred while distributing rewards, review logs")
	}

	return nil
}

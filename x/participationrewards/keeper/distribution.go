package keeper

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/internal/multierror"
	"github.com/ingenuity-build/quicksilver/utils"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

type tokenValues map[string]sdk.Dec

func (k *Keeper) calcTokenValues(ctx sdk.Context) (tokenValues, error) {
	k.Logger(ctx).Info("calcTokenValues")

	data, found := k.GetProtocolData(ctx, types.ProtocolDataTypeOsmosisParams, "osmosisparams")
	if !found {
		return tokenValues{}, errors.New("could not find osmosisparams protocol data")
	}
	osmoParams, err := types.UnmarshalProtocolData(types.ProtocolDataTypeOsmosisParams, data.Data)
	if err != nil {
		return tokenValues{}, err
	}

	baseDenom := osmoParams.(*types.OsmosisParamsProtocolData).BaseDenom
	baseChain := osmoParams.(*types.OsmosisParamsProtocolData).BaseChain

	tvs := make(map[string]sdk.Dec)

	// add base value
	tvs[baseDenom] = sdk.OneDec()

	// capture errors from iteratora
	errs := make(map[string]error)
	k.IteratePrefixedProtocolDatas(ctx, types.GetPrefixProtocolDataKey(types.ProtocolDataTypeOsmosisPool), func(idx int64, _ []byte, data types.ProtocolData) bool {
		idxLabel := fmt.Sprintf("index[%d]", idx)
		ipool, err := types.UnmarshalProtocolData(types.ProtocolDataTypeOsmosisPool, data.Data)
		if err != nil {
			errs[idxLabel] = err
			return true
		}
		pool, _ := ipool.(*types.OsmosisPoolProtocolData)

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
		isBasePair := false

		for chainID, denom := range pool.Zones {
			if chainID == baseChain {
				isBasePair = true
				baseIBCDenom = denom
				continue
			}

			zone, ok := k.icsKeeper.GetZone(ctx, chainID)
			if !ok {
				errs[idxLabel] = fmt.Errorf("zone not found, %s", chainID)
				return true
			}

			queryIBCDenom = denom
			valueDenom = zone.BaseDenom
		}

		if isBasePair {
			if pool.PoolData == nil {
				errs[idxLabel] = fmt.Errorf("pool data is nil, awaiting OsmosisPoolUpdateCallback")
				return true
			}
			pool, err := pool.GetPool()
			if err != nil {
				errs[idxLabel] = err
				return true
			}
			value, err := pool.SpotPrice(ctx, baseIBCDenom, queryIBCDenom)
			if err != nil {
				errs[idxLabel] = err
				return true
			}

			tvs[valueDenom] = value
		}

		return false
	})

	if len(errs) > 0 {
		return nil, multierror.New(errs)
	}

	return tvs, nil
}

// AllocateZoneRewards executes zone based rewards allocation. This entails
// rewards that are proportionally distributed to zones based on the tvl for
// each zone relative to the tvl of the QS protocol.
func (k *Keeper) AllocateZoneRewards(ctx sdk.Context, tvs tokenValues, allocation types.RewardsAllocation) error {
	k.Logger(ctx).Info("allocateZoneRewards", "token values", tvs, "allocation", allocation)

	if err := k.SetZoneAllocations(ctx, tvs, allocation); err != nil {
		return err
	}

	k.AllocateValidatorSelectionRewards(ctx)

	return k.AllocateHoldingsRewards(ctx)
}

// SetZoneAllocations returns the proportional zone rewards allocations as a
// map indexed by the zone id.
func (k *Keeper) SetZoneAllocations(ctx sdk.Context, tvs tokenValues, allocation types.RewardsAllocation) error {
	k.Logger(ctx).Info("setZoneAllocations", "allocation", allocation)

	otvl := sdk.ZeroDec()
	// pass 1: iterate zones - set tvl & calc overall tvl
	for _, zone := range k.icsKeeper.AllZones(ctx) {
		// explicit memory referencing
		zone := zone

		tv, exists := tvs[zone.BaseDenom]
		if !exists {
			k.Logger(ctx).Error(fmt.Sprintf("unable to obtain token value for zone %s", zone.ChainId))
			continue
		}
		ztvl := sdk.NewDecFromInt(k.icsKeeper.GetDelegatedAmount(ctx, &zone).Amount.Add(k.icsKeeper.GetDelegationsInProcess(ctx, &zone))).Mul(tv)

		zone.Tvl = ztvl
		k.icsKeeper.SetZone(ctx, &zone)

		k.Logger(ctx).Info("zone tvl", "zone", zone.ChainId, "tvl", ztvl)

		otvl = otvl.Add(ztvl)
	}

	// check overall protocol tvl
	if otvl.IsZero() {
		err := errors.New("protocol tvl is zero")
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

// DistributeToUsers sends the allocated user rewards to the user address.
func (k *Keeper) DistributeToUsers(ctx sdk.Context, userAllocations []types.UserAllocation) error {
	k.Logger(ctx).Info("distributeToUsers", "allocations", userAllocations)
	hasError := false

	for _, ua := range userAllocations {
		if ua.Amount.IsZero() {
			continue
		}

		coins := sdk.NewCoins(
			sdk.NewCoin(
				k.stakingKeeper.BondDenom(ctx),
				ua.Amount,
			),
		)

		addrBytes, err := utils.AccAddressFromBech32(ua.Address, "")
		if err != nil {
			k.Logger(ctx).Error("unmarshalling address", "address", ua.Address)
			hasError = true
		}

		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addrBytes, coins)
		if err != nil {
			k.Logger(ctx).Error("distribute to user", "address", ua.Address, "coins", coins)
			hasError = true
		} else {
			k.Logger(ctx).Info("distribute to user", "address", ua.Address, "coins", coins, "remaining", k.GetModuleBalance(ctx))
		}
	}

	if hasError {
		return errors.New("errors occurred while distributing rewards, review logs")
	}

	return nil
}

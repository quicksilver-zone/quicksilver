package keeper

import (
	"errors"
	"fmt"

	"go.uber.org/multierr"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/gamm/pool-models/stableswap"
	poolmanager "github.com/quicksilver-zone/quicksilver/third-party-chains/osmosis-types/poolmanager/types"
	"github.com/quicksilver-zone/quicksilver/utils"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	"github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

type TokenValues map[string]sdk.Dec

type (
	AssetGraph      map[string]map[string]sdk.Dec
	AssetGraphSlice map[string]map[string][]sdk.Dec
)

func DepthFirstSearch(graph AssetGraph, visited map[string]struct{}, asset string, price sdk.Dec, result TokenValues) {
	visited[asset] = struct{}{}
	result[asset] = price

	for _, neighbour := range utils.Keys(graph[asset]) {
		if _, ok := visited[neighbour]; !ok {
			DepthFirstSearch(graph, visited, neighbour, graph[asset][neighbour].Mul(price), result)
		}
	}
}

func (k *Keeper) CalcTokenValues(ctx sdk.Context) (TokenValues, error) {
	k.Logger(ctx).Info("calcTokenValues")

	_, osmoParams, err := GetAndUnmarshalProtocolData[*types.OsmosisParamsProtocolData](ctx, k, "osmosisparams", types.ProtocolDataTypeOsmosisParams)
	if err != nil {
		return TokenValues{}, err
	}

	baseDenom := osmoParams.BaseDenom

	tvs := make(TokenValues)
	graph := make(AssetGraphSlice)
	graph2 := make(AssetGraph)

	// add base value
	tvs[baseDenom] = sdk.OneDec()

	// capture errors from iterator
	errs := make(map[string]error)
	k.IteratePrefixedProtocolDatas(ctx, types.GetPrefixProtocolDataKey(types.ProtocolDataTypeOsmosisPool), func(idx int64, _ []byte, data types.ProtocolData) bool {
		idxLabel := fmt.Sprintf("index[%d]", idx)
		ipool, err := types.UnmarshalProtocolData(types.ProtocolDataTypeOsmosisPool, data.Data)
		if err != nil {
			errs[idxLabel] = err
			return true
		}
		pool, _ := ipool.(*types.OsmosisPoolProtocolData)

		// pool must be a base pair
		if len(pool.Denoms) != 2 {
			return false
		}

		if pool.PoolData == nil {
			errs[idxLabel] = errors.New("pool data is nil, awaiting OsmosisPoolUpdateCallback")
			return false
		}
		gammPool, err := pool.GetPool()
		if err != nil {
			errs[idxLabel] = err
			return false
		}

		denoms := utils.Keys(pool.Denoms)

		prettyDenom0 := pool.Denoms[denoms[0]].Denom
		prettyDenom1 := pool.Denoms[denoms[1]].Denom

		for _, ibcDenom := range denoms {
			if _, ok := graph[pool.Denoms[ibcDenom].Denom]; !ok {
				graph[pool.Denoms[ibcDenom].Denom] = make(map[string][]sdk.Dec)
			}
		}

		if gammPool.GetType() == poolmanager.Stableswap {
			// be defensive. if scaling_factors are missing, avoid panic.
			ss, ok := gammPool.(*stableswap.Pool)
			if !ok {
				errs[idxLabel] = fmt.Errorf("gammPool %d cannot be cast to StableswapPool", pool.PoolID)
				return false
			}
			if len(ss.GetScalingFactors()) != 2 {
				errs[idxLabel] = fmt.Errorf("gammPool %d is missing scaling factors", pool.PoolID)
				return false
			}
		}
		value, err := gammPool.SpotPrice(ctx, denoms[0], denoms[1])
		if err != nil {
			errs[idxLabel] = err
			return false
		}

		decVal := sdk.NewDecFromBigIntWithPrec(value.Dec().BigInt(), 18)

		graph[prettyDenom0][prettyDenom1] = append(graph[prettyDenom0][prettyDenom1], decVal)
		graph[prettyDenom1][prettyDenom0] = append(graph[prettyDenom1][prettyDenom0], sdk.OneDec().Quo(decVal))

		return false
	})

	k.IteratePrefixedProtocolDatas(ctx, types.GetPrefixProtocolDataKey(types.ProtocolDataTypeOsmosisCLPool), func(idx int64, _ []byte, data types.ProtocolData) bool {
		idxLabel := fmt.Sprintf("index[%d]", idx)
		ipool, err := types.UnmarshalProtocolData(types.ProtocolDataTypeOsmosisCLPool, data.Data)
		if err != nil {
			errs[idxLabel] = err
			return false
		}
		pool, _ := ipool.(*types.OsmosisClPoolProtocolData)

		// pool must be a base pair
		if len(pool.Denoms) != 2 {
			return false
		}

		if pool.PoolData == nil {
			errs[idxLabel] = errors.New("pool data is nil, awaiting OsmosisClPoolUpdateCallback")
			return false
		}
		clPool, err := pool.GetPool()
		if err != nil {
			errs[idxLabel] = err
			return false
		}

		denoms := utils.Keys(pool.Denoms)
		prettyDenom0 := pool.Denoms[denoms[0]].Denom
		prettyDenom1 := pool.Denoms[denoms[1]].Denom

		for _, ibcDenom := range denoms {
			if _, ok := graph[pool.Denoms[ibcDenom].Denom]; !ok {
				graph[pool.Denoms[ibcDenom].Denom] = make(map[string][]sdk.Dec)
			}
		}

		value, err := clPool.SpotPrice(ctx, denoms[0], denoms[1])
		if err != nil {
			errs[idxLabel] = err
			return false
		}

		decVal := sdk.NewDecFromBigIntWithPrec(value.Dec().BigInt(), 18)

		graph[prettyDenom0][prettyDenom1] = append(graph[prettyDenom0][prettyDenom1], decVal)
		graph[prettyDenom1][prettyDenom0] = append(graph[prettyDenom1][prettyDenom0], sdk.OneDec().Quo(decVal))

		return false
	})

	for _, denom0 := range utils.Keys(graph) {
		graph2[denom0] = make(map[string]sdk.Dec)
		values := graph[denom0]
		for _, denom1 := range utils.Keys(values) {
			value := sdk.ZeroDec()
			count := math.ZeroInt()
			for _, asset := range values[denom1] {
				value = value.Add(asset)
				count = count.Add(math.OneInt())
			}
			graph2[denom0][denom1] = value.QuoInt(count)
		}
	}

	visited := make(map[string]struct{})
	DepthFirstSearch(graph2, visited, baseDenom, sdk.OneDec(), tvs)

	if len(errs) > 0 {
		return TokenValues{}, multierr.Combine(utils.ErrorMapToSlice(errs)...)
	}

	return tvs, nil
}

// AllocateZoneRewards executes zone based rewards allocation. This entails
// rewards that are proportionally distributed to zones based on the tvl for
// each zone relative to the tvl of the QS protocol.
func (k *Keeper) AllocateZoneRewards(ctx sdk.Context, tvs TokenValues, allocation types.RewardsAllocation) error {
	k.Logger(ctx).Info("allocateZoneRewards", "token values", tvs, "allocation", allocation)

	if err := k.SetZoneAllocations(ctx, tvs, allocation); err != nil {
		return err
	}

	k.AllocateValidatorSelectionRewards(ctx)

	return k.AllocateHoldingsRewards(ctx)
}

// SetZoneAllocations returns the proportional zone rewards allocations as a
// map indexed by the zone id.
func (k *Keeper) SetZoneAllocations(ctx sdk.Context, tvs TokenValues, allocation types.RewardsAllocation) error {
	k.Logger(ctx).Info("setZoneAllocations", "allocation", allocation)

	otvl := sdk.ZeroDec()
	// pass 1: iterate zones - set tvl & calc overall tvl
	k.icsKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
		tv, exists := tvs[zone.BaseDenom]
		if !exists {
			k.Logger(ctx).Error(fmt.Sprintf("unable to obtain token value for zone %s", zone.ChainId))
			return false
		}
		ztvl := sdk.NewDecFromInt(k.icsKeeper.GetDelegatedAmount(ctx, zone).Amount.Add(k.icsKeeper.GetDelegationsInProcess(ctx, zone.ChainId))).Mul(tv)
		zone.Tvl = ztvl
		k.icsKeeper.SetZone(ctx, zone)

		k.Logger(ctx).Info("zone tvl", "zone", zone.ChainId, "tvl", ztvl)

		otvl = otvl.Add(ztvl)
		return false
	})

	// check overall protocol tvl
	if otvl.IsZero() {
		err := errors.New("protocol tvl is zero")
		return err
	}

	// pass 2: iterate zones - calc zone tvl proportion & set allocations
	k.icsKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
		if zone.Tvl.IsNil() {
			zone.Tvl = sdk.ZeroDec()
		}

		zp := zone.Tvl.Quo(otvl)
		k.Logger(ctx).Info("zone proportion", "zone", zone.ChainId, "proportion", zp)

		zone.ValidatorSelectionAllocation = sdk.NewDecFromInt(allocation.ValidatorSelection).Mul(zp).TruncateInt().Uint64()
		zone.HoldingsAllocation = sdk.NewDecFromInt(allocation.Holdings).Mul(zp).TruncateInt().Uint64()
		k.icsKeeper.SetZone(ctx, zone)
		return false
	})

	return nil
}

// DistributeToUsersFromModule sends the allocated user rewards to the user address.
func (k *Keeper) DistributeToUsersFromModule(ctx sdk.Context, userAllocations []types.UserAllocation) error {
	k.Logger(ctx).Info("distribute to users from module", "allocations", userAllocations)

	for _, ua := range userAllocations {
		if !ua.Amount.IsPositive() {
			continue
		}

		coins := sdk.NewCoins(ua.Amount)

		addrBytes, err := addressutils.AccAddressFromBech32(ua.Address, "")
		if err != nil {
			return err
		}

		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addrBytes, coins)
		if err != nil {
			return err
		}
		k.Logger(ctx).Info("distribute to user", "address", ua.Address, "coins", coins, "remaining", k.GetModuleBalance(ctx))

	}

	return nil
}

// DistributeToUsers sends the allocated user rewards to the user address.
func (k *Keeper) DistributeToUsersFromAddress(ctx sdk.Context, userAllocations []types.UserAllocation, fromAddress string) error {
	k.Logger(ctx).Info("distribute to users from account", "allocations", userAllocations)

	fromAddrBytes, err := addressutils.AccAddressFromBech32(fromAddress, "")
	if err != nil {
		return err
	}

	for _, ua := range userAllocations {
		if !ua.Amount.IsPositive() {
			continue
		}

		coins := sdk.NewCoins(
			ua.Amount,
		)

		addrBytes, err := addressutils.AccAddressFromBech32(ua.Address, "")
		if err != nil {
			return err
		}

		err = k.bankKeeper.SendCoins(ctx, fromAddrBytes, addrBytes, coins)
		if err != nil {
			return err
		}
		k.Logger(ctx).Info("distribute to user", "address", ua.Address, "coins", coins, "remaining", k.GetModuleBalance(ctx))
	}

	return nil
}

package keeper

import (
	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"

	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

// GetDelegation returns a specific delegation.
func (k *Keeper) GetDelegation(ctx sdk.Context, chainID string, delegatorAddress, validatorAddress string) (delegation types.Delegation, found bool) {
	store := ctx.KVStore(k.storeKey)

	_, delAddr, _ := bech32.DecodeAndConvert(delegatorAddress)
	_, valAddr, _ := bech32.DecodeAndConvert(validatorAddress)

	key := types.GetDelegationKey(chainID, delAddr, valAddr)

	value := store.Get(key)
	if value == nil {
		return delegation, false
	}

	delegation = types.MustUnmarshalDelegation(k.cdc, value)

	return delegation, true
}

// GetPerformanceDelegation returns a specific delegation.
func (k *Keeper) GetPerformanceDelegation(ctx sdk.Context, chainID string, performanceAddress *types.ICAAccount, validatorAddress string) (delegation types.Delegation, found bool) {
	if performanceAddress == nil {
		return types.Delegation{}, false
	}

	store := ctx.KVStore(k.storeKey)

	_, delAddr, _ := bech32.DecodeAndConvert(performanceAddress.Address)
	_, valAddr, _ := bech32.DecodeAndConvert(validatorAddress)

	key := types.GetPerformanceDelegationKey(chainID, delAddr, valAddr)

	value := store.Get(key)
	if value == nil {
		return delegation, false
	}

	delegation = types.MustUnmarshalDelegation(k.cdc, value)

	return delegation, true
}

// IterateAllDelegations iterates through all of the delegations.
func (k *Keeper) IterateAllDelegations(ctx sdk.Context, chainID string, cb func(delegation types.Delegation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, append(types.KeyPrefixDelegation, chainID...))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		if cb(delegation) {
			break
		}
	}
}

// GetAllDelegations returns all delegations used during genesis dump.
func (k *Keeper) GetAllDelegations(ctx sdk.Context, chainID string) (delegations []types.Delegation) {
	k.IterateAllDelegations(ctx, chainID, func(delegation types.Delegation) bool {
		delegations = append(delegations, delegation)
		return false
	})

	return delegations
}

// IterateAllPerformanceDelegations iterates through all of the delegations.
func (k *Keeper) IterateAllPerformanceDelegations(ctx sdk.Context, chainID string, cb func(delegation types.Delegation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, append(types.KeyPrefixPerformanceDelegation, chainID...))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		if cb(delegation) {
			break
		}
	}
}

// GetAllDelegations returns all delegations used during genesis dump.
func (k *Keeper) GetAllPerformanceDelegations(ctx sdk.Context, chainID string) (delegations []types.Delegation) {
	k.IterateAllPerformanceDelegations(ctx, chainID, func(delegation types.Delegation) bool {
		delegations = append(delegations, delegation)
		return false
	})

	return delegations
}

// GetAllDelegations returns all delegations used during genesis dump.
func (k *Keeper) GetAllDelegationsAsPointer(ctx sdk.Context, chainID string) (delegations []*types.Delegation) {
	k.IterateAllDelegations(ctx, chainID, func(delegation types.Delegation) bool {
		delegations = append(delegations, &delegation)
		return false
	})

	return delegations
}

// GetAllDelegations returns all delegations used during genesis dump.
func (k *Keeper) GetAllPerformanceDelegationsAsPointer(ctx sdk.Context, chainID string) (delegations []*types.Delegation) {
	k.IterateAllPerformanceDelegations(ctx, chainID, func(delegation types.Delegation) bool {
		delegations = append(delegations, &delegation)
		return false
	})

	return delegations
}

// GetDelegatorDelegations returns a given amount of all the delegations from a
// delegator.
func (k *Keeper) GetDelegatorDelegations(ctx sdk.Context, chainID string, delegator sdk.AccAddress) (delegations []types.Delegation) {
	k.IterateDelegatorDelegations(ctx, chainID, delegator, func(delegation types.Delegation) bool {
		delegations = append(delegations, delegation)
		return false
	})

	return delegations
}

// SetDelegation sets a delegation.
func (k *Keeper) SetDelegation(ctx sdk.Context, chainID string, delegation types.Delegation) {
	delegatorAddress := delegation.GetDelegatorAddr()

	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalDelegation(k.cdc, delegation)
	store.Set(types.GetDelegationKey(chainID, delegatorAddress, delegation.GetValidatorAddr()), b)
}

// SetPerformanceDelegation sets a delegation.
func (k *Keeper) SetPerformanceDelegation(ctx sdk.Context, chainID string, delegation types.Delegation) {
	delegatorAddress := delegation.GetDelegatorAddr()

	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalDelegation(k.cdc, delegation)
	store.Set(types.GetPerformanceDelegationKey(chainID, delegatorAddress, delegation.GetValidatorAddr()), b)
}

// RemoveDelegation removes a delegation.
func (k *Keeper) RemoveDelegation(ctx sdk.Context, chainID string, delegation types.Delegation) error {
	delegatorAddress := delegation.GetDelegatorAddr()

	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetDelegationKey(chainID, delegatorAddress, delegation.GetValidatorAddr()))
	return nil
}

// RemovePerformanceDelegation removes a performance delegation.
func (k *Keeper) RemovePerformanceDelegation(ctx sdk.Context, chainID string, delegation types.Delegation) error {
	delegatorAddress := delegation.GetDelegatorAddr()

	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetPerformanceDelegationKey(chainID, delegatorAddress, delegation.GetValidatorAddr()))
	return nil
}

// IterateDelegatorDelegations iterates through one delegator's delegations.
func (k *Keeper) IterateDelegatorDelegations(ctx sdk.Context, chainID string, delegator sdk.AccAddress, cb func(delegation types.Delegation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := types.GetDelegationsKey(chainID, delegator)
	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		if cb(delegation) {
			break
		}
	}
}


func (k *Keeper) GetDelegationMap(ctx sdk.Context, chainID string) (out map[string]sdkmath.Int, sum sdkmath.Int, locked map[string]bool, lockedSum sdkmath.Int) {
	out = make(map[string]sdkmath.Int)
	locked = make(map[string]bool)
	sum = sdk.ZeroInt()
	lockedSum = sdk.ZeroInt()

	k.IterateAllDelegations(ctx, chainID, func(delegation types.Delegation) bool {
		out[delegation.ValidatorAddress] = delegation.Amount.Amount
		if delegation.RedelegationEnd >= ctx.BlockTime().Unix() {
			locked[delegation.ValidatorAddress] = true
			lockedSum = lockedSum.Add(delegation.Amount.Amount)
		}
		sum = sum.Add(delegation.Amount.Amount)
		return false
	})

	return out, sum, locked, lockedSum
}

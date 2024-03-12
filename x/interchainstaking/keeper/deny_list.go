package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

// SetZoneValidatorToDenyList sets the zone validator deny list
func (k *Keeper) SetZoneValidatorToDenyList(ctx sdk.Context, chainID string, validator types.Validator) error {
	store := ctx.KVStore(k.storeKey)

	key := types.GetZoneDeniedValidatorKey(chainID)
	b := types.MustMarshalValidator(k.cdc, validator)
	store.Set(key, b)
	return nil
}

// GetZoneValidatorDenyList get the validator deny list of a specific zone
func (k *Keeper) GetZoneValidatorDenyList(ctx sdk.Context, chainID string) (types.ValidatorDenyList, bool) {
	store := ctx.KVStore(k.storeKey)

	key := []byte(chainID)
	value := store.Get(key)
	if value == nil {
		return types.ValidatorDenyList{}, false
	}

	denyList := types.MustUnmarshalDenyList(k.cdc, value)

	return denyList, true
}

// // IterateDelegatorDelegations iterates through one delegator's delegations.
// func (k *Keeper) IterateDelegatorDelegations(ctx sdk.Context, chainID string, delegator sdk.AccAddress, cb func(delegation types.Delegation) (stop bool)) {
// 	store := ctx.KVStore(k.storeKey)
// 	delegatorPrefixKey := types.GetDelegationsKey(chainID, delegator)
// 	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey)
// 	defer iterator.Close()

//		for ; iterator.Valid(); iterator.Next() {
//			delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
//			if cb(delegation) {
//				break
//			}
//		}
//	}
func (k *Keeper) IterateZoneDeniedValidator(ctx sdk.Context, chainID string, cb func(validator types.Validator) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	deniedValPrefixKey := types.GetZoneDeniedValidatorKey(chainID)
	iterator := sdk.KVStorePrefixIterator(store, deniedValPrefixKey)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		validator := types.MustUnmarshalValidator(k.cdc, iterator.Value())
		if cb(validator) {
			break
		}

	}
}

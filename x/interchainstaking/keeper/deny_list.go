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
func (k *Keeper) GetZoneValidatorDenyList(ctx sdk.Context, chainID string) ([]types.Validator, bool) {
	store := ctx.KVStore(k.storeKey)

	key := []byte(chainID)
	value := store.Get(key)
	if value == nil {
		return []types.Validator{}, false
	}

	denyList := []types.Validator{}
	k.IterateZoneDeniedValidator(ctx, chainID, func(validator types.Validator) bool {
		denyList = append(denyList, validator)
		return false
	})

	return denyList, true
}

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

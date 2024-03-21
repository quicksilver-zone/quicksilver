package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

// SetZoneValidatorToDenyList sets the zone validator deny list
func (k *Keeper) SetZoneValidatorToDenyList(ctx sdk.Context, chainID string, validatorAddress sdk.ValAddress) error {
	store := ctx.KVStore(k.storeKey)

	key := types.GetDeniedValidatorKey(chainID, validatorAddress)
	store.Set(key, validatorAddress)
	return nil
}

// GetZoneValidatorDenyList get the validator deny list of a specific zone
func (k *Keeper) GetZoneValidatorDenyList(ctx sdk.Context, chainID string) (denyList []string) {
	zone, found := k.GetZone(ctx, chainID)
	if !found {
		return denyList
	}
	k.IterateZoneDeniedValidator(ctx, chainID, func(validator sdk.ValAddress) bool {
		denyList = append(denyList, addressutils.MustEncodeAddressToBech32(zone.GetValoperPrefix(), validator))
		return false
	})

	return denyList
}

func (k *Keeper) GetDeniedValidatorInDenyList(ctx sdk.Context, chainID string, validatorAddress sdk.ValAddress) bool {
	key := types.GetDeniedValidatorKey(chainID, validatorAddress)
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(key)
	return bz != nil
}

// RemoveValidatorFromDenyList removes a validator from the deny list. Panic if the validator is not in the deny list
func (k *Keeper) RemoveValidatorFromDenyList(ctx sdk.Context, chainID string, validator sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetDeniedValidatorKey(chainID, validator)
	store.Delete(key)
}

func (k *Keeper) IterateZoneDeniedValidator(ctx sdk.Context, chainID string, cb func(validator sdk.ValAddress) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	deniedValPrefixKey := types.GetZoneDeniedValidatorKey(chainID)

	iterator := sdk.KVStorePrefixIterator(store, deniedValPrefixKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		validator := sdk.ValAddress(iterator.Value())
		if cb(validator) {
			break
		}
	}
}

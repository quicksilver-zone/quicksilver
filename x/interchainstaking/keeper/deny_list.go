package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

// SetZoneValidatorToDenyList sets the zone validator deny list
func (k *Keeper) SetZoneValidatorToDenyList(ctx sdk.Context, chainID string, validator types.Validator) error {
	store := ctx.KVStore(k.storeKey)

	key := types.GetDeniedValidatorKey(chainID, validator.ValoperAddress)
	addrBytes, err := sdk.ValAddressFromBech32(validator.ValoperAddress)
	if err != nil {
		return err
	}
	store.Set(key, addrBytes.Bytes())
	return nil
}

// GetZoneValidatorDenyList get the validator deny list of a specific zone
func (k *Keeper) GetZoneValidatorDenyList(ctx sdk.Context, chainID string) (denyList []string) {
	k.IterateZoneDeniedValidator(ctx, chainID, func(validator string) bool {
		denyList = append(denyList, validator)
		return false
	})

	return denyList
}

func (k *Keeper) GetDeniedValidatorInDenyList(ctx sdk.Context, chainID string, validatorAddress string) (types.Validator, bool) {
	key := types.GetDeniedValidatorKey(chainID, validatorAddress)
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(key)
	if bz == nil {
		return types.Validator{}, false
	}

	val, found := k.GetValidator(ctx, chainID, bz)
	if !found {
		return types.Validator{}, false
	}

	return val, true
}

// RemoveValidatorFromDenyList removes a validator from the deny list. Panic if the validator is not in the deny list
func (k *Keeper) RemoveValidatorFromDenyList(ctx sdk.Context, chainID string, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetDeniedValidatorKey(chainID, validator.ValoperAddress)
	store.Delete(key)
}

func (k *Keeper) IterateZoneDeniedValidator(ctx sdk.Context, chainID string, cb func(validator string) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	deniedValPrefixKey := types.GetZoneDeniedValidatorKey(chainID)

	iterator := sdk.KVStorePrefixIterator(store, deniedValPrefixKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		validator := sdk.ValAddress(iterator.Value()).String()
		if cb(validator) {
			break
		}
	}
}

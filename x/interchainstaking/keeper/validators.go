package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// GetValidators returns validators info by chainID
func (k Keeper) GetValidators(ctx sdk.Context, chainID string) []types.Validator {
	validators := []types.Validator{}
	k.IterateValidators(ctx, chainID, func(_ int64, validator types.Validator) (stop bool) {
		validators = append(validators, validator)
		return false
	})
	return validators
}

// GetValidators returns validators info by chainID
func (k Keeper) GetValidator(ctx sdk.Context, chainID string, address []byte) (types.Validator, bool) {
	val := types.Validator{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), GetZoneValidatorsKey(chainID))
	bz := store.Get(address)
	if len(bz) == 0 {
		return val, false
	}

	k.cdc.MustUnmarshal(bz, &val)
	return val, true
}

// SetValidators set validators info
func (k Keeper) SetValidator(ctx sdk.Context, chainID string, val types.Validator) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), GetZoneValidatorsKey(chainID))
	bz := k.cdc.MustMarshal(&val)
	valAddr, err := val.GetAddressBytes()
	if err != nil {
		return err
	}
	store.Set(valAddr, bz)
	return nil
}

// DeleteValidators delete validators info
func (k Keeper) DeleteValidator(ctx sdk.Context, chainID string, address []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), GetZoneValidatorsKey(chainID))
	store.Delete(address)
}

// IterateZones iterate through zones
func (k Keeper) IterateValidators(ctx sdk.Context, chainID string, fn func(index int64, validator types.Validator) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), nil)

	iterator := sdk.KVStorePrefixIterator(store, GetZoneValidatorsKey(chainID))
	defer iterator.Close()

	i := int64(0)

	for ; iterator.Valid(); iterator.Next() {
		validator := types.Validator{}
		k.cdc.MustUnmarshal(iterator.Value(), &validator)

		stop := fn(i, validator)

		if stop {
			break
		}
		i++
	}
}

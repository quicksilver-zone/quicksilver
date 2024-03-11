package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

// SetZoneValidatorToDenyList sets the zone validator deny list
func (k *Keeper) SetZoneValidatorToDenyList(ctx sdk.Context, chainID string, validator types.Validator) error {
	store := ctx.KVStore(k.storeKey)
	key := []byte(chainID)
	denyList, found := k.GetZoneValidatorDenyList(ctx, chainID)
	if !found {
		denyList = types.NewValidatorDenyListForZone(chainID)
	}
	// Append if not already in the list
	for _, v := range denyList.DeniedVals {
		if v.ValoperAddress == validator.ValoperAddress {
			return types.ErrValidatorAlreadyInDenyList
		}
	}
	denyList.DeniedVals = append(denyList.DeniedVals, validator)
	store.Set(key, types.MustMarshalDenyList(k.cdc, denyList))
	return nil
}

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

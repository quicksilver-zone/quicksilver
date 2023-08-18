package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// GetValidators returns validators by zone.
func (k *Keeper) GetValidators(ctx sdk.Context, zone *types.Zone) []types.Validator {
	var validators []types.Validator
	k.IterateValidators(ctx, zone, func(_ int64, validator types.Validator) (stop bool) {
		validators = append(validators, validator)
		return false
	})
	return validators
}

// GetValidatorAddresses returns a slice of validator addresses by zone.
func (k *Keeper) GetValidatorAddresses(ctx sdk.Context, zone *types.Zone) []string {
	var validators []string
	k.IterateValidators(ctx, zone, func(_ int64, validator types.Validator) (stop bool) {
		validators = append(validators, validator.ValoperAddress)
		return false
	})
	return validators
}

// GetActiveValidators returns validators by zone where status = BONDED.
func (k *Keeper) GetActiveValidators(ctx sdk.Context, zone *types.Zone) []types.Validator {
	var validators []types.Validator
	k.IterateValidators(ctx, zone, func(_ int64, validator types.Validator) (stop bool) {
		if validator.Status == stakingtypes.BondStatusBonded {
			validators = append(validators, validator)
		}
		return false
	})
	return validators
}

// GetValidator returns validator by zone.
func (k *Keeper) GetValidator(ctx sdk.Context, zone *types.Zone, address []byte) (types.Validator, bool) {
	val := types.Validator{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetZoneValidatorsKey(zone.BaseChainID()))
	bz := store.Get(address)
	if len(bz) == 0 {
		return val, false
	}

	k.cdc.MustUnmarshal(bz, &val)
	return val, true
}

// SetValidator sets a validator.
func (k *Keeper) SetValidator(ctx sdk.Context, zone *types.Zone, val types.Validator) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetZoneValidatorsKey(zone.BaseChainID()))
	bz := k.cdc.MustMarshal(&val)
	valAddr, err := val.GetAddressBytes()
	if err != nil {
		return err
	}
	store.Set(valAddr, bz)
	return nil
}

// DeleteValidator deletes a validator.
func (k *Keeper) DeleteValidator(ctx sdk.Context, zone *types.Zone, address []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetZoneValidatorsKey(zone.BaseChainID()))
	store.Delete(address)
}

// IterateValidators iterates through validators.
func (k *Keeper) IterateValidators(ctx sdk.Context, zone *types.Zone, fn func(index int64, validator types.Validator) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), nil)

	iterator := sdk.KVStorePrefixIterator(store, types.GetZoneValidatorsKey(zone.BaseChainID()))
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

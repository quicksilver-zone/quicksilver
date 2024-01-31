package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/quicksilver-zone/quicksilver/v7/x/interchainstaking/types"
)

// GetValidators returns validators by chainID.
func (k Keeper) GetValidators(ctx sdk.Context, chainID string) []types.Validator {
	validators := []types.Validator{}
	k.IterateValidators(ctx, chainID, func(_ int64, validator types.Validator) (stop bool) {
		validators = append(validators, validator)
		return false
	})
	return validators
}

// GetValidatorAddresses returns a slice of validator addresses by chainID.
func (k Keeper) GetValidatorAddresses(ctx sdk.Context, chainID string) []string {
	validators := []string{}
	k.IterateValidators(ctx, chainID, func(_ int64, validator types.Validator) (stop bool) {
		validators = append(validators, validator.ValoperAddress)
		return false
	})
	return validators
}

// GetActiveValidators returns validators by chainID where status = BONDED.
func (k Keeper) GetActiveValidators(ctx sdk.Context, chainID string) []types.Validator {
	validators := []types.Validator{}
	k.IterateValidators(ctx, chainID, func(_ int64, validator types.Validator) (stop bool) {
		if validator.Status == stakingtypes.BondStatusBonded {
			validators = append(validators, validator)
		}
		return false
	})
	return validators
}

// GetValidator returns validator by chainID and address.
func (k Keeper) GetValidator(ctx sdk.Context, chainID string, address []byte) (types.Validator, bool) {
	val := types.Validator{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetZoneValidatorsKey(chainID))
	bz := store.Get(address)
	if len(bz) == 0 {
		return val, false
	}

	k.cdc.MustUnmarshal(bz, &val)
	return val, true
}

// SetValidators set validators.
func (k Keeper) SetValidator(ctx sdk.Context, chainID string, val types.Validator) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetZoneValidatorsKey(chainID))
	bz := k.cdc.MustMarshal(&val)
	valAddr, err := val.GetAddressBytes()
	if err != nil {
		return err
	}
	store.Set(valAddr, bz)
	return nil
}

// DeleteValidator delete validator by chainID and address.
func (k Keeper) DeleteValidator(ctx sdk.Context, chainID string, address []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetZoneValidatorsKey(chainID))
	store.Delete(address)
}

// IterateZones iterates through zones.
func (k Keeper) IterateValidators(ctx sdk.Context, chainID string, fn func(index int64, validator types.Validator) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), nil)

	iterator := storetypes.KVStorePrefixIterator(store, types.GetZoneValidatorsKey(chainID))
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

// GetValidatorAddrByConsAddr returns validator address by Consensus address.
func (k Keeper) GetValidatorAddrByConsAddr(ctx sdk.Context, chainID string, consAddr []byte) (string, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetZoneValidatorAddrsByConsAddrKey(chainID))
	bz := store.Get(consAddr)
	if len(bz) == 0 {
		return "", false
	}

	return string(bz), true
}

// SetValidatorAddrByConsAddr set validator address by Consensus address.
func (k Keeper) SetValidatorAddrByConsAddr(ctx sdk.Context, chainID string, valAddr string, consAddr sdk.ConsAddress) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetZoneValidatorAddrsByConsAddrKey(chainID))
	store.Set(consAddr, []byte(valAddr))
}

// DeleteValidatorAddrByConsAddr delete validator address by Consensus address.
func (k Keeper) DeleteValidatorAddrByConsAddr(ctx sdk.Context, chainID string, consAddr []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetZoneValidatorAddrsByConsAddrKey(chainID))
	store.Delete(consAddr)
}

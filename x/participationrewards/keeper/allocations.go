package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

// GetHoldingAllocation returns sdk.Coin allocated to the given identifier.
func (k *Keeper) GetHoldingAllocation(ctx sdk.Context, chainID string) sdk.Coin {
	value := sdk.Coin{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixHoldingAllocation)
	bz := store.Get([]byte(chainID))
	if len(bz) == 0 {
		return value
	}

	k.cdc.MustUnmarshal(bz, &value)
	return value
}

// SetHoldingAllocation sets sdk.Coin allocated as the given identifier.
func (k Keeper) SetHoldingAllocation(ctx sdk.Context, chainID string, value sdk.Coin) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixHoldingAllocation)
	bz := k.cdc.MustMarshal(&value)
	store.Set([]byte(chainID), bz)
}

// GetValidatorAllocation returns sdk.Coin allocated to the given identifier.
func (k *Keeper) GetValidatorAllocation(ctx sdk.Context, chainID string) sdk.Coin {
	value := sdk.Coin{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixValidatorAllocation)
	bz := store.Get([]byte(chainID))
	if len(bz) == 0 {
		return value
	}

	k.cdc.MustUnmarshal(bz, &value)
	return value
}

// SetValidatorAllocation sets sdk.Coin allocated as the given identifier.
func (k Keeper) SetValidatorAllocation(ctx sdk.Context, chainID string, value sdk.Coin) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixValidatorAllocation)
	bz := k.cdc.MustMarshal(&value)
	store.Set([]byte(chainID), bz)
}

func (k Keeper) DetermineAllocations(ctx sdk.Context, moduleBalance sdk.Coin, proportions types.DistributionProportions) error {
	if moduleBalance.IsNil() || moduleBalance.IsZero() {
		return types.ErrNothingToAllocate
	}

	if sum := proportions.Total(); !sum.Equal(sdk.OneDec()) {
		return fmt.Errorf("%w: got %v", types.ErrInvalidTotalProportions, sum)
	}

	// split participation rewards allocations
	validatorAllocation := sdk.NewDecFromInt(moduleBalance.Amount).Mul(proportions.ValidatorSelectionAllocation).TruncateInt()
	holdingAllocation := sdk.NewDecFromInt(moduleBalance.Amount).Mul(proportions.HoldingsAllocation).TruncateInt()

	// use sum to check total distribution to collect and allocate dust
	sum := validatorAllocation.Add(holdingAllocation)
	dust := moduleBalance.Amount.Sub(sum)

	// Add dust to validator choice allocation (favors decentralization)
	validatorAllocation = validatorAllocation.Add(dust)

	k.SetHoldingAllocation(ctx, types.ModuleName, sdk.NewCoin(moduleBalance.Denom, holdingAllocation))
	k.SetValidatorAllocation(ctx, types.ModuleName, sdk.NewCoin(moduleBalance.Denom, validatorAllocation))

	return nil
}

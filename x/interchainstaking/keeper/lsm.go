package keeper

import (
	"errors"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// GetCap returns Cap info by zone and delegator
func (k Keeper) GetLsmCaps(ctx sdk.Context, zone *types.Zone) (types.LsmCaps, bool) {
	cap := types.LsmCaps{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixLsmCaps)
	bz := store.Get([]byte(zone.ChainId))
	if len(bz) == 0 {
		return cap, false
	}
	k.cdc.MustUnmarshal(bz, &cap)
	return cap, true
}

// SetCap store the delegator Cap
func (k Keeper) SetLsmCaps(ctx sdk.Context, zone *types.Zone, cap types.LsmCaps) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixLsmCaps)
	bz := k.cdc.MustMarshal(&cap)
	store.Set([]byte(zone.ChainId), bz)
}

// DeleteCap deletes delegator Cap
func (k Keeper) DeleteLsmCaps(ctx sdk.Context, zone *types.Zone, delegator string, snapshot bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixLsmCaps)
	store.Delete([]byte(delegator))
}

// IterateCaps iterate through Caps for a given zone
func (k Keeper) IterateLsmCaps(ctx sdk.Context, fn func(index int64, chainId string, cap types.LsmCaps) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixLsmCaps)

	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	i := int64(0)

	for ; iterator.Valid(); iterator.Next() {
		cap := types.LsmCaps{}
		k.cdc.MustUnmarshal(iterator.Value(), &cap)

		stop := fn(i, string(iterator.Key()), cap)

		if stop {
			break
		}
		i++
	}
}

// AllCaps returns every Cap in the store for the specified zone
func (k Keeper) AllLsmCaps(ctx sdk.Context, snapshot bool) map[string]types.LsmCaps {
	caps := map[string]types.LsmCaps{}
	k.IterateLsmCaps(ctx, func(_ int64, chainId string, cap types.LsmCaps) (stop bool) {
		caps[chainId] = cap
		return false
	})
	return caps
}

func (k Keeper) GetLiquidStakedSupply(zone *types.Zone) sdk.Dec {
	out := sdk.ZeroDec()
	for _, val := range zone.Validators {
		if val.Status == stakingtypes.BondStatusBonded {
			out = out.Add(val.LiquidShares)
		}
	}
	return out
}

func (k Keeper) GetTotalStakedSupply(zone *types.Zone) math.Int {
	out := sdk.ZeroInt()
	for _, val := range zone.Validators {
		if val.Status == stakingtypes.BondStatusBonded {
			out = out.Add(val.VotingPower)
		}
	}
	return out
}

func (k Keeper) CheckExceedsGlobalCap(ctx sdk.Context, zone *types.Zone, amount math.Int) bool {
	cap, found := k.GetLsmCaps(ctx, zone)
	if !found {
		// no caps found, permit
		return false
	}

	liquidSupply := k.GetLiquidStakedSupply(zone)
	totalSupply := sdk.NewDecFromInt(k.GetTotalStakedSupply(zone))
	amountDec := sdk.NewDecFromInt(amount)
	return liquidSupply.Add(amountDec).Quo(totalSupply).GT(cap.GlobalCap)
}

func (k Keeper) CheckExceedsValidatorCap(ctx sdk.Context, zone *types.Zone, validator string, amount math.Int) error {
	// Retrieve the cap for the given zone
	cap, found := k.GetLsmCaps(ctx, zone)
	if !found {
		// No cap found, permit the transaction
		return nil
	}

	// Retrieve the validator's information
	val, found := zone.GetValidatorByValoper(validator)
	if !found {
		// Validator not found, throw an error
		return errors.New("validator not found")
	}

	// Calculate the liquid shares and tokens
	amountDec := sdk.NewDecFromInt(amount)
	liquidShares := val.LiquidShares.Add(amountDec)
	tokens := sdk.NewDecFromInt(val.VotingPower).Add(amountDec)

	if liquidShares.Quo(tokens).GT(cap.ValidatorCap) {
		return errors.New("exceeds validator cap")
	}

	return nil
}

func (k Keeper) CheckExceedsValidatorBondCap(ctx sdk.Context, zone *types.Zone, validator string, amount math.Int) error {
	cap, found := k.GetLsmCaps(ctx, zone)
	if !found {
		// no caps found, permit
		return nil
	}

	val, found := zone.GetValidatorByValoper(validator)
	if !found {
		// cannot find validator, do not allow to proceed.
		return errors.New("validator not found")
	}

	var maxShares sdk.Dec
	if val != nil {
		maxShares = val.ValidatorBondShares.Mul(cap.ValidatorBondCap)
	} else {
		return errors.New("validator is nil")
	}

	amountDec := sdk.NewDecFromInt(amount)
	liquidShares := val.LiquidShares.Add(amountDec)

	if liquidShares.GT(maxShares) {
		return errors.New("exceeds validator bond cap")
	}

	return nil
}

package keeper

import (
	"errors"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

// GetCap returns Cap info by zone and delegator
func (k Keeper) GetLsmCaps(ctx sdk.Context, chainID string) (*types.LsmCaps, bool) {
	caps := types.LsmCaps{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixLsmCaps)
	bz := store.Get([]byte(chainID))
	if len(bz) == 0 {
		return nil, false
	}
	k.cdc.MustUnmarshal(bz, &caps)
	return &caps, true
}

// SetCap store the delegator Cap
func (k Keeper) SetLsmCaps(ctx sdk.Context, chainID string, caps types.LsmCaps) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixLsmCaps)
	bz := k.cdc.MustMarshal(&caps)
	store.Set([]byte(chainID), bz)
}

// DeleteCap deletes delegator Cap
func (k Keeper) DeleteLsmCaps(ctx sdk.Context, chainID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixLsmCaps)
	store.Delete([]byte(chainID))
}

// IterateCaps iterate through Caps for a given zone
func (k Keeper) IterateLsmCaps(ctx sdk.Context, fn func(index int64, chainID string, cap types.LsmCaps) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixLsmCaps)

	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	i := int64(0)

	for ; iterator.Valid(); iterator.Next() {
		caps := types.LsmCaps{}
		k.cdc.MustUnmarshal(iterator.Value(), &caps)

		stop := fn(i, string(iterator.Key()), caps)

		if stop {
			break
		}
		i++
	}
}

// AllCaps returns every Cap in the store for the specified zone
func (k Keeper) AllLsmCaps(ctx sdk.Context) map[string]types.LsmCaps {
	allCaps := map[string]types.LsmCaps{}
	k.IterateLsmCaps(ctx, func(_ int64, chainID string, caps types.LsmCaps) (stop bool) {
		allCaps[chainID] = caps
		return false
	})
	return allCaps
}

func (k Keeper) GetLiquidStakedSupply(ctx sdk.Context, zone *types.Zone) sdk.Dec {
	out := sdk.ZeroDec()
	for _, val := range k.GetActiveValidators(ctx, zone.ChainId) {
		if val.Status == stakingtypes.BondStatusBonded {
			out = out.Add(val.LiquidShares)
		}
	}
	return out
}

func (k Keeper) GetTotalStakedSupply(ctx sdk.Context, zone *types.Zone) math.Int {
	out := sdk.ZeroInt()
	for _, val := range k.GetActiveValidators(ctx, zone.ChainId) {
		if val.Status == stakingtypes.BondStatusBonded {
			out = out.Add(val.VotingPower)
		}
	}
	return out
}

func (k Keeper) CheckExceedsGlobalCap(ctx sdk.Context, zone *types.Zone, amount math.Int) bool {
	caps, found := k.GetLsmCaps(ctx, zone.ChainId)
	if !found {
		// no caps found, permit
		return false
	}

	liquidSupply := k.GetLiquidStakedSupply(ctx, zone)
	totalSupply := sdk.NewDecFromInt(k.GetTotalStakedSupply(ctx, zone))
	amountDec := sdk.NewDecFromInt(amount)
	return liquidSupply.Add(amountDec).Quo(totalSupply).GT(caps.GlobalCap)
}

func (k Keeper) CheckExceedsValidatorCap(ctx sdk.Context, zone *types.Zone, validator string, amount math.Int) error {
	// Retrieve the caps for the given zone
	caps, found := k.GetLsmCaps(ctx, zone.ChainId)
	if !found {
		// No cap found, permit the transaction
		return nil
	}

	// Retrieve the validator's information
	valAddrBytes, err := addressutils.ValAddressFromBech32(validator, zone.GetValoperPrefix())
	if err != nil {
		return err
	}

	val, found := k.GetValidator(ctx, zone.ChainId, valAddrBytes)
	if !found {
		// Validator not found, throw an error
		return errors.New("validator not found")
	}

	// Calculate the liquid shares and tokens
	amountDec := sdk.NewDecFromInt(amount)
	liquidShares := val.LiquidShares.Add(amountDec)
	tokens := sdk.NewDecFromInt(val.VotingPower).Add(amountDec)

	if liquidShares.Quo(tokens).GT(caps.ValidatorCap) {
		return errors.New("exceeds validator cap")
	}

	return nil
}

func (k Keeper) CheckExceedsValidatorBondCap(ctx sdk.Context, zone *types.Zone, validator string, amount math.Int) error {
	caps, found := k.GetLsmCaps(ctx, zone.ChainId)
	if !found {
		// no caps found, permit
		return nil
	}

	// Retrieve the validator's information
	valAddrBytes, err := addressutils.ValAddressFromBech32(validator, zone.GetValoperPrefix())
	if err != nil {
		return err
	}

	val, found := k.GetValidator(ctx, zone.ChainId, valAddrBytes)
	if !found {
		// cannot find validator, do not allow to proceed.
		return errors.New("validator not found")
	}

	maxShares := val.ValidatorBondShares.Mul(caps.ValidatorBondCap)

	amountDec := sdk.NewDecFromInt(amount)
	liquidShares := val.LiquidShares.Add(amountDec)

	if liquidShares.GT(maxShares) {
		return errors.New("exceeds validator bond cap")
	}

	return nil
}

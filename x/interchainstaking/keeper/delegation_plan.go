package keeper

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// gets the key for delegator bond with validator
// VALUE: staking/DelegationPlan
func GetDelegationPlanKey(zone *types.Zone, txhash string, delAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	return append(GetDelegationPlansKey(zone, txhash, delAddr), valAddr.Bytes()...)
}

// gets the prefix for a delegator for all validators
func GetDelegationPlansKey(zone *types.Zone, txhash string, delAddr sdk.AccAddress) []byte {
	return append(append(append(types.KeyPrefixDelegationPlan, []byte(zone.ChainId)...), []byte(txhash)...), delAddr.Bytes()...)
}

// GetDelegationPlan returns a specific delegation.
func (k Keeper) GetDelegationPlan(ctx sdk.Context, zone *types.Zone, txhash string, delegatorAddress string, validatorAddress string) (delegationPlan types.DelegationPlan, found bool) {
	store := ctx.KVStore(k.storeKey)

	_, delAddr, _ := bech32.DecodeAndConvert(delegatorAddress)
	_, valAddr, _ := bech32.DecodeAndConvert(validatorAddress)

	key := GetDelegationPlanKey(zone, txhash, delAddr, valAddr)

	value := store.Get(key)
	if value == nil {
		return delegationPlan, false
	}

	delegationPlan = types.MustUnmarshalDelegationPlan(k.cdc, value)

	return delegationPlan, true
}

// IterateAllDelegationPlansForHash iterates through all of the delegations for a given transaction.
func (k Keeper) IterateAllDelegationPlans(ctx sdk.Context, zone *types.Zone, cb func(delegationPlan types.DelegationPlan, key []byte) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, append(types.KeyPrefixDelegationPlan, []byte(zone.ChainId)...))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		delegationPlan := types.MustUnmarshalDelegationPlan(k.cdc, iterator.Value())
		if cb(delegationPlan, iterator.Key()) {
			break
		}
	}
}

func (k Keeper) GetAllDelegationPlans(ctx sdk.Context, zone *types.Zone) []types.DelegationPlan {
	out := []types.DelegationPlan{}
	k.IterateAllDelegationPlans(ctx, zone, func(delegationPlan types.DelegationPlan, _ []byte) bool {
		out = append(out, delegationPlan)
		return false
	})
	return out
}

func (k Keeper) GetAllDelegationPlansWithKey(ctx sdk.Context, zone *types.Zone) map[string]*types.DelegationPlan {
	out := map[string]*types.DelegationPlan{}
	k.IterateAllDelegationPlans(ctx, zone, func(delegationPlan types.DelegationPlan, key []byte) bool {
		keyString := string(key)
		parts := strings.Split(keyString, "/")
		out[parts[1]] = &delegationPlan
		return false
	})
	return out
}

// IterateAllDelegationPlansForHash iterates through all of the delegations for a given transaction.
func (k Keeper) IterateAllDelegationPlansForHash(ctx sdk.Context, zone *types.Zone, txhash string, cb func(delegationPlan types.DelegationPlan) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, append(append(types.KeyPrefixDelegationPlan, []byte(zone.ChainId)...), []byte(txhash)...))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		delegationPlan := types.MustUnmarshalDelegationPlan(k.cdc, iterator.Value())
		if cb(delegationPlan) {
			break
		}
	}
}

// IterateAllDelegationPlansForHashAndDelegator iterates through all of the delegations for a given transaction and delegator tuple.
func (k Keeper) IterateAllDelegationPlansForHashAndDelegator(ctx sdk.Context, zone *types.Zone, txhash string, delegatorAddr sdk.AccAddress, cb func(delegationPlan types.DelegationPlan) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, GetDelegationPlansKey(zone, txhash, delegatorAddr))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		delegationPlan := types.MustUnmarshalDelegationPlan(k.cdc, iterator.Value())
		if cb(delegationPlan) {
			break
		}
	}
}

// SetDelegationPlan sets a delegation.
func (k Keeper) SetDelegationPlan(ctx sdk.Context, zone *types.Zone, txhash string, delegationPlan types.DelegationPlan) {
	delegatorAddress := delegationPlan.GetDelegatorAddr()
	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalDelegationPlan(k.cdc, delegationPlan)
	store.Set(GetDelegationPlanKey(zone, txhash, delegatorAddress, delegationPlan.GetValidatorAddr()), b)
}

// RemoveDelegationPlan removes a delegation
func (k Keeper) RemoveDelegationPlan(ctx sdk.Context, zone *types.Zone, txhash string, delegationPlan types.DelegationPlan) error {
	delegatorAddress := delegationPlan.GetDelegatorAddr()

	store := ctx.KVStore(k.storeKey)
	store.Delete(GetDelegationPlanKey(zone, txhash, delegatorAddress, delegationPlan.GetValidatorAddr()))
	return nil
}

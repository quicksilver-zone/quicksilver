package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

func (*Keeper) getStoreKey(zone *types.Zone, snapshot bool) []byte {
	if snapshot {
		return append(types.KeyPrefixSnapshotIntent, zone.ChainId...)
	}
	return append(types.KeyPrefixIntent, zone.ChainId...)
}

// GetDelegatorIntent returns intent info by zone and delegator.
func (k *Keeper) GetDelegatorIntent(ctx sdk.Context, zone *types.Zone, delegator string, snapshot bool) (types.DelegatorIntent, bool) {
	intent := types.DelegatorIntent{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), k.getStoreKey(zone, snapshot))
	bz := store.Get([]byte(delegator))
	if len(bz) == 0 {
		// usually we'd return false here, but we always want to return an empty intent if one doesn't exist; keep standard Get* interface for consistency.
		return types.DelegatorIntent{Delegator: delegator, Intents: nil}, true
	}
	k.cdc.MustUnmarshal(bz, &intent)
	return intent, true
}

// SetDelegatorIntent store the delegator intent.
func (k *Keeper) SetDelegatorIntent(ctx sdk.Context, zone *types.Zone, intent types.DelegatorIntent, snapshot bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), k.getStoreKey(zone, snapshot))
	bz := k.cdc.MustMarshal(&intent)
	store.Set([]byte(intent.Delegator), bz)
}

// DeleteDelegatorIntent deletes delegator intent.
func (k *Keeper) DeleteDelegatorIntent(ctx sdk.Context, zone *types.Zone, delegator string, snapshot bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), k.getStoreKey(zone, snapshot))
	store.Delete([]byte(delegator))
}

// IterateDelegatorIntents iterate through delegator intents for a given zone.
func (k *Keeper) IterateDelegatorIntents(ctx sdk.Context, zone *types.Zone, snapshot bool, fn func(index int64, intent types.DelegatorIntent) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), k.getStoreKey(zone, snapshot))

	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	i := int64(0)

	for ; iterator.Valid(); iterator.Next() {
		intent := types.DelegatorIntent{}
		k.cdc.MustUnmarshal(iterator.Value(), &intent)

		stop := fn(i, intent)

		if stop {
			break
		}
		i++
	}
}

// AllDelegatorIntents returns every intent in the store for the specified zone.
func (k *Keeper) AllDelegatorIntents(ctx sdk.Context, zone *types.Zone, snapshot bool) []types.DelegatorIntent {
	var intents []types.DelegatorIntent
	k.IterateDelegatorIntents(ctx, zone, snapshot, func(_ int64, intent types.DelegatorIntent) (stop bool) {
		intents = append(intents, intent)
		return false
	})
	return intents
}

// AllDelegatorIntentsAsPointer returns every intent in the store for the specified zone.
func (k *Keeper) AllDelegatorIntentsAsPointer(ctx sdk.Context, zone *types.Zone, snapshot bool) []*types.DelegatorIntent {
	var intents []*types.DelegatorIntent
	k.IterateDelegatorIntents(ctx, zone, snapshot, func(_ int64, intent types.DelegatorIntent) (stop bool) {
		intents = append(intents, &intent)
		return false
	})
	return intents
}

// UpdateDelegatorIntent updates delegator intents.
func (k *Keeper) UpdateDelegatorIntent(ctx sdk.Context, delegator sdk.AccAddress, zone *types.Zone, inAmount sdk.Coins, memoIntent types.ValidatorIntents) error {
	return nil
}

func (k msgServer) validateValidatorIntents(ctx sdk.Context, zone types.Zone, intents []*types.ValidatorIntent) error {
	return nil
}

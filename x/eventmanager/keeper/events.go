package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/x/eventmanager/types"
)

// ----------------------------------------------------------------

func GenerateEventKey(module, chainId, id string) []byte {
	return []byte(module + chainId + id)
}

// GetEvent returns event.
func (k Keeper) GetEvent(ctx sdk.Context, module, chainId, id string) (types.Event, bool) {
	event := types.Event{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixEvent)
	bz := store.Get(GenerateEventKey(module, chainId, id))
	if len(bz) == 0 {
		return event, false
	}
	k.cdc.MustUnmarshal(bz, &event)
	return event, true
}

// SetEvent set event.
func (k Keeper) SetEvent(ctx sdk.Context, event types.Event) {
	key := GenerateEventKey(event.Module, event.ChainId, event.Identifier)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixEvent)
	bz := k.cdc.MustMarshal(&event)
	store.Set(key, bz)
}

// DeleteEvent delete event.
func (k Keeper) DeleteEvent(ctx sdk.Context, module, chainId, id string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixEvent)
	store.Delete(GenerateEventKey(module, chainId, id))
}

// IterateEvents iterate through queries.
func (k Keeper) IteratePrefixedEvents(ctx sdk.Context, prefixBytes []byte, fn func(index int64, event types.Event) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixEvent)
	iterator := sdk.KVStorePrefixIterator(store, prefixBytes)
	defer iterator.Close()

	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		event := types.Event{}
		k.cdc.MustUnmarshal(iterator.Value(), &event)
		stop := fn(i, event)

		if stop {
			break
		}
		i++
	}
}

func (k Keeper) IterateModuleEvents(ctx sdk.Context, module string, fn func(index int64, event types.Event) (stop bool)) {
	k.IteratePrefixedEvents(ctx, []byte(module), fn)
}

func (k Keeper) IterateModuleChainEvents(ctx sdk.Context, module string, chainId string, fn func(index int64, event types.Event) (stop bool)) {
	k.IteratePrefixedEvents(ctx, []byte(module+chainId), fn)
}

// AllEvents returns every eventInfo in the store.
func (k Keeper) AllEvents(ctx sdk.Context) []types.Event {
	queries := []types.Event{}
	k.IteratePrefixedEvents(ctx, nil, func(_ int64, eventInfo types.Event) (stop bool) {
		queries = append(queries, eventInfo)
		return false
	})
	return queries
}

func (k Keeper) MarkCompleted(ctx sdk.Context, module string, chainID string, identifier string) {
	k.DeleteEvent(ctx, module, chainID, identifier)
	k.Trigger(ctx, module, chainID)
}

func (k Keeper) Trigger(ctx sdk.Context, module string, chainID string) {
	k.IterateModuleChainEvents(ctx, module, chainID, func(_ int64, e types.Event) (stop bool) {
		if e.EventStatus == types.EventStatusPending {
			err := k.Call(ctx, e.Module, e.Callback, e.Payload)
			if err != nil {
				k.Logger(ctx).Error("unable to execute callback", "module", e.Module, "id", e.Identifier, "callback", e.Callback, "error", err)
			}
		}
		return false
	})
}

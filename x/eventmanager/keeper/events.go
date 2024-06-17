package keeper

import (
	"fmt"
	"time"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/x/eventmanager/types"
)

// ----------------------------------------------------------------

func GenerateEventKey(module, chainID, id string) []byte {
	return []byte(module + chainID + id)
}

// GetEvent returns event.
func (k Keeper) GetEvent(ctx sdk.Context, module, chainID, id string) (types.Event, bool) {
	event := types.Event{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixEvent)
	bz := store.Get(GenerateEventKey(module, chainID, id))
	if len(bz) == 0 {
		return event, false
	}
	k.cdc.MustUnmarshal(bz, &event)
	return event, true
}

func (k Keeper) GetEvents(ctx sdk.Context, module, chainID, prefix string) ([]types.Event, int) {
	events := make([]types.Event, 0)

	k.IteratePrefixedEvents(ctx, []byte(module+chainID+prefix), func(index int64, event types.Event) (stop bool) {
		events = append(events, event)
		return false
	})

	return events, len(events)
}

// SetEvent set event.
func (k Keeper) SetEvent(ctx sdk.Context, event types.Event) {
	key := GenerateEventKey(event.Module, event.ChainId, event.Identifier)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixEvent)
	bz := k.cdc.MustMarshal(&event)
	store.Set(key, bz)
}

// DeleteEvent delete event.
func (k Keeper) DeleteEvent(ctx sdk.Context, module, chainID, id string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixEvent)
	store.Delete(GenerateEventKey(module, chainID, id))
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

func (k Keeper) IterateModuleChainEvents(ctx sdk.Context, module string, chainID string, fn func(index int64, event types.Event) (stop bool)) {
	k.IteratePrefixedEvents(ctx, []byte(module+chainID), fn)
}

// AllEvents returns every eventInfo in the store.
func (k Keeper) AllEvents(ctx sdk.Context) []types.Event {
	events := []types.Event{}
	k.IteratePrefixedEvents(ctx, nil, func(_ int64, eventInfo types.Event) (stop bool) {
		events = append(events, eventInfo)
		return false
	})
	return events
}

func (k Keeper) MarkCompleted(ctx sdk.Context, module string, chainID string, identifier string) {
	k.Logger(ctx).Info(fmt.Sprintf("marking event %s/%s/%s as complete!", module, chainID, identifier))
	k.DeleteEvent(ctx, module, chainID, identifier)
	k.Trigger(ctx, module, chainID)
}

func (k Keeper) GetTriggerFn(ctx sdk.Context) func(_ int64, e types.Event) (stop bool) {
	return func(_ int64, e types.Event) (stop bool) {
		if e.EventStatus == types.EventStatusPending {
			res, err := e.CanExecute(ctx, &k)
			if err != nil {
				k.Logger(ctx).Error("unable to determine if event can execute callback", "module", e.Module, "id", e.Identifier, "callback", e.Callback, "error", err)
			}
			if res {
				k.Logger(ctx).Info(fmt.Sprintf("triggered event callback %s for event %s (%s)", e.Callback, e.Identifier, e.ChainId))
				err := k.Call(ctx, e.Module, e.Callback, e.Payload)
				if err != nil {
					k.Logger(ctx).Error("unable to execute callback", "module", e.Module, "id", e.Identifier, "callback", e.Callback, "error", err)
				}
				e.EventStatus = types.EventStatusActive
				k.SetEvent(ctx, e)
			}
		}
		return false
	}
}

func (k Keeper) Trigger(ctx sdk.Context, module string, chainID string) {
	k.IterateModuleChainEvents(ctx, module, chainID, k.GetTriggerFn(ctx))
}

func (k Keeper) TriggerAll(ctx sdk.Context) {
	k.IteratePrefixedEvents(ctx, nil, k.GetTriggerFn(ctx))
}

func (k Keeper) AddEventWithExpiry(ctx sdk.Context, module, chainID, identifier string, eventType, status int32, expiry time.Time) {
	// expiring events cannot have callbacks
	event := types.Event{
		ChainId:          chainID,
		Module:           module,
		Identifier:       identifier,
		EventType:        eventType,
		Callback:         "",
		Payload:          nil,
		EventStatus:      status,
		ExecuteCondition: nil,
		EmittedHeight:    ctx.BlockHeight(),
		ExpiryTime:       &expiry,
	}

	k.SetEvent(ctx, event)
}

func (k Keeper) AddEvent(ctx sdk.Context,
	module, chainID, identifier, callback string,
	eventType, status int32,
	condition types.ConditionI,
	payload []byte,
) {
	var err error
	var conditionAny *codectypes.Any
	if condition != nil {
		conditionAny, err = codectypes.NewAnyWithValue(condition)
		if err != nil {
			panic(err)
		}
	}

	event := types.Event{
		ChainId:          chainID,
		Module:           module,
		Identifier:       identifier,
		EventType:        eventType,
		Callback:         callback,
		Payload:          payload,
		EventStatus:      status,
		ExecuteCondition: conditionAny,
		EmittedHeight:    ctx.BlockHeight(),
	}

	k.SetEvent(ctx, event)
}

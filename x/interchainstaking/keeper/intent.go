package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// GetIntent returns intent info by zone and delegator
func (k Keeper) GetIntent(ctx sdk.Context, zone types.RegisteredZone, delegator string) (types.DelegatorIntent, bool) {
	intent := types.DelegatorIntent{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), append(types.KeyPrefixIntent, []byte(zone.ChainId)...))
	bz := store.Get([]byte(delegator))
	if len(bz) == 0 {
		// usually we'd return false here, but we always want to return an empty intent if one doesn't exist; keep standard Get* interface for consistency.
		return types.DelegatorIntent{Delegator: delegator, Intents: []*types.ValidatorIntent{}}, true
	}
	k.cdc.MustUnmarshal(bz, &intent)
	return intent, true
}

// SetIntent store the delegator intent
func (k Keeper) SetIntent(ctx sdk.Context, zone types.RegisteredZone, intent types.DelegatorIntent) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), append(types.KeyPrefixIntent, []byte(zone.ChainId)...))
	bz := k.cdc.MustMarshal(&intent)
	ctx.Logger().Error(fmt.Sprintf("Writing the intent for chain %s for delegator %s: %v", zone.ChainId, intent.Delegator, intent))
	store.Set([]byte(intent.Delegator), bz)
}

// DeleteIntent deletes delegator intent
func (k Keeper) DeleteIntent(ctx sdk.Context, zone types.RegisteredZone, delegator string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), append(types.KeyPrefixIntent, []byte(zone.ChainId)...))
	store.Delete([]byte(delegator))
}

// IterateIntents iterate through intents for a given zone
func (k Keeper) IterateIntents(ctx sdk.Context, zone types.RegisteredZone, fn func(index int64, intent types.DelegatorIntent) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), append(types.KeyPrefixIntent, []byte(zone.ChainId)...))

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

// AllIntents returns every intent in the store for the specified zone
func (k Keeper) AllIntents(ctx sdk.Context, zone types.RegisteredZone) []types.DelegatorIntent {
	intents := []types.DelegatorIntent{}
	k.IterateIntents(ctx, zone, func(_ int64, intent types.DelegatorIntent) (stop bool) {
		intents = append(intents, intent)
		return false
	})
	return intents
}

// AllOrdinalizedIntents returns every intent in the store for the specified zone
func (k Keeper) AllOrdinalizedIntents(ctx sdk.Context, zone types.RegisteredZone) []types.DelegatorIntent {
	intents := []types.DelegatorIntent{}
	k.IterateIntents(ctx, zone, func(_ int64, intent types.DelegatorIntent) (stop bool) {
		query := bankTypes.QueryBalanceRequest{Address: intent.Delegator, Denom: zone.LocalDenom}
		balance, err := k.BankKeeper.Balance(sdk.WrapSDKContext(ctx), &query)
		if err != nil {
			panic(err)
		}
		baseBalance := zone.RedemptionRate.Mul(sdk.NewDecFromInt(balance.Balance.Amount)).TruncateInt()
		intents = append(intents, intent.Ordinalize(baseBalance))
		return false
	})
	return intents
}

func (k *Keeper) UpdateIntent(ctx sdk.Context, sender sdk.AccAddress, zone types.RegisteredZone, inAmount sdk.Coins) {
	// this is here because we need access to the bankKeeper to ordinalize intent
	intent, _ := k.GetIntent(ctx, zone, sender.String())

	// ordinalize
	query := bankTypes.QueryBalanceRequest{Address: sender.String(), Denom: zone.LocalDenom}
	balance, err := k.BankKeeper.Balance(sdk.WrapSDKContext(ctx), &query)
	if err != nil {
		panic(err)
	}
	baseBalance := zone.RedemptionRate.Mul(sdk.NewDecFromInt(balance.Balance.Amount)).TruncateInt()
	intent = intent.AddOrdinal(baseBalance, zone.ConvertCoinsToOrdinalIntents(ctx, inAmount))
	k.SetIntent(ctx, zone, intent)
}

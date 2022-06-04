package keeper

import (
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

func (k *Keeper) AggregateIntents(ctx sdk.Context, zone types.RegisteredZone) {
	intents := map[string]*types.ValidatorIntent{}
	ordinalizedIntentSum := sdk.ZeroDec()
	k.IterateIntents(ctx, zone, func(_ int64, intent types.DelegatorIntent) (stop bool) {
		query := bankTypes.QueryBalanceRequest{Address: intent.Delegator, Denom: zone.LocalDenom}
		balance, err := k.BankKeeper.Balance(sdk.WrapSDKContext(ctx), &query)
		if err != nil {
			panic(err)
		}
		baseBalance := zone.RedemptionRate.Mul(sdk.NewDecFromInt(balance.Balance.Amount)).TruncateInt()
		for _, vIntent := range intent.Ordinalize(baseBalance).Intents {
			thisIntent, ok := intents[vIntent.ValoperAddress]
			ordinalizedIntentSum = ordinalizedIntentSum.Add(vIntent.Weight)
			if !ok {
				intents[vIntent.ValoperAddress] = vIntent
			} else {
				thisIntent.Weight = thisIntent.Weight.Add(vIntent.Weight)
				intents[vIntent.ValoperAddress] = thisIntent
			}
		}

		return false
	})

	for key, val := range intents {
		val.Weight = val.Weight.Quo(ordinalizedIntentSum)
		intents[key] = val
	}

	zone.AggregateIntent = intents
	k.SetRegisteredZone(ctx, zone)
}

func (k *Keeper) UpdateIntent(ctx sdk.Context, sender sdk.AccAddress, zone types.RegisteredZone, inAmount sdk.Coins, memo string) {
	// this is here because we need access to the bankKeeper to ordinalize intent
	intent, _ := k.GetIntent(ctx, zone, sender.String())

	// ordinalize
	query := bankTypes.QueryBalanceRequest{Address: sender.String(), Denom: zone.LocalDenom}
	balance, err := k.BankKeeper.Balance(sdk.WrapSDKContext(ctx), &query)
	if err != nil {
		panic(err)
	}
	baseBalance := zone.RedemptionRate.Mul(sdk.NewDecFromInt(balance.Balance.Amount)).TruncateInt()
	intent = intent.AddOrdinal(baseBalance, zone.ConvertCoinsToOrdinalIntents(inAmount))
	intent = intent.AddOrdinal(baseBalance, zone.ConvertMemoToOrdinalIntents(inAmount, memo))
	if len(intent.Intents) == 0 {
		return
	}
	k.SetIntent(ctx, zone, intent)
}

package keeper

import (
	"errors"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	prtypes "github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

func (k Keeper) getStoreKey(zone types.Zone, snapshot bool) []byte {
	if snapshot {
		return append(types.KeyPrefixSnapshotIntent, []byte(zone.ChainId)...)
	}
	return append(types.KeyPrefixIntent, []byte(zone.ChainId)...)
}

// GetIntent returns intent info by zone and delegator
func (k Keeper) GetIntent(ctx sdk.Context, zone types.Zone, delegator string, snapshot bool) (types.DelegatorIntent, bool) {
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

// SetIntent store the delegator intent
func (k Keeper) SetIntent(ctx sdk.Context, zone types.Zone, intent types.DelegatorIntent, snapshot bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), k.getStoreKey(zone, snapshot))
	bz := k.cdc.MustMarshal(&intent)
	store.Set([]byte(intent.Delegator), bz)
}

// DeleteIntent deletes delegator intent
func (k Keeper) DeleteIntent(ctx sdk.Context, zone types.Zone, delegator string, snapshot bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), k.getStoreKey(zone, snapshot))
	store.Delete([]byte(delegator))
}

// IterateIntents iterate through intents for a given zone
func (k Keeper) IterateIntents(ctx sdk.Context, zone types.Zone, snapshot bool, fn func(index int64, intent types.DelegatorIntent) (stop bool)) {
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

// AllIntents returns every intent in the store for the specified zone
func (k Keeper) AllIntents(ctx sdk.Context, zone types.Zone, snapshot bool) []types.DelegatorIntent {
	intents := []types.DelegatorIntent{}
	k.IterateIntents(ctx, zone, snapshot, func(_ int64, intent types.DelegatorIntent) (stop bool) {
		intents = append(intents, intent)
		return false
	})
	return intents
}

// AllIntents returns every intent in the store for the specified zone
func (k Keeper) AllIntentsAsPointer(ctx sdk.Context, zone types.Zone, snapshot bool) []*types.DelegatorIntent {
	intents := []*types.DelegatorIntent{}
	k.IterateIntents(ctx, zone, snapshot, func(_ int64, intent types.DelegatorIntent) (stop bool) {
		intents = append(intents, &intent)
		return false
	})
	return intents
}

// AllOrdinalizedIntents returns every intent in the store for the specified zone.
func (k Keeper) AllOrdinalizedIntents(ctx sdk.Context, zone types.Zone, snapshot bool) ([]types.DelegatorIntent, error) {
	intents := []types.DelegatorIntent{}
	var err error
	k.IterateIntents(ctx, zone, snapshot, func(_ int64, intent types.DelegatorIntent) (stop bool) {
		addr, localErr := sdk.AccAddressFromBech32(intent.Delegator)
		if localErr != nil {
			err = localErr
			return true
		}
		balance := k.BankKeeper.GetBalance(ctx, addr, zone.LocalDenom)

		intents = append(intents, intent.Ordinalize(sdk.NewDecFromInt(balance.Amount)))
		return false
	})
	if err != nil {
		// check on nil here to ensure we don't return half a slice of intents
		return []types.DelegatorIntent{}, err
	}
	return intents, nil
}

func (k *Keeper) AggregateIntents(ctx sdk.Context, zone types.Zone) error {
	var err error
	snapshot := false
	aggregate := make(types.ValidatorIntents, 0)
	ordinalizedIntentSum := sdk.ZeroDec()
	// reduce intents

	k.IterateIntents(ctx, zone, snapshot, func(_ int64, intent types.DelegatorIntent) (stop bool) {
		addr, localErr := sdk.AccAddressFromBech32(intent.Delegator)
		if localErr != nil {
			err = localErr
			return true
		}
		balance := k.BankKeeper.GetBalance(ctx, addr, zone.LocalDenom)

		// grab offchain asset value, and raise the users' base value by this amount.
		k.ParticipationRewardsKeeper.IterateLastEpochUserClaims(ctx, zone.ChainId, intent.Delegator, func(index int64, data prtypes.Claim) (stop bool) {
			balance.Amount = balance.Amount.Add(math.NewIntFromUint64(data.Amount))
			return false
		})

		intents := intent.Ordinalize(sdk.NewDecFromInt(balance.Amount)).Intents
		for vIntent := range intents.Sort() {
			thisIntent, ok := aggregate.GetForValoper(intents[vIntent].ValoperAddress)
			ordinalizedIntentSum = ordinalizedIntentSum.Add(intents[vIntent].Weight)
			if !ok {
				aggregate = append(aggregate, intents[vIntent])
			} else {
				thisIntent.Weight = thisIntent.Weight.Add(intents[vIntent].Weight)
				aggregate = aggregate.SetForValoper(intents[vIntent].ValoperAddress, thisIntent)
			}
		}

		return false
	})
	if err != nil {
		return err
	}

	if len(aggregate) > 0 && ordinalizedIntentSum.IsZero() {
		return errors.New("ordinalized intent sum is zero, this should never happen")
	}

	// normalise aggregated intents again.
	for idx, intent := range aggregate.Sort() {
		intent.Weight = intent.Weight.Quo(ordinalizedIntentSum)
		aggregate[idx] = intent
	}

	k.Logger(ctx).Info("aggregates", "agg", aggregate, "chain", zone.ChainId)

	zone.AggregateIntent = aggregate
	k.SetZone(ctx, &zone)
	return nil
}

func (k *Keeper) UpdateIntent(ctx sdk.Context, sender sdk.AccAddress, zone types.Zone, inAmount sdk.Coins, memo string) error {
	snapshot := false
	// this is here because we need access to the bankKeeper to ordinalize intent
	intent, _ := k.GetIntent(ctx, zone, sender.String(), snapshot)

	// ordinalize
	balance := k.BankKeeper.GetBalance(ctx, sender, zone.BaseDenom)
	if balance.Amount.IsNil() {
		balance.Amount = math.NewInt(0)
	}

	k.Logger(ctx).Error("DEBUG", "zone", zone.ChainId)
	k.Logger(ctx).Error("DEBUG", "sender", sender.String())
	k.Logger(ctx).Error("DEBUG", "balance", balance)
	k.Logger(ctx).Error("DEBUG", "prkeeper", k.ParticipationRewardsKeeper)

	// grab offchain asset value, and raise the users' base value by this amount.
	k.ParticipationRewardsKeeper.IterateLastEpochUserClaims(ctx, zone.ChainId, sender.String(), func(index int64, data prtypes.Claim) (stop bool) {
		k.Logger(ctx).Error("DEBUG", "claim", data)
		balance.Amount = balance.Amount.Add(math.NewIntFromUint64(data.Amount))
		return false
	})

	// inAmount is ordinal with respect to the redemption rate, so we must scale
	baseBalance := zone.RedemptionRate.Mul(sdk.NewDecFromInt(balance.Amount))

	if inAmount.IsValid() {
		intent = zone.UpdateIntentWithCoins(intent, baseBalance, inAmount)
	}

	if len(memo) > 0 {
		var err error
		intent, err = zone.UpdateIntentWithMemo(intent, memo, baseBalance, inAmount)
		if err != nil {
			return err
		}
	}

	if len(intent.Intents) == 0 {
		return nil
	}

	k.SetIntent(ctx, zone, intent, snapshot)
	return nil
}

package keeper

import (
	"errors"
	"fmt"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/internal/multierror"
	prtypes "github.com/ingenuity-build/quicksilver/x/claimsmanager/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
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

func (k *Keeper) AggregateIntents(ctx sdk.Context, zone *types.Zone) error {
	var err error
	snapshot := false
	aggregate := make(types.ValidatorIntents, 0)
	ordinalizedIntentSum := sdk.ZeroDec()
	// reduce intents

	k.IterateIntents(ctx, *zone, snapshot, func(_ int64, intent types.DelegatorIntent) (stop bool) {
		// addr, localErr := sdk.AccAddressFromBech32(intent.Delegator)
		// if localErr != nil {
		// 	err = localErr
		// 	return true
		// }
		// balance := k.BankKeeper.GetBalance(ctx, addr, zone.LocalDenom)
		balance := sdk.NewCoin(zone.LocalDenom, math.ZeroInt())

		// grab offchain asset value, and raise the users' base value by this amount.
		// currently ignoring base value (locally held assets)
		k.ClaimsManagerKeeper.IterateLastEpochUserClaims(ctx, zone.ChainId, intent.Delegator, func(index int64, data prtypes.Claim) (stop bool) {
			balance.Amount = balance.Amount.Add(math.NewIntFromUint64(data.Amount))
			// claim amounts are in zone.baseDenom - but given weights are all relative to one another this okay.
			k.Logger(ctx).Error("Intents - found claim for user", "user", intent.Delegator, "claim amount", data.Amount, "new balance", balance.Amount)
			return false
		})

		intents := intent.Ordinalize(sdk.NewDecFromInt(balance.Amount)).Intents
		k.Logger(ctx).Error("Intents - ordinalized", "user", intent.Delegator, "new balance", balance.Amount, "normal intents", intent.Intents, "intents", intents)

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
		return errors.New("ordinalized intent sum is zero, this may happen if no claims are recorded")
	}

	// normalise aggregated intents again.
	newAggregate := make(types.ValidatorIntents, 0)
	for _, intent := range aggregate.Sort() {
		if !intent.Weight.IsZero() && intent.Weight.IsPositive() {
			intent.Weight = intent.Weight.Quo(ordinalizedIntentSum)
			newAggregate = append(newAggregate, intent)
		}
	}

	k.Logger(ctx).Info("aggregates", "agg", newAggregate, "chain", zone.ChainId)

	zone.AggregateIntent = newAggregate
	k.SetZone(ctx, zone)
	return nil
}

func (k *Keeper) UpdateIntent(ctx sdk.Context, sender sdk.AccAddress, zone types.Zone, inAmount sdk.Coins, memo string) error {
	snapshot := false
	// this is here because we need access to the bankKeeper to ordinalize intent
	intent, _ := k.GetIntent(ctx, zone, sender.String(), snapshot)

	// ordinalize
	// this is the currently held amount
	// not aligned with last epoch claims
	// balance := k.BankKeeper.GetBalance(ctx, sender, zone.BaseDenom)
	// if balance.Amount.IsNil() {
	// 	balance.Amount = math.ZeroInt()
	// }
	balance := sdk.NewCoin(zone.BaseDenom, math.ZeroInt())

	// grab offchain asset value, and raise the users' base value by this amount.
	k.ClaimsManagerKeeper.IterateLastEpochUserClaims(ctx, zone.ChainId, sender.String(), func(index int64, data prtypes.Claim) (stop bool) {
		k.Logger(ctx).Error("Update intents - found claim for user", "user", intent.Delegator, "claim amount", data.Amount, "new balance", balance.Amount)

		balance.Amount = balance.Amount.Add(math.NewIntFromUint64(data.Amount))
		return false
	})

	// inAmount is ordinal with respect to the redemption rate, so we must scale
	baseBalance := zone.RedemptionRate.Mul(sdk.NewDecFromInt(balance.Amount))
	if baseBalance.IsZero() {
		return nil
	}

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

func (k msgServer) validateIntents(zone types.Zone, intents []*types.ValidatorIntent) error {
	errors := make(map[string]error)

	for i, intent := range intents {
		_, found := zone.GetValidatorByValoper(intent.ValoperAddress)
		if !found {
			errors[fmt.Sprintf("intent[%v]", i)] = fmt.Errorf("unable to find valoper %s", intent.ValoperAddress)
		}
	}

	if len(errors) > 0 {
		return multierror.New(errors)
	}

	return nil
}

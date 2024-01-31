package keeper

import (
	"errors"
	"fmt"

	"github.com/ingenuity-build/multierror"

	sdkmath "cosmossdk.io/math"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/quicksilver-zone/quicksilver/v7/utils"
	"github.com/quicksilver-zone/quicksilver/v7/utils/addressutils"
	prtypes "github.com/quicksilver-zone/quicksilver/v7/x/claimsmanager/types"
	"github.com/quicksilver-zone/quicksilver/v7/x/interchainstaking/types"
)

func (*Keeper) getStoreKey(zone *types.Zone, snapshot bool) []byte {
	if snapshot {
		return append(types.KeyPrefixSnapshotIntent, []byte(zone.ChainId)...)
	}
	return append(types.KeyPrefixIntent, []byte(zone.ChainId)...)
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

	iterator := storetypes.KVStorePrefixIterator(store, nil)
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

// AggregateDelegatorIntents takes a snapshot of delegator intents for a given zone.
func (k *Keeper) AggregateDelegatorIntents(ctx sdk.Context, zone *types.Zone) error {
	snapshot := false
	aggregate := make(types.ValidatorIntents, 0)
	ordinalizedIntentSum := sdkmath.LegacyZeroDec()

	k.IterateDelegatorIntents(ctx, zone, snapshot, func(_ int64, delIntent types.DelegatorIntent) (stop bool) {
		balance := sdk.NewCoin(zone.LocalDenom, sdkmath.ZeroInt())
		// grab offchain asset value, and raise the users' base value by this amount.
		// currently ignoring base value (locally held assets)
		k.ClaimsManagerKeeper.IterateLastEpochUserClaims(ctx, zone.ChainId, delIntent.Delegator, func(index int64, data prtypes.Claim) (stop bool) {
			balance.Amount = balance.Amount.Add(sdkmath.NewIntFromUint64(data.Amount))
			// claim amounts are in zone.baseDenom - but given weights are all relative to one another this okay.
			k.Logger(ctx).Debug(
				"intents - found claim for user",
				"user", delIntent.Delegator,
				"claim amount", data.Amount,
				"new balance", balance.Amount,
			)
			return false
		})

		valIntents := delIntent.Ordinalize(sdkmath.LegacyNewDecFromInt(balance.Amount)).Intents
		k.Logger(ctx).Debug(
			"intents - ordinalized",
			"user", delIntent.Delegator,
			"new balance", balance.Amount,
			"normal intents", delIntent.Intents,
			"intents", valIntents,
		)

		for idx := range valIntents.Sort() {
			valIntent, found := aggregate.GetForValoper(valIntents[idx].ValoperAddress)
			ordinalizedIntentSum = ordinalizedIntentSum.Add(valIntents[idx].Weight)
			if !found {
				aggregate = append(aggregate, valIntents[idx])
			} else {
				valIntent.Weight = valIntent.Weight.Add(valIntents[idx].Weight)
				aggregate = aggregate.SetForValoper(valIntents[idx].ValoperAddress, valIntent)
			}
		}

		return false
	})

	// weight supply for which we do not have claim equally across active validators.
	// this stops a small number of claimants exercising a disproportionate amount of
	// power, in the event claims cannot be made properly.
	supply := k.BankKeeper.GetSupply(ctx, zone.LocalDenom)
	defaults := k.DefaultAggregateIntents(ctx, zone.ChainId)
	nonVotingSupply := sdkmath.LegacyNewDecFromInt(supply.Amount).Sub(ordinalizedIntentSum)
	di := types.DelegatorIntent{Delegator: "", Intents: defaults}
	di = di.Ordinalize(nonVotingSupply)
	defaults = di.Intents

	for idx := range defaults.Sort() {
		valIntent, found := aggregate.GetForValoper(defaults[idx].ValoperAddress)
		ordinalizedIntentSum = ordinalizedIntentSum.Add(defaults[idx].Weight)
		if !found {
			aggregate = append(aggregate, defaults[idx])
		} else {
			valIntent.Weight = valIntent.Weight.Add(defaults[idx].Weight)
			aggregate = aggregate.SetForValoper(defaults[idx].ValoperAddress, valIntent)
		}
	}

	if len(aggregate) > 0 && !ordinalizedIntentSum.IsPositive() {
		return errors.New("ordinalized intent sum is zero, this may happen if no claims are recorded")
	}

	// normalise aggregated intents again.
	newAggregate := make(types.ValidatorIntents, 0)
	for _, valIntent := range aggregate.Sort() {
		if valIntent.Weight.IsPositive() {
			valIntent.Weight = valIntent.Weight.Quo(ordinalizedIntentSum)
			newAggregate = append(newAggregate, valIntent)
		}
	}

	k.Logger(ctx).Info(
		"aggregates",
		"agg", newAggregate,
		"chain", zone.ChainId,
	)

	zone.AggregateIntent = newAggregate
	k.SetZone(ctx, zone)
	return nil
}

// UpdateDelegatorIntent updates delegator intents.
func (k *Keeper) UpdateDelegatorIntent(ctx sdk.Context, delegator sdk.AccAddress, zone *types.Zone, inAmount sdk.Coins, memoIntent types.ValidatorIntents) error {
	snapshot := false
	updateWithCoin := inAmount.IsValid()
	updateWithMemo := memoIntent != nil

	// this is here because we need access to the bankKeeper to ordinalize intent
	delIntent, _ := k.GetDelegatorIntent(ctx, zone, delegator.String(), snapshot)

	// ordinalize
	// this is the currently held amount
	// not aligned with last epoch claims
	// balance := k.BankKeeper.GetBalance(ctx, sender, zone.BaseDenom)
	// if balance.Amount.IsNil() {
	// 	balance.Amount = math.ZeroInt()
	// }
	claimAmt := sdkmath.ZeroInt()

	// grab offchain asset value, and raise the users' base value by this amount.
	k.ClaimsManagerKeeper.IterateLastEpochUserClaims(ctx, zone.ChainId, delegator.String(), func(index int64, claim prtypes.Claim) (stop bool) {
		claimAmt = claimAmt.Add(sdkmath.NewIntFromUint64(claim.Amount))
		k.Logger(ctx).Error("Update intents - found claim for user", "user", delIntent.Delegator, "claim amount", claim.Amount, "new balance", claimAmt)
		return false
	})

	// inAmount is ordinal with respect to the redemption rate, so we must scale
	baseBalance := zone.RedemptionRate.Mul(sdkmath.LegacyNewDecFromInt(claimAmt))
	if baseBalance.IsZero() {
		return nil
	}

	if updateWithCoin {
		delIntent = zone.UpdateIntentWithCoins(delIntent, baseBalance, inAmount, utils.StringSliceToMap(k.GetValidatorAddresses(ctx, zone.ChainId)))
	}

	if updateWithMemo {
		delIntent = zone.UpdateZoneIntentWithMemo(memoIntent, delIntent, baseBalance)
	}

	if len(delIntent.Intents) == 0 {
		return nil
	}

	k.SetDelegatorIntent(ctx, zone, delIntent, snapshot)

	return nil
}

func (k msgServer) validateValidatorIntents(ctx sdk.Context, zone types.Zone, intents []*types.ValidatorIntent) error {
	errMap := make(map[string]error)

	for i, intent := range intents {
		var valAddrBytes []byte
		valAddrBytes, err := addressutils.ValAddressFromBech32(intent.ValoperAddress, zone.GetValoperPrefix())
		if err != nil {
			return err
		}
		_, found := k.GetValidator(ctx, zone.ChainId, valAddrBytes)
		if !found {
			errMap[fmt.Sprintf("intent[%v]", i)] = fmt.Errorf("unable to find valoper %s", intent.ValoperAddress)
		}
	}

	if len(errMap) > 0 {
		return multierror.New(errMap)
	}

	return nil
}

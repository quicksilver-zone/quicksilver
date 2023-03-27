package keeper

import (
	"errors"
	"math"
	"sort"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	lsmstakingTypes "github.com/iqlusioninc/liquidity-staking-module/x/staking/types"

	"github.com/ingenuity-build/quicksilver/utils"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// gets the key for delegator bond with validator
// VALUE: staking/Delegation
func GetDelegationKey(zone *types.Zone, delAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	return append(GetDelegationsKey(zone, delAddr), valAddr.Bytes()...)
}

// gets the prefix for a delegator for all validators
func GetDelegationsKey(zone *types.Zone, delAddr sdk.AccAddress) []byte {
	return append(append(types.KeyPrefixDelegation, []byte(zone.ChainId)...), delAddr.Bytes()...)
}

// gets the key for delegator bond with validator
// VALUE: staking/Delegation
func GetPerformanceDelegationKey(zone *types.Zone, delAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	return append(GetPerformanceDelegationsKey(zone, delAddr), valAddr.Bytes()...)
}

// gets the prefix for a delegator for all validators
func GetPerformanceDelegationsKey(zone *types.Zone, delAddr sdk.AccAddress) []byte {
	return append(append(types.KeyPrefixPerformanceDelegation, []byte(zone.ChainId)...), delAddr.Bytes()...)
}

// GetDelegation returns a specific delegation.
func (k Keeper) GetDelegation(ctx sdk.Context, zone *types.Zone, delegatorAddress string, validatorAddress string) (delegation types.Delegation, found bool) {
	store := ctx.KVStore(k.storeKey)

	_, delAddr, _ := bech32.DecodeAndConvert(delegatorAddress)
	_, valAddr, _ := bech32.DecodeAndConvert(validatorAddress)

	key := GetDelegationKey(zone, delAddr, valAddr)

	value := store.Get(key)
	if value == nil {
		return delegation, false
	}

	delegation = types.MustUnmarshalDelegation(k.cdc, value)

	return delegation, true
}

// GetDelegation returns a specific delegation.
func (k Keeper) GetPerformanceDelegation(ctx sdk.Context, zone *types.Zone, validatorAddress string) (delegation types.Delegation, found bool) {
	if zone.PerformanceAddress == nil {
		return types.Delegation{}, false
	}

	store := ctx.KVStore(k.storeKey)

	_, delAddr, _ := bech32.DecodeAndConvert(zone.PerformanceAddress.Address)
	_, valAddr, _ := bech32.DecodeAndConvert(validatorAddress)

	key := GetPerformanceDelegationKey(zone, delAddr, valAddr)

	value := store.Get(key)
	if value == nil {
		return delegation, false
	}

	delegation = types.MustUnmarshalDelegation(k.cdc, value)

	return delegation, true
}

// IterateAllDelegations iterates through all of the delegations.
func (k Keeper) IterateAllDelegations(ctx sdk.Context, zone *types.Zone, cb func(delegation types.Delegation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, append(types.KeyPrefixDelegation, []byte(zone.ChainId)...))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		if cb(delegation) {
			break
		}
	}
}

// GetAllDelegations returns all delegations used during genesis dump.
func (k Keeper) GetAllDelegations(ctx sdk.Context, zone *types.Zone) (delegations []types.Delegation) {
	k.IterateAllDelegations(ctx, zone, func(delegation types.Delegation) bool {
		delegations = append(delegations, delegation)
		return false
	})

	return delegations
}

// IterateAllPerformanceDelegations iterates through all of the delegations.
func (k Keeper) IterateAllPerformanceDelegations(ctx sdk.Context, zone *types.Zone, cb func(delegation types.Delegation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, append(types.KeyPrefixPerformanceDelegation, []byte(zone.ChainId)...))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		if cb(delegation) {
			break
		}
	}
}

// GetAllDelegations returns all delegations used during genesis dump.
func (k Keeper) GetAllPerformanceDelegations(ctx sdk.Context, zone *types.Zone) (delegations []types.Delegation) {
	k.IterateAllPerformanceDelegations(ctx, zone, func(delegation types.Delegation) bool {
		delegations = append(delegations, delegation)
		return false
	})

	return delegations
}

// GetAllDelegations returns all delegations used during genesis dump.
func (k Keeper) GetAllDelegationsAsPointer(ctx sdk.Context, zone *types.Zone) (delegations []*types.Delegation) {
	k.IterateAllDelegations(ctx, zone, func(delegation types.Delegation) bool {
		delegations = append(delegations, &delegation)
		return false
	})

	return delegations
}

// GetAllDelegations returns all delegations used during genesis dump.
func (k Keeper) GetAllPerformanceDelegationsAsPointer(ctx sdk.Context, zone *types.Zone) (delegations []*types.Delegation) {
	k.IterateAllPerformanceDelegations(ctx, zone, func(delegation types.Delegation) bool {
		delegations = append(delegations, &delegation)
		return false
	})

	return delegations
}

// GetDelegatorDelegations returns a given amount of all the delegations from a
// delegator.
func (k Keeper) GetDelegatorDelegations(ctx sdk.Context, zone *types.Zone, delegator sdk.AccAddress) (delegations []types.Delegation) {
	k.IterateDelegatorDelegations(ctx, zone, delegator, func(delegation types.Delegation) bool {
		delegations = append(delegations, delegation)
		return false
	})

	return delegations
}

// SetDelegation sets a delegation.
func (k Keeper) SetDelegation(ctx sdk.Context, zone *types.Zone, delegation types.Delegation) {
	delegatorAddress := delegation.GetDelegatorAddr()

	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalDelegation(k.cdc, delegation)
	store.Set(GetDelegationKey(zone, delegatorAddress, delegation.GetValidatorAddr()), b)
}

// SetPerformanceDelegation sets a delegation.
func (k Keeper) SetPerformanceDelegation(ctx sdk.Context, zone *types.Zone, delegation types.Delegation) {
	delegatorAddress := delegation.GetDelegatorAddr()

	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalDelegation(k.cdc, delegation)
	store.Set(GetPerformanceDelegationKey(zone, delegatorAddress, delegation.GetValidatorAddr()), b)
}

// RemoveDelegation removes a delegation
func (k Keeper) RemoveDelegation(ctx sdk.Context, zone *types.Zone, delegation types.Delegation) error {
	delegatorAddress := delegation.GetDelegatorAddr()

	store := ctx.KVStore(k.storeKey)
	store.Delete(GetDelegationKey(zone, delegatorAddress, delegation.GetValidatorAddr()))
	return nil
}

// IterateDelegatorDelegations iterates through one delegator's delegations.
func (k Keeper) IterateDelegatorDelegations(ctx sdk.Context, zone *types.Zone, delegator sdk.AccAddress, cb func(delegation types.Delegation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := GetDelegationsKey(zone, delegator)
	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		if cb(delegation) {
			break
		}
	}
}

func (k *Keeper) PrepareDelegationMessagesForCoins(_ sdk.Context, zone *types.Zone, allocations map[string]sdkmath.Int) []sdk.Msg {
	var msgs []sdk.Msg
	for _, valoper := range utils.Keys(allocations) {
		if !allocations[valoper].IsZero() {
			msgs = append(msgs, &stakingTypes.MsgDelegate{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: valoper, Amount: sdk.NewCoin(zone.BaseDenom, allocations[valoper])})
		}
	}
	return msgs
}

func (k *Keeper) PrepareDelegationMessagesForShares(_ sdk.Context, zone *types.Zone, coins sdk.Coins) []sdk.Msg {
	var msgs []sdk.Msg
	for _, coin := range coins.Sort() {
		if !coin.IsZero() {
			msgs = append(msgs, &lsmstakingTypes.MsgRedeemTokensforShares{DelegatorAddress: zone.DelegationAddress.Address, Amount: coin})
		}
	}
	return msgs
}

func (k Keeper) DeterminePlanForDelegation(ctx sdk.Context, zone *types.Zone, amount sdk.Coins) map[string]sdkmath.Int {
	currentAllocations, currentSum := k.GetDelegationMap(ctx, zone)
	targetAllocations := zone.GetAggregateIntentOrDefault()
	allocations := DetermineAllocationsForDelegation(currentAllocations, currentSum, targetAllocations, amount)
	return allocations
}

// CalculateDeltas determines, for the current delegations, in delta between actual allocations and the target intent.
func CalculateDeltas(currentAllocations map[string]sdkmath.Int, currentSum sdkmath.Int, targetAllocations types.ValidatorIntents) types.ValidatorIntents {
	deltas := make(types.ValidatorIntents, 0)

	targetValopers := func(in types.ValidatorIntents) []string {
		out := make([]string, 0)
		for _, i := range in.Sort() {
			out = append(out, i.ValoperAddress)
		}
		return out
	}(targetAllocations)

	keySet := utils.Unique(append(targetValopers, utils.Keys(currentAllocations)...))
	sort.Strings(keySet)
	// for target allocations, raise the intent weight by the total delegated value to get target amount
	for _, valoper := range keySet {
		current, ok := currentAllocations[valoper]
		if !ok {
			current = sdk.ZeroInt()
		}

		target, ok := targetAllocations.GetForValoper(valoper)
		if !ok {
			target = &types.ValidatorIntent{ValoperAddress: valoper, Weight: sdk.ZeroDec()}
		}
		targetAmount := target.Weight.MulInt(currentSum).TruncateInt()
		// diff between target and current allocations
		// positive == below target, negative == above target
		delta := targetAmount.Sub(current)
		deltas = append(deltas, &types.ValidatorIntent{Weight: sdk.NewDecFromInt(delta), ValoperAddress: valoper})
	}

	// sort keys by relative value of delta
	sort.SliceStable(deltas, func(i, j int) bool {
		return deltas[i].ValoperAddress > deltas[j].ValoperAddress
	})

	// sort keys by relative value of delta
	sort.SliceStable(deltas, func(i, j int) bool {
		return deltas[i].Weight.GT(deltas[j].Weight)
	})

	return deltas
}

// minDeltas returns the lowest value in a slice of Deltas.
func minDeltas(deltas types.ValidatorIntents) sdkmath.Int {
	minValue := sdk.NewInt(math.MaxInt64)
	for _, intent := range deltas {
		if minValue.GT(intent.Weight.TruncateInt()) {
			minValue = intent.Weight.TruncateInt()
		}
	}

	return minValue
}

func DetermineAllocationsForDelegation(currentAllocations map[string]sdkmath.Int, currentSum sdkmath.Int, targetAllocations types.ValidatorIntents, amount sdk.Coins) map[string]sdkmath.Int {
	input := amount[0].Amount
	deltas := CalculateDeltas(currentAllocations, currentSum, targetAllocations)
	minValue := minDeltas(deltas)
	sum := sdk.ZeroInt()

	// // sort keys by relative value of delta
	// sort.SliceStable(deltas, func(i, j int) bool {
	// 	return deltas[i].ValoperAddress > deltas[j].ValoperAddress
	// })

	// // sort keys by relative value of delta
	// sort.SliceStable(deltas, func(i, j int) bool {
	// 	return deltas[i].Weight.GT(deltas[j].Weight)
	// })

	// raise all deltas such that the minimum value is zero.
	for idx := range deltas {
		deltas[idx].Weight = deltas[idx].Weight.Add(sdk.NewDecFromInt(minValue.Abs()))
		sum = sum.Add(deltas[idx].Weight.TruncateInt())
	}

	// unequalSplit is the portion of input that should be distributed in attempt to make targets == 0
	unequalSplit := sdk.MinInt(sum, input)

	if !unequalSplit.IsZero() {
		for idx := range deltas {
			deltas[idx].Weight = deltas[idx].Weight.QuoInt(sum).MulInt(unequalSplit)
		}
	}

	// equalSplit is the portion of input that should be distributed equally across all validators, once targets are zero.
	equalSplit := sdk.NewDecFromInt(input.Sub(unequalSplit))

	if !equalSplit.IsZero() {
		each := equalSplit.Quo(sdk.NewDec(int64(len(deltas))))
		for idx := range deltas {
			deltas[idx].Weight = deltas[idx].Weight.Add(each)
		}
	}

	// dust is the portion of the input that was truncated in previous calculations; add this to the first validator in the list,
	// once sorted alphabetically. This will always be a small amount, and will count toward the delta calculations on the next run.

	outSum := sdk.ZeroInt()
	outWeights := make(map[string]sdkmath.Int)
	for _, delta := range deltas {
		outWeights[delta.ValoperAddress] = delta.Weight.TruncateInt()
		outSum = outSum.Add(delta.Weight.TruncateInt())
	}
	dust := input.Sub(outSum)
	outWeights[deltas[0].ValoperAddress] = outWeights[deltas[0].ValoperAddress].Add(dust)

	return outWeights
}

func (k *Keeper) WithdrawDelegationRewardsForResponse(ctx sdk.Context, zone *types.Zone, delegator string, response []byte) error {
	var msgs []sdk.Msg

	delegatorRewards := distrTypes.QueryDelegationTotalRewardsResponse{}
	err := k.cdc.Unmarshal(response, &delegatorRewards)
	if err != nil {
		return err
	}

	if zone.DelegationAddress.Address != delegator {
		return errors.New("failed attempting to withdraw rewards from non-delegation account")
	}

	for _, del := range delegatorRewards.Rewards {
		if !del.Reward.IsZero() && !del.Reward.Empty() {
			k.Logger(ctx).Info("Withdraw rewards", "delegator", delegator, "validator", del.ValidatorAddress, "amount", del.Reward)

			msgs = append(msgs, &distrTypes.MsgWithdrawDelegatorReward{DelegatorAddress: delegator, ValidatorAddress: del.ValidatorAddress})
		}
	}

	if len(msgs) == 0 {
		// always setZone here because calling method update waitgroup.
		k.SetZone(ctx, zone)
		return nil
	}
	// increment withdrawal waitgroup for every withdrawal msg sent
	// this allows us to track individual msg responses and ensure all
	// responses have been received and handled...
	// HandleWithdrawRewards contains the opposing decrement.
	zone.WithdrawalWaitgroup += uint32(len(msgs))
	k.SetZone(ctx, zone)
	k.Logger(ctx).Info("Received WithdrawDelegationRewardsForResponse acknowledgement", "wg", zone.WithdrawalWaitgroup, "address", delegator)

	return k.SubmitTx(ctx, msgs, zone.DelegationAddress, "", zone.MessagesPerTx)
}

func (k *Keeper) GetDelegationMap(ctx sdk.Context, zone *types.Zone) (map[string]sdkmath.Int, sdkmath.Int) {
	out := make(map[string]sdkmath.Int)
	sum := sdk.ZeroInt()

	k.IterateAllDelegations(ctx, zone, func(delegation types.Delegation) bool {
		existing, found := out[delegation.ValidatorAddress]
		if !found {
			out[delegation.ValidatorAddress] = delegation.Amount.Amount
		} else {
			out[delegation.ValidatorAddress] = existing.Add(delegation.Amount.Amount)
		}
		sum = sum.Add(delegation.Amount.Amount)
		return false
	})

	return out, sum
}

func (k *Keeper) MakePerformanceDelegation(ctx sdk.Context, zone *types.Zone, validator string) error {
	// create delegation record in MsgDelegate acknowledgement callback
	if zone.PerformanceAddress != nil {
		k.SetPerformanceDelegation(ctx, zone, types.NewDelegation(zone.PerformanceAddress.Address, validator, sdk.NewInt64Coin(zone.BaseDenom, 0))) // intentionally zero; we add a record here to stop race conditions
		msg := stakingTypes.MsgDelegate{DelegatorAddress: zone.PerformanceAddress.Address, ValidatorAddress: validator, Amount: sdk.NewInt64Coin(zone.BaseDenom, 10000)}
		return k.SubmitTx(ctx, []sdk.Msg{&msg}, zone.PerformanceAddress, "perf/"+validator, zone.MessagesPerTx)
	}
	return nil
}

func (k *Keeper) FlushOutstandingDelegations(ctx sdk.Context, zone *types.Zone) error {
	var err error
	k.IterateReceipts(ctx, func(_ int64, receiptInfo types.Receipt) (stop bool) {
		if receiptInfo.ChainId == zone.ChainId && receiptInfo.Completed == nil {
			sendMsg := banktypes.MsgSend{
				FromAddress: "",
				ToAddress:   "",
				Amount:      receiptInfo.Amount,
			}
			err = k.handleSendToDelegate(ctx, zone, &sendMsg, receiptInfo.Txhash)
			if err != nil {
				k.Logger(ctx).Error("error in processing pending delegations", "chain", zone.ChainId, "receipt", receiptInfo.Txhash, "error", err)
				return true
			}
		}
		return false
	})
	return err
}

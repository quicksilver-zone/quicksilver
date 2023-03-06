package keeper

import (
	"errors"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	distrTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	lsmstakingTypes "github.com/iqlusioninc/liquidity-staking-module/x/staking/types"

	"github.com/ingenuity-build/quicksilver/utils"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// GetDelegation returns a specific delegation.
func (k *Keeper) GetDelegation(ctx sdk.Context, zone *types.Zone, delegatorAddress string, validatorAddress string) (delegation types.Delegation, found bool) {
	store := ctx.KVStore(k.storeKey)

	_, delAddr, _ := bech32.DecodeAndConvert(delegatorAddress)
	_, valAddr, _ := bech32.DecodeAndConvert(validatorAddress)

	key := types.GetDelegationKey(zone, delAddr, valAddr)

	value := store.Get(key)
	if value == nil {
		return delegation, false
	}

	delegation = types.MustUnmarshalDelegation(k.cdc, value)

	return delegation, true
}

// GetDelegation returns a specific delegation.
func (k *Keeper) GetPerformanceDelegation(ctx sdk.Context, zone *types.Zone, validatorAddress string) (delegation types.Delegation, found bool) {
	if zone.PerformanceAddress == nil {
		return types.Delegation{}, false
	}

	store := ctx.KVStore(k.storeKey)

	_, delAddr, _ := bech32.DecodeAndConvert(zone.PerformanceAddress.Address)
	_, valAddr, _ := bech32.DecodeAndConvert(validatorAddress)

	key := types.GetPerformanceDelegationKey(zone, delAddr, valAddr)

	value := store.Get(key)
	if value == nil {
		return delegation, false
	}

	delegation = types.MustUnmarshalDelegation(k.cdc, value)

	return delegation, true
}

// IterateAllDelegations iterates through all of the delegations.
func (k *Keeper) IterateAllDelegations(ctx sdk.Context, zone *types.Zone, cb func(delegation types.Delegation) (stop bool)) {
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
func (k *Keeper) GetAllDelegations(ctx sdk.Context, zone *types.Zone) (delegations []types.Delegation) {
	k.IterateAllDelegations(ctx, zone, func(delegation types.Delegation) bool {
		delegations = append(delegations, delegation)
		return false
	})

	return delegations
}

// IterateAllPerformanceDelegations iterates through all of the delegations.
func (k *Keeper) IterateAllPerformanceDelegations(ctx sdk.Context, zone *types.Zone, cb func(delegation types.Delegation) (stop bool)) {
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
func (k *Keeper) GetAllPerformanceDelegations(ctx sdk.Context, zone *types.Zone) (delegations []types.Delegation) {
	k.IterateAllPerformanceDelegations(ctx, zone, func(delegation types.Delegation) bool {
		delegations = append(delegations, delegation)
		return false
	})

	return delegations
}

// GetAllDelegations returns all delegations used during genesis dump.
func (k *Keeper) GetAllDelegationsAsPointer(ctx sdk.Context, zone *types.Zone) (delegations []*types.Delegation) {
	k.IterateAllDelegations(ctx, zone, func(delegation types.Delegation) bool {
		delegations = append(delegations, &delegation)
		return false
	})

	return delegations
}

// GetAllDelegations returns all delegations used during genesis dump.
func (k *Keeper) GetAllPerformanceDelegationsAsPointer(ctx sdk.Context, zone *types.Zone) (delegations []*types.Delegation) {
	k.IterateAllPerformanceDelegations(ctx, zone, func(delegation types.Delegation) bool {
		delegations = append(delegations, &delegation)
		return false
	})

	return delegations
}

// GetDelegatorDelegations returns a given amount of all the delegations from a
// delegator.
func (k *Keeper) GetDelegatorDelegations(ctx sdk.Context, zone *types.Zone, delegator sdk.AccAddress) (delegations []types.Delegation) {
	k.IterateDelegatorDelegations(ctx, zone, delegator, func(delegation types.Delegation) bool {
		delegations = append(delegations, delegation)
		return false
	})

	return delegations
}

// SetDelegation sets a delegation.
func (k *Keeper) SetDelegation(ctx sdk.Context, zone *types.Zone, delegation types.Delegation) {
	delegatorAddress := delegation.GetDelegatorAddr()

	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalDelegation(k.cdc, delegation)
	store.Set(types.GetDelegationKey(zone, delegatorAddress, delegation.GetValidatorAddr()), b)
}

// SetPerformanceDelegation sets a delegation.
func (k *Keeper) SetPerformanceDelegation(ctx sdk.Context, zone *types.Zone, delegation types.Delegation) {
	delegatorAddress := delegation.GetDelegatorAddr()

	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalDelegation(k.cdc, delegation)
	store.Set(types.GetPerformanceDelegationKey(zone, delegatorAddress, delegation.GetValidatorAddr()), b)
}

// RemoveDelegation removes a delegation
func (k *Keeper) RemoveDelegation(ctx sdk.Context, zone *types.Zone, delegation types.Delegation) error {
	delegatorAddress := delegation.GetDelegatorAddr()

	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetDelegationKey(zone, delegatorAddress, delegation.GetValidatorAddr()))
	return nil
}

// IterateDelegatorDelegations iterates through one delegator's delegations.
func (k *Keeper) IterateDelegatorDelegations(ctx sdk.Context, zone *types.Zone, delegator sdk.AccAddress, cb func(delegation types.Delegation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := types.GetDelegationsKey(zone, delegator)
	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		if cb(delegation) {
			break
		}
	}
}

func (k *Keeper) PrepareDelegationMessagesForCoins(ctx sdk.Context, zone *types.Zone, allocations map[string]sdkmath.Int) []sdk.Msg {
	var msgs []sdk.Msg
	for _, valoper := range utils.Keys(allocations) {
		if !allocations[valoper].IsZero() {
			msgs = append(msgs, &stakingTypes.MsgDelegate{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: valoper, Amount: sdk.NewCoin(zone.BaseDenom, allocations[valoper])})
		}
	}
	return msgs
}

func (k *Keeper) PrepareDelegationMessagesForShares(ctx sdk.Context, zone *types.Zone, coins sdk.Coins) []sdk.Msg {
	var msgs []sdk.Msg
	for _, coin := range coins.Sort() {
		if !coin.IsZero() {
			msgs = append(msgs, &lsmstakingTypes.MsgRedeemTokensforShares{DelegatorAddress: zone.DelegationAddress.Address, Amount: coin})
		}
	}
	return msgs
}

func (k *Keeper) DeterminePlanForDelegation(ctx sdk.Context, zone *types.Zone, amount sdk.Coins) map[string]sdkmath.Int {
	currentAllocations, currentSum, _ := k.GetDelegationMap(ctx, zone)
	targetAllocations := zone.GetAggregateIntentOrDefault()
	allocations := types.DetermineAllocationsForDelegation(currentAllocations, currentSum, targetAllocations, amount)
	return allocations
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

	return k.SubmitTx(ctx, msgs, zone.DelegationAddress, "")
}

func (k *Keeper) GetDelegationMap(ctx sdk.Context, zone *types.Zone) (map[string]sdkmath.Int, sdkmath.Int, map[string]bool) {
	out := make(map[string]sdkmath.Int)
	locked := make(map[string]bool)
	sum := sdk.ZeroInt()

	k.IterateAllDelegations(ctx, zone, func(delegation types.Delegation) bool {
		existing, found := out[delegation.ValidatorAddress]
		if !found {
			out[delegation.ValidatorAddress] = delegation.Amount.Amount
			locked[delegation.ValidatorAddress] = delegation.RedelegationEnd != 0 && delegation.RedelegationEnd >= ctx.BlockTime().Unix()
		} else {
			out[delegation.ValidatorAddress] = existing.Add(delegation.Amount.Amount)
			locked[delegation.ValidatorAddress] = locked[delegation.ValidatorAddress] || (delegation.RedelegationEnd != 0 && delegation.RedelegationEnd >= ctx.BlockTime().Unix())
		}
		sum = sum.Add(delegation.Amount.Amount)
		return false
	})

	return out, sum, locked
}

func (k *Keeper) MakePerformanceDelegation(ctx sdk.Context, zone *types.Zone, validator string) error {
	// create delegation record in MsgDelegate acknowledgement callback
	if zone.PerformanceAddress != nil {
		k.SetPerformanceDelegation(ctx, zone, types.NewDelegation(zone.PerformanceAddress.Address, validator, sdk.NewInt64Coin(zone.BaseDenom, 0))) // intentionally zero; we add a record here to stop race conditions
		msg := stakingTypes.MsgDelegate{DelegatorAddress: zone.PerformanceAddress.Address, ValidatorAddress: validator, Amount: sdk.NewInt64Coin(zone.BaseDenom, 10000)}
		return k.SubmitTx(ctx, []sdk.Msg{&msg}, zone.PerformanceAddress, "perf/"+validator)
	}
	return nil
}

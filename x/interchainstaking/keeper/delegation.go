package keeper

import (
	"errors"
	"fmt"

	lsmstakingTypes "github.com/iqlusioninc/liquidity-staking-module/x/staking/types"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/quicksilver-zone/quicksilver/utils"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

// GetDelegation returns a specific delegation.
func (k *Keeper) GetDelegation(ctx sdk.Context, chainID string, delegatorAddress, validatorAddress string) (delegation types.Delegation, found bool) {
	store := ctx.KVStore(k.storeKey)

	_, delAddr, _ := bech32.DecodeAndConvert(delegatorAddress)
	_, valAddr, _ := bech32.DecodeAndConvert(validatorAddress)

	key := types.GetDelegationKey(chainID, delAddr, valAddr)

	value := store.Get(key)
	if value == nil {
		return delegation, false
	}

	delegation = types.MustUnmarshalDelegation(k.cdc, value)

	return delegation, true
}

// GetPerformanceDelegation returns a specific delegation.
func (k *Keeper) GetPerformanceDelegation(ctx sdk.Context, chainID string, performanceAddress *types.ICAAccount, validatorAddress string) (delegation types.Delegation, found bool) {
	if performanceAddress == nil {
		return types.Delegation{}, false
	}

	store := ctx.KVStore(k.storeKey)

	_, delAddr, _ := bech32.DecodeAndConvert(performanceAddress.Address)
	_, valAddr, _ := bech32.DecodeAndConvert(validatorAddress)

	key := types.GetPerformanceDelegationKey(chainID, delAddr, valAddr)

	value := store.Get(key)
	if value == nil {
		return delegation, false
	}

	delegation = types.MustUnmarshalDelegation(k.cdc, value)

	return delegation, true
}

// IterateAllDelegations iterates through all of the delegations.
func (k *Keeper) IterateAllDelegations(ctx sdk.Context, chainID string, cb func(delegation types.Delegation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, append(types.KeyPrefixDelegation, []byte(chainID)...))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		if cb(delegation) {
			break
		}
	}
}

// GetAllDelegations returns all delegations used during genesis dump.
func (k *Keeper) GetAllDelegations(ctx sdk.Context, chainID string) (delegations []types.Delegation) {
	k.IterateAllDelegations(ctx, chainID, func(delegation types.Delegation) bool {
		delegations = append(delegations, delegation)
		return false
	})

	return delegations
}

// IterateAllPerformanceDelegations iterates through all of the delegations.
func (k *Keeper) IterateAllPerformanceDelegations(ctx sdk.Context, chainID string, cb func(delegation types.Delegation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, append(types.KeyPrefixPerformanceDelegation, []byte(chainID)...))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		if cb(delegation) {
			break
		}
	}
}

// GetAllDelegations returns all delegations used during genesis dump.
func (k *Keeper) GetAllPerformanceDelegations(ctx sdk.Context, chainID string) (delegations []types.Delegation) {
	k.IterateAllPerformanceDelegations(ctx, chainID, func(delegation types.Delegation) bool {
		delegations = append(delegations, delegation)
		return false
	})

	return delegations
}

// GetAllDelegations returns all delegations used during genesis dump.
func (k *Keeper) GetAllDelegationsAsPointer(ctx sdk.Context, chainID string) (delegations []*types.Delegation) {
	k.IterateAllDelegations(ctx, chainID, func(delegation types.Delegation) bool {
		delegations = append(delegations, &delegation)
		return false
	})

	return delegations
}

// GetAllDelegations returns all delegations used during genesis dump.
func (k *Keeper) GetAllPerformanceDelegationsAsPointer(ctx sdk.Context, chainID string) (delegations []*types.Delegation) {
	k.IterateAllPerformanceDelegations(ctx, chainID, func(delegation types.Delegation) bool {
		delegations = append(delegations, &delegation)
		return false
	})

	return delegations
}

// GetDelegatorDelegations returns a given amount of all the delegations from a
// delegator.
func (k *Keeper) GetDelegatorDelegations(ctx sdk.Context, chainID string, delegator sdk.AccAddress) (delegations []types.Delegation) {
	k.IterateDelegatorDelegations(ctx, chainID, delegator, func(delegation types.Delegation) bool {
		delegations = append(delegations, delegation)
		return false
	})

	return delegations
}

// SetDelegation sets a delegation.
func (k *Keeper) SetDelegation(ctx sdk.Context, chainID string, delegation types.Delegation) {
	delegatorAddress := delegation.GetDelegatorAddr()

	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalDelegation(k.cdc, delegation)
	store.Set(types.GetDelegationKey(chainID, delegatorAddress, delegation.GetValidatorAddr()), b)
}

// SetPerformanceDelegation sets a delegation.
func (k *Keeper) SetPerformanceDelegation(ctx sdk.Context, chainID string, delegation types.Delegation) {
	delegatorAddress := delegation.GetDelegatorAddr()

	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalDelegation(k.cdc, delegation)
	store.Set(types.GetPerformanceDelegationKey(chainID, delegatorAddress, delegation.GetValidatorAddr()), b)
}

// RemoveDelegation removes a delegation.
func (k *Keeper) RemoveDelegation(ctx sdk.Context, chainID string, delegation types.Delegation) error {
	delegatorAddress := delegation.GetDelegatorAddr()

	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetDelegationKey(chainID, delegatorAddress, delegation.GetValidatorAddr()))
	return nil
}

// RemovePerformanceDelegation removes a performance delegation.
func (k *Keeper) RemovePerformanceDelegation(ctx sdk.Context, chainID string, delegation types.Delegation) error {
	delegatorAddress := delegation.GetDelegatorAddr()

	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetPerformanceDelegationKey(chainID, delegatorAddress, delegation.GetValidatorAddr()))
	return nil
}

// IterateDelegatorDelegations iterates through one delegator's delegations.
func (k *Keeper) IterateDelegatorDelegations(ctx sdk.Context, chainID string, delegator sdk.AccAddress, cb func(delegation types.Delegation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := types.GetDelegationsKey(chainID, delegator)
	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		if cb(delegation) {
			break
		}
	}
}

func (*Keeper) PrepareDelegationMessagesForCoins(zone *types.Zone, allocations map[string]sdkmath.Int) []sdk.Msg {
	var msgs []sdk.Msg
	for _, valoper := range utils.Keys(allocations) {
		if !allocations[valoper].IsZero() {
			msgs = append(msgs, &stakingTypes.MsgDelegate{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: valoper, Amount: sdk.NewCoin(zone.BaseDenom, allocations[valoper])})
		}
	}
	return msgs
}

func (*Keeper) PrepareDelegationMessagesForShares(zone *types.Zone, coins sdk.Coins) []sdk.Msg {
	var msgs []sdk.Msg
	for _, coin := range coins.Sort() {
		if !coin.IsZero() {
			msgs = append(msgs, &lsmstakingTypes.MsgRedeemTokensforShares{DelegatorAddress: zone.DelegationAddress.Address, Amount: coin})
		}
	}
	return msgs
}

func (k *Keeper) DeterminePlanForDelegation(ctx sdk.Context, zone *types.Zone, amount sdk.Coins) (map[string]sdkmath.Int, error) {
	currentAllocations, currentSum, _, _ := k.GetDelegationMap(ctx, zone.ChainId)
	targetAllocations, err := k.GetAggregateIntentOrDefault(ctx, zone)
	if err != nil {
		return nil, err
	}
	return types.DetermineAllocationsForDelegation(currentAllocations, currentSum, targetAllocations, amount)
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

func (k *Keeper) GetDelegationMap(ctx sdk.Context, chainID string) (out map[string]sdkmath.Int, sum sdkmath.Int, locked map[string]bool, lockedSum sdkmath.Int) {
	out = make(map[string]sdkmath.Int)
	locked = make(map[string]bool)
	sum = sdk.ZeroInt()
	lockedSum = sdk.ZeroInt()

	k.IterateAllDelegations(ctx, chainID, func(delegation types.Delegation) bool {
		out[delegation.ValidatorAddress] = delegation.Amount.Amount
		if delegation.RedelegationEnd >= ctx.BlockTime().Unix() {
			locked[delegation.ValidatorAddress] = true
			lockedSum = lockedSum.Add(delegation.Amount.Amount)
		}
		sum = sum.Add(delegation.Amount.Amount)
		return false
	})

	return out, sum, locked, lockedSum
}

func (k *Keeper) MakePerformanceDelegation(ctx sdk.Context, zone *types.Zone, validator string) error {
	// create delegation record in MsgDelegate acknowledgement callback
	if zone.PerformanceAddress != nil {
		k.SetPerformanceDelegation(ctx, zone.ChainId, types.NewDelegation(zone.PerformanceAddress.Address, validator, sdk.NewInt64Coin(zone.BaseDenom, 0))) // intentionally zero; we add a record here to stop race conditions
		msg := stakingTypes.MsgDelegate{DelegatorAddress: zone.PerformanceAddress.Address, ValidatorAddress: validator, Amount: sdk.NewInt64Coin(zone.BaseDenom, 10000)}
		return k.SubmitTx(ctx, []sdk.Msg{&msg}, zone.PerformanceAddress, fmt.Sprintf("%s/%s", types.MsgTypePerformance, validator), zone.MessagesPerTx)
	}
	return nil
}

func (k *Keeper) FlushOutstandingDelegations(ctx sdk.Context, zone *types.Zone, delAddrBalance sdk.Coin) error {
	var pendingAmount sdk.Coins
	exclusionTime := ctx.BlockTime().AddDate(0, 0, -1)
	k.IterateZoneReceipts(ctx, zone.ChainId, func(_ int64, receiptInfo types.Receipt) (stop bool) {
		if (receiptInfo.FirstSeen.After(exclusionTime) || receiptInfo.FirstSeen.Equal(exclusionTime)) && receiptInfo.Completed == nil {
			pendingAmount = pendingAmount.Add(receiptInfo.Amount...)
		}
		return false
	})

	coinsToFlush, hasNeg := sdk.NewCoins(delAddrBalance).SafeSub(pendingAmount...)
	if hasNeg || coinsToFlush.IsZero() {
		k.Logger(ctx).Debug("delegate account balance negative, setting outdated reciepts")
		k.SetReceiptsCompleted(ctx, zone.ChainId, exclusionTime, ctx.BlockTime())
		return nil
	}

	// set the zone amount to the coins to be flushed.
	zone.DelegationAddress.Balance = coinsToFlush
	k.Logger(ctx).Info("flush delegations ", "total", coinsToFlush)
	k.SetZone(ctx, zone)

	sendMsg := banktypes.MsgSend{
		FromAddress: "",
		ToAddress:   "",
		Amount:      coinsToFlush,
	}
	return k.handleSendToDelegate(ctx, zone, &sendMsg, fmt.Sprintf("batch/%d", exclusionTime.Unix()))
}

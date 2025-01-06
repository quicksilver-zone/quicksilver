package keeper

import (
	"errors"
	"fmt"
	"math"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/quicksilver-zone/quicksilver/utils"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	lsmstakingtypes "github.com/quicksilver-zone/quicksilver/x/lsmtypes"
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

	iterator := sdk.KVStorePrefixIterator(store, append(types.KeyPrefixDelegation, chainID...))
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

	iterator := sdk.KVStorePrefixIterator(store, append(types.KeyPrefixPerformanceDelegation, chainID...))
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

func (*Keeper) PrepareDelegationMessagesForCoins(zone *types.Zone, allocations map[string]sdkmath.Int, isFlush bool) []sdk.Msg {
	var msgs []sdk.Msg
	for _, valoper := range utils.Keys(allocations) {
		if allocations[valoper].IsPositive() {
			if allocations[valoper].GTE(zone.DustThreshold) || isFlush {
				msgs = append(msgs, &stakingtypes.MsgDelegate{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: valoper, Amount: sdk.NewCoin(zone.BaseDenom, allocations[valoper])})
			}
		}
	}
	return msgs
}

func (*Keeper) PrepareDelegationMessagesForShares(zone *types.Zone, coins sdk.Coins) []sdk.Msg {
	var msgs []sdk.Msg
	for _, coin := range coins.Sort() {
		if coin.IsPositive() {
			// no min amount here.
			msgs = append(msgs, &lsmstakingtypes.MsgRedeemTokensForShares{DelegatorAddress: zone.DelegationAddress.Address, Amount: coin})
		}
	}
	return msgs
}

func (k Keeper) DetermineMaximumValidatorAllocations(ctx sdk.Context, zone *types.Zone) map[string]sdkmath.Int {
	out := make(map[string]sdkmath.Int)
	caps, found := k.GetLsmCaps(ctx, zone.ChainId)
	if !found {
		// No cap found, permit the transaction
		return out
	}

	for _, val := range k.GetValidators(ctx, zone.ChainId) {
		// validator bond max
		maxBondShares := val.ValidatorBondShares.Mul(caps.ValidatorBondCap).Sub(val.LiquidShares)

		// validator pc max
		maxLiquidStakedShares := sdk.NewDecFromInt(val.VotingPower).Mul(caps.ValidatorCap).Sub(val.LiquidShares)
		out[val.ValoperAddress] = sdkmath.MaxInt(sdk.ZeroInt(), sdkmath.MinInt(maxBondShares.TruncateInt(), maxLiquidStakedShares.TruncateInt()))
	}

	return out
}

func (k *Keeper) DeterminePlanForDelegation(ctx sdk.Context, zone *types.Zone, amount sdk.Coins) (map[string]sdkmath.Int, error) {
	currentAllocations, currentSum, _, _ := k.GetDelegationMap(ctx, zone.ChainId)
	targetAllocations, err := k.GetAggregateIntentOrDefault(ctx, zone)
	if err != nil {
		return nil, err
	}
	maxCanAllocate := k.DetermineMaximumValidatorAllocations(ctx, zone)
	return types.DetermineAllocationsForDelegation(currentAllocations, currentSum, targetAllocations, amount, maxCanAllocate)
}

func (k *Keeper) WithdrawDelegationRewardsForResponse(ctx sdk.Context, zone *types.Zone, delegator string, response []byte) error {
	var msgs []sdk.Msg

	delegatorRewards := distrtypes.QueryDelegationTotalRewardsResponse{}
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

			msgs = append(msgs, &distrtypes.MsgWithdrawDelegatorReward{DelegatorAddress: delegator, ValidatorAddress: del.ValidatorAddress})
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
	if len(msgs) > math.MaxUint32 {
		return fmt.Errorf("number of messages exceeds uint32 range: %d", len(msgs))
	}
	if err = zone.IncrementWithdrawalWaitgroup(k.Logger(ctx), uint32(len(msgs)), "WithdrawDelegationRewardsForResponse"); err != nil { //nolint:gosec
		return err
	}
	k.SetZone(ctx, zone)
	k.Logger(ctx).Info("Received WithdrawDelegationRewardsForResponse acknowledgement", "wg", zone.GetWithdrawalWaitgroup(), "address", delegator)

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
		msg := stakingtypes.MsgDelegate{DelegatorAddress: zone.PerformanceAddress.Address, ValidatorAddress: validator, Amount: sdk.NewInt64Coin(zone.BaseDenom, 10000)}
		return k.SubmitTx(ctx, []sdk.Msg{&msg}, zone.PerformanceAddress, fmt.Sprintf("%s/%s", types.MsgTypePerformance, validator), zone.MessagesPerTx)
	}
	return nil
}

func (k *Keeper) FlushOutstandingDelegations(ctx sdk.Context, zone *types.Zone, delAddrBalance sdk.Coin) error {
	var pendingAmount sdk.Coins
	exclusionTime := ctx.BlockTime().AddDate(0, 0, -1)
	k.IterateZoneReceipts(ctx, zone.ChainId, func(_ int64, receiptInfo types.Receipt) (stop bool) {
		if (receiptInfo.FirstSeen.After(exclusionTime) || receiptInfo.FirstSeen.Equal(exclusionTime)) && receiptInfo.Completed == nil && receiptInfo.Amount[0].Denom == delAddrBalance.Denom {
			k.Logger(ctx).Info("adding to pending amount", "pending receipt", receiptInfo)
			pendingAmount = pendingAmount.Add(receiptInfo.Amount...)
		}
		return false
	})

	pendingAmount = pendingAmount.Add(k.GetInflightUnbondingAmount(ctx, zone))

	coinsToFlush, hasNeg := sdk.NewCoins(delAddrBalance).SafeSub(pendingAmount...)
	if hasNeg || coinsToFlush.IsZero() {
		k.Logger(ctx).Info("delegate account balance negative, or nothing to flush, setting outdated receipts")
		k.SetReceiptsCompleted(ctx, zone.ChainId, exclusionTime, ctx.BlockTime(), delAddrBalance.Denom)
		if zone.GetWithdrawalWaitgroup() == 0 {
			// we won't be sending any messages when we exit here; so if WG==0, then trigger RR update
			k.Logger(ctx).Info("triggering redemption rate calc in lieu of delegation flush (non-positive coins)")
			if err := k.TriggerRedemptionRate(ctx, zone); err != nil {
				return err
			}
		}
		return nil
	}

	// set the zone amount to the coins to be flushed.
	k.Logger(ctx).Info("flush delegations ", "total", coinsToFlush)

	sendMsg := banktypes.MsgSend{
		FromAddress: "",
		ToAddress:   "",
		Amount:      coinsToFlush,
	}
	numMsgs, err := k.handleSendToDelegate(ctx, zone, &sendMsg, fmt.Sprintf("batch/%d", exclusionTime.Unix()))
	if err != nil {
		return err
	}
	if numMsgs > math.MaxUint32 {
		return fmt.Errorf("number of messages exceeds uint32 range: %d", numMsgs)
	}
	if err = zone.IncrementWithdrawalWaitgroup(k.Logger(ctx), uint32(numMsgs), "sending flush messages"); err != nil { //nolint:gosec
		return err
	}

	// if we didn't send any messages (thus no acks will happen), and WG==0, then trigger RR update
	if numMsgs == 0 && zone.GetWithdrawalWaitgroup() == 0 {
		k.Logger(ctx).Info("triggering redemption rate calc in lieu of delegation flush (no messages to send)")
		if err := k.TriggerRedemptionRate(ctx, zone); err != nil {
			return err
		}
	}

	k.SetZone(ctx, zone)
	return nil
}

package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	distrTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"

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

// GetAllDelegations returns all delegations used during genesis dump.
func (k Keeper) GetAllDelegationsAsPointer(ctx sdk.Context, zone *types.Zone) (delegations []*types.Delegation) {
	k.IterateAllDelegations(ctx, zone, func(delegation types.Delegation) bool {
		delegations = append(delegations, &delegation)
		return false
	})

	return delegations
}

// GetValidatorDelegations returns all delegations to a specific validator.
// Useful for querier.
func (k Keeper) GetValidatorDelegations(ctx sdk.Context, zone *types.Zone, valAddr sdk.ValAddress) (delegations []types.Delegation) { //nolint:interfacer
	k.IterateAllDelegations(ctx, zone, func(delegation types.Delegation) bool {
		if delegation.GetValidatorAddr().Equals(valAddr) {
			delegations = append(delegations, delegation)
		}
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

// Delegate determines how the balance of a DelegateAccount should be distributed across validators.
func (k *Keeper) Delegate(ctx sdk.Context, zone types.Zone, account *types.ICAAccount, allocations types.Allocations) error {
	var msgs []sdk.Msg

	for _, allocation := range allocations.Sorted() {
		for _, coin := range allocation.Amount {
			if coin.Denom == zone.BaseDenom {
				msgs = append(msgs, &stakingTypes.MsgDelegate{DelegatorAddress: account.GetAddress(), ValidatorAddress: allocation.Address, Amount: coin})
			} else {
				msgs = append(msgs, &stakingTypes.MsgRedeemTokensforShares{DelegatorAddress: account.GetAddress(), Amount: coin})
			}
		}
	}
	return k.SubmitTx(ctx, msgs, account, "")
}

func (k Keeper) DeterminePlanForDelegation(ctx sdk.Context, zone types.Zone, amount sdk.Coins, delegator string, txhash string) (types.Allocations, error) {
	bins := k.GetDelegationBinsMap(ctx, &zone)

	sendPlan := types.Allocations{}

	for _, coin := range amount {
		var delPlan types.Allocations
		var err error
		if coin.Denom == zone.BaseDenom {
			var valPlan types.ValidatorIntents
			plan, found := k.GetIntent(ctx, zone, delegator, false)
			if !found || len(plan.Intents) == 0 {
				valPlan = zone.GetAggregateIntentOrDefault()
				delPlan, err = types.DelegationPlanFromGlobalIntent(k.GetDelegatedAmount(ctx, &zone), bins, coin, valPlan)
				if err != nil {
					return types.Allocations{}, err
				}
			} else {
				valPlan = plan.ToValidatorIntents()
				delPlan = types.DelegationPlanFromUserIntent(zone, coin, valPlan) // TODO: does it make sense to do this? why don't we just use the global intent?
				if err != nil {
					return types.Allocations{}, err
				}
			}

		} else {
			delPlan = types.DelegationPlanFromCoins(zone, coin)
		}

		for _, allocation := range delPlan.Sorted() {
			var delegatorAddress string
			for _, coin := range allocation.Amount {
				delegatorAddress, bins = bins.FindAccountForDelegation(allocation.Address, sdk.NewCoin(zone.BaseDenom, coin.Amount))

				delegationPlan := types.NewDelegationPlan(delegatorAddress, allocation.Address, sdk.NewCoins(coin))
				sendPlan = sendPlan.Allocate(delegatorAddress, sdk.NewCoins(coin))
				k.SetDelegationPlan(ctx, &zone, txhash, delegationPlan)
			}
		}
	}

	return sendPlan, nil
}

func (k *Keeper) WithdrawDelegationRewardsForResponse(ctx sdk.Context, zone *types.Zone, delegator string, response []byte) error {
	var msgs []sdk.Msg

	delegatorRewards := distrTypes.QueryDelegationTotalRewardsResponse{}
	err := k.cdc.Unmarshal(response, &delegatorRewards)
	if err != nil {
		return err
	}
	account, err := zone.GetDelegationAccountByAddress(delegator)
	if err != nil {
		return err
	}

	var delAddr sdk.AccAddress
	_, delAddr, _ = bech32.DecodeAndConvert(delegator)

	// send withdrawal msg for each delegation (delegator:validator pairs)
	k.IterateDelegatorDelegations(ctx, zone, delAddr, func(delegation types.Delegation) bool {
		amount := rewardsForDelegation(delegatorRewards, delegation.ValidatorAddress)
		k.Logger(ctx).Info("Withdraw rewards", "delegator", delegation.DelegationAddress, "validator", delegation.ValidatorAddress, "amount", amount)
		if !amount.IsZero() || !amount.Empty() {
			msgs = append(msgs, &distrTypes.MsgWithdrawDelegatorReward{DelegatorAddress: delegation.GetDelegationAddress(), ValidatorAddress: delegation.GetValidatorAddress()})
		}
		return false
	})

	if len(msgs) == 0 {
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

	return k.SubmitTx(ctx, msgs, account, "")
}

func rewardsForDelegation(delegatorRewards distrTypes.QueryDelegationTotalRewardsResponse, validator string) sdk.DecCoins {
	for _, reward := range delegatorRewards.Rewards {
		if reward.ValidatorAddress == validator {
			return reward.Reward
		}
	}
	return sdk.NewDecCoins()
}

func (k *Keeper) GetDelegationBinsMap(ctx sdk.Context, zone *types.Zone) types.Allocations {
	out := types.Allocations{}
	for _, da := range zone.DelegationAddresses {
		out = out.Allocate(da.Address, sdk.Coins{sdk.Coin{Denom: zone.BaseDenom, Amount: sdk.ZeroInt()}})
	}

	k.IterateAllDelegations(ctx, zone, func(delegation types.Delegation) bool {
		out = out.Allocate(delegation.DelegationAddress, sdk.Coins{sdk.Coin{Denom: delegation.ValidatorAddress, Amount: delegation.Amount.Amount}})
		return false
	})

	return out.Sorted()
}

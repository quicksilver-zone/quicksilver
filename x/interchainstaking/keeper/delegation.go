package keeper

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	distrTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

// gets the key for delegator bond with validator
// VALUE: staking/Delegation
func GetDelegationKey(zone *types.RegisteredZone, delAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	return append(GetDelegationsKey(zone, delAddr), valAddr.Bytes()...)
}

// gets the prefix for a delegator for all validators
func GetDelegationsKey(zone *types.RegisteredZone, delAddr sdk.AccAddress) []byte {
	return append(append(types.KeyPrefixDelegation, []byte(zone.ChainId)...), delAddr.Bytes()...)
}

// GetDelegation returns a specific delegation.
func (k Keeper) GetDelegation(ctx sdk.Context, zone *types.RegisteredZone, delegatorAddress string, validatorAddress string) (delegation types.Delegation, found bool) {
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
func (k Keeper) IterateAllDelegations(ctx sdk.Context, zone *types.RegisteredZone, cb func(delegation types.Delegation) (stop bool)) {
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
func (k Keeper) GetAllDelegations(ctx sdk.Context, zone *types.RegisteredZone) (delegations []types.Delegation) {
	k.IterateAllDelegations(ctx, zone, func(delegation types.Delegation) bool {
		delegations = append(delegations, delegation)
		return false
	})

	return delegations
}

// GetValidatorDelegations returns all delegations to a specific validator.
// Useful for querier.
func (k Keeper) GetValidatorDelegations(ctx sdk.Context, zone *types.RegisteredZone, valAddr sdk.ValAddress) (delegations []types.Delegation) { //nolint:interfacer
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, append(types.KeyPrefixDelegation, []byte(zone.ChainId)...))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		if delegation.GetValidatorAddr().Equals(valAddr) {
			delegations = append(delegations, delegation)
		}
	}

	return delegations
}

// GetDelegatorDelegations returns a given amount of all the delegations from a
// delegator.
func (k Keeper) GetDelegatorDelegations(ctx sdk.Context, zone *types.RegisteredZone, delegator sdk.AccAddress) (delegations []types.Delegation) {
	delegations = []types.Delegation{}
	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := GetDelegationsKey(zone, delegator)

	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		delegations = append(delegations, delegation)
	}

	return delegations
}

// SetDelegation sets a delegation.
func (k Keeper) SetDelegation(ctx sdk.Context, zone *types.RegisteredZone, delegation types.Delegation) {
	delegatorAddress := delegation.GetDelegatorAddr()

	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalDelegation(k.cdc, delegation)
	store.Set(GetDelegationKey(zone, delegatorAddress, delegation.GetValidatorAddr()), b)
}

// RemoveDelegation removes a delegation
func (k Keeper) RemoveDelegation(ctx sdk.Context, zone *types.RegisteredZone, delegation types.Delegation) error {
	delegatorAddress := delegation.GetDelegatorAddr()

	store := ctx.KVStore(k.storeKey)
	store.Delete(GetDelegationKey(zone, delegatorAddress, delegation.GetValidatorAddr()))
	return nil
}

// IterateDelegatorDelegations iterates through one delegator's delegations.
func (k Keeper) IterateDelegatorDelegations(ctx sdk.Context, zone *types.RegisteredZone, delegator sdk.AccAddress, cb func(delegation types.Delegation) (stop bool)) {
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
func (k *Keeper) Delegate(ctx sdk.Context, zone types.RegisteredZone, account *types.ICAAccount) error {
	var msgs []sdk.Msg

	balance := account.Balance

	// deterministically sort balance
	sort.Slice(balance, func(i, j int) bool { return balance[i].Denom > balance[j].Denom })

	for _, asset := range balance {
		if asset.Denom == zone.GetBaseDenom() {
			keys, validators, err := k.DetermineValidatorsForDelegation(ctx, zone, asset)
			// TODO: return multiple validators here; consider the size of the delegation too - are we going to increase balance 'too far'?
			// given that we pass in the account balance, we should be able to return a map of valoper:balance and send the requisite MsgDelegates.
			// this is less important for rewards, but far more important for deposits of native assets.
			if err != nil {
				k.Logger(ctx).Error("Unable to determine validators for delegation: %v", err)
				continue
			}
			for _, valoper_address := range keys {
				amount := validators[valoper_address]
				if !amount.Amount.IsZero() {
					k.Logger(ctx).Info("Sending a MsgDelegate!", "asset", amount, "valoper", valoper_address)
					msgs = append(msgs, &stakingTypes.MsgDelegate{DelegatorAddress: account.GetAddress(), ValidatorAddress: valoper_address, Amount: amount})
				}
			}
		} else {
			k.Logger(ctx).Info("Sending a MsgRedeemTokensforShares!", "asset", asset)

			// TODO: validate this against validators?
			// if validator is not active, then redelegate msg too?
			msgs = append(msgs, &stakingTypes.MsgRedeemTokensforShares{DelegatorAddress: account.GetAddress(), Amount: asset})
		}
	}
	return k.SubmitTx(ctx, msgs, account)
}

func (k Keeper) DetermineStateIntentDiff(ctx sdk.Context, zone types.RegisteredZone) map[string]sdk.Int {
	totalAggregateIntent := sdk.ZeroDec()
	currentState := make(map[string]sdk.Int)
	totalDelegations := sdk.ZeroInt()
	diff := make(map[string]sdk.Int)

	// sum total aggregate intent
	for _, val := range zone.AggregateIntent {
		totalAggregateIntent = totalAggregateIntent.Add(val.Weight)

	}

	validators := zone.GetValidatorsSorted()

	if totalAggregateIntent.IsZero() {
		// if totalAggregateIntent is zero (that is, we have no intent set - which can happen
		// if we have only ever have native tokens staked and nbody has signalled intent) give
		// every validator an equal intent artificially.

		// this can be removed when we cache intent.
		if zone.AggregateIntent == nil {
			zone.AggregateIntent = make(map[string]*types.ValidatorIntent)
		}

		for _, val := range validators {
			zone.AggregateIntent[val.ValoperAddress] = &types.ValidatorIntent{ValoperAddress: val.ValoperAddress, Weight: sdk.OneDec()}
			totalAggregateIntent = totalAggregateIntent.Add(sdk.OneDec())
		}
	}

	for _, i := range validators {
		stake := sdk.ZeroInt()
		_, valAddr, _ := bech32.DecodeAndConvert(i.ValoperAddress)
		for _, delegation := range k.GetValidatorDelegations(ctx, &zone, valAddr) {
			stake = stake.Add(delegation.Amount.Amount)
		}
		currentState[i.ValoperAddress] = stake
		totalDelegations = totalDelegations.Add(stake)
	}
	ratio := totalDelegations.ToDec().Quo(totalAggregateIntent) // will always be >= 1.0

	for _, i := range validators {
		current, found := currentState[i.ValoperAddress]
		if !found {
			// this probably can happen if we have intent for a validator not in the set
			// (although we _should_ have all validators, current and past in the set).
			panic("this shouldn't happen...")
		}
		desired, found := zone.AggregateIntent[i.ValoperAddress]
		if !found {
			desired = &types.ValidatorIntent{ValoperAddress: i.ValoperAddress, Weight: sdk.ZeroDec()} // this is okay! just means nobody wants this validator anymore!
		}
		thisDiff := desired.Weight.Mul(ratio).TruncateInt().Sub(current)
		if !thisDiff.Equal(sdk.ZeroInt()) {
			diff[i.ValoperAddress] = thisDiff
		}
	}
	return diff
}

func (k Keeper) DetermineValidatorsForDelegation(ctx sdk.Context, zone types.RegisteredZone, amount sdk.Coin) ([]string, map[string]sdk.Coin, error) {
	out := make(map[string]sdk.Coin)

	coinAmount := amount.Amount
	aggregateIntents := zone.GetAggregateIntent()

	if len(aggregateIntents) == 0 {
		aggregateIntents = defaultAggregateIntents(ctx, zone)
	}

	keys := make([]string, 0)
	for valoper, intent := range aggregateIntents {
		keys = append(keys, valoper)
		if !coinAmount.IsZero() {
			// while there is some balance left to distribute
			// calculate the int value of weight * amount to distribute.
			thisAmount := intent.Weight.MulInt(amount.Amount).TruncateInt()
			// set distrubtion amount
			out[valoper] = sdk.Coin{Denom: amount.Denom, Amount: thisAmount}
			// reduce outstanding pool
			coinAmount = coinAmount.Sub(thisAmount)
		}
	}

	sort.Strings(keys)
	v0 := keys[0]
	out[v0] = out[v0].AddAmount(coinAmount)

	k.Logger(ctx).Info("Validator weightings without diffs", "weights", out)

	// calculate diff between current state and intended state.
	//diffs := k.DetermineStateIntentDiff(ctx, zone)

	// apply diff to distrubtion of delegation.
	// out, remaining := zone.ApplyDiffsToDistribution(out, diffs)
	// if !remaining.IsZero() {
	// 	for _, valoper := range keys {
	// 		intent := aggregateIntents[valoper]
	// 		thisAmount := intent.Weight.MulInt(remaining).TruncateInt()
	// 		thisOutAmount, ok := out[valoper]
	// 		if !ok {
	// 			thisOutAmount = sdk.NewCoin(amount.Denom, sdk.ZeroInt())
	// 		}

	// 		out[valoper] = thisOutAmount.AddAmount(thisAmount)
	// 		remaining = remaining.Sub(thisAmount)
	// 	}

	// 	v0 := keys[0]
	// 	out[v0] = out[v0].AddAmount(remaining)
	// }

	//k.Logger(ctx).Info("Determined validators from aggregated intents +/- rebalance diffs", "amount", amount.Amount, "out", out)

	return keys, out, nil
}

func (k *Keeper) WithdrawDelegationRewardsForResponse(ctx sdk.Context, zone *types.RegisteredZone, delegator string, response []byte) error {
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
		amount := rewardsForDelegation(delegatorRewards, delegation.DelegationAddress, delegation.ValidatorAddress)
		k.Logger(ctx).Info("Withdraw rewards", "delegator", delegation.DelegationAddress, "validator", delegation.ValidatorAddress, "amount", amount)
		if !amount.IsZero() {
			msgs = append(msgs, &distrTypes.MsgWithdrawDelegatorReward{DelegatorAddress: delegation.GetDelegationAddress(), ValidatorAddress: delegation.GetValidatorAddress()})
		}
		return false
	})

	if len(msgs) == 0 {
		return nil
	}
	// add withdrawal waitgroup tally
	zone.WithdrawalWaitgroup += uint32(len(msgs))
	k.SetRegisteredZone(ctx, *zone)

	k.Logger(ctx).Info("Withdraw delegation messages", "msgs", msgs)

	return k.SubmitTx(ctx, msgs, account)
}

func rewardsForDelegation(delegatorRewards distrTypes.QueryDelegationTotalRewardsResponse, delegator string, validator string) sdk.DecCoins {
	for _, reward := range delegatorRewards.Rewards {
		if reward.ValidatorAddress == validator {
			return reward.Reward
		}
	}
	return sdk.NewDecCoins()
}

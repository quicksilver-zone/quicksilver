package keeper

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

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

func (k *Keeper) WithdrawDelegationRewardsForResponse(ctx sdk.Context, zone types.RegisteredZone, delegator string, response []byte) error {
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
	// send withdrawal msg for each delegation (delegator:validator pairs)
	for _, delegation := range zone.GetDelegationsForDelegator(delegator) {
		amount := rewardsForDelegation(delegatorRewards, delegation.DelegationAddress, delegation.ValidatorAddress)
		k.Logger(ctx).Info("Withdraw rewards", "delegator", delegation.DelegationAddress, "validator", delegation.ValidatorAddress, "amount", amount)
		msgs = append(msgs, &distrTypes.MsgWithdrawDelegatorReward{DelegatorAddress: delegation.GetDelegationAddress(), ValidatorAddress: delegation.GetValidatorAddress()})
	}
	if len(msgs) == 0 {
		return nil
	}
	// add withdrawal waitgroup tally
	zone.WithdrawalWaitgroup += uint32(len(msgs))
	k.SetRegisteredZone(ctx, zone)

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

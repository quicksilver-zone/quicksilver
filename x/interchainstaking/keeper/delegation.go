package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	queryKeeper "github.com/ingenuity-build/quicksilver/x/interchainquery/keeper"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func (k *Keeper) Delegate(ctx sdk.Context, zone types.RegisteredZone, account *types.ICAAccount) error {
	var msgs []sdk.Msg

	for _, asset := range account.Balance {
		if asset.Denom == zone.GetBaseDenom() {
			validators, err := k.DetermineValidatorsForDelegation(ctx, zone, asset)
			// TODO: return multiple validators here; consider the size of the delegation too - are we going to increase balance 'too far'?
			// given that we pass in the account balance, we should be able to return a map of valoper:balance and send the requisite MsgDelegates.
			// this is less important for rewards, but far more important for deposits of native assets.
			if err != nil {
				k.Logger(ctx).Error("Unable to determine validators for delegation: %v", err)
				continue
			}
			for valoper_address, amount := range validators {
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

func (k *Keeper) WithdrawDelegationRewards(ctx sdk.Context, zone types.RegisteredZone, account *types.ICAAccount) error {
	k.Logger(ctx).Debug("Withdrawing rewards for delegate account", "account", account.GetAddress(), "zone", zone.ChainId)
	var msgs []sdk.Msg
	delegatorRewardsDatapoint, err := k.ICQKeeper.GetDatapointForId(ctx, queryKeeper.GenerateQueryHash(zone.ConnectionId, zone.ChainId, "cosmos.distribution.v1beta1.Query/DelegationTotalRewards", map[string]string{"delegator": account.GetAddress()}))
	delegatorRewards := distrTypes.QueryDelegationTotalRewardsResponse{}
	if err == nil {
		k.cdc.MustUnmarshalJSON(delegatorRewardsDatapoint.Value, &delegatorRewards)
	}
	for _, delegation := range zone.GetDelegationsForDelegator(account.GetAddress()) {
		// maybe check if there are rewards to withdraw here? if there are we can delegate them in the same tx.
		amount := rewardsForDelegation(delegatorRewards, delegation.DelegationAddress, delegation.ValidatorAddress)
		fmt.Printf("Rewards for %s and %s: %v\n", delegation.DelegationAddress, delegation.ValidatorAddress, amount)
		msgs = append(msgs, &distrTypes.MsgWithdrawDelegatorReward{DelegatorAddress: delegation.GetDelegationAddress(), ValidatorAddress: delegation.GetValidatorAddress()})
		// send fee to fee collection account.
		// redelegate balance.
	}
	if len(msgs) == 0 {
		return nil
	}
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

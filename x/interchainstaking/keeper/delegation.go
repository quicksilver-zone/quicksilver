package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func (k *Keeper) Delegate(ctx sdk.Context, zone types.RegisteredZone, account *types.ICAAccount) error {
	var msgs []sdk.Msg

	for _, asset := range account.Balance {
		if asset.Denom == zone.GetBaseDenom() {
			k.Logger(ctx).Info("Sending a MsgDelegate!", "asset", asset)
			valoper_address, err := k.DetermineValidatorForDelegation(ctx, zone, account.Balance)
			if err != nil {
				panic("impossible!")
			}
			msgs = append(msgs, &stakingTypes.MsgDelegate{DelegatorAddress: account.GetAddress(), ValidatorAddress: valoper_address, Amount: asset})
		} else {
			k.Logger(ctx).Info("Sending a MsgRedeemTokensforShares!", "asset", asset)

			// validate this against validators?
			// if validator is not active, then redelegate msg too?
			msgs = append(msgs, &stakingTypes.MsgRedeemTokensforShares{DelegatorAddress: account.GetAddress(), Amount: asset})
		}
	}

	return k.SubmitTx(ctx, msgs, account)
}

func (k *Keeper) WithdrawDelegationRewards(ctx sdk.Context, zone types.RegisteredZone, account *types.ICAAccount) error {
	k.Logger(ctx).Info("Withdrawing rewards for delegate account", "account", account.GetAddress(), "zone", zone.ChainId)
	var msgs []sdk.Msg
	for _, delegation := range zone.GetDelegationsForDelegator(account.GetAddress()) {
		// maybe check if there are rewards to withdraw here? if there are we can delegate them in the same tx.
		msgs = append(msgs, &distrTypes.MsgWithdrawDelegatorReward{DelegatorAddress: delegation.GetDelegationAddress(), ValidatorAddress: delegation.GetValidatorAddress()})
	}
	if len(msgs) == 0 {
		return nil
	}
	k.Logger(ctx).Info("Submitting messages", "msgs", msgs)

	return k.SubmitTx(ctx, msgs, account)
}

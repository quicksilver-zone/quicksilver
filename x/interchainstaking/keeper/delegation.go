package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func (k *Keeper) Delegate(ctx sdk.Context, zone types.RegisteredZone, account *types.ICAAccount) error {

	var msgs []sdk.Msg

	for _, asset := range account.Balance {
		if asset.Denom == zone.GetDenom() {
			k.Logger(ctx).Info("Sending a MsgDelegate!", "asset", asset)
			// staking!
			// determine what the correct delegation is?
			msgs = append(msgs, &stakingTypes.MsgDelegate{DelegatorAddress: account.GetAddress(), ValidatorAddress: account.GetAddress(), Amount: asset})
		} else {
			k.Logger(ctx).Info("Sending a MsgRedeemTokensforShares!", "asset", asset)

			// validate this against validators?
			// if validator is not active, then redelegate msg too?
			msgs = append(msgs, &stakingTypes.MsgRedeemTokensforShares{DelegatorAddress: account.GetAddress(), Amount: asset})

		}
	}

	return k.SubmitTx(ctx, msgs, account)
}

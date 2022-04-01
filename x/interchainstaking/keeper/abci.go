package keeper

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
)

// BeginBlocker of interchainstaking module
func (k Keeper) BeginBlocker(ctx sdk.Context) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	if ctx.BlockHeight()%types.ValidatorSetInterval == 0 {
		k.IterateRegisteredZones(ctx, func(index int64, zoneInfo types.RegisteredZone) (stop bool) {
			k.Logger(ctx).Info("Setting validators for zone", "zone", zoneInfo.ChainId)
			// we must populate validators first, else the next piece fails :)
			validator_data, err := k.ICQKeeper.GetDatapoint(ctx, zoneInfo.ConnectionId, zoneInfo.ChainId, "cosmos.staking.v1beta1.Query/Validators", map[string]string{"status": stakingTypes.BondStatusBonded})
			if err != nil {
				k.Logger(ctx).Error("Unable to query validators for zone", "zone", zoneInfo.ChainId)
				return false
			}
			if validator_data.LocalHeight.LT(sdk.NewInt(ctx.BlockHeight() - types.DelegateDelegationsInterval)) {
				k.Logger(ctx).Error(fmt.Sprintf("Validators Info for zone is older than %d blocks", types.DelegateDelegationsInterval), "zone", zoneInfo.ChainId)
				return false
			}
			validatorsRes := stakingTypes.QueryValidatorsResponse{}
			err = k.cdc.UnmarshalJSON(validator_data.Value, &validatorsRes)
			if err != nil {
				k.Logger(ctx).Error("Unable to unmarshal validators info for zone", "zone", zoneInfo.ChainId, "err", err)
			}
			for _, validator := range validatorsRes.Validators {
				val, err := zoneInfo.GetValidatorByValoper(validator.OperatorAddress)
				if err != nil {
					k.Logger(ctx).Info("Unable to find validator - adding...", "valoper", validator.OperatorAddress)
					zoneInfo.Validators = append(zoneInfo.Validators, &types.Validator{
						ValoperAddress: validator.OperatorAddress,
						CommissionRate: validator.GetCommission(),
						VotingPower:    sdk.NewDecFromInt(validator.Tokens),
						Delegations:    []*types.Delegation{},
					})
				} else {
					if !val.CommissionRate.Equal(validator.GetCommission()) {
						val.CommissionRate = validator.GetCommission()
						k.Logger(ctx).Info("Validator commission rate change; updating...", "valoper", validator.OperatorAddress, "oldRate", val.CommissionRate, "newRate", validator.GetCommission())
					}

					if !val.VotingPower.Equal(sdk.NewDecFromInt(validator.Tokens)) {
						val.VotingPower = sdk.NewDecFromInt(validator.Tokens)
						k.Logger(ctx).Info("Validator voting power change; updating", "valoper", validator.OperatorAddress, "oldPower", val.VotingPower, "newPower", validator.Tokens.ToDec())
					}
				}

			}
			// also do this for Unbonded and Unbonding
			k.SetRegisteredZone(ctx, zoneInfo)

			return false
		})
	}

	// every N blocks, emit QueryAccountBalances event.
	if ctx.BlockHeight()%types.DepositInterval == 0 {
		k.IterateRegisteredZones(ctx, func(index int64, zoneInfo types.RegisteredZone) (stop bool) {
			// refactor me!
			balance_data, err := k.ICQKeeper.GetDatapoint(ctx, zoneInfo.ConnectionId, zoneInfo.ChainId, "cosmos.bank.v1beta1.Query/AllBalances", map[string]string{"address": zoneInfo.DepositAddress.GetAddress()})
			if err != nil {
				k.Logger(ctx).Error("Unable to query balance for deposit account", "deposit_address", zoneInfo.DepositAddress.GetAddress())
				return false
			}
			balanceRes := bankTypes.QueryAllBalancesResponse{}
			err = k.cdc.UnmarshalJSON(balance_data.Value, &balanceRes)
			if err != nil {
				k.Logger(ctx).Error("Unable to unmarshal balance for deposit account", "deposit_address", zoneInfo.DepositAddress.GetAddress(), "err", err)
			}
			balance := balanceRes.Balances
			if !balance.Empty() {
				k.Logger(ctx).Info("Balance is non zero", "existing", zoneInfo.DepositAddress.Balance, "current", balance)
				tx_data, err := k.ICQKeeper.GetDatapointOrRequest(ctx, zoneInfo.ConnectionId, zoneInfo.ChainId, "cosmos.tx.v1beta1.Query/GetTxEvents", map[string]string{"transfer.recipient": zoneInfo.DepositAddress.GetAddress()})
				if err != nil {
					// this happens, it's okay, we fetch the data async. we'll hit this loop again next iteration :)
					k.Logger(ctx).Info("No data yet. Ignoring...")
					return false
				}
				txs := coretypes.ResultTxSearch{}
				err = json.Unmarshal(tx_data.Value, &txs)
				if err != nil {
					k.Logger(ctx).Error("Unable to unmarshal txs for deposit account", "deposit_address", zoneInfo.DepositAddress.GetAddress())

				}
				for _, tx := range txs.Txs {
					k.HandleReceiptTransaction(ctx, tx, zoneInfo)
				}
				// update balance
				zoneInfo.DepositAddress.Balance = balance
				k.SetRegisteredZone(ctx, zoneInfo)

			}
			return false
		})
	}

	if ctx.BlockHeight()%types.DelegateInterval == 0 {
		// refactor me!
		k.IterateRegisteredZones(ctx, func(index int64, zoneInfo types.RegisteredZone) (stop bool) {
			for _, da := range zoneInfo.DelegationAddresses {
				balance_data, err := k.ICQKeeper.GetDatapoint(ctx, zoneInfo.ConnectionId, zoneInfo.ChainId, "cosmos.bank.v1beta1.Query/AllBalances", map[string]string{"address": da.GetAddress()})
				if err != nil {
					k.Logger(ctx).Error("Unable to query balance for delegate account", "delegate_address", da.GetAddress())
					continue
				}
				if balance_data.LocalHeight.LT(sdk.NewInt(ctx.BlockHeight() - types.DelegateInterval)) {
					k.Logger(ctx).Info(fmt.Sprintf("Balance for delegate account is older than %d blocks", types.DelegateInterval), "delegate_address", da.GetAddress())
					continue
				}
				balanceRes := bankTypes.QueryAllBalancesResponse{}
				err = k.cdc.UnmarshalJSON(balance_data.Value, &balanceRes)
				if err != nil {
					k.Logger(ctx).Error("Unable to unmarshal balance for delegate account", "delegation_address", zoneInfo.DepositAddress.GetAddress(), "err", err)
				}
				balance := balanceRes.Balances

				if !balance.Empty() {
					da.Balance = balance
					k.SetRegisteredZone(ctx, zoneInfo)
					k.Logger(ctx).Info("Delegate account balance is non-zero; delegating!", "current", balance)
					err := k.Delegate(ctx, zoneInfo, da)
					if err != nil {
						k.Logger(ctx).Error("Unable to delegate balances", "delegation_address", zoneInfo.DepositAddress.GetAddress(), "zone_identifier", zoneInfo.Identifier, "err", err)
					}
				}
			}
			return false
		})
	}

	if ctx.BlockHeight()%types.DelegateDelegationsInterval == 0 {
		k.IterateRegisteredZones(ctx, func(index int64, zoneInfo types.RegisteredZone) (stop bool) {
			// populate / handle delegations
			for _, da := range zoneInfo.DelegationAddresses {
				delegation_data, err := k.ICQKeeper.GetDatapoint(ctx, zoneInfo.ConnectionId, zoneInfo.ChainId, "cosmos.staking.v1beta1.Query/DelegatorDelegations", map[string]string{"address": da.GetAddress()})
				if err != nil {
					k.Logger(ctx).Error("Unable to query balance for delegate account", "delegate_address", da.GetAddress())
					continue
				}
				if delegation_data.LocalHeight.LT(sdk.NewInt(ctx.BlockHeight() - types.DelegateDelegationsInterval)) {
					k.Logger(ctx).Info(fmt.Sprintf("Delegations Info for delegate account is older than %d blocks", types.DelegateDelegationsInterval), "delegate_address", da.GetAddress())
					continue
				}
				delegationsRes := stakingTypes.QueryDelegatorDelegationsResponse{}
				err = k.cdc.UnmarshalJSON(delegation_data.Value, &delegationsRes)
				if err != nil {
					k.Logger(ctx).Error("Unable to unmarshal delegations info for delegate account", "delegation_address", zoneInfo.DepositAddress.GetAddress(), "err", err)
				}
				delegations := delegationsRes.DelegationResponses
				daBalance := sdk.Coin{Amount: sdk.ZeroInt(), Denom: zoneInfo.BaseDenom}
				for _, d := range delegations {
					delegator := d.Delegation.DelegatorAddress
					if delegator != da.GetAddress() {
						k.Logger(ctx).Error("Delegator mismatch", "d1", delegator, "d2", da.GetAddress())
						//panic("Delegator address mismatch") // is this a panic()????
					}
					delegatedCoins := d.Balance
					val, err := zoneInfo.GetValidatorByValoper(d.Delegation.ValidatorAddress)
					if err != nil {
						k.Logger(ctx).Error("Unable to find validator for delegation", "valoper", d.Delegation.ValidatorAddress)
					}
					delegation, err := val.GetDelegationForDelegator(da.GetAddress())
					if err != nil {
						k.Logger(ctx).Info("Adding delegation tuple", "delegator", da.GetAddress(), "validator", val.ValoperAddress, "amount", delegatedCoins.Amount)
						val.Delegations = append(val.Delegations, &types.Delegation{
							DelegationAddress: da.GetAddress(),
							ValidatorAddress:  val.ValoperAddress,
							Amount:            d.Balance.Amount.ToDec(),
							Rewards:           sdk.Coins{},
							RedelegationEnd:   0,
						})
					} else {
						if !delegation.Amount.Equal(delegatedCoins.Amount.ToDec()) {
							k.Logger(ctx).Info("Updating delegation tuple amount", "delegator", da.GetAddress(), "validator", val.ValoperAddress, "our_amount", delegation.Amount, "chain_amount", delegatedCoins.Amount)
							delegation.Amount = delegatedCoins.Amount.ToDec()
						}
					}
					daBalance.Add(delegatedCoins)
				}
				da.DelegatedBalance = daBalance
			}
			k.SetRegisteredZone(ctx, zoneInfo)

			return false
		})
	}
}

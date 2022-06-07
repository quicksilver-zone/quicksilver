package keeper

import (
	"encoding/json"
	"fmt"
	"strings"

	//lint:ignore SA1019 ignore this!
	"github.com/golang/protobuf/proto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icqtypes "github.com/ingenuity-build/quicksilver/x/interchainquery/types"
	queryTypes "github.com/ingenuity-build/quicksilver/x/interchainquery/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"

	"time"
)

func (k *Keeper) HandleAcknowledgement(ctx sdk.Context, packet channeltypes.Packet, acknowledgement []byte) error {
	ack := channeltypes.Acknowledgement_Result{}
	err := json.Unmarshal(acknowledgement, &ack)
	if err != nil {
		ackErr := channeltypes.Acknowledgement_Error{}
		err := json.Unmarshal(acknowledgement, &ackErr)
		if err != nil {
			k.Logger(ctx).Error("Unable to unmarshal acknowledgement error", "error", err, "data", acknowledgement)
			return err
		}
		k.Logger(ctx).Error("Unable to unmarshal acknowledgement result", "error", err, "remote_err", ackErr, "data", acknowledgement)
		return err
	}

	txMsgData := &sdk.TxMsgData{}
	err = proto.Unmarshal(ack.Result, txMsgData)
	if err != nil {
		k.Logger(ctx).Error("Unable to unmarshal acknowledgement", "error", err, "ack", ack.Result)
		return err
	}

	var packetData icatypes.InterchainAccountPacketData
	err = icatypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &packetData)
	if err != nil {
		k.Logger(ctx).Error("unable to unmarshal acknowledgement packet data", "error", err, "data", packetData)
		return err
	}
	msgs, err := icatypes.DeserializeCosmosTx(k.cdc, packetData.Data)
	if err != nil {
		k.Logger(ctx).Error("unable to decode messages", "err", err)
		return err
	}

	for msgIndex, msgData := range txMsgData.Data {
		src := msgs[msgIndex]
		switch msgData.MsgType {
		case "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward":
			response := distrtypes.MsgWithdrawDelegatorRewardResponse{}
			err := proto.Unmarshal(msgData.Data, &response)
			if err != nil {
				k.Logger(ctx).Error("Unable to unmarshal MsgWithdrawDelegatorReward response", "error", err)
				return err
			}
			k.Logger(ctx).Info("Rewards withdrawn", "response", response)
			if err := k.HandleWithdrawRewards(ctx, src, response.Amount); err != nil {
				return err
			}
			continue
		case "/cosmos.staking.v1beta1.MsgRedeemTokensforShares":
			response := stakingtypes.MsgRedeemTokensforSharesResponse{}
			err := proto.Unmarshal(msgData.Data, &response)
			if err != nil {
				k.Logger(ctx).Error("Unable to unmarshal MsgRedeemTokensforShares response", "error", err)
				return err
			}
			k.Logger(ctx).Debug("Tokens redeemed for shares", "response", response)
			// we should update delegation records here.
			if err := k.HandleRedeemTokens(ctx, src, response.Amount); err != nil {
				return err
			}
			continue
		case "/cosmos.staking.v1beta1.MsgTokenizeShares":
			response := stakingtypes.MsgTokenizeSharesResponse{}
			err := proto.Unmarshal(msgData.Data, &response)
			if err != nil {
				k.Logger(ctx).Error("Unable to unmarshal MsgTokenizeShares response", "error", err)
				return err
			}
			k.Logger(ctx).Info("Shares tokenized", "response", response)
			// check tokenizedShareTransfers (inc. rebalance and unbond)
			if err := k.HandleTokenizedShares(ctx, src, response.Amount); err != nil {
				return err
			}
			continue
		case "/cosmos.staking.v1beta1.MsgDelegate":
			response := stakingtypes.MsgDelegateResponse{}
			err := proto.Unmarshal(msgData.Data, &response)
			if err != nil {
				k.Logger(ctx).Error("Unable to unmarshal MsgDelegate response", "error", err)
				return err
			}
			k.Logger(ctx).Info("Delegated", "response", response)
			// we should update delegation records here.
			if err := k.HandleDelegate(ctx, src); err != nil {
				return err
			}
			continue
		case "/cosmos.staking.v1beta1.MsgBeginRedelegate":
			response := stakingtypes.MsgBeginRedelegateResponse{}
			err := proto.Unmarshal(msgData.Data, &response)
			if err != nil {
				k.Logger(ctx).Error("Unable to unmarshal MsgBeginRedelegate response", "error", err)
				return err
			}
			k.Logger(ctx).Debug("Redelegation initiated", "response", response)
			if err := k.HandleBeginRedelegate(ctx, src, response.CompletionTime); err != nil {
				return err
			}
			continue
		case "/cosmos.bank.v1beta1.MsgSend":
			response := banktypes.MsgSendResponse{}
			err := proto.Unmarshal(msgData.Data, &response)
			if err != nil {
				k.Logger(ctx).Error("Unable to unmarshal MsgSend response", "error", err)
				return err
			}
			k.Logger(ctx).Debug("Funds Transferred", "response", response)
			// check tokenTransfers - if end user unescrow and burn txs
			if err := k.HandleCompleteSend(ctx, src, packetData.Memo); err != nil {
				return err
			}
			continue
		case "/cosmos.bank.v1beta1.MsgMultiSend":
			response := banktypes.MsgMultiSendResponse{}
			err := proto.Unmarshal(msgData.Data, &response)
			if err != nil {
				k.Logger(ctx).Error("Unable to unmarshal MsgMultiSend response", "error", err)
				return err
			}
			k.Logger(ctx).Debug("Funds Transferred (Multi)", "response", response)
			if err := k.HandleCompleteMultiSend(ctx, src, packetData.Memo); err != nil {
				return err
			}
			continue
		case "/cosmos.distribution.v1beta1.MsgSetWithdrawAddress":
			response := distrtypes.MsgSetWithdrawAddressResponse{}
			err := proto.Unmarshal(msgData.Data, &response)
			if err != nil {
				k.Logger(ctx).Error("Unable to unmarshal MsgMultiSend response", "error", err)
				return err
			}
			k.Logger(ctx).Debug("Withdraw Address Updated", "response", response)
			if err := k.HandleUpdatedWithdrawAddress(ctx, src); err != nil {
				return err
			}
			continue
		case "/ibc.applications.transfer.v1.MsgTransfer":
			response := ibctransfertypes.MsgTransferResponse{}
			err := proto.Unmarshal(msgData.Data, &response)
			if err != nil {
				k.Logger(ctx).Error("Unable to unmarshal MsgTransfer response", "error", err)
				return err
			}
			k.Logger(ctx).Debug("MsgTranfer acknowledgement received")
			if err := k.HandleMsgTransfer(ctx, src); err != nil {
				return err
			}
			continue
		default:
			k.Logger(ctx).Error("Unhandled acknowledgement packet", "type", msgData.MsgType)
		}
	}

	return nil
}

func (k *Keeper) HandleTimeout(ctx sdk.Context, packet channeltypes.Packet) error {
	return nil
}

//----------------------------------------------------------------

func (k *Keeper) HandleMsgTransfer(ctx sdk.Context, msg sdk.Msg) error {
	k.Logger(ctx).Info("Received MsgTransfer acknowledgement")
	// first, type assertion. we should have ibctransfertypes.MsgTransfer
	sMsg, ok := msg.(*ibctransfertypes.MsgTransfer)
	if !ok {
		k.Logger(ctx).Error("unable to cast source message to MsgTransfer")
		return fmt.Errorf("unable to cast source message to MsgTransfer")
	}

	// check if destination is interchainstaking module account (spoiler: it was)
	if sMsg.Receiver != k.AccountKeeper.GetModuleAddress(types.ModuleName).String() {
		k.Logger(ctx).Error("MsgTransfer to unknown account!")
		return nil
	}

	return k.HandleDistributeFeesFromModuleAccount(ctx)
}

func (k *Keeper) HandleDistributeFeesFromModuleAccount(ctx sdk.Context) error {
	// what do we have in the account?
	balance := k.BankKeeper.GetAllBalances(ctx, k.AccountKeeper.GetModuleAddress(types.ModuleName))
	k.Logger(ctx).Info("distributing collected fees to stakers", "amount", balance)
	return k.BankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, authtypes.FeeCollectorName, balance) // Fee collector name needs to be passed in to keeper constructor.
}

func (k *Keeper) HandleCompleteMultiSend(ctx sdk.Context, msg sdk.Msg, memo string) error {
	k.Logger(ctx).Info("Received MsgMultiSend acknowledgement")
	// first, type assertion. we should have banktypes.MsgMultiSend
	sMsg, ok := msg.(*banktypes.MsgMultiSend)
	if !ok {
		k.Logger(ctx).Error("unable to cast source message to MsgMultiSend")
		return fmt.Errorf("unable to cast source message to MsgMultiSend")
	}

	// check for sending of tokens from deposit -> delegate.
	zone, err := k.GetZoneFromContext(ctx)
	if err != nil {
		k.Logger(ctx).Error(err.Error())
		return err
	}

	for _, out := range sMsg.Outputs {
		accAddr, err := types.AccAddressFromBech32(out.Address, "")
		plan := types.Allocations{}

		k.IterateAllDelegationPlansForHashAndDelegator(ctx, zone, memo, accAddr, func(delegationPlan types.DelegationPlan) bool {
			plan = plan.Allocate(delegationPlan.ValidatorAddress, delegationPlan.Value)
			k.RemoveDelegationPlan(ctx, zone, memo, delegationPlan)
			return false
		})

		da, err := zone.GetDelegationAccountByAddress(out.Address)
		if err != nil {
			k.Logger(ctx).Error(err.Error())
			return err
		}
		da.Balance = da.Balance.Add(out.Coins...)
		k.Delegate(ctx, *zone, da, plan)

	}

	return nil
}

func (k *Keeper) HandleCompleteSend(ctx sdk.Context, msg sdk.Msg, memo string) error {
	k.Logger(ctx).Info("Received MsgSend acknowledgement")
	// first, type assertion. we should have banktypes.MsgSend
	sMsg, ok := msg.(*banktypes.MsgSend)
	if !ok {
		err := fmt.Errorf("unable to cast source message to MsgSend")
		k.Logger(ctx).Error(err.Error())
		return err
	}

	// get zone
	zone, err := k.GetZoneFromContext(ctx)
	if err != nil {
		return err
	}

	// checks here are specific to ensure future extensibility;
	switch {
	case sMsg.FromAddress == zone.WithdrawalAddress.GetAddress():
		// WithdrawalAddress (for rewards) only send to DelegationAddresses.
		// Target here is one of the DelegationAddresses.
		if err := k.handleRewardsDelegation(ctx, *zone, sMsg); err != nil {
			return err
		}
	default:
		if err := k.handleWithdrawForUser(ctx, sMsg, memo); err != nil {
			return err
		}
	}

	return nil
}

func (k *Keeper) handleRewardsDelegation(ctx sdk.Context, zone types.RegisteredZone, msg *banktypes.MsgSend) error {
	da, err := zone.GetDelegationAccountByAddress(msg.ToAddress)
	if err != nil {
		return err
	}
	da.Balance = msg.Amount

	plan, err := types.DelegationPlanFromGlobalIntent(k.GetDelegationBinsMap(ctx, &zone), zone, sdk.NewCoin(zone.BaseDenom, msg.Amount.AmountOf(zone.BaseDenom)), zone.GetAggregateIntentOrDefault())
	if err != nil {
		return err
	}
	return k.Delegate(ctx, zone, da, plan)
}

func (k *Keeper) handleWithdrawForUser(ctx sdk.Context, msg *banktypes.MsgSend, memo string) error {
	var err error = nil
	var done bool = false

	// first check for withdrawals (if FromAddress is a DelegateAccount)
	k.IterateWithdrawalRecords(ctx, msg.FromAddress, func(idx int64, withdrawal types.WithdrawalRecord) bool {
		k.Logger(ctx).Debug("iterating withdraw record", "idx", idx, "record", withdrawal)
		if withdrawal.Recipient == msg.ToAddress {
			k.Logger(ctx).Debug("matched the recipient", "val", withdrawal.Delegator, "recipient", withdrawal.Recipient)
			z := k.GetZoneForDelegateAccount(ctx, withdrawal.Delegator)
			if msg.Amount.AmountOf(z.BaseDenom).Equal(withdrawal.Amount.Amount) {
				k.Logger(ctx).Debug("matched the amount", "amount", msg.Amount, "record.amount", withdrawal.Amount.Amount)
				if withdrawal.Status == WITHDRAW_STATUS_SEND {
					k.Logger(ctx).Info("Found matching withdrawal; withdrawal marked as completed")
					k.DeleteWithdrawalRecord(ctx, withdrawal.Delegator, withdrawal.Validator, withdrawal.Recipient)
					da, err := z.GetDelegationAccountByAddress(withdrawal.Delegator)
					if err != nil {
						return true
					}
					da.Balance = da.Balance.Sub(msg.Amount)

					done = true
					return true
				}
			}
		}
		return false
	})
	// after iteration, if we are marked done, exit cleanly.
	if done {
		return nil
	}

	// after iteration if we have an err, return it.
	if err != nil {
		return err
	}

	// second check for sending of tokens from deposit -> delegate.
	zone := k.GetZoneForDelegateAccount(ctx, msg.ToAddress)

	if zone != nil { // this _is_ a delegate account
		accAddr, err := types.AccAddressFromBech32(msg.ToAddress, "")
		plan := types.Allocations{}

		k.IterateAllDelegationPlansForHashAndDelegator(ctx, zone, memo, accAddr, func(delegationPlan types.DelegationPlan) bool {
			plan = plan.Allocate(delegationPlan.ValidatorAddress, delegationPlan.Value)
			k.RemoveDelegationPlan(ctx, zone, memo, delegationPlan)
			return false
		})

		da, err := zone.GetDelegationAccountByAddress(msg.ToAddress)
		if err != nil {
			return err
		}
		da.Balance = da.Balance.Add(msg.Amount...)
		k.Delegate(ctx, *zone, da, plan)
	}

	return nil
}

func (k *Keeper) HandleTokenizedShares(ctx sdk.Context, msg sdk.Msg, amount sdk.Coin) error {
	k.Logger(ctx).Info("Received MsgTokenizeShares acknowledgement")
	// first, type assertion. we should have stakingtypes.MsgTokenizeShares
	var err error = nil
	tsMsg, ok := msg.(*stakingtypes.MsgTokenizeShares)
	if !ok {
		k.Logger(ctx).Error("unable to cast source message to MsgTokenizeShares")
		return fmt.Errorf("unable to cast source message to MsgTokenizeShares")
	}
	// here we are either withdrawing for a user _or_ rebalancing internally. lets check both action queues:
	k.IterateWithdrawalRecords(ctx, tsMsg.DelegatorAddress, func(idx int64, withdrawal types.WithdrawalRecord) bool {
		k.Logger(ctx).Debug("iterating withdraw record", "idx", idx, "record", withdrawal)
		if strings.HasPrefix(amount.Denom, withdrawal.Validator) {
			k.Logger(ctx).Debug("matched the prefix", "token", amount.Denom, "denom", "val", withdrawal.Validator)
			if amount.Amount.Equal(withdrawal.Amount.Amount) {
				k.Logger(ctx).Debug("matched the amount", "amount", amount.Amount, "record.amount", withdrawal.Amount.Amount)
				if withdrawal.Status == WITHDRAW_STATUS_TOKENIZE {
					k.Logger(ctx).Info("Found matching withdrawal", "request_amount", withdrawal.Amount, "actual_amount", amount)
					// bingo!
					_, delegatorIca := k.GetICAForDelegateAccount(ctx, withdrawal.Delegator)
					if delegatorIca == nil {
						k.Logger(ctx).Error("unable to find delegator account for withdrawal; this shouldn't happen", err)
						return true
					}
					sendMsg := &banktypes.MsgSend{FromAddress: withdrawal.Delegator, ToAddress: withdrawal.Recipient, Amount: sdk.Coins{amount}}

					err = k.SubmitTx(ctx, []sdk.Msg{sendMsg}, delegatorIca, "")
					if err != nil {
						k.Logger(ctx).Error("error", err)
						return true
					}
					k.Logger(ctx).Info("Sending funds", "from", withdrawal.Delegator, "to", withdrawal.Recipient, "amount", amount)
					withdrawal.Status = WITHDRAW_STATUS_SEND
					k.SetWithdrawalRecord(ctx, &withdrawal)
					return true
				}
			}
		}
		return false
	})

	return err
}

func (k *Keeper) HandleBeginRedelegate(ctx sdk.Context, msg sdk.Msg, completion time.Time) error {
	panic("not implemented")
}

func (k *Keeper) HandleRedeemTokens(ctx sdk.Context, msg sdk.Msg, amount sdk.Coin) error {
	k.Logger(ctx).Info("Received MsgRedeemTokensforShares acknowledgement")
	// first, type assertion. we should have stakingtypes.MsgRedeemTokensforShares
	redeemMsg, ok := msg.(*stakingtypes.MsgRedeemTokensforShares)
	if !ok {
		k.Logger(ctx).Error("unable to cast source message to MsgRedeemTokensforShares")
		return fmt.Errorf("unable to cast source message to MsgRedeemTokensforShares")
	}
	validatorAddress, err := k.GetValidatorForToken(ctx, redeemMsg.DelegatorAddress, redeemMsg.Amount)
	if err != nil {
		return err
	}
	zone := k.GetZoneForDelegateAccount(ctx, redeemMsg.DelegatorAddress)

	err = k.UpdateDelegationRecordForAddress(ctx, redeemMsg.DelegatorAddress, validatorAddress, amount, zone, false)
	if err != nil {
		return err
	}

	k.SetRegisteredZone(ctx, *zone)
	return nil
}

func (k *Keeper) HandleDelegate(ctx sdk.Context, msg sdk.Msg) error {
	k.Logger(ctx).Info("Received MsgDelegate acknowledgement")
	// first, type assertion. we should have stakingtypes.MsgDelegate
	delegateMsg, ok := msg.(*stakingtypes.MsgDelegate)
	if !ok {
		k.Logger(ctx).Error("unable to cast source message to MsgDelegate")
		return fmt.Errorf("unable to cast source message to MsgDelegate")
	}
	zone := k.GetZoneForDelegateAccount(ctx, delegateMsg.DelegatorAddress)
	if zone == nil {
		// most likely a performance account...
		if zone := k.GetZoneForPerformanceAccount(ctx, delegateMsg.DelegatorAddress); zone != nil {
			return nil
		} else {
			return fmt.Errorf("unable to find zone for address %s", delegateMsg.DelegatorAddress)
		}
	}

	err := k.UpdateDelegationRecordForAddress(ctx, delegateMsg.DelegatorAddress, delegateMsg.ValidatorAddress, delegateMsg.Amount, zone, false)
	if err != nil {
		return err
	}

	k.SetRegisteredZone(ctx, *zone)
	return nil
}

func (k *Keeper) HandleUpdatedWithdrawAddress(ctx sdk.Context, msg sdk.Msg) error {
	k.Logger(ctx).Info("Received MsgSetWithdrawAddress acknowledgement")
	// first, type assertion. we should have distrtypes.MsgSetWithdrawAddress
	_, ok := msg.(*distrtypes.MsgSetWithdrawAddress)
	if !ok {
		k.Logger(ctx).Error("unable to cast source message to MsgSetWithdrawAddress")
		return fmt.Errorf("unable to cast source message to MsgSetWithdrawAddress")
	}

	return nil
}

// TODO: this should be part of Keeper, but part of zone. Refactor me.
func (k *Keeper) GetValidatorForToken(ctx sdk.Context, delegatorAddress string, amount sdk.Coin) (string, error) {
	zone, err := k.GetZoneFromContext(ctx)
	if err != nil {
		k.Logger(ctx).Error(err.Error())
		return "", err
	}

	for _, val := range zone.GetValidatorsAddressesAsSlice() {
		if strings.HasPrefix(amount.Denom, val) {
			// match!
			return val, nil
		}
	}

	return "", fmt.Errorf("unable to find validator for token %s", amount.Denom)

}

var delegationCb Callback = func(k Keeper, ctx sdk.Context, args []byte, query icqtypes.Query) error {
	zone, found := k.GetRegisteredZoneInfo(ctx, query.GetChainId())
	if !found {
		return fmt.Errorf("no registered zone for chain id: %s", query.GetChainId())
	}

	delegation := stakingtypes.Delegation{}
	err := k.cdc.Unmarshal(args, &delegation)
	if err != nil {
		return err
	}

	val, err := zone.GetValidatorByValoper(delegation.ValidatorAddress)
	if err != nil {
		return err
	}

	return k.UpdateDelegationRecordForAddress(ctx, delegation.DelegatorAddress, delegation.ValidatorAddress, sdk.NewCoin(zone.BaseDenom, val.SharesToTokens(delegation.Shares)), &zone, true)
}

func (k *Keeper) UpdateDelegationRecordsForAddress(ctx sdk.Context, zone *types.RegisteredZone, delegatorAddress string, args []byte) error {
	var response stakingtypes.QueryDelegatorDelegationsResponse
	err := k.cdc.Unmarshal(args, &response)
	if err != nil {
		return err
	}

	_, delAddr, _ := bech32.DecodeAndConvert(delegatorAddress)
	delegatorDelegations := k.GetDelegatorDelegations(ctx, zone, delAddr)
	delMap := make(map[string]types.Delegation, len(delegatorDelegations))
	for _, del := range delegatorDelegations {
		delMap[del.ValidatorAddress] = del
	}

	da, err := zone.GetDelegationAccountByAddress(delegatorAddress)
	if err != nil {
		return err
	}

	for _, delegationRecord := range response.DelegationResponses {

		_, valAddr, _ := bech32.DecodeAndConvert(delegationRecord.Delegation.ValidatorAddress)
		data := stakingtypes.GetDelegationKey(delAddr, valAddr)

		delegation, ok := delMap[delegationRecord.Delegation.ValidatorAddress]
		if !ok || !delegation.Amount.Equal(delegationRecord.GetBalance()) { // new or updated delegation
			k.Logger(ctx).Info("Outdated delegation record - fetching proof...", "valoper", delegationRecord.Delegation.ValidatorAddress)

			k.ICQKeeper.MakeRequest(
				ctx,
				zone.ConnectionId,
				zone.ChainId,
				"store/staking/key",
				data,
				sdk.NewInt(-1),
				types.ModuleName,
				delegationCb,
			)
		}

		if ok {
			delete(delMap, delegationRecord.Delegation.ValidatorAddress)
		}
	}

	for _, existingDelegation := range delMap {
		_, valAddr, _ := bech32.DecodeAndConvert(existingDelegation.ValidatorAddress)
		data := stakingtypes.GetDelegationKey(delAddr, valAddr)

		k.RemoveDelegation(ctx, zone, existingDelegation)
		da.DelegatedBalance = da.DelegatedBalance.Sub(existingDelegation.Amount) // remove old delegation from da.DelegatedBalance
		// send request to prove delegation no longer exists.
		k.ICQKeeper.MakeRequest(
			ctx,
			zone.ConnectionId,
			zone.ChainId,
			"store/staking/key",
			data,
			sdk.NewInt(-1),
			types.ModuleName,
			delegationCb,
		)
	}

	k.SetRegisteredZone(ctx, *zone)

	return nil
}

func (k *Keeper) UpdateDelegationRecordForAddress(ctx sdk.Context, delegatorAddress string, validatorAddress string, amount sdk.Coin, zone *types.RegisteredZone, absolute bool) error {

	delegation, found := k.GetDelegation(ctx, zone, delegatorAddress, validatorAddress)
	da, _ := zone.GetDelegationAccountByAddress(delegatorAddress)

	if !found {
		k.Logger(ctx).Info("Adding delegation tuple", "delegator", delegatorAddress, "validator", validatorAddress, "amount", amount.Amount)
		delegation = types.NewDelegation(delegatorAddress, validatorAddress, amount)
		da.DelegatedBalance = da.DelegatedBalance.Add(amount)
	} else {
		if !delegation.Amount.Equal(amount.Amount.ToDec()) {
			k.Logger(ctx).Info("Updating delegation tuple amount", "delegator", delegatorAddress, "validator", validatorAddress, "old_amount", delegation.Amount, "inbound_amount", amount.Amount, "abs", absolute)
			if !absolute {
				da.DelegatedBalance = da.DelegatedBalance.Add(amount)
				delegation.Amount = delegation.Amount.Add(amount)
			} else {
				da.DelegatedBalance = da.DelegatedBalance.Sub(delegation.Amount).Add(amount)
				delegation.Amount = amount
			}
		}
	}
	k.SetDelegation(ctx, zone, delegation)
	k.EmitValsetRequery(ctx, zone.ConnectionId, zone.ChainId, -1)
	k.SetRegisteredZone(ctx, *zone)

	return nil
}

func (k *Keeper) HandleWithdrawRewards(ctx sdk.Context, msg sdk.Msg, amount sdk.Coins) error {
	k.Logger(ctx).Info("Received MsgWithdrawDelegatorReward acknowledgement")
	// ? we don't actually need the distrtypes.MsgWithdrawDelegatorReward here
	// as it just returns the delegator:validator tuple that we sent ?
	zone, err := k.GetZoneFromContext(ctx)
	if err != nil {
		k.Logger(ctx).Error(err.Error())
		return err
	}
	// decrement withdrawal waitgroup
	zone.WithdrawalWaitgroup--
	k.SetRegisteredZone(ctx, *zone)

	switch zone.WithdrawalWaitgroup {
	case 0:
		// interface assertion
		var cb Callback = DistributeRewardsFromWithdrawAccount

		balanceQuery := banktypes.QueryAllBalancesRequest{Address: zone.WithdrawalAddress.Address}
		bz, err := k.cdc.Marshal(&balanceQuery)
		if err != nil {
			return err
		}

		// total rewards balance withdrawn
		k.ICQKeeper.MakeRequest(
			ctx,
			zone.ConnectionId,
			zone.ChainId,
			"cosmos.bank.v1beta1.Query/AllBalances",
			bz,
			sdk.NewInt(int64(-1)),
			types.ModuleName,
			cb,
		)
		return nil
	default:
		return nil
	}
}

func DistributeRewardsFromWithdrawAccount(k Keeper, ctx sdk.Context, args []byte, query queryTypes.Query) error {
	zone, found := k.GetRegisteredZoneInfo(ctx, query.ChainId)
	if !found {
		return fmt.Errorf("unable to find zone for %s", query.ChainId)
	}

	// query all balances as chains can accumulate fees in different denoms.
	withdrawBalance := banktypes.QueryAllBalancesResponse{}

	err := k.cdc.Unmarshal(args, &withdrawBalance)
	if err != nil {
		return err
	}
	baseDenomAmount := withdrawBalance.Balances.AmountOf(zone.BaseDenom)
	// calculate fee (fee = amount * rate)

	baseDenomFee := baseDenomAmount.ToDec().
		Mul(k.GetCommissionRate(ctx)).
		TruncateInt()

	// prepare rewards distribution
	rewards := sdk.NewCoin(zone.BaseDenom, baseDenomAmount.Sub(baseDenomFee))

	dust, msgs := k.prepareRewardsDistributionMsgs(ctx, zone, rewards)

	// subtract dust from rewards
	rewards = rewards.SubAmount(dust)

	// multiDenomFee is the balance of withdrawal account minus the redelegated rewards.
	multiDenomFee := withdrawBalance.Balances.Sub(sdk.Coins{rewards})

	channelReq := channeltypes.QueryConnectionChannelsRequest{Connection: zone.ConnectionId}
	localChannelResp, err := k.IBCKeeper.ChannelKeeper.ConnectionChannels(sdk.WrapSDKContext(ctx), &channelReq)
	if err != nil {
		return err
	}
	var remotePort string
	var remoteChannel string
	for _, localChannel := range localChannelResp.Channels {
		if localChannel.PortId == "transfer" {
			remoteChannel = localChannel.Counterparty.ChannelId
			remotePort = localChannel.Counterparty.PortId
			break
		}
	}
	if remotePort == "" {
		return fmt.Errorf("unable to find remote transfer connection")
	}

	for _, coin := range multiDenomFee {
		msgs = append(
			msgs,
			&ibctransfertypes.MsgTransfer{
				SourcePort:       remotePort,
				SourceChannel:    remoteChannel,
				Token:            coin,
				Sender:           zone.WithdrawalAddress.Address,
				Receiver:         k.AccountKeeper.GetModuleAddress(types.ModuleName).String(),
				TimeoutTimestamp: uint64(ctx.BlockTime().UnixNano() + 5*time.Minute.Nanoseconds()),
				TimeoutHeight:    clienttypes.Height{RevisionNumber: 0, RevisionHeight: 0},
			},
		)
	}

	// update redemption rate
	k.updateRedemptionRate(ctx, zone, rewards)

	// send tx
	return k.SubmitTx(ctx, msgs, zone.WithdrawalAddress, "")

}

func (k *Keeper) updateRedemptionRate(ctx sdk.Context, zone types.RegisteredZone, epochRewards sdk.Coin) {
	ratio := zone.GetDelegatedAmount().Add(epochRewards).Amount.ToDec().Quo(k.BankKeeper.GetSupply(ctx, zone.LocalDenom).Amount.ToDec())
	k.Logger(ctx).Info("Epochly rewards", "coins", epochRewards)
	k.Logger(ctx).Info("Last redemption rate", "rate", zone.LastRedemptionRate)
	k.Logger(ctx).Info("Current redemption rate", "rate", zone.RedemptionRate)
	k.Logger(ctx).Info("New redemption rate", "rate", ratio, "supply", k.BankKeeper.GetSupply(ctx, zone.LocalDenom).Amount.ToDec(), "lv", zone.GetDelegatedAmount().Add(epochRewards).Amount.ToDec())

	zone.LastRedemptionRate = zone.RedemptionRate
	zone.RedemptionRate = ratio
	k.SetRegisteredZone(ctx, zone)
}

func (k *Keeper) prepareRewardsDistributionMsgs(ctx sdk.Context, zone types.RegisteredZone, rewards sdk.Coin) (sdk.Int, []sdk.Msg) {
	// todo: use multisend.
	// todo: this will probably not want to be an equal distribution. we want to use this to even out the distribution between accounts.
	var msgs []sdk.Msg

	dust := rewards.Amount
	portion := rewards.Amount.ToDec().Quo(sdk.NewDec(int64(len(zone.DelegationAddresses)))).TruncateInt()
	for _, da := range zone.GetDelegationAccounts() {
		msgs = append(
			msgs,
			&banktypes.MsgSend{
				FromAddress: zone.WithdrawalAddress.GetAddress(),
				ToAddress:   da.GetAddress(),
				Amount:      sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, portion)),
			},
		)
		dust = dust.Sub(portion)
	}

	return dust, msgs
}

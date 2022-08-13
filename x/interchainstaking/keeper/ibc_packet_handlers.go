package keeper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	//nolint:staticcheck
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

	"github.com/ingenuity-build/quicksilver/utils"
	queryTypes "github.com/ingenuity-build/quicksilver/x/interchainquery/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func (k *Keeper) HandleAcknowledgement(ctx sdk.Context, packet channeltypes.Packet, acknowledgement []byte) error {
	ack := channeltypes.Acknowledgement_Result{}
	err := json.Unmarshal(acknowledgement, &ack)
	if err != nil {
		ackErr := channeltypes.Acknowledgement_Error{}
		err := json.Unmarshal(acknowledgement, &ackErr)
		if err != nil {
			k.Logger(ctx).Error("unable to unmarshal acknowledgement error", "error", err, "data", acknowledgement)
			return err
		}
		k.Logger(ctx).Error("unable to unmarshal acknowledgement result", "error", err, "remote_err", ackErr, "data", acknowledgement)
		return err
	}

	txMsgData := &sdk.TxMsgData{}
	err = proto.Unmarshal(ack.Result, txMsgData)
	if err != nil {
		k.Logger(ctx).Error("unable to unmarshal acknowledgement", "error", err, "ack", ack.Result)
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
			k.Logger(ctx).Info("Rewards withdrawn")
			if err := k.HandleWithdrawRewards(ctx, src); err != nil {
				return err
			}
			continue
		case "/cosmos.staking.v1beta1.MsgRedeemTokensforShares":
			response := stakingtypes.MsgRedeemTokensforSharesResponse{}
			err := proto.Unmarshal(msgData.Data, &response)
			if err != nil {
				k.Logger(ctx).Error("unable to unmarshal MsgRedeemTokensforShares response", "error", err)
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
				k.Logger(ctx).Error("unable to unmarshal MsgTokenizeShares response", "error", err)
				return err
			}
			k.Logger(ctx).Info("Shares tokenized", "response", response)
			// check tokenizedShareTransfers (inc. rebalance and unbond)
			if err := k.HandleTokenizedShares(ctx, src, response.Amount, packetData.Memo); err != nil {
				return err
			}
			continue
		case "/cosmos.staking.v1beta1.MsgDelegate":
			response := stakingtypes.MsgDelegateResponse{}
			err := proto.Unmarshal(msgData.Data, &response)
			if err != nil {
				k.Logger(ctx).Error("unable to unmarshal MsgDelegate response", "error", err)
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
				k.Logger(ctx).Error("unable to unmarshal MsgBeginRedelegate response", "error", err)
				return err
			}
			k.Logger(ctx).Debug("Redelegation initiated", "response", response)
			if err := k.HandleBeginRedelegate(ctx, src, response.CompletionTime); err != nil {
				return err
			}
			continue
		case "/cosmos.staking.v1beta1.MsgUndelegate":
			response := stakingtypes.MsgUndelegateResponse{}
			err := proto.Unmarshal(msgData.Data, &response)
			if err != nil {
				k.Logger(ctx).Error("unable to unmarshal MsgDelegate response", "error", err)
				return err
			}
			k.Logger(ctx).Info("Undelegation started", "response", response)
			// we should update delegation records here.
			if err := k.HandleUndelegate(ctx, src, response.CompletionTime, packetData.Memo); err != nil {
				return err
			}
			continue
		case "/cosmos.bank.v1beta1.MsgSend":
			response := banktypes.MsgSendResponse{}
			err := proto.Unmarshal(msgData.Data, &response)
			if err != nil {
				k.Logger(ctx).Error("unable to unmarshal MsgSend response", "error", err)
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
				k.Logger(ctx).Error("unable to unmarshal MsgMultiSend response", "error", err)
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
				k.Logger(ctx).Error("unable to unmarshal MsgMultiSend response", "error", err)
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
				k.Logger(ctx).Error("unable to unmarshal MsgTransfer response", "error", err)
				return err
			}
			k.Logger(ctx).Debug("MsgTranfer acknowledgement received")
			if err := k.HandleMsgTransfer(ctx, src); err != nil {
				return err
			}
			continue
		default:
			k.Logger(ctx).Error("unhandled acknowledgement packet", "type", msgData.MsgType)
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
		k.Logger(ctx).Error("msgTransfer to unknown account!")
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
		err = fmt.Errorf("1: %w", err)
		k.Logger(ctx).Error(err.Error())
		return err
	}

	for _, out := range sMsg.Outputs {
		accAddr, err := utils.AccAddressFromBech32(out.Address, zone.AccountPrefix)
		if err != nil {
			return err
		}

		plan := types.Allocations{}

		// NOTE: deleting mid-iteration breaks the iterator; cache the results and delete retrospectively.
		toDelete := []types.DelegationPlan{}

		k.IterateAllDelegationPlansForHashAndDelegator(ctx, zone, memo, accAddr, func(delegationPlan types.DelegationPlan) bool {
			plan = plan.Allocate(delegationPlan.ValidatorAddress, delegationPlan.Value)
			toDelete = append(toDelete, delegationPlan)
			return false
		})

		for _, delegationPlan := range toDelete {
			if err := k.RemoveDelegationPlan(ctx, zone, memo, delegationPlan); err != nil {
				return err
			}
		}

		da, err := zone.GetDelegationAccountByAddress(out.Address)
		if err != nil {
			return err
		}
		da.Balance = da.Balance.Add(out.Coins...)
		if err = k.Delegate(ctx, *zone, da, plan); err != nil {
			return err
		}
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
		err = fmt.Errorf("2: %w", err)
		k.Logger(ctx).Error(err.Error())
		return err
	}

	// checks here are specific to ensure future extensibility;
	switch {
	case sMsg.FromAddress == zone.WithdrawalAddress.GetAddress():
		// WithdrawalAddress (for rewards) only send to DelegationAddresses.
		// Target here is one of the DelegationAddresses.
		return k.handleRewardsDelegation(ctx, *zone, sMsg)
	case zone.IsDelegateAddress(sMsg.FromAddress):
		return k.handleWithdrawForUser(ctx, zone, sMsg, memo)
	case zone.IsDelegateAddress(sMsg.ToAddress) && zone.DepositAddress.Address == sMsg.FromAddress:
		return k.handleSendToDelegate(ctx, zone, sMsg, memo)
	default:
		err = fmt.Errorf("unexpected completed send")
		k.Logger(ctx).Error(err.Error())
		return err
	}
}

func (k *Keeper) handleRewardsDelegation(ctx sdk.Context, zone types.Zone, msg *banktypes.MsgSend) error {
	da, err := zone.GetDelegationAccountByAddress(msg.ToAddress)
	if err != nil {
		return err
	}
	da.Balance = msg.Amount

	plan, err := types.DelegationPlanFromGlobalIntent(k.GetDelegatedAmount(ctx, &zone), k.GetDelegationBinsMap(ctx, &zone), sdk.NewCoin(zone.BaseDenom, msg.Amount.AmountOf(zone.BaseDenom)), zone.GetAggregateIntentOrDefault())
	if err != nil {
		return err
	}
	return k.Delegate(ctx, zone, da, plan)
}

func (k *Keeper) handleSendToDelegate(ctx sdk.Context, zone *types.Zone, msg *banktypes.MsgSend, memo string) error {
	accAddr, err := utils.AccAddressFromBech32(msg.ToAddress, zone.AccountPrefix)
	if err != nil {
		return err
	}
	plan := types.Allocations{}

	// NOTE: deleting mid-iteration breaks the iterator; cache the results and delete retrospectively.
	toDelete := []types.DelegationPlan{}

	k.IterateAllDelegationPlansForHashAndDelegator(ctx, zone, memo, accAddr, func(delegationPlan types.DelegationPlan) bool {
		plan = plan.Allocate(delegationPlan.ValidatorAddress, delegationPlan.Value)
		toDelete = append(toDelete, delegationPlan)
		return false
	})

	for _, delegationPlan := range toDelete {
		if err := k.RemoveDelegationPlan(ctx, zone, memo, delegationPlan); err != nil {
			return err
		}
	}

	da, err := zone.GetDelegationAccountByAddress(msg.ToAddress)
	if err != nil {
		return err
	}
	da.Balance = da.Balance.Add(msg.Amount...)
	return k.Delegate(ctx, *zone, da, plan)
}

func (k *Keeper) handleWithdrawForUser(ctx sdk.Context, zone *types.Zone, msg *banktypes.MsgSend, memo string) error {
	var err error
	// first check for withdrawals (if FromAddress is a DelegateAccount)
	k.IterateZoneDelegatorHashWithdrawalRecords(ctx, zone, memo, msg.FromAddress, func(idx int64, withdrawal types.WithdrawalRecord) bool {
		if withdrawal.Recipient == msg.ToAddress {
			k.Logger(ctx).Info("matched the recipient", "val", withdrawal.Delegator, "recipient", withdrawal.Recipient)
			if msg.Amount[0].Amount.Equal(withdrawal.Amount.Amount) {
				k.Logger(ctx).Info("matched the amount", "amount", msg.Amount, "record.amount", withdrawal.Amount.Amount)
				if withdrawal.Status == WithdrawStatusSend {
					k.Logger(ctx).Info("Found matching withdrawal; withdrawal marked as completed")
					k.DeleteWithdrawalRecord(ctx, zone, memo, withdrawal.Delegator, withdrawal.Validator)
					if len(k.AllZoneDelegatorHashWithdrawalRecords(ctx, zone, memo, withdrawal.Delegator)) == 0 {
						err = k.BankKeeper.BurnCoins(ctx, types.ModuleName, sdk.Coins{withdrawal.BurnAmount})
						if err != nil {
							return false
						}
						k.Logger(ctx).Info("burned coins post-withdrawal", "coins", withdrawal.BurnAmount)
					}

					err = k.EmitValsetRequery(ctx, zone.ConnectionId, zone.ChainId)
					return err != nil
				}
			}
		}
		return false
	})
	return err
}

func (k *Keeper) HandleCompletedUnbondings(ctx sdk.Context, zone *types.Zone) error {
	var err error
	k.IterateZoneWithdrawalRecords(ctx, zone, func(idx int64, withdrawal types.WithdrawalRecord) bool {
		k.Logger(ctx).Info("iterating unbondings")
		if withdrawal.Status == WithdrawStatusUnbond && withdrawal.CompletionTime.After(ctx.BlockTime()) { // completion date has passed.
			k.Logger(ctx).Info("matched unbonding")

			// bingo!
			_, delegatorIca := k.GetICAForDelegateAccount(ctx, withdrawal.Delegator)
			if delegatorIca == nil {
				k.Logger(ctx).Error("unable to find delegator account for withdrawal; this shouldn't happen", err)
				return true
			}
			sendMsg := &banktypes.MsgSend{FromAddress: withdrawal.Delegator, ToAddress: withdrawal.Recipient, Amount: sdk.Coins{withdrawal.Amount}}

			err = k.SubmitTx(ctx, []sdk.Msg{sendMsg}, delegatorIca, withdrawal.Txhash)
			if err != nil {
				k.Logger(ctx).Error("error", err)
				return true
			}
			k.Logger(ctx).Info("sending funds", "from", withdrawal.Delegator, "to", withdrawal.Recipient, "amount", withdrawal.Amount)
			withdrawal.Status = WithdrawStatusSend
			k.SetWithdrawalRecord(ctx, &withdrawal)
		}
		return false
	})
	return err
}

func (k *Keeper) HandleTokenizedShares(ctx sdk.Context, msg sdk.Msg, amount sdk.Coin, memo string) error {
	k.Logger(ctx).Info("Received MsgTokenizeShares acknowledgement")
	// first, type assertion. we should have stakingtypes.MsgTokenizeShares
	var err error
	tsMsg, ok := msg.(*stakingtypes.MsgTokenizeShares)
	if !ok {
		k.Logger(ctx).Error("unable to cast source message to MsgTokenizeShares")
		return fmt.Errorf("unable to cast source message to MsgTokenizeShares")
	}

	zone := k.GetZoneForDelegateAccount(ctx, tsMsg.DelegatorAddress)
	// here we are either withdrawing for a user _or_ rebalancing internally. lets check both action queues:
	k.IterateZoneDelegatorHashWithdrawalRecords(ctx, zone, memo, tsMsg.DelegatorAddress, func(idx int64, withdrawal types.WithdrawalRecord) bool {
		k.Logger(ctx).Debug("iterating withdraw record", "idx", idx, "record", withdrawal)
		if strings.HasPrefix(amount.Denom, withdrawal.Validator) {
			k.Logger(ctx).Debug("matched the prefix", "token", amount.Denom, "denom", "val", withdrawal.Validator)
			if amount.Amount.Equal(withdrawal.Amount.Amount) {
				k.Logger(ctx).Debug("matched the amount", "amount", amount.Amount, "record.amount", withdrawal.Amount.Amount)
				if withdrawal.Status == WithdrawStatusTokenize {
					k.Logger(ctx).Info("Found matching withdrawal", "request_amount", withdrawal.Amount, "actual_amount", amount)
					// bingo!
					_, delegatorIca := k.GetICAForDelegateAccount(ctx, withdrawal.Delegator)
					if delegatorIca == nil {
						k.Logger(ctx).Error("unable to find delegator account for withdrawal; this shouldn't happen", err)
						return true
					}
					sendMsg := &banktypes.MsgSend{FromAddress: withdrawal.Delegator, ToAddress: withdrawal.Recipient, Amount: sdk.Coins{amount}}

					err = k.SubmitTx(ctx, []sdk.Msg{sendMsg}, delegatorIca, memo)
					if err != nil {
						k.Logger(ctx).Error("error", err)
						return true
					}
					k.Logger(ctx).Info("sending funds", "from", withdrawal.Delegator, "to", withdrawal.Recipient, "amount", amount)
					withdrawal.Status = WithdrawStatusSend
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

func (k *Keeper) HandleUndelegate(ctx sdk.Context, msg sdk.Msg, completion time.Time, hash string) error {
	k.Logger(ctx).Info("Received MsgUndelegate acknowledgement")
	// first, type assertion. we should have stakingtypes.MsgUndelegate
	undelegateMsg, ok := msg.(*stakingtypes.MsgUndelegate)
	if !ok {
		k.Logger(ctx).Error("unable to cast source message to MsgUndelegate")
		return fmt.Errorf("unable to cast source message to MsgUndelegate")
	}
	zone := k.GetZoneForDelegateAccount(ctx, undelegateMsg.DelegatorAddress)
	k.Logger(ctx).Info("MsgUndelegate", "del", undelegateMsg.DelegatorAddress, "val", undelegateMsg.ValidatorAddress, "hash", hash, "chain", zone.ChainId)
	record, found := k.GetWithdrawalRecord(ctx, zone, hash, undelegateMsg.DelegatorAddress, undelegateMsg.ValidatorAddress)
	if !found {
		return fmt.Errorf("unable to lookup withdrawal record")
	}
	record.CompletionTime = completion
	k.Logger(ctx).Error("record to save", "rcd", record)
	k.SetWithdrawalRecord(ctx, &record)

	delegationQuery := stakingtypes.QueryDelegatorDelegationsRequest{DelegatorAddr: undelegateMsg.DelegatorAddress}
	bz := k.cdc.MustMarshal(&delegationQuery)

	k.ICQKeeper.MakeRequest(
		ctx,
		zone.ConnectionId,
		zone.ChainId,
		"cosmos.staking.v1beta1.Query/DelegatorDelegations",
		bz,
		sdk.NewInt(-1),
		types.ModuleName,
		"delegations",
		0,
	)
	return nil
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

	return k.UpdateDelegationRecordForAddress(ctx, redeemMsg.DelegatorAddress, validatorAddress, amount, zone, false)
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
		}
		return fmt.Errorf("unable to find zone for address %s", delegateMsg.DelegatorAddress)

	}

	return k.UpdateDelegationRecordForAddress(ctx, delegateMsg.DelegatorAddress, delegateMsg.ValidatorAddress, delegateMsg.Amount, zone, false)
}

func (k *Keeper) HandleUpdatedWithdrawAddress(ctx sdk.Context, msg sdk.Msg) error {
	k.Logger(ctx).Info("Received MsgSetWithdrawAddress acknowledgement")
	// first, type assertion. we should have distrtypes.MsgSetWithdrawAddress
	original, ok := msg.(*distrtypes.MsgSetWithdrawAddress)
	if !ok {
		k.Logger(ctx).Error("unable to cast source message to MsgSetWithdrawAddress")
		return fmt.Errorf("unable to cast source message to MsgSetWithdrawAddress")
	}
	zone, ica := k.GetICAForDelegateAccount(ctx, original.DelegatorAddress)
	if zone == nil {
		zone = k.GetZoneForPerformanceAccount(ctx, original.DelegatorAddress)
		if zone == nil {
			return fmt.Errorf("unable to find zone")
		}
		if err := zone.PerformanceAddress.SetWithdrawalAddress(original.WithdrawAddress); err != nil {
			return err
		}
	} else {
		if err := ica.SetWithdrawalAddress(original.WithdrawAddress); err != nil {
			return err
		}
	}
	k.SetZone(ctx, zone)

	return nil
}

// TODO: this should be part of Keeper, but part of zone. Refactor me.
func (k *Keeper) GetValidatorForToken(ctx sdk.Context, delegatorAddress string, amount sdk.Coin) (string, error) {
	zone, err := k.GetZoneFromContext(ctx)
	if err != nil {
		err = fmt.Errorf("3: %w", err)
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

func parseDelegationKey(key []byte) ([]byte, []byte, error) {
	if !bytes.Equal(key[0:1], []byte{0x31}) {
		return []byte{}, []byte{}, fmt.Errorf("not a valid delegation key")
	}
	delAddrLen := key[1]
	delAddr := key[2:delAddrLen]
	// valAddrLen := key[2+delAddrLen]
	valAddr := key[3+delAddrLen:]
	return delAddr, valAddr, nil
}

func (k *Keeper) UpdateDelegationRecordsForAddress(ctx sdk.Context, zone *types.Zone, delegatorAddress string, args []byte) error {
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
				"delegation",
				0,
			)
			da.IncrementBalanceWaitgroup()
		}

		if ok {
			delete(delMap, delegationRecord.Delegation.ValidatorAddress)
		}
	}

	sortedLeftAddrs := make([]string, 0, len(delMap))
	for valAddr := range delMap {
		sortedLeftAddrs = append(sortedLeftAddrs, valAddr)
	}
	sort.Strings(sortedLeftAddrs)

	for _, existingValAddr := range sortedLeftAddrs {
		existingDelegation := delMap[existingValAddr]
		_, valAddr, _ := bech32.DecodeAndConvert(existingDelegation.ValidatorAddress)
		data := stakingtypes.GetDelegationKey(delAddr, valAddr)

		if err := k.RemoveDelegation(ctx, zone, existingDelegation); err != nil {
			return err
		}

		// send request to prove delegation no longer exists.
		k.ICQKeeper.MakeRequest(
			ctx,
			zone.ConnectionId,
			zone.ChainId,
			"store/staking/key",
			data,
			sdk.NewInt(-1),
			types.ModuleName,
			"delegation",
			0,
		)
	}

	k.SetZone(ctx, zone)

	return nil
}

func (k *Keeper) UpdateDelegationRecordForAddress(ctx sdk.Context, delegatorAddress string, validatorAddress string, amount sdk.Coin, zone *types.Zone, absolute bool) error {
	delegation, found := k.GetDelegation(ctx, zone, delegatorAddress, validatorAddress)

	if !found {
		k.Logger(ctx).Info("Adding delegation tuple", "delegator", delegatorAddress, "validator", validatorAddress, "amount", amount.Amount)
		delegation = types.NewDelegation(delegatorAddress, validatorAddress, amount)
	} else if !delegation.Amount.Equal(amount.Amount.ToDec()) {
		oldAmount := delegation.Amount
		if !absolute {
			delegation.Amount = delegation.Amount.Add(amount)
		} else {
			delegation.Amount = amount
		}
		k.Logger(ctx).Info("Updating delegation tuple amount", "delegator", delegatorAddress, "validator", validatorAddress, "old_amount", oldAmount, "inbound_amount", amount.Amount, "new_amount", delegation.Amount, "abs", absolute)
	}
	k.SetDelegation(ctx, zone, delegation)
	if err := k.EmitValsetRequery(ctx, zone.ConnectionId, zone.ChainId); err != nil {
		return err
	}
	return nil
}

func (k *Keeper) HandleWithdrawRewards(ctx sdk.Context, msg sdk.Msg) error {
	withdrawalMsg, ok := msg.(*distrtypes.MsgWithdrawDelegatorReward)
	if !ok {
		k.Logger(ctx).Error("unable to cast source message to MsgWithdrawDelegatorReward")
		return fmt.Errorf("unable to cast source message to MsgWithdrawDelegatorReward")
	}

	zone, err := k.GetZoneFromContext(ctx)
	if err != nil {
		err = fmt.Errorf("4: %w", err)
		k.Logger(ctx).Error(err.Error())
		return err
	}
	// decrement withdrawal waitgroup
	if withdrawalMsg.DelegatorAddress != zone.PerformanceAddress.Address {
		zone.WithdrawalWaitgroup--
		k.SetZone(ctx, zone)
	}
	k.Logger(ctx).Info("Received MsgWithdrawDelegatorReward acknowledgement", "wg", zone.WithdrawalWaitgroup, "delegator", withdrawalMsg.DelegatorAddress)
	switch zone.WithdrawalWaitgroup {
	case 0:
		// interface assertion
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
			"distributerewards",
			0,
		)
		return nil
	default:
		return nil
	}
}

func DistributeRewardsFromWithdrawAccount(k Keeper, ctx sdk.Context, args []byte, query queryTypes.Query) error {
	zone, found := k.GetZone(ctx, query.ChainId)
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

	dust, msgs := k.prepareRewardsDistributionMsgs(zone, rewards)

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

func (k *Keeper) updateRedemptionRate(ctx sdk.Context, zone types.Zone, epochRewards sdk.Coin) {
	ratio := k.GetDelegatedAmount(ctx, &zone).Add(epochRewards).Amount.ToDec().Quo(k.BankKeeper.GetSupply(ctx, zone.LocalDenom).Amount.ToDec())
	k.Logger(ctx).Info("Epochly rewards", "coins", epochRewards)
	k.Logger(ctx).Info("Last redemption rate", "rate", zone.LastRedemptionRate)
	k.Logger(ctx).Info("Current redemption rate", "rate", zone.RedemptionRate)
	k.Logger(ctx).Info("New redemption rate", "rate", ratio, "supply", k.BankKeeper.GetSupply(ctx, zone.LocalDenom).Amount.ToDec(), "lv", k.GetDelegatedAmount(ctx, &zone).Add(epochRewards).Amount.ToDec())

	zone.LastRedemptionRate = zone.RedemptionRate
	zone.RedemptionRate = ratio
	k.SetZone(ctx, &zone)
}

func (k *Keeper) prepareRewardsDistributionMsgs(zone types.Zone, rewards sdk.Coin) (sdk.Int, []sdk.Msg) {
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

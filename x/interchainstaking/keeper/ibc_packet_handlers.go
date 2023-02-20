package keeper

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"cosmossdk.io/math"
	"github.com/golang/protobuf/proto" //nolint:staticcheck

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	icatypes "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v5/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v5/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v5/modules/core/04-channel/types"

	"github.com/cosmos/cosmos-sdk/telemetry"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	lsmstakingtypes "github.com/iqlusioninc/liquidity-staking-module/x/staking/types"

	"github.com/ingenuity-build/quicksilver/utils"
	queryTypes "github.com/ingenuity-build/quicksilver/x/interchainquery/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

const (
	transferPort = "transfer"
	withdrawal   = "withdrawal"
)

type TypedMsg struct {
	Msg  sdk.Msg
	Type string
}

func DeserializeCosmosTxTyped(cdc codec.BinaryCodec, data []byte) ([]TypedMsg, error) {
	var cosmosTx icatypes.CosmosTx
	if err := cdc.Unmarshal(data, &cosmosTx); err != nil {
		return nil, err
	}

	msgs := make([]TypedMsg, len(cosmosTx.Messages))

	for i, any := range cosmosTx.Messages {
		var msg sdk.Msg

		err := cdc.UnpackAny(any, &msg)
		if err != nil {
			return nil, err
		}

		msgs[i] = TypedMsg{Msg: msg, Type: any.TypeUrl}

	}

	return msgs, nil
}

func (k *Keeper) HandleAcknowledgement(ctx sdk.Context, packet channeltypes.Packet, acknowledgement []byte) error {
	ack := channeltypes.Acknowledgement_Result{}
	err := json.Unmarshal(acknowledgement, &ack)
	txMsgData := &sdk.TxMsgData{}
	var success bool
	if err != nil {
		k.Logger(ctx).Error("unable to unmarshal acknowledgement", "error", err, "data", acknowledgement)
		return err
	}
	if reflect.DeepEqual(ack, channeltypes.Acknowledgement_Result{}) {
		ackErr := channeltypes.Acknowledgement_Error{}
		err := json.Unmarshal(acknowledgement, &ackErr)
		if err != nil {
			k.Logger(ctx).Error("unable to unmarshal acknowledgement error", "error", err, "data", acknowledgement)
			return err
		}

		k.Logger(ctx).Error("received an acknowledgement error", "error", err, "remote_err", ackErr, "data", acknowledgement)
		defer telemetry.IncrCounter(1, types.ModuleName, "ica_acknowledgement_errors")
		success = false
	} else {
		defer telemetry.IncrCounter(1, types.ModuleName, "ica_acknowledgement_success")

		err = proto.Unmarshal(ack.Result, txMsgData)
		if err != nil {
			k.Logger(ctx).Error("unable to unmarshal acknowledgement", "error", err, "ack", ack.Result)
			return err
		}
		success = true
	}

	var packetData icatypes.InterchainAccountPacketData
	err = icatypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &packetData)
	if err != nil {
		k.Logger(ctx).Error("unable to unmarshal acknowledgement packet data", "error", err, "data", packetData)
		return err
	}
	if reflect.DeepEqual(packetData, icatypes.InterchainAccountPacketData{}) {
		return errors.New("unable to unmarshal packet data; got empty JSON object")
	}

	msgs, err := DeserializeCosmosTxTyped(k.cdc, packetData.Data)
	if err != nil {
		k.Logger(ctx).Error("unable to decode messages", "err", err)
		return err
	}

	for msgIndex, msg := range msgs {
		// use msgData for v0.45 and below and msgResponse for v0.46+
		//nolint:staticcheck // SA1019 ignore this!
		var msgData *sdk.MsgData
		var msgResponse []byte
		if len(txMsgData.MsgResponses) > 0 {
			msgResponse = txMsgData.MsgResponses[msgIndex].GetValue()
		} else if len(txMsgData.Data) > 0 {
			msgData = txMsgData.Data[msgIndex]
		}

		src := msg.Msg

		switch msg.Type {
		case "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward":
			// TODO: if this fails, it's okay. log but continue.
			if !success {
				return nil
			}
			k.Logger(ctx).Info("Rewards withdrawn")
			if err := k.HandleWithdrawRewards(ctx, src); err != nil {
				return err
			}
			continue
		case "/cosmos.staking.v1beta1.MsgRedeemTokensforShares":
			// TODO: handle this before LSM
			if !success {
				return nil
			}
			response := lsmstakingtypes.MsgRedeemTokensforSharesResponse{}

			if msgResponse != nil {
				err = proto.Unmarshal(msgResponse, &response)
				if err != nil {
					k.Logger(ctx).Error("unable to unpack MsgRedeemTokensforShares response", "error", err)
					return err
				}
			} else {
				err := proto.Unmarshal(msgData.Data, &response)
				if err != nil {
					k.Logger(ctx).Error("unable to unmarshal MsgRedeemTokensforShares response", "error", err)
					return err
				}
			}
			k.Logger(ctx).Info("Tokens redeemed for shares", "response", response)
			// we should update delegation records here.
			if err := k.HandleRedeemTokens(ctx, src, response.Amount); err != nil {
				return err
			}
			continue
		case "/cosmos.staking.v1beta1.MsgTokenizeShares":
			// TODO: handle this before LSM
			if !success {
				return nil
			}
			response := lsmstakingtypes.MsgTokenizeSharesResponse{}
			if msgResponse != nil {
				err = proto.Unmarshal(msgResponse, &response)
				if err != nil {
					k.Logger(ctx).Error("unable to unpack MsgTokenizeShares response", "error", err)
					return err
				}
			} else {
				err := proto.Unmarshal(msgData.Data, &response)
				if err != nil {
					k.Logger(ctx).Error("unable to unmarshal MsgTokenizeShares response", "error", err)
					return err
				}
			}
			k.Logger(ctx).Info("Shares tokenized", "response", response)
			if err := k.HandleTokenizedShares(ctx, src, response.Amount, packetData.Memo); err != nil {
				return err
			}
			continue
		case "/cosmos.staking.v1beta1.MsgDelegate":
			// TODO: can we safely ignore this?
			if !success {
				return nil
			}
			response := stakingtypes.MsgDelegateResponse{}
			if msgResponse != nil {
				err = proto.Unmarshal(msgResponse, &response)
				if err != nil {
					k.Logger(ctx).Error("unable to unpack MsgDelegate response", "error", err)
					return err
				}
			} else {
				err := proto.Unmarshal(msgData.Data, &response)
				if err != nil {
					k.Logger(ctx).Error("unable to unmarshal MsgDelegate response", "error", err)
					return err
				}
			}
			k.Logger(ctx).Info("Delegated", "response", response)
			// we should update delegation records here.
			if err := k.HandleDelegate(ctx, src, packetData.Memo); err != nil {
				return err
			}
			continue
		case "/cosmos.staking.v1beta1.MsgBeginRedelegate":
			if success {
				response := stakingtypes.MsgBeginRedelegateResponse{}
				if msgResponse != nil {
					err = proto.Unmarshal(msgResponse, &response)
					if err != nil {
						k.Logger(ctx).Error("unable to unpack MsgBeginRedelegate response", "error", err)
						return err
					}
				} else {
					err := proto.Unmarshal(msgData.Data, &response)
					if err != nil {
						k.Logger(ctx).Error("unable to unmarshal MsgBeginRedelegate response", "error", err)
						return err
					}
				}
				k.Logger(ctx).Info("Redelegation initiated", "response", response)
				if err := k.HandleBeginRedelegate(ctx, src, response.CompletionTime, packetData.Memo); err != nil {
					return err
				}
			} else {
				if err := k.HandleFailedBeginRedelegate(ctx, src, packetData.Memo); err != nil {
					return err
				}
			}
			continue
		case "/cosmos.staking.v1beta1.MsgUndelegate":
			if success {
				response := stakingtypes.MsgUndelegateResponse{}
				if msgResponse != nil {
					err = proto.Unmarshal(msgResponse, &response)
					if err != nil {
						k.Logger(ctx).Error("unable to unpack MsgUndelegate response", "error", err)
						return err
					}
				} else {
					err := proto.Unmarshal(msgData.Data, &response)
					if err != nil {
						k.Logger(ctx).Error("unable to unmarshal MsgUndelegate response", "error", err)
						return err
					}
				}
				k.Logger(ctx).Info("Undelegation started", "response", response)
				if err := k.HandleUndelegate(ctx, src, response.CompletionTime, packetData.Memo); err != nil {
					return err
				}
			} else {
				if err := k.HandleFailedUndelegate(ctx, src, packetData.Memo); err != nil {
					return err
				}
			}
			continue

		case "/cosmos.bank.v1beta1.MsgSend":
			if !success {
				// TODO: handle this.
				return nil
			}
			response := banktypes.MsgSendResponse{}
			if msgResponse != nil {
				err = proto.Unmarshal(msgResponse, &response)
				if err != nil {
					k.Logger(ctx).Error("unable to unpack MsgSend response", "error", err)
					return err
				}
			} else {
				err := proto.Unmarshal(msgData.Data, &response)
				if err != nil {
					k.Logger(ctx).Error("unable to unmarshal MsgSend response", "error", err)
					return err
				}
			}
			k.Logger(ctx).Info("Funds Transferred", "response", response)
			// check tokenTransfers - if end user unescrow and burn txs
			if err := k.HandleCompleteSend(ctx, src, packetData.Memo); err != nil {
				return err
			}
		case "/cosmos.distribution.v1beta1.MsgSetWithdrawAddress":
			if !success {
				// safely ignore this, as we'll try again anyway.
				return nil
			}
			response := distrtypes.MsgSetWithdrawAddressResponse{}
			if msgResponse != nil {
				err = proto.Unmarshal(msgResponse, &response)
				if err != nil {
					k.Logger(ctx).Error("unable to unpack MsgSetWithdrawAddress response", "error", err)
					return err
				}
			} else {
				err := proto.Unmarshal(msgData.Data, &response)
				if err != nil {
					k.Logger(ctx).Error("unable to unmarshal MsgSetWithdrawAddress response", "error", err)
					return err
				}
			}
			k.Logger(ctx).Info("Withdraw Address Updated", "response", response)
			if err := k.HandleUpdatedWithdrawAddress(ctx, src); err != nil {
				return err
			}
		case "/ibc.applications.transfer.v1.MsgTransfer":
			// this should be okay to fail; we'll pick it up next time around.
			if !success {
				return nil
			}
			response := ibctransfertypes.MsgTransferResponse{}
			if msgResponse != nil {
				err = proto.Unmarshal(msgResponse, &response)
				if err != nil {
					k.Logger(ctx).Error("unable to unpack MsgTransfer response", "error", err)
					return err
				}
			} else {
				err := proto.Unmarshal(msgData.Data, &response)
				if err != nil {
					k.Logger(ctx).Error("unable to unmarshal MsgTransfer response", "error", err)
					return err
				}
			}
			k.Logger(ctx).Info("MsgTranfer acknowledgement received")
			if err := k.HandleMsgTransfer(ctx, src); err != nil {
				return err
			}
		default:
			k.Logger(ctx).Error("unhandled acknowledgement packet", "type", reflect.TypeOf(src).Name())
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
		return errors.New("unable to cast source message to MsgTransfer")
	}

	// check if destination is interchainstaking module account (spoiler: it was)
	if sMsg.Receiver != k.AccountKeeper.GetModuleAddress(types.ModuleName).String() {
		k.Logger(ctx).Error("msgTransfer to unknown account!")
		return errors.New("unexpected recipient")
	}

	return k.HandleDistributeFeesFromModuleAccount(ctx)
}

func (k *Keeper) HandleDistributeFeesFromModuleAccount(ctx sdk.Context) error {
	// what do we have in the account?
	balance := k.BankKeeper.GetAllBalances(ctx, k.AccountKeeper.GetModuleAddress(types.ModuleName))
	k.Logger(ctx).Info("distributing collected fees to stakers", "amount", balance)
	return k.BankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, authtypes.FeeCollectorName, balance) // Fee collector name needs to be passed in to keeper constructor.
}

func (k *Keeper) HandleCompleteSend(ctx sdk.Context, msg sdk.Msg, memo string) error {
	k.Logger(ctx).Info("Received MsgSend acknowledgement")
	// first, type assertion. we should have banktypes.MsgSend
	sMsg, ok := msg.(*banktypes.MsgSend)
	if !ok {
		err := errors.New("unable to cast source message to MsgSend")
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
		// Target here is the DelegationAddresses.
		return k.handleRewardsDelegation(ctx, *zone, sMsg)
	case zone.IsDelegateAddress(sMsg.FromAddress):
		return k.HandleWithdrawForUser(ctx, zone, sMsg, memo)
	case zone.IsDelegateAddress(sMsg.ToAddress) && zone.DepositAddress.Address == sMsg.FromAddress:
		return k.handleSendToDelegate(ctx, zone, sMsg, memo)
	default:
		err = errors.New("unexpected completed send")
		k.Logger(ctx).Error(err.Error())
		return err
	}
}

func (k *Keeper) handleRewardsDelegation(ctx sdk.Context, zone types.Zone, msg *banktypes.MsgSend) error {
	return k.handleSendToDelegate(ctx, &zone, msg, "rewards")
}

func (k *Keeper) handleSendToDelegate(ctx sdk.Context, zone *types.Zone, msg *banktypes.MsgSend, memo string) error {
	var msgs []sdk.Msg
	for _, coin := range msg.Amount {
		if coin.Denom == zone.BaseDenom {
			allocations := k.DeterminePlanForDelegation(ctx, zone, msg.Amount)
			msgs = append(msgs, k.PrepareDelegationMessagesForCoins(ctx, zone, allocations)...)
		} else {
			msgs = append(msgs, k.PrepareDelegationMessagesForShares(ctx, zone, msg.Amount)...)
		}
	}

	k.Logger(ctx).Info("messages to send", "messages", msgs)

	return k.SubmitTx(ctx, msgs, zone.DelegationAddress, memo)
}

// withdraw for user will check that the msgSend we have successfully executed matches an existing withdrawal record.
// on a match (recipient = msg.ToAddress + amount + status == SEND), we mark the record as complete.
// if no other withdrawal records exist for this triple (i.e. no further withdrawal from this delegator account for this user (i.e. different validator))
// then burn the withdrawal_record's burn_amount.
func (k *Keeper) HandleWithdrawForUser(ctx sdk.Context, zone *types.Zone, msg *banktypes.MsgSend, memo string) error {
	var err error
	var withdrawalRecord types.WithdrawalRecord

	withdrawalRecord, found := k.GetWithdrawalRecord(ctx, zone.ChainId, memo, WithdrawStatusSend)

	if !found {
		return errors.New("no matching withdrawal record found")
	}

	// case 1: total amount - native unbonding
	// this statement is ridiculous, but currently calling coins.Equals against coins with different denoms panics; which is pretty useless.
	if len(withdrawalRecord.Amount) == 1 && len(msg.Amount) == 1 && msg.Amount[0].Denom == withdrawalRecord.Amount[0].Denom && withdrawalRecord.Amount.IsEqual(msg.Amount) {
		k.Logger(ctx).Info("found matching withdrawal; marking as completed")
		k.UpdateWithdrawalRecordStatus(ctx, &withdrawalRecord, WithdrawStatusCompleted)
		if err = k.BankKeeper.BurnCoins(ctx, types.EscrowModuleAccount, sdk.NewCoins(withdrawalRecord.BurnAmount)); err != nil {
			// if we can't burn the coins, fail.
			return err
		}
		k.SetWithdrawalRecord(ctx, withdrawalRecord)
		k.Logger(ctx).Info("burned coins post-withdrawal", "coins", withdrawalRecord.BurnAmount)
	} else {

		// case 2: per validator amounts - LSM unbonding

		dlist := make(map[int]struct{})
		for i, dist := range withdrawalRecord.Distribution {
			if msg.Amount[0].Amount.Equal(sdk.NewIntFromUint64(dist.Amount)) { // check valoper here too?
				dlist[i] = struct{}{}
				// matched amount
				if len(withdrawalRecord.Distribution) == len(dlist) {
					// we just removed the last element
					k.Logger(ctx).Info("found matching withdrawal; marking as completed")
					k.UpdateWithdrawalRecordStatus(ctx, &withdrawalRecord, WithdrawStatusCompleted)
					if err = k.BankKeeper.BurnCoins(ctx, types.EscrowModuleAccount, sdk.NewCoins(withdrawalRecord.BurnAmount)); err != nil {
						// if we can't burn the coins, fail.
						return err
					}
					k.SetWithdrawalRecord(ctx, withdrawalRecord)
					k.Logger(ctx).Info("burned coins post-withdrawal", "coins", withdrawalRecord.BurnAmount)
				}
				break
			}
		}

		if len(dlist) > 0 {
			newDist := make([]*types.Distribution, 0)
			i := 0
			for idx := range withdrawalRecord.Distribution {
				if _, delete := dlist[idx]; !delete {
					newDist = append(newDist, withdrawalRecord.Distribution[idx])
				}
				i++
			}
			k.Logger(ctx).Info("found matching withdrawal; awaiting additional messages")
			withdrawalRecord.Distribution = newDist
			k.SetWithdrawalRecord(ctx, withdrawalRecord)
		}
	}

	return k.EmitValsetRequery(ctx, zone.ConnectionId, zone.ChainId)
}

func (k *Keeper) GCCompletedRedelegations(ctx sdk.Context) error {
	var err error

	k.IterateRedelegationRecords(ctx, func(idx int64, key []byte, redelegation types.RedelegationRecord) bool {
		// if the redelegation completion time was in the past AND is not 0000-00-00T00:00:00Z, then delete it.
		if ctx.BlockTime().After(redelegation.CompletionTime) && !redelegation.CompletionTime.Equal(time.Time{}) {
			k.Logger(ctx).Info("garbage collecting completed redelegations", "key", key, "completion", redelegation.CompletionTime)
			k.DeleteRedelegationRecordByKey(ctx, append(types.KeyPrefixRedelegationRecord, key...))
		}
		return false
	})

	return err
}

func (k *Keeper) HandleMaturedUnbondings(ctx sdk.Context, zone *types.Zone) error {
	var err error

	k.IterateZoneStatusWithdrawalRecords(ctx, zone.ChainId, WithdrawStatusUnbond, func(idx int64, withdrawal types.WithdrawalRecord) bool {
		if ctx.BlockTime().After(withdrawal.CompletionTime) && !withdrawal.CompletionTime.Equal(time.Time{}) { // completion date has passed.
			k.Logger(ctx).Info("found completed unbonding")
			sendMsg := &banktypes.MsgSend{FromAddress: zone.DelegationAddress.GetAddress(), ToAddress: withdrawal.Recipient, Amount: sdk.Coins{withdrawal.Amount[0]}}
			err = k.SubmitTx(ctx, []sdk.Msg{sendMsg}, zone.DelegationAddress, withdrawal.Txhash)
			if err != nil {
				k.Logger(ctx).Error("error", err)
				return true
			}
			k.Logger(ctx).Info("sending funds", "for", withdrawal.Delegator, "delegate_account", zone.DelegationAddress.GetAddress(), "to", withdrawal.Recipient, "amount", withdrawal.Amount)
			k.UpdateWithdrawalRecordStatus(ctx, &withdrawal, WithdrawStatusSend)
		}
		return false
	})
	return err
}

func (k *Keeper) HandleTokenizedShares(ctx sdk.Context, msg sdk.Msg, sharesAmount sdk.Coin, memo string) error {
	var err error
	k.Logger(ctx).Info("received MsgTokenizeShares acknowledgement")
	// first, type assertion. we should have stakingtypes.MsgTokenizeShares
	tsMsg, ok := msg.(*lsmstakingtypes.MsgTokenizeShares)
	if !ok {
		k.Logger(ctx).Error("unable to cast source message to MsgTokenizeShares")
		return errors.New("unable to cast source message to MsgTokenizeShares")
	}

	zone := k.GetZoneForDelegateAccount(ctx, tsMsg.DelegatorAddress)

	withdrawalRecord, found := k.GetWithdrawalRecord(ctx, zone.ChainId, memo, WithdrawStatusTokenize)

	if !found {
		return errors.New("no matching withdrawal record found")
	}

	for _, dist := range withdrawalRecord.Distribution {
		if sharesAmount.Equal(dist.Amount) {
			withdrawalRecord.Amount.Add(sharesAmount)
			// matched amount
			if len(withdrawalRecord.Distribution) == len(withdrawalRecord.Amount) {
				// we just added the last tokens
				k.Logger(ctx).Info("Found matching withdrawal; marking for send")
				k.DeleteWithdrawalRecord(ctx, zone.ChainId, memo, WithdrawStatusTokenize)
				withdrawalRecord.Status = WithdrawStatusSend
				sendMsg := &banktypes.MsgSend{FromAddress: zone.DelegationAddress.Address, ToAddress: withdrawalRecord.Recipient, Amount: withdrawalRecord.Amount}
				err = k.SubmitTx(ctx, []sdk.Msg{sendMsg}, zone.DelegationAddress, memo)
				if err != nil {
					return err
				}
			} else {
				k.Logger(ctx).Info("Found matching withdrawal; awaiting additional messages")
			}
			k.SetWithdrawalRecord(ctx, withdrawalRecord)
			break
		}
	}
	return nil
}

func (k *Keeper) HandleBeginRedelegate(ctx sdk.Context, msg sdk.Msg, completion time.Time, memo string) error {
	parts := strings.Split(memo, "/")
	if len(parts) != 2 || parts[0] != "rebalance" {
		return errors.New("unexpected epoch rebalance memo format")
	}

	epochNumber, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return errors.New("unexpected epoch rebalance memo format (2)")
	}

	k.Logger(ctx).Info("Received MsgBeginRedelegate acknowledgement")
	// first, type assertion. we should have stakingtypes.MsgBeginRedelegate
	redelegateMsg, ok := msg.(*stakingtypes.MsgBeginRedelegate)
	if !ok {
		return errors.New("unable to unmarshal MsgBeginRedelegate")
	}
	zone := k.GetZoneForDelegateAccount(ctx, redelegateMsg.DelegatorAddress)
	record, found := k.GetRedelegationRecord(ctx, zone.ChainId, redelegateMsg.ValidatorSrcAddress, redelegateMsg.ValidatorDstAddress, epochNumber)
	if !found {
		k.Logger(ctx).Error("unable to find redelegation record", "chain", zone.ChainId, "source", redelegateMsg.ValidatorSrcAddress, "dst", redelegateMsg.ValidatorDstAddress, "epoch", epochNumber)
		return fmt.Errorf("unable to find redelegation record for chain %s, src: %s, dst: %s, at epoch %d", zone.ChainId, redelegateMsg.ValidatorSrcAddress, redelegateMsg.ValidatorDstAddress, epochNumber)
	}
	k.Logger(ctx).Info("updating redelegation record with completion time", "completion", completion)
	record.CompletionTime = completion
	k.SetRedelegationRecord(ctx, record)

	delegation, found := k.GetDelegation(ctx, zone, redelegateMsg.DelegatorAddress, redelegateMsg.ValidatorDstAddress)
	if !found {
		k.Logger(ctx).Error("unable to find delegation record", "chain", zone.ChainId, "source", redelegateMsg.ValidatorSrcAddress, "dst", redelegateMsg.ValidatorDstAddress, "epoch", epochNumber)
		return fmt.Errorf("unable to find delegation record for chain %s, src: %s, dst: %s, at epoch %d", zone.ChainId, redelegateMsg.ValidatorSrcAddress, redelegateMsg.ValidatorDstAddress, epochNumber)
	}
	delegation.RedelegationEnd = completion.Unix() // this field should be a timestamp, but lets avoid unnecessary state changes.
	k.SetDelegation(ctx, zone, delegation)
	return nil
}

func (k *Keeper) HandleFailedBeginRedelegate(ctx sdk.Context, msg sdk.Msg, memo string) error {
	parts := strings.Split(memo, "/")
	if len(parts) != 2 || parts[0] != "rebalance" {
		return errors.New("unexpected epoch rebalance memo format")
	}

	epochNumber, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return errors.New("unexpected epoch rebalance memo format (2)")
	}

	k.Logger(ctx).Error("Received MsgBeginRedelegate acknowledgement error")
	// first, type assertion. we should have stakingtypes.MsgBeginRedelegate
	redelegateMsg, ok := msg.(*stakingtypes.MsgBeginRedelegate)
	if !ok {
		return errors.New("unable to unmarshal MsgBeginRedelegate")
	}
	zone := k.GetZoneForDelegateAccount(ctx, redelegateMsg.DelegatorAddress)
	k.DeleteRedelegationRecord(ctx, zone.ChainId, redelegateMsg.ValidatorSrcAddress, redelegateMsg.ValidatorDstAddress, epochNumber)
	k.Logger(ctx).Error("Cleaning up redelegation record")
	return nil
}

func (k *Keeper) HandleUndelegate(ctx sdk.Context, msg sdk.Msg, completion time.Time, memo string) error {
	k.Logger(ctx).Info("Received MsgUndelegate acknowledgement")
	// first, type assertion. we should have stakingtypes.MsgUndelegate
	undelegateMsg, ok := msg.(*stakingtypes.MsgUndelegate)
	if !ok {
		k.Logger(ctx).Error("unable to cast source message to MsgUndelegate")
		return errors.New("unable to cast source message to MsgUndelegate")
	}
	memoParts := strings.Split(memo, "/")
	if len(memoParts) != 2 || memoParts[0] != withdrawal {
		return errors.New("unexpected memo form")
	}

	epochNumber, err := strconv.ParseInt(memoParts[1], 10, 64)
	if err != nil {
		return err
	}
	zone := k.GetZoneForDelegateAccount(ctx, undelegateMsg.DelegatorAddress)

	ubr, found := k.GetUnbondingRecord(ctx, zone.ChainId, undelegateMsg.ValidatorAddress, epochNumber)
	if !found {
		return fmt.Errorf("unbonding record for %s not found for epoch %d", undelegateMsg.ValidatorAddress, epochNumber)
	}

	for _, hash := range ubr.RelatedTxhash {
		k.Logger(ctx).Info("MsgUndelegate", "del", undelegateMsg.DelegatorAddress, "val", undelegateMsg.ValidatorAddress, "hash", hash, "chain", zone.ChainId)

		record, found := k.GetWithdrawalRecord(ctx, zone.ChainId, hash, WithdrawStatusUnbond)
		if !found {
			return fmt.Errorf("unable to lookup withdrawal record; chain: %s, hash: %s", zone.ChainId, hash)
		}
		if completion.After(record.CompletionTime) {
			record.CompletionTime = completion
		}
		k.Logger(ctx).Info("withdrawal record to save", "rcd", record)
		k.UpdateWithdrawalRecordStatus(ctx, &record, WithdrawStatusUnbond)
	}

	delAddr, err := utils.AccAddressFromBech32(undelegateMsg.DelegatorAddress, "")
	if err != nil {
		return err
	}
	valAddr, err := utils.ValAddressFromBech32(undelegateMsg.ValidatorAddress, "")
	if err != nil {
		return err
	}

	data := stakingtypes.GetDelegationKey(delAddr, valAddr)

	// send request to update delegation record for undelegated del/val tuple.
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

	return nil
}

func (k *Keeper) HandleFailedUndelegate(ctx sdk.Context, msg sdk.Msg, memo string) error {
	parts := strings.Split(memo, "/")
	if len(parts) != 2 || parts[0] != withdrawal {
		return errors.New("unexpected epoch undelegate memo format")
	}

	epochNumber, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return errors.New("unexpected epoch undelegate memo format (2)")
	}

	k.Logger(ctx).Error("Received MsgUndelegate acknowledgement error")
	// first, type assertion. we should have stakingtypes.MsgBeginRedelegate
	undelegateMsg, ok := msg.(*stakingtypes.MsgUndelegate)
	if !ok {
		return errors.New("unable to unmarshal MsgUndelegate")
	}

	zone := k.GetZoneForDelegateAccount(ctx, undelegateMsg.DelegatorAddress)
	ubr, found := k.GetUnbondingRecord(ctx, zone.ChainId, undelegateMsg.ValidatorAddress, epochNumber)
	if !found {
		return fmt.Errorf("cannot find unbonding record for %s/%s/%d", zone.ChainId, undelegateMsg.ValidatorAddress, epochNumber)
	}

	for _, hash := range ubr.RelatedTxhash {
		wdr, found := k.GetWithdrawalRecord(ctx, zone.ChainId, hash, WithdrawStatusUnbond)
		if !found {
			return fmt.Errorf("cannot find withdrawal record for %s/%s", zone.ChainId, hash)
		}
		if len(wdr.Distribution) == 1 {
			// sanity check
			if wdr.Distribution[0].Valoper != ubr.Validator {
				return fmt.Errorf("unable to requeue withdrawal record for failed unbonding; expected %s, got %s", ubr.Validator, wdr.Distribution[0].Valoper)
			}
			wdr.Distribution = nil
			wdr.Requeued = true
			k.UpdateWithdrawalRecordStatus(ctx, &wdr, WithdrawStatusQueued)
		} else {
			// remove this validator from distribution; amend amounts; requeue.
			newDistribution := make([]*types.Distribution, 0)
			relatedAmount := uint64(0)
			for _, dist := range wdr.Distribution {
				if dist.Valoper != ubr.Validator {
					newDistribution = append(newDistribution, dist)
				} else {
					relatedAmount = dist.Amount
				}
			}
			wdr.Distribution = newDistribution
			amount := wdr.Amount.AmountOf(zone.BaseDenom)
			wdr.Amount = wdr.Amount.Sub(sdk.NewCoin(zone.BaseDenom, sdk.NewIntFromUint64(relatedAmount)))
			rr := sdk.NewDecFromInt(wdr.BurnAmount.Amount).Quo(sdk.NewDecFromInt(amount))
			relatedQAsset := sdk.NewDec(int64(relatedAmount)).Mul(rr).TruncateInt()
			wdr.BurnAmount = wdr.BurnAmount.SubAmount(relatedQAsset)
			k.SetWithdrawalRecord(ctx, wdr)
			// create a new record with the failed amount
			newWdr := types.WithdrawalRecord{
				ChainId:      zone.ChainId,
				Delegator:    wdr.Delegator,
				Recipient:    wdr.Recipient,
				Distribution: nil,
				Amount:       sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewIntFromUint64(relatedAmount))),
				BurnAmount:   sdk.NewCoin(zone.LocalDenom, relatedQAsset),
				Txhash:       fmt.Sprintf("%064d", k.GetNextWithdrawalRecordSequence(ctx)),
				Status:       WithdrawStatusQueued,
				Requeued:     true,
			}
			k.SetWithdrawalRecord(ctx, newWdr)
		}
	}

	k.DeleteUnbondingRecord(ctx, zone.ChainId, undelegateMsg.ValidatorAddress, epochNumber)
	k.Logger(ctx).Error("Cleaning up redelegation record")
	return nil
}

func (k *Keeper) HandleRedeemTokens(ctx sdk.Context, msg sdk.Msg, amount sdk.Coin) error {
	k.Logger(ctx).Info("Received MsgRedeemTokensforShares acknowledgement")
	// first, type assertion. we should have stakingtypes.MsgRedeemTokensforShares
	redeemMsg, ok := msg.(*lsmstakingtypes.MsgRedeemTokensforShares)
	if !ok {
		k.Logger(ctx).Error("unable to cast source message to MsgRedeemTokensforShares")
		return errors.New("unable to cast source message to MsgRedeemTokensforShares")
	}
	validatorAddress, err := k.GetValidatorForToken(ctx, redeemMsg.DelegatorAddress, redeemMsg.Amount)
	if err != nil {
		return err
	}
	zone := k.GetZoneForDelegateAccount(ctx, redeemMsg.DelegatorAddress)

	return k.UpdateDelegationRecordForAddress(ctx, redeemMsg.DelegatorAddress, validatorAddress, amount, zone, false)
}

func (k *Keeper) HandleDelegate(ctx sdk.Context, msg sdk.Msg, memo string) error {
	k.Logger(ctx).Info("Received MsgDelegate acknowledgement")
	// first, type assertion. we should have stakingtypes.MsgDelegate
	delegateMsg, ok := msg.(*stakingtypes.MsgDelegate)
	if !ok {
		k.Logger(ctx).Error("unable to cast source message to MsgDelegate")
		return errors.New("unable to cast source message to MsgDelegate")
	}
	zone := k.GetZoneForDelegateAccount(ctx, delegateMsg.DelegatorAddress)
	if zone == nil {
		// most likely a performance account...
		if zone := k.GetZoneForPerformanceAccount(ctx, delegateMsg.DelegatorAddress); zone != nil {
			return nil
		}
		return fmt.Errorf("unable to find zone for address %s", delegateMsg.DelegatorAddress)
	}

	if memo != "rewards" {
		receipt, found := k.GetReceipt(ctx, GetReceiptKey(zone.ChainId, memo))
		if !found {
			return fmt.Errorf("unable to find receipt for hash %s", memo)
		}
		t := ctx.BlockTime()
		receipt.Completed = &t
		k.SetReceipt(ctx, receipt)
	}
	return k.UpdateDelegationRecordForAddress(ctx, delegateMsg.DelegatorAddress, delegateMsg.ValidatorAddress, delegateMsg.Amount, zone, false)
}

func (k *Keeper) HandleUpdatedWithdrawAddress(ctx sdk.Context, msg sdk.Msg) error {
	k.Logger(ctx).Info("Received MsgSetWithdrawAddress acknowledgement")
	// first, type assertion. we should have distrtypes.MsgSetWithdrawAddress
	original, ok := msg.(*distrtypes.MsgSetWithdrawAddress)
	if !ok {
		k.Logger(ctx).Error("unable to cast source message to MsgSetWithdrawAddress")
		return errors.New("unable to cast source message to MsgSetWithdrawAddress")
	}
	zone := k.GetZoneForDelegateAccount(ctx, original.DelegatorAddress)
	if zone == nil {
		zone = k.GetZoneForPerformanceAccount(ctx, original.DelegatorAddress)
		if zone == nil {
			zone = k.GetZoneForDepositAccount(ctx, original.DelegatorAddress)
			if zone == nil {
				return errors.New("unable to find zone")
			}
			if err := zone.DepositAddress.SetWithdrawalAddress(original.WithdrawAddress); err != nil {
				return err
			}
		}
		if err := zone.PerformanceAddress.SetWithdrawalAddress(original.WithdrawAddress); err != nil {
			return err
		}
	} else {
		if err := zone.DelegationAddress.SetWithdrawalAddress(original.WithdrawAddress); err != nil {
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

func (k *Keeper) UpdateDelegationRecordsForAddress(ctx sdk.Context, zone types.Zone, delegatorAddress string, args []byte) error {
	var response stakingtypes.QueryDelegatorDelegationsResponse
	err := k.cdc.Unmarshal(args, &response)
	if err != nil {
		return err
	}
	k.Logger(ctx).Info("Delegation query response", "response", response)
	_, delAddr, err := bech32.DecodeAndConvert(delegatorAddress)
	if err != nil {
		return err
	}
	delegatorDelegations := k.GetDelegatorDelegations(ctx, &zone, delAddr)

	delMap := make(map[string]types.Delegation, len(delegatorDelegations))
	for _, del := range delegatorDelegations {
		delMap[del.ValidatorAddress] = del
	}
	for _, delegationRecord := range response.DelegationResponses {

		_, valAddr, err := bech32.DecodeAndConvert(delegationRecord.Delegation.ValidatorAddress)
		if err != nil {
			return err
		}
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
			// zone.DelegationAddress.IncrementBalanceWaitgroup() // does this get decremented?
		}

		if ok {
			delete(delMap, delegationRecord.Delegation.ValidatorAddress)
		}
	}
	for _, existingValAddr := range utils.Keys(delMap) {
		existingDelegation := delMap[existingValAddr]
		_, valAddr, err := bech32.DecodeAndConvert(existingDelegation.ValidatorAddress)
		if err != nil {
			return err
		}
		data := stakingtypes.GetDelegationKey(delAddr, valAddr)

		// send request to prove delegation no longer exists. If the response is nil (i.e. no delegation), then
		// the delegation record is removed by the callback.
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

	return nil
}

func (k *Keeper) UpdateDelegationRecordForAddress(ctx sdk.Context, delegatorAddress string, validatorAddress string, amount sdk.Coin, zone *types.Zone, absolute bool) error {
	delegation, found := k.GetDelegation(ctx, zone, delegatorAddress, validatorAddress)

	if !found {
		k.Logger(ctx).Info("Adding delegation tuple", "delegator", delegatorAddress, "validator", validatorAddress, "amount", amount.Amount)
		delegation = types.NewDelegation(delegatorAddress, validatorAddress, amount)
	} else {
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
		return errors.New("unable to cast source message to MsgWithdrawDelegatorReward")
	}

	zone, err := k.GetZoneFromContext(ctx)
	if err != nil {
		err = fmt.Errorf("4: %w", err)
		k.Logger(ctx).Error(err.Error())
		return err
	}
	// decrement withdrawal waitgroup
	// We are specifically looking for protocol delegator:validator pairs
	// and must not decrement the waitgroup for the performance address as it
	// is not part of the waitgroup set. It is a special delegator address that
	// operates outside the delegator set, its purpose is to track validator
	// performance only.
	if withdrawalMsg.DelegatorAddress != zone.PerformanceAddress.Address {
		zone.WithdrawalWaitgroup--
		k.Logger(ctx).Info("Decremented waitgroup", "wg", zone.WithdrawalWaitgroup)
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
		k.Logger(ctx).Info("Distributing rewards")
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

	baseDenomFee := sdk.NewDecFromInt(baseDenomAmount).
		Mul(k.GetCommissionRate(ctx)).
		TruncateInt()

	// prepare rewards distribution
	rewards := sdk.NewCoin(zone.BaseDenom, baseDenomAmount.Sub(baseDenomFee))

	var msgs []sdk.Msg
	msgs = append(msgs, k.prepareRewardsDistributionMsgs(zone, rewards.Amount))

	// multiDenomFee is the balance of withdrawal account minus the redelegated rewards.
	multiDenomFee := withdrawBalance.Balances.Sub(sdk.Coins{rewards}...)

	var remotePort string
	var remoteChannel string
	k.IBCKeeper.ChannelKeeper.IterateChannels(ctx, func(channel channeltypes.IdentifiedChannel) bool {
		if channel.ConnectionHops[0] == zone.ConnectionId && channel.PortId == transferPort && channel.State == channeltypes.OPEN {
			remoteChannel = channel.Counterparty.ChannelId
			remotePort = channel.Counterparty.PortId
			return true
		}
		return false
	})

	if remotePort == "" {
		return errors.New("unable to find remote transfer connection")
	}

	for _, coin := range multiDenomFee.Sort() {
		msgs = append(
			msgs,
			&ibctransfertypes.MsgTransfer{
				SourcePort:       remotePort,
				SourceChannel:    remoteChannel,
				Token:            coin,
				Sender:           zone.WithdrawalAddress.Address,
				Receiver:         k.AccountKeeper.GetModuleAddress(types.ModuleName).String(),
				TimeoutTimestamp: uint64(ctx.BlockTime().UnixNano() + 6*time.Hour.Nanoseconds()),
				TimeoutHeight:    clienttypes.Height{RevisionNumber: 0, RevisionHeight: 0},
			},
		)
	}

	// update redemption rate
	k.UpdateRedemptionRate(ctx, zone, rewards.Amount)

	// send tx
	return k.SubmitTx(ctx, msgs, zone.WithdrawalAddress, "")
}

func (k *Keeper) prepareRewardsDistributionMsgs(zone types.Zone, rewards math.Int) sdk.Msg {
	return &banktypes.MsgSend{
		FromAddress: zone.WithdrawalAddress.GetAddress(),
		ToAddress:   zone.DelegationAddress.GetAddress(),
		Amount:      sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, rewards)),
	}
}

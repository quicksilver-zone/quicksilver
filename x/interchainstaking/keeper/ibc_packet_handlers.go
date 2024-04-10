package keeper

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/proto" // nolint:staticcheck

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	icatypes "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v5/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v5/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v5/modules/core/04-channel/types"

	"github.com/quicksilver-zone/quicksilver/utils"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	cmtypes "github.com/quicksilver-zone/quicksilver/x/claimsmanager/types"
	emtypes "github.com/quicksilver-zone/quicksilver/x/eventmanager/types"
	querytypes "github.com/quicksilver-zone/quicksilver/x/interchainquery/types"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	lsmstakingtypes "github.com/quicksilver-zone/quicksilver/x/lsmtypes"
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

func (k *Keeper) HandleAcknowledgement(ctx sdk.Context, packet channeltypes.Packet, acknowledgement []byte, connectionID string) error {
	var (
		ack        channeltypes.Acknowledgement
		success    bool
		txMsgData  sdk.TxMsgData
		packetData icatypes.InterchainAccountPacketData
	)

	err := icatypes.ModuleCdc.UnmarshalJSON(acknowledgement, &ack)
	if err != nil {
		k.Logger(ctx).Error("unable to unmarshal acknowledgement", "error", err, "data", acknowledgement)
		return err
	}

	if !ack.Success() {
		ackErr := ack.GetError()
		k.Logger(ctx).Error("received an acknowledgement error", "remote_err", ackErr, "data", ack.String())
		defer telemetry.IncrCounter(1, types.ModuleName, "ica_acknowledgement_errors")
		success = false
	} else {
		defer telemetry.IncrCounter(1, types.ModuleName, "ica_acknowledgement_success")
		err = proto.Unmarshal(ack.GetResult(), &txMsgData)
		if err != nil {
			k.Logger(ctx).Error("unable to unmarshal acknowledgement", "error", err, "ack", ack.GetResult())
			return err
		}
		success = true
	}

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
		var msgResponse []byte

		// check that the msgResponses slice is at least the length of the current index.
		switch {
		case !success:
			// no-op - there is no msgresponse for a AckErr
		case len(txMsgData.MsgResponses) > msgIndex:
			msgResponse = txMsgData.MsgResponses[msgIndex].GetValue()
		case len(txMsgData.Data) > msgIndex:
			msgResponse = txMsgData.Data[msgIndex].GetData()
		default:
			return fmt.Errorf("could not find msgresponse for index %d", msgIndex)
		}

		switch msg.Type {
		case "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward":
			if !success {
				withdrawalMsg, ok := msg.Msg.(*distrtypes.MsgWithdrawDelegatorReward)
				if !ok {
					return errors.New("unable to unmarshal MsgWithdrawDelegatorReward")
				}
				k.Logger(ctx).Error("failed to withdraw rewards; will try again next epoch", "validator", withdrawalMsg.ValidatorAddress)
				return nil
			}
			k.Logger(ctx).Info("Rewards withdrawn")
			if err := k.HandleWithdrawRewards(ctx, msg.Msg, connectionID); err != nil {
				return err
			}
			continue
		case "/cosmos.staking.v1beta1.MsgRedeemTokensForShares":
			if !success {
				if err := k.HandleFailedRedeemTokens(ctx, msg.Msg, packetData.Memo); err != nil {
					return err
				}
				continue
			}
			response := lsmstakingtypes.MsgRedeemTokensForSharesResponse{}

			err = proto.Unmarshal(msgResponse, &response)
			if err != nil {
				k.Logger(ctx).Error("unable to unmarshal MsgRedeemTokensForShares response", "error", err)
				return err
			}

			k.Logger(ctx).Info("Tokens redeemed for shares", "response", response)
			// we should update delegation records here.
			if err := k.HandleRedeemTokens(ctx, msg.Msg, response.Amount, packetData.Memo, connectionID); err != nil {
				return err
			}
			continue
		case "/cosmos.staking.v1beta1.MsgTokenizeShares":
			if !success {
				// We can safely ignore this, as this can reasonably fail, and we cater for this in the flush logic.
				return nil
			}
			response := lsmstakingtypes.MsgTokenizeSharesResponse{}

			err = proto.Unmarshal(msgResponse, &response)
			if err != nil {
				k.Logger(ctx).Error("unable to unpack MsgTokenizeShares response", "error", err)
				return err
			}

			k.Logger(ctx).Info("Shares tokenized", "response", response)
			if err := k.HandleTokenizedShares(ctx, msg.Msg, response.Amount, packetData.Memo); err != nil {
				return err
			}
			continue
		case "/cosmos.staking.v1beta1.MsgDelegate":
			if !success {
				if err := k.HandleFailedDelegate(ctx, msg.Msg, packetData.Memo); err != nil {
					return err
				}
				continue
			}
			response := stakingtypes.MsgDelegateResponse{}
			err = proto.Unmarshal(msgResponse, &response)
			if err != nil {
				k.Logger(ctx).Error("unable to unpack MsgDelegate response", "error", err)
				return err
			}

			k.Logger(ctx).Info("Delegated", "response", response)
			// we should update delegation records here.
			if err := k.HandleDelegate(ctx, msg.Msg, packetData.Memo); err != nil {
				return err
			}
			continue
		case "/cosmos.staking.v1beta1.MsgBeginRedelegate":
			if success {
				response := stakingtypes.MsgBeginRedelegateResponse{}
				err = proto.Unmarshal(msgResponse, &response)
				k.Logger(ctx).Info("unmarshalling msgResponse", "response", response)
				if err != nil {
					k.Logger(ctx).Error("unable to unpack MsgBeginRedelegate response", "error", err)
					return err
				}

				k.Logger(ctx).Info("Redelegation initiated", "response", response)
				if err := k.HandleBeginRedelegate(ctx, msg.Msg, response.CompletionTime, packetData.Memo); err != nil {
					return err
				}
			} else {
				if err := k.HandleFailedBeginRedelegate(ctx, msg.Msg, packetData.Memo); err != nil {
					return err
				}
			}
			continue
		case "/cosmos.staking.v1beta1.MsgUndelegate":
			if success {
				response := stakingtypes.MsgUndelegateResponse{}
				err = proto.Unmarshal(msgResponse, &response)
				if err != nil {
					k.Logger(ctx).Error("unable to unpack MsgUndelegate response", "error", err)
					return err
				}

				k.Logger(ctx).Info("Undelegation started", "response", response)
				if err := k.HandleUndelegate(ctx, msg.Msg, response.CompletionTime, packetData.Memo); err != nil {
					return err
				}
			} else {
				if err := k.HandleFailedUndelegate(ctx, msg.Msg, packetData.Memo); err != nil {
					return err
				}
			}
			continue

		case "/cosmos.bank.v1beta1.MsgSend":
			if !success {
				if err := k.HandleFailedBankSend(ctx, msg.Msg, packetData.Memo, connectionID); err != nil {
					k.Logger(ctx).Error("unable to handle failed MsgSend", "error", err)
					return err
				}
				continue
			}
			response := banktypes.MsgSendResponse{}
			err = proto.Unmarshal(msgResponse, &response)
			if err != nil {
				k.Logger(ctx).Error("unable to unpack MsgSend response", "error", err)
				return err
			}

			k.Logger(ctx).Info("Funds Transferred", "response", response)
			// check tokenTransfers - if end user unescrow and burn txs
			if err := k.HandleCompleteSend(ctx, msg.Msg, packetData.Memo, connectionID); err != nil {
				return err
			}
		case "/cosmos.distribution.v1beta1.MsgSetWithdrawAddress":
			if !success {
				// safely ignore this, as we'll try again anyway.
				return nil
			}
			response := distrtypes.MsgSetWithdrawAddressResponse{}
			err = proto.Unmarshal(msgResponse, &response)
			if err != nil {
				k.Logger(ctx).Error("unable to unpack MsgSetWithdrawAddress response", "error", err)
				return err
			}

			k.Logger(ctx).Info("Withdraw Address Updated", "response", response)
			if err := k.HandleUpdatedWithdrawAddress(ctx, msg.Msg); err != nil {
				return err
			}
		case "/ibc.applications.transfer.v1.MsgTransfer":
			k.Logger(ctx).Debug("Received MsgTransfer acknowledgement; no action")
			return nil

		default:
			k.Logger(ctx).Error("unhandled acknowledgement packet", "type", reflect.TypeOf(msg.Msg).Name())
		}
	}

	return nil
}

func (*Keeper) HandleTimeout(_ sdk.Context, _ channeltypes.Packet) error {
	return nil
}

// ----------------------------------------------------------------

func (k *Keeper) HandleMsgTransfer(ctx sdk.Context, msg ibctransfertypes.FungibleTokenPacketData, ibcDenom string) error {
	k.Logger(ctx).Info("Received MsgTransfer acknowledgement")
	// first, type assertion. we should have ibctransfertypes.MsgTransfer

	// check if destination is interchainstaking module account (spoiler: it was)
	if msg.Receiver != k.AccountKeeper.GetModuleAddress(types.ModuleName).String() {
		k.Logger(ctx).Error("msgTransfer to unknown account!")
		return errors.New("unexpected recipient")
	}

	receivedAmount, ok := math.NewIntFromString(msg.Amount)
	if !ok {
		return fmt.Errorf("unable to marshal amount into math.Int: %s", msg.Amount)
	}
	receivedCoin := sdk.NewCoin(ibcDenom, receivedAmount)

	zone, found := k.GetZoneForWithdrawalAccount(ctx, msg.Sender)
	if !found {
		return fmt.Errorf("zone not found for withdrawal account %s", msg.Sender)
	}

	if found && msg.Denom != zone.BaseDenom {
		feeAmount := sdk.NewDecFromInt(receivedCoin.Amount).Mul(k.GetCommissionRate(ctx)).TruncateInt()
		rewardCoin := receivedCoin.SubAmount(feeAmount)
		zoneAddress, err := addressutils.AccAddressFromBech32(zone.WithdrawalAddress.Address, "")
		if err != nil {
			return err
		}
		k.Logger(ctx).Info("distributing collected rewards to users", "amount", rewardCoin)
		remaining, err := k.DistributeToClaimants(ctx, zone, zoneAddress, rewardCoin)
		if err != nil {
			return err
		}
		receivedCoin = sdk.NewCoin(receivedCoin.Denom, feeAmount).Add(remaining)
	}

	balance := sdk.NewCoins(receivedCoin)
	k.Logger(ctx).Info("distributing collected fees to stakers", "amount", balance)
	return k.BankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, authtypes.FeeCollectorName, balance)
}

func (k *Keeper) DistributeToClaimants(ctx sdk.Context, zone *types.Zone, zoneAddress sdk.AccAddress, rewardsCoin sdk.Coin) (sdk.Coin, error) {
	var err error
	toDistribute := rewardsCoin.Amount
	supply := k.BankKeeper.GetSupply(ctx, zone.LocalDenom).Amount
	claimTotal := math.ZeroInt()
	k.ClaimsManagerKeeper.IterateLastEpochClaims(ctx, zone.ChainId, func(index int64, data cmtypes.Claim) (stop bool) {
		claimTotal = claimTotal.Add(data.Amount)
		return false
	})

	ratio := math.LegacyOneDec()
	if claimTotal.GT(supply) {
		ratio = math.LegacyNewDecFromInt(supply).Quo(math.LegacyNewDecFromInt(claimTotal))
	}

	k.ClaimsManagerKeeper.IterateLastEpochClaims(ctx, zone.ChainId, func(index int64, data cmtypes.Claim) (stop bool) {
		claimAmount := math.LegacyNewDecFromInt(data.Amount).Mul(ratio).Quo(math.LegacyNewDecFromInt(supply)).Mul(rewardsCoin.Amount.ToLegacyDec()).TruncateInt()
		err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addressutils.MustAccAddressFromBech32(data.UserAddress, ""), sdk.NewCoins(sdk.NewCoin(rewardsCoin.Denom, claimAmount)))
		toDistribute = toDistribute.Sub(claimAmount)
		return err != nil
	})

	if toDistribute.IsNegative() {
		return sdk.Coin{}, fmt.Errorf("unexpected negative value")
	}

	return sdk.NewCoin(rewardsCoin.Denom, toDistribute), err
}

func (k *Keeper) HandleCompleteSend(ctx sdk.Context, msg sdk.Msg, memo string, connectionID string) error {
	k.Logger(ctx).Info("Received MsgSend acknowledgement")
	// first, type assertion. we should have banktypes.MsgSend
	sMsg, ok := msg.(*banktypes.MsgSend)
	if !ok {
		err := errors.New("unable to cast source message to MsgSend")
		k.Logger(ctx).Error(err.Error())
		return err
	}

	// get zone
	zone, err := k.GetZoneFromConnectionID(ctx, connectionID)
	if err != nil {
		err = fmt.Errorf("2: %w", err)
		k.Logger(ctx).Error(err.Error())
		return err
	}

	// checks here are specific to ensure future extensibility;
	switch {
	case zone.IsDelegateAddress(sMsg.ToAddress) && zone.IsWithdrawalAddress(sMsg.FromAddress):
		k.Logger(ctx).Info("delegate account received tokens from withdrawal account; delegating rewards", "amount", sMsg.Amount)
		return k.handleRewardsDelegation(ctx, *zone, sMsg)

	case zone.IsWithdrawalAddress(sMsg.ToAddress):
		k.Logger(ctx).Info("withdrawal account received tokens to disburse", "amount", sMsg.Amount)
		return nil

	case zone.IsDelegateAddress(sMsg.ToAddress) && zone.IsDepositAddress(sMsg.FromAddress):
		k.Logger(ctx).Info("delegate account received tokens from deposit account; delegating deposit", "amount", sMsg.Amount, "memo", memo)
		_, err := k.handleSendToDelegate(ctx, zone, sMsg, memo)
		return err

	case zone.IsDelegateAddress(sMsg.FromAddress):
		k.Logger(ctx).Info("delegate account send tokens; handling withdrawal", "amount", sMsg.Amount, "memo", memo)
		return k.HandleWithdrawForUser(ctx, zone, sMsg, memo)

	case zone.IsDepositAddress(sMsg.FromAddress) && memo == "refund":
		k.Logger(ctx).Info("unable to process deposit, returning funds to sender", "recipient", sMsg.ToAddress, "amount", sMsg.Amount, "memo", memo)
		return nil

	default:
		err = fmt.Errorf("unexpected completed send (2) from %s to %s (amount: %s)", sMsg.FromAddress, sMsg.ToAddress, sMsg.Amount)
		k.Logger(ctx).Error(err.Error())
		return err
	}
}

func (k *Keeper) handleRewardsDelegation(ctx sdk.Context, zone types.Zone, msg *banktypes.MsgSend) error {
	_, err := k.handleSendToDelegate(ctx, &zone, msg, "rewards")
	return err
}

func (k *Keeper) handleSendToDelegate(ctx sdk.Context, zone *types.Zone, msg *banktypes.MsgSend, memo string) (int, error) {
	var msgs []sdk.Msg
	for _, coin := range msg.Amount {
		if coin.Denom == zone.BaseDenom {
			allocations, err := k.DeterminePlanForDelegation(ctx, zone, msg.Amount)
			if err != nil {
				return 0, err
			}
			msgs = append(msgs, k.PrepareDelegationMessagesForCoins(zone, allocations, isBatchOrRewards(memo))...)
		} else {
			msgs = append(msgs, k.PrepareDelegationMessagesForShares(zone, msg.Amount)...)
		}
	}

	k.Logger(ctx).Info("messages to send", "messages", msgs)

	return len(msgs), k.SubmitTx(ctx, msgs, zone.DelegationAddress, memo, zone.MessagesPerTx)
}

func isBatchOrRewards(memo string) bool {
	if memo == "rewards" {
		return true
	}
	return strings.HasPrefix(memo, "batch")
}

// HandleWithdrawForUser handles withdraw for user will check that the msgSend we have successfully executed matches an existing withdrawal record.
// on a match (recipient = msg.ToAddress + amount + status == SEND), we mark the record as complete.
// if no other withdrawal records exist for this triple (i.e. no further withdrawal from this delegator account for this user (i.e. different validator))
// then burn the withdrawal_record's burn_amount.
func (k *Keeper) HandleWithdrawForUser(ctx sdk.Context, zone *types.Zone, msg *banktypes.MsgSend, memo string) error {
	txHash, err := types.ParseTxMsgMemo(memo, types.MsgTypeUnbondSend)
	if err != nil {
		return err
	}

	withdrawalRecord, found := k.GetWithdrawalRecord(ctx, zone.ChainId, txHash, types.WithdrawStatusSend)
	if !found {
		return errors.New("no matching withdrawal record found")
	}

	// case 1: total amount - native unbonding
	// this statement is ridiculous, but currently calling coins.Equals against coins with different denoms panics; which is pretty useless.
	if len(withdrawalRecord.Amount) == 1 && len(msg.Amount) == 1 && msg.Amount[0].Denom == withdrawalRecord.Amount[0].Denom && withdrawalRecord.Amount.IsEqual(msg.Amount) {
		k.Logger(ctx).Info("found matching withdrawal; marking as completed")
		k.UpdateWithdrawalRecordStatus(ctx, &withdrawalRecord, types.WithdrawStatusCompleted)
		if err := k.BankKeeper.BurnCoins(ctx, types.EscrowModuleAccount, sdk.NewCoins(withdrawalRecord.BurnAmount)); err != nil {
			// if we can't burn the coins, fail.
			return err
		}
		k.Logger(ctx).Info("burned coins post-withdrawal", "coins", withdrawalRecord.BurnAmount)
	} else {

		// case 2: per validator amounts - LSM unbonding

		dlist := make(map[int]struct{})
		for i, dist := range withdrawalRecord.Distribution {
			if msg.Amount[0].Amount.Equal(dist.Amount) { // check valoper here too?
				dlist[i] = struct{}{}
				// matched amount
				if len(withdrawalRecord.Distribution) == len(dlist) {
					// we just removed the last element
					k.Logger(ctx).Info("found matching withdrawal; marking as completed")
					k.UpdateWithdrawalRecordStatus(ctx, &withdrawalRecord, types.WithdrawStatusCompleted)
					if err := k.BankKeeper.BurnCoins(ctx, types.EscrowModuleAccount, sdk.NewCoins(withdrawalRecord.BurnAmount)); err != nil {
						// if we can't burn the coins, fail.
						return err
					}
					k.Logger(ctx).Info("burned coins post-withdrawal", "coins", withdrawalRecord.BurnAmount)
				}
				break
			}
		}

		if len(dlist) > 0 {
			newDist := make([]*types.Distribution, 0)
			for idx := range withdrawalRecord.Distribution {
				if _, remove := dlist[idx]; !remove {
					newDist = append(newDist, withdrawalRecord.Distribution[idx])
				}
			}
			k.Logger(ctx).Info("found matching withdrawal; awaiting additional messages")
			withdrawalRecord.Distribution = newDist
			err = k.SetWithdrawalRecord(ctx, withdrawalRecord)
			if err != nil {
				return err
			}
		}
	}

	period := int64(k.GetParam(ctx, types.KeyValidatorSetInterval))
	query := stakingtypes.QueryValidatorsRequest{}
	return k.EmitValSetQuery(ctx, zone.ConnectionId, zone.ChainId, query, math.NewInt(period))
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
	k.IterateZoneStatusWithdrawalRecords(ctx, zone.ChainId, types.WithdrawStatusUnbond, func(idx int64, withdrawal types.WithdrawalRecord) bool {
		if ctx.BlockTime().After(withdrawal.CompletionTime) && withdrawal.Acknowledged { // completion date has passed.
			k.Logger(ctx).Info("found completed unbonding")
			sendMsg := &banktypes.MsgSend{FromAddress: zone.DelegationAddress.GetAddress(), ToAddress: withdrawal.Recipient, Amount: sdk.Coins{withdrawal.Amount[0]}}
			err := k.SubmitTx(ctx, []sdk.Msg{sendMsg}, zone.DelegationAddress, types.TxUnbondSendMemo(withdrawal.Txhash), zone.MessagesPerTx)

			if err != nil {
				k.Logger(ctx).Error("error submitting transaction - requeue withdrawal", "error", err)

				// do not update status and increment completion time
				withdrawal.DelayCompletion(ctx, types.DefaultWithdrawalRequeueDelay)
				err = k.SetWithdrawalRecord(ctx, withdrawal)
				if err != nil {
					k.Logger(ctx).Error("error updating withdrawal record", "error", err)
				}

			} else {
				k.Logger(ctx).Info("sending funds", "for", withdrawal.Delegator, "delegate_account", zone.DelegationAddress.GetAddress(), "to", withdrawal.Recipient, "amount", withdrawal.Amount)
				k.UpdateWithdrawalRecordStatus(ctx, &withdrawal, types.WithdrawStatusSend)
			}
		}
		return false
	})
	return nil
}

func (k *Keeper) GetInflightUnbondingAmount(ctx sdk.Context, zone *types.Zone) sdk.Coin {
	outCoin := sdk.NewCoin(zone.BaseDenom, sdk.ZeroInt())
	k.IterateZoneWithdrawalRecords(ctx, zone.ChainId, func(idx int64, withdrawal types.WithdrawalRecord) bool {
		if (withdrawal.Status == types.WithdrawStatusUnbond && ctx.BlockTime().After(withdrawal.CompletionTime) && withdrawal.Acknowledged) || // status unbond, completion has pass
			withdrawal.Status == types.WithdrawStatusSend { // already in state send.
			outCoin = outCoin.Add(withdrawal.Amount[0])
		}
		return false
	})
	return outCoin
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

	zone, found := k.GetZoneForDelegateAccount(ctx, tsMsg.DelegatorAddress)
	if !found {
		return fmt.Errorf("zone for delegate account %s not found", tsMsg.DelegatorAddress)
	}
	withdrawalRecord, found := k.GetWithdrawalRecord(ctx, zone.ChainId, memo, types.WithdrawStatusTokenize)

	if !found {
		return errors.New("no matching withdrawal record found")
	}

	for _, dist := range withdrawalRecord.Distribution {
		if equalLsmCoin(dist.Valoper, dist.Amount, sharesAmount) {
			withdrawalRecord.Amount = withdrawalRecord.Amount.Add(sharesAmount)
			break
		}
	}

	err = k.SetWithdrawalRecord(ctx, withdrawalRecord)
	if err != nil {
		return err
	}

	if len(withdrawalRecord.Distribution) != len(withdrawalRecord.Amount) {
		k.Logger(ctx).Info(fmt.Sprintf("Found matching withdrawal (%d/%d); awaiting additional messages", len(withdrawalRecord.Amount), len(withdrawalRecord.Distribution)))
	} else {
		k.Logger(ctx).Info("Found matching withdrawal; marking for send")
		k.DeleteWithdrawalRecord(ctx, zone.ChainId, memo, types.WithdrawStatusTokenize)
		withdrawalRecord.Status = types.WithdrawStatusSend
		err = k.SetWithdrawalRecord(ctx, withdrawalRecord)
		if err != nil {
			return err
		}
		sendMsg := &banktypes.MsgSend{FromAddress: zone.DelegationAddress.Address, ToAddress: withdrawalRecord.Recipient, Amount: withdrawalRecord.Amount}
		err = k.SubmitTx(ctx, []sdk.Msg{sendMsg}, zone.DelegationAddress, memo, zone.MessagesPerTx)
	}
	return err
}

func (k *Keeper) HandleBeginRedelegate(ctx sdk.Context, msg sdk.Msg, completion time.Time, memo string) error {
	epochNumber, err := types.ParseEpochMsgMemo(memo, types.MsgTypeRebalance)
	if err != nil {
		return err
	}

	k.Logger(ctx).Info("Received MsgBeginRedelegate acknowledgement")
	// first, type assertion. we should have stakingtypes.MsgBeginRedelegate
	redelegateMsg, ok := msg.(*stakingtypes.MsgBeginRedelegate)
	if !ok {
		return errors.New("unable to unmarshal MsgBeginRedelegate")
	}

	zone, found := k.GetZoneForDelegateAccount(ctx, redelegateMsg.DelegatorAddress)
	if !found {
		return fmt.Errorf("zone for delegate account %s not found", redelegateMsg.DelegatorAddress)
	}

	if completion.IsZero() {
		// a zero completion time can only happen when the validator is unbonded; this means the redelegation has _already_ completed and can be removed.
		k.DeleteRedelegationRecord(ctx, zone.ChainId, redelegateMsg.ValidatorSrcAddress, redelegateMsg.ValidatorDstAddress, epochNumber)
	} else {
		record, found := k.GetRedelegationRecord(ctx, zone.ChainId, redelegateMsg.ValidatorSrcAddress, redelegateMsg.ValidatorDstAddress, epochNumber)
		if !found {
			// it is possible that the record was cleaned up if there was a long delay in processing acknowledgements.
			// just create a new one
			record = types.RedelegationRecord{
				ChainId:        zone.ChainId,
				EpochNumber:    epochNumber,
				Source:         redelegateMsg.ValidatorSrcAddress,
				Destination:    redelegateMsg.ValidatorDstAddress,
				Amount:         redelegateMsg.Amount.Amount,
				CompletionTime: completion,
			}
		}

		k.Logger(ctx).Info("updating redelegation record with completion time", "completion", completion)
		record.CompletionTime = completion
		k.SetRedelegationRecord(ctx, record)
	}

	tgtDelegation, found := k.GetDelegation(ctx, zone.ChainId, redelegateMsg.DelegatorAddress, redelegateMsg.ValidatorDstAddress)
	if !found {
		tgtDelegation = types.NewDelegation(redelegateMsg.DelegatorAddress, redelegateMsg.ValidatorDstAddress, redelegateMsg.Amount)
	} else {
		tgtDelegation.Amount = tgtDelegation.Amount.Add(redelegateMsg.Amount)
	}
	// RedelegationEnd is used to determine whether the delegation is 'locked' for transient redelegations.
	tgtDelegation.RedelegationEnd = completion.Unix() // this field should be a timestamp, but let's avoid unnecessary state changes.
	k.SetDelegation(ctx, zone.ChainId, tgtDelegation)

	delAddr, err := addressutils.AccAddressFromBech32(redelegateMsg.DelegatorAddress, zone.AccountPrefix)
	if err != nil {
		return err
	}
	valAddr, err := addressutils.ValAddressFromBech32(redelegateMsg.ValidatorDstAddress, zone.AccountPrefix+"valoper")
	if err != nil {
		return err
	}
	data := stakingtypes.GetDelegationKey(delAddr, valAddr)

	// send request to update delegation record for target del/val tuple.
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

	srcDelegation, found := k.GetDelegation(ctx, zone.ChainId, redelegateMsg.DelegatorAddress, redelegateMsg.ValidatorSrcAddress)
	if !found {
		k.Logger(ctx).Error("unable to find delegation record", "chain", zone.ChainId, "source", redelegateMsg.ValidatorSrcAddress, "dst", redelegateMsg.ValidatorDstAddress, "epoch_number", epochNumber)
		return fmt.Errorf("unable to find delegation record for chain %s, src: %s, dst: %s, at epoch %d", zone.ChainId, redelegateMsg.ValidatorSrcAddress, redelegateMsg.ValidatorDstAddress, epochNumber)
	}
	srcDelegation.Amount, err = srcDelegation.Amount.SafeSub(redelegateMsg.Amount)
	if err != nil {
		if strings.Contains(err.Error(), "negative coin amount") {
			// we received a negative srcDelegation. Obviously this cannot happen, but we can get a crossed re/un/delegation, all which fetch absolute values.
			k.Logger(ctx).Error("possible race condition; unable to sub redelegation amount. requerying delegation anyway")
		} else {
			// we got some other, unrecoverable err
			return err
		}
	} else {
		k.SetDelegation(ctx, zone.ChainId, srcDelegation)
	}

	valAddr, err = addressutils.ValAddressFromBech32(redelegateMsg.ValidatorDstAddress, zone.AccountPrefix+"valoper")
	if err != nil {
		return err
	}
	data = stakingtypes.GetDelegationKey(delAddr, valAddr)

	// send request to update delegation record for src del/val tuple.
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

func (k *Keeper) HandleFailedBeginRedelegate(ctx sdk.Context, msg sdk.Msg, memo string) error {
	epochNumber, err := types.ParseEpochMsgMemo(memo, types.MsgTypeRebalance)
	if err != nil {
		return err
	}

	k.Logger(ctx).Error("received MsgBeginRedelegate acknowledgement error")
	// first, type assertion. we should have stakingtypes.MsgBeginRedelegate
	redelegateMsg, ok := msg.(*stakingtypes.MsgBeginRedelegate)
	if !ok {
		return errors.New("unable to unmarshal MsgBeginRedelegate")
	}
	zone, found := k.GetZoneForDelegateAccount(ctx, redelegateMsg.DelegatorAddress)
	if !found {
		return fmt.Errorf("zone for delegate account %s not found", redelegateMsg.DelegatorAddress)
	}
	k.DeleteRedelegationRecord(ctx, zone.ChainId, redelegateMsg.ValidatorSrcAddress, redelegateMsg.ValidatorDstAddress, epochNumber)
	k.Logger(ctx).Info("cleaning up redelegation record")
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

	epochNumber, err := types.ParseEpochMsgMemo(memo, types.MsgTypeWithdrawal)
	if err != nil {
		return err
	}

	zone, found := k.GetZoneForDelegateAccount(ctx, undelegateMsg.DelegatorAddress)
	if !found {
		return fmt.Errorf("zone for delegate account %s not found", undelegateMsg.DelegatorAddress)
	}

	if err := zone.DecrementWithdrawalWaitgroup(k.Logger(ctx), 1, "unbonding message ack"); err != nil {
		// given that there _could_ be a backlog of message, we don't want to bail here, else they will remain undeliverable.
		k.Logger(ctx).Error(err.Error())
	}
	ubr, found := k.GetUnbondingRecord(ctx, zone.ChainId, undelegateMsg.ValidatorAddress, epochNumber)
	if !found {
		return fmt.Errorf("unbonding record for %s not found for epoch %d", undelegateMsg.ValidatorAddress, epochNumber)
	}

	for _, hash := range ubr.RelatedTxhash {
		k.Logger(ctx).Info("MsgUndelegate", "del", undelegateMsg.DelegatorAddress, "val", undelegateMsg.ValidatorAddress, "hash", hash, "chain", zone.ChainId)

		record, found := k.GetWithdrawalRecord(ctx, zone.ChainId, hash, types.WithdrawStatusUnbond)
		if !found {
			return fmt.Errorf("unable to lookup withdrawal record; chain: %s, hash: %s", zone.ChainId, hash)
		}

		record.Acknowledged = true

		if completion.After(record.CompletionTime) {
			record.CompletionTime = completion
		}
		k.Logger(ctx).Info("withdrawal record to save", "rcd", record)
		k.UpdateWithdrawalRecordStatus(ctx, &record, types.WithdrawStatusUnbond)
	}

	delAddr, err := addressutils.AccAddressFromBech32(undelegateMsg.DelegatorAddress, "")
	if err != nil {
		return err
	}
	valAddr, err := addressutils.ValAddressFromBech32(undelegateMsg.ValidatorAddress, "")
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
		"delegation_epoch",
		0,
	)

	if err = zone.IncrementWithdrawalWaitgroup(k.Logger(ctx), 1, "unbonding message ack emit delegation_epoch query"); err != nil {
		return err
	}
	k.SetZone(ctx, zone)

	return nil
}

func (k *Keeper) HandleFailedBankSend(ctx sdk.Context, msg sdk.Msg, memo string, connectionID string) error {
	sMsg, ok := msg.(*banktypes.MsgSend)
	if !ok {
		err := errors.New("unable to cast source message to MsgSend")
		k.Logger(ctx).Error(err.Error())
		return err
	}

	// get zone
	zone, err := k.GetZoneFromConnectionID(ctx, connectionID)
	if err != nil {
		k.Logger(ctx).Error(err.Error())
		return err
	}

	// checks here are specific to ensure future extensibility;
	switch {
	case zone.IsDelegateAddress(sMsg.ToAddress) && zone.IsWithdrawalAddress(sMsg.FromAddress):
		// MsgSend from Withdrawal account to delegate account was not completed. We can ignore this.
		k.Logger(ctx).Info("MsgSend to delegate account from withdrawal account failed", "amount", sMsg.Amount)
	case zone.IsWithdrawalAddress(sMsg.ToAddress):
		k.Logger(ctx).Info("MsgSend to withdrawal account for disbursal failed", "amount", sMsg.Amount)
	case zone.IsDelegateAddress(sMsg.ToAddress) && zone.IsDepositAddress(sMsg.FromAddress):
		// MsgSend from deposit account to delegate account for deposit.
		k.Logger(ctx).Error("MsgSend from deposit account to delegate account failed", "amount", sMsg.Amount)
	case zone.IsDelegateAddress(sMsg.FromAddress):
		k.Logger(ctx).Info("MsgSend from delegate account failed; updating withdrawal", "amount", sMsg.Amount, "memo", memo)
		return k.HandleFailedUnbondSend(ctx, sMsg, memo)
	default:
		err = fmt.Errorf("unexpected failed send (1) from %s to %s (amount: %s)", sMsg.FromAddress, sMsg.ToAddress, sMsg.Amount)
		k.Logger(ctx).Error(err.Error())
	}

	return nil
}

func (k *Keeper) HandleFailedUnbondSend(ctx sdk.Context, sendMsg *banktypes.MsgSend, memo string) error {
	txHash, err := types.ParseTxMsgMemo(memo, types.MsgTypeUnbondSend)
	if err != nil {
		return err
	}

	// get chainID for the remote zone using msg addresses (ICA acc)
	chainID, found := k.GetAddressZoneMapping(ctx, sendMsg.FromAddress)
	if !found {
		return fmt.Errorf("unable to find address mapping for address %s: txHash %s", sendMsg.FromAddress, txHash)
	}

	wdr, found := k.GetWithdrawalRecord(ctx, chainID, txHash, types.WithdrawStatusSend)
	if !found {
		return fmt.Errorf("unable to find withdrawal record for %s: txHash %s", sendMsg.ToAddress, txHash)
	}

	// update delayed record with status
	wdr.DelayCompletion(ctx, types.DefaultWithdrawalRequeueDelay)
	k.UpdateWithdrawalRecordStatus(ctx, &wdr, types.WithdrawStatusUnbond)

	return nil
}

func (k *Keeper) HandleFailedUndelegate(ctx sdk.Context, msg sdk.Msg, memo string) error {
	epochNumber, err := types.ParseEpochMsgMemo(memo, types.MsgTypeWithdrawal)
	if err != nil {
		return err
	}

	k.Logger(ctx).Error("received MsgUndelegate acknowledgement error")
	// first, type assertion. we should have stakingtypes.MsgBeginRedelegate
	undelegateMsg, ok := msg.(*stakingtypes.MsgUndelegate)
	if !ok {
		return errors.New("unable to unmarshal MsgUndelegate")
	}

	zone, found := k.GetZoneForDelegateAccount(ctx, undelegateMsg.DelegatorAddress)
	if !found {
		return fmt.Errorf("zone for delegate account %s not found", undelegateMsg.DelegatorAddress)
	}
	ubr, found := k.GetUnbondingRecord(ctx, zone.ChainId, undelegateMsg.ValidatorAddress, epochNumber)
	if !found {
		return fmt.Errorf("cannot find unbonding record for %s/%s/%d", zone.ChainId, undelegateMsg.ValidatorAddress, epochNumber)
	}

	for _, hash := range ubr.RelatedTxhash {
		wdr, found := k.GetWithdrawalRecord(ctx, zone.ChainId, hash, types.WithdrawStatusUnbond)
		if !found {
			return fmt.Errorf("cannot find withdrawal record for %s/%s", zone.ChainId, hash)
		}
		// if multi val then:
		// - remove this validator from distribution
		// - related amount = amount from this val
		// - determine RR paid
		// - mult RR by related amount, sub this from burn amount
		// - save old record
		// - create new record for unhandled burn amount
		newDistribution := make([]*types.Distribution, 0)
		relatedAmount := math.ZeroInt()
		for _, dist := range wdr.Distribution {
			if dist.Valoper != ubr.Validator {
				newDistribution = append(newDistribution, dist)
			} else {
				relatedAmount = dist.Amount
			}
		}

		amount := wdr.Amount.AmountOf(zone.BaseDenom)
		rr := sdk.NewDecFromInt(wdr.BurnAmount.Amount).Quo(sdk.NewDecFromInt(amount))
		relatedQAsset := sdk.NewDecFromInt(relatedAmount).Mul(rr).TruncateInt()

		if len(newDistribution) == 0 {
			// if this was the final record, delete the withdrawal record
			k.DeleteWithdrawalRecord(ctx, wdr.ChainId, wdr.Txhash, wdr.Status)
		} else {
			// else update it
			wdr.Distribution = newDistribution
			wdr.Amount = wdr.Amount.Sub(sdk.NewCoin(zone.BaseDenom, relatedAmount))
			wdr.BurnAmount = wdr.BurnAmount.SubAmount(relatedQAsset)
			err = k.SetWithdrawalRecord(ctx, wdr)
			if err != nil {
				return err
			}

		}

		record := k.GetUserChainRequeuedWithdrawalRecord(ctx, zone.ChainId, wdr.Delegator)
		if record.Txhash == "" {
			// create a new record with the failed amount
			record = types.WithdrawalRecord{
				ChainId:      zone.ChainId,
				Delegator:    wdr.Delegator,
				Recipient:    wdr.Recipient,
				Distribution: nil,
				BurnAmount:   sdk.NewCoin(zone.LocalDenom, relatedQAsset),
				Txhash:       fmt.Sprintf("%064d", k.GetNextWithdrawalRecordSequence(ctx)),
				Status:       types.WithdrawStatusQueued,
				Requeued:     true,
				EpochNumber:  wdr.EpochNumber,
			}
		} else {
			record.BurnAmount = record.BurnAmount.Add(sdk.NewCoin(zone.LocalDenom, relatedQAsset))
		}
		err = k.SetWithdrawalRecord(ctx, record)
		if err != nil {
			return err
		}
	}

	k.DeleteUnbondingRecord(ctx, zone.ChainId, undelegateMsg.ValidatorAddress, epochNumber)
	k.Logger(ctx).Info("cleaning up unbonding record")
	return nil
}

func (k *Keeper) HandleRedeemTokens(ctx sdk.Context, msg sdk.Msg, amount sdk.Coin, memo string, connectionID string) error {
	k.Logger(ctx).Info("Received MsgRedeemTokensforShares acknowledgement")
	// first, type assertion. we should have stakingtypes.MsgRedeemTokensforShares
	redeemMsg, ok := msg.(*lsmstakingtypes.MsgRedeemTokensForShares)
	if !ok {
		k.Logger(ctx).Error("unable to cast source message to MsgRedeemTokensforShares")
		return errors.New("unable to cast source message to MsgRedeemTokensforShares")
	}
	validatorAddress, err := k.GetValidatorForToken(ctx, redeemMsg.Amount, connectionID)
	if err != nil {
		return err
	}
	zone, found := k.GetZoneForDelegateAccount(ctx, redeemMsg.DelegatorAddress)
	if !found {
		return fmt.Errorf("zone for delegate account %s not found", redeemMsg.DelegatorAddress)
	}

	switch {
	case strings.HasPrefix(memo, "batch"):
		k.Logger(ctx).Debug("batch delegation", "memo", memo, "tx", redeemMsg)
		exclusionTimestampUnix, err := strconv.ParseInt(strings.Split(memo, "/")[1], 10, 64)
		if err != nil {
			return err
		}
		k.Logger(ctx).Debug("outstanding delegations ack-received")
		k.SetReceiptsCompleted(ctx, zone.ChainId, time.Unix(exclusionTimestampUnix, 0), ctx.BlockTime(), redeemMsg.Amount.Denom)
		balance, negative := zone.DelegationAddress.Balance.SafeSub(redeemMsg.Amount)
		if negative {
			k.Logger(ctx).Error("unexpected negative balance; likely due to stale ack")
			return nil
		}
		zone.DelegationAddress.Balance = balance
		k.SetZone(ctx, zone)
		if zone.GetWithdrawalWaitgroup() == 0 {
			k.Logger(ctx).Info("Triggering redemption rate calc after delegation flush")
			if err = k.TriggerRedemptionRate(ctx, zone); err != nil {
				return err
			}
		}

	default:
		receipt, found := k.GetReceipt(ctx, zone.ChainId, memo)
		if !found {
			return fmt.Errorf("unable to find receipt for hash %s", memo)
		}
		t := ctx.BlockTime()
		receipt.Completed = &t
		k.SetReceipt(ctx, receipt)
	}
	return k.UpdateDelegationRecordForAddress(ctx, redeemMsg.DelegatorAddress, validatorAddress, amount, zone, false, false)
}

func (k *Keeper) HandleFailedRedeemTokens(ctx sdk.Context, msg sdk.Msg, memo string) error {
	k.Logger(ctx).Info("Received MsgRedeemTokensForShares failure acknowledgement")
	// first, type assertion. we should have lsmstakingtypes.MsgRedeemTokensForShares
	redeemMsg, ok := msg.(*lsmstakingtypes.MsgRedeemTokensForShares)
	if !ok {
		k.Logger(ctx).Error("unable to cast source message to MsgRedeemTokensForShares")
		return errors.New("unable to cast source message to MsgRedeemTokensForShares")
	}
	zone, found := k.GetZoneForDelegateAccount(ctx, redeemMsg.DelegatorAddress)
	if !found {
		// most likely a performance account...
		if _, found := k.GetZoneForPerformanceAccount(ctx, redeemMsg.DelegatorAddress); !found {
			return nil
		}
		return fmt.Errorf("unable to find zone for address %s", redeemMsg.DelegatorAddress)
	}

	switch {
	case strings.HasPrefix(memo, "batch"):
		k.Logger(ctx).Error("batch token redemption failed", "memo", memo, "tx", redeemMsg)
		if err := zone.DecrementWithdrawalWaitgroup(k.Logger(ctx), uint32(1), "batch token redemption failure ack"); err != nil {
			k.Logger(ctx).Error(err.Error())
			return nil
			// return nil here so we don't reject the incoming tx, but log the error and don't trigger RR update for repeated zero.
		}
		k.Logger(ctx).Info("Decremented waitgroup after failed batch token redemption", "wg", zone.GetWithdrawalWaitgroup())
		k.SetZone(ctx, zone)
		if zone.GetWithdrawalWaitgroup() == 0 {
			k.Logger(ctx).Info("Triggering redemption rate calc after delegation flush")
			if err := k.TriggerRedemptionRate(ctx, zone); err != nil {
				return err
			}
		}

	default:
		// no-op
	}
	return nil
}

func (k *Keeper) HandleDelegate(ctx sdk.Context, msg sdk.Msg, memo string) error {
	k.Logger(ctx).Info("Received MsgDelegate acknowledgement")
	// first, type assertion. we should have stakingtypes.MsgDelegate
	delegateMsg, ok := msg.(*stakingtypes.MsgDelegate)
	if !ok {
		k.Logger(ctx).Error("unable to cast source message to MsgDelegate")
		return errors.New("unable to cast source message to MsgDelegate")
	}
	zone, found := k.GetZoneForDelegateAccount(ctx, delegateMsg.DelegatorAddress)
	if !found {
		// most likely a performance account...
		if _, found := k.GetZoneForPerformanceAccount(ctx, delegateMsg.DelegatorAddress); found {
			return nil
		}
		return fmt.Errorf("unable to find zone for address %s", delegateMsg.DelegatorAddress)
	}
	switch {
	case memo == "rewards":
	case strings.HasPrefix(memo, "batch"):
		k.Logger(ctx).Info("batch delegation", "memo", memo, "tx", delegateMsg)
		exclusionTimestampUnix, err := strconv.ParseInt(strings.Split(memo, "/")[1], 10, 64)
		if err != nil {
			return err
		}
		k.Logger(ctx).Debug("outstanding delegations ack-received")
		k.SetReceiptsCompleted(ctx, zone.ChainId, time.Unix(exclusionTimestampUnix, 0), ctx.BlockTime(), delegateMsg.Amount.Denom)
		balance, negative := zone.DelegationAddress.Balance.SafeSub(delegateMsg.Amount)
		if negative {
			k.Logger(ctx).Error("unexpected negative balance; likely a stale ack")
			return nil
		}
		zone.DelegationAddress.Balance = balance
		if err := zone.DecrementWithdrawalWaitgroup(k.Logger(ctx), uint32(1), "batch/reward delegation success ack"); err != nil {
			k.Logger(ctx).Error(err.Error())
			return nil
			// return nil here so we don't reject the incoming tx, but log the error and don't trigger RR update for repeated zero.
		}
		k.SetZone(ctx, zone)
		if zone.GetWithdrawalWaitgroup() == 0 {
			k.Logger(ctx).Info("Triggering redemption rate calc after delegation flush")
			if err := k.TriggerRedemptionRate(ctx, zone); err != nil {
				return err
			}
		}
	default:
		receipt, found := k.GetReceipt(ctx, zone.ChainId, memo)
		if !found {
			return fmt.Errorf("unable to find receipt for hash %s", memo)
		}
		t := ctx.BlockTime()
		receipt.Completed = &t
		k.SetReceipt(ctx, receipt)

	}

	return k.UpdateDelegationRecordForAddress(ctx, delegateMsg.DelegatorAddress, delegateMsg.ValidatorAddress, delegateMsg.Amount, zone, false, false)
}

func (k *Keeper) HandleFailedDelegate(ctx sdk.Context, msg sdk.Msg, memo string) error {
	k.Logger(ctx).Info("Received MsgDelegate failure acknowledgement")
	// first, type assertion. we should have stakingtypes.MsgDelegate
	delegateMsg, ok := msg.(*stakingtypes.MsgDelegate)
	if !ok {
		k.Logger(ctx).Error("unable to cast source message to MsgDelegate")
		return errors.New("unable to cast source message to MsgDelegate")
	}
	zone, found := k.GetZoneForDelegateAccount(ctx, delegateMsg.DelegatorAddress)
	if !found {
		// most likely a performance account...
		if _, found := k.GetZoneForPerformanceAccount(ctx, delegateMsg.DelegatorAddress); found {
			return nil
		}
		return fmt.Errorf("unable to find zone for address %s", delegateMsg.DelegatorAddress)
	}

	switch {
	case strings.HasPrefix(memo, "batch"):
		k.Logger(ctx).Error("batch delegation failed", "memo", memo, "tx", delegateMsg)
		if err := zone.DecrementWithdrawalWaitgroup(k.Logger(ctx), 1, "batch delegation failed ack"); err != nil {
			k.Logger(ctx).Error(err.Error())
			return nil
			// return nil here so we don't reject the incoming tx, but log the error and don't trigger RR update for repeated zero.
		}
		k.SetZone(ctx, zone)
		if zone.GetWithdrawalWaitgroup() == 0 {
			k.Logger(ctx).Info("Triggering redemption rate calc after delegation flush")
			if err := k.TriggerRedemptionRate(ctx, zone); err != nil {
				return err
			}
		}

	default:
		// no-op
	}
	return nil
}

func (k *Keeper) HandleUpdatedWithdrawAddress(ctx sdk.Context, msg sdk.Msg) error {
	k.Logger(ctx).Info("Received MsgSetWithdrawAddress acknowledgement")
	// first, type assertion. we should have distrtypes.MsgSetWithdrawAddress
	original, ok := msg.(*distrtypes.MsgSetWithdrawAddress)
	if !ok {
		k.Logger(ctx).Error("unable to cast source message to MsgSetWithdrawAddress")
		return errors.New("unable to cast source message to MsgSetWithdrawAddress")
	}
	zone, found := k.GetZoneForDelegateAccount(ctx, original.DelegatorAddress)
	if !found {
		zone, found = k.GetZoneForPerformanceAccount(ctx, original.DelegatorAddress)
		if !found {
			zone, found = k.GetZoneForDepositAccount(ctx, original.DelegatorAddress)
			if !found {
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

func (k *Keeper) GetValidatorForToken(ctx sdk.Context, amount sdk.Coin, connectionID string) (string, error) {
	zone, err := k.GetZoneFromConnectionID(ctx, connectionID)
	if err != nil {
		err = fmt.Errorf("3: %w", err)
		k.Logger(ctx).Error(err.Error())
		return "", err
	}

	for _, val := range k.GetValidatorAddresses(ctx, zone.ChainId) {
		if strings.HasPrefix(amount.Denom, val) {
			// match!
			return val, nil
		}
	}

	return "", fmt.Errorf("unable to find validator for token %s", amount.Denom)
}

// UpdateDelegationRecordsForAddress accepts a QueryDelegatorDelegationsResponse and for new, or changed delegation records will
// trigger an ICQ request for that record. If this was triggered by an epoch, the withdrawal waitgroup should be decremented once,
// (for the incoming message) and incremented for each outgoing message.
func (k *Keeper) UpdateDelegationRecordsForAddress(ctx sdk.Context, zone types.Zone, delegatorAddress string, args []byte, isEpoch bool) error {
	var response stakingtypes.QueryDelegatorDelegationsResponse
	err := k.cdc.Unmarshal(args, &response)
	if err != nil {
		return err
	}
	k.Logger(ctx).Info("Delegation query response", "isEpoch", isEpoch, "response", response)
	_, delAddr, err := bech32.DecodeAndConvert(delegatorAddress)
	if err != nil {
		return err
	}
	delegatorDelegations := k.GetDelegatorDelegations(ctx, zone.ChainId, delAddr)

	delMap := make(map[string]types.Delegation, len(delegatorDelegations))
	for _, del := range delegatorDelegations {
		delMap[del.ValidatorAddress] = del
	}

	cb := "delegation"
	if isEpoch {
		if err := zone.DecrementWithdrawalWaitgroup(k.Logger(ctx), 1, "delegations_epoch callback succeeded"); err != nil {
			k.Logger(ctx).Error(err.Error())
			// don't return here, catch and squash err.
		}
		cb = "delegation_epoch"
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
				cb,
				0,
			)
			if isEpoch {
				err = zone.IncrementWithdrawalWaitgroup(k.Logger(ctx), 1, fmt.Sprintf("delegation callback emit %s query", cb))
				if err != nil {
					return err
				}
			}
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
			cb,
			0,
		)
		if isEpoch {
			err = zone.IncrementWithdrawalWaitgroup(k.Logger(ctx), 1, fmt.Sprintf("delegations callback emit %s query", cb))
			if err != nil {
				return err
			}
		}
	}

	if isEpoch {
		k.SetZone(ctx, &zone)
	}

	return nil
}

func (k *Keeper) UpdateDelegationRecordForAddress(
	ctx sdk.Context,
	delegatorAddress,
	validatorAddress string,
	amount sdk.Coin,
	zone *types.Zone,
	absolute bool,
	isEpoch bool,
) error {
	delegation, found := k.GetDelegation(ctx, zone.ChainId, delegatorAddress, validatorAddress)

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
	k.SetDelegation(ctx, zone.ChainId, delegation)

	period := int64(k.GetParam(ctx, types.KeyValidatorSetInterval))
	query := stakingtypes.QueryValidatorsRequest{}
	err := k.EmitValSetQuery(ctx, zone.ConnectionId, zone.ChainId, query, math.NewInt(period))
	if err != nil {
		return err
	}

	if isEpoch {
		err = zone.DecrementWithdrawalWaitgroup(k.Logger(ctx), 1, "delegation_epoch success")
		if err != nil {
			k.Logger(ctx).Error(err.Error())
			// return nil here as to not fail the ack, but don't trigger RR multiple times.
			return nil
		}

		k.SetZone(ctx, zone)

		if zone.GetWithdrawalWaitgroup() == 0 {
			k.Logger(ctx).Info("Triggering redemption rate upgrade after delegation updates")
			err = k.TriggerRedemptionRate(ctx, zone)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (k *Keeper) HandleWithdrawRewards(ctx sdk.Context, msg sdk.Msg, connectionID string) error {
	withdrawalMsg, ok := msg.(*distrtypes.MsgWithdrawDelegatorReward)
	if !ok {
		k.Logger(ctx).Error("unable to cast source message to MsgWithdrawDelegatorReward")
		return errors.New("unable to cast source message to MsgWithdrawDelegatorReward")
	}

	zone, err := k.GetZoneFromConnectionID(ctx, connectionID)
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
		defer k.EventManagerKeeper.MarkCompleted(ctx, types.ModuleName, zone.ChainId, fmt.Sprintf("%s/%s", "withdraw_rewards_epoch", withdrawalMsg.ValidatorAddress))
	}
	return nil
}

func (k *Keeper) TriggerRedemptionRate(ctx sdk.Context, zone *types.Zone) error {
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

	k.EventManagerKeeper.AddEvent(
		ctx,
		types.ModuleName,
		zone.ChainId,
		"query_withdrawal_balance",
		"",
		emtypes.EventTypeICQAccountBalances,
		emtypes.EventStatusActive,
		nil,
		nil,
	)

	return nil
}

func DistributeRewardsFromWithdrawAccount(k *Keeper, ctx sdk.Context, args []byte, query querytypes.Query) error {
	defer k.EventManagerKeeper.MarkCompleted(ctx, types.ModuleName, query.ChainId, "query_withdrawal_balance")
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

	msgs := make([]sdk.Msg, 0)

	if rewards.Amount.IsPositive() {
		msgs = append(msgs, k.prepareRewardsDistributionMsgs(zone, rewards.Amount))
	}

	// multiDenomFee is the balance of withdrawal account minus the redelegated rewards.
	multiDenomFee := withdrawBalance.Balances.Sub(sdk.Coins{rewards}...)

	var remotePort string
	var remoteChannel string
	k.IBCKeeper.ChannelKeeper.IterateChannels(ctx, func(channel channeltypes.IdentifiedChannel) bool {
		if channel.ConnectionHops[0] == zone.ConnectionId && channel.PortId == types.TransferPort && channel.State == channeltypes.OPEN {
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
	k.UpdateRedemptionRate(ctx, &zone, rewards.Amount)

	// send tx
	return k.SubmitTx(ctx, msgs, zone.WithdrawalAddress, "", zone.MessagesPerTx)
}

func (*Keeper) prepareRewardsDistributionMsgs(zone types.Zone, rewards math.Int) sdk.Msg {
	return &banktypes.MsgSend{
		FromAddress: zone.WithdrawalAddress.GetAddress(),
		ToAddress:   zone.DelegationAddress.GetAddress(),
		Amount:      sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, rewards)),
	}
}

func isNumericString(in string) bool {
	// It is okay to use strconv.ParseInt to test if a value is numeric
	// because the total supply of QCK is:
	//      400_000_000 (400 million) qck aka 400_000_000_000_000 uqck
	// and to parse numeric values, say in the smallest unit of uqck
	//      MaxInt64: (1<<63)-1 = 9_223_372_036_854_775_807 uqck aka
	//                            9_223_372_036_854.775 (9.223 Trillion) qck
	// so the function is appropriate as its range won't be exceeded.
	_, err := strconv.ParseInt(in, 10, 64)
	return err == nil
}

func equalLsmCoin(valoper string, amount math.Int, lsmAmount sdk.Coin) bool {
	parts := strings.Split(lsmAmount.Denom, "/")
	if len(parts) == 2 && strings.HasPrefix(parts[0], valoper) && isNumericString(parts[1]) {
		return lsmAmount.Amount.Equal(amount)
	}
	return false
}

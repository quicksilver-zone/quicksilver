package keeper

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"cosmossdk.io/math"
	"github.com/golang/protobuf/proto" //nolint:staticcheck

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	icatypes "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v5/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v5/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v5/modules/core/04-channel/types"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	lsmstakingtypes "github.com/iqlusioninc/liquidity-staking-module/x/staking/types"

	"github.com/ingenuity-build/quicksilver/utils"
	queryTypes "github.com/ingenuity-build/quicksilver/x/interchainquery/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

const transferPort = "transfer"

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
	if reflect.DeepEqual(packetData, icatypes.InterchainAccountPacketData{}) {
		return errors.New("unable to unmarshal packet data; got empty JSON object")
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
			response := lsmstakingtypes.MsgRedeemTokensforSharesResponse{}
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
			response := lsmstakingtypes.MsgTokenizeSharesResponse{}
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
			k.Logger(ctx).Error("Redelegation initiated", "response", response)
			if err := k.HandleBeginRedelegate(ctx, src, response.CompletionTime, packetData.Memo); err != nil {
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

func (k *Keeper) HandleCompleteMultiSend(ctx sdk.Context, msg sdk.Msg, memo string) error {
	k.Logger(ctx).Info("Received MsgMultiSend acknowledgement")
	// first, type assertion. we should have banktypes.MsgMultiSend
	sMsg, ok := msg.(*banktypes.MsgMultiSend)
	if !ok {
		k.Logger(ctx).Error("unable to cast source message to MsgMultiSend")
		return errors.New("unable to cast source message to MsgMultiSend")
	}

	// check for sending of tokens from deposit -> delegate.
	zone, err := k.GetZoneFromContext(ctx)
	if err != nil {
		err = fmt.Errorf("1: %w", err)
		k.Logger(ctx).Error(err.Error())
		return err
	}

	for _, out := range sMsg.Outputs {
		// coerce banktype.Output to banktype.MsgSend
		// to use in handleSendToDelegate
		msg := banktypes.MsgSend{
			FromAddress: "",
			ToAddress:   out.Address,
			Amount:      out.Coins,
		}
		if err := k.handleSendToDelegate(ctx, zone, &msg, memo); err != nil {
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

	k.Logger(ctx).Error("messages to send", "messages", msgs)

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
		if err = k.BankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(withdrawalRecord.BurnAmount)); err != nil {
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
					if err = k.BankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(withdrawalRecord.BurnAmount)); err != nil {
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

// GetUnlockedTokensForZone will iterate over all delegation records for a zone, and then remove the
// locked tokens (those actively being redelegated), returning a slice of int64 staking tokens that
// are unlocked and free to redelegate or unbond.
func (k *Keeper) GetUnlockedTokensForZone(ctx sdk.Context, zone *types.Zone) map[string]int64 {
	availablePerValidator := map[string]int64{}
	for _, delegation := range k.GetAllDelegations(ctx, zone) {
		thisAvailable, found := availablePerValidator[delegation.ValidatorAddress]
		if !found {
			thisAvailable = 0
		}
		availablePerValidator[delegation.ValidatorAddress] = thisAvailable + delegation.Amount.Amount.Int64()
	}
	for _, redelegation := range k.ZoneRedelegationRecords(ctx, zone.ChainId) {
		thisAvailable, found := availablePerValidator[redelegation.Destination]
		if found {
			availablePerValidator[redelegation.Destination] = thisAvailable - redelegation.Amount
		}
	}
	return availablePerValidator
}

// handle queued unbondings is called once per epoch to aggregate all queued unbondings into
// a single unbond transaction per delegation.
func (k *Keeper) HandleQueuedUnbondings(ctx sdk.Context, zone *types.Zone, epoch int64) error {
	// out here will only ever be in native bond denom
	out := make(map[string]sdk.Coin, 0)
	txhashes := make(map[string][]string, 0)

	availablePerValidator := k.GetUnlockedTokensForZone(ctx, zone)

	var err error
	k.IterateZoneStatusWithdrawalRecords(ctx, zone.ChainId, WithdrawStatusQueued, func(idx int64, withdrawal types.WithdrawalRecord) bool {
		// copy this so we can rollback on fail
		thisAvail := availablePerValidator
		thisOut := make(map[string]sdk.Coin, 0)
		k.Logger(ctx).Info("unbonding funds", "from", withdrawal.Delegator, "to", withdrawal.Recipient, "amount", withdrawal.Amount)
		for _, dist := range withdrawal.Distribution {
			if thisAvail[dist.Valoper] < int64(dist.Amount) {
				// we cannot satisfy this unbond this epoch.
				k.Logger(ctx).Error("unable to satisfy unbonding for this epoch, due to locked tokens.", "txhash", withdrawal.Txhash, "user", withdrawal.Delegator, "chain", zone.ChainId)
				return false
			}
			thisOut[dist.Valoper] = sdk.NewCoin(zone.BaseDenom, math.NewIntFromUint64(dist.Amount))
			thisAvail[dist.Valoper] -= int64(dist.Amount)

			// if the validator has been historically slashed, and delegatorShares does not match tokens, then we end up with 'clipping'.
			// clipping is the truncation of the expected unbonding amount because of the need to have whole integer tokens.
			// the amount unbonded is emitted as an event, but not in the response, so we never _know_ this has happened.
			// as such, if we know the validator has hisotrical slashing, we remove 1 utoken from the distribution for this validator, with
			// the expectation that clipping will occur. We do not reduce the amount requested to unbond.
			val, found := zone.GetValidatorByValoper(dist.Valoper)
			if !found {
				// something kooky is going on...
				err = fmt.Errorf("unable to find a validator we expected to exist [%s]", dist.Valoper)
				return true
			}
			if val.DelegatorShares.Equal(sdk.NewDecFromInt(val.VotingPower)) {
				dist.Amount--
			}
		}

		// update record of available balances.
		availablePerValidator = thisAvail

		for valoper, amount := range thisOut {
			existing, found := out[valoper]
			if !found {
				out[valoper] = amount
				txhashes[valoper] = []string{withdrawal.Txhash}

			} else {
				out[valoper] = existing.Add(amount)
				txhashes[valoper] = append(txhashes[valoper], withdrawal.Txhash)

			}
		}

		k.UpdateWithdrawalRecordStatus(ctx, &withdrawal, WithdrawStatusUnbond)
		return false
	})
	if err != nil {
		return err
	}

	if len(txhashes) == 0 {
		// no records to handle.
		return nil
	}

	var msgs []sdk.Msg
	for _, valoper := range utils.Keys(out) {
		sort.Strings(txhashes[valoper])
		k.SetUnbondingRecord(ctx, types.UnbondingRecord{ChainId: zone.ChainId, EpochNumber: epoch, Validator: valoper, RelatedTxhash: txhashes[valoper]})
		msgs = append(msgs, &stakingtypes.MsgUndelegate{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: valoper, Amount: out[valoper]})
	}

	k.Logger(ctx).Error("unbonding messages", "msg", msgs)

	return k.SubmitTx(ctx, msgs, zone.DelegationAddress, fmt.Sprintf("withdrawal/%d", epoch))
}

func (k *Keeper) GCCompletedUnbondings(ctx sdk.Context, zone *types.Zone) error {
	var err error

	k.IterateZoneStatusWithdrawalRecords(ctx, zone.ChainId, WithdrawStatusCompleted, func(idx int64, withdrawal types.WithdrawalRecord) bool {
		if ctx.BlockTime().After(withdrawal.CompletionTime.Add(24 * time.Hour)) {
			k.Logger(ctx).Info("garbage collecting completed unbondings")
			k.DeleteWithdrawalRecord(ctx, zone.ChainId, withdrawal.Txhash, WithdrawStatusCompleted)
		}
		return false
	})

	return err
}

func (k *Keeper) GCCompletedRedelegations(ctx sdk.Context) error {
	var err error

	k.IterateRedelegationRecords(ctx, func(idx int64, key []byte, redelegation types.RedelegationRecord) bool {
		if ctx.BlockTime().After(redelegation.CompletionTime) {
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
		if ctx.BlockTime().After(withdrawal.CompletionTime) && !withdrawal.CompletionTime.IsZero() { // completion date has passed.
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
	if len(parts) != 2 {
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
	record, found := k.GetRedelegationRecord(ctx, zone.ChainId, redelegateMsg.ValidatorSrcAddress, redelegateMsg.DelegatorAddress, epochNumber)
	if !found {
		k.Logger(ctx).Error("unable to find redelegation record")
		return errors.New("unable to find redelegation record")
	}
	k.Logger(ctx).Error("updating redelegation record with completion time")
	record.CompletionTime = completion
	k.SetRedelegationRecord(ctx, record)
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
	if len(memoParts) != 2 {
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
			return errors.New("unable to lookup withdrawal record")
		}
		if completion.After(record.CompletionTime) {
			record.CompletionTime = completion
		}
		k.Logger(ctx).Error("withdrawal record to save", "rcd", record)
		k.UpdateWithdrawalRecordStatus(ctx, &record, WithdrawStatusUnbond)
	}
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

func (k *Keeper) HandleDelegate(ctx sdk.Context, msg sdk.Msg) error {
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
			return errors.New("unable to find zone")
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
	k.Logger(ctx).Error("ERROR 1", "response", response)
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

		if err := k.RemoveDelegation(ctx, &zone, existingDelegation); err != nil {
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

	// k.SetZone(ctx, &zone)

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
		k.Logger(ctx).Error("Updating delegation tuple amount", "delegator", delegatorAddress, "validator", validatorAddress, "old_amount", oldAmount, "inbound_amount", amount.Amount, "new_amount", delegation.Amount, "abs", absolute)
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
		k.Logger(ctx).Error("WAITGROUP DECREMENTED", "wg", zone.WithdrawalWaitgroup)
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
		k.Logger(ctx).Error("TRIGGER DISTRIBUTE REWARDS")
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
				TimeoutTimestamp: uint64(ctx.BlockTime().UnixNano() + 5*time.Minute.Nanoseconds()),
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

package keeper

import (
	"encoding/json"
	"fmt"
	"strings"

	//lint:ignore SA1019 ignore this!
	"github.com/golang/protobuf/proto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
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

	packetData := icatypes.InterchainAccountPacketData{}
	if err := json.Unmarshal(packet.Data, &packetData); err != nil {
		return err
	}
	msgs, err := icatypes.DeserializeCosmosTx(k.cdc, packetData.Data)
	if err != nil {
		k.Logger(ctx).Info("Error decoding messages", "err", err)
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
			k.Logger(ctx).Info("Tokens redeemed for shares", "response", response)
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
			k.Logger(ctx).Info("Redelegation initiated", "response", response)
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
			k.Logger(ctx).Info("Funds Transferred", "response", response)
			// check tokenTransfers - if end user unescrow and burn txs
			if err := k.HandleCompleteSend(ctx, src); err != nil {
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
			k.Logger(ctx).Info("Funds Transferred (Multi)", "response", response)
			if err := k.HandleCompleteMultiSend(ctx, src); err != nil {
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

func (k *Keeper) HandleCompleteMultiSend(ctx sdk.Context, msg sdk.Msg) error {
	k.Logger(ctx).Info("Received MsgMultiSend acknowledgement")
	// first, type assertion. we should have banktypes.MsgMultiSend
	sMsg, ok := msg.(*banktypes.MsgMultiSend)
	if !ok {
		k.Logger(ctx).Error("unable to cast source message to MsgMultiSend")
		return fmt.Errorf("unable to cast source message to MsgMultiSend")
	}

	// check for sending of tokens from deposit -> delegate.
	zone := k.GetZoneForDelegateAccount(ctx, sMsg.Outputs[0].Address) // do this once, save multiple lookups.
	if zone != nil {
		for _, out := range sMsg.Outputs {
			da, err := zone.GetDelegationAccountByAddress(out.Address)
			if err != nil {
				return err
			}
			da.Balance = da.Balance.Add(out.Coins...)
			k.Delegate(ctx, *zone, da)
		}
	}

	return nil
}

func (k *Keeper) HandleCompleteSend(ctx sdk.Context, msg sdk.Msg) error {
	k.Logger(ctx).Info("Received MsgSend acknowledgement")
	// first, type assertion. we should have banktypes.MsgSend
	var err error = nil
	var done bool = false
	sMsg, ok := msg.(*banktypes.MsgSend)
	if !ok {
		k.Logger(ctx).Error("unable to cast source message to MsgSend")
		return fmt.Errorf("unable to cast source message to MsgSend")
	}

	// first check for withdrawals (if FromAddress is a DelegateAccount)
	k.IterateWithdrawalRecords(ctx, sMsg.FromAddress, func(idx int64, withdrawal types.WithdrawalRecord) bool {
		k.Logger(ctx).Debug("iterating withdraw record", "idx", idx, "record", withdrawal)
		if withdrawal.Recipient == sMsg.ToAddress {
			k.Logger(ctx).Debug("matched the recipient", "val", withdrawal.Delegator, "recipient", withdrawal.Recipient)
			z := k.GetZoneForDelegateAccount(ctx, withdrawal.Delegator)
			if sMsg.Amount.AmountOf(z.BaseDenom).Equal(withdrawal.Amount.Amount) {
				k.Logger(ctx).Debug("matched the amount", "amount", sMsg.Amount, "record.amount", withdrawal.Amount.Amount)
				if withdrawal.Status == WITHDRAW_STATUS_SEND {
					k.Logger(ctx).Info("Found matching withdrawal; withdrawal marked as completed")
					k.DeleteWithdrawalRecord(ctx, withdrawal.Delegator, withdrawal.Validator, withdrawal.Recipient)
					da, err := z.GetDelegationAccountByAddress(withdrawal.Delegator)
					if err != nil {
						return true
					}
					da.Balance = da.Balance.Sub(sMsg.Amount)

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
	zone := k.GetZoneForDelegateAccount(ctx, sMsg.ToAddress)
	if zone != nil { // this _is_ a delegate account
		da, err := zone.GetDelegationAccountByAddress(sMsg.ToAddress)
		if err != nil {
			return err
		}
		da.Balance = da.Balance.Add(sMsg.Amount...)
		k.Delegate(ctx, *zone, da)
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
					delegatorIca := k.GetICAForDelegateAccount(ctx, withdrawal.Delegator)
					if delegatorIca == nil {
						k.Logger(ctx).Error("unable to find delegator account for withdrawal; this shouldn't happen", err)
						return true
					}
					sendMsg := &banktypes.MsgSend{FromAddress: withdrawal.Delegator, ToAddress: withdrawal.Recipient, Amount: sdk.Coins{amount}}

					err = k.SubmitTx(ctx, []sdk.Msg{sendMsg}, delegatorIca)
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
	return k.UpdateDelegationRecordForAddress(ctx, redeemMsg.DelegatorAddress, validatorAddress, amount)
}

func (k *Keeper) HandleDelegate(ctx sdk.Context, msg sdk.Msg) error {
	k.Logger(ctx).Info("Received MsgDelegate acknowledgement")
	// first, type assertion. we should have stakingtypes.MsgDelegate
	delegateMsg, ok := msg.(*stakingtypes.MsgDelegate)
	if !ok {
		k.Logger(ctx).Error("unable to cast source message to MsgDelegate")
		return fmt.Errorf("unable to cast source message to MsgDelegate")
	}
	return k.UpdateDelegationRecordForAddress(ctx, delegateMsg.DelegatorAddress, delegateMsg.ValidatorAddress, delegateMsg.Amount)
}

func (k *Keeper) GetValidatorForToken(ctx sdk.Context, delegatorAddress string, amount sdk.Coin) (string, error) {
	zone := k.GetZoneForDelegateAccount(ctx, delegatorAddress)
	if zone == nil {
		return "", fmt.Errorf("unable to fetch zone for delegate address %s", delegatorAddress)
	}

	for _, val := range zone.Validators {
		if strings.HasPrefix(amount.Denom, val.ValoperAddress) {
			// match!
			return val.ValoperAddress, nil
		}
	}

	return "", fmt.Errorf("unable to find validator for token %s", amount.Denom)

}

func (k *Keeper) UpdateDelegationRecordForAddress(ctx sdk.Context, delegatorAddress string, validatorAddress string, amount sdk.Coin) error {

	var validator *types.Validator
	var err error

	zone := k.GetZoneForDelegateAccount(ctx, delegatorAddress)
	if zone == nil {
		return fmt.Errorf("unable to fetch zone for delegate address %s", delegatorAddress)
	}

	validator, err = zone.GetValidatorByValoper(validatorAddress)
	if err != nil {
		return err
	}

	delegation, err := validator.GetDelegationForDelegator(delegatorAddress)
	if err != nil {
		if validator.Delegations == nil {
			validator.Delegations = []*types.Delegation{}
		}
		k.Logger(ctx).Info("Adding delegation tuple", "delegator", delegatorAddress, "validator", validator.ValoperAddress, "amount", amount.Amount)
		delegation = &types.Delegation{
			DelegationAddress: delegatorAddress,
			ValidatorAddress:  validator.ValoperAddress,
			Amount:            amount.Amount.ToDec(),
			Rewards:           sdk.Coins{},
			RedelegationEnd:   0,
		}
		validator.Delegations = append(validator.Delegations, delegation)
	} else {
		k.Logger(ctx).Info("Updating delegation tuple amount", "delegator", delegatorAddress, "validator", validator.ValoperAddress, "old_amount", delegation.Amount, "inbound_amount", amount.Amount)
		delegation.Amount = delegation.Amount.Add(amount.Amount.ToDec())
	}

	da, err := zone.GetDelegationAccountByAddress(delegation.DelegationAddress)

	if err != nil {
		k.Logger(ctx).Error("Unable to retrieve delegation account", "delegator", delegatorAddress)
		return err
	}

	if da.DelegatedBalance.IsNil() || da.DelegatedBalance.IsZero() {
		da.DelegatedBalance = amount
	} else {
		da.DelegatedBalance = da.DelegatedBalance.Add(amount)
	}

	zone.UpdateDelegatedAmount()
	k.SetRegisteredZone(ctx, *zone)
	return nil
}

func (k *Keeper) HandleWithdrawRewards(ctx sdk.Context, msg sdk.Msg, amount sdk.Coins) error {
	k.Logger(ctx).Info("Received MsgWithdrawDelegatorReward acknowledgement")
	// first, type assertion. we should have distrtypes.MsgWithdrawDelegatorReward
	withdrawMsg, ok := msg.(*distrtypes.MsgWithdrawDelegatorReward)
	if !ok {
		k.Logger(ctx).Error("unable to cast source message to MsgWithdrawDelegatorReward")
		return fmt.Errorf("unable to cast source message to MsgWithdrawDelegatorReward")
	}
	zone := k.GetZoneForDelegateAccount(ctx, withdrawMsg.DelegatorAddress)
	if zone == nil {
		return fmt.Errorf("unable to find zone for delegator account %s", withdrawMsg.DelegatorAddress)
	}
	da, err := zone.GetDelegationAccountByAddress(withdrawMsg.DelegatorAddress)
	if err != nil {
		return err
	}
	return k.Delegate(ctx, *zone, da)
}

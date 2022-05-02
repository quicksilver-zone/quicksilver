package keeper

import (
	"encoding/json"
	"fmt"
	"strings"

	//lint:ignore SA1019 ignore this!
	"github.com/golang/protobuf/proto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
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
		k.Logger(ctx).Error("Unable to unmarshal acknowledgement packet data", "error", err, "data", packetData)
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
			k.Logger(ctx).Debug("Shares tokenized", "response", response)
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
			k.Logger(ctx).Debug("Delegated", "response", response)
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
			k.Logger(ctx).Debug("Funds Transferred (Multi)", "response", response)
			if err := k.HandleCompleteMultiSend(ctx, src); err != nil {
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

// TODO: rework to reflect changes to HandleWithdrawRewards:
//   1. handle MsgSend from WithdrawalAccount to FeeAccount;
//   2. handle MsgSend from WithdrawalAccount to Delegation Accounts;
func (k *Keeper) HandleCompleteSend(ctx sdk.Context, msg sdk.Msg) error {
	k.Logger(ctx).Info("Received MsgSend acknowledgement")
	// first, type assertion. we should have banktypes.MsgSend
	sMsg, ok := msg.(*banktypes.MsgSend)
	if !ok {
		err := fmt.Errorf("unable to cast source message to MsgSend")
		k.Logger(ctx).Error(err.Error())
		return err
	}

	// get zone
	var zone *types.RegisteredZone
	zone = k.GetZoneForDelegateAccount(ctx, sMsg.ToAddress)
	if zone == nil {
		zone = k.GetZoneForDelegateAccount(ctx, sMsg.FromAddress)
		if zone == nil {
			err := fmt.Errorf("unable to find delegate account for %s or %s", sMsg.ToAddress, sMsg.FromAddress)
			k.Logger(ctx).Error(err.Error())
			return err
		}
	}

	// checks here are specific to ensure future extensibility;
	switch {
	case sMsg.FromAddress == zone.WithdrawalAddress.GetAddress() && sMsg.ToAddress == zone.FeeAddress.GetAddress():
		// WithdrawalAddress (for rewards) only send to FeeAddress or DelegationAddresses.
		// Target here is FeeAddress.
		// shouldn't be called.
		panic("unexpected")
	case sMsg.FromAddress == zone.WithdrawalAddress.GetAddress() && sMsg.ToAddress != zone.FeeAddress.GetAddress():
		// WithdrawalAddress (for rewards) only send to FeeAddress or DelegationAddresses.
		// Target here is one of the DelegationAddresses.
		if err := k.handleRewardsDelegation(ctx, *zone, sMsg); err != nil {
			return err
		}
	default:
		if err := k.handleWithdrawForUser(ctx, sMsg); err != nil {
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
	return k.Delegate(ctx, zone, da)
}

func (k *Keeper) handleWithdrawForUser(ctx sdk.Context, msg *banktypes.MsgSend) error {
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
		da, err := zone.GetDelegationAccountByAddress(msg.ToAddress)
		if err != nil {
			return err
		}
		da.Balance = da.Balance.Add(msg.Amount...)
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
					_, delegatorIca := k.GetICAForDelegateAccount(ctx, withdrawal.Delegator)
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
	err = k.UpdateDelegationRecordForAddress(ctx, redeemMsg.DelegatorAddress, validatorAddress, amount)
	if err != nil {
		return err
	}
	zone, da := k.GetICAForDelegateAccount(ctx, redeemMsg.DelegatorAddress)
	da.DelegatedBalance = da.DelegatedBalance.Add(amount)
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
	err := k.UpdateDelegationRecordForAddress(ctx, delegateMsg.DelegatorAddress, delegateMsg.ValidatorAddress, delegateMsg.Amount)
	if err != nil {
		return err
	}
	zone, da := k.GetICAForDelegateAccount(ctx, delegateMsg.DelegatorAddress)
	da.DelegatedBalance = da.DelegatedBalance.Add(delegateMsg.Amount)
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

func (k *Keeper) GetValidatorForToken(ctx sdk.Context, delegatorAddress string, amount sdk.Coin) (string, error) {
	zone := k.GetZoneForDelegateAccount(ctx, delegatorAddress)
	if zone == nil {
		return "", fmt.Errorf("unable to fetch zone for delegate address %s", delegatorAddress)
	}

	for _, val := range zone.GetValidatorsAddressesAsSlice() {
		if strings.HasPrefix(amount.Denom, val) {
			// match!
			return val, nil
		}
	}

	return "", fmt.Errorf("unable to find validator for token %s", amount.Denom)

}

func (k *Keeper) UpdateDelegationRecordsForAddress(ctx sdk.Context, zone types.RegisteredZone, delegatorAddress string, args []byte) error {
	var response stakingtypes.QueryDelegatorDelegationsResponse
	err := k.cdc.UnmarshalJSON(args, &response)
	if err != nil {
		return err
	}

	delegatorSum := sdk.NewCoin(zone.BaseDenom, sdk.ZeroInt())
	for _, delegation := range response.DelegationResponses {
		err = k.UpdateDelegationRecordForAddress(ctx, delegatorAddress, delegation.Delegation.ValidatorAddress, delegation.Balance)
		delegatorSum = delegatorSum.Add(delegation.Balance)
		if err != nil {
			return err
		}
	}

	zone.WithdrawalWaitgroup--
	k.Logger(ctx).Info("Decrementing waitgroup", "value", zone.WithdrawalWaitgroup)
	da, err := zone.GetDelegationAccountByAddress(delegatorAddress)
	if err != nil {
		return err
	}
	da.DelegatedBalance = delegatorSum

	k.SetRegisteredZone(ctx, zone)

	return nil
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
		if !delegation.Amount.Equal(amount.Amount.ToDec()) {
			k.Logger(ctx).Info("Updating delegation tuple amount", "delegator", delegatorAddress, "validator", validator.ValoperAddress, "old_amount", delegation.Amount, "inbound_amount", amount.Amount)
			delegation.Amount = delegation.Amount.Add(amount.Amount.ToDec())
		}

	}

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
	// decrement withdrawal waitgroup
	zone.WithdrawalWaitgroup--
	k.Logger(ctx).Error("Withdrawal waitgroup DOWN", "value", zone.WithdrawalWaitgroup)

	k.SetRegisteredZone(ctx, *zone)

	switch zone.WithdrawalWaitgroup {
	case 0:
		// interface assertion
		var cb Callback = DistributeRewardsFromWithdrawAccount

		// total rewards balance withdrawn
		k.ICQKeeper.MakeRequest(
			ctx,
			zone.ConnectionId,
			zone.ChainId, "cosmos.bank.v1beta1.Query/AllBalances",
			map[string]string{"address": zone.WithdrawalAddress.Address},
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

	err := k.cdc.UnmarshalJSON(args, &withdrawBalance)
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
	return k.SubmitTx(ctx, msgs, zone.WithdrawalAddress)

}

func (k *Keeper) updateRedemptionRate(ctx sdk.Context, zone types.RegisteredZone, epochRewards sdk.Coin) {
	ratio := zone.GetDelegatedAmount().Add(epochRewards).Amount.ToDec().Quo(k.BankKeeper.GetSupply(ctx, zone.LocalDenom).Amount.ToDec())
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

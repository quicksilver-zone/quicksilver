package keeper

import (
	"encoding/json"

	"github.com/golang/protobuf/proto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

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
	json.Unmarshal(packet.Data, &packetData)
	msgs, err := icatypes.DeserializeCosmosTx(k.cdc, packetData.Data)
	if err != nil {
		k.Logger(ctx).Info("Error decoding messages", "err", err)
	}

	for msgIndex, msgData := range txMsgData.Data {
		src := msgs[msgIndex]
		switch msgData.MsgType {
		case "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward":
			// response := distrtypes.MsgWithdrawDelegatorRewardResponse{}
			// err := proto.Unmarshal(msgData.Data, &response)
			// if err != nil {
			// 	k.Logger(ctx).Error("Unable to unmarshal MsgWithdrawDelegatorReward response", "error", err)
			// 	return err
			// }
			// k.Logger(ctx).Info("Rewards withdrawn", "response", response)
			// noop here - we can plausibl
			continue
		case "/cosmos.staking.v1beta1.MsgRedeemTokensforShares":
			// response := stakingtypes.MsgRedeemTokensforSharesResponse{}
			// err := proto.Unmarshal(msgData.Data, &response)
			// if err != nil {
			// 	k.Logger(ctx).Error("Unable to unmarshal MsgRedeemTokensforShares response", "error", err)
			// 	return err
			// }
			// k.Logger(ctx).Info("Tokens redeemed for shares", "response", response)
			// noop
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
			k.HandleTokenizedShares(ctx, src, response.Amount)
			continue
		case "/cosmos.staking.v1beta1.MsgDelegate":
			// response := stakingtypes.MsgDelegateResponse{}
			// err := proto.Unmarshal(msgData.Data, &response)
			// if err != nil {
			// 	k.Logger(ctx).Error("Unable to unmarshal MsgDelegate response", "error", err)
			// 	return err
			// }
			// k.Logger(ctx).Info("Delegated", "response", response)
			// no action
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

	return nil
}

func (k *Keeper) HandleCompleteSend(ctx sdk.Context, msg sdk.Msg) error {

	return nil
}

func (k *Keeper) HandleTokenizedShares(ctx sdk.Context, msg sdk.Msg, amount sdk.Coin) error {

	return nil
}

func (k *Keeper) HandleBeginRedelegate(ctx sdk.Context, msg sdk.Msg, completion time.Time) error {

	return nil
}

package keeper

import (
	"errors"
	"time"

	"github.com/cosmos/gogoproto/proto"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	icacontrollerkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"

	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

const (
	Unset           = "unset"
	ICAMsgChunkSize = 5
	ICATimeout      = time.Hour * 6
	// AuthzAutoClaimAddress = "quick1psevptdp90jad76zt9y9x2nga686hutgmasmwd"
)

// SendTokenIBC is a helper function that finds the zone channel and performs an ibc transfer from senderAccAddress
// to receiver.
func (k *Keeper) SendTokenIBC(ctx sdk.Context, senderAccAddress sdk.AccAddress, receiver string, zone *types.Zone, coin sdk.Coin) error {
	var srcPort string
	var srcChannel string

	k.IBCKeeper.ChannelKeeper.IterateChannels(ctx, func(channel channeltypes.IdentifiedChannel) bool {
		if channel.ConnectionHops[0] == zone.ConnectionId && channel.PortId == types.TransferPort && channel.State == channeltypes.OPEN {
			srcChannel = channel.ChannelId
			srcPort = channel.PortId
			return true
		}
		return false
	})
	if srcPort == "" {
		return errors.New("unable to find remote transfer connection")
	}

	_, err := k.TransferKeeper.Transfer(ctx, &transfertypes.MsgTransfer{
		SourcePort:    srcPort,
		SourceChannel: srcChannel,
		Token:         coin,
		Sender:        senderAccAddress.String(),
		Receiver:      receiver,
		TimeoutHeight: clienttypes.Height{
			RevisionNumber: 0,
			RevisionHeight: 0,
		},
		TimeoutTimestamp: uint64(ctx.BlockTime().UnixNano() + 5*time.Minute.Nanoseconds()), //nolint:gosec
		Memo:             "",
	})
	return err
}

func (k *Keeper) SubmitTx(ctx sdk.Context, msgs []sdk.Msg, account *types.ICAAccount, memo string, messagesPerTx int64) error {
	return k.txSubmit(ctx, k, msgs, account, memo, messagesPerTx)
}

// SubmitTx submits a Tx on behalf of an ICAAccount to a remote chain.
func ProdSubmitTx(ctx sdk.Context, k *Keeper, msgs []sdk.Msg, account *types.ICAAccount, memo string, messagesPerTx int64) error {
	// if no messages, do nothing
	if len(msgs) == 0 {
		return nil
	}

	portID := account.GetPortName()
	connectionID, err := k.GetConnectionForPort(ctx, portID)
	if err != nil {
		return err
	}

	chunkSize := int(messagesPerTx)
	if chunkSize < 1 {
		chunkSize = ICAMsgChunkSize
	}

	timeoutTimestamp := uint64(ctx.BlockTime().Add(ICATimeout).UnixNano()) //nolint:gosec

	for len(msgs) > 0 {
		// if the last chunk, make chunksize the number of messages
		if len(msgs) < chunkSize {
			chunkSize = len(msgs)
		}

		// remove chunk from original msg slice
		msgsChunk := msgs[0:chunkSize]
		msgs = msgs[chunkSize:]

		protoMsgs := make([]proto.Message, len(msgsChunk))
		for i, msg := range msgsChunk {
			protoMsgs[i] = msg.(proto.Message)
		}
		// build and submit message for this chunk
		data, err := icatypes.SerializeCosmosTx(k.cdc, protoMsgs)
		if err != nil {
			return err
		}

		// validate memo < 256 bytes
		packetData := icatypes.InterchainAccountPacketData{
			Type: icatypes.EXECUTE_TX,
			Data: data,
			Memo: memo,
		}

		ckMsgServer := icacontrollerkeeper.NewMsgServerImpl(&k.ICAControllerKeeper)
		portOwner := portID[len(icatypes.ControllerPortPrefix):] // TODO: this is a hack to get the port owner; change PortName() to PortOwner() sans prefix.
		msgSendTx := icacontrollertypes.NewMsgSendTx(portOwner, connectionID, timeoutTimestamp, packetData)
		_, err = ckMsgServer.SendTx(ctx, msgSendTx)
		if err != nil {
			return err
		}
	}

	return nil
}

// ---------------------------------------------------------------

func (k Keeper) NilReceipt(ctx sdk.Context, zone *types.Zone, txhash string) {
	t := ctx.BlockTime()
	r := types.Receipt{ChainId: zone.ChainId, Sender: "", Txhash: txhash, Amount: sdk.Coins{}, FirstSeen: &t, Completed: &t}
	k.SetReceipt(ctx, r)
}

func (Keeper) NewCompletedReceipt(ctx sdk.Context, zone *types.Zone, sender, txhash string, amount sdk.Coins) *types.Receipt {
	t := ctx.BlockTime()
	return &types.Receipt{ChainId: zone.ChainId, Sender: sender, Txhash: txhash, Amount: amount, FirstSeen: &t, Completed: &t}
}

func (Keeper) NewReceipt(ctx sdk.Context, zone *types.Zone, sender, txhash string, amount sdk.Coins) *types.Receipt {
	t := ctx.BlockTime()
	return &types.Receipt{ChainId: zone.ChainId, Sender: sender, Txhash: txhash, Amount: amount, FirstSeen: &t}
}

// GetReceipt returns receipt for the given key.
func (k *Keeper) GetReceipt(ctx sdk.Context, chainID, txHash string) (types.Receipt, bool) {
	key := types.GetReceiptKey(chainID, txHash)
	receipt := types.Receipt{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixReceipt)
	bz := store.Get([]byte(key))
	if len(bz) == 0 {
		return receipt, false
	}

	k.cdc.MustUnmarshal(bz, &receipt)
	return receipt, true
}

// SetReceipt sets receipt info.
func (k *Keeper) SetReceipt(ctx sdk.Context, receipt types.Receipt) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixReceipt)
	bz := k.cdc.MustMarshal(&receipt)
	store.Set([]byte(types.GetReceiptKey(receipt.ChainId, receipt.Txhash)), bz)
}

// DeleteReceipt delete receipt info.
func (k *Keeper) DeleteReceipt(ctx sdk.Context, chainID, txHash string) {
	key := types.GetReceiptKey(chainID, txHash)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixReceipt)
	store.Delete([]byte(key))
}

// IterateReceipts iterate through receipts.
func (k *Keeper) IterateReceipts(ctx sdk.Context, fn func(index int64, receiptInfo types.Receipt) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixReceipt)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		receipt := types.Receipt{}
		k.cdc.MustUnmarshal(iterator.Value(), &receipt)
		stop := fn(i, receipt)
		if stop {
			break
		}
		i++
	}
}

func (k *Keeper) AllReceipts(ctx sdk.Context) []types.Receipt {
	receipts := make([]types.Receipt, 0)
	k.IterateReceipts(ctx, func(_ int64, receiptInfo types.Receipt) (stop bool) {
		receipts = append(receipts, receiptInfo)
		return false
	})
	return receipts
}

// IterateZoneReceipts iterates through receipts of the given zone.
func (k *Keeper) IterateZoneReceipts(ctx sdk.Context, chainID string, fn func(index int64, receiptInfo types.Receipt) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixReceipt)
	iterator := sdk.KVStorePrefixIterator(store, []byte(chainID))
	defer iterator.Close()

	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		receipt := types.Receipt{}
		k.cdc.MustUnmarshal(iterator.Value(), &receipt)
		stop := fn(i, receipt)
		if stop {
			break
		}
		i++
	}
}

// UserZoneReceipts returns all receipts of the given user for the given zone.
func (k *Keeper) UserZoneReceipts(ctx sdk.Context, zone *types.Zone, addr sdk.AccAddress) ([]types.Receipt, error) {
	receipts := make([]types.Receipt, 0)

	bech32Address, err := addressutils.EncodeAddressToBech32(zone.AccountPrefix, addr)
	if err != nil {
		return receipts, err
	}

	k.IterateZoneReceipts(ctx, zone.ChainId, func(_ int64, receipt types.Receipt) (stop bool) {
		if receipt.Sender == bech32Address {
			receipts = append(receipts, receipt)
		}
		return false
	})

	return receipts, nil
}

func (k *Keeper) SetReceiptsCompleted(ctx sdk.Context, chainID string, qualifyingTime, completionTime time.Time, denom string) {
	k.IterateZoneReceipts(ctx, chainID, func(_ int64, receiptInfo types.Receipt) (stop bool) {
		if receiptInfo.FirstSeen.Before(qualifyingTime) && receiptInfo.Completed == nil && denom == receiptInfo.Amount[0].Denom {
			receiptInfo.Completed = &completionTime
			k.SetReceipt(ctx, receiptInfo)

		}
		return false
	})
}

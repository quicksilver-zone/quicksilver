package keeper

import (
	"errors"
	"fmt"
	"time"

	sdkioerrors "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/types/tx"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	"github.com/gogo/protobuf/proto"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

const (
	Unset           = "unset"
	ICAMsgChunkSize = 5
)

func (k *Keeper) HandleReceiptForTransaction(ctx sdk.Context, txr *sdk.TxResponse, txn *tx.Tx, zone *types.Zone) error {
	k.Logger(ctx).Info("deposit receipt.", "ischeck", ctx.IsCheckTx(), "isrecheck", ctx.IsReCheckTx())
	hash := txr.TxHash
	memo := txn.Body.Memo

	senderAddress := Unset
	assets := sdk.Coins{}

	for _, event := range txr.Events {
		if event.Type == types.TransferPort {
			attrs := types.AttributesToMap(event.Attributes)
			sender := attrs["sender"]
			amount := attrs["amount"]
			if attrs["recipient"] == zone.DepositAddress.GetAddress() { // negate case where sender sends to multiple addresses in one tx
				if senderAddress == Unset {
					senderAddress = sender
				}

				if sender != senderAddress {
					k.Logger(ctx).Error("sender mismatch", "expected", senderAddress, "received", sender)
					return fmt.Errorf("sender mismatch: expected %q, got %q", senderAddress, sender)
				}

				k.Logger(ctx).Info("deposit receipt", "deposit_address", zone.DepositAddress.GetAddress(), "sender", sender, "amount", amount)
				thisCoins, err := sdk.ParseCoinsNormalized(amount)
				if err != nil {
					k.Logger(ctx).Error("unable to parse coin", "string", amount)
				}
				assets = assets.Add(thisCoins...)
			}
		}
	}

	if senderAddress == Unset {
		k.Logger(ctx).Error("no sender found. Ignoring.")
		return fmt.Errorf("no sender found. Ignoring")
	}

	// sdk.AccAddressFromBech32 doesn't work here as it expects the local HRP
	_, addressBytes, err := bech32.DecodeAndConvert(senderAddress)
	if err != nil {
		k.Logger(ctx).Error("unable to decode sender address. Ignoring.", "senderAddress", senderAddress)
		return fmt.Errorf("unable to decode sender address. Ignoring. senderAddress=%q", senderAddress)
	}
	var senderAccAddress sdk.AccAddress = addressBytes

	if err := zone.ValidateCoinsForZone(ctx, assets); err != nil {
		// we expect this to trigger if the valset has changed recently (i.e. we haven't seen the validator before.
		// That is okay, we'll catch it next round!)
		k.Logger(ctx).Error("unable to validate coins. Ignoring.", "senderAddress", senderAddress)
		return fmt.Errorf("unable to validate coins. Ignoring. senderAddress=%q", senderAddress)
	}

	k.Logger(ctx).Info("found new deposit tx", "deposit_address", zone.DepositAddress.GetAddress(), "senderAddress", senderAddress, "local", senderAccAddress.String(), "chain id", zone.ChainId, "assets", assets, "hash", hash)

	// update state
	if err := k.UpdateDelegatorIntent(ctx, senderAccAddress, zone, assets, memo); err != nil {
		k.Logger(ctx).Error("unable to update intent. Ignoring.", "senderAddress", senderAddress, "zone", zone.ChainId, "err", err.Error())
		return fmt.Errorf("unable to update intent. Ignoring. senderAddress=%q zone=%q err: %w", senderAddress, zone.ChainId, err)
	}

	if err := k.MintQAsset(ctx, senderAccAddress, senderAddress, zone, assets); err != nil {
		k.Logger(ctx).Error("unable to mint QAsset. Ignoring.", "senderAddress", senderAddress, "zone", zone.ChainId, "err", err.Error())
		return fmt.Errorf("unable to mint QAsset. Ignoring. senderAddress=%q zone=%q err: %w", senderAddress, zone.ChainId, err)
	}

	if err := k.TransferToDelegate(ctx, zone, assets, hash); err != nil {
		k.Logger(ctx).Error("unable to transfer to delegate. Ignoring.", "senderAddress", senderAddress, "zone", zone.ChainId, "err", err)
		return fmt.Errorf("unable to transfer to delegate. Ignoring. senderAddress=%q zone=%q err: %w", senderAddress, zone.ChainId, err)
	}

	// create receipt
	receipt := k.NewReceipt(ctx, zone, senderAddress, hash, assets)
	k.SetReceipt(ctx, *receipt)

	return nil
}

// MintQAsset mints qAssets based on the native asset redemption rate.  Tokens are then transferred to the given user.
func (k *Keeper) MintQAsset(ctx sdk.Context, sender sdk.AccAddress, senderAddress string, zone *types.Zone, assets sdk.Coins) error {
	if zone.RedemptionRate.IsZero() {
		return errors.New("zero redemption rate")
	}

	qAssets := sdk.Coins{}
	for _, asset := range assets.Sort() {
		amount := sdk.NewDecFromInt(asset.Amount).Quo(zone.RedemptionRate).TruncateInt()
		qAssets = qAssets.Add(sdk.NewCoin(zone.LocalDenom, amount))
	}

	k.Logger(ctx).Info("Minting qAssets for receipt", "assets", qAssets)
	err := k.BankKeeper.MintCoins(ctx, types.ModuleName, qAssets)
	if err != nil {
		return err
	}

	if zone.ReturnToSender {
		var srcPort string
		var srcChannel string

		k.IBCKeeper.ChannelKeeper.IterateChannels(ctx, func(channel channeltypes.IdentifiedChannel) bool {
			if channel.ConnectionHops[0] == zone.ConnectionId && channel.PortId == types.TransferPort && channel.State == channeltypes.OPEN {
				srcChannel = channel.Counterparty.ChannelId
				srcPort = channel.Counterparty.PortId
				return true
			}
			return false
		})
		if srcPort == "" {
			return errors.New("unable to find remote transfer connection")
		}

		err = k.TransferKeeper.SendTransfer(
			ctx,
			srcPort,
			srcChannel,
			qAssets[0],
			k.AccountKeeper.GetModuleAddress(types.ModuleName),
			senderAddress,
			clienttypes.Height{
				RevisionNumber: 0,
				RevisionHeight: 0,
			},
			uint64(ctx.BlockTime().UnixNano()+5*time.Minute.Nanoseconds()),
		)
	} else {
		err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, qAssets)
	}

	if err != nil {
		return fmt.Errorf("unable to transfer coins: %w", err)
	}

	k.Logger(ctx).Info("Transferred qAssets to sender", "assets", qAssets, "sender", sender)
	return nil
}

// TransferToDelegate transfers tokens from the zone deposit account address to the zone delegate account address.
func (k *Keeper) TransferToDelegate(ctx sdk.Context, zone *types.Zone, coins sdk.Coins, memo string) error {
	msg := &bankTypes.MsgSend{FromAddress: zone.DepositAddress.GetAddress(), ToAddress: zone.DelegationAddress.GetAddress(), Amount: coins}
	return k.SubmitTx(ctx, []proto.Message{msg}, zone.DepositAddress, memo)
}

// SubmitTx submits a Tx on behalf of an ICAAccount to a remote chain.
func (k *Keeper) SubmitTx(ctx sdk.Context, msgs []proto.Message, account *types.ICAAccount, memo string) error {
	// if no messages, do nothing
	if len(msgs) == 0 {
		return nil
	}

	portID := account.GetPortName()
	connectionID, err := k.GetConnectionForPort(ctx, portID)
	if err != nil {
		return err
	}
	channelID, found := k.ICAControllerKeeper.GetActiveChannelID(ctx, connectionID, portID)
	if !found {
		return sdkioerrors.Wrapf(icatypes.ErrActiveChannelNotFound, "failed to retrieve active channel for port %s in submittx", portID)
	}

	chanCap, found := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(portID, channelID))
	if !found {
		return sdkioerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	chunkSize := ICAMsgChunkSize
	timeoutTimestamp := uint64(ctx.BlockTime().Add(24 * time.Hour).UnixNano())

	for {
		// if no messages, no chunks!
		if len(msgs) == 0 {
			break
		}

		// if the last chunk, make chunksize the number of messages
		if len(msgs) < chunkSize {
			chunkSize = len(msgs)
		}
		msgsChunk := msgs[0:chunkSize]

		// build and submit message for this chunk
		data, err := icatypes.SerializeCosmosTx(k.cdc, msgsChunk)
		if err != nil {
			return err
		}

		// validate memo < 256 bytes
		packetData := icatypes.InterchainAccountPacketData{
			Type: icatypes.EXECUTE_TX,
			Data: data,
			Memo: memo,
		}

		_, err = k.ICAControllerKeeper.SendTx(ctx, chanCap, connectionID, portID, packetData, timeoutTimestamp)
		if err != nil {
			return err
		}

		// remove chunk from original msg slice
		msgs = msgs[chunkSize:]
	}

	return nil
}

// ---------------------------------------------------------------

func (k *Keeper) NewReceipt(ctx sdk.Context, zone *types.Zone, sender, txhash string, amount sdk.Coins) *types.Receipt {
	t := ctx.BlockTime()
	return &types.Receipt{ChainId: zone.ChainId, Sender: sender, Txhash: txhash, Amount: amount, FirstSeen: &t}
}

// GetReceipt returns receipt for the given key.
func (k *Keeper) GetReceipt(ctx sdk.Context, key string) (types.Receipt, bool) {
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
func (k *Keeper) DeleteReceipt(ctx sdk.Context, key string) {
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
func (k *Keeper) IterateZoneReceipts(ctx sdk.Context, zone *types.Zone, fn func(index int64, receiptInfo types.Receipt) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixReceipt)
	iterator := sdk.KVStorePrefixIterator(store, []byte(zone.ChainId))
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

	bech32Address, err := bech32.ConvertAndEncode(zone.AccountPrefix, addr)
	if err != nil {
		return receipts, err
	}

	k.IterateZoneReceipts(ctx, zone, func(_ int64, receipt types.Receipt) (stop bool) {
		if receipt.Sender == bech32Address {
			receipts = append(receipts, receipt)
		}
		return false
	})

	return receipts, nil
}

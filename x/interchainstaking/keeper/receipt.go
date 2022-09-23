package keeper

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	icatypes "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v5/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v5/modules/core/24-host"
	abcitypes "github.com/tendermint/tendermint/abci/types"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

const UNSET = "unset"

func (k Keeper) HandleReceiptTransaction(ctx sdk.Context, txr *sdk.TxResponse, txn *tx.Tx, zone types.Zone) {
	k.Logger(ctx).Info("Deposit receipt.", "ischeck", ctx.IsCheckTx(), "isrecheck", ctx.IsReCheckTx())
	hash := txr.TxHash
	memo := txn.Body.Memo

	senderAddress := UNSET
	coins := sdk.Coins{}

	k.Logger(ctx).Error("txr.Events", "events", txr.Events)

	for _, event := range txr.Events {
		if event.Type == "transfer" {
			attrs := attributesToMap(event.Attributes)
			sender := attrs["sender"]
			amount := attrs["amount"]
			if attrs["recipient"] == zone.DepositAddress.GetAddress() { // negate case where sender sends to multiple addresses in one tx
				if senderAddress == UNSET {
					senderAddress = sender
				}

				if sender != senderAddress {
					k.Logger(ctx).Error("sender mismatch", "expected", senderAddress, "received", sender)
				}

				k.Logger(ctx).Error("Deposit receipt", "deposit_address", zone.DepositAddress.GetAddress(), "sender", sender, "amount", amount)
				thisCoins, err := sdk.ParseCoinsNormalized(amount)
				if err != nil {
					k.Logger(ctx).Error("unable to parse coin", "string", amount)
				}
				coins = coins.Add(thisCoins...)
			}
		}
	}

	if senderAddress == UNSET {
		k.Logger(ctx).Error("no sender found. Ignoring.")
		return
	}

	// sdk.AccAddressFromBech32 doesn't work here as it expects the local HRP
	_, addressBytes, err := bech32.DecodeAndConvert(senderAddress)
	if err != nil {
		k.Logger(ctx).Error("unable to decode sender address. Ignoring.", "sender", senderAddress)
		return
	}

	if err := zone.ValidateCoinsForZone(ctx, coins); err != nil {
		// we expect this to trigger if the validatorset has changed recently (i.e. we haven't seen the validator before. That is okay, we'll catch it next round!)
		k.Logger(ctx).Error("unable to validate coins. Ignoring.", "sender", senderAddress)
		return
	}

	var accAddress sdk.AccAddress = addressBytes

	k.Logger(ctx).Info("Found new deposit tx", "deposit_address", zone.DepositAddress.GetAddress(), "sender", senderAddress, "local", accAddress.String(), "chain id", zone.ChainId, "amount", coins, "hash", hash)
	// create receipt

	if err := k.UpdateIntent(ctx, accAddress, zone, coins, memo); err != nil {
		k.Logger(ctx).Error("unable to update intent. Ignoring.", "sender", senderAddress, "zone", zone.ChainId, "err", err)
		return
	}
	if err := k.MintQAsset(ctx, accAddress, zone, coins); err != nil {
		k.Logger(ctx).Error("unable to mint QAsset. Ignoring.", "sender", senderAddress, "zone", zone.ChainId, "err", err)
		return
	}

	if err := k.TransferToDelegate(ctx, zone, coins, hash); err != nil {
		k.Logger(ctx).Error("unable to transfer to delegate. Ignoring.", "sender", senderAddress, "zone", zone.ChainId, "err", err)
		return
	}
	receipt := k.NewReceipt(ctx, zone, senderAddress, hash, coins)

	k.SetReceipt(ctx, *receipt)
}

func attributesToMap(attrs []abcitypes.EventAttribute) map[string]string {
	out := make(map[string]string)
	for _, attr := range attrs {
		out[string(attr.Key)] = string(attr.Value)
	}
	return out
}

func (k *Keeper) MintQAsset(ctx sdk.Context, sender sdk.AccAddress, zone types.Zone, inCoins sdk.Coins) error {
	if zone.RedemptionRate.IsZero() {
		return fmt.Errorf("zero redemption rate")
	}

	outCoins := sdk.Coins{}
	for _, inCoin := range inCoins.Sort() {
		outAmount := sdk.NewDecFromInt(inCoin.Amount).Quo(zone.RedemptionRate).TruncateInt()
		outCoin := sdk.NewCoin(zone.LocalDenom, outAmount)
		outCoins = outCoins.Add(outCoin)
	}
	k.Logger(ctx).Info("Minting qAssets for receipt", "assets", outCoins)
	err := k.BankKeeper.MintCoins(ctx, types.ModuleName, outCoins)
	if err != nil {
		return err
	}

	err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, outCoins)
	if err != nil {
		return err
	}
	k.Logger(ctx).Info("Transferred qAssets to sender", "assets", outCoins, "sender", sender)

	return nil
}

func (k *Keeper) TransferToDelegate(ctx sdk.Context, zone types.Zone, coins sdk.Coins, memo string) error {
	msg := &bankTypes.MsgSend{FromAddress: zone.DepositAddress.GetAddress(), ToAddress: zone.DelegationAddress.GetAddress(), Amount: coins}
	return k.SubmitTx(ctx, []sdk.Msg{msg}, zone.DepositAddress, memo)
}

func (k *Keeper) SubmitTx(ctx sdk.Context, msgs []sdk.Msg, account *types.ICAAccount, memo string) error {
	portID := account.GetPortName()
	connectionID, err := k.GetConnectionForPort(ctx, portID)
	if err != nil {
		return err
	}
	channelID, found := k.ICAControllerKeeper.GetActiveChannelID(ctx, connectionID, portID)
	if !found {
		return sdkerrors.Wrapf(icatypes.ErrActiveChannelNotFound, "failed to retrieve active channel for port %s in submittx", portID)
	}

	chanCap, found := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(portID, channelID))
	if !found {
		return sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	data, err := icatypes.SerializeCosmosTx(k.cdc, msgs)
	if err != nil {
		return err
	}

	// validate memo < 256 bytes
	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
		Memo: memo,
	}

	timeoutTimestamp := uint64(ctx.BlockTime().Add(24 * time.Hour).UnixNano())
	_, err = k.ICAControllerKeeper.SendTx(ctx, chanCap, connectionID, portID, packetData, timeoutTimestamp)
	if err != nil {
		return err
	}

	return nil
}

// ---------------------------------------------------------------

func (k Keeper) NewReceipt(ctx sdk.Context, zone types.Zone, sender string, txhash string, amount sdk.Coins) *types.Receipt {
	return &types.Receipt{ChainId: zone.ChainId, Sender: sender, Txhash: txhash, Amount: amount}
}

// GetReceipt returns receipt
func (k Keeper) GetReceipt(ctx sdk.Context, key string) (types.Receipt, bool) {
	receipt := types.Receipt{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixReceipt)
	bz := store.Get([]byte(key))
	if len(bz) == 0 {
		return receipt, false
	}

	k.cdc.MustUnmarshal(bz, &receipt)
	return receipt, true
}

// SetReceipt set receipt info
func (k Keeper) SetReceipt(ctx sdk.Context, receipt types.Receipt) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixReceipt)
	bz := k.cdc.MustMarshal(&receipt)
	store.Set([]byte(GetReceiptKey(receipt.ChainId, receipt.Txhash)), bz)
}

// DeleteReceipt delete receipt info
func (k Keeper) DeleteReceipt(ctx sdk.Context, key string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixReceipt)
	store.Delete([]byte(key))
}

// IterateQueries iterate through receipts
func (k Keeper) IterateReceipts(ctx sdk.Context, fn func(index int64, receiptInfo types.Receipt) (stop bool)) {
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

func (k Keeper) AllReceipts(ctx sdk.Context) []types.Receipt {
	receipts := make([]types.Receipt, 0)
	k.IterateReceipts(ctx, func(_ int64, receiptInfo types.Receipt) (stop bool) {
		receipts = append(receipts, receiptInfo)
		return false
	})
	return receipts
}

// IterateZoneReceipts iterate through receipts of the given zone
func (k Keeper) IterateZoneReceipts(ctx sdk.Context, zone *types.Zone, fn func(index int64, receiptInfo types.Receipt) (stop bool)) {
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

// UserZoneReceipts returns all receipts of the given user for the given zone
func (k Keeper) UserZoneReceipts(ctx sdk.Context, zone *types.Zone, addr sdk.AccAddress) ([]types.Receipt, error) {
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

func GetReceiptKey(chainID string, txhash string) string {
	return fmt.Sprintf("%s/%s", chainID, txhash)
}

package keeper

import (
	"fmt"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	cosmostx "github.com/cosmos/cosmos-sdk/types/tx"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	tmtypes "github.com/tendermint/tendermint/abci/types"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
)

func (k Keeper) HandleReceiptTransaction(ctx sdk.Context, tx *coretypes.ResultTx, zone types.RegisteredZone) {
	hash := tx.Hash.String()

	_, found := k.GetReceipt(ctx, GetReceiptKey(zone, hash))
	if found {
		k.Logger(ctx).Info("Found previously handled tx. Ignoring.", "txhash", hash)
		return
	}

	var raw cosmostx.TxRaw
	// replicates txDecoder logic, because it's easier than trying to inject the dependency :/
	err := k.cdc.Unmarshal(tx.Tx, &raw)
	if err != nil {
		k.Logger(ctx).Error("Unable to unmarshal tx_raw for deposit account receipt", "deposit_address", zone.DepositAddress.GetAddress(), "tx data", tx.Tx)
	}

	var body cosmostx.TxBody

	err = k.cdc.Unmarshal(raw.BodyBytes, &body)
	if err != nil {
		k.Logger(ctx).Error("Unable to unmarshal tx_body for deposit account receipt", "deposit_address", zone.DepositAddress.GetAddress(), "tx body", raw.BodyBytes)
	}

	senderAddress := "unset"
	coins := sdk.Coins{}

	for _, event := range tx.TxResult.Events {
		if event.Type == "transfer" {
			attrs := attributesToMap(event.Attributes)
			sender := attrs["sender"]
			amount := attrs["amount"]
			if attrs["recipient"] == zone.DepositAddress.GetAddress() { // negate case where sender sends to multiple addresses in one tx
				if senderAddress == "unset" {
					senderAddress = sender
				}

				if sender != senderAddress {
					k.Logger(ctx).Error("Sender mismatch", "expected", senderAddress, "received", sender)
				}

				k.Logger(ctx).Info("Deposit receipt", "deposit_address", zone.DepositAddress.GetAddress(), "sender", sender, "amount", amount)
				thisCoins, err := sdk.ParseCoinsNormalized(amount)
				if err != nil {
					k.Logger(ctx).Error("Unable to parse coin", "string", amount)
				}
				coins = coins.Add(thisCoins...)
			}
		}
	}

	if senderAddress == "unset" {
		k.Logger(ctx).Error("No sender found. Ignoring.")
		return
	}

	// sdk.AccAddressFromBech32 doesn't work here as it expects the local HRP
	_, addressBytes, err := bech32.DecodeAndConvert(senderAddress)
	if err != nil {
		k.Logger(ctx).Error("Unable to decode sender address. Ignoring.", "sender", senderAddress)
		return
	}

	if err := zone.ValidateCoinsForZone(ctx, coins); err != nil {
		// we expect this to trigger if the validatorset has changed recently (i.e. we haven't seen the validator before. That is okay, we'll catch it next round!)
		k.Logger(ctx).Error("Unable to validate coins. Ignoring.", "sender", senderAddress)
		return
	}

	var accAddress sdk.AccAddress = addressBytes

	k.Logger(ctx).Info("Found new deposit tx", "deposit_address", zone.DepositAddress.GetAddress(), "sender", senderAddress, "local", accAddress.String(), "chain id", zone.ChainId, "amount", coins, "hash", hash)
	// create receipt
	receipt := k.NewReceipt(ctx, zone, senderAddress, hash, coins)

	k.UpdateIntent(ctx, accAddress, zone, coins)
	k.MintQAsset(ctx, accAddress, zone, coins)
	k.TransferToDelegate(ctx, zone, coins)
	k.SetReceipt(ctx, *receipt)
}

func attributesToMap(attrs []tmtypes.EventAttribute) map[string]string {
	out := make(map[string]string)
	for _, attr := range attrs {
		out[string(attr.Key)] = string(attr.Value)
	}
	return out
}

func (k *Keeper) MintQAsset(ctx sdk.Context, sender sdk.AccAddress, zone types.RegisteredZone, inCoins sdk.Coins) error {
	outCoins := sdk.Coins{}
	for _, inCoin := range inCoins {
		outAmount := inCoin.Amount.ToDec().Quo(zone.RedemptionRate).TruncateInt()
		outCoin := sdk.NewCoin(zone.LocalDenom, outAmount)
		outCoins = outCoins.Add(outCoin)
	}
	k.Logger(ctx).Info("Minting qAssets for receipt", "assets", outCoins)
	err := k.BankKeeper.MintCoins(ctx, types.ModuleName, outCoins)
	if err != nil {
		panic(err)
	}

	err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, outCoins)
	if err != nil {
		panic(err)
	}
	k.Logger(ctx).Info("Transferred qAssets to sender", "assets", outCoins, "sender", sender)

	return nil
}

func (k *Keeper) TransferToDelegate(ctx sdk.Context, zone types.RegisteredZone, inAmount sdk.Coins) error {
	if zone.SupportMultiSend() {
		return k.TransferToDelegateMulti(ctx, zone, inAmount)
	} else {
		eachAmount := sdk.Coins{}
		splits := int64(math.Min(types.DelegationAccountSplit, float64(len(zone.DelegationAddresses))))

		for _, asset := range inAmount {
			thisAsset := sdk.Coin{Denom: asset.Denom, Amount: asset.Amount.Quo(sdk.NewInt(splits))}
			eachAmount = eachAmount.Add(thisAsset)
		}
		accounts := zone.GetDelegationAccountsByLowestBalance(splits)

		// reverse this slice; because we want the smallest account last, so the remainder ends in it.
		for i, j := 0, len(accounts)-1; i < j; i, j = i+1, j-1 {
			accounts[i], accounts[j] = accounts[j], accounts[i]
		}

		var msgs []sdk.Msg

		for idx, account := range accounts {
			var amount sdk.Coins
			if idx < len(accounts)-1 {
				// if not the last account in sequence then sub each amount from inAmount
				amount = eachAmount
				inAmount = inAmount.Sub(amount)
			} else {
				amount = inAmount
			}
			msgs = append(msgs, &bankTypes.MsgSend{FromAddress: zone.DepositAddress.GetAddress(), ToAddress: account.GetAddress(), Amount: amount})
		}

		return k.SubmitTx(ctx, msgs, zone.DepositAddress)
	}
}

func (k *Keeper) TransferToDelegateMulti(ctx sdk.Context, zone types.RegisteredZone, inAmount sdk.Coins) error {
	eachAmount := sdk.Coins{}
	splits := int64(math.Min(types.DelegationAccountSplit, float64(len(zone.DelegationAddresses))))
	for _, asset := range inAmount {
		thisAsset := sdk.Coin{Denom: asset.Denom, Amount: asset.Amount.Quo(sdk.NewInt(splits))}
		eachAmount = eachAmount.Add(thisAsset)
	}

	in := []bankTypes.Input{}
	out := []bankTypes.Output{}

	in = append(in, bankTypes.Input{Address: zone.DepositAddress.GetAddress(), Coins: inAmount})

	accounts := zone.GetDelegationAccountsByLowestBalance(splits)
	for _, account := range accounts {
		out = append(out, bankTypes.Output{Address: account.GetAddress(), Coins: eachAmount})
		inAmount = inAmount.Sub(eachAmount)
	}

	// ensure any remainder gets deposited in the first account (as it will have the lowest balance)
	out[0].Coins.Add(inAmount...)

	sendMsg := bankTypes.NewMsgMultiSend(in, out)
	var msg sdk.Msg = sendMsg
	// send from deposit to accounts

	return k.SubmitTx(ctx, []sdk.Msg{msg}, zone.DepositAddress)
}

func (k *Keeper) SubmitTx(ctx sdk.Context, msgs []sdk.Msg, account *types.ICAAccount) error {

	portID := account.GetPortName()
	connectionID, err := k.GetConnectionForPort(ctx, portID)
	if err != nil {
		return err
	}

	channelID, found := k.ICAControllerKeeper.GetActiveChannelID(ctx, connectionID, portID)
	if !found {
		return sdkerrors.Wrapf(icatypes.ErrActiveChannelNotFound, "failed to retrieve active channel for port %s", portID)
	}

	chanCap, found := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(portID, channelID))
	if !found {
		return sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	data, err := icatypes.SerializeCosmosTx(k.cdc, msgs)
	if err != nil {
		return err
	}

	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
	}

	// timeoutTimestamp set to max value with the unsigned bit shifted to sastisfy hermes timestamp conversion
	// it is the responsibility of the auth module developer to ensure an appropriate timeout timestamp
	timeoutTimestamp := ^uint64(0) >> 1
	_, err = k.ICAControllerKeeper.SendTx(ctx, chanCap, connectionID, portID, packetData, timeoutTimestamp)
	if err != nil {
		return err
	}

	return nil
}

// ---------------------------------------------------------------

func (k Keeper) NewReceipt(ctx sdk.Context, zone types.RegisteredZone, sender string, txhash string, amount sdk.Coins) *types.Receipt {
	return &types.Receipt{Zone: &zone, Sender: sender, Txhash: txhash, Amount: amount}
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
	store.Set([]byte(GetReceiptKey(*receipt.Zone, receipt.Txhash)), bz)
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

func GetReceiptKey(zone types.RegisteredZone, txhash string) string {
	return fmt.Sprintf("%s/%s", zone.ChainId, txhash)
}

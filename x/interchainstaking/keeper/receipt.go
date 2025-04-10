package keeper

import (
	"errors"
	"fmt"
	"time"

	"github.com/cosmos/gogoproto/proto"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	icacontrollerkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"

	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
	minttypes "github.com/quicksilver-zone/quicksilver/x/mint/types"
	prtypes "github.com/quicksilver-zone/quicksilver/x/participationrewards/types"
)

const (
	Unset           = "unset"
	ICAMsgChunkSize = 5
	ICATimeout      = time.Hour * 6
	// AuthzAutoClaimAddress = "quick1psevptdp90jad76zt9y9x2nga686hutgmasmwd"
)

func (k *Keeper) HandleReceiptTransaction(ctx sdk.Context, txn *tx.Tx, hash string, zone types.Zone) error {
	k.Logger(ctx).Info("Deposit receipt.", "ischeck", ctx.IsCheckTx(), "isrecheck", ctx.IsReCheckTx())
	memo := txn.Body.Memo

	senderAddress := Unset
	assets := sdk.Coins{}

	for _, msg := range txn.GetMsgs() {
		msgSend, ok := msg.(*banktypes.MsgSend)
		if !ok {
			k.Logger(ctx).Error("got message that wasn't MsgSend!")
			continue
		}
		sender := msgSend.FromAddress
		amount := msgSend.Amount

		if msgSend.ToAddress == zone.DepositAddress.GetAddress() { // negate case where sender sends to multiple addresses in one tx
			if senderAddress == Unset {
				senderAddress = sender
			}

			if sender != senderAddress {
				// TODO: technically nothing wrong with this; just need to make sure we _only_ consider assets sent to QS
				k.Logger(ctx).Error("sender mismatch", "expected", senderAddress, "received", sender)
				k.NilReceipt(ctx, &zone, hash) // nil receipt will stop this hash being submitted again
				return k.SendToWithdrawal(ctx, &zone, zone.DepositAddress, assets)
			}

			k.Logger(ctx).Info("Deposit receipt", "deposit_address", zone.DepositAddress.GetAddress(), "sender", sender, "amount", amount)

			assets = assets.Add(amount...)
		}

	}

	if senderAddress == Unset { // not sure this is ever reachable. A valid MsgSend must always have a valid sender.
		k.Logger(ctx).Error("no sender found. Ignoring.")
		k.NilReceipt(ctx, &zone, hash) // nil receipt will stop this hash being submitted again
		return k.SendToWithdrawal(ctx, &zone, zone.DepositAddress, assets)

	}
	senderAccAddress, err := addressutils.AccAddressFromBech32(senderAddress, zone.GetAccountPrefix())
	if err != nil { // not sure this is ever reachable. A valid MsgSend must always have a valid sender.
		k.Logger(ctx).Error("unable to decode sender address. Ignoring.", "senderAddress", senderAddress, "error", err)
		k.NilReceipt(ctx, &zone, hash) // nil receipt will stop this hash being submitted again
		return k.SendToWithdrawal(ctx, &zone, zone.DepositAddress, assets)

	}

	valid, matchesVals := zone.ValidateCoinsForZone(assets, k.GetValidatorAddressesAsMap(ctx, zone.ChainId))

	if !valid {
		k.Logger(ctx).Error("unable to validate coins. Ignoring.", "senderAddress", senderAddress)
		k.NewCompletedReceipt(ctx, &zone, senderAddress, hash, assets) // nil receipt will stop this hash being submitted again
		// send tokens to withdrawal for disbursal.
		return k.SendToWithdrawal(ctx, &zone, zone.DepositAddress, assets)
	} else if !matchesVals {
		k.Logger(ctx).Error("unable to validate coins for this valset.", "senderAddress", senderAddress)
		// Do not set a nil receipt so we can revisit this tx.
		// Don't return an error as to not clog queue.
		return nil
	}

	k.Logger(ctx).Info("found new deposit tx", "deposit_address", zone.DepositAddress.GetAddress(), "senderAddress", senderAddress, "local", senderAccAddress.String(), "chain id", zone.ChainId, "assets", assets, "hash", hash)

	var (
		memoIntent    types.ValidatorIntents
		memoFields    types.MemoFields
		memoRTS       bool
		memoAutoClaim bool
		mappedAddress []byte
	)

	if len(memo) > 0 {
		// process memo
		memoFields, err = zone.DecodeMemo(memo)
		if err != nil {
			// What should we do on error here? just log?
			k.Logger(ctx).Error("error decoding memo", "error", err.Error(), "memo", memo)
		}
		memoRTS = memoFields.RTS()
		mappedAddress, _ = memoFields.AccountMap()
		memoIntent, _ = memoFields.Intent(assets, &zone)
		memoAutoClaim = memoFields.AutoClaim()
	}

	// update state
	if err := k.UpdateDelegatorIntent(ctx, senderAccAddress, &zone, assets, memoIntent); err != nil {
		k.Logger(ctx).Error("unable to update intent. Ignoring.", "senderAddress", senderAddress, "zone", zone.ChainId, "err", err.Error())
		return fmt.Errorf("unable to update intent. Ignoring. senderAddress=%q zone=%q err: %w", senderAddress, zone.ChainId, err)
	}

	success, err := k.MintAndSendQAsset(ctx, senderAccAddress, senderAddress, &zone, assets, memoRTS, mappedAddress)
	if err != nil {
		k.Logger(ctx).Error("unable to mint QAsset. Ignoring.", "senderAddress", senderAddress, "zone", zone.ChainId, "err", err)
		return fmt.Errorf("unable to mint QAsset. Ignoring. senderAddress=%q zone=%q err: %w", senderAddress, zone.ChainId, err)
	}

	if success {
		if err := k.TransferToDelegate(ctx, &zone, assets, hash); err != nil {
			k.Logger(ctx).Error("unable to transfer to delegate. Ignoring.", "senderAddress", senderAddress, "zone", zone.ChainId, "err", err)
			return fmt.Errorf("unable to transfer to delegate. Ignoring. senderAddress=%q zone=%q err: %w", senderAddress, zone.ChainId, err)
		}
		if memoAutoClaim {
			if err := k.HandleAutoClaim(ctx, senderAccAddress); err != nil {
				k.Logger(ctx).Error("unable to handle auto claim. Ignoring.", "senderAddress", senderAddress, "zone", zone.ChainId, "err", err)
				return fmt.Errorf("unable to handle auto claim. Ignoring. senderAddress=%q zone=%q err: %w", senderAddress, zone.ChainId, err)
			}
		}
	}
	// create receipt
	receipt := k.NewReceipt(ctx, &zone, senderAddress, hash, assets)
	k.SetReceipt(ctx, *receipt)

	return nil
}

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

func (k *Keeper) HandleAutoClaim(ctx sdk.Context, senderAddress sdk.AccAddress) error {
	authzAutoClaimAddress := k.GetAuthzAutoClaimAddress(ctx)
	if authzAutoClaimAddress == "" {
		return errors.New("no auto claim address set")
	}

	return k.AuthzKeeper.SaveGrant(
		ctx,
		addressutils.MustAccAddressFromBech32(authzAutoClaimAddress, "quick"),
		senderAddress,
		&authz.GenericAuthorization{
			Msg: sdk.MsgTypeURL(&prtypes.MsgSubmitClaim{}),
		},
		nil,
	)
}

// MintAndSendQAsset mints qAssets based on the native asset redemption rate.  Tokens are then transferred to the given user.
// The function handles the following cases:
//  1. If the zone is labeled "return to sender" or the Tx memo contains "return to sender" flag:
//     - Mint QAssets and IBC transfer to the corresponding zone acc
//  2. If there is no mapped account but the zone is labeled as non-118 coin type:
//     - Do not mint QAssets and refund assets
//  3. If a mapped account is set for a non-118 coin type zone:
//     - Mint QAssets and send to corresponding mapped address
//  4. If a new mapped account is provided to the function and the zone is labeled as non-118 coin type:
//     - Mint QAssets, set new mapping for the mapped account in the keeper, and send to corresponding mapped account.
//  5. If the zone is 118 and no other flags are set:
//     - Mint QAssets and transfer to send to msg creator.
func (k *Keeper) MintAndSendQAsset(ctx sdk.Context, sender sdk.AccAddress, senderAddress string, zone *types.Zone, assets sdk.Coins, memoRTS bool, mappedAddress sdk.AccAddress) (bool, error) {
	if zone.RedemptionRate.IsZero() {
		return false, errors.New("zero redemption rate")
	}

	qAssets := sdk.Coins{}
	for _, asset := range assets.Sort() {
		amount := sdk.NewDecFromInt(asset.Amount).Quo(zone.RedemptionRate).TruncateInt()
		qAssets = qAssets.Add(sdk.NewCoin(zone.LocalDenom, amount))
	}

	// check if a remote address exists for a non 118 coin type zone
	setMappedAddress := true
	if mappedAddress == nil && !zone.Is_118 && !zone.ReturnToSender && !memoRTS {
		var found bool
		mappedAddress, found = k.GetRemoteAddressMap(ctx, sender, zone.ChainId)
		if !found {
			// if not found, skip minting and refund assets
			msg := &banktypes.MsgSend{FromAddress: zone.DepositAddress.GetAddress(), ToAddress: senderAddress, Amount: assets}
			return false, k.SubmitTx(ctx, []sdk.Msg{msg}, zone.DepositAddress, "refund", zone.MessagesPerTx)
		}
		// do not set, since mapped address already exists
		setMappedAddress = false
	}

	k.Logger(ctx).Info("Minting qAssets for receipt", "assets", qAssets)
	err := k.BankKeeper.MintCoins(ctx, types.ModuleName, qAssets)
	if err != nil {
		return false, err
	}

	switch {
	case zone.ReturnToSender || memoRTS:
		err = k.SendTokenIBC(ctx, k.AccountKeeper.GetModuleAddress(types.ModuleName), senderAddress, zone, qAssets[0])
		k.Logger(ctx).Info("Transferred qAssets via rts", "address", senderAddress, "assets", qAssets)

	case mappedAddress != nil:
		// set mapped account
		if setMappedAddress {
			k.SetAddressMapPair(ctx, mappedAddress, sender, zone.ChainId)
		}

		// set send to mapped account
		err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, mappedAddress, qAssets)
		k.Logger(ctx).Info("Transferred qAssets to mapped account", "address", mappedAddress, "assets", qAssets)
	default:
		err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, qAssets)
		k.Logger(ctx).Info("Transferred qAssets to sender", "address", sender, "assets", qAssets)

	}

	if err != nil {
		return false, fmt.Errorf("unable to transfer coins: %w", err)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			minttypes.EventTypeMint,
			sdk.NewAttribute(sdk.AttributeKeyAmount, qAssets.String()),
		),
	)
	return true, nil
}

// TransferToDelegate transfers tokens from the zone deposit account address to the zone delegate account address.
func (k *Keeper) TransferToDelegate(ctx sdk.Context, zone *types.Zone, coins sdk.Coins, memo string) error {
	msg := &banktypes.MsgSend{FromAddress: zone.DepositAddress.GetAddress(), ToAddress: zone.DelegationAddress.GetAddress(), Amount: coins}
	return k.SubmitTx(ctx, []sdk.Msg{msg}, zone.DepositAddress, memo, zone.MessagesPerTx)
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

		// DEBUG:
		activeChannelID, found := k.ICAControllerKeeper.GetOpenActiveChannel(ctx, connectionID, portID)
		if !found {
			k.Logger(ctx).Error("DEBUG - NO ACTIVE CHANNEL FOUND")
		}
		module, cap, err := k.IBCKeeper.ChannelKeeper.LookupModuleByChannel(ctx, portID, activeChannelID)
		k.Logger(ctx).Error("DEBUG - CAPS FOR CHANNEL", "module", module)
		k.Logger(ctx).Error("DEBUG - CAPS FOR CHANNEL", "cap", cap)
		k.Logger(ctx).Error("DEBUG - CAPS FOR CHANNEL", "err", err)

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

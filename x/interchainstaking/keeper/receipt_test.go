package keeper_test

import (
	"fmt"
	"time"

	"cosmossdk.io/math"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"

	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	"github.com/quicksilver-zone/quicksilver/utils/randomutils"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

func (suite *KeeperTestSuite) TestHandleReceiptTransactionGood() {
	suite.SetupTest()
	suite.setupTestZones()

	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

	// get test zone
	zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.True(found)

	fromAddress := addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix)

	msg := banktypes.MsgSend{FromAddress: fromAddress, ToAddress: zone.DepositAddress.Address, Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(1000000)))}
	anymsg, err := codectypes.NewAnyWithValue(&msg)
	suite.NoError(err)

	transaction := &tx.Tx{Body: &tx.TxBody{Messages: []*codectypes.Any{anymsg}}}
	hash := randomutils.GenerateRandomHashAsHex(64)
	hash2 := randomutils.GenerateRandomHashAsHex(64)

	before := suite.GetQuicksilverApp(suite.chainA).BankKeeper.GetSupply(ctx, zone.LocalDenom)
	suite.Equal(sdk.NewCoin(zone.LocalDenom, sdk.ZeroInt()), before)
	// rr is 1.0
	err = icsKeeper.HandleReceiptTransaction(ctx, transaction, hash, zone)
	suite.NoError(err)

	after := suite.GetQuicksilverApp(suite.chainA).BankKeeper.GetSupply(ctx, zone.LocalDenom)
	suite.Equal(sdk.NewCoin(zone.LocalDenom, math.NewInt(1000000)), after)

	zone.RedemptionRate = sdk.NewDecWithPrec(12, 1)
	err = icsKeeper.HandleReceiptTransaction(ctx, transaction, hash2, zone)
	suite.NoError(err)

	after2 := suite.GetQuicksilverApp(suite.chainA).BankKeeper.GetSupply(ctx, zone.LocalDenom)
	suite.Equal(sdk.NewCoin(zone.LocalDenom, math.NewInt(1833333)), after2)
}

func (suite *KeeperTestSuite) TestHandleReceiptTransactionBadRecipient() {
	suite.SetupTest()
	suite.setupTestZones()

	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

	// get test zone
	zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.True(found)

	fromAddress := addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix)

	msg := banktypes.MsgSend{FromAddress: fromAddress, ToAddress: zone.DelegationAddress.Address, Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(1000000)))}
	anymsg, err := codectypes.NewAnyWithValue(&msg)
	suite.NoError(err)

	transaction := &tx.Tx{Body: &tx.TxBody{Messages: []*codectypes.Any{anymsg}}}
	hash := randomutils.GenerateRandomHashAsHex(64)

	before := suite.GetQuicksilverApp(suite.chainA).BankKeeper.GetSupply(ctx, zone.LocalDenom)
	suite.Equal(sdk.NewCoin(zone.LocalDenom, sdk.ZeroInt()), before)

	err = icsKeeper.HandleReceiptTransaction(ctx, transaction, hash, zone)
	// suite.ErrorContains(err, "no sender found. Ignoring")
	nilReceipt, found := icsKeeper.GetReceipt(ctx, zone.ChainId, hash)
	suite.True(found)                  // check nilReceipt is found for hash
	suite.Equal("", nilReceipt.Sender) // check nilReceipt has empty sender
	suite.Nil(nilReceipt.Amount)       // check nilReceipt has nil amount
	suite.NoError(err)
}

func (suite *KeeperTestSuite) TestHandleReceiptTransactionBadMessageType() {
	suite.SetupTest()
	suite.setupTestZones()

	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

	// get test zone
	zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.True(found)

	fromAddress := addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix)

	msg := stakingtypes.MsgDelegate{DelegatorAddress: fromAddress, ValidatorAddress: zone.DelegationAddress.Address, Amount: sdk.NewCoin(zone.BaseDenom, math.NewInt(1000000))}
	anymsg, err := codectypes.NewAnyWithValue(&msg)
	suite.NoError(err)

	transaction := &tx.Tx{Body: &tx.TxBody{Messages: []*codectypes.Any{anymsg}}}
	hash := randomutils.GenerateRandomHashAsHex(64)

	before := suite.GetQuicksilverApp(suite.chainA).BankKeeper.GetSupply(ctx, zone.LocalDenom)
	suite.Equal(sdk.NewCoin(zone.LocalDenom, sdk.ZeroInt()), before)

	err = icsKeeper.HandleReceiptTransaction(ctx, transaction, hash, zone)
	// suite.ErrorContains(err, "no sender found. Ignoring")
	nilReceipt, found := icsKeeper.GetReceipt(ctx, zone.ChainId, hash)
	suite.True(found)                  // check nilReceipt is found for hash
	suite.Equal("", nilReceipt.Sender) // check nilReceipt has empty sender
	suite.Nil(nilReceipt.Amount)       // check nilReceipt has nil amount
	suite.NoError(err)
}

func (suite *KeeperTestSuite) TestHandleReceiptMixedMessageTypeGood() {
	suite.SetupTest()
	suite.setupTestZones()

	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

	// get test zone
	zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.True(found)

	fromAddress := addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix)

	msg := banktypes.MsgSend{FromAddress: fromAddress, ToAddress: zone.DepositAddress.Address, Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(1000000)))}
	anymsg, err := codectypes.NewAnyWithValue(&msg)
	suite.NoError(err)

	msg2 := stakingtypes.MsgDelegate{DelegatorAddress: fromAddress, ValidatorAddress: zone.DelegationAddress.Address, Amount: sdk.NewCoin(zone.BaseDenom, math.NewInt(1000000))}
	anymsg2, err := codectypes.NewAnyWithValue(&msg2)
	suite.NoError(err)

	transaction := &tx.Tx{Body: &tx.TxBody{Messages: []*codectypes.Any{anymsg, anymsg2}}}
	hash := randomutils.GenerateRandomHashAsHex(64)

	before := suite.GetQuicksilverApp(suite.chainA).BankKeeper.GetSupply(ctx, zone.LocalDenom)
	suite.Equal(sdk.NewCoin(zone.LocalDenom, sdk.ZeroInt()), before)

	err = icsKeeper.HandleReceiptTransaction(ctx, transaction, hash, zone)
	suite.NoError(err)

	after := suite.GetQuicksilverApp(suite.chainA).BankKeeper.GetSupply(ctx, zone.LocalDenom)
	suite.Equal(sdk.NewCoin(zone.LocalDenom, math.NewInt(1000000)), after)
}

func (suite *KeeperTestSuite) TestHandleReceiptTransactionBadMixedSender() { // this shouldn't be possibly in theory, but hey!
	suite.SetupTest()
	suite.setupTestZones()

	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

	// get test zone
	zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.True(found)

	fromAddress := addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix)
	fromAddress2 := addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix)

	msg := banktypes.MsgSend{FromAddress: fromAddress, ToAddress: zone.DepositAddress.Address, Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(1000000)))}
	anymsg, err := codectypes.NewAnyWithValue(&msg)
	suite.NoError(err)
	msg2 := banktypes.MsgSend{FromAddress: fromAddress2, ToAddress: zone.DepositAddress.Address, Amount: sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, math.NewInt(1000000)))}
	anymsg2, err := codectypes.NewAnyWithValue(&msg2)
	suite.NoError(err)

	transaction := &tx.Tx{Body: &tx.TxBody{Messages: []*codectypes.Any{anymsg, anymsg2}}}
	hash := randomutils.GenerateRandomHashAsHex(64)

	before := suite.GetQuicksilverApp(suite.chainA).BankKeeper.GetSupply(ctx, zone.LocalDenom)
	suite.Equal(sdk.NewCoin(zone.LocalDenom, sdk.ZeroInt()), before)

	err = icsKeeper.HandleReceiptTransaction(ctx, transaction, hash, zone)
	// suite.ErrorContains(err, "sender mismatch: expected")
	nilReceipt, found := icsKeeper.GetReceipt(ctx, zone.ChainId, hash)
	suite.True(found)                  // check nilReceipt is found for hash
	suite.Equal("", nilReceipt.Sender) // check nilReceipt has empty sender
	suite.Nil(nilReceipt.Amount)       // check nilReceipt has nil amount
	suite.NoError(err)
}

func (suite *KeeperTestSuite) TestHandleReceiptTransactionBadDenom() {
	suite.SetupTest()
	suite.setupTestZones()

	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

	// get test zone
	zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.True(found)

	fromAddress := addressutils.GenerateAddressForTestWithPrefix(zone.AccountPrefix)

	msg := banktypes.MsgSend{FromAddress: fromAddress, ToAddress: zone.DepositAddress.Address, Amount: sdk.NewCoins(sdk.NewCoin("ushit", math.NewInt(1000000)))}
	anymsg, err := codectypes.NewAnyWithValue(&msg)
	suite.NoError(err)

	transaction := &tx.Tx{Body: &tx.TxBody{Messages: []*codectypes.Any{anymsg}}}
	hash := randomutils.GenerateRandomHashAsHex(64)

	before := suite.GetQuicksilverApp(suite.chainA).BankKeeper.GetSupply(ctx, zone.LocalDenom)
	suite.Equal(sdk.NewCoin(zone.LocalDenom, sdk.ZeroInt()), before)

	err = icsKeeper.HandleReceiptTransaction(ctx, transaction, hash, zone)
	suite.NoError(err)

	after := suite.GetQuicksilverApp(suite.chainA).BankKeeper.GetSupply(ctx, zone.LocalDenom)
	suite.Equal(sdk.NewCoin(zone.LocalDenom, sdk.ZeroInt()), after)
}

// func (suite *KeeperTestSuite) TestMintQAsset() {
// }

// test all getters, setters, deleters, iterators.
func (suite *KeeperTestSuite) TestReceiptStore() {
	suite.SetupTest()
	suite.setupTestZones()

	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

	// get test zone
	zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.True(found)

	account1 := addressutils.GenerateAccAddressForTest()
	account2 := addressutils.GenerateAccAddressForTest()
	hash1 := randomutils.GenerateRandomHashAsHex(64)
	hash2 := randomutils.GenerateRandomHashAsHex(64)
	hash3 := randomutils.GenerateRandomHashAsHex(64)
	hash4 := randomutils.GenerateRandomHashAsHex(64)

	zone2 := types.Zone{ChainId: "test-1"}

	suite.Zero(len(icsKeeper.AllReceipts(ctx)))

	receipt1 := icsKeeper.NewReceipt(ctx, &zone, account1.String(), hash1, sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(100))))
	receipt2 := icsKeeper.NewReceipt(ctx, &zone, account1.String(), hash2, sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(200))))
	receipt3 := icsKeeper.NewReceipt(ctx, &zone, account2.String(), hash3, sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(300))))
	receipt4 := icsKeeper.NewReceipt(ctx, &zone2, account2.String(), hash4, sdk.NewCoins(sdk.NewCoin("uosmo", math.NewInt(500))))

	icsKeeper.SetReceipt(ctx, *receipt1)
	icsKeeper.SetReceipt(ctx, *receipt2)
	icsKeeper.SetReceipt(ctx, *receipt3)
	icsKeeper.SetReceipt(ctx, *receipt4)

	suite.Equal(4, len(icsKeeper.AllReceipts(ctx)))

	count := 0
	coins := sdk.Coins{}
	icsKeeper.IterateReceipts(ctx, func(index int64, receiptInfo types.Receipt) (stop bool) {
		count++
		coins = coins.Add(receiptInfo.Amount...)
		return false
	})

	suite.Equal(4, count)
	suite.Equal(600, int(coins.AmountOf("uatom").Int64()))
	suite.Equal(500, int(coins.AmountOf("uosmo").Int64()))

	count = 0
	sum := 0
	icsKeeper.IterateZoneReceipts(ctx, zone.ChainId, func(index int64, receiptInfo types.Receipt) (stop bool) {
		count++
		sum += int(receiptInfo.Amount.AmountOf("uatom").Int64())
		return false
	})

	suite.Equal(3, count)
	suite.Equal(600, sum)

	count = 0
	sum = 0
	icsKeeper.IterateZoneReceipts(ctx, zone2.ChainId, func(index int64, receiptInfo types.Receipt) (stop bool) {
		count++
		sum += int(receiptInfo.Amount.AmountOf("uosmo").Int64())
		return false
	})

	suite.Equal(1, count)
	suite.Equal(500, sum)

	out, err := icsKeeper.UserZoneReceipts(ctx, &zone, account1)
	suite.NoError(err)
	suite.Equal(2, len(out))

	receipt, found := icsKeeper.GetReceipt(ctx, zone.ChainId, hash1)
	suite.True(found)
	suite.Equal(receipt1, &receipt)
	now := ctx.BlockTime().Add(time.Second)
	receipt.Completed = &now
	icsKeeper.SetReceipt(ctx, receipt)
	icsKeeper.DeleteReceipt(ctx, zone.ChainId, hash2)

	out, err = icsKeeper.UserZoneReceipts(ctx, &zone, account1)
	suite.NoError(err)
	suite.Equal(1, len(out))
	suite.Equal(&now, out[0].Completed)

	icsKeeper.SetReceiptsCompleted(ctx, zone.ChainId, now, now, zone.BaseDenom)

	receipt, found = icsKeeper.GetReceipt(ctx, zone.ChainId, hash3)
	suite.True(found)

	suite.Equal(&now, receipt.Completed)
}

func (suite *KeeperTestSuite) TestSendTokenIBC() {
	suite.Run("test", func() {
		suite.SetupTest()

		// setup transfer channel
		suite.path.EndpointA.ChannelConfig.Version = transfertypes.Version
		suite.path.EndpointB.ChannelConfig.Version = transfertypes.Version

		// set the counter party channel sequence to a different channel to prevent regression
		suite.GetQuicksilverApp(suite.chainB).IBCKeeper.ChannelKeeper.SetNextChannelSequence(suite.chainB.GetContext(), 10)

		suite.coordinator.CreateTransferChannels(suite.path)
		suite.setupTestZones()

		quicksilver := suite.GetQuicksilverApp(suite.chainA)
		ctx := suite.chainA.GetContext()

		zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
		suite.True(found)

		sender := suite.chainA.SenderAccount.GetAddress()
		receiver := addressutils.GenerateAddressForTestWithPrefix("cosmos")

		amount := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))
		err := quicksilver.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(amount))
		suite.NoError(err)
		err = quicksilver.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.NewCoins(amount))
		suite.NoError(err)

		// Try to send native token but wrong connection id on zone
		wrongZone := zone
		wrongZone.ConnectionId = "connection-10"
		err = quicksilver.InterchainstakingKeeper.SendTokenIBC(ctx, sender, receiver, &wrongZone, amount)
		suite.ErrorContains(err, "unable to find remote transfer connection")

		// Try to send the native token
		err = quicksilver.InterchainstakingKeeper.SendTokenIBC(ctx, sender, receiver, &zone, amount)
		suite.NoError(err)

		portID := types.TransferPort
		channelID := suite.path.EndpointA.ChannelID

		ibcAmount := transfertypes.GetTransferCoin(portID, channelID, sdk.DefaultBondDenom, sdk.NewInt(100))

		err = quicksilver.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(ibcAmount))
		suite.NoError(err)
		err = quicksilver.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.NewCoins(ibcAmount))
		suite.NoError(err)

		quicksilver.TransferKeeper.SetDenomTrace(
			ctx,
			transfertypes.DenomTrace{
				Path:      fmt.Sprintf("%s/%s", portID, channelID),
				BaseDenom: sdk.DefaultBondDenom,
			},
		)

		// Try to send the ibc token
		err = quicksilver.InterchainstakingKeeper.SendTokenIBC(ctx, sender, receiver, &zone, ibcAmount)
		suite.NoError(err)
	})
}

func (suite *KeeperTestSuite) TestMintAndSendQAsset1RR() {
	suite.SetupTest()
	suite.setupTestZones()

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()

	zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.True(found)

	senderAddress := addressutils.GenerateAddressForTestWithPrefix("cosmos")
	sender := addressutils.MustAccAddressFromBech32(senderAddress, "")

	amount := sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(5000)))

	// Test sending QAsset
	err := quicksilver.InterchainstakingKeeper.MintAndSendQAsset(ctx, sender, senderAddress, &zone, amount, false, nil)
	suite.NoError(err)

	// Verify balance of receiver
	receiverBalance := quicksilver.BankKeeper.GetBalance(ctx, sender, zone.LocalDenom)
	suite.Equal(sdk.NewCoin(zone.LocalDenom, sdk.NewInt(5000)), receiverBalance)
}

func (suite *KeeperTestSuite) TestMintAndSendQAssetNon1RR() {
	suite.SetupTest()
	suite.setupTestZones()

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()

	zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	zone.RedemptionRate = sdk.NewDecWithPrec(110, 2)
	suite.True(found)

	senderAddress := addressutils.GenerateAddressForTestWithPrefix("cosmos")
	sender := addressutils.MustAccAddressFromBech32(senderAddress, "")

	amount := sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(5000)))

	// Test sending QAsset
	err := quicksilver.InterchainstakingKeeper.MintAndSendQAsset(ctx, sender, senderAddress, &zone, amount, false, nil)
	suite.NoError(err)

	// Verify balance of receiver
	receiverBalance := quicksilver.BankKeeper.GetBalance(ctx, sender, zone.LocalDenom)
	suite.Equal(sdk.NewCoin(zone.LocalDenom, sdk.NewInt(4545)), receiverBalance)
}

func (suite *KeeperTestSuite) TestMintAndSendQAssetSub1RR() {
	suite.SetupTest()
	suite.setupTestZones()

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()

	zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	zone.RedemptionRate = sdk.NewDecWithPrec(90, 2)
	suite.True(found)

	senderAddress := addressutils.GenerateAddressForTestWithPrefix("cosmos")
	sender := addressutils.MustAccAddressFromBech32(senderAddress, "")

	amount := sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(5000)))

	// Test sending QAsset
	err := quicksilver.InterchainstakingKeeper.MintAndSendQAsset(ctx, sender, senderAddress, &zone, amount, false, nil)
	suite.NoError(err)

	// Verify balance of receiver
	receiverBalance := quicksilver.BankKeeper.GetBalance(ctx, sender, zone.LocalDenom)
	suite.Equal(sdk.NewCoin(zone.LocalDenom, sdk.NewInt(5555)), receiverBalance)
}

func (suite *KeeperTestSuite) TestMintAndSendQAssetNon1RRMappedAccount() {
	suite.SetupTest()
	suite.setupTestZones()

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()

	zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	zone.RedemptionRate = sdk.NewDecWithPrec(110, 2)
	suite.True(found)

	senderAddress := addressutils.GenerateAddressForTestWithPrefix("cosmos")
	sender := addressutils.MustAccAddressFromBech32(senderAddress, "")
	mappedAccount := addressutils.GenerateAccAddressForTest()

	amount := sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(5000)))

	// Test sending QAsset
	err := quicksilver.InterchainstakingKeeper.MintAndSendQAsset(ctx, sender, senderAddress, &zone, amount, false, mappedAccount)
	suite.NoError(err)

	// Verify balance of receiver
	receiverBalance := quicksilver.BankKeeper.GetBalance(ctx, sender, zone.LocalDenom)
	suite.Equal(sdk.NewCoin(zone.LocalDenom, sdk.NewInt(0)), receiverBalance)

	mappedBalance := quicksilver.BankKeeper.GetBalance(ctx, mappedAccount, zone.LocalDenom)
	suite.Equal(sdk.NewCoin(zone.LocalDenom, sdk.NewInt(4545)), mappedBalance)

	remoteAddress, found := quicksilver.InterchainstakingKeeper.GetRemoteAddressMap(ctx, mappedAccount, suite.chainB.ChainID)
	suite.True(found)
	suite.Equal(senderAddress, remoteAddress.String())

	localAddress, found := quicksilver.InterchainstakingKeeper.GetLocalAddressMap(ctx, sender, suite.chainB.ChainID)
	suite.True(found)
	suite.Equal(mappedAccount, localAddress)
}

func (suite *KeeperTestSuite) TestMintAndSendQAssetNon1RTS() {
	suite.SetupTest()
	// this is required because the ibc-go test suite CreateTransferChannels defaults to a value that causes executing a message to error.
	suite.path.EndpointA.ChannelConfig.Version = "ics20-1"
	suite.path.EndpointA.Counterparty.ChannelConfig.Version = "ics20-1"
	suite.coordinator.CreateTransferChannels(suite.path)

	suite.setupTestZones()

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	ctx := suite.chainA.GetContext()

	zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	zone.RedemptionRate = sdk.NewDecWithPrec(110, 2)
	suite.True(found)

	senderAddress := addressutils.GenerateAddressForTestWithPrefix("cosmos")
	sender := addressutils.MustAccAddressFromBech32(senderAddress, "")

	amount := sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(5000)))

	// Test sending QAsset
	err := quicksilver.InterchainstakingKeeper.MintAndSendQAsset(ctx, sender, senderAddress, &zone, amount, true, nil)
	suite.NoError(err)

	// Verify balance of receiver
	receiverBalance := quicksilver.BankKeeper.GetBalance(ctx, sender, zone.LocalDenom)
	suite.Equal(sdk.NewCoin(zone.LocalDenom, sdk.NewInt(0)), receiverBalance)

	ibcEscrowAddress := transfertypes.GetEscrowAddress("transfer", "channel-0")
	ibcEscrowAccountBalance := quicksilver.BankKeeper.GetBalance(ctx, ibcEscrowAddress, zone.LocalDenom)
	suite.Equal(sdk.NewCoin(zone.LocalDenom, sdk.NewInt(4545)), ibcEscrowAccountBalance)
}

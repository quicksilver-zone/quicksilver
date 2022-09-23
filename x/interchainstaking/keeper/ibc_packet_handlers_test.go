package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v5/modules/apps/transfer/types"
	"github.com/ingenuity-build/quicksilver/app"
	"github.com/ingenuity-build/quicksilver/utils"
	icskeeper "github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	"github.com/stretchr/testify/require"
)

func TestHandleMsgTransferGood(t *testing.T) {
	app, ctx := app.GetAppWithContext(t, true)
	app.BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(sdk.NewCoin("denom", sdk.NewInt(100))))

	sender := utils.GenerateAccAddressForTest()
	senderAddr, _ := sdk.Bech32ifyAddressBytes("cosmos", sender)

	txMacc := app.AccountKeeper.GetModuleAddress(icstypes.ModuleName)
	feeMacc := app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	txMaccBalance := app.BankKeeper.GetAllBalances(ctx, txMacc)
	feeMaccBalance := app.BankKeeper.GetAllBalances(ctx, feeMacc)

	transferMsg := ibctransfertypes.MsgTransfer{
		SourcePort:    "transfer",
		SourceChannel: "channel-0",
		Token:         sdk.NewCoin("denom", sdk.NewInt(100)),
		Sender:        senderAddr,
		Receiver:      app.AccountKeeper.GetModuleAddress(icstypes.ModuleName).String(),
	}
	require.NoError(t, app.InterchainstakingKeeper.HandleMsgTransfer(ctx, &transferMsg))

	txMaccBalance2 := app.BankKeeper.GetAllBalances(ctx, txMacc)
	feeMaccBalance2 := app.BankKeeper.GetAllBalances(ctx, feeMacc)

	// assert that ics module balance is now 100denom less than before HandleMsgTransfer()
	require.Equal(t, txMaccBalance.AmountOf("denom").Sub(txMaccBalance2.AmountOf("denom")), sdk.NewInt(100))
	// assert that fee collector module balance is now 100denom more than before HandleMsgTransfer()
	require.Equal(t, feeMaccBalance2.AmountOf("denom").Sub(feeMaccBalance.AmountOf("denom")), sdk.NewInt(100))
}

func TestHandleMsgTransferBadType(t *testing.T) {
	app, ctx := app.GetAppWithContext(t, true)
	app.BankKeeper.MintCoins(ctx, ibctransfertypes.ModuleName, sdk.NewCoins(sdk.NewCoin("denom", sdk.NewInt(100))))

	transferMsg := banktypes.MsgSend{}
	require.Error(t, app.InterchainstakingKeeper.HandleMsgTransfer(ctx, &transferMsg))
}

func TestHandleMsgTransferBadRecipient(t *testing.T) {
	recipient := utils.GenerateAccAddressForTest()
	app, ctx := app.GetAppWithContext(t, true)

	sender := utils.GenerateAccAddressForTest()
	senderAddr, _ := sdk.Bech32ifyAddressBytes("cosmos", sender)

	transferMsg := ibctransfertypes.MsgTransfer{
		SourcePort:    "transfer",
		SourceChannel: "channel-0",
		Token:         sdk.NewCoin("denom", sdk.NewInt(100)),
		Sender:        senderAddr,
		Receiver:      recipient.String(),
	}
	require.Error(t, app.InterchainstakingKeeper.HandleMsgTransfer(ctx, &transferMsg))
}

// func (s *KeeperTestSuite) TestHandleSendToDelegate() {
// 	tests := []struct {
// 		name string
// 	}{
// 		{
// 			name: "valid",
// 		},
// 	}

// 	for _, test := range tests {
// 		s.Run(test.name, func() {

// 			s.SetupTest()
// 			s.SetupZones()

// 			recipient := utils.GenerateAccAddressForTest()
// 			app := s.GetQuicksilverApp(s.chainA)
// 			ctx := s.chainA.GetContext()
// 			ctx = ctx.WithContext(context.WithValue(ctx.Context(), utils.ContextKey("connectionID"), s.path.EndpointA.ConnectionID))

// 			sender := utils.GenerateAccAddressForTest()
// 			senderAddr, _ := sdk.Bech32ifyAddressBytes("cosmos", sender)

// 			sendMsg := banktypes.MsgSend{
// 				FromAddress: senderAddr,
// 				ToAddress:   recipient.String(),
// 				Amount:      sdk.NewCoins(sdk.NewCoin("denom", sdk.NewInt(100))),
// 			}
// 			s.Require().NoError(app.InterchainstakingKeeper.HandleCompleteSend(ctx, &sendMsg, ""))
// 		})
// 	}
// }

func mustGetTestBech32Address(hrp string) string {
	outAddr, err := bech32.ConvertAndEncode(hrp, utils.GenerateAccAddressForTest())
	if err != nil {
		panic(err)
	}
	return outAddr
}

func (s *KeeperTestSuite) TestHandleQueuedUnbondings() {
	tests := []struct {
		name    string
		records func(chainID string, hrp string) []icstypes.WithdrawalRecord
	}{
		{
			name: "valid",
			records: func(chainID string, hrp string) []icstypes.WithdrawalRecord {
				return []icstypes.WithdrawalRecord{
					{
						ChainId:   chainID,
						Delegator: utils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{utils.GenerateValAddressForTest().String(), 1000000},
							{utils.GenerateValAddressForTest().String(), 1000000},
							{utils.GenerateValAddressForTest().String(), 1000000},
							{utils.GenerateValAddressForTest().String(), 1000000},
						},
						Recipient:  mustGetTestBech32Address(hrp),
						Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(4000000))),
						BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(4000000)),
						Txhash:     "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
						Status:     icskeeper.WithdrawStatusQueued,
					},
				}
			},
		},
		{
			name: "valid - two",
			records: func(chainID string, hrp string) []icstypes.WithdrawalRecord {
				return []icstypes.WithdrawalRecord{
					{
						ChainId:   chainID,
						Delegator: utils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{utils.GenerateValAddressForTest().String(), 1000000},
							{utils.GenerateValAddressForTest().String(), 1000000},
							{utils.GenerateValAddressForTest().String(), 1000000},
							{utils.GenerateValAddressForTest().String(), 1000000},
						},
						Recipient:  mustGetTestBech32Address(hrp),
						Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(4000000))),
						BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(4000000)),
						Txhash:     "7C8B95EEE82CB63771E02EBEB05E6A80076D70B2E0A1C457F1FD1A0EF2EA961D",
						Status:     icskeeper.WithdrawStatusQueued,
					},
					{
						ChainId:   chainID,
						Delegator: utils.GenerateAccAddressForTest().String(),
						Distribution: []*icstypes.Distribution{
							{utils.GenerateValAddressForTest().String(), 5000000},
							{utils.GenerateValAddressForTest().String(), 1250000},
							{utils.GenerateValAddressForTest().String(), 5000000},
							{utils.GenerateValAddressForTest().String(), 1250000},
						},
						Recipient:  mustGetTestBech32Address(hrp),
						Amount:     sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(15000000))),
						BurnAmount: sdk.NewCoin("uqatom", sdk.NewInt(15000000)),
						Txhash:     "d786f7d4c94247625c2882e921a790790eb77a00d0534d5c3154d0a9c5ab68f5",
						Status:     icskeeper.WithdrawStatusQueued,
					},
				}
			},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {

			s.SetupTest()
			s.SetupZones()

			app := s.GetQuicksilverApp(s.chainA)
			ctx := s.chainA.GetContext()

			zone, found := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
			if !found {
				s.Fail("unable to retrieve zone for test")
			}

			records := test.records(s.chainB.ChainID, zone.AccountPrefix)

			// set up zones
			for _, record := range records {
				app.InterchainstakingKeeper.SetWithdrawalRecord(ctx, record)
			}

			// trigger handler
			err := app.InterchainstakingKeeper.HandleQueuedUnbondings(ctx, &zone, 1)
			s.Require().NoError(err)

			for _, record := range records {
				// check record with old status cannot be found
				_, found := app.InterchainstakingKeeper.GetWithdrawalRecord(ctx, zone.ChainId, record.Txhash, icskeeper.WithdrawStatusQueued)
				s.Require().False(found)
				// check record with new status can be found
				_, found = app.InterchainstakingKeeper.GetWithdrawalRecord(ctx, zone.ChainId, record.Txhash, icskeeper.WithdrawStatusUnbond)
				s.Require().True(found)
			}
		})
	}
}

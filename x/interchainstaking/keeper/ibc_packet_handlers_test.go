package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	"github.com/ingenuity-build/quicksilver/app"
	"github.com/ingenuity-build/quicksilver/utils"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	"github.com/stretchr/testify/require"
)

func TestHandleMsgTransferGood(t *testing.T) {
	app, ctx := app.GetAppWithContext(true)
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
	require.Equal(t, txMaccBalance.Sub(txMaccBalance2), sdk.NewCoins(sdk.NewCoin("denom", sdk.NewInt(100))))
	// assert that fee collector module balance is now 100denom more than before HandleMsgTransfer()
	require.Equal(t, feeMaccBalance2.Sub(feeMaccBalance), sdk.NewCoins(sdk.NewCoin("denom", sdk.NewInt(100))))
}

func TestHandleMsgTransferBadType(t *testing.T) {
	app, ctx := app.GetAppWithContext(true)
	app.BankKeeper.MintCoins(ctx, ibctransfertypes.ModuleName, sdk.NewCoins(sdk.NewCoin("denom", sdk.NewInt(100))))

	transferMsg := banktypes.MsgSend{}
	require.Error(t, app.InterchainstakingKeeper.HandleMsgTransfer(ctx, &transferMsg))
}

func TestHandleMsgTransferBadRecipient(t *testing.T) {
	recipient := utils.GenerateAccAddressForTest()
	app, ctx := app.GetAppWithContext(true)

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

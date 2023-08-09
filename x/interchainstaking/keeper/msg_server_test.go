package keeper_test

import (
	"fmt"
	"time"

	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ingenuity-build/quicksilver/utils/addressutils"
	icskeeper "github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func (suite *KeeperTestSuite) TestRequestRedemption() {
	var (
		msg         icstypes.MsgRequestRedemption
		zoneID      string
		testAccount sdk.AccAddress
		err         error
		denom       string
	)

	tests := []struct {
		name         string
		init         func()
		malleate     func()
		expectErr    string
		expectErrLsm string
	}{
		{
			"valid - full claim",
			func() {
				testAccount, err = addressutils.AccAddressFromBech32(testAddress, "")
				suite.NoError(err)
				zoneID = testzoneID
				denom = "uqatom"
			},
			func() {
				addr, err := addressutils.EncodeAddressToBech32("cosmos", addressutils.GenerateAccAddressForTest())
				suite.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin(denom, sdk.NewInt(10000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}
			},
			"",
			"",
		},
		{
			"valid - full claim for subzone",
			func() {
				testAccount, err = addressutils.AccAddressFromBech32(subzoneAddress, "")
				suite.NoError(err)
				zoneID = subzoneID
				denom = "usqatom"
			},
			func() {
				addr, err := addressutils.EncodeAddressToBech32("cosmos", addressutils.GenerateAccAddressForTest())
				suite.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin(denom, sdk.NewInt(10000000)),
					DestinationAddress: addr,
					FromAddress:        subzoneAddress,
				}
			},
			"",
			"",
		},
		{
			"invalid - incorrect authority for subzone",
			func() {
				testAccount, err = addressutils.AccAddressFromBech32(testAddress, "")
				suite.NoError(err)
				zoneID = subzoneID
				denom = "usqatom"
			},
			func() {
				addr, err := addressutils.EncodeAddressToBech32("cosmos", addressutils.GenerateAccAddressForTest())
				suite.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin(denom, sdk.NewInt(10000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}
			},
			"invalid authority for subzone",
			"invalid authority for subzone",
		},
		{
			"valid - full claim (discounted)",
			func() {
				testAccount, err = addressutils.AccAddressFromBech32(testAddress, "")
				suite.NoError(err)
				zoneID = testzoneID
				denom = "uqatom"
			},
			func() {
				addr, err := addressutils.EncodeAddressToBech32("cosmos", addressutils.GenerateAccAddressForTest())
				suite.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin(denom, sdk.NewInt(10000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}

				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), testzoneID)
				suite.True(found)
				zone.RedemptionRate = sdk.MustNewDecFromStr("0.95")
				suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetZone(suite.chainA.GetContext(), &zone)
			},
			"",
			"",
		},
		{
			"valid - full claim (interest)",
			func() {
				testAccount, err = addressutils.AccAddressFromBech32(testAddress, "")
				suite.NoError(err)
				zoneID = testzoneID
				denom = "uqatom"
			},
			func() {
				addr, err := addressutils.EncodeAddressToBech32("cosmos", addressutils.GenerateAccAddressForTest())
				suite.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin(denom, sdk.NewInt(10000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}

				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), testzoneID)
				suite.True(found)
				zone.LastRedemptionRate = sdk.MustNewDecFromStr("1.05")
				zone.RedemptionRate = sdk.MustNewDecFromStr("1.1")
				suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetZone(suite.chainA.GetContext(), &zone)
			},
			"",
			"",
		},
		{
			"valid - full claim (interest)",
			func() {
				testAccount, err = addressutils.AccAddressFromBech32(testAddress, "")
				suite.NoError(err)
				zoneID = testzoneID
				denom = "uqatom"
			},
			func() {
				addr, err := addressutils.EncodeAddressToBech32("cosmos", addressutils.GenerateAccAddressForTest())
				suite.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin(denom, sdk.NewInt(10000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}

				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), testzoneID)
				suite.True(found)
				zone.LastRedemptionRate = sdk.MustNewDecFromStr("1.1")
				zone.RedemptionRate = sdk.MustNewDecFromStr("1.05")
				suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetZone(suite.chainA.GetContext(), &zone)
			},
			"",
			"",
		},
		{
			"valid - partial claim",
			func() {
				testAccount, err = addressutils.AccAddressFromBech32(testAddress, "")
				suite.NoError(err)
				zoneID = testzoneID
				denom = "uqatom"
			},
			func() {
				addr, err := addressutils.EncodeAddressToBech32("cosmos", addressutils.GenerateAccAddressForTest())
				suite.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin(denom, sdk.NewInt(5000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}
			},
			"",
			"",
		},
		{
			"valid - partial claim (discounted)",
			func() {
				testAccount, err = addressutils.AccAddressFromBech32(testAddress, "")
				suite.NoError(err)
				zoneID = testzoneID
				denom = "uqatom"
			},
			func() {
				addr, err := addressutils.EncodeAddressToBech32("cosmos", addressutils.GenerateAccAddressForTest())
				suite.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin(denom, sdk.NewInt(5000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}

				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), testzoneID)
				suite.True(found)
				zone.RedemptionRate = sdk.MustNewDecFromStr("0.99999")
				suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetZone(suite.chainA.GetContext(), &zone)
			},
			"",
			"",
		},
		{
			"valid - partial claim (interest)",
			func() {
				testAccount, err = addressutils.AccAddressFromBech32(testAddress, "")
				suite.NoError(err)
				zoneID = testzoneID
				denom = "uqatom"
			},
			func() {
				addr, err := addressutils.EncodeAddressToBech32("cosmos", addressutils.GenerateAccAddressForTest())
				suite.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin(denom, sdk.NewInt(5000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}

				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), testzoneID)
				suite.True(found)
				zone.LastRedemptionRate = sdk.MustNewDecFromStr("1.049999")
				zone.RedemptionRate = sdk.MustNewDecFromStr("1.099999")
				suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetZone(suite.chainA.GetContext(), &zone)
			},
			"",
			"",
		},
		{
			"invalid - unbonding not enabled for zone",
			func() {
				testAccount, err = addressutils.AccAddressFromBech32(testAddress, "")
				suite.NoError(err)
				zoneID = testzoneID
				denom = "uqatom"
			},
			func() {
				addr, err := addressutils.EncodeAddressToBech32("cosmos", addressutils.GenerateAccAddressForTest())
				suite.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin(denom, sdk.NewInt(5000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}

				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), testzoneID)
				suite.True(found)
				zone.UnbondingEnabled = false
				suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetZone(suite.chainA.GetContext(), &zone)
			},
			"unbonding currently disabled for zone testchain2",
			"unbonding currently disabled for zone testchain2",
		},
		{
			"invalid - wrong denom",
			func() {
				testAccount, err = addressutils.AccAddressFromBech32(testAddress, "")
				suite.NoError(err)
				zoneID = testzoneID
				denom = "uqatom"
			},
			func() {
				addr, err := addressutils.EncodeAddressToBech32("cosmos", addressutils.GenerateAccAddressForTest())
				suite.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uatom", sdk.NewInt(10000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}
			},
			"unable to find matching zone for denom uatom",
			"unable to find matching zone for denom uatom",
		},
		{
			"invalid - insufficient funds",
			func() {
				testAccount, err = addressutils.AccAddressFromBech32(testAddress, "")
				suite.NoError(err)
				zoneID = testzoneID
				denom = "uqatom"
			},
			func() {
				addr, err := addressutils.EncodeAddressToBech32("cosmos", addressutils.GenerateAccAddressForTest())
				suite.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin(denom, sdk.NewInt(1000000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}
			},
			"account has insufficient balance of qasset to burn",
			"account has insufficient balance of qasset to burn",
		},
		{
			"invalid - bad prefix",
			func() {
				testAccount, err = addressutils.AccAddressFromBech32(testAddress, "")
				suite.NoError(err)
				zoneID = testzoneID
				denom = "uqatom"
			},
			func() {
				addr, err := addressutils.EncodeAddressToBech32("bob", addressutils.GenerateAccAddressForTest())
				suite.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin(denom, sdk.OneInt()),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}
			},
			"destination address bob",
			"destination address bob",
		},
		{
			"invalid - bad from address",
			func() {
				testAccount, err = addressutils.AccAddressFromBech32(testAddress, "")
				suite.NoError(err)
				zoneID = testzoneID
				denom = "uqatom"
			},
			func() {
				addr, err := addressutils.EncodeAddressToBech32("cosmos", addressutils.GenerateAccAddressForTest())
				suite.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin(denom, sdk.OneInt()),
					DestinationAddress: addr,
					FromAddress:        addr,
				}
			},
			"account has insufficient balance of qasset to burn",
			"account has insufficient balance of qasset to burn",
		},
		{
			"invalid - too many locked tokens",
			func() {
				testAccount, err = addressutils.AccAddressFromBech32(testAddress, "")
				suite.NoError(err)
				zoneID = testzoneID
				denom = "uqatom"
			},
			func() {
				addr, err := addressutils.EncodeAddressToBech32("cosmos", addressutils.GenerateAccAddressForTest())
				suite.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin(denom, sdk.NewInt(10000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}

				ctx := suite.chainA.GetContext()
				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(ctx, testzoneID)
				suite.True(found)
				zoneVals := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetValidatorAddresses(ctx, &zone)
				suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetRedelegationRecord(ctx, icstypes.RedelegationRecord{
					ChainId:        testzoneID,
					EpochNumber:    1,
					Source:         zoneVals[0],
					Destination:    zoneVals[1],
					Amount:         3000000,
					CompletionTime: suite.chainA.GetContext().BlockTime().Add(time.Hour),
				})
			},
			"",
			"unable to satisfy unbond request; delegations may be locked",
		},
		{
			"invalid - unbonding is disabled",
			func() {
				testAccount, err = addressutils.AccAddressFromBech32(testAddress, "")
				suite.NoError(err)
				zoneID = testzoneID
				denom = "uqatom"
			},
			func() {
				ctx := suite.chainA.GetContext()

				addr, err := addressutils.EncodeAddressToBech32("cosmos", addressutils.GenerateAccAddressForTest())
				suite.Require().NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin(denom, sdk.NewInt(10000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}
				params := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetParams(ctx)
				params.UnbondingEnabled = false
				suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetParams(ctx, params)
			},
			"unbonding is currently disabled",
			"unbonding is currently disabled",
		},
	}

	for _, tt := range tests {
		tt := tt

		// run tests with LSM disabled.
		suite.Run(tt.name, func() {
			suite.SetupTest()
			suite.setupTestZones()
			tt.init()

			quicksilver := suite.GetQuicksilverApp(suite.chainA)
			ctx := suite.chainA.GetContext()

			params := quicksilver.InterchainstakingKeeper.GetParams(ctx)
			params.UnbondingEnabled = true
			quicksilver.InterchainstakingKeeper.SetParams(ctx, params)

			err := quicksilver.BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(sdk.NewCoin(denom, math.NewInt(10000000))))
			suite.NoError(err)
			err = quicksilver.BankKeeper.SendCoinsFromModuleToAccount(ctx, icstypes.ModuleName, testAccount, sdk.NewCoins(sdk.NewCoin(denom, math.NewInt(10000000))))
			suite.NoError(err)

			quicksilver.InterchainstakingKeeper.IterateZones(ctx, func(index int64, zone *icstypes.Zone) (stop bool) {
				_ = zone
				return false
			})

			// disable LSM
			zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, zoneID)
			suite.True(found)
			zone.LiquidityModule = false
			zone.UnbondingEnabled = true
			quicksilver.InterchainstakingKeeper.SetZone(ctx, &zone)

			tt.malleate()

			msgSrv := icskeeper.NewMsgServerImpl(*quicksilver.InterchainstakingKeeper)
			res, err := msgSrv.RequestRedemption(sdk.WrapSDKContext(ctx), &msg)

			if tt.expectErr != "" {
				suite.ErrorContains(err, tt.expectErr)
				suite.Nil(res)
				suite.T().Logf("Error: %v", err)
			} else {
				suite.NoError(err)
				suite.NotNil(res)
			}
		})

		// run tests with LSM enabled.
		tt.name += "_LSM_enabled"
		suite.Run(tt.name, func() {
			suite.SetupTest()
			suite.setupTestZones()

			quicksilver := suite.GetQuicksilverApp(suite.chainA)
			ctx := suite.chainA.GetContext()

			tt.init()

			params := quicksilver.InterchainstakingKeeper.GetParams(ctx)
			params.UnbondingEnabled = true
			quicksilver.InterchainstakingKeeper.SetParams(ctx, params)

			err := quicksilver.BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqatom", math.NewInt(10000000))))
			suite.NoError(err)
			err = quicksilver.BankKeeper.SendCoinsFromModuleToAccount(ctx, icstypes.ModuleName, testAccount, sdk.NewCoins(sdk.NewCoin("uqatom", math.NewInt(10000000))))
			suite.NoError(err)

			// enable LSM
			zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, zoneID)
			suite.True(found)
			zone.LiquidityModule = true
			zone.UnbondingEnabled = true
			quicksilver.InterchainstakingKeeper.SetZone(ctx, &zone)

			validators := quicksilver.InterchainstakingKeeper.GetValidatorAddresses(ctx, &zone)
			for _, delegation := range func(zone icstypes.Zone) []icstypes.Delegation {
				out := make([]icstypes.Delegation, 0)
				for _, valoper := range validators {
					out = append(out, icstypes.NewDelegation(zone.DelegationAddress.Address, valoper, sdk.NewCoin(zone.BaseDenom, sdk.NewInt(3000000))))
				}
				return out
			}(zone) {
				quicksilver.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegation)
			}

			tt.malleate()

			msgSrv := icskeeper.NewMsgServerImpl(*quicksilver.InterchainstakingKeeper)
			res, err := msgSrv.RequestRedemption(sdk.WrapSDKContext(ctx), &msg)

			if tt.expectErrLsm != "" {
				suite.Errorf(err, tt.expectErrLsm)
				suite.Nil(res)
				suite.T().Logf("Error: %v", err)
			} else {
				suite.NoError(err)
				suite.NotNil(res)
			}
		})

	}
}

func (suite *KeeperTestSuite) TestSignalIntent() {
	tests := []struct {
		name             string
		malleate         func(suite *KeeperTestSuite) *icstypes.MsgSignalIntent
		fromAddress      string
		zoneID           string
		expected         []sdk.Dec
		failsValidations bool
		expectErr        bool
	}{
		{
			"invalid - weight sum < 1",
			func(suite *KeeperTestSuite) *icstypes.MsgSignalIntent {
				val1, err := sdk.ValAddressFromHex(suite.chainB.Vals.Validators[0].Address.String())
				suite.NoError(err)

				return &icstypes.MsgSignalIntent{
					ChainId:     testzoneID,
					Intents:     fmt.Sprintf("0.3%s", val1.String()),
					FromAddress: testAddress,
				}
			},
			testAddress,
			testzoneID,
			[]sdk.Dec{},
			true,
			false,
		},
		{
			"invalid - weight sum > 1",
			func(suite *KeeperTestSuite) *icstypes.MsgSignalIntent {
				val1, err := sdk.ValAddressFromHex(suite.chainB.Vals.Validators[0].Address.String())
				suite.NoError(err)

				return &icstypes.MsgSignalIntent{
					ChainId:     testzoneID,
					Intents:     fmt.Sprintf("3.0%s", val1.String()),
					FromAddress: testAddress,
				}
			},
			testAddress,
			testzoneID,
			[]sdk.Dec{},
			true,
			false,
		},
		{
			"invalid - invalid authority for subzone",
			func(suite *KeeperTestSuite) *icstypes.MsgSignalIntent {
				val1, err := sdk.ValAddressFromHex(suite.chainB.Vals.Validators[0].Address.String())
				suite.NoError(err)

				return &icstypes.MsgSignalIntent{
					ChainId:     subzoneID,
					Intents:     fmt.Sprintf("1.0%s", val1.String()),
					FromAddress: testAddress,
				}
			},
			subzoneID,
			testzoneID,
			[]sdk.Dec{},
			false,
			true,
		},
		{
			"invalid - chain id",
			func(suite *KeeperTestSuite) *icstypes.MsgSignalIntent {
				val1, err := sdk.ValAddressFromHex(suite.chainB.Vals.Validators[0].Address.String())
				suite.NoError(err)

				return &icstypes.MsgSignalIntent{
					ChainId:     suite.chainA.ChainID,
					Intents:     fmt.Sprintf("1.0%s", val1.String()),
					FromAddress: testAddress,
				}
			},
			testAddress,
			testzoneID,
			[]sdk.Dec{},
			false,
			true,
		},
		{
			"valid - single weight",
			func(suite *KeeperTestSuite) *icstypes.MsgSignalIntent {
				val1, err := sdk.ValAddressFromHex(suite.chainB.Vals.Validators[0].Address.String())
				suite.NoError(err)

				return &icstypes.MsgSignalIntent{
					ChainId:     testzoneID,
					Intents:     fmt.Sprintf("1.0%s", val1.String()),
					FromAddress: testAddress,
				}
			},
			testAddress,
			testzoneID,
			[]sdk.Dec{sdk.NewDecWithPrec(1, 0)},
			false,
			false,
		},
		{
			"valid - multi weight",
			func(suite *KeeperTestSuite) *icstypes.MsgSignalIntent {
				val1, err := sdk.ValAddressFromHex(suite.chainB.Vals.Validators[0].Address.String())
				suite.NoError(err)
				val2, err := sdk.ValAddressFromHex(suite.chainB.Vals.Validators[1].Address.String())
				suite.NoError(err)
				val3, err := sdk.ValAddressFromHex(suite.chainB.Vals.Validators[2].Address.String())
				suite.NoError(err)

				return &icstypes.MsgSignalIntent{
					ChainId:     testzoneID,
					Intents:     fmt.Sprintf("0.5%s,0.2%s,0.3%s", val1.String(), val2.String(), val3.String()),
					FromAddress: testAddress,
				}
			},
			testAddress,
			testzoneID,
			[]sdk.Dec{
				sdk.NewDecWithPrec(5, 1),
				sdk.NewDecWithPrec(2, 1),
				sdk.NewDecWithPrec(3, 1),
			},
			false,
			false,
		},
		{
			"valid - single weight  for subzone",
			func(suite *KeeperTestSuite) *icstypes.MsgSignalIntent {
				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), subzoneID)
				suite.True(found)

				val1, err := sdk.ValAddressFromHex(suite.chainB.Vals.Validators[0].Address.String())
				suite.NoError(err)

				return &icstypes.MsgSignalIntent{
					ChainId:     zone.ZoneID(),
					Intents:     fmt.Sprintf("1.0%s", val1.String()),
					FromAddress: subzoneAddress,
				}
			},
			subzoneAddress,
			subzoneID,
			[]sdk.Dec{sdk.NewDecWithPrec(1, 0)},
			false,
			false,
		},
		{
			"valid - multi weight for subzone",
			func(suite *KeeperTestSuite) *icstypes.MsgSignalIntent {
				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), subzoneID)
				suite.True(found)

				val1, err := sdk.ValAddressFromHex(suite.chainB.Vals.Validators[0].Address.String())
				suite.NoError(err)
				val2, err := sdk.ValAddressFromHex(suite.chainB.Vals.Validators[1].Address.String())
				suite.NoError(err)
				val3, err := sdk.ValAddressFromHex(suite.chainB.Vals.Validators[2].Address.String())
				suite.NoError(err)

				return &icstypes.MsgSignalIntent{
					ChainId:     zone.ZoneID(),
					Intents:     fmt.Sprintf("0.5%s,0.2%s,0.3%s", val1.String(), val2.String(), val3.String()),
					FromAddress: subzoneAddress,
				}
			},
			subzoneAddress,
			subzoneID,
			[]sdk.Dec{
				sdk.NewDecWithPrec(5, 1),
				sdk.NewDecWithPrec(2, 1),
				sdk.NewDecWithPrec(3, 1),
			},
			false,
			false,
		},
	}

	for _, tt := range tests {
		tt := tt

		suite.Run(tt.name, func() {
			suite.SetupTest()
			suite.setupTestZones()

			msg := tt.malleate(suite)
			// validateBasic not explicitly tested here - but we don't call it inside msgSrv.SignalIntent
			// so call here to make sure out tests are sane.
			err := msg.ValidateBasic()
			if tt.failsValidations {
				suite.Error(err)
				return
			}
			suite.NoError(err)

			msgSrv := icskeeper.NewMsgServerImpl(*suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper)
			res, err := msgSrv.SignalIntent(sdk.WrapSDKContext(suite.chainA.GetContext()), msg)
			if tt.expectErr {
				suite.Error(err)
				suite.Nil(res)
				return
			}

			suite.NoError(err)
			suite.NotNil(res)

			quicksilver := suite.GetQuicksilverApp(suite.chainA)
			icsKeeper := quicksilver.InterchainstakingKeeper
			zone, found := icsKeeper.GetZone(suite.chainA.GetContext(), tt.zoneID)
			suite.True(found)

			intent, found := icsKeeper.GetDelegatorIntent(suite.chainA.GetContext(), &zone, tt.fromAddress, false)
			suite.True(found)
			intents := intent.GetIntents()

			for idx, weight := range tt.expected {
				val, err := sdk.ValAddressFromHex(suite.chainB.Vals.Validators[idx].Address.String())
				suite.NoError(err)

				valIntent, found := intents.GetForValoper(val.String())
				suite.True(found)

				suite.Equal(weight, valIntent.Weight)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestRegisterZone() {
	var msg *icstypes.MsgRegisterZone

	testAccount, err := addressutils.AccAddressFromBech32(testAddress, "")
	suite.NoError(err)

	tests := []struct {
		name      string
		malleate  func()
		expectErr string
	}{
		{
			"invalid: duplicate zone",
			func() {
				msg = &icstypes.MsgRegisterZone{
					Authority:        suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetGovAuthority(),
					ConnectionID:     suite.path.EndpointA.ConnectionID,
					LocalDenom:       "uqatom",
					BaseDenom:        "uatom",
					AccountPrefix:    "cosmos",
					ReturnToSender:   false,
					UnbondingEnabled: false,
					LiquidityModule:  true,
					DepositsEnabled:  true,
					Decimals:         6,
					Is_118:           true,
					SubzoneInfo:      nil,
				}
			},
			"invalid chain id",
		},
		{
			"invalid: unknown connectionID",
			func() {
				msg = &icstypes.MsgRegisterZone{
					Authority:        suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetGovAuthority(),
					ConnectionID:     "invalid",
					LocalDenom:       "uqatom",
					BaseDenom:        "uatom",
					AccountPrefix:    "cosmos",
					ReturnToSender:   false,
					UnbondingEnabled: false,
					LiquidityModule:  true,
					DepositsEnabled:  true,
					Decimals:         6,
					Is_118:           true,
					SubzoneInfo:      nil,
				}
			},
			"unable to obtain chain id",
		},
		{
			"invalid: incorrect authority",
			func() {
				msg = &icstypes.MsgRegisterZone{
					Authority:        "invalid",
					ConnectionID:     suite.path.EndpointA.ConnectionID,
					LocalDenom:       "uqatom",
					BaseDenom:        "uatom",
					AccountPrefix:    "cosmos",
					ReturnToSender:   false,
					UnbondingEnabled: false,
					LiquidityModule:  true,
					DepositsEnabled:  true,
					Decimals:         6,
					Is_118:           true,
					SubzoneInfo:      nil,
				}
			},
			"invalid authority",
		},
		{
			"invalid: invalid subzone info: ID mismatch",
			func() {
				msg = &icstypes.MsgRegisterZone{
					Authority:        suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetGovAuthority(),
					ConnectionID:     suite.path.EndpointA.ConnectionID,
					LocalDenom:       "uqatom",
					BaseDenom:        "uatom",
					AccountPrefix:    "cosmos",
					ReturnToSender:   false,
					UnbondingEnabled: false,
					LiquidityModule:  true,
					DepositsEnabled:  true,
					Decimals:         6,
					Is_118:           true,
					SubzoneInfo: &icstypes.SubzoneInfo{
						Authority:   "test",
						BaseChainID: "invalid",
						ChainID:     "test-1",
					},
				}
			},
			"incorrect ID",
		},
		{
			"invalid: invalid subzone info: subzone ID taken",
			func() {
				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), testzoneID)
				suite.True(found)

				msg = &icstypes.MsgRegisterZone{
					Authority:        suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetGovAuthority(),
					ConnectionID:     suite.path.EndpointA.ConnectionID,
					LocalDenom:       "uqatom",
					BaseDenom:        "uatom",
					AccountPrefix:    "cosmos",
					ReturnToSender:   false,
					UnbondingEnabled: false,
					LiquidityModule:  true,
					DepositsEnabled:  true,
					Decimals:         6,
					Is_118:           true,
					SubzoneInfo: &icstypes.SubzoneInfo{
						Authority:   subzoneAddress,
						BaseChainID: zone.BaseChainID(),
						ChainID:     zone.BaseChainID(),
					},
				}
			},
			"subzone ID already exists",
		},
		{
			"invalid: invalid subzone info: invalid subzone authority info",
			func() {
				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), testzoneID)
				suite.True(found)

				msg = &icstypes.MsgRegisterZone{
					Authority:        suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetGovAuthority(),
					ConnectionID:     suite.path.EndpointA.ConnectionID,
					LocalDenom:       "uqatom",
					BaseDenom:        "uatom",
					AccountPrefix:    "cosmos",
					ReturnToSender:   false,
					UnbondingEnabled: false,
					LiquidityModule:  true,
					DepositsEnabled:  true,
					Decimals:         6,
					Is_118:           true,
					SubzoneInfo: &icstypes.SubzoneInfo{
						Authority:   "",
						BaseChainID: zone.BaseChainID(),
						ChainID:     "test-1234",
					},
				}
			},
			"all subzone info must be populated",
		},
	}

	for _, tt := range tests {
		tt := tt
		suite.Run(tt.name, func() {
			suite.SetupTest()
			suite.setupTestZones()

			ctx := suite.chainA.GetContext()

			err := suite.GetQuicksilverApp(suite.chainA).BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqatom", math.NewInt(10000000))))
			suite.NoError(err)
			err = suite.GetQuicksilverApp(suite.chainA).BankKeeper.SendCoinsFromModuleToAccount(ctx, icstypes.ModuleName, testAccount, sdk.NewCoins(sdk.NewCoin("uqatom", math.NewInt(10000000))))
			suite.NoError(err)

			tt.malleate()

			msgSrv := icskeeper.NewMsgServerImpl(*suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper)
			res, err := msgSrv.RegisterZone(sdk.WrapSDKContext(ctx), msg)

			if tt.expectErr != "" {
				suite.ErrorContains(err, tt.expectErr)
				suite.T().Logf("Error: %v", err)
			} else {
				suite.NoError(err)
				suite.NotNil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestUpdateZone() {
	var msg *icstypes.MsgUpdateZone

	tests := []struct {
		name      string
		malleate  func()
		expectErr string
	}{
		{
			"invalid: incorrect authority",
			func() {
				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), testzoneID)
				suite.True(found)

				msg = &icstypes.MsgUpdateZone{
					Authority: "invalid",
					ZoneID:    zone.BaseChainID(),
					Changes:   nil,
				}
			},
			"invalid authority",
		},
		{
			"invalid: zone does not exist",
			func() {
				msg = &icstypes.MsgUpdateZone{
					Authority: suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetGovAuthority(),
					ZoneID:    "invalid",
					Changes:   nil,
				}
			},
			"unable to get registered zone for zone id",
		},
		{
			"valid: no changes",
			func() {
				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), testzoneID)
				suite.True(found)

				msg = &icstypes.MsgUpdateZone{
					Authority: suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetGovAuthority(),
					ZoneID:    zone.BaseChainID(),
					Changes:   nil,
				}
			},
			"",
		},
		{
			"valid: update base denom",
			func() {
				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), testzoneID)
				suite.True(found)

				msg = &icstypes.MsgUpdateZone{
					Authority: suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetGovAuthority(),
					ZoneID:    zone.BaseChainID(),
					Changes: []*icstypes.UpdateZoneValue{
						{
							Key:   icstypes.UpdateZoneKeyBaseDenom,
							Value: "valid",
						},
					},
				}
			},
			"",
		},
		{
			"invalid: update base denom: invalid denom",
			func() {
				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), testzoneID)
				suite.True(found)

				msg = &icstypes.MsgUpdateZone{
					Authority: suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetGovAuthority(),
					ZoneID:    zone.BaseChainID(),
					Changes: []*icstypes.UpdateZoneValue{
						{
							Key:   icstypes.UpdateZoneKeyBaseDenom,
							Value: "",
						},
					},
				}
			},
			"invalid denom",
		},
		{
			"valid: update local denom",
			func() {
				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), testzoneID)
				suite.True(found)

				msg = &icstypes.MsgUpdateZone{
					Authority: suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetGovAuthority(),
					ZoneID:    zone.BaseChainID(),
					Changes: []*icstypes.UpdateZoneValue{
						{
							Key:   icstypes.UpdateZoneKeyLocalDenom,
							Value: "valid",
						},
					},
				}
			},
			"",
		},
		{
			"invalid: update local denom: invalid denom",
			func() {
				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), testzoneID)
				suite.True(found)

				msg = &icstypes.MsgUpdateZone{
					Authority: suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetGovAuthority(),
					ZoneID:    zone.BaseChainID(),
					Changes: []*icstypes.UpdateZoneValue{
						{
							Key:   icstypes.UpdateZoneKeyLocalDenom,
							Value: "",
						},
					},
				}
			},
			"invalid denom",
		},
		{
			"valid: update liquidity module",
			func() {
				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), testzoneID)
				suite.True(found)

				msg = &icstypes.MsgUpdateZone{
					Authority: suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetGovAuthority(),
					ZoneID:    zone.BaseChainID(),
					Changes: []*icstypes.UpdateZoneValue{
						{
							Key:   icstypes.UpdateZoneKeyLiquidityModule,
							Value: "true",
						},
					},
				}
			},
			"",
		},
		{
			"invalid: update liquidity module: invalid syntax",
			func() {
				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), testzoneID)
				suite.True(found)

				msg = &icstypes.MsgUpdateZone{
					Authority: suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetGovAuthority(),
					ZoneID:    zone.BaseChainID(),
					Changes: []*icstypes.UpdateZoneValue{
						{
							Key:   icstypes.UpdateZoneKeyLiquidityModule,
							Value: "",
						},
					},
				}
			},
			"invalid syntax",
		},
		{
			"valid: update unbonding enabled",
			func() {
				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), testzoneID)
				suite.True(found)

				msg = &icstypes.MsgUpdateZone{
					Authority: suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetGovAuthority(),
					ZoneID:    zone.BaseChainID(),
					Changes: []*icstypes.UpdateZoneValue{
						{
							Key:   icstypes.UpdateZoneKeyUnbondingEnabled,
							Value: "true",
						},
					},
				}
			},
			"",
		},
		{
			"invalid: update unbonding enabled: invalid syntax",
			func() {
				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), testzoneID)
				suite.True(found)

				msg = &icstypes.MsgUpdateZone{
					Authority: suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetGovAuthority(),
					ZoneID:    zone.BaseChainID(),
					Changes: []*icstypes.UpdateZoneValue{
						{
							Key:   icstypes.UpdateZoneKeyUnbondingEnabled,
							Value: "",
						},
					},
				}
			},
			"invalid syntax",
		},
		{
			"valid: update deposits enabled",
			func() {
				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), testzoneID)
				suite.True(found)

				msg = &icstypes.MsgUpdateZone{
					Authority: suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetGovAuthority(),
					ZoneID:    zone.BaseChainID(),
					Changes: []*icstypes.UpdateZoneValue{
						{
							Key:   icstypes.UpdateZoneKeyDepositsEnabled,
							Value: "true",
						},
					},
				}
			},
			"",
		},
		{
			"invalid: update deposits enabled: invalid syntax",
			func() {
				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), testzoneID)
				suite.True(found)

				msg = &icstypes.MsgUpdateZone{
					Authority: suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetGovAuthority(),
					ZoneID:    zone.BaseChainID(),
					Changes: []*icstypes.UpdateZoneValue{
						{
							Key:   icstypes.UpdateZoneKeyDepositsEnabled,
							Value: "",
						},
					},
				}
			},
			"invalid syntax",
		},
		{
			"valid: update return to sender",
			func() {
				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), testzoneID)
				suite.True(found)

				msg = &icstypes.MsgUpdateZone{
					Authority: suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetGovAuthority(),
					ZoneID:    zone.BaseChainID(),
					Changes: []*icstypes.UpdateZoneValue{
						{
							Key:   icstypes.UpdateZoneKeyReturnToSender,
							Value: "true",
						},
					},
				}
			},
			"",
		},
		{
			"invalid: update return to sender: invalid syntax",
			func() {
				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), testzoneID)
				suite.True(found)

				msg = &icstypes.MsgUpdateZone{
					Authority: suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetGovAuthority(),
					ZoneID:    zone.BaseChainID(),
					Changes: []*icstypes.UpdateZoneValue{
						{
							Key:   icstypes.UpdateZoneKeyReturnToSender,
							Value: "",
						},
					},
				}
			},
			"invalid syntax",
		},
		{
			"valid: update messages per tx",
			func() {
				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), testzoneID)
				suite.True(found)

				msg = &icstypes.MsgUpdateZone{
					Authority: suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetGovAuthority(),
					ZoneID:    zone.BaseChainID(),
					Changes: []*icstypes.UpdateZoneValue{
						{
							Key:   icstypes.UpdateZoneKeyMessagesPerTx,
							Value: "10",
						},
					},
				}
			},
			"",
		},
		{
			"invalid: update messages per tx: invalid syntax",
			func() {
				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), testzoneID)
				suite.True(found)

				msg = &icstypes.MsgUpdateZone{
					Authority: suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetGovAuthority(),
					ZoneID:    zone.BaseChainID(),
					Changes: []*icstypes.UpdateZoneValue{
						{
							Key:   icstypes.UpdateZoneKeyMessagesPerTx,
							Value: "",
						},
					},
				}
			},
			"invalid syntax",
		},
		{
			"valid: update account prefix",
			func() {
				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), testzoneID)
				suite.True(found)

				msg = &icstypes.MsgUpdateZone{
					Authority: suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetGovAuthority(),
					ZoneID:    zone.BaseChainID(),
					Changes: []*icstypes.UpdateZoneValue{
						{
							Key:   icstypes.UpdateZoneKeyAccountPrefix,
							Value: "test",
						},
					},
				}
			},
			"",
		},
		{
			"valid: update is 188",
			func() {
				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), testzoneID)
				suite.True(found)

				msg = &icstypes.MsgUpdateZone{
					Authority: suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetGovAuthority(),
					ZoneID:    zone.BaseChainID(),
					Changes: []*icstypes.UpdateZoneValue{
						{
							Key:   icstypes.UpdateZoneKeyIs118,
							Value: "false",
						},
					},
				}
			},
			"",
		},
		{
			"invalid: update is 188 invalid syntax",
			func() {
				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), testzoneID)
				suite.True(found)

				msg = &icstypes.MsgUpdateZone{
					Authority: suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetGovAuthority(),
					ZoneID:    zone.BaseChainID(),
					Changes: []*icstypes.UpdateZoneValue{
						{
							Key:   icstypes.UpdateZoneKeyIs118,
							Value: "",
						},
					},
				}
			},
			"invalid syntax",
		},
		{
			"invalid: update connection ID: already initialised",
			func() {
				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), testzoneID)
				suite.True(found)

				msg = &icstypes.MsgUpdateZone{
					Authority: suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetGovAuthority(),
					ZoneID:    zone.BaseChainID(),
					Changes: []*icstypes.UpdateZoneValue{
						{
							Key:   icstypes.UpdateZoneKeyConnectionID,
							Value: "connection-10",
						},
					},
				}
			},
			"zone already intialised",
		},
		{
			"invalid: update connection ID invalid syntax",
			func() {
				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), testzoneID)
				suite.True(found)

				msg = &icstypes.MsgUpdateZone{
					Authority: suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetGovAuthority(),
					ZoneID:    zone.BaseChainID(),
					Changes: []*icstypes.UpdateZoneValue{
						{
							Key:   icstypes.UpdateZoneKeyConnectionID,
							Value: "",
						},
					},
				}
			},
			"unexpected connection format",
		},
		{
			"invalid: unknown key",
			func() {
				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), testzoneID)
				suite.True(found)

				msg = &icstypes.MsgUpdateZone{
					Authority: suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetGovAuthority(),
					ZoneID:    zone.BaseChainID(),
					Changes: []*icstypes.UpdateZoneValue{
						{
							Key:   "invalid",
							Value: "",
						},
					},
				}
			},
			"unexpected key",
		},
	}
	for _, tt := range tests {
		tt := tt
		suite.Run(tt.name, func() {
			suite.SetupTest()
			suite.setupTestZones()

			tt.malleate()

			ctx := suite.chainA.GetContext()

			msgSrv := icskeeper.NewMsgServerImpl(*suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper)
			res, err := msgSrv.UpdateZone(sdk.WrapSDKContext(ctx), msg)

			if tt.expectErr != "" {
				suite.ErrorContains(err, tt.expectErr)
				suite.T().Logf("Error: %v", err)
			} else {
				suite.NoError(err)
				suite.NotNil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGovReopenChannel() {
	var msg *icstypes.MsgGovReopenChannel

	testAccount, err := addressutils.AccAddressFromBech32(testAddress, "")
	suite.NoError(err)

	tests := []struct {
		name      string
		malleate  func()
		expectErr string
	}{
		{
			"invalid: invalid connection ID",
			func() {
				msg = &icstypes.MsgGovReopenChannel{
					ConnectionId: "invalid",
					PortId:       "",
					Authority:    "",
				}
			},
			"unable to obtain chain id",
		},
		{
			"invalid: invalid connection ID",
			func() {
				msg = &icstypes.MsgGovReopenChannel{
					ConnectionId: suite.path.EndpointA.ConnectionID,
					PortId:       "",
					Authority:    "",
				}
			},
			"chainID / connectionID mismatch",
		},
		{
			"invalid: existing active channel",
			func() {
				msg = &icstypes.MsgGovReopenChannel{
					ConnectionId: suite.path.EndpointA.ConnectionID,
					PortId:       "testchain2-1.delegate",
					Authority:    "",
				}
			},
			"existing active channel",
		},
	}

	for _, tt := range tests {
		tt := tt
		suite.Run(tt.name, func() {
			suite.SetupTest()
			suite.setupTestZones()

			ctx := suite.chainA.GetContext()

			err := suite.GetQuicksilverApp(suite.chainA).BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqatom", math.NewInt(10000000))))
			suite.NoError(err)
			err = suite.GetQuicksilverApp(suite.chainA).BankKeeper.SendCoinsFromModuleToAccount(ctx, icstypes.ModuleName, testAccount, sdk.NewCoins(sdk.NewCoin("uqatom", math.NewInt(10000000))))
			suite.NoError(err)

			tt.malleate()

			msgSrv := icskeeper.NewMsgServerImpl(*suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper)
			res, err := msgSrv.GovReopenChannel(sdk.WrapSDKContext(suite.chainA.GetContext()), msg)

			if tt.expectErr != "" {
				suite.ErrorContains(err, tt.expectErr)
				suite.T().Logf("Error: %v", err)
			} else {
				suite.NoError(err)
				suite.NotNil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGovCloseChannel() {
	var msg *icstypes.MsgGovCloseChannel

	testAccount, err := addressutils.AccAddressFromBech32(testAddress, "")
	suite.NoError(err)

	tests := []struct {
		name      string
		malleate  func()
		expectErr string
	}{
		{
			"invalid: invalid authority",
			func() {
				msg = &icstypes.MsgGovCloseChannel{
					ChannelId: "",
					PortId:    "",
					Authority: "invalid",
				}
			},
			"invalid authority",
		},
		{
			"invalid: capability not found",
			func() {
				msg = &icstypes.MsgGovCloseChannel{
					ChannelId: "invalid",
					PortId:    "invalid",
					Authority: suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetGovAuthority(),
				}
			},
			"capability not found",
		},
		{
			"valid close",
			func() {
				msg = &icstypes.MsgGovCloseChannel{
					ChannelId: "channel-4",
					PortId:    "icacontroller-testchain2-1.deposit",
					Authority: suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetGovAuthority(),
				}
			},
			"",
		},
	}

	for _, tt := range tests {
		tt := tt
		suite.Run(tt.name, func() {
			suite.SetupTest()
			suite.setupTestZones()

			ctx := suite.chainA.GetContext()

			err := suite.GetQuicksilverApp(suite.chainA).BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqatom", math.NewInt(10000000))))
			suite.NoError(err)
			err = suite.GetQuicksilverApp(suite.chainA).BankKeeper.SendCoinsFromModuleToAccount(ctx, icstypes.ModuleName, testAccount, sdk.NewCoins(sdk.NewCoin("uqatom", math.NewInt(10000000))))
			suite.NoError(err)

			tt.malleate()

			msgSrv := icskeeper.NewMsgServerImpl(*suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper)
			res, err := msgSrv.GovCloseChannel(sdk.WrapSDKContext(suite.chainA.GetContext()), msg)

			if tt.expectErr != "" {
				suite.ErrorContains(err, tt.expectErr)
				suite.T().Logf("Error: %v", err)
			} else {
				suite.NoError(err)
				suite.NotNil(res)

				// verify channel is found but closed
				channel, found := suite.GetQuicksilverApp(suite.chainA).IBCKeeper.ChannelKeeper.GetChannel(ctx, msg.PortId, msg.ChannelId)
				suite.True(found)
				suite.Equal(channeltypes.CLOSED, channel.State)
			}
		})
	}
}

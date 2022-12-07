package keeper_test

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"

	"github.com/ingenuity-build/quicksilver/utils"
	icskeeper "github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func (s *KeeperTestSuite) TestRequestRedemption() {
	var msg icstypes.MsgRequestRedemption

	testAccount, err := utils.AccAddressFromBech32(testAddress, "")
	s.Require().NoError(err)

	tests := []struct {
		name      string
		malleate  func()
		expectErr bool
	}{
		{
			"valid - full claim",
			func() {
				addr, err := bech32.ConvertAndEncode("cosmos", utils.GenerateAccAddressForTest())
				s.Require().NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uqatom", sdk.NewInt(10000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}
			},
			false,
		},
		{
			"valid - full claim (discounted)",
			func() {
				addr, err := bech32.ConvertAndEncode("cosmos", utils.GenerateAccAddressForTest())
				s.Require().NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uqatom", sdk.NewInt(10000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}

				zone, found := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.GetZone(s.chainA.GetContext(), s.chainB.ChainID)
				s.Require().True(found)
				zone.RedemptionRate = sdk.MustNewDecFromStr("0.95")
				s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)
			},
			false,
		},
		{
			"valid - full claim (interest)",
			func() {
				addr, err := bech32.ConvertAndEncode("cosmos", utils.GenerateAccAddressForTest())
				s.Require().NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uqatom", sdk.NewInt(10000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}

				zone, found := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.GetZone(s.chainA.GetContext(), s.chainB.ChainID)
				s.Require().True(found)
				zone.LastRedemptionRate = sdk.MustNewDecFromStr("1.05")
				zone.RedemptionRate = sdk.MustNewDecFromStr("1.1")
				s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)
			},
			false,
		},
		{
			"valid - full claim (interest)",
			func() {
				addr, err := bech32.ConvertAndEncode("cosmos", utils.GenerateAccAddressForTest())
				s.Require().NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uqatom", sdk.NewInt(10000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}

				zone, found := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.GetZone(s.chainA.GetContext(), s.chainB.ChainID)
				s.Require().True(found)
				zone.LastRedemptionRate = sdk.MustNewDecFromStr("1.1")
				zone.RedemptionRate = sdk.MustNewDecFromStr("1.05")
				s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)
			},
			false,
		},
		{
			"valid - partial claim",
			func() {
				addr, err := bech32.ConvertAndEncode("cosmos", utils.GenerateAccAddressForTest())
				s.Require().NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uqatom", sdk.NewInt(5000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}
			},
			false,
		},
		{
			"valid - partial claim (discounted)",
			func() {
				addr, err := bech32.ConvertAndEncode("cosmos", utils.GenerateAccAddressForTest())
				s.Require().NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uqatom", sdk.NewInt(5000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}

				zone, found := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.GetZone(s.chainA.GetContext(), s.chainB.ChainID)
				s.Require().True(found)
				zone.RedemptionRate = sdk.MustNewDecFromStr("0.99999")
				s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)
			},
			false,
		},
		{
			"valid - partial claim (interest)",
			func() {
				addr, err := bech32.ConvertAndEncode("cosmos", utils.GenerateAccAddressForTest())
				s.Require().NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uqatom", sdk.NewInt(5000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}

				zone, found := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.GetZone(s.chainA.GetContext(), s.chainB.ChainID)
				s.Require().True(found)
				zone.LastRedemptionRate = sdk.MustNewDecFromStr("1.049999")
				zone.RedemptionRate = sdk.MustNewDecFromStr("1.099999")
				s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)
			},
			false,
		},
		{
			"invalid - wrong denom",
			func() {
				addr, err := bech32.ConvertAndEncode("cosmos", utils.GenerateAccAddressForTest())
				s.Require().NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uatom", sdk.NewInt(10000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}
			},
			true,
		},
		{
			"invalid - insufficient funds",
			func() {
				addr, err := bech32.ConvertAndEncode("cosmos", utils.GenerateAccAddressForTest())
				s.Require().NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uqatom", sdk.NewInt(1000000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}
			},
			true,
		},
		{
			"invalid - zero coins",
			func() {
				addr, err := bech32.ConvertAndEncode("cosmos", utils.GenerateAccAddressForTest())
				s.Require().NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uqatom", sdk.ZeroInt()),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}
			},
			true,
		},
		{
			"invalid - negative coins",
			func() {
				addr, err := bech32.ConvertAndEncode("cosmos", utils.GenerateAccAddressForTest())
				s.Require().NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.Coin{Denom: "uqatom", Amount: sdk.NewInt(-1)},
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}
			},
			true,
		},
		{
			"invalid - bad prefix",
			func() {
				addr, err := bech32.ConvertAndEncode("bob", utils.GenerateAccAddressForTest())
				s.Require().NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uqatom", sdk.OneInt()),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}
			},
			true,
		},
		{
			"invalid - bad from address",
			func() {
				addr, err := bech32.ConvertAndEncode("cosmos", utils.GenerateAccAddressForTest())
				s.Require().NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uqatom", sdk.OneInt()),
					DestinationAddress: addr,
					FromAddress:        addr,
				}
			},
			true,
		},
		{
			"invalid - nil recipient address",
			func() {
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uqatom", sdk.OneInt()),
					DestinationAddress: "",
					FromAddress:        testAddress,
				}
			},
			true,
		},
		{
			"invalid - nil from address",
			func() {
				addr, err := bech32.ConvertAndEncode("cosmos", utils.GenerateAccAddressForTest())
				s.Require().NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uqatom", sdk.OneInt()),
					DestinationAddress: addr,
					FromAddress:        "",
				}
			},
			true,
		},
	}

	for _, tt := range tests {
		tt := tt

		// run tests with LSM disabled.
		s.Run(tt.name, func() {
			s.SetupTest()
			s.setupTestZones()

			ctx := s.chainA.GetContext()

			s.GetQuicksilverApp(s.chainA).BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqatom", math.NewInt(10000000))))
			s.GetQuicksilverApp(s.chainA).BankKeeper.SendCoinsFromModuleToAccount(ctx, icstypes.ModuleName, testAccount, sdk.NewCoins(sdk.NewCoin("uqatom", math.NewInt(10000000))))

			// disable LSM
			zone, found := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
			s.Require().True(found)
			zone.LiquidityModule = false
			s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(ctx, &zone)

			tt.malleate()

			msgSrv := icskeeper.NewMsgServerImpl(s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper)
			res, err := msgSrv.RequestRedemption(sdk.WrapSDKContext(s.chainA.GetContext()), &msg)

			if tt.expectErr {
				s.Require().Error(err)
				s.Require().Nil(res)
				s.T().Logf("Error: %v", err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})

		// run tests with LSM enabled.
		tt.name = tt.name + "_LSM_enabled"
		s.Run(tt.name, func() {
			s.SetupTest()
			s.setupTestZones()

			ctx := s.chainA.GetContext()

			s.GetQuicksilverApp(s.chainA).BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqatom", math.NewInt(10000000))))
			s.GetQuicksilverApp(s.chainA).BankKeeper.SendCoinsFromModuleToAccount(ctx, icstypes.ModuleName, testAccount, sdk.NewCoins(sdk.NewCoin("uqatom", math.NewInt(10000000))))

			// enable LSM
			zone, found := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
			s.Require().True(found)
			zone.LiquidityModule = true
			s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(ctx, &zone)

			tt.malleate()

			msgSrv := icskeeper.NewMsgServerImpl(s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper)
			res, err := msgSrv.RequestRedemption(sdk.WrapSDKContext(s.chainA.GetContext()), &msg)

			if tt.expectErr {
				s.Require().Error(err)
				s.Require().Nil(res)
				s.T().Logf("Error: %v", err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})

	}
}

func (s *KeeperTestSuite) TestSignalIntent() {
	tests := []struct {
		name             string
		malleate         func(s *KeeperTestSuite) *icstypes.MsgSignalIntent
		expected         []sdk.Dec
		failsValidations bool
		expectErr        bool
	}{
		{
			"invalid - weight sum < 1",
			func(s *KeeperTestSuite) *icstypes.MsgSignalIntent {
				val1, err := sdk.ValAddressFromHex(s.chainB.Vals.Validators[0].Address.String())
				s.Require().NoError(err)

				return &icstypes.MsgSignalIntent{
					ChainId:     s.chainB.ChainID,
					Intents:     fmt.Sprintf("0.3%s", val1.String()),
					FromAddress: testAddress,
				}
			},
			[]sdk.Dec{},
			true,
			false,
		},
		{
			"invalid - weight sum > 1",
			func(s *KeeperTestSuite) *icstypes.MsgSignalIntent {
				val1, err := sdk.ValAddressFromHex(s.chainB.Vals.Validators[0].Address.String())
				s.Require().NoError(err)

				return &icstypes.MsgSignalIntent{
					ChainId:     s.chainB.ChainID,
					Intents:     fmt.Sprintf("3.0%s", val1.String()),
					FromAddress: testAddress,
				}
			},
			[]sdk.Dec{},
			true,
			false,
		},
		{
			"invalid - chain id",
			func(s *KeeperTestSuite) *icstypes.MsgSignalIntent {
				val1, err := sdk.ValAddressFromHex(s.chainB.Vals.Validators[0].Address.String())
				s.Require().NoError(err)

				return &icstypes.MsgSignalIntent{
					ChainId:     s.chainA.ChainID,
					Intents:     fmt.Sprintf("1.0%s", val1.String()),
					FromAddress: testAddress,
				}
			},
			[]sdk.Dec{},
			false,
			true,
		},
		{
			"valid - single weight",
			func(s *KeeperTestSuite) *icstypes.MsgSignalIntent {
				val1, err := sdk.ValAddressFromHex(s.chainB.Vals.Validators[0].Address.String())
				s.Require().NoError(err)

				return &icstypes.MsgSignalIntent{
					ChainId:     s.chainB.ChainID,
					Intents:     fmt.Sprintf("1.0%s", val1.String()),
					FromAddress: testAddress,
				}
			},
			[]sdk.Dec{sdk.NewDecWithPrec(1, 0)},
			false,
			false,
		},
		{
			"valid - multi weight",
			func(s *KeeperTestSuite) *icstypes.MsgSignalIntent {
				val1, err := sdk.ValAddressFromHex(s.chainB.Vals.Validators[0].Address.String())
				s.Require().NoError(err)
				val2, err := sdk.ValAddressFromHex(s.chainB.Vals.Validators[1].Address.String())
				s.Require().NoError(err)
				val3, err := sdk.ValAddressFromHex(s.chainB.Vals.Validators[2].Address.String())
				s.Require().NoError(err)

				return &icstypes.MsgSignalIntent{
					ChainId:     s.chainB.ChainID,
					Intents:     fmt.Sprintf("0.5%s,0.2%s,0.3%s", val1.String(), val2.String(), val3.String()),
					FromAddress: testAddress,
				}
			},
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

		s.Run(tt.name, func() {
			s.SetupTest()
			s.setupTestZones()

			msg := tt.malleate(s)
			// validateBasic not explicitly tested here - but we don't call it inside msgSrv.SignalIntent
			// so call here to make sure out tests are sane.
			err := msg.ValidateBasic()
			if tt.failsValidations {
				s.Require().Error(err)
				return
			} else {
				s.Require().NoError(err)
			}

			msgSrv := icskeeper.NewMsgServerImpl(s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper)
			res, err := msgSrv.SignalIntent(sdk.WrapSDKContext(s.chainA.GetContext()), msg)
			if tt.expectErr {
				s.Require().Error(err)
				s.Require().Nil(res)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(res)
			}

			qapp := s.GetQuicksilverApp(s.chainA)
			icsKeeper := qapp.InterchainstakingKeeper
			zone, found := icsKeeper.GetZone(s.chainA.GetContext(), s.chainB.ChainID)
			s.Require().True(found)

			intent, found := icsKeeper.GetIntent(s.chainA.GetContext(), zone, testAddress, false)
			s.Require().True(found)
			intents := intent.GetIntents()

			for idx, weight := range tt.expected {
				val, err := sdk.ValAddressFromHex(s.chainB.Vals.Validators[idx].Address.String())
				s.Require().NoError(err)

				valIntent, found := intents.GetForValoper(val.String())
				s.Require().True(found)

				s.Require().Equal(weight, valIntent.Weight)
			}
		})
	}
}

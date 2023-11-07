package keeper_test

import (
	"fmt"
	"time"

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
	s.NoError(err)

	tests := []struct {
		name         string
		malleate     func()
		expectErr    string
		expectErrLsm string
	}{
		{
			"valid - full claim",
			func() {
				addr, err := bech32.ConvertAndEncode("cosmos", utils.GenerateAccAddressForTest())
				s.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uqatom", sdk.NewInt(10000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}
			},
			"",
			"",
		},
		{
			"valid - full claim (discounted)",
			func() {
				addr, err := bech32.ConvertAndEncode("cosmos", utils.GenerateAccAddressForTest())
				s.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uqatom", sdk.NewInt(10000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}

				zone, found := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.GetZone(s.chainA.GetContext(), s.chainB.ChainID)
				s.True(found)
				zone.RedemptionRate = sdk.MustNewDecFromStr("0.95")
				s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)
			},
			"",
			"",
		},
		{
			"valid - full claim (interest)",
			func() {
				addr, err := bech32.ConvertAndEncode("cosmos", utils.GenerateAccAddressForTest())
				s.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uqatom", sdk.NewInt(10000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}

				zone, found := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.GetZone(s.chainA.GetContext(), s.chainB.ChainID)
				s.True(found)
				zone.LastRedemptionRate = sdk.MustNewDecFromStr("1.05")
				zone.RedemptionRate = sdk.MustNewDecFromStr("1.1")
				s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)
			},
			"",
			"",
		},
		{
			"valid - full claim (interest)",
			func() {
				addr, err := bech32.ConvertAndEncode("cosmos", utils.GenerateAccAddressForTest())
				s.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uqatom", sdk.NewInt(10000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}

				zone, found := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.GetZone(s.chainA.GetContext(), s.chainB.ChainID)
				s.True(found)
				zone.LastRedemptionRate = sdk.MustNewDecFromStr("1.1")
				zone.RedemptionRate = sdk.MustNewDecFromStr("1.05")
				s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)
			},
			"",
			"",
		},
		{
			"valid - partial claim",
			func() {
				addr, err := bech32.ConvertAndEncode("cosmos", utils.GenerateAccAddressForTest())
				s.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uqatom", sdk.NewInt(5000000)),
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
				addr, err := bech32.ConvertAndEncode("cosmos", utils.GenerateAccAddressForTest())
				s.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uqatom", sdk.NewInt(5000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}

				zone, found := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.GetZone(s.chainA.GetContext(), s.chainB.ChainID)
				s.True(found)
				zone.RedemptionRate = sdk.MustNewDecFromStr("0.99999")
				s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)
			},
			"",
			"",
		},
		{
			"valid - partial claim (interest)",
			func() {
				addr, err := bech32.ConvertAndEncode("cosmos", utils.GenerateAccAddressForTest())
				s.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uqatom", sdk.NewInt(5000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}

				zone, found := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.GetZone(s.chainA.GetContext(), s.chainB.ChainID)
				s.True(found)
				zone.LastRedemptionRate = sdk.MustNewDecFromStr("1.049999")
				zone.RedemptionRate = sdk.MustNewDecFromStr("1.099999")
				s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)
			},
			"",
			"",
		},
		{
			"invalid - wrong denom",
			func() {
				addr, err := bech32.ConvertAndEncode("cosmos", utils.GenerateAccAddressForTest())
				s.NoError(err)
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
				addr, err := bech32.ConvertAndEncode("cosmos", utils.GenerateAccAddressForTest())
				s.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uqatom", sdk.NewInt(1000000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}
			},
			"account has insufficient balance of qasset to burn",
			"account has insufficient balance of qasset to burn",
		},
		{
			"invalid - zero coins",
			func() {
				addr, err := bech32.ConvertAndEncode("cosmos", utils.GenerateAccAddressForTest())
				s.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uqatom", sdk.ZeroInt()),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}
			},
			"cannot redeem zero-value coins",
			"cannot redeem zero-value coins",
		},
		{
			"invalid - negative coins",
			func() {
				addr, err := bech32.ConvertAndEncode("cosmos", utils.GenerateAccAddressForTest())
				s.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.Coin{Denom: "uqatom", Amount: sdk.NewInt(-1)},
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}
			},
			"negative coin amount: -1",
			"negative coin amount: -1",
		},
		{
			"invalid - bad prefix",
			func() {
				addr, err := bech32.ConvertAndEncode("bob", utils.GenerateAccAddressForTest())
				s.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uqatom", sdk.OneInt()),
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
				addr, err := bech32.ConvertAndEncode("cosmos", utils.GenerateAccAddressForTest())
				s.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uqatom", sdk.OneInt()),
					DestinationAddress: addr,
					FromAddress:        addr,
				}
			},
			"account has insufficient balance of qasset to burn",
			"account has insufficient balance of qasset to burn",
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
			"recipient address not provided",
			"recipient address not provided",
		},
		{
			"invalid - nil from address",
			func() {
				addr, err := bech32.ConvertAndEncode("cosmos", utils.GenerateAccAddressForTest())
				s.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uqatom", sdk.OneInt()),
					DestinationAddress: addr,
					FromAddress:        "",
				}
			},
			"empty address string is not allowed",
			"empty address string is not allowed",
		},
		{
			"invalid - too many locked tokens",
			func() {
				addr, err := bech32.ConvertAndEncode("cosmos", utils.GenerateAccAddressForTest())
				s.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uqatom", sdk.NewInt(10000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}

				zone, _ := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.GetZone(s.chainA.GetContext(), s.chainB.ChainID)
				s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetRedelegationRecord(s.chainA.GetContext(), icstypes.RedelegationRecord{
					ChainId:        s.chainB.ChainID,
					EpochNumber:    1,
					Source:         zone.GetValidatorsAddressesAsSlice()[0],
					Destination:    zone.GetValidatorsAddressesAsSlice()[1],
					Amount:         3000000,
					CompletionTime: time.Time(s.chainA.GetContext().BlockTime().Add(time.Hour)),
				})
			},
			"",
			"unable to satisfy unbond request; delegations may be locked",
		},
	}

	for _, tt := range tests {
		tt := tt

		// run tests with LSM disabled.
		s.Run(tt.name, func() {
			s.SetupTest()
			s.setupTestZones()

			ctx := s.chainA.GetContext()

			params := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.GetParams(ctx)
			params.UnbondingEnabled = true
			s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetParams(ctx, params)

			s.GetQuicksilverApp(s.chainA).BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqatom", math.NewInt(10000000))))
			s.GetQuicksilverApp(s.chainA).BankKeeper.SendCoinsFromModuleToAccount(ctx, icstypes.ModuleName, testAccount, sdk.NewCoins(sdk.NewCoin("uqatom", math.NewInt(10000000))))

			// disable LSM
			zone, found := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
			s.True(found)
			zone.LiquidityModule = false
			s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(ctx, &zone)

			tt.malleate()

			msgSrv := icskeeper.NewMsgServerImpl(s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper)
			res, err := msgSrv.RequestRedemption(sdk.WrapSDKContext(s.chainA.GetContext()), &msg)

			if tt.expectErr != "" {
				s.ErrorContains(err, tt.expectErr)
				s.Nil(res)
				s.T().Logf("Error: %v", err)
			} else {
				s.NoError(err)
				s.NotNil(res)
			}
		})

		// run tests with LSM enabled.
		tt.name = tt.name + "_LSM_enabled"
		s.Run(tt.name, func() {
			s.SetupTest()
			s.setupTestZones()

			ctx := s.chainA.GetContext()

			params := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.GetParams(ctx)
			params.UnbondingEnabled = true
			s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetParams(ctx, params)

			s.GetQuicksilverApp(s.chainA).BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqatom", math.NewInt(10000000))))
			s.GetQuicksilverApp(s.chainA).BankKeeper.SendCoinsFromModuleToAccount(ctx, icstypes.ModuleName, testAccount, sdk.NewCoins(sdk.NewCoin("uqatom", math.NewInt(10000000))))

			// enable LSM
			zone, found := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
			s.True(found)
			zone.LiquidityModule = true
			s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(ctx, &zone)

			for _, delegation := range func(zone icstypes.Zone) []icstypes.Delegation {
				validators := zone.GetValidatorsAddressesAsSlice()
				out := make([]icstypes.Delegation, 0)
				for _, valoper := range validators {
					out = append(out, icstypes.NewDelegation(zone.DelegationAddress.Address, valoper, sdk.NewCoin(zone.BaseDenom, sdk.NewInt(3000000))))
				}
				return out
			}(zone) {
				s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetDelegation(ctx, &zone, delegation)
			}

			tt.malleate()

			msgSrv := icskeeper.NewMsgServerImpl(s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper)
			res, err := msgSrv.RequestRedemption(sdk.WrapSDKContext(s.chainA.GetContext()), &msg)

			if tt.expectErrLsm != "" {
				s.Errorf(err, tt.expectErrLsm)
				s.Nil(res)
				s.T().Logf("Error: %v", err)
			} else {
				s.NoError(err)
				s.NotNil(res)
			}
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
				s.NoError(err)

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
				s.NoError(err)

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
				s.NoError(err)

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
				s.NoError(err)

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
				s.NoError(err)
				val2, err := sdk.ValAddressFromHex(s.chainB.Vals.Validators[1].Address.String())
				s.NoError(err)
				val3, err := sdk.ValAddressFromHex(s.chainB.Vals.Validators[2].Address.String())
				s.NoError(err)

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
				s.Error(err)
				return
			} else {
				s.NoError(err)
			}

			msgSrv := icskeeper.NewMsgServerImpl(s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper)
			res, err := msgSrv.SignalIntent(sdk.WrapSDKContext(s.chainA.GetContext()), msg)
			if tt.expectErr {
				s.Error(err)
				s.Nil(res)
			} else {
				s.NoError(err)
				s.NotNil(res)
			}

			qapp := s.GetQuicksilverApp(s.chainA)
			icsKeeper := qapp.InterchainstakingKeeper
			zone, found := icsKeeper.GetZone(s.chainA.GetContext(), s.chainB.ChainID)
			s.True(found)

			intent, found := icsKeeper.GetIntent(s.chainA.GetContext(), zone, testAddress, false)
			s.True(found)
			intents := intent.GetIntents()

			for idx, weight := range tt.expected {
				val, err := sdk.ValAddressFromHex(s.chainB.Vals.Validators[idx].Address.String())
				s.NoError(err)

				valIntent, found := intents.GetForValoper(val.String())
				s.True(found)

				s.Equal(weight, valIntent.Weight)
			}
		})
	}
}

func (s *KeeperTestSuite) TestSetLsmCaps() {
	tests := []struct {
		name      string
		malleate  func(s *KeeperTestSuite) *icstypes.MsgGovSetLsmCaps
		expectErr bool
	}{

		{
			"invalid authority",
			func(s *KeeperTestSuite) *icstypes.MsgGovSetLsmCaps {
				return &icstypes.MsgGovSetLsmCaps{
					ChainId: s.chainB.ChainID,
					Caps: &icstypes.LsmCaps{
						ValidatorCap:     sdk.NewDecWithPrec(50, 2),
						ValidatorBondCap: sdk.NewDec(250),
						GlobalCap:        sdk.NewDecWithPrec(25, 2),
					},
					Authority: testAddress,
				}
			},
			true,
		},
		{
			"valid",
			func(s *KeeperTestSuite) *icstypes.MsgGovSetLsmCaps {
				return &icstypes.MsgGovSetLsmCaps{
					ChainId: s.chainB.ChainID,
					Caps: &icstypes.LsmCaps{
						ValidatorCap:     sdk.NewDecWithPrec(50, 2),
						ValidatorBondCap: sdk.NewDec(250),
						GlobalCap:        sdk.NewDecWithPrec(25, 2),
					},
					Authority: "quick10d07y265gmmuvt4z0w9aw880jnsr700j3xrh0p",
				}
			},
			false,
		},
	}

	for _, tt := range tests {
		tt := tt

		s.Run(tt.name, func() {
			s.SetupTest()
			s.setupTestZones()

			msg := tt.malleate(s)

			msgSrv := icskeeper.NewMsgServerImpl(s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper)
			res, err := msgSrv.GovSetLsmCaps(sdk.WrapSDKContext(s.chainA.GetContext()), msg)
			if tt.expectErr {
				s.Error(err)
				s.Nil(res)
			} else {
				s.NoError(err)
				s.NotNil(res)
			}

			qapp := s.GetQuicksilverApp(s.chainA)
			icsKeeper := qapp.InterchainstakingKeeper
			zone, found := icsKeeper.GetZone(s.chainA.GetContext(), s.chainB.ChainID)
			s.True(found)

			caps, found := icsKeeper.GetLsmCaps(s.chainA.GetContext(), zone.ChainId)
			if tt.expectErr {
				s.False(found)
				s.Nil(caps)
			} else {
				s.True(found)
				s.Equal(caps, msg.Caps)

			}
		})
	}
}

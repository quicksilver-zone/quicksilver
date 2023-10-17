package keeper_test

import (
	"fmt"
	"time"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	connectiontypes "github.com/cosmos/ibc-go/v5/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v5/modules/core/04-channel/types"

	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	icskeeper "github.com/quicksilver-zone/quicksilver/x/interchainstaking/keeper"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

func (suite *KeeperTestSuite) TestRequestRedemption() {
	var msg icstypes.MsgRequestRedemption

	testAccount, err := addressutils.AccAddressFromBech32(testAddress, "")
	suite.NoError(err)

	tests := []struct {
		name         string
		malleate     func()
		expectErr    string
		expectErrLsm string
	}{
		{
			"valid - full claim",
			func() {
				addr, err := addressutils.EncodeAddressToBech32("cosmos", addressutils.GenerateAccAddressForTest())
				suite.NoError(err)
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
				addr, err := addressutils.EncodeAddressToBech32("cosmos", addressutils.GenerateAccAddressForTest())
				suite.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uqatom", sdk.NewInt(10000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}

				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), suite.chainB.ChainID)
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
				addr, err := addressutils.EncodeAddressToBech32("cosmos", addressutils.GenerateAccAddressForTest())
				suite.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uqatom", sdk.NewInt(10000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}

				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), suite.chainB.ChainID)
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
				addr, err := addressutils.EncodeAddressToBech32("cosmos", addressutils.GenerateAccAddressForTest())
				suite.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uqatom", sdk.NewInt(10000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}

				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), suite.chainB.ChainID)
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
				addr, err := addressutils.EncodeAddressToBech32("cosmos", addressutils.GenerateAccAddressForTest())
				suite.NoError(err)
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
				addr, err := addressutils.EncodeAddressToBech32("cosmos", addressutils.GenerateAccAddressForTest())
				suite.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uqatom", sdk.NewInt(5000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}

				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), suite.chainB.ChainID)
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
				addr, err := addressutils.EncodeAddressToBech32("cosmos", addressutils.GenerateAccAddressForTest())
				suite.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uqatom", sdk.NewInt(5000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}

				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), suite.chainB.ChainID)
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
				addr, err := addressutils.EncodeAddressToBech32("cosmos", addressutils.GenerateAccAddressForTest())
				suite.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uqatom", sdk.NewInt(5000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}

				zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), suite.chainB.ChainID)
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
				addr, err := addressutils.EncodeAddressToBech32("cosmos", addressutils.GenerateAccAddressForTest())
				suite.NoError(err)
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
			"invalid - bad prefix",
			func() {
				addr, err := addressutils.EncodeAddressToBech32("bob", addressutils.GenerateAccAddressForTest())
				suite.NoError(err)
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
				addr, err := addressutils.EncodeAddressToBech32("cosmos", addressutils.GenerateAccAddressForTest())
				suite.NoError(err)
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
			"invalid - too many locked tokens",
			func() {
				addr, err := addressutils.EncodeAddressToBech32("cosmos", addressutils.GenerateAccAddressForTest())
				suite.NoError(err)
				msg = icstypes.MsgRequestRedemption{
					Value:              sdk.NewCoin("uqatom", sdk.NewInt(10000000)),
					DestinationAddress: addr,
					FromAddress:        testAddress,
				}

				ctx := suite.chainA.GetContext()
				zoneVals := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetValidatorAddresses(ctx, suite.chainB.ChainID)
				suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetRedelegationRecord(ctx, icstypes.RedelegationRecord{
					ChainId:        suite.chainB.ChainID,
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
	}

	for _, tt := range tests {
		tt := tt

		// run tests with LSM disabled.
		suite.Run(tt.name, func() {
			suite.SetupTest()
			suite.setupTestZones()

			ctx := suite.chainA.GetContext()

			params := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetParams(ctx)
			params.UnbondingEnabled = true
			suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetParams(ctx, params)

			err := suite.GetQuicksilverApp(suite.chainA).BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqatom", math.NewInt(10000000))))
			suite.NoError(err)
			err = suite.GetQuicksilverApp(suite.chainA).BankKeeper.SendCoinsFromModuleToAccount(ctx, icstypes.ModuleName, testAccount, sdk.NewCoins(sdk.NewCoin("uqatom", math.NewInt(10000000))))
			suite.NoError(err)

			// disable LSM
			zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
			suite.True(found)
			zone.LiquidityModule = false
			zone.UnbondingEnabled = true
			suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetZone(ctx, &zone)

			tt.malleate()

			msgSrv := icskeeper.NewMsgServerImpl(*suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper)
			res, err := msgSrv.RequestRedemption(sdk.WrapSDKContext(suite.chainA.GetContext()), &msg)

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

			ctx := suite.chainA.GetContext()

			params := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetParams(ctx)
			params.UnbondingEnabled = true
			suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetParams(ctx, params)

			err := suite.GetQuicksilverApp(suite.chainA).BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqatom", math.NewInt(10000000))))
			suite.NoError(err)
			err = suite.GetQuicksilverApp(suite.chainA).BankKeeper.SendCoinsFromModuleToAccount(ctx, icstypes.ModuleName, testAccount, sdk.NewCoins(sdk.NewCoin("uqatom", math.NewInt(10000000))))
			suite.NoError(err)

			// enable LSM
			zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
			suite.True(found)
			zone.LiquidityModule = true
			zone.UnbondingEnabled = true
			suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetZone(ctx, &zone)

			validators := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetValidatorAddresses(ctx, suite.chainB.ChainID)
			for _, delegation := range func(zone icstypes.Zone) []icstypes.Delegation {
				out := make([]icstypes.Delegation, 0)
				for _, valoper := range validators {
					out = append(out, icstypes.NewDelegation(zone.DelegationAddress.Address, valoper, sdk.NewCoin(zone.BaseDenom, sdk.NewInt(3000000))))
				}
				return out
			}(zone) {
				suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetDelegation(ctx, &zone, delegation)
			}

			tt.malleate()

			msgSrv := icskeeper.NewMsgServerImpl(*suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper)
			res, err := msgSrv.RequestRedemption(sdk.WrapSDKContext(suite.chainA.GetContext()), &msg)

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
					ChainId:     suite.chainB.ChainID,
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
			func(suite *KeeperTestSuite) *icstypes.MsgSignalIntent {
				val1, err := sdk.ValAddressFromHex(suite.chainB.Vals.Validators[0].Address.String())
				suite.NoError(err)

				return &icstypes.MsgSignalIntent{
					ChainId:     suite.chainB.ChainID,
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
			func(suite *KeeperTestSuite) *icstypes.MsgSignalIntent {
				val1, err := sdk.ValAddressFromHex(suite.chainB.Vals.Validators[0].Address.String())
				suite.NoError(err)

				return &icstypes.MsgSignalIntent{
					ChainId:     suite.chainA.ChainID,
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
			func(suite *KeeperTestSuite) *icstypes.MsgSignalIntent {
				val1, err := sdk.ValAddressFromHex(suite.chainB.Vals.Validators[0].Address.String())
				suite.NoError(err)

				return &icstypes.MsgSignalIntent{
					ChainId:     suite.chainB.ChainID,
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
			func(suite *KeeperTestSuite) *icstypes.MsgSignalIntent {
				val1, err := sdk.ValAddressFromHex(suite.chainB.Vals.Validators[0].Address.String())
				suite.NoError(err)
				val2, err := sdk.ValAddressFromHex(suite.chainB.Vals.Validators[1].Address.String())
				suite.NoError(err)
				val3, err := sdk.ValAddressFromHex(suite.chainB.Vals.Validators[2].Address.String())
				suite.NoError(err)

				return &icstypes.MsgSignalIntent{
					ChainId:     suite.chainB.ChainID,
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
			} else {
				suite.NoError(err)
				suite.NotNil(res)
			}

			quicksilver := suite.GetQuicksilverApp(suite.chainA)
			icsKeeper := quicksilver.InterchainstakingKeeper
			zone, found := icsKeeper.GetZone(suite.chainA.GetContext(), suite.chainB.ChainID)
			suite.True(found)

			intent, found := icsKeeper.GetDelegatorIntent(suite.chainA.GetContext(), &zone, testAddress, false)
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

func (suite *KeeperTestSuite) TestGovCloseChannel() {
	testCase := []struct {
		name     string
		malleate func(suite *KeeperTestSuite) *icstypes.MsgGovCloseChannel
		expecErr error
	}{
		{
			name: "invalid authority",
			malleate: func(suite *KeeperTestSuite) *icstypes.MsgGovCloseChannel {
				return &icstypes.MsgGovCloseChannel{
					ChannelId: "",
					PortId:    "",
					Authority: testAddress,
				}
			},
			expecErr: govtypes.ErrInvalidSigner,
		},
		{
			name: "capability not found",
			malleate: func(suite *KeeperTestSuite) *icstypes.MsgGovCloseChannel {
				k := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper

				return &icstypes.MsgGovCloseChannel{
					ChannelId: "",
					PortId:    "",
					Authority: sdk.MustBech32ifyAddressBytes(sdk.GetConfig().GetBech32AccountAddrPrefix(), k.AccountKeeper.GetModuleAddress(govtypes.ModuleName)),
				}
			},
			expecErr: capabilitytypes.ErrCapabilityNotFound,
		},
		{
			name: "invalid connection state",
			malleate: func(suite *KeeperTestSuite) *icstypes.MsgGovCloseChannel {
				ctx := suite.chainA.GetContext()
				k := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
				channels := suite.GetQuicksilverApp(suite.chainA).IBCKeeper.ChannelKeeper.GetAllChannels(ctx)

				return &icstypes.MsgGovCloseChannel{
					ChannelId: channels[0].ChannelId,
					PortId:    channels[0].PortId,
					Authority: sdk.MustBech32ifyAddressBytes(sdk.GetConfig().GetBech32AccountAddrPrefix(), k.AccountKeeper.GetModuleAddress(govtypes.ModuleName)),
				}
			},
			expecErr: connectiontypes.ErrInvalidConnectionState,
		},
		{
			name: "closes an ICA channel success",
			malleate: func(suite *KeeperTestSuite) *icstypes.MsgGovCloseChannel {
				ctx := suite.chainA.GetContext()
				suite.GetQuicksilverApp(suite.chainA).IBCKeeper.ConnectionKeeper.SetConnection(ctx, suite.path.EndpointA.ConnectionID, connectiontypes.ConnectionEnd{ClientId: "07-tendermint-0", State: connectiontypes.OPEN})
				k := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
				channels := suite.GetQuicksilverApp(suite.chainA).IBCKeeper.ChannelKeeper.GetAllChannels(ctx)

				return &icstypes.MsgGovCloseChannel{
					ChannelId: channels[0].ChannelId,
					PortId:    channels[0].PortId,
					Authority: sdk.MustBech32ifyAddressBytes(sdk.GetConfig().GetBech32AccountAddrPrefix(), k.AccountKeeper.GetModuleAddress(govtypes.ModuleName)),
				}
			},
			expecErr: nil,
		},
	}
	for _, tc := range testCase {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			suite.setupTestZones()

			msg := tc.malleate(suite)
			msgSrv := icskeeper.NewMsgServerImpl(*suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper)
			ctx := suite.chainA.GetContext()

			_, err := msgSrv.GovCloseChannel(ctx, msg)
			if tc.expecErr != nil {
				suite.ErrorIs(tc.expecErr, err)
				return
			}
			suite.NoError(err)

			// check state channel is CLOSED
			channel, found := suite.GetQuicksilverApp(suite.chainA).IBCKeeper.ChannelKeeper.GetChannel(ctx, msg.PortId, msg.ChannelId)
			suite.True(found)
			suite.True(channel.State == channeltypes.CLOSED)
		})
	}
}

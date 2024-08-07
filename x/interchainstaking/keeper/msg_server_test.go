package keeper_test

import (
	"errors"
	"fmt"
	"time"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v6/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	tmclienttypes "github.com/cosmos/ibc-go/v6/modules/light-clients/07-tendermint/types"

	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	"github.com/quicksilver-zone/quicksilver/utils/randomutils"
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
					Amount:         math.NewInt(3000000),
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

			msgSrv := icskeeper.NewMsgServerImpl(suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper)
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

		// run tests with LSM enabled.- disabled until we decide to use LSM unbonding.
		// tt.name += "_LSM_enabled"
		// suite.Run(tt.name, func() {
		// 	suite.SetupTest()
		// 	suite.setupTestZones()

		// 	ctx := suite.chainA.GetContext()

		// 	params := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetParams(ctx)
		// 	params.UnbondingEnabled = true
		// 	suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetParams(ctx, params)

		// 	err := suite.GetQuicksilverApp(suite.chainA).BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqatom", math.NewInt(10000000))))
		// 	suite.NoError(err)
		// 	err = suite.GetQuicksilverApp(suite.chainA).BankKeeper.SendCoinsFromModuleToAccount(ctx, icstypes.ModuleName, testAccount, sdk.NewCoins(sdk.NewCoin("uqatom", math.NewInt(10000000))))
		// 	suite.NoError(err)

		// 	// enable LSM
		// 	zone, found := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
		// 	suite.True(found)
		// 	zone.LiquidityModule = true
		// 	zone.UnbondingEnabled = true
		// 	suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetZone(ctx, &zone)

		// 	validators := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetValidatorAddresses(ctx, suite.chainB.ChainID)
		// 	for _, delegation := range func(zone icstypes.Zone) []icstypes.Delegation {
		// 		out := make([]icstypes.Delegation, 0)
		// 		for _, valoper := range validators {
		// 			out = append(out, icstypes.NewDelegation(zone.DelegationAddress.Address, valoper, sdk.NewCoin(zone.BaseDenom, sdk.NewInt(3000000))))
		// 		}
		// 		return out
		// 	}(zone) {
		// 		suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetDelegation(ctx, zone.ChainId, delegation)
		// 	}

		// 	tt.malleate()

		// 	msgSrv := icskeeper.NewMsgServerImpl(suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper)
		// 	res, err := msgSrv.RequestRedemption(sdk.WrapSDKContext(suite.chainA.GetContext()), &msg)

		// 	if tt.expectErrLsm != "" {
		// 		suite.Errorf(err, tt.expectErrLsm)
		// 		suite.Nil(res)
		// 		suite.T().Logf("Error: %v", err)
		// 	} else {
		// 		suite.NoError(err)
		// 		suite.NotNil(res)
		// 	}
		// })

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

			msgSrv := icskeeper.NewMsgServerImpl(suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper)
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
		name      string
		malleate  func(suite *KeeperTestSuite) *icstypes.MsgGovCloseChannel
		expectErr error
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
			expectErr: govtypes.ErrInvalidSigner,
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
			expectErr: capabilitytypes.ErrCapabilityNotFound,
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
			expectErr: connectiontypes.ErrInvalidConnectionState,
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
			expectErr: nil,
		},
	}
	for _, tc := range testCase {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			suite.setupTestZones()

			msg := tc.malleate(suite)
			msgSrv := icskeeper.NewMsgServerImpl(suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper)
			ctx := suite.chainA.GetContext()

			_, err := msgSrv.GovCloseChannel(ctx, msg)
			if tc.expectErr != nil {
				suite.ErrorIs(tc.expectErr, err)
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

func (suite *KeeperTestSuite) TestGovReopenChannel() {
	testCase := []struct {
		name     string
		malleate func(suite *KeeperTestSuite) *icstypes.MsgGovReopenChannel
		expecErr error
	}{
		{
			name: "invalid connection id",
			malleate: func(suite *KeeperTestSuite) *icstypes.MsgGovReopenChannel {
				return &icstypes.MsgGovReopenChannel{
					ConnectionId: "",
					PortId:       "",
					Authority:    "",
				}
			},
			expecErr: fmt.Errorf("unable to obtain chain id: invalid connection id, \"%s\" not found", ""),
		},
		{
			name: "chainID / connectsionID mismatch",
			malleate: func(suite *KeeperTestSuite) *icstypes.MsgGovReopenChannel {
				return &icstypes.MsgGovReopenChannel{
					ConnectionId: suite.path.EndpointA.ConnectionID,
					PortId:       "",
					Authority:    "",
				}
			},
			expecErr: fmt.Errorf("chainID / connectionID mismatch. Connection: %s, Port: %s", "testchain2", ""),
		},
		{
			name: "existing active channel",
			malleate: func(suite *KeeperTestSuite) *icstypes.MsgGovReopenChannel {
				k := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
				ctx := suite.chainA.GetContext()
				channels := suite.GetQuicksilverApp(suite.chainA).IBCKeeper.ChannelKeeper.GetAllChannels(ctx)
				return &icstypes.MsgGovReopenChannel{
					ConnectionId: suite.path.EndpointA.ConnectionID,
					PortId:       channels[0].PortId,
					Authority:    sdk.MustBech32ifyAddressBytes(sdk.GetConfig().GetBech32AccountAddrPrefix(), k.AccountKeeper.GetModuleAddress(govtypes.ModuleName)),
				}
			},
			expecErr: errors.New("existing active channel channel-7 for portID icacontroller-testchain2.delegate on connection connection-0: active channel already set for this owner"),
		},
		{
			name: "pass",
			malleate: func(suite *KeeperTestSuite) *icstypes.MsgGovReopenChannel {
				quicksilver := suite.GetQuicksilverApp(suite.chainA)
				ctx := suite.chainA.GetContext()
				connectionID := "connection-1"
				portID := "icacontroller-testchain2.delegate"
				channelID := "channel-9"

				version := []*connectiontypes.Version{
					{Identifier: "1", Features: []string{"ORDER_ORDERED", "ORDER_UNORDERED"}},
				}
				connectionEnd := connectiontypes.ConnectionEnd{ClientId: "09-tendermint-1", State: connectiontypes.OPEN, Versions: version}
				quicksilver.IBCKeeper.ConnectionKeeper.SetConnection(ctx, connectionID, connectionEnd)

				_, f := quicksilver.IBCKeeper.ConnectionKeeper.GetConnection(ctx, connectionID)
				suite.True(f)

				channelSet := channeltypes.Channel{
					State:          channeltypes.TRYOPEN,
					Ordering:       channeltypes.NONE,
					Counterparty:   channeltypes.NewCounterparty(portID, channelID),
					ConnectionHops: []string{connectionID},
				}
				quicksilver.IBCKeeper.ChannelKeeper.SetChannel(ctx, portID, channelID, channelSet)

				quicksilver.IBCKeeper.ClientKeeper.SetClientState(ctx, connectionEnd.ClientId, &tmclienttypes.ClientState{ChainId: suite.chainB.ChainID, TrustingPeriod: time.Hour, LatestHeight: clienttypes.Height{RevisionNumber: 1, RevisionHeight: 100}})

				return &icstypes.MsgGovReopenChannel{
					ConnectionId: connectionID,
					PortId:       portID,
					Authority:    sdk.MustBech32ifyAddressBytes(sdk.GetConfig().GetBech32AccountAddrPrefix(), quicksilver.InterchainstakingKeeper.AccountKeeper.GetModuleAddress(govtypes.ModuleName)),
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
			msgSrv := icskeeper.NewMsgServerImpl(suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper)
			ctx := suite.chainA.GetContext()

			_, err := msgSrv.GovReopenChannel(ctx, msg)
			if tc.expecErr != nil {
				suite.Equal(tc.expecErr.Error(), err.Error())
				return
			}
			suite.NoError(err)

			// Check connection for port has been set
			conn, err := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.GetConnectionForPort(ctx, msg.PortId)
			suite.NoError(err)
			suite.Equal(conn, msg.ConnectionId)
		})
	}
}

func (suite *KeeperTestSuite) TestSetLsmCaps() {
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
			"invalid zone",
			func(s *KeeperTestSuite) *icstypes.MsgGovSetLsmCaps {
				return &icstypes.MsgGovSetLsmCaps{
					ChainId: "unknownzone-1",
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
			"non lsm zone",
			func(s *KeeperTestSuite) *icstypes.MsgGovSetLsmCaps {
				zone, _ := s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.GetZone(s.chainA.GetContext(), s.chainB.ChainID)
				zone.LiquidityModule = false
				s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper.SetZone(s.chainA.GetContext(), &zone)

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
					Authority: "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
				}
			},
			false,
		},
	}

	for _, tt := range tests {
		tt := tt

		suite.Run(tt.name, func() {
			suite.SetupTest()
			suite.setupTestZones()

			msg := tt.malleate(suite)

			msgSrv := icskeeper.NewMsgServerImpl(suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper)
			res, err := msgSrv.GovSetLsmCaps(sdk.WrapSDKContext(suite.chainA.GetContext()), msg)
			if tt.expectErr {
				suite.Error(err)
				suite.Nil(res)
			} else {
				suite.NoError(err)
				suite.NotNil(res)
			}

			qapp := suite.GetQuicksilverApp(suite.chainA)
			icsKeeper := qapp.InterchainstakingKeeper
			zone, found := icsKeeper.GetZone(suite.chainA.GetContext(), suite.chainB.ChainID)
			suite.True(found)

			caps, found := icsKeeper.GetLsmCaps(suite.chainA.GetContext(), zone.ChainId)
			if tt.expectErr {
				suite.False(found)
				suite.Nil(caps)
			} else {
				suite.True(found)
				suite.Equal(caps, msg.Caps)

			}
		})
	}
}

func (suite *KeeperTestSuite) TestMsgCancelRedemeption() {
	hash := randomutils.GenerateRandomHashAsHex(32)
	tests := []struct {
		name      string
		malleate  func(s *KeeperTestSuite) *icstypes.MsgCancelRedemption
		expectErr string
	}{
		{
			"no zone exists",
			func(s *KeeperTestSuite) *icstypes.MsgCancelRedemption {
				return &icstypes.MsgCancelRedemption{
					ChainId:     "bob",
					Hash:        hash,
					FromAddress: addressutils.GenerateAddressForTestWithPrefix("quick"),
				}
			},
			fmt.Sprintf("no queued record with hash \"%s\" found", hash),
		},
		{
			"no hash exists",
			func(s *KeeperTestSuite) *icstypes.MsgCancelRedemption {
				return &icstypes.MsgCancelRedemption{
					ChainId:     s.chainB.ChainID,
					Hash:        hash,
					FromAddress: addressutils.GenerateAddressForTestWithPrefix("quick"),
				}
			},
			fmt.Sprintf("no queued record with hash \"%s\" found", hash),
		},
		{
			"hash exists but in unbond status, no errors",
			func(s *KeeperTestSuite) *icstypes.MsgCancelRedemption {
				ctx := s.chainA.GetContext()
				k := s.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
				address := addressutils.GenerateAddressForTestWithPrefix("quick")
				_ = k.SetWithdrawalRecord(ctx, icstypes.WithdrawalRecord{
					ChainId:        s.chainB.ChainID,
					Delegator:      address,
					BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(500)),
					Status:         icstypes.WithdrawStatusUnbond,
					CompletionTime: ctx.BlockHeader().Time.Add(time.Hour * 72),
					Txhash:         hash,
				})

				return &icstypes.MsgCancelRedemption{
					ChainId:     s.chainB.ChainID,
					Hash:        hash,
					FromAddress: address,
				}
			},
			fmt.Sprintf("cannot cancel unbond \"%s\" with no errors", hash),
		},
		{
			"hash exists in queued status, with errors",
			func(s *KeeperTestSuite) *icstypes.MsgCancelRedemption {
				ctx := s.chainA.GetContext()
				k := s.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
				address := addressutils.GenerateAddressForTestWithPrefix("quick")
				_ = k.SetWithdrawalRecord(ctx, icstypes.WithdrawalRecord{
					ChainId:        s.chainB.ChainID,
					Delegator:      address,
					BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(500)),
					Status:         icstypes.WithdrawStatusUnbond,
					CompletionTime: ctx.BlockHeader().Time.Add(time.Hour * 72),
					Txhash:         hash,
					SendErrors:     1,
				})

				suite.NoError(k.BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqatom", math.NewInt(500)))))
				suite.NoError(k.BankKeeper.SendCoinsFromModuleToModule(ctx, icstypes.ModuleName, icstypes.EscrowModuleAccount, sdk.NewCoins(sdk.NewCoin("uqatom", math.NewInt(500)))))

				return &icstypes.MsgCancelRedemption{
					ChainId:     s.chainB.ChainID,
					Hash:        hash,
					FromAddress: address,
				}
			},
			"",
		},
		{
			"hash exists in correct status but different user",
			func(s *KeeperTestSuite) *icstypes.MsgCancelRedemption {
				ctx := s.chainA.GetContext()
				k := s.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
				// generate two addresses, one for the withdrawal record we are looking up; one for the tx.
				withdrawalAddress := addressutils.GenerateAddressForTestWithPrefix("quick")
				address := addressutils.GenerateAddressForTestWithPrefix("quick")
				_ = k.SetWithdrawalRecord(ctx, icstypes.WithdrawalRecord{
					ChainId:        s.chainB.ChainID,
					Delegator:      withdrawalAddress,
					BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(500)),
					Status:         icstypes.WithdrawStatusQueued,
					CompletionTime: ctx.BlockHeader().Time.Add(time.Hour * 72),
					Txhash:         hash,
				})

				return &icstypes.MsgCancelRedemption{
					ChainId:     s.chainB.ChainID,
					Hash:        hash,
					FromAddress: address,
				}
			},
			fmt.Sprintf("incorrect user for record with hash \"%s\"", hash),
		},
		{
			"valid",
			func(s *KeeperTestSuite) *icstypes.MsgCancelRedemption {
				ctx := s.chainA.GetContext()
				k := s.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
				address := addressutils.GenerateAddressForTestWithPrefix("quick")
				_ = k.SetWithdrawalRecord(ctx, icstypes.WithdrawalRecord{
					ChainId:        s.chainB.ChainID,
					Delegator:      address,
					BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(500)),
					Status:         icstypes.WithdrawStatusQueued,
					CompletionTime: ctx.BlockHeader().Time.Add(time.Hour * 72),
					Txhash:         hash,
				})

				suite.NoError(k.BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqatom", math.NewInt(500)))))
				suite.NoError(k.BankKeeper.SendCoinsFromModuleToModule(ctx, icstypes.ModuleName, icstypes.EscrowModuleAccount, sdk.NewCoins(sdk.NewCoin("uqatom", math.NewInt(500)))))

				return &icstypes.MsgCancelRedemption{
					ChainId:     s.chainB.ChainID,
					Hash:        hash,
					FromAddress: address,
				}
			},
			"",
		},
		{
			"valid - governance",
			func(s *KeeperTestSuite) *icstypes.MsgCancelRedemption {
				ctx := s.chainA.GetContext()
				k := s.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
				address := addressutils.GenerateAddressForTestWithPrefix("quick")
				_ = k.SetWithdrawalRecord(ctx, icstypes.WithdrawalRecord{
					ChainId:        s.chainB.ChainID,
					Delegator:      address,
					BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(500)),
					Status:         icstypes.WithdrawStatusQueued,
					CompletionTime: ctx.BlockHeader().Time.Add(time.Hour * 72),
					Txhash:         hash,
				})

				suite.NoError(k.BankKeeper.MintCoins(ctx, icstypes.ModuleName, sdk.NewCoins(sdk.NewCoin("uqatom", math.NewInt(500)))))
				suite.NoError(k.BankKeeper.SendCoinsFromModuleToModule(ctx, icstypes.ModuleName, icstypes.EscrowModuleAccount, sdk.NewCoins(sdk.NewCoin("uqatom", math.NewInt(500)))))

				return &icstypes.MsgCancelRedemption{
					ChainId:     s.chainB.ChainID,
					Hash:        hash,
					FromAddress: k.GetGovAuthority(ctx),
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

			msg := tt.malleate(suite)

			msgSrv := icskeeper.NewMsgServerImpl(suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper)
			res, err := msgSrv.CancelRedemption(sdk.WrapSDKContext(suite.chainA.GetContext()), msg)
			if len(tt.expectErr) != 0 {
				suite.ErrorContains(err, tt.expectErr)
				suite.Nil(res)
			} else {
				suite.NoError(err)
				suite.NotNil(res)
			}

			qapp := suite.GetQuicksilverApp(suite.chainA)
			icsKeeper := qapp.InterchainstakingKeeper
			_, found := icsKeeper.GetZone(suite.chainA.GetContext(), suite.chainB.ChainID)
			suite.True(found)
		})
	}
}

func (suite *KeeperTestSuite) TestMsgRequeueRedemeption() {
	hash := randomutils.GenerateRandomHashAsHex(32)
	tests := []struct {
		name      string
		malleate  func(s *KeeperTestSuite) *icstypes.MsgRequeueRedemption
		expectErr string
	}{
		{
			"no zone exists",
			func(s *KeeperTestSuite) *icstypes.MsgRequeueRedemption {
				return &icstypes.MsgRequeueRedemption{
					ChainId:     "bob",
					Hash:        hash,
					FromAddress: addressutils.GenerateAddressForTestWithPrefix("quick"),
				}
			},
			fmt.Sprintf("no unbonding record with hash \"%s\" found", hash),
		},
		{
			"no hash exists",
			func(s *KeeperTestSuite) *icstypes.MsgRequeueRedemption {
				return &icstypes.MsgRequeueRedemption{
					ChainId:     s.chainB.ChainID,
					Hash:        hash,
					FromAddress: addressutils.GenerateAddressForTestWithPrefix("quick"),
				}
			},
			fmt.Sprintf("no unbonding record with hash \"%s\" found", hash),
		},
		{
			"hash exists but in unbond status, no errors",
			func(s *KeeperTestSuite) *icstypes.MsgRequeueRedemption {
				ctx := s.chainA.GetContext()
				k := s.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
				address := addressutils.GenerateAddressForTestWithPrefix("quick")
				_ = k.SetWithdrawalRecord(ctx, icstypes.WithdrawalRecord{
					ChainId:        s.chainB.ChainID,
					Delegator:      address,
					BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(500)),
					Status:         icstypes.WithdrawStatusUnbond,
					CompletionTime: ctx.BlockHeader().Time.Add(time.Hour * 72),
					Txhash:         hash,
				})

				return &icstypes.MsgRequeueRedemption{
					ChainId:     s.chainB.ChainID,
					Hash:        hash,
					FromAddress: address,
				}
			},
			fmt.Sprintf("cannot requeue unbond \"%s\" with no errors", hash),
		},
		{
			"hash exists in queued status, with errors",
			func(s *KeeperTestSuite) *icstypes.MsgRequeueRedemption {
				ctx := s.chainA.GetContext()
				k := s.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
				address := addressutils.GenerateAddressForTestWithPrefix("quick")
				_ = k.SetWithdrawalRecord(ctx, icstypes.WithdrawalRecord{
					ChainId:        s.chainB.ChainID,
					Delegator:      address,
					BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(500)),
					Status:         icstypes.WithdrawStatusUnbond,
					CompletionTime: ctx.BlockHeader().Time.Add(time.Hour * 72),
					Txhash:         hash,
					SendErrors:     1,
				})

				return &icstypes.MsgRequeueRedemption{
					ChainId:     s.chainB.ChainID,
					Hash:        hash,
					FromAddress: address,
				}
			},
			"",
		},
		{
			"hash exists in correct status but different user",
			func(s *KeeperTestSuite) *icstypes.MsgRequeueRedemption {
				ctx := s.chainA.GetContext()
				k := s.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
				// generate two addresses, one for the withdrawal record we are looking up; one for the tx.
				withdrawalAddress := addressutils.GenerateAddressForTestWithPrefix("quick")
				address := addressutils.GenerateAddressForTestWithPrefix("quick")
				_ = k.SetWithdrawalRecord(ctx, icstypes.WithdrawalRecord{
					ChainId:        s.chainB.ChainID,
					Delegator:      withdrawalAddress,
					BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(500)),
					Status:         icstypes.WithdrawStatusUnbond,
					CompletionTime: ctx.BlockHeader().Time.Add(time.Hour * 72),
					Txhash:         hash,
					SendErrors:     1,
				})

				return &icstypes.MsgRequeueRedemption{
					ChainId:     s.chainB.ChainID,
					Hash:        hash,
					FromAddress: address,
				}
			},
			fmt.Sprintf("incorrect user for record with hash \"%s\"", hash),
		},
		{
			"valid - governance",
			func(s *KeeperTestSuite) *icstypes.MsgRequeueRedemption {
				ctx := s.chainA.GetContext()
				k := s.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
				address := addressutils.GenerateAddressForTestWithPrefix("quick")
				_ = k.SetWithdrawalRecord(ctx, icstypes.WithdrawalRecord{
					ChainId:        s.chainB.ChainID,
					Delegator:      address,
					BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(500)),
					Status:         icstypes.WithdrawStatusUnbond,
					CompletionTime: ctx.BlockHeader().Time.Add(time.Hour * 72),
					Txhash:         hash,
					SendErrors:     1,
				})

				return &icstypes.MsgRequeueRedemption{
					ChainId:     s.chainB.ChainID,
					Hash:        hash,
					FromAddress: k.GetGovAuthority(ctx),
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

			msg := tt.malleate(suite)

			msgSrv := icskeeper.NewMsgServerImpl(suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper)
			res, err := msgSrv.RequeueRedemption(sdk.WrapSDKContext(suite.chainA.GetContext()), msg)
			if len(tt.expectErr) != 0 {
				suite.ErrorContains(err, tt.expectErr)
				suite.Nil(res)
			} else {
				suite.NoError(err)
				suite.NotNil(res)
			}

			qapp := suite.GetQuicksilverApp(suite.chainA)
			icsKeeper := qapp.InterchainstakingKeeper
			_, found := icsKeeper.GetZone(suite.chainA.GetContext(), suite.chainB.ChainID)
			suite.True(found)
		})
	}
}

func (suite *KeeperTestSuite) TestMsgUpdateRedemption() {
	hash := randomutils.GenerateRandomHashAsHex(32)
	tests := []struct {
		name      string
		malleate  func(s *KeeperTestSuite) *icstypes.MsgUpdateRedemption
		expectErr string
		assert    func(s *KeeperTestSuite) bool
	}{
		{
			"no zone exists",
			func(s *KeeperTestSuite) *icstypes.MsgUpdateRedemption {
				k := s.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
				ctx := s.chainA.GetContext()
				return &icstypes.MsgUpdateRedemption{
					ChainId:     "bob",
					Hash:        hash,
					NewStatus:   icstypes.WithdrawStatusUnbond,
					FromAddress: k.GetGovAuthority(ctx),
				}
			},
			fmt.Sprintf("no unbonding record with hash \"%s\" found", hash),
			nil,
		},
		{
			"no hash exists",
			func(s *KeeperTestSuite) *icstypes.MsgUpdateRedemption {
				k := s.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
				ctx := s.chainA.GetContext()
				return &icstypes.MsgUpdateRedemption{
					ChainId:     s.chainB.ChainID,
					Hash:        hash,
					NewStatus:   icstypes.WithdrawStatusUnbond,
					FromAddress: k.GetGovAuthority(ctx),
				}
			},
			fmt.Sprintf("no unbonding record with hash \"%s\" found", hash),
			nil,
		},
		{
			"invalid - cannot transition to send",
			func(s *KeeperTestSuite) *icstypes.MsgUpdateRedemption {
				ctx := s.chainA.GetContext()
				k := s.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
				address := addressutils.GenerateAddressForTestWithPrefix("quick")
				_ = k.SetWithdrawalRecord(ctx, icstypes.WithdrawalRecord{
					ChainId:        s.chainB.ChainID,
					Delegator:      address,
					BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(500)),
					Status:         icstypes.WithdrawStatusQueued,
					CompletionTime: ctx.BlockHeader().Time.Add(time.Hour * 72),
					Txhash:         hash,
					SendErrors:     1,
				})

				return &icstypes.MsgUpdateRedemption{
					ChainId:     s.chainB.ChainID,
					Hash:        hash,
					NewStatus:   icstypes.WithdrawStatusSend,
					FromAddress: k.GetGovAuthority(ctx),
				}
			},
			"new status WithdrawStatusSend not supported",
			nil,
		},
		{
			"invalid - cannot transition to tokenize",
			func(s *KeeperTestSuite) *icstypes.MsgUpdateRedemption {
				ctx := s.chainA.GetContext()
				k := s.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
				address := addressutils.GenerateAddressForTestWithPrefix("quick")
				_ = k.SetWithdrawalRecord(ctx, icstypes.WithdrawalRecord{
					ChainId:        s.chainB.ChainID,
					Delegator:      address,
					BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(500)),
					Status:         icstypes.WithdrawStatusQueued,
					CompletionTime: ctx.BlockHeader().Time.Add(time.Hour * 72),
					Txhash:         hash,
					SendErrors:     1,
				})

				return &icstypes.MsgUpdateRedemption{
					ChainId:     s.chainB.ChainID,
					Hash:        hash,
					NewStatus:   icstypes.WithdrawStatusTokenize,
					FromAddress: k.GetGovAuthority(ctx),
				}
			},
			"new status WithdrawStatusTokenize not supported",
			nil,
		},
		{
			"invalid - cannot transition to send",
			func(s *KeeperTestSuite) *icstypes.MsgUpdateRedemption {
				ctx := s.chainA.GetContext()
				k := s.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
				address := addressutils.GenerateAddressForTestWithPrefix("quick")
				_ = k.SetWithdrawalRecord(ctx, icstypes.WithdrawalRecord{
					ChainId:        s.chainB.ChainID,
					Delegator:      address,
					BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(500)),
					Status:         icstypes.WithdrawStatusQueued,
					CompletionTime: ctx.BlockHeader().Time.Add(time.Hour * 72),
					Txhash:         hash,
					SendErrors:     1,
				})

				return &icstypes.MsgUpdateRedemption{
					ChainId:     s.chainB.ChainID,
					Hash:        hash,
					NewStatus:   0,
					FromAddress: k.GetGovAuthority(ctx),
				}
			},
			"new status not provided or invalid",
			nil,
		},
		{
			"valid - target complete",
			func(s *KeeperTestSuite) *icstypes.MsgUpdateRedemption {
				ctx := s.chainA.GetContext()
				k := s.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
				address := addressutils.GenerateAddressForTestWithPrefix("quick")
				_ = k.SetWithdrawalRecord(ctx, icstypes.WithdrawalRecord{
					ChainId:        s.chainB.ChainID,
					Delegator:      address,
					BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(500)),
					Status:         icstypes.WithdrawStatusUnbond,
					CompletionTime: ctx.BlockHeader().Time.Add(time.Hour * 72),
					Txhash:         hash,
					SendErrors:     1,
				})

				return &icstypes.MsgUpdateRedemption{
					ChainId:     s.chainB.ChainID,
					Hash:        hash,
					NewStatus:   icstypes.WithdrawStatusCompleted,
					FromAddress: k.GetGovAuthority(ctx),
				}
			},
			"",
			func(s *KeeperTestSuite) bool {
				qapp := suite.GetQuicksilverApp(suite.chainA)
				icsKeeper := qapp.InterchainstakingKeeper
				_, found := icsKeeper.GetWithdrawalRecord(suite.chainA.GetContext(), s.chainB.ChainID, hash, icstypes.WithdrawStatusCompleted)
				return found
			},
		},
		{
			"valid - target unbonding",
			func(s *KeeperTestSuite) *icstypes.MsgUpdateRedemption {
				ctx := s.chainA.GetContext()
				k := s.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
				address := addressutils.GenerateAddressForTestWithPrefix("quick")
				_ = k.SetWithdrawalRecord(ctx, icstypes.WithdrawalRecord{
					ChainId:        s.chainB.ChainID,
					Delegator:      address,
					BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(500)),
					Status:         icstypes.WithdrawStatusSend,
					CompletionTime: ctx.BlockHeader().Time.Add(time.Hour * 72),
					Txhash:         hash,
					SendErrors:     1,
				})

				return &icstypes.MsgUpdateRedemption{
					ChainId:     s.chainB.ChainID,
					Hash:        hash,
					NewStatus:   icstypes.WithdrawStatusUnbond,
					FromAddress: k.GetGovAuthority(ctx),
				}
			},
			"",
			func(s *KeeperTestSuite) bool {
				qapp := suite.GetQuicksilverApp(suite.chainA)
				icsKeeper := qapp.InterchainstakingKeeper
				_, found := icsKeeper.GetWithdrawalRecord(suite.chainA.GetContext(), s.chainB.ChainID, hash, icstypes.WithdrawStatusUnbond)
				return found
			},
		},
		{
			"valid - target queued",
			func(s *KeeperTestSuite) *icstypes.MsgUpdateRedemption {
				ctx := s.chainA.GetContext()
				k := s.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
				address := addressutils.GenerateAddressForTestWithPrefix("quick")
				_ = k.SetWithdrawalRecord(ctx, icstypes.WithdrawalRecord{
					ChainId:        s.chainB.ChainID,
					Delegator:      address,
					BurnAmount:     sdk.NewCoin("uqatom", math.NewInt(500)),
					Status:         icstypes.WithdrawStatusUnbond,
					CompletionTime: ctx.BlockHeader().Time.Add(time.Hour * 72),
					Txhash:         hash,
					SendErrors:     1,
				})

				return &icstypes.MsgUpdateRedemption{
					ChainId:     s.chainB.ChainID,
					Hash:        hash,
					NewStatus:   icstypes.WithdrawStatusQueued,
					FromAddress: k.GetGovAuthority(ctx),
				}
			},
			"",
			func(s *KeeperTestSuite) bool {
				qapp := suite.GetQuicksilverApp(suite.chainA)
				icsKeeper := qapp.InterchainstakingKeeper
				record, found := icsKeeper.GetWithdrawalRecord(suite.chainA.GetContext(), s.chainB.ChainID, hash, icstypes.WithdrawStatusQueued)

				return found && record.Distribution == nil && record.Amount == nil && record.SendErrors == 0
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		suite.Run(tt.name, func() {
			suite.SetupTest()
			suite.setupTestZones()

			msg := tt.malleate(suite)

			msgSrv := icskeeper.NewMsgServerImpl(suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper)
			res, err := msgSrv.UpdateRedemption(sdk.WrapSDKContext(suite.chainA.GetContext()), msg)
			if len(tt.expectErr) != 0 {
				suite.ErrorContains(err, tt.expectErr)
				suite.Nil(res)
			} else {
				suite.NoError(err)
				suite.NotNil(res)
				suite.True(tt.assert(suite))
			}

			qapp := suite.GetQuicksilverApp(suite.chainA)
			icsKeeper := qapp.InterchainstakingKeeper
			_, found := icsKeeper.GetZone(suite.chainA.GetContext(), suite.chainB.ChainID)
			suite.True(found)
		})
	}
}

func (suite *KeeperTestSuite) TestMsgGovAddValidatorToDenyList() {
	dummyChainID := "dummychain"
	testValAddr, _ := sdk.ValAddressFromHex(suite.chainB.Vals.Validators[0].Address.String())
	bech32ValoperAddr := addressutils.MustEncodeAddressToBech32(sdk.GetConfig().GetBech32ValidatorAddrPrefix(), testValAddr)
	govModuleAddr := sdk.MustBech32ifyAddressBytes(sdk.GetConfig().GetBech32AccountAddrPrefix(), suite.GetQuicksilverApp(suite.chainA).AccountKeeper.GetModuleAddress(govtypes.ModuleName))

	tests := []struct {
		name      string
		malleate  func(s *KeeperTestSuite) *icstypes.MsgGovAddValidatorDenyList
		expectErr string
	}{
		{
			"invalid authority",
			func(s *KeeperTestSuite) *icstypes.MsgGovAddValidatorDenyList {
				return &icstypes.MsgGovAddValidatorDenyList{
					ChainId:         dummyChainID,
					OperatorAddress: bech32ValoperAddr,
					Authority:       testAddress,
				}
			},
			"expected gov account as only signer for proposal message",
		},
		{
			"invalid chain-id",
			func(s *KeeperTestSuite) *icstypes.MsgGovAddValidatorDenyList {
				return &icstypes.MsgGovAddValidatorDenyList{
					ChainId:         dummyChainID,
					OperatorAddress: bech32ValoperAddr,
					Authority:       govModuleAddr,
				}
			},
			fmt.Sprintf("no zone found for: %s", dummyChainID),
		},
		{
			"invalid operator address",
			func(s *KeeperTestSuite) *icstypes.MsgGovAddValidatorDenyList {
				return &icstypes.MsgGovAddValidatorDenyList{
					ChainId:         s.chainB.ChainID,
					OperatorAddress: "invalid",
					Authority:       govModuleAddr,
				}
			},
			"decoding bech32 failed",
		},
		{
			"valid",
			func(s *KeeperTestSuite) *icstypes.MsgGovAddValidatorDenyList {
				return &icstypes.MsgGovAddValidatorDenyList{
					ChainId:         s.chainB.ChainID,
					OperatorAddress: bech32ValoperAddr,
					Authority:       govModuleAddr,
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

			msg := tt.malleate(suite)

			msgSrv := icskeeper.NewMsgServerImpl(suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper)
			res, err := msgSrv.GovAddValidatorDenyList(sdk.WrapSDKContext(suite.chainA.GetContext()), msg)
			if len(tt.expectErr) != 0 {
				suite.ErrorContains(err, tt.expectErr)
				suite.Nil(res)
			} else {
				suite.NoError(err)
				suite.NotNil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMsgGovRemoveValidatorToDenyList() {
	dummyChainID := "dummychain"
	testValAddr, _ := sdk.ValAddressFromHex(suite.chainB.Vals.Validators[0].Address.String())
	bech32ValoperAddr := addressutils.MustEncodeAddressToBech32(sdk.GetConfig().GetBech32ValidatorAddrPrefix(), testValAddr)
	govModuleAddr := sdk.MustBech32ifyAddressBytes(sdk.GetConfig().GetBech32AccountAddrPrefix(), suite.GetQuicksilverApp(suite.chainA).AccountKeeper.GetModuleAddress(govtypes.ModuleName))
	tests := []struct {
		name      string
		malleate  func(s *KeeperTestSuite) *icstypes.MsgGovRemoveValidatorDenyList
		expectErr string
	}{
		{
			"invalid authority",
			func(s *KeeperTestSuite) *icstypes.MsgGovRemoveValidatorDenyList {
				return &icstypes.MsgGovRemoveValidatorDenyList{
					ChainId:         dummyChainID,
					OperatorAddress: bech32ValoperAddr,
					Authority:       testAddress,
				}
			},
			"expected gov account as only signer for proposal message",
		},
		{
			"invalid chain-id",
			func(s *KeeperTestSuite) *icstypes.MsgGovRemoveValidatorDenyList {
				return &icstypes.MsgGovRemoveValidatorDenyList{
					ChainId:         dummyChainID,
					OperatorAddress: bech32ValoperAddr,
					Authority:       govModuleAddr,
				}
			},
			fmt.Sprintf("no zone found for: %s", dummyChainID),
		},
		{
			"invalid operator address",
			func(s *KeeperTestSuite) *icstypes.MsgGovRemoveValidatorDenyList {
				return &icstypes.MsgGovRemoveValidatorDenyList{
					ChainId:         s.chainB.ChainID,
					OperatorAddress: "invalid",
					Authority:       govModuleAddr,
				}
			},
			"decoding bech32 failed",
		},
		{
			"valid msg, but not in deny list",
			func(s *KeeperTestSuite) *icstypes.MsgGovRemoveValidatorDenyList {
				return &icstypes.MsgGovRemoveValidatorDenyList{
					ChainId:         s.chainB.ChainID,
					OperatorAddress: bech32ValoperAddr,
					Authority:       govModuleAddr,
				}
			},
			fmt.Sprintf("validator %s not found in deny list", bech32ValoperAddr),
		},
		{
			"valid msg, validator in deny list",
			func(s *KeeperTestSuite) *icstypes.MsgGovRemoveValidatorDenyList {
				k := s.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
				err := k.SetZoneValidatorToDenyList(s.chainA.GetContext(), s.chainB.ChainID, testValAddr)
				suite.NoError(err)
				return &icstypes.MsgGovRemoveValidatorDenyList{
					ChainId:         s.chainB.ChainID,
					OperatorAddress: bech32ValoperAddr,
					Authority:       govModuleAddr,
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

			msg := tt.malleate(suite)

			msgSrv := icskeeper.NewMsgServerImpl(suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper)
			res, err := msgSrv.GovRemoveValidatorDenyList(sdk.WrapSDKContext(suite.chainA.GetContext()), msg)
			if len(tt.expectErr) != 0 {
				suite.ErrorContains(err, tt.expectErr)
				suite.Nil(res)
			} else {
				suite.NoError(err)
				suite.NotNil(res)
			}
		})
	}
}

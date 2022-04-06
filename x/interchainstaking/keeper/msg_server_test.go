package keeper_test

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	ibctesting "github.com/cosmos/ibc-go/v3/testing"
	icskeeper "github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func (s *KeeperTestSuite) TestRegisterZone() {
	var (
		path *ibctesting.Path
		msg  icstypes.MsgRegisterZone
	)

	tests := []struct {
		name      string
		malleate  func()
		expectErr bool
	}{
		{
			"invalid connection",
			func() {
				msg = icstypes.MsgRegisterZone{
					Identifier:   "cosmosquitto",
					ConnectionId: "unknown",
					LocalDenom:   "uqatom",
					BaseDenom:    "uatom",
					FromAddress:  TestOwnerAddress,
				}
			},
			true,
		},
		{
			"duplicate",
			func() {
				s.SetupRegisteredZones()
				msg = icstypes.MsgRegisterZone{
					Identifier:   "cosmos",
					ConnectionId: path.EndpointA.ConnectionID,
					LocalDenom:   "uqatom",
					BaseDenom:    "uatom",
					FromAddress:  TestOwnerAddress,
				}
			},
			true,
		},
		{
			"valid",
			func() {
				msg = icstypes.MsgRegisterZone{
					Identifier:   "cosmos",
					ConnectionId: path.EndpointA.ConnectionID,
					LocalDenom:   "uqatom",
					BaseDenom:    "uatom",
					FromAddress:  TestOwnerAddress,
				}
			},
			false,
		},
	}

	for _, tt := range tests {
		tt := tt

		s.Run(tt.name, func() {
			s.SetupTest()

			path = NewQuicksilverPath(s.chainA, s.chainB)
			s.coordinator.SetupConnections(path)

			tt.malleate()

			msgSrv := icskeeper.NewMsgServerImpl(s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper)
			res, err := msgSrv.RegisterZone(sdktypes.WrapSDKContext(s.chainA.GetContext()), &msg)

			if tt.expectErr {
				s.Require().Error(err)
				s.Require().Nil(res)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(res)
			}
		})
	}
}

func (s *KeeperTestSuite) TestRequestRedemption() {
	var (
		path *ibctesting.Path
		msg  icstypes.MsgRequestRedemption
	)

	tests := []struct {
		name      string
		malleate  func()
		expectErr bool
	}{
		{
			"valid",
			func() {
				msg = icstypes.MsgRequestRedemption{
					Coin:               "uatom",
					DestinationAddress: TestOwnerAddress,
					FromAddress:        TestOwnerAddress,
				}
			},
			false,
		},
	}

	for _, tt := range tests {
		tt := tt

		s.Run(tt.name, func() {
			s.SetupTest()

			path = NewQuicksilverPath(s.chainA, s.chainB)
			s.coordinator.SetupConnections(path)

			tt.malleate()

			msgSrv := icskeeper.NewMsgServerImpl(s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper)
			res, err := msgSrv.RequestRedemption(sdktypes.WrapSDKContext(s.chainA.GetContext()), &msg)

			if tt.expectErr {
				s.Require().Error(err)
				s.Require().Nil(res)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(res)
			}
		})
	}
}

func (s *KeeperTestSuite) TestSignalIntent() {
	var (
		msg icstypes.MsgSignalIntent
	)

	tests := []struct {
		name      string
		malleate  func()
		expectErr bool
	}{
		{
			"valid",
			func() {
				valAddress, err := sdktypes.ValAddressFromHex(s.chainB.Vals.Validators[0].Address.String())
				s.Require().NoError(err)
				msg = icstypes.MsgSignalIntent{
					ChainId: s.chainB.ChainID,
					Intents: []*icstypes.ValidatorIntent{
						{
							ValoperAddress: valAddress.String(),
							Weight:         sdktypes.MustNewDecFromStr("0.3"),
						},
					},
					FromAddress: TestOwnerAddress,
				}
			},
			false,
		},
	}

	for _, tt := range tests {
		tt := tt

		s.Run(tt.name, func() {
			s.SetupTest()
			s.SetupRegisteredZones()

			tt.malleate()

			msgSrv := icskeeper.NewMsgServerImpl(s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper)
			res, err := msgSrv.SignalIntent(sdktypes.WrapSDKContext(s.chainA.GetContext()), &msg)

			if tt.expectErr {
				s.Require().Error(err)
				s.Require().Nil(res)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(res)
			}
		})
	}
}

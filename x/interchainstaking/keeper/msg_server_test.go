package keeper_test

import (
	"context"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctesting "github.com/cosmos/ibc-go/v3/testing"
	icqkeeper "github.com/ingenuity-build/quicksilver/x/interchainquery/keeper"
	icqtypes "github.com/ingenuity-build/quicksilver/x/interchainquery/types"
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
					ChainId:      s.chainB.ChainID,
					LocalDenom:   "uqatom",
					BaseDenom:    "uatom",
					FromAddress:  TestOwnerAddress,
				}
			},
			true,
		},
		// This test does not fail as RegisterZone does no validation of the ChainId...
		/*{
			"invalid chain",
			func() {
				msg = icstypes.MsgRegisterZone{
					Identifier:   "cosmosquitto",
					ConnectionId: path.EndpointA.ConnectionID,
					ChainId:      "boguschain",
					LocalDenom:   "uqatom",
					BaseDenom:    "uatom",
					FromAddress:  TestOwnerAddress,
				}
			},
			true,
		},*/
		{
			"valid",
			func() {
				msg = icstypes.MsgRegisterZone{
					Identifier:   "cosmos",
					ConnectionId: path.EndpointA.ConnectionID,
					ChainId:      s.chainB.ChainID,
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
		path *ibctesting.Path
		msg  icstypes.MsgSignalIntent
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

			zonemsg := icstypes.MsgRegisterZone{
				Identifier:   "cosmos",
				ConnectionId: path.EndpointA.ConnectionID,
				ChainId:      s.chainB.ChainID,
				LocalDenom:   "uqatom",
				BaseDenom:    "uatom",
				FromAddress:  TestOwnerAddress,
			}

			tt.malleate()

			msgSrv := icskeeper.NewMsgServerImpl(s.GetQuicksilverApp(s.chainA).InterchainstakingKeeper)
			ctx := s.chainA.GetContext()
			ctx = ctx.WithContext(context.WithValue(ctx.Context(), "TEST", "TEST"))
			_, err := msgSrv.RegisterZone(sdktypes.WrapSDKContext(ctx), &zonemsg)
			s.Require().NoError(err)

			qvr := stakingtypes.QueryValidatorsResponse{
				Validators: s.GetQuicksilverApp(s.chainB).StakingKeeper.GetBondedValidatorsByPower(s.chainB.GetContext()),
			}

			icqmsgSrv := icqkeeper.NewMsgServerImpl(s.GetQuicksilverApp(s.chainA).InterchainQueryKeeper)
			qmsg := icqtypes.MsgSubmitQueryResponse{
				// target or source chain_id?
				ChainId: s.chainB.ChainID,
				QueryId: icqkeeper.GenerateQueryHash(
					path.EndpointA.ConnectionID,
					s.chainB.ChainID,
					"cosmos.staking.v1beta1.Query/Validators",
					map[string]string{"status": stakingtypes.BondStatusBonded},
				),
				Result:      s.GetQuicksilverApp(s.chainB).AppCodec().MustMarshalJSON(&qvr),
				Height:      s.chainB.CurrentHeader.Height,
				FromAddress: TestOwnerAddress,
			}
			_, err = icqmsgSrv.SubmitQueryResponse(sdktypes.WrapSDKContext(ctx), &qmsg)
			s.Require().NoError(err)

			s.coordinator.CommitNBlocks(s.chainA, 25)
			s.coordinator.CommitNBlocks(s.chainB, 25)

			//fmt.Printf("qmsg: %v\n", qmsg)

			res, err := msgSrv.SignalIntent(sdktypes.WrapSDKContext(ctx), &msg)

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

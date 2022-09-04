package keeper_test

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"

	icskeeper "github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

func (s *KeeperTestSuite) TestRequestRedemption() {
	var msg icstypes.MsgRequestRedemption

	tests := []struct {
		name      string
		malleate  func()
		expectErr bool
	}{
		// TODO: setup test cases for RequestRedemption
		/*{
			"valid",
			func() {
				msg = icstypes.MsgRequestRedemption{
					Value:              sdktypes.NewCoin("uatom", sdktypes.NewInt(10000000)),
					DestinationAddress: TestOwnerAddress,
					FromAddress:        TestOwnerAddress,
				}
			},
			false,
		},*/
	}

	for _, tt := range tests {
		tt := tt

		s.Run(tt.name, func() {
			s.SetupTest()
			s.SetupZones()

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
	var msg icstypes.MsgSignalIntent

	tests := []struct {
		name      string
		malleate  func()
		expectErr bool
	}{
		/*{
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
		},*/
	}

	for _, tt := range tests {
		tt := tt

		s.Run(tt.name, func() {
			s.SetupTest()
			s.SetupZones()

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

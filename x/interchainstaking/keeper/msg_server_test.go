package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
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
				val1, err := sdktypes.ValAddressFromHex(s.chainB.Vals.Validators[0].Address.String())
				s.Require().NoError(err)

				return &icstypes.MsgSignalIntent{
					ChainId: s.chainB.ChainID,
					Intents: []*icstypes.ValidatorIntent{
						{
							ValoperAddress: val1.String(),
							Weight:         sdktypes.MustNewDecFromStr("0.3"),
						},
					},
					FromAddress: TestOwnerAddress,
				}
			},
			[]sdk.Dec{},
			true,
			false,
		},
		{
			"invalid - weight sum > 1",
			func(s *KeeperTestSuite) *icstypes.MsgSignalIntent {
				val1, err := sdktypes.ValAddressFromHex(s.chainB.Vals.Validators[0].Address.String())
				s.Require().NoError(err)

				return &icstypes.MsgSignalIntent{
					ChainId: s.chainB.ChainID,
					Intents: []*icstypes.ValidatorIntent{
						{
							ValoperAddress: val1.String(),
							Weight:         sdktypes.MustNewDecFromStr("3.0"),
						},
					},
					FromAddress: TestOwnerAddress,
				}
			},
			[]sdk.Dec{},
			true,
			false,
		},
		{
			"invalid - chain id",
			func(s *KeeperTestSuite) *icstypes.MsgSignalIntent {
				val1, err := sdktypes.ValAddressFromHex(s.chainB.Vals.Validators[0].Address.String())
				s.Require().NoError(err)

				return &icstypes.MsgSignalIntent{
					ChainId: s.chainA.ChainID,
					Intents: []*icstypes.ValidatorIntent{
						{
							ValoperAddress: val1.String(),
							Weight:         sdktypes.MustNewDecFromStr("1.0"),
						},
					},
					FromAddress: TestOwnerAddress,
				}
			},
			[]sdk.Dec{},
			false,
			true,
		},
		{
			"valid - single weight",
			func(s *KeeperTestSuite) *icstypes.MsgSignalIntent {
				val1, err := sdktypes.ValAddressFromHex(s.chainB.Vals.Validators[0].Address.String())
				s.Require().NoError(err)

				return &icstypes.MsgSignalIntent{
					ChainId: s.chainB.ChainID,
					Intents: []*icstypes.ValidatorIntent{
						{
							ValoperAddress: val1.String(),
							Weight:         sdktypes.MustNewDecFromStr("1.0"),
						},
					},
					FromAddress: TestOwnerAddress,
				}
			},
			[]sdk.Dec{sdk.NewDecWithPrec(1, 0)},
			false,
			false,
		},
		{
			"valid - multi weight",
			func(s *KeeperTestSuite) *icstypes.MsgSignalIntent {
				val1, err := sdktypes.ValAddressFromHex(s.chainB.Vals.Validators[0].Address.String())
				s.Require().NoError(err)
				val2, err := sdktypes.ValAddressFromHex(s.chainB.Vals.Validators[1].Address.String())
				s.Require().NoError(err)
				val3, err := sdktypes.ValAddressFromHex(s.chainB.Vals.Validators[2].Address.String())
				s.Require().NoError(err)

				return &icstypes.MsgSignalIntent{
					ChainId: s.chainB.ChainID,
					Intents: []*icstypes.ValidatorIntent{
						{
							ValoperAddress: val1.String(),
							Weight:         sdktypes.MustNewDecFromStr("0.5"),
						},
						{
							ValoperAddress: val2.String(),
							Weight:         sdktypes.MustNewDecFromStr("0.2"),
						},
						{
							ValoperAddress: val3.String(),
							Weight:         sdktypes.MustNewDecFromStr("0.3"),
						},
					},
					FromAddress: TestOwnerAddress,
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
			s.SetupZones()

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
			res, err := msgSrv.SignalIntent(sdktypes.WrapSDKContext(s.chainA.GetContext()), msg)
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

			intent, found := icsKeeper.GetIntent(s.chainA.GetContext(), zone, TestOwnerAddress, false)
			s.Require().True(found)
			intents := intent.GetIntents()

			for idx, weight := range tt.expected {
				val, err := sdktypes.ValAddressFromHex(s.chainB.Vals.Validators[idx].Address.String())
				s.Require().NoError(err)

				valIntent, found := intents[val.String()]
				s.Require().True(found)

				s.Require().Equal(weight, valIntent.Weight)
			}
		})
	}
}

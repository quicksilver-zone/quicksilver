package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/quicksilver-zone/quicksilver/x/tokenfactory/types"
)

// TestMintDenomMsg tests TypeMsgMint message is emitted on a successful mint.
func (s *KeeperTestSuite) TestMintDenomMsg() {
	// Create a denom
	s.CreateDefaultDenom()

	for _, tc := range []struct {
		desc                  string
		amount                int64
		mintDenom             string
		admin                 string
		valid                 bool
		expectedError         error
		expectedMessageEvents int
	}{
		{
			desc:          "denom does not exist",
			amount:        10,
			mintDenom:     "factory/QCK1t7egva48prqmzl59x5ngv4zx0dtrwewc9m7z44/evmos",
			admin:         s.TestAccs[0].String(),
			valid:         false,
			expectedError: types.ErrDenomDoesNotExist,
		},
		{
			desc:                  "success case",
			amount:                10,
			mintDenom:             s.defaultDenom,
			admin:                 s.TestAccs[0].String(),
			valid:                 true,
			expectedMessageEvents: 1,
		},
	} {
		s.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			ctx := s.Ctx.WithEventManager(sdk.NewEventManager())
			s.Require().Equal(0, len(ctx.EventManager().Events()))
			// Test mint message
			_, err := s.msgServer.Mint(sdk.WrapSDKContext(ctx), types.NewMsgMint(tc.admin, sdk.NewInt64Coin(tc.mintDenom, 10)))
			s.Require().ErrorIs(err, tc.expectedError)
			// Ensure current number and type of event is emitted
			s.AssertEventEmitted(ctx, types.TypeMsgMint, tc.expectedMessageEvents)
		})
	}
}

// TestBurnDenomMsg test TypeMsgBurn message is emitted on a successful burn.
func (s *KeeperTestSuite) TestBurnDenomMsg() {
	// Create a denom.
	s.CreateDefaultDenom()
	// mint 10 default token for testAcc[0]
	_, err := s.msgServer.Mint(sdk.WrapSDKContext(s.Ctx), types.NewMsgMint(s.TestAccs[0].String(), sdk.NewInt64Coin(s.defaultDenom, 10)))
	s.Require().NoError(err)

	for _, tc := range []struct {
		desc                  string
		amount                int64
		burnDenom             string
		admin                 string
		valid                 bool
		expectedError         error
		expectedMessageEvents int
	}{
		{
			desc:          "denom does not exist",
			burnDenom:     "factory/quick1vprpg84y4c50fxpf9ngza2y0p0q3k7yrw2q8tf/evmos",
			admin:         s.TestAccs[0].String(),
			valid:         false,
			expectedError: types.ErrUnauthorized,
		},
		{
			desc:                  "success case",
			burnDenom:             s.defaultDenom,
			admin:                 s.TestAccs[0].String(),
			valid:                 true,
			expectedMessageEvents: 1,
		},
	} {
		s.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			ctx := s.Ctx.WithEventManager(sdk.NewEventManager())
			s.Require().Equal(0, len(ctx.EventManager().Events()))
			// Test burn message
			_, err := s.msgServer.Burn(sdk.WrapSDKContext(ctx), types.NewMsgBurn(tc.admin, sdk.NewInt64Coin(tc.burnDenom, 10)))
			s.Require().ErrorIs(err, tc.expectedError)
			// Ensure current number and type of event is emitted
			s.AssertEventEmitted(ctx, types.TypeMsgBurn, tc.expectedMessageEvents)
		})
	}
}

// TestCreateDenomMsg test TypeMsgCreateDenom message is emitted on a successful denom creation.
func (s *KeeperTestSuite) TestCreateDenomMsg() {
	defaultDenomCreationFee := types.Params{DenomCreationFee: sdk.NewCoins(sdk.NewCoin("uqck", sdk.NewInt(50000000)))}
	for _, tc := range []struct {
		desc                  string
		denomCreationFee      types.Params
		subdenom              string
		valid                 bool
		expectedError         error
		expectedMessageEvents int
	}{
		{
			desc:             "subdenom too long",
			denomCreationFee: defaultDenomCreationFee,
			subdenom:         "assadsadsadasdasdsadsadsadsadsadsadsklkadaskkkdasdasedskhanhassyeunganassfnlksdflksafjlkasd",
			valid:            false,
			expectedError:    types.ErrSubdenomTooLong,
		},
		{
			desc:                  "success case: defaultDenomCreationFee",
			denomCreationFee:      defaultDenomCreationFee,
			subdenom:              "evmos",
			valid:                 true,
			expectedMessageEvents: 1,
		},
	} {
		s.SetupTest()
		s.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			tokenFactoryKeeper := s.App.TokenFactoryKeeper
			ctx := s.Ctx.WithEventManager(sdk.NewEventManager())
			s.Require().Equal(0, len(ctx.EventManager().Events()))
			// Set denom creation fee in params
			tokenFactoryKeeper.SetParams(s.Ctx, tc.denomCreationFee)
			// Test create denom message
			_, err := s.msgServer.CreateDenom(sdk.WrapSDKContext(ctx), types.NewMsgCreateDenom(s.TestAccs[0].String(), tc.subdenom))
			s.Require().ErrorIs(err, tc.expectedError)
			// Ensure current number and type of event is emitted
			s.AssertEventEmitted(ctx, types.TypeMsgCreateDenom, tc.expectedMessageEvents)
		})
	}
}

// TestChangeAdminDenomMsg test TypeMsgChangeAdmin message is emitted on a successful admin change.
func (s *KeeperTestSuite) TestChangeAdminDenomMsg() {
	for _, tc := range []struct {
		desc                    string
		msgChangeAdmin          func(denom string) *types.MsgChangeAdmin
		expectedChangeAdminPass bool
		expectedAdminIndex      int
		msgMint                 func(denom string) *types.MsgMint
		expectedMintPass        bool
		expectedError           error
		expectedMessageEvents   int
	}{
		{
			desc: "non-admins can't change the existing admin",
			msgChangeAdmin: func(denom string) *types.MsgChangeAdmin {
				return types.NewMsgChangeAdmin(s.TestAccs[1].String(), denom, s.TestAccs[2].String())
			},
			expectedChangeAdminPass: false,
			expectedError:           types.ErrUnauthorized,
			expectedAdminIndex:      0,
		},
		{
			desc: "success change admin",
			msgChangeAdmin: func(denom string) *types.MsgChangeAdmin {
				return types.NewMsgChangeAdmin(s.TestAccs[0].String(), denom, s.TestAccs[1].String())
			},
			expectedAdminIndex:      1,
			expectedChangeAdminPass: true,
			expectedMessageEvents:   1,
			msgMint: func(denom string) *types.MsgMint {
				return types.NewMsgMint(s.TestAccs[1].String(), sdk.NewInt64Coin(denom, 5))
			},
			expectedMintPass: true,
		},
	} {
		s.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			// setup test
			s.SetupTest()
			ctx := s.Ctx.WithEventManager(sdk.NewEventManager())
			s.Require().Equal(0, len(ctx.EventManager().Events()))
			// Create a denom and mint
			res, err := s.msgServer.CreateDenom(sdk.WrapSDKContext(ctx), types.NewMsgCreateDenom(s.TestAccs[0].String(), "bitcoin"))
			s.Require().NoError(err)
			testDenom := res.GetNewTokenDenom()
			_, err = s.msgServer.Mint(sdk.WrapSDKContext(ctx), types.NewMsgMint(s.TestAccs[0].String(), sdk.NewInt64Coin(testDenom, 10)))
			s.Require().NoError(err)

			// Test change admin message
			_, err = s.msgServer.ChangeAdmin(sdk.WrapSDKContext(ctx), tc.msgChangeAdmin(testDenom))
			s.Require().ErrorIs(err, tc.expectedError)

			// Ensure current number and type of event is emitted
			s.AssertEventEmitted(ctx, types.TypeMsgChangeAdmin, tc.expectedMessageEvents)
		})
	}
}

// TestSetDenomMetaDataMsg test TypeMsgSetDenomMetadata message is emitted on a successful denom metadata change.
func (s *KeeperTestSuite) TestSetDenomMetaDataMsg() {
	// setup test
	s.SetupTest()
	s.CreateDefaultDenom()

	for _, tc := range []struct {
		desc                  string
		msgSetDenomMetadata   types.MsgSetDenomMetadata
		expectedPass          bool
		expectedError         error
		expectedMessageEvents int
	}{
		{
			desc: "successful set denom metadata",
			msgSetDenomMetadata: *types.NewMsgSetDenomMetadata(s.TestAccs[0].String(), banktypes.Metadata{
				Description: "yeehaw",
				DenomUnits: []*banktypes.DenomUnit{
					{
						Denom:    s.defaultDenom,
						Exponent: 0,
					},
					{
						Denom:    "uqck",
						Exponent: 6,
					},
				},
				Base:    s.defaultDenom,
				Display: "uqck",
				Name:    "QCK",
				Symbol:  "QCK",
			}),
			expectedPass:          true,
			expectedMessageEvents: 1,
		},
		{
			desc: "non existent factory denom name",
			msgSetDenomMetadata: *types.NewMsgSetDenomMetadata(s.TestAccs[0].String(), banktypes.Metadata{
				Description: "yeehaw",
				DenomUnits: []*banktypes.DenomUnit{
					{
						Denom:    fmt.Sprintf("factory/%s/litecoin", s.TestAccs[0].String()),
						Exponent: 0,
					},
					{
						Denom:    "uqck",
						Exponent: 6,
					},
				},
				Base:    fmt.Sprintf("factory/%s/litecoin", s.TestAccs[0].String()),
				Display: "uqck",
				Name:    "QCK",
				Symbol:  "QCK",
			}),
			expectedPass:  false,
			expectedError: types.ErrUnauthorized,
		},
	} {
		s.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			tc := tc
			ctx := s.Ctx.WithEventManager(sdk.NewEventManager())
			s.Require().Equal(0, len(ctx.EventManager().Events()))
			// Test set denom metadata message
			_, err := s.msgServer.SetDenomMetadata(sdk.WrapSDKContext(ctx), &tc.msgSetDenomMetadata)
			s.Require().ErrorIs(err, tc.expectedError)
			// Ensure current number and type of event is emitted
			s.AssertEventEmitted(ctx, types.TypeMsgSetDenomMetadata, tc.expectedMessageEvents)
		})
	}
}

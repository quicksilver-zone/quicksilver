package keeper_test

import (
	"fmt"

	"github.com/quicksilver-zone/quicksilver/x/tokenfactory/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *KeeperTestSuite) TestMsgCreateDenom() {
	var (
		tokenFactoryKeeper = s.App.TokenFactoryKeeper
		bankKeeper         = s.App.BankKeeper
		denomCreationFee   = tokenFactoryKeeper.GetParams(s.Ctx).DenomCreationFee
	)

	// Get balance of acc 0 before creating a denom
	preCreateBalance := bankKeeper.GetBalance(s.Ctx, s.TestAccs[0], denomCreationFee[0].Denom)

	// Creating a denom should work
	res, err := s.msgServer.CreateDenom(sdk.WrapSDKContext(s.Ctx), types.NewMsgCreateDenom(s.TestAccs[0].String(), "bitcoin"))
	s.Require().NoError(err)
	s.Require().NotEmpty(res.GetNewTokenDenom())

	// Make sure that the admin is set correctly
	queryRes, err := s.queryClient.DenomAuthorityMetadata(s.Ctx.Context(), &types.QueryDenomAuthorityMetadataRequest{
		Denom: res.GetNewTokenDenom(),
	})
	s.Require().NoError(err)
	s.Require().Equal(s.TestAccs[0].String(), queryRes.AuthorityMetadata.Admin)

	// Make sure that creation fee was deducted
	postCreateBalance := bankKeeper.GetBalance(s.Ctx, s.TestAccs[0], tokenFactoryKeeper.GetParams(s.Ctx).DenomCreationFee[0].Denom)
	s.Require().True(preCreateBalance.Sub(postCreateBalance).IsEqual(denomCreationFee[0]))

	// Make sure that a second version of the same denom can't be recreated
	_, err = s.msgServer.CreateDenom(sdk.WrapSDKContext(s.Ctx), types.NewMsgCreateDenom(s.TestAccs[0].String(), "bitcoin"))
	s.Require().Error(err)

	// Creating a second denom should work
	res, err = s.msgServer.CreateDenom(sdk.WrapSDKContext(s.Ctx), types.NewMsgCreateDenom(s.TestAccs[0].String(), "litecoin"))
	s.Require().NoError(err)
	s.Require().NotEmpty(res.GetNewTokenDenom())

	// Try querying all the denoms created by suite.TestAccs[0]
	queryRes2, err := s.queryClient.DenomsFromCreator(s.Ctx.Context(), &types.QueryDenomsFromCreatorRequest{
		Creator: s.TestAccs[0].String(),
	})
	s.Require().NoError(err)
	s.Require().Len(queryRes2.Denoms, 2)

	// Make sure that a second account can create a denom with the same subdenom
	res, err = s.msgServer.CreateDenom(sdk.WrapSDKContext(s.Ctx), types.NewMsgCreateDenom(s.TestAccs[1].String(), "bitcoin"))
	s.Require().NoError(err)
	s.Require().NotEmpty(res.GetNewTokenDenom())

	// Make sure that an address with a "/" in it can't create denoms
	_, err = s.msgServer.CreateDenom(sdk.WrapSDKContext(s.Ctx), types.NewMsgCreateDenom("quick.eth/creator", "bitcoin"))
	s.Require().Error(err)
}

func (s *KeeperTestSuite) TestCreateDenom() {
	var (
		primaryDenom            = types.DefaultParams().DenomCreationFee[0].Denom
		secondaryDenom          = SecondaryDenom
		defaultDenomCreationFee = types.Params{DenomCreationFee: sdk.NewCoins(sdk.NewCoin(primaryDenom, sdk.NewInt(50000000)))}
		twoDenomCreationFee     = types.Params{DenomCreationFee: sdk.NewCoins(sdk.NewCoin(primaryDenom, sdk.NewInt(50000000)), sdk.NewCoin(secondaryDenom, sdk.NewInt(50000000)))}
		nilCreationFee          = types.Params{DenomCreationFee: nil}
		largeCreationFee        = types.Params{DenomCreationFee: sdk.NewCoins(sdk.NewCoin(primaryDenom, sdk.NewInt(5000000000)))}
	)

	for _, tc := range []struct {
		desc             string
		denomCreationFee types.Params
		setup            func()
		subdenom         string
		valid            bool
	}{
		{
			desc:             "subdenom too long",
			denomCreationFee: defaultDenomCreationFee,
			subdenom:         "assadsadsadasdasdsadsadsadsadsadsadsklkadaskkkdasdasedskhanhassyeunganassfnlksdflksafjlkasd",
			valid:            false,
		},
		{
			desc:             "subdenom and creator pair already exists",
			denomCreationFee: defaultDenomCreationFee,
			setup: func() {
				_, err := s.msgServer.CreateDenom(sdk.WrapSDKContext(s.Ctx), types.NewMsgCreateDenom(s.TestAccs[0].String(), "bitcoin"))
				s.Require().NoError(err)
			},
			subdenom: "bitcoin",
			valid:    false,
		},
		{
			desc:             "success case: defaultDenomCreationFee",
			denomCreationFee: defaultDenomCreationFee,
			subdenom:         "evmos",
			valid:            true,
		},
		{
			desc:             "success case: twoDenomCreationFee",
			denomCreationFee: twoDenomCreationFee,
			subdenom:         "catcoin",
			valid:            true,
		},
		{
			desc:             "success case: nilCreationFee",
			denomCreationFee: nilCreationFee,
			subdenom:         "czcoin",
			valid:            true,
		},
		{
			desc:             "account doesn't have enough to pay for denom creation fee",
			denomCreationFee: largeCreationFee,
			subdenom:         "tooexpensive",
			valid:            false,
		},
		{
			desc:             "subdenom having invalid characters",
			denomCreationFee: defaultDenomCreationFee,
			subdenom:         "bit/***///&&&/coin",
			valid:            false,
		},
	} {
		s.SetupTest()
		s.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			if tc.setup != nil {
				tc.setup()
			}
			tokenFactoryKeeper := s.App.TokenFactoryKeeper
			bankKeeper := s.App.BankKeeper
			// Set denom creation fee in params
			tokenFactoryKeeper.SetParams(s.Ctx, tc.denomCreationFee)
			denomCreationFee := tokenFactoryKeeper.GetParams(s.Ctx).DenomCreationFee
			s.Require().Equal(tc.denomCreationFee.DenomCreationFee, denomCreationFee)

			// note balance, create a tokenfactory denom, then note balance again
			preCreateBalance := bankKeeper.GetAllBalances(s.Ctx, s.TestAccs[0])
			res, err := s.msgServer.CreateDenom(sdk.WrapSDKContext(s.Ctx), types.NewMsgCreateDenom(s.TestAccs[0].String(), tc.subdenom))
			postCreateBalance := bankKeeper.GetAllBalances(s.Ctx, s.TestAccs[0])
			if tc.valid {
				s.Require().NoError(err)
				s.Require().True(preCreateBalance.Sub(postCreateBalance...).IsEqual(denomCreationFee))

				// Make sure that the admin is set correctly
				queryRes, err := s.queryClient.DenomAuthorityMetadata(s.Ctx.Context(), &types.QueryDenomAuthorityMetadataRequest{
					Denom: res.GetNewTokenDenom(),
				})

				s.Require().NoError(err)
				s.Require().Equal(s.TestAccs[0].String(), queryRes.AuthorityMetadata.Admin)

			} else {
				s.Require().Error(err)
				// Ensure we don't charge if we expect an error
				s.Require().True(preCreateBalance.IsEqual(postCreateBalance))
			}
		})
	}
}

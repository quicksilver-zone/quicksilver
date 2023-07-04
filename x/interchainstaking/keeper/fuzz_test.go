package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/ingenuity-build/quicksilver/utils/addressutils"
	icqtypes "github.com/ingenuity-build/quicksilver/x/interchainquery/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

type FuzzingTestSuite struct {
	KeeperTestSuite
}

func FuzzZones(f *testing.F) {
	if testing.Short() {
		f.Skip("In -short mode")
	}

	// 1. Generate the seeds.
	suite := new(FuzzingTestSuite)
	suite.SetT(new(testing.T))
	suite.SetupTest()

	suite.setupTestZones()
	app := suite.GetQuicksilverApp(suite.chainA)
	app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()

	seeds := []*icstypes.QueryZonesRequest{
		{},
		nil,
		{
			Pagination: &query.PageRequest{},
		},
		{
			Pagination: &query.PageRequest{},
		},
		{
			Pagination: &query.PageRequest{
				Offset:  10,
				Reverse: true,
				Limit:   icstypes.TxRetrieveCount,
			},
		},
	}

	for _, seed := range seeds {
		bz, err := app.AppCodec().Marshal(seed)
		suite.Require().NoError(err)
		f.Add(bz)
	}

	// 2. Now fuzz the code.
	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

	f.Fuzz(func(t *testing.T, reqBz []byte) {
		switch str := string(reqBz); str {
		// Manually skipping over known and reported vectors
		// as we know they cause crashes.
		case "\n\t\n\x01000 0(0", "\n\t\n\x03000 0(0": // https://github.com/ingenuity-build/quicksilver-incognito/issues/88
			return
		case "\n\t\n\x01K\x10\x0000(0", "\n\t\n\x030D0 0(0", "\n\t\n\x0301000(0":
			return
		}

		suite := new(FuzzingTestSuite)
		suite.SetT(new(testing.T))
		suite.SetupTest()
		suite.setupTestZones()
		app := suite.GetQuicksilverApp(suite.chainA)
		app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()

		req := new(icstypes.QueryZonesRequest)
		if err := app.AppCodec().Unmarshal(reqBz, req); err != nil {
			// Do nothing with an invalid ZonesInfoRequest.
			return
		}

		if pag := req.Pagination; pag != nil {
			// A case already seen.
			return
		}
		_, err := icsKeeper.Zones(ctx, req)
		require.NoError(t, err)
	})
}

func FuzzValsetCallback(f *testing.F) {
	// 1. Generate the seeds.
	newVal := addressutils.GenerateValAddressForTest()
	valSetFuncs := []func(in stakingtypes.Validators) stakingtypes.QueryValidatorsResponse{
		func(in stakingtypes.Validators) stakingtypes.QueryValidatorsResponse {
			val := in[0]
			val.OperatorAddress = newVal.String()
			in = append(in, val)
			return stakingtypes.QueryValidatorsResponse{Validators: in}
		},
		func(in stakingtypes.Validators) stakingtypes.QueryValidatorsResponse {
			in[1].DelegatorShares = in[1].DelegatorShares.Add(sdk.NewDec(1000))
			in[2].DelegatorShares = in[2].DelegatorShares.Add(sdk.NewDec(2000))
			return stakingtypes.QueryValidatorsResponse{Validators: in}
		},
		func(in stakingtypes.Validators) stakingtypes.QueryValidatorsResponse {
			in[0].Tokens = in[0].Tokens.Add(sdk.NewInt(1000))
			return stakingtypes.QueryValidatorsResponse{Validators: in}
		},
		func(in stakingtypes.Validators) stakingtypes.QueryValidatorsResponse {
			in[1].Tokens = in[1].Tokens.Add(sdk.NewInt(1000))
			in[2].Tokens = in[2].Tokens.Add(sdk.NewInt(2000))
			return stakingtypes.QueryValidatorsResponse{Validators: in}
		},
		func(in stakingtypes.Validators) stakingtypes.QueryValidatorsResponse {
			in[1].Tokens = in[1].Tokens.Sub(sdk.NewInt(10))
			in[2].Tokens = in[2].Tokens.Sub(sdk.NewInt(20))
			return stakingtypes.QueryValidatorsResponse{Validators: in}
		},
		func(in stakingtypes.Validators) stakingtypes.QueryValidatorsResponse {
			in[0].Commission.CommissionRates.Rate = sdk.NewDecWithPrec(5, 1)
			in[2].Commission.CommissionRates.Rate = sdk.NewDecWithPrec(5, 2)
			return stakingtypes.QueryValidatorsResponse{Validators: in}
		},
		func(in stakingtypes.Validators) stakingtypes.QueryValidatorsResponse {
			val := in[0]
			val.OperatorAddress = newVal.String()
			in = append(in, val)
			return stakingtypes.QueryValidatorsResponse{Validators: in}
		},
	}

	suite := new(FuzzingTestSuite)
	suite.SetT(new(testing.T))
	suite.SetupTest()
	suite.setupTestZones()

	for _, valFunc := range valSetFuncs {
		// 1.5. Set up a fresh test suite given that valFunc can mutate inputs.
		chainBVals := suite.GetQuicksilverApp(suite.chainB).StakingKeeper.GetValidators(suite.chainB.GetContext(), 300)
		queryRes := valFunc(chainBVals)
		app := suite.GetQuicksilverApp(suite.chainA)
		bz, err := app.AppCodec().Marshal(&queryRes)
		suite.Require().NoError(err)
		f.Add(bz)
	}

	// 2. Now fuzz.
	f.Fuzz(func(t *testing.T, args []byte) {
		suite.SetT(t)
		suite.FuzzValsetCallback(args)
	})
}

func (suite *FuzzingTestSuite) FuzzValsetCallback(args []byte) {
	app := suite.GetQuicksilverApp(suite.chainA)
	app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
	ctx := suite.chainA.GetContext()

	err := keeper.ValsetCallback(app.InterchainstakingKeeper, ctx, args, icqtypes.Query{ChainId: suite.chainB.ChainID})
	suite.Require().NoError(err)
}

package keeper_test

import (
	"bytes"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/ingenuity-build/quicksilver/utils"
	icqtypes "github.com/ingenuity-build/quicksilver/x/interchainquery/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

type FuzzingTestSuite struct {
	KeeperTestSuite
}

func FuzzZoneInfos(f *testing.F) {
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

	seeds := []*icstypes.QueryZonesInfoRequest{
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
				Key:     []byte("key"),
				Offset:  10,
				Reverse: true,
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

		req := new(icstypes.QueryZonesInfoRequest)
		if err := app.AppCodec().Unmarshal(reqBz, req); err != nil {
			// Do nothing with an invalid ZonesInfoRequest.
			return
		}

		if pag := req.Pagination; pag != nil && bytes.Count(pag.Key, []byte("0")) == len(pag.Key) {
			// A case already seen.
			return
		}
		_, err := icsKeeper.Zones(ctx, req)
		require.NoError(t, err)
	})
}

func TestInvalidPaginationForQueryZones(t *testing.T) {
	t.Skip("Not yet fixed per https://github.com/ingenuity-build/quicksilver-incognito/issues/88")

	suite := new(FuzzingTestSuite)
	suite.SetT(t)
	suite.SetupTest()
	suite.setupTestZones()
	app := suite.GetQuicksilverApp(suite.chainA)
	app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
	icsKeeper := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper
	ctx := suite.chainA.GetContext()

	reqBz := []byte("\n\t\n\x03000 0(0")
	req := new(icstypes.QueryZonesInfoRequest)
	if err := app.AppCodec().Unmarshal(reqBz, req); err != nil {
		// Do nothing with an invalid ZonesInfoRequest.
		return
	}

	_, err := icsKeeper.Zones(ctx, req)
	require.NoError(t, err)
}

func FuzzValsetCallback(f *testing.F) {
	// 1. Generate the seeds.
	newVal := utils.GenerateValAddressForTest()
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

func FuzzDelegationsCallback(f *testing.F) {
	// 1. Add the samples firstly.
	suite := new(FuzzingTestSuite)
	suite.SetT(new(testing.T))
	suite.SetupTest()
	suite.setupTestZones()

	app := suite.GetQuicksilverApp(suite.chainA)
	app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()

	// 1.5. Create the queries.
	zone, ok := app.InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), suite.chainB.ChainID)
	if !ok {
		f.Fatalf("Could not retrieve zone for chainB: %q", suite.chainB.ChainID)
	}
	queries := make([]*stakingtypes.QueryDelegatorDelegationsRequest, 2)
	for i, addr := range []string{zone.DepositAddress.Address, zone.WithdrawalAddress.Address} {
		accAddr, err := sdk.AccAddressFromBech32(addr)
		suite.Require().NoError(err)
		queries[i] = &stakingtypes.QueryDelegatorDelegationsRequest{
			DelegatorAddr: accAddr.String(),
		}
	}

	for _, query := range queries {
		bz, err := app.AppCodec().Marshal(query)
		suite.Require().NoError(err)
		f.Add(bz)
	}

	f.Fuzz(func(t *testing.T, args []byte) {
		suite.SetT(t)
		suite.FuzzDelegationsCallback(args)
	})
}

func FuzzAccountBalanceCallback(f *testing.F) {
	// 1. Add the samples firstly.
	suite := new(FuzzingTestSuite)
	suite.SetT(new(testing.T))
	suite.SetupTest()
	suite.setupTestZones()

	app := suite.GetQuicksilverApp(suite.chainA)

	values := []int64{10, 0, 100, 1_000, 1_000_000}
	for _, val := range values {
		response := sdk.NewCoin("qck", sdk.NewInt(val))
		respbz, err := app.AppCodec().Marshal(&response)
		suite.Require().NoError(err)
		f.Add(respbz)
	}

	// 2. Now fuzz.
	f.Fuzz(func(t *testing.T, respbz []byte) {
		suite.SetT(t)
		suite.FuzzAccountBalanceCallback(respbz)
	})
}

func FuzzAllBalancesCallback(f *testing.F) {
	// 1. Add the samples firstly.
	suite := new(FuzzingTestSuite)
	suite.SetT(new(testing.T))
	suite.SetupTest()
	suite.setupTestZones()

	// 1. Add corpus from chainA.
	app := suite.GetQuicksilverApp(suite.chainA)
	zone, ok := app.InterchainstakingKeeper.GetZone(suite.chainA.GetContext(), suite.chainB.ChainID)
	if !ok {
		f.Fatalf("Could not retrieve zone for chainB: %q", suite.chainB.ChainID)
	}
	reqbz, err := app.AppCodec().Marshal(&banktypes.QueryAllBalancesRequest{
		Address: zone.DepositAddress.Address,
	})
	suite.Require().NoError(err)
	f.Add(reqbz)

	if false {
		// 2. Add corpus from chainB.
		app = suite.GetQuicksilverApp(suite.chainB)
		zone, ok = app.InterchainstakingKeeper.GetZone(suite.chainB.GetContext(), suite.chainA.ChainID)
		if !ok {
			f.Fatalf("Could not retrieve zone for chainA: %q", suite.chainA.ChainID)
		}
		reqbz, err = app.AppCodec().Marshal(&banktypes.QueryAllBalancesRequest{
			Address: zone.DepositAddress.Address,
		})
		suite.Require().NoError(err)
		f.Add(reqbz)
	}

	// 3. Now fuzz.
	f.Fuzz(func(t *testing.T, respbz []byte) {
		suite.SetT(t)
		suite.FuzzAllBalancesCallback(respbz)
	})
}

func (s *FuzzingTestSuite) FuzzAccountBalanceCallback(respbz []byte) {
	if testing.Short() {
		s.T().Skip("In -short mode")
	}

	app := s.GetQuicksilverApp(s.chainA)
	app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
	ctx := s.chainA.GetContext()

	zone, _ := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
	zone.DepositAddress.BalanceWaitgroup++
	zone.WithdrawalAddress.BalanceWaitgroup++
	app.InterchainstakingKeeper.SetZone(ctx, &zone)

	for _, addr := range []string{zone.DepositAddress.Address, zone.WithdrawalAddress.Address} {
		accAddr, err := sdk.AccAddressFromBech32(addr)
		s.Require().NoError(err)
		data := append(banktypes.CreateAccountBalancesPrefix(accAddr), []byte("stake")...)

		err = keeper.AccountBalanceCallback(&app.InterchainstakingKeeper, ctx, respbz, icqtypes.Query{ChainID: s.chainB.ChainID, Request: data})
		s.Require().NoError(err)
	}
}

func (s *FuzzingTestSuite) FuzzDelegationsCallback(respbz []byte) {
	if testing.Short() {
		s.T().Skip("In -short mode")
	}

	app := s.GetQuicksilverApp(s.chainA)
	app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
	ctx := s.chainA.GetContext()

	zone, _ := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
	queryReq := stakingtypes.QueryDelegatorDelegationsRequest{
		DelegatorAddr: zone.DelegationAddress.Address,
	}
	reqbz, err := app.AppCodec().Marshal(&queryReq)
	s.Require().NoError(err)

	err = keeper.DelegationsCallback(&app.InterchainstakingKeeper, ctx, respbz, icqtypes.Query{ChainID: s.chainB.ChainID, Request: reqbz})
	s.Require().NoError(err)
}

func (s *FuzzingTestSuite) FuzzAllBalancesCallback(respbz []byte) {
	if testing.Short() {
		s.T().Skip("In -short mode")
	}

	app := s.GetQuicksilverApp(s.chainA)
	app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
	ctx := s.chainA.GetContext()

	zone, _ := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)

	queryReq := banktypes.QueryAllBalancesRequest{
		Address: zone.DepositAddress.Address,
	}
	reqbz, err := app.AppCodec().Marshal(&queryReq)
	s.Require().NoError(err)

	err = keeper.AllBalancesCallback(&app.InterchainstakingKeeper, ctx, respbz, icqtypes.Query{ChainID: s.chainB.ChainID, Request: reqbz})
	s.Require().NoError(err)
}

func (s *FuzzingTestSuite) FuzzValsetCallback(args []byte) {
	app := s.GetQuicksilverApp(s.chainA)
	app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
	ctx := s.chainA.GetContext()

	err := keeper.ValsetCallback(&app.InterchainstakingKeeper, ctx, args, icqtypes.Query{ChainID: s.chainB.ChainID})
	s.Require().NoError(err)
}

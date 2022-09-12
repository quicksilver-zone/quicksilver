package keeper_test

import (
	"bytes"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ingenuity-build/quicksilver/app"
	"github.com/ingenuity-build/quicksilver/utils"
	icqtypes "github.com/ingenuity-build/quicksilver/x/interchainquery/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	"github.com/stretchr/testify/require"
)

func TestCoinFromRequestKey(t *testing.T) {
	accAddr := utils.GenerateAccAddressForTest()
	prefix := banktypes.CreateAccountBalancesPrefix(accAddr.Bytes())
	query := append(prefix, []byte("denom")...)

	coin, err := keeper.CoinFromRequestKey(query, accAddr)
	require.NoError(t, err)
	require.Equal(t, "denom", coin.Denom)
}

// ValSetCallback

func (s *KeeperTestSuite) TestHandleValsetCallback() {
	newVal := utils.GenerateValAddressForTest()

	tests := []struct {
		name   string
		valset func(in stakingtypes.Validators) stakingtypes.QueryValidatorsResponse
		checks func(require *require.Assertions, ctx sdk.Context, app *app.Quicksilver, in stakingtypes.Validators)
	}{
		{
			name: "valid - no-op",
			valset: func(in stakingtypes.Validators) stakingtypes.QueryValidatorsResponse {
				return stakingtypes.QueryValidatorsResponse{Validators: in}
			},
			checks: func(_ *require.Assertions, _ sdk.Context, _ *app.Quicksilver, _ stakingtypes.Validators) {
				// no op
			},
		},
		{
			name: "valid - shares +1000 val[0]",
			valset: func(in stakingtypes.Validators) stakingtypes.QueryValidatorsResponse {
				in[0].DelegatorShares = in[0].DelegatorShares.Add(sdk.NewDec(1000))
				return stakingtypes.QueryValidatorsResponse{Validators: in}
			},
			checks: func(require *require.Assertions, ctx sdk.Context, app *app.Quicksilver, in stakingtypes.Validators) {
				foundQuery := false
				_, addr, _ := bech32.DecodeAndConvert(in[0].OperatorAddress)
				data := stakingtypes.GetValidatorKey(addr)
				for _, i := range app.InterchainQueryKeeper.AllQueries(ctx) {
					if i.QueryType == "store/staking/key" && bytes.Equal(i.Request, data) {
						foundQuery = true
						break
					}
				}
				require.True(foundQuery)
			},
		},
		{
			name: "valid - shares +1000 val[1], +2000 val[2]",
			valset: func(in stakingtypes.Validators) stakingtypes.QueryValidatorsResponse {
				in[1].DelegatorShares = in[1].DelegatorShares.Add(sdk.NewDec(1000))
				in[2].DelegatorShares = in[2].DelegatorShares.Add(sdk.NewDec(2000))
				return stakingtypes.QueryValidatorsResponse{Validators: in}
			},
			checks: func(require *require.Assertions, ctx sdk.Context, app *app.Quicksilver, in stakingtypes.Validators) {
				foundQuery := false
				foundQuery2 := false
				_, addr, _ := bech32.DecodeAndConvert(in[1].OperatorAddress)
				data := stakingtypes.GetValidatorKey(addr)
				_, addr2, _ := bech32.DecodeAndConvert(in[2].OperatorAddress)
				data2 := stakingtypes.GetValidatorKey(addr2)
				for _, i := range app.InterchainQueryKeeper.AllQueries(ctx) {
					if i.QueryType == "store/staking/key" && bytes.Equal(i.Request, data) {
						foundQuery = true
					}
					if i.QueryType == "store/staking/key" && bytes.Equal(i.Request, data2) {
						foundQuery2 = true
					}
				}
				require.True(foundQuery)
				require.True(foundQuery2)
			},
		},
		{
			name: "valid - tokens +1000 val[0]",
			valset: func(in stakingtypes.Validators) stakingtypes.QueryValidatorsResponse {
				in[0].Tokens = in[0].Tokens.Add(sdk.NewInt(1000))
				return stakingtypes.QueryValidatorsResponse{Validators: in}
			},
			checks: func(require *require.Assertions, ctx sdk.Context, app *app.Quicksilver, in stakingtypes.Validators) {
				foundQuery := false
				_, addr, _ := bech32.DecodeAndConvert(in[0].OperatorAddress)
				data := stakingtypes.GetValidatorKey(addr)
				for _, i := range app.InterchainQueryKeeper.AllQueries(ctx) {
					if i.QueryType == "store/staking/key" && bytes.Equal(i.Request, data) {
						foundQuery = true
						break
					}
				}
				require.True(foundQuery)
			},
		},
		{
			name: "valid - tokens +1000 val[1], +2000 val[2]",
			valset: func(in stakingtypes.Validators) stakingtypes.QueryValidatorsResponse {
				in[1].Tokens = in[1].Tokens.Add(sdk.NewInt(1000))
				in[2].Tokens = in[2].Tokens.Add(sdk.NewInt(2000))
				return stakingtypes.QueryValidatorsResponse{Validators: in}
			},
			checks: func(require *require.Assertions, ctx sdk.Context, app *app.Quicksilver, in stakingtypes.Validators) {
				foundQuery := false
				foundQuery2 := false
				_, addr, _ := bech32.DecodeAndConvert(in[1].OperatorAddress)
				data := stakingtypes.GetValidatorKey(addr)
				_, addr2, _ := bech32.DecodeAndConvert(in[2].OperatorAddress)
				data2 := stakingtypes.GetValidatorKey(addr2)
				for _, i := range app.InterchainQueryKeeper.AllQueries(ctx) {
					if i.QueryType == "store/staking/key" && bytes.Equal(i.Request, data) {
						foundQuery = true
					}
					if i.QueryType == "store/staking/key" && bytes.Equal(i.Request, data2) {
						foundQuery2 = true
					}
				}
				require.True(foundQuery)
				require.True(foundQuery2)
			},
		},
		{
			name: "valid - tokens -10 val[1], -20 val[2]",
			valset: func(in stakingtypes.Validators) stakingtypes.QueryValidatorsResponse {
				in[1].Tokens = in[1].Tokens.Sub(sdk.NewInt(10))
				in[2].Tokens = in[2].Tokens.Sub(sdk.NewInt(20))
				return stakingtypes.QueryValidatorsResponse{Validators: in}
			},
			checks: func(require *require.Assertions, ctx sdk.Context, app *app.Quicksilver, in stakingtypes.Validators) {
				foundQuery := false
				foundQuery2 := false
				_, addr, _ := bech32.DecodeAndConvert(in[1].OperatorAddress)
				data := stakingtypes.GetValidatorKey(addr)
				_, addr2, _ := bech32.DecodeAndConvert(in[2].OperatorAddress)
				data2 := stakingtypes.GetValidatorKey(addr2)
				for _, i := range app.InterchainQueryKeeper.AllQueries(ctx) {
					if i.QueryType == "store/staking/key" && bytes.Equal(i.Request, data) {
						foundQuery = true
					}
					if i.QueryType == "store/staking/key" && bytes.Equal(i.Request, data2) {
						foundQuery2 = true
					}
				}
				require.True(foundQuery)
				require.True(foundQuery2)
			},
		},
		{
			name: "valid - commission 0.5 val[0], 0.05 val[2]",
			valset: func(in stakingtypes.Validators) stakingtypes.QueryValidatorsResponse {
				in[0].Commission.CommissionRates.Rate = sdk.NewDecWithPrec(5, 1)
				in[2].Commission.CommissionRates.Rate = sdk.NewDecWithPrec(5, 2)
				return stakingtypes.QueryValidatorsResponse{Validators: in}
			},
			checks: func(require *require.Assertions, ctx sdk.Context, app *app.Quicksilver, in stakingtypes.Validators) {
				foundQuery := false
				foundQuery2 := false
				_, addr, _ := bech32.DecodeAndConvert(in[0].OperatorAddress)
				data := stakingtypes.GetValidatorKey(addr)
				_, addr2, _ := bech32.DecodeAndConvert(in[2].OperatorAddress)
				data2 := stakingtypes.GetValidatorKey(addr2)
				for _, i := range app.InterchainQueryKeeper.AllQueries(ctx) {
					if i.QueryType == "store/staking/key" && bytes.Equal(i.Request, data) {
						foundQuery = true
					}
					if i.QueryType == "store/staking/key" && bytes.Equal(i.Request, data2) {
						foundQuery2 = true
					}
				}
				require.True(foundQuery)
				require.True(foundQuery2)
			},
		},
		{
			name: "valid - new validator",
			valset: func(in stakingtypes.Validators) stakingtypes.QueryValidatorsResponse {
				new := in[0]
				new.OperatorAddress = newVal.String()
				in = append(in, new)
				return stakingtypes.QueryValidatorsResponse{Validators: in}
			},
			checks: func(require *require.Assertions, ctx sdk.Context, app *app.Quicksilver, in stakingtypes.Validators) {
				foundQuery := false
				data := stakingtypes.GetValidatorKey(newVal)
				for _, i := range app.InterchainQueryKeeper.AllQueries(ctx) {
					if i.QueryType == "store/staking/key" && bytes.Equal(i.Request, data) {
						foundQuery = true
					}
				}
				require.True(foundQuery)
			},
		},
		// TODO: trigger callback on status change.
		// {
		// 	name: "valid - status unbonding val[0]",
		// 	valset: func(in stakingtypes.Validators) stakingtypes.QueryValidatorsResponse {
		// 		in[0].Status = stakingtypes.Unbonding
		// 		return stakingtypes.QueryValidatorsResponse{Validators: in}
		// 	},
		// 	checks: func(require *require.Assertions, ctx sdk.Context, app *app.Quicksilver, in stakingtypes.Validators) {
		// 		foundQuery := false
		// 		_, addr, _ := bech32.DecodeAndConvert(in[0].OperatorAddress)
		// 		data := stakingtypes.GetValidatorKey(addr)
		// 		for _, i := range app.InterchainQueryKeeper.AllQueries(ctx) {
		// 			if i.QueryType == "store/staking/key" && bytes.Equal(i.Request, data) {
		// 				foundQuery = true
		// 			}
		// 		}
		// 		require.True(foundQuery)
		// 	},
		// },
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.SetupTest()
			s.SetupZones()

			app := s.GetQuicksilverApp(s.chainA)
			app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
			ctx := s.chainA.GetContext()

			chainBVals := s.GetQuicksilverApp(s.chainB).StakingKeeper.GetValidators(s.chainB.GetContext(), 300)

			query := test.valset(chainBVals)
			bz, err := app.AppCodec().Marshal(&query)
			s.Require().NoError(err)

			err = keeper.ValsetCallback(app.InterchainstakingKeeper, ctx, bz, icqtypes.Query{ChainId: s.chainB.ChainID})
			s.Require().NoError(err)
			// valset callback doesn't actually update validators, but does emit icq callbacks.
			test.checks(s.Require(), ctx, app, chainBVals)
		})
	}
}

func (s *KeeperTestSuite) TestHandleValsetCallbackBadChain() {
	s.Run("bad chain", func() {
		s.SetupTest()
		s.SetupZones()

		app := s.GetQuicksilverApp(s.chainA)
		app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := s.chainA.GetContext()

		query := stakingtypes.QueryValidatorsResponse{Validators: []stakingtypes.Validator{}}
		bz, err := app.AppCodec().Marshal(&query)
		s.Require().NoError(err)

		err = keeper.ValsetCallback(app.InterchainstakingKeeper, ctx, bz, icqtypes.Query{ChainId: "badchain"})
		// this should bail on a non-matching chain id.
		s.Require().Error(err)
	})
}

func (s *KeeperTestSuite) TestHandleValsetCallbackNilValset() {
	s.Run("nil valset", func() {
		s.SetupTest()
		s.SetupZones()

		app := s.GetQuicksilverApp(s.chainA)
		app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := s.chainA.GetContext()

		query := stakingtypes.QueryValidatorsResponse{Validators: []stakingtypes.Validator{}}
		bz, err := app.AppCodec().Marshal(&query)
		s.Require().NoError(err)

		err = keeper.ValsetCallback(app.InterchainstakingKeeper, ctx, bz, icqtypes.Query{ChainId: s.chainB.ChainID})
		// this should error on unmarshalling an empty slice, which is not a valid response here.
		s.Require().Error(err)
	})
}

func (s *KeeperTestSuite) TestHandleValsetCallbackInvalidResponse() {
	s.Run("bad payload type", func() {
		s.SetupTest()
		s.SetupZones()

		app := s.GetQuicksilverApp(s.chainA)
		app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := s.chainA.GetContext()

		query := stakingtypes.QueryValidatorsRequest{Status: stakingtypes.BondStatusBonded}
		bz, err := app.AppCodec().Marshal(&query)
		s.Require().NoError(err)

		err = keeper.ValsetCallback(app.InterchainstakingKeeper, ctx, bz, icqtypes.Query{ChainId: s.chainB.ChainID})
		// this should error on unmarshalling an empty slice, which is not a valid response here.
		s.Require().Error(err)
	})
}

// ValidatorCallback

// func (s *KeeperTestSuite) TestHandleValidatorCallbackInvalidResponse() {
// 	s.Run("bad payload type", func() {
// 		s.SetupTest()
// 		s.SetupZones()

// 		app := s.GetQuicksilverApp(s.chainA)
// 		app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
// 		ctx := s.chainA.GetContext()

// 		query := stakingtypes.QueryValidatorsRequest{Status: stakingtypes.BondStatusBonded}
// 		bz, err := app.AppCodec().Marshal(&query)
// 		s.Require().NoError(err)

// 		err = keeper.ValidatorCallback(app.InterchainstakingKeeper, ctx, bz, icqtypes.Query{ChainId: s.chainB.ChainID})
// 		// this should error on unmarshalling an empty slice, which is not a valid response here.
// 		s.Require().Error(err)
// 	})
// }

func (s *KeeperTestSuite) TestHandleValidatorCallbackBadChain() {
	s.Run("bad chain", func() {
		s.SetupTest()
		s.SetupZones()

		app := s.GetQuicksilverApp(s.chainA)
		app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := s.chainA.GetContext()

		query := stakingtypes.QueryValidatorsResponse{Validators: []stakingtypes.Validator{}}
		bz, err := app.AppCodec().Marshal(&query)
		s.Require().NoError(err)

		err = keeper.ValidatorCallback(app.InterchainstakingKeeper, ctx, bz, icqtypes.Query{ChainId: "badchain"})
		// this should bail on a non-matching chain id.
		s.Require().Error(err)
	})
}

// func (s *KeeperTestSuite) TestHandleValidatorCallbackEmptyValue() {
// 	s.Run("empty value", func() {
// 		s.SetupTest()
// 		s.SetupZones()

// 		app := s.GetQuicksilverApp(s.chainA)
// 		app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
// 		ctx := s.chainA.GetContext()

// 		query := stakingtypes.QueryValidatorResponse{Validator: stakingtypes.Validator{}}
// 		bz, err := app.AppCodec().Marshal(&query)
// 		s.Require().NoError(err)

// 		err = keeper.ValidatorCallback(app.InterchainstakingKeeper, ctx, bz, icqtypes.Query{ChainId: s.chainB.ChainID})
// 		// this should error on unmarshalling an empty slice, which is not a valid response here.
// 		s.Require().Error(err)
// 	})
// }

// func (s *KeeperTestSuite) TestHandleValidatorCallback() {
// 	newVal := utils.GenerateValAddressForTest()

// 	tests := []struct {
// 		name      string
// 		validator func(in stakingtypes.Validators) stakingtypes.QueryValidatorResponse
// 		expected  icstypes.Validator
// 	}{
// 		{
// 			name: "valid - no-op",
// 			validator: func(in stakingtypes.Validators) stakingtypes.QueryValidatorResponse {
// 				return stakingtypes.QueryValidatorResponse{Validator: in[0]}
// 			},
// 			expected: icstypes.Validator{},
// 		},
// 		{
// 			name: "valid - shares +1000 val[0]",
// 			validator: func(in stakingtypes.Validators) stakingtypes.QueryValidatorResponse {
// 				in[0].DelegatorShares = in[0].DelegatorShares.Add(sdk.NewDec(1000))
// 				return stakingtypes.QueryValidatorResponse{Validator: in[0]}
// 			},
// 			expected: icstypes.Validator{},
// 		},
// 		{
// 			name: "valid - shares +1000 val[1], +2000 val[2]",
// 			validator: func(in stakingtypes.Validators) stakingtypes.QueryValidatorResponse {
// 				in[1].DelegatorShares = in[1].DelegatorShares.Add(sdk.NewDec(1000))
// 				in[2].DelegatorShares = in[2].DelegatorShares.Add(sdk.NewDec(2000))
// 				return stakingtypes.QueryValidatorResponse{Validator: in[0]}
// 			},
// 			expected: icstypes.Validator{},
// 		},
// 		{
// 			name: "valid - tokens +1000 val[0]",
// 			validator: func(in stakingtypes.Validators) stakingtypes.QueryValidatorResponse {
// 				in[0].Tokens = in[0].Tokens.Add(sdk.NewInt(1000))
// 				return stakingtypes.QueryValidatorResponse{Validator: in[0]}
// 			},
// 			expected: icstypes.Validator{},
// 		},
// 		{
// 			name: "valid - tokens +1000 val[1], +2000 val[2]",
// 			validator: func(in stakingtypes.Validators) stakingtypes.QueryValidatorResponse {
// 				in[1].Tokens = in[1].Tokens.Add(sdk.NewInt(1000))
// 				in[2].Tokens = in[2].Tokens.Add(sdk.NewInt(2000))
// 				return stakingtypes.QueryValidatorResponse{Validator: in[0]}
// 			},
// 			expected: icstypes.Validator{},
// 		},
// 		{
// 			name: "valid - tokens -10 val[1], -20 val[2]",
// 			validator: func(in stakingtypes.Validators) stakingtypes.QueryValidatorResponse {
// 				in[1].Tokens = in[1].Tokens.Sub(sdk.NewInt(10))
// 				in[2].Tokens = in[2].Tokens.Sub(sdk.NewInt(20))
// 				return stakingtypes.QueryValidatorResponse{Validator: in[0]}
// 			},
// 			expected: icstypes.Validator{},
// 		},
// 		{
// 			name: "valid - commission 0.5 val[0], 0.05 val[2]",
// 			validator: func(in stakingtypes.Validators) stakingtypes.QueryValidatorResponse {
// 				in[0].Commission.CommissionRates.Rate = sdk.NewDecWithPrec(5, 1)
// 				in[2].Commission.CommissionRates.Rate = sdk.NewDecWithPrec(5, 2)
// 				return stakingtypes.QueryValidatorResponse{Validator: in[0]}
// 			},
// 			expected: icstypes.Validator{},
// 		},
// 		{
// 			name: "valid - new validator",
// 			validator: func(in stakingtypes.Validators) stakingtypes.QueryValidatorResponse {
// 				new := in[0]
// 				new.OperatorAddress = newVal.String()
// 				in = append(in, new)
// 				return stakingtypes.QueryValidatorResponse{Validator: in[0]}
// 			},
// 			expected: icstypes.Validator{},
// 		},
// 	}

// 	for _, test := range tests {
// 		s.Run(test.name, func() {
// 			s.SetupTest()
// 			s.SetupZones()

// 			app := s.GetQuicksilverApp(s.chainA)
// 			app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
// 			ctx := s.chainA.GetContext()

// 			chainBVals := s.GetQuicksilverApp(s.chainB).StakingKeeper.GetValidators(s.chainB.GetContext(), 300)

// 			query := test.validator(chainBVals)
// 			bz, err := app.AppCodec().Marshal(&query)
// 			s.Require().NoError(err)

// 			err = keeper.ValsetCallback(app.InterchainstakingKeeper, ctx, bz, icqtypes.Query{ChainId: s.chainB.ChainID})
// 			s.Require().NoError(err)

// 			zone, found := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
// 			s.Require().True(found)
// 			// valset callback doesn't actually update validators, but does emit icq callbacks.
// 			expected := false
// 			for _, val := range zone.Validators {
// 				if val.IsEqual(test.expected) {
// 					expected = true
// 					break
// 				}
// 			}
// 			s.Require().True(expected)
// 		})
// 	}
// }

// func (s *KeeperTestSuite) TestHandleDelegationCallback() {
// 	type TestCase struct {
// 		name     string
// 		setup    func(vals []*types.Validator) []types.Delegation
// 		callback func(vals []*types.Validator) stakingtypes.Delegation
// 		expected func(vals []*types.Validator) types.Delegation
// 	}

// 	tests := []TestCase{
// 		func() TestCase {
// 			d1 := utils.GenerateValAddressForTest()
// 			return TestCase{
// 				name: "valid - no-op",
// 				setup: func(vals []*types.Validator) []types.Delegation {
// 					return []types.Delegation{
// 						{DelegationAddress: d1.String(), ValidatorAddress: vals[0].ValoperAddress, Amount: sdk.NewCoin("uatom", sdk.NewInt(5000000))},
// 						{DelegationAddress: d1.String(), ValidatorAddress: vals[1].ValoperAddress, Amount: sdk.NewCoin("raa", sdk.NewInt(2000000))},
// 					}
// 				},
// 				callback: func(vals []*types.Validator) stakingtypes.Delegation {
// 					return stakingtypes.Delegation{DelegatorAddress: d1.String(), ValidatorAddress: vals[0].ValoperAddress, Shares: sdk.NewDec(1000)}
// 				},
// 				expected: func(vals []*types.Validator) types.Delegation {
// 					return types.Delegation{DelegationAddress: d1.String(), ValidatorAddress: vals[0].ValoperAddress}
// 				},
// 			}
// 		}(),
// 	}

// 	for _, test := range tests {
// 		s.Run(test.name, func() {
// 			s.SetupTest()
// 			s.SetupZones()

// 			app := s.GetQuicksilverApp(s.chainA)
// 			app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
// 			ctx := s.chainA.GetContext()

// 			zone, found := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
// 			s.Require().True(found)

// 			for _, dg := range test.setup(zone.Validators) {
// 				app.InterchainstakingKeeper.SetDelegation(ctx, &zone, dg)
// 			}

// 			payload := test.callback(zone.Validators)
// 			bz, err := app.AppCodec().Marshal(&payload)
// 			s.Require().NoError(err)

// 			err = keeper.DelegationCallback(app.InterchainstakingKeeper, ctx, bz, icqtypes.Query{ChainId: s.chainB.ChainID})
// 			s.Require().NoError(err)

// 			expected := test.expected(zone.Validators)
// 			fmt.Println(app.InterchainstakingKeeper.GetAllDelegations(ctx, &zone))
// 			_, found = app.InterchainstakingKeeper.GetDelegation(ctx, &zone, expected.DelegationAddress, expected.ValidatorAddress)
// 			s.Require().True(found)
// 		})
// 	}
// }

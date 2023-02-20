package keeper_test

import (
	"bytes"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ingenuity-build/quicksilver/app"
	"github.com/ingenuity-build/quicksilver/utils"
	icqtypes "github.com/ingenuity-build/quicksilver/x/interchainquery/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	"github.com/stretchr/testify/require"
)

func TestCoinFromRequestKey(t *testing.T) {
	accAddr := utils.GenerateAccAddressForTest()
	prefix := banktypes.CreateAccountBalancesPrefix(accAddr.Bytes())
	query := append(prefix, []byte("denom")...)

	coin, err := utils.CoinFromRequestKey(query, accAddr)
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
			s.setupTestZones()

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
		s.setupTestZones()

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
		s.setupTestZones()

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
		s.setupTestZones()

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
		s.setupTestZones()

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

func (s *KeeperTestSuite) TestHandleValidatorCallbackNilValue() {
	s.Run("empty value", func() {
		s.SetupTest()
		s.setupTestZones()

		app := s.GetQuicksilverApp(s.chainA)
		app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := s.chainA.GetContext()

		bz := []byte{}

		err := keeper.ValidatorCallback(app.InterchainstakingKeeper, ctx, bz, icqtypes.Query{ChainId: s.chainB.ChainID})
		// this should error on unmarshalling an empty slice, which is not a valid response here.
		s.Require().Error(err)
	})
}

func (s *KeeperTestSuite) TestHandleValidatorCallback() {
	newVal := utils.GenerateAccAddressForTestWithPrefix("cosmosvaloper")

	zone := icstypes.Zone{ConnectionId: "connection-0", ChainId: "cosmoshub-4", AccountPrefix: "cosmos", LocalDenom: "uqatom", BaseDenom: "uatom"}
	zone.Validators = append(zone.Validators, &icstypes.Validator{ValoperAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), DelegatorShares: sdk.NewDec(2000)})
	zone.Validators = append(zone.Validators, &icstypes.Validator{ValoperAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), DelegatorShares: sdk.NewDec(2000)})
	zone.Validators = append(zone.Validators, &icstypes.Validator{ValoperAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), DelegatorShares: sdk.NewDec(2000)})
	zone.Validators = append(zone.Validators, &icstypes.Validator{ValoperAddress: "cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), DelegatorShares: sdk.NewDec(2000)})
	zone.Validators = append(zone.Validators, &icstypes.Validator{ValoperAddress: "cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), DelegatorShares: sdk.NewDec(2000)})

	tests := []struct {
		name      string
		validator stakingtypes.Validator
		expected  *icstypes.Validator
	}{
		{
			name:      "valid - no-op",
			validator: stakingtypes.Validator{OperatorAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", Jailed: false, Status: stakingtypes.Bonded, Tokens: sdk.NewInt(2000), DelegatorShares: sdk.NewDec(2000), Commission: stakingtypes.NewCommission(sdk.MustNewDecFromStr("0.2"), sdk.MustNewDecFromStr("0.2"), sdk.MustNewDecFromStr("0.2"))},
			expected:  &icstypes.Validator{ValoperAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), DelegatorShares: sdk.NewDec(2000), Score: sdk.ZeroDec(), Status: "BOND_STATUS_BONDED"},
		},
		{
			name:      "valid - +2000 tokens/shares",
			validator: stakingtypes.Validator{OperatorAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", Jailed: false, Status: stakingtypes.Bonded, Tokens: sdk.NewInt(4000), DelegatorShares: sdk.NewDec(4000), Commission: stakingtypes.NewCommission(sdk.MustNewDecFromStr("0.2"), sdk.MustNewDecFromStr("0.2"), sdk.MustNewDecFromStr("0.2"))},
			expected:  &icstypes.Validator{ValoperAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(4000), DelegatorShares: sdk.NewDec(4000), Score: sdk.ZeroDec(), Status: "BOND_STATUS_BONDED"},
		},
		{
			name:      "valid - inc. commission",
			validator: stakingtypes.Validator{OperatorAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", Jailed: false, Status: stakingtypes.Bonded, Tokens: sdk.NewInt(2000), DelegatorShares: sdk.NewDec(2000), Commission: stakingtypes.NewCommission(sdk.MustNewDecFromStr("0.5"), sdk.MustNewDecFromStr("0.2"), sdk.MustNewDecFromStr("0.2"))},
			expected:  &icstypes.Validator{ValoperAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", CommissionRate: sdk.MustNewDecFromStr("0.5"), VotingPower: sdk.NewInt(2000), DelegatorShares: sdk.NewDec(2000), Score: sdk.ZeroDec(), Status: "BOND_STATUS_BONDED"},
		},
		{
			name:      "valid - new validator",
			validator: stakingtypes.Validator{OperatorAddress: newVal, Jailed: false, Status: stakingtypes.Bonded, Tokens: sdk.NewInt(3000), DelegatorShares: sdk.NewDec(3050), Commission: stakingtypes.NewCommission(sdk.MustNewDecFromStr("0.25"), sdk.MustNewDecFromStr("0.2"), sdk.MustNewDecFromStr("0.2"))},
			expected:  &icstypes.Validator{ValoperAddress: newVal, CommissionRate: sdk.MustNewDecFromStr("0.25"), VotingPower: sdk.NewInt(3000), DelegatorShares: sdk.NewDec(3050), Score: sdk.ZeroDec(), Status: "BOND_STATUS_BONDED"},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			s.SetupTest()
			s.setupTestZones()

			app := s.GetQuicksilverApp(s.chainA)
			app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
			ctx := s.chainA.GetContext()

			app.InterchainstakingKeeper.SetZone(ctx, &zone)

			bz, err := app.AppCodec().Marshal(&test.validator)
			s.Require().NoError(err)

			err = keeper.ValidatorCallback(app.InterchainstakingKeeper, ctx, bz, icqtypes.Query{ChainId: zone.ChainId})
			s.Require().NoError(err)

			zone, found := app.InterchainstakingKeeper.GetZone(ctx, zone.ChainId)
			s.True(found)

			// valset callback doesn't actually update validators, but does emit icq callbacks.
			valFromZone, found := zone.GetValidatorByValoper(test.expected.ValoperAddress)
			s.True(found)
			s.Equal(test.expected, valFromZone)
		})
	}
}

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

func (s *KeeperTestSuite) TestHandleRewardsCallbackBadChain() {
	s.Run("bad chain", func() {
		s.SetupTest()
		s.setupTestZones()

		app := s.GetQuicksilverApp(s.chainA)
		app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := s.chainA.GetContext()

		query := distrtypes.QueryDelegationTotalRewardsResponse{}
		bz, err := app.AppCodec().Marshal(&query)
		s.Require().NoError(err)

		err = keeper.RewardsCallback(app.InterchainstakingKeeper, ctx, bz, icqtypes.Query{ChainId: "badchain"})
		// this should bail on a non-matching chain id.
		s.Require().Error(err)
	})
}

func (s *KeeperTestSuite) TestHandleRewardsEmptyRequestCallback() {
	s.Run("empty request", func() {
		s.SetupTest()
		s.setupTestZones()

		app := s.GetQuicksilverApp(s.chainA)
		app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := s.chainA.GetContext()

		query := distrtypes.QueryDelegationTotalRewardsRequest{}
		bz, err := app.AppCodec().Marshal(&query)
		s.Require().NoError(err)

		err = keeper.RewardsCallback(app.InterchainstakingKeeper, ctx, bz, icqtypes.Query{ChainId: s.chainB.ChainID})
		// this should fail because the waitgroup becomes negative.
		s.Require().Errorf(err, "attempted to unmarshal zero length byte slice (2)")
	})
}

func (s *KeeperTestSuite) TestHandleRewardsCallbackNonDelegator() {
	s.Run("valid response, bad delegator", func() {
		s.SetupTest()
		s.setupTestZones()

		app := s.GetQuicksilverApp(s.chainA)
		app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := s.chainA.GetContext()

		zone, _ := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
		zone.WithdrawalWaitgroup++
		app.InterchainstakingKeeper.SetZone(ctx, &zone)

		user := utils.GenerateAccAddressForTest()

		query := distrtypes.QueryDelegationTotalRewardsRequest{
			DelegatorAddress: user.String(),
		}

		response := distrtypes.QueryDelegationTotalRewardsResponse{
			Rewards: []distrtypes.DelegationDelegatorReward{
				{ValidatorAddress: s.chainB.Vals.Validators[0].String(), Reward: sdk.NewDecCoins(sdk.NewDecCoin("uatom", sdk.NewInt((1000))))},
			},
			Total: sdk.NewDecCoins(sdk.NewDecCoin("uatom", sdk.NewInt((1000)))),
		}
		reqbz, err := app.AppCodec().Marshal(&query)
		s.Require().NoError(err)

		respbz, err := app.AppCodec().Marshal(&response)
		s.Require().NoError(err)
		err = keeper.RewardsCallback(app.InterchainstakingKeeper, ctx, respbz, icqtypes.Query{ChainId: s.chainB.ChainID, Request: reqbz})
		//
		s.Require().Errorf(err, "failed attempting to withdraw rewards from non-delegation account")
	})
}

func (s *KeeperTestSuite) TestHandleRewardsCallbackEmptyResponse() {
	s.Run("empty response", func() {
		s.SetupTest()
		s.setupTestZones()

		app := s.GetQuicksilverApp(s.chainA)
		app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := s.chainA.GetContext()

		zone, _ := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
		zone.WithdrawalWaitgroup++
		app.InterchainstakingKeeper.SetZone(ctx, &zone)

		query := distrtypes.QueryDelegationTotalRewardsRequest{
			DelegatorAddress: zone.DelegationAddress.Address,
		}

		response := distrtypes.QueryDelegationTotalRewardsResponse{}
		reqbz, err := app.AppCodec().Marshal(&query)
		s.Require().NoError(err)

		respbz, err := app.AppCodec().Marshal(&response)
		s.Require().NoError(err)
		err = keeper.RewardsCallback(app.InterchainstakingKeeper, ctx, respbz, icqtypes.Query{ChainId: s.chainB.ChainID, Request: reqbz})
		//
		s.Require().NoError(err)
	})
}

func (s *KeeperTestSuite) TestHandleValideRewardsCallback() {
	s.Run("valid response, negative waitgroup", func() {
		s.SetupTest()
		s.setupTestZones()

		app := s.GetQuicksilverApp(s.chainA)
		app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := s.chainA.GetContext()

		zone, _ := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
		zone.WithdrawalWaitgroup++
		app.InterchainstakingKeeper.SetZone(ctx, &zone)

		query := distrtypes.QueryDelegationTotalRewardsRequest{
			DelegatorAddress: zone.DelegationAddress.Address,
		}

		response := distrtypes.QueryDelegationTotalRewardsResponse{
			Rewards: []distrtypes.DelegationDelegatorReward{
				{ValidatorAddress: s.chainB.Vals.Validators[0].String(), Reward: sdk.NewDecCoins(sdk.NewDecCoin("uatom", sdk.NewInt((1000))))},
			},
			Total: sdk.NewDecCoins(sdk.NewDecCoin("uatom", sdk.NewInt((1000)))),
		}
		reqbz, err := app.AppCodec().Marshal(&query)
		s.Require().NoError(err)

		respbz, err := app.AppCodec().Marshal(&response)
		s.Require().NoError(err)
		err = keeper.RewardsCallback(app.InterchainstakingKeeper, ctx, respbz, icqtypes.Query{ChainId: s.chainB.ChainID, Request: reqbz})
		//
		s.Require().NoError(err)
	})
}

func (s *KeeperTestSuite) TestAllBalancesCallback() {
	s.Run("all balances non-zero)", func() {
		s.SetupTest()
		s.setupTestZones()

		app := s.GetQuicksilverApp(s.chainA)
		app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := s.chainA.GetContext()

		zone, _ := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)

		query := banktypes.QueryAllBalancesRequest{
			Address: zone.DepositAddress.Address,
		}
		reqbz, err := app.AppCodec().Marshal(&query)
		s.Require().NoError(err)

		response := banktypes.QueryAllBalancesResponse{Balances: sdk.NewCoins(sdk.NewCoin("uqck", sdk.OneInt()))}
		respbz, err := app.AppCodec().Marshal(&response)
		s.Require().NoError(err)

		err = keeper.AllBalancesCallback(app.InterchainstakingKeeper, ctx, respbz, icqtypes.Query{ChainId: s.chainB.ChainID, Request: reqbz})
		s.Require().NoError(err)

		// refetch zone
		zone, _ = app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
		s.Require().Equal(uint32(1), zone.DepositAddress.BalanceWaitgroup)

		_, addr, err := bech32.DecodeAndConvert(zone.DepositAddress.Address)
		s.Require().NoError(err)
		data := banktypes.CreateAccountBalancesPrefix(addr)

		// check a ICQ request was made
		found := false
		app.InterchainQueryKeeper.IterateQueries(ctx, func(index int64, queryInfo icqtypes.Query) (stop bool) {
			if queryInfo.ChainId == zone.ChainId &&
				queryInfo.ConnectionId == zone.ConnectionId &&
				queryInfo.QueryType == icstypes.BankStoreKey &&
				bytes.Equal(queryInfo.Request, append(data, []byte(response.Balances[0].GetDenom())...)) {
				found = true
				return true
			}
			return false
		})
		s.Require().True(found)
	})
}

func (s *KeeperTestSuite) TestAllBalancesCallbackWithExistingWg() {
	s.Run("all balances non-zero)", func() {
		s.SetupTest()
		s.setupTestZones()

		app := s.GetQuicksilverApp(s.chainA)
		app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := s.chainA.GetContext()

		zone, _ := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
		zone.DepositAddress.BalanceWaitgroup = 2
		app.InterchainstakingKeeper.SetZone(ctx, &zone)

		query := banktypes.QueryAllBalancesRequest{
			Address: zone.DepositAddress.Address,
		}
		reqbz, err := app.AppCodec().Marshal(&query)
		s.Require().NoError(err)

		response := banktypes.QueryAllBalancesResponse{Balances: sdk.NewCoins(sdk.NewCoin("uqck", sdk.OneInt()))}
		respbz, err := app.AppCodec().Marshal(&response)
		s.Require().NoError(err)

		err = keeper.AllBalancesCallback(app.InterchainstakingKeeper, ctx, respbz, icqtypes.Query{ChainId: s.chainB.ChainID, Request: reqbz})
		s.Require().NoError(err)

		// refetch zone
		zone, _ = app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
		s.Require().Equal(uint32(1), zone.DepositAddress.BalanceWaitgroup)

		_, addr, err := bech32.DecodeAndConvert(zone.DepositAddress.Address)
		s.Require().NoError(err)
		data := banktypes.CreateAccountBalancesPrefix(addr)

		// check a ICQ request was made
		found := false
		app.InterchainQueryKeeper.IterateQueries(ctx, func(index int64, queryInfo icqtypes.Query) (stop bool) {
			if queryInfo.ChainId == zone.ChainId &&
				queryInfo.ConnectionId == zone.ConnectionId &&
				queryInfo.QueryType == icstypes.BankStoreKey &&
				bytes.Equal(queryInfo.Request, append(data, []byte(response.Balances[0].GetDenom())...)) {
				found = true
				return true
			}
			return false
		})
		s.Require().True(found)
	})
}

// tests where we have an existing balance and that balance is now reported as zero.
// we expect that an icq query will be emitted to assert with proof that the balance
// is now zero.
func (s *KeeperTestSuite) TestAllBalancesCallbackExistingBalanceNowNil() {
	s.Run("existing balance - now zero - deposit", func() {
		s.SetupTest()
		s.setupTestZones()

		app := s.GetQuicksilverApp(s.chainA)
		app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := s.chainA.GetContext()

		zone, _ := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
		zone.DepositAddress.Balance = sdk.NewCoins(sdk.NewCoin("uqck", sdk.OneInt()))
		app.InterchainstakingKeeper.SetZone(ctx, &zone)

		query := banktypes.QueryAllBalancesRequest{
			Address: zone.DepositAddress.Address,
		}
		reqbz, err := app.AppCodec().Marshal(&query)
		s.Require().NoError(err)

		response := banktypes.QueryAllBalancesResponse{Balances: sdk.Coins{}}
		respbz, err := app.AppCodec().Marshal(&response)
		s.Require().NoError(err)

		err = keeper.AllBalancesCallback(app.InterchainstakingKeeper, ctx, respbz, icqtypes.Query{ChainId: s.chainB.ChainID, Request: reqbz})
		s.Require().NoError(err)

		// refetch zone
		zone, _ = app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
		s.Require().Equal(uint32(1), zone.DepositAddress.BalanceWaitgroup)

		_, addr, err := bech32.DecodeAndConvert(zone.DepositAddress.Address)
		s.Require().NoError(err)
		data := banktypes.CreateAccountBalancesPrefix(addr)

		// check a ICQ request was made
		found := false
		app.InterchainQueryKeeper.IterateQueries(ctx, func(index int64, queryInfo icqtypes.Query) (stop bool) {
			if queryInfo.ChainId == zone.ChainId &&
				queryInfo.ConnectionId == zone.ConnectionId &&
				queryInfo.QueryType == icstypes.BankStoreKey &&
				bytes.Equal(queryInfo.Request, append(data, []byte("uqck")...)) {
				found = true
				return true
			}
			return false
		})
		s.Require().True(found)
	})

	s.Run("existing balance - now zero - withdrawal", func() {
		s.SetupTest()
		s.setupTestZones()

		app := s.GetQuicksilverApp(s.chainA)
		app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := s.chainA.GetContext()

		zone, _ := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
		zone.WithdrawalAddress.Balance = sdk.NewCoins(sdk.NewCoin("uqck", sdk.OneInt()))
		app.InterchainstakingKeeper.SetZone(ctx, &zone)

		query := banktypes.QueryAllBalancesRequest{
			Address: zone.WithdrawalAddress.Address,
		}
		reqbz, err := app.AppCodec().Marshal(&query)
		s.Require().NoError(err)

		response := banktypes.QueryAllBalancesResponse{Balances: sdk.Coins{}}
		respbz, err := app.AppCodec().Marshal(&response)
		s.Require().NoError(err)

		err = keeper.AllBalancesCallback(app.InterchainstakingKeeper, ctx, respbz, icqtypes.Query{ChainId: s.chainB.ChainID, Request: reqbz})
		s.Require().NoError(err)

		// refetch zone
		zone, _ = app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
		s.Require().Equal(uint32(1), zone.WithdrawalAddress.BalanceWaitgroup)

		_, addr, err := bech32.DecodeAndConvert(zone.WithdrawalAddress.Address)
		s.Require().NoError(err)
		data := banktypes.CreateAccountBalancesPrefix(addr)

		// check a ICQ request was made
		found := false
		app.InterchainQueryKeeper.IterateQueries(ctx, func(index int64, queryInfo icqtypes.Query) (stop bool) {
			if queryInfo.ChainId == zone.ChainId &&
				queryInfo.ConnectionId == zone.ConnectionId &&
				queryInfo.QueryType == icstypes.BankStoreKey &&
				bytes.Equal(queryInfo.Request, append(data, []byte("uqck")...)) {
				found = true
				return true
			}
			return false
		})
		s.Require().True(found)
	})
}

func (s *KeeperTestSuite) TestAllBalancesCallbackMulti() {
	s.Run("all balances non-zero)", func() {
		s.SetupTest()
		s.setupTestZones()

		app := s.GetQuicksilverApp(s.chainA)
		app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := s.chainA.GetContext()

		zone, _ := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)

		query := banktypes.QueryAllBalancesRequest{
			Address: zone.DepositAddress.Address,
		}
		reqbz, err := app.AppCodec().Marshal(&query)
		s.Require().NoError(err)

		response := banktypes.QueryAllBalancesResponse{Balances: sdk.NewCoins(sdk.NewCoin("uqck", sdk.OneInt()), sdk.NewCoin("stake", sdk.OneInt()))}
		respbz, err := app.AppCodec().Marshal(&response)
		s.Require().NoError(err)

		err = keeper.AllBalancesCallback(app.InterchainstakingKeeper, ctx, respbz, icqtypes.Query{ChainId: s.chainB.ChainID, Request: reqbz})
		s.Require().NoError(err)

		// refetch zone
		zone, _ = app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
		s.Require().Equal(uint32(2), zone.DepositAddress.BalanceWaitgroup)

		_, addr, err := bech32.DecodeAndConvert(zone.DepositAddress.Address)
		s.Require().NoError(err)
		data := banktypes.CreateAccountBalancesPrefix(addr)

		// check a ICQ request was made for each denom
		for _, coin := range response.Balances {
			found := false
			app.InterchainQueryKeeper.IterateQueries(ctx, func(index int64, queryInfo icqtypes.Query) (stop bool) {
				if queryInfo.ChainId == zone.ChainId &&
					queryInfo.ConnectionId == zone.ConnectionId &&
					queryInfo.QueryType == icstypes.BankStoreKey &&
					bytes.Equal(queryInfo.Request, append(data, []byte(coin.GetDenom())...)) {
					found = true
					return true
				}
				return false
			})
			s.Require().True(found)
		}
	})
}

func (s *KeeperTestSuite) TestAccountBalanceCallback() {
	s.Run("account balance", func() {
		s.SetupTest()
		s.setupTestZones()

		app := s.GetQuicksilverApp(s.chainA)
		app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := s.chainA.GetContext()

		zone, _ := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
		zone.DepositAddress.BalanceWaitgroup++
		zone.WithdrawalAddress.BalanceWaitgroup++
		app.InterchainstakingKeeper.SetZone(ctx, &zone)

		response := sdk.NewCoin("qck", sdk.NewInt(10))
		respbz, err := app.AppCodec().Marshal(&response)
		s.Require().NoError(err)

		for _, addr := range []string{zone.DepositAddress.Address, zone.WithdrawalAddress.Address} {
			accAddr, err := sdk.AccAddressFromBech32(addr)
			s.Require().NoError(err)
			data := append(banktypes.CreateAccountBalancesPrefix(accAddr), []byte("qck")...)

			err = keeper.AccountBalanceCallback(app.InterchainstakingKeeper, ctx, respbz, icqtypes.Query{ChainId: s.chainB.ChainID, Request: data})
			s.Require().NoError(err)
		}
	})
}

func (s *KeeperTestSuite) TestAccountBalanceCallbackMismatch() {
	s.Run("account balance", func() {
		s.SetupTest()
		s.setupTestZones()

		app := s.GetQuicksilverApp(s.chainA)
		app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := s.chainA.GetContext()

		zone, _ := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
		zone.DepositAddress.BalanceWaitgroup++
		zone.WithdrawalAddress.BalanceWaitgroup++
		app.InterchainstakingKeeper.SetZone(ctx, &zone)

		response := sdk.NewCoin("qck", sdk.NewInt(10))
		respbz, err := app.AppCodec().Marshal(&response)
		s.Require().NoError(err)

		for _, addr := range []string{zone.DepositAddress.Address, zone.WithdrawalAddress.Address} {
			accAddr, err := sdk.AccAddressFromBech32(addr)
			s.Require().NoError(err)
			data := append(banktypes.CreateAccountBalancesPrefix(accAddr), []byte("stake")...)

			err = keeper.AccountBalanceCallback(app.InterchainstakingKeeper, ctx, respbz, icqtypes.Query{ChainId: s.chainB.ChainID, Request: data})
			s.Require().ErrorContains(err, "received coin denom qck does not match requested denom stake")
		}
	})
}

func (s *KeeperTestSuite) TestAccountBalance046Callback() {
	s.Run("account balance", func() {
		s.SetupTest()
		s.setupTestZones()

		app := s.GetQuicksilverApp(s.chainA)
		app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := s.chainA.GetContext()

		zone, _ := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
		zone.DepositAddress.IncrementBalanceWaitgroup()
		zone.WithdrawalAddress.IncrementBalanceWaitgroup()
		app.InterchainstakingKeeper.SetZone(ctx, &zone)

		response := sdk.NewInt(10)

		respbz, err := response.Marshal()
		s.Require().NoError(err)

		for _, addr := range []string{zone.DepositAddress.Address, zone.WithdrawalAddress.Address} {
			accAddr, err := sdk.AccAddressFromBech32(addr)
			s.Require().NoError(err)
			data := append(banktypes.CreateAccountBalancesPrefix(accAddr), []byte("qck")...)

			err = keeper.AccountBalanceCallback(app.InterchainstakingKeeper, ctx, respbz, icqtypes.Query{ChainId: s.chainB.ChainID, Request: data})
			s.Require().NoError(err)
		}
	})
}

func (s *KeeperTestSuite) TestAccountBalanceCallbackNil() {
	s.Run("account balance", func() {
		s.SetupTest()
		s.setupTestZones()

		app := s.GetQuicksilverApp(s.chainA)
		app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := s.chainA.GetContext()

		zone, _ := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
		zone.DepositAddress.BalanceWaitgroup++
		zone.WithdrawalAddress.BalanceWaitgroup++
		app.InterchainstakingKeeper.SetZone(ctx, &zone)

		var response *sdk.Coin = nil
		respbz, err := app.AppCodec().Marshal(response)
		s.Require().NoError(err)

		for _, addr := range []string{zone.DepositAddress.Address, zone.WithdrawalAddress.Address} {
			accAddr, err := sdk.AccAddressFromBech32(addr)
			s.Require().NoError(err)
			data := append(banktypes.CreateAccountBalancesPrefix(accAddr), []byte("stake")...)

			err = keeper.AccountBalanceCallback(app.InterchainstakingKeeper, ctx, respbz, icqtypes.Query{ChainId: s.chainB.ChainID, Request: data})
			s.Require().NoError(err)
		}
	})
}

// Ensures that a fuzz vector which resulted in a crash of ValidatorReq.Pagination crashing
// doesn't creep back up. Please see https://github.com/ingenuity-build/quicksilver-incognito/issues/82
func TestValsetCallbackNilValidatorReqPagination(t *testing.T) {
	s := new(KeeperTestSuite)
	s.SetT(t)
	s.SetupTest()
	s.setupTestZones()

	app := s.GetQuicksilverApp(s.chainA)
	app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
	ctx := s.chainA.GetContext()

	data := []byte("\x12\"\n 00000000000000000000000000000000")
	_ = keeper.ValsetCallback(app.InterchainstakingKeeper, ctx, data, icqtypes.Query{ChainId: s.chainB.ChainID})
}

func TestDelegationsCallbackAllPresentNoChange(t *testing.T) {
	s := new(KeeperTestSuite)
	s.SetT(t)
	s.SetupTest()
	s.setupTestZones()

	app := s.GetQuicksilverApp(s.chainA)
	app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
	ctx := s.chainA.GetContext()
	cdc := app.IBCKeeper.Codec()

	zone, found := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
	s.Require().True(found)

	vals := s.GetQuicksilverApp(s.chainB).StakingKeeper.GetAllValidators(s.chainB.GetContext())
	delegationA := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[0].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationB := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[1].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationC := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[2].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}

	app.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationA)
	app.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationB)
	app.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationC)

	response := stakingtypes.QueryDelegatorDelegationsResponse{DelegationResponses: []stakingtypes.DelegationResponse{
		{Delegation: stakingtypes.Delegation{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[0].OperatorAddress, Shares: sdk.NewDec(1000)}, Balance: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))},
		{Delegation: stakingtypes.Delegation{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[1].OperatorAddress, Shares: sdk.NewDec(1000)}, Balance: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))},
		{Delegation: stakingtypes.Delegation{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[2].OperatorAddress, Shares: sdk.NewDec(1000)}, Balance: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))},
	}}

	data := cdc.MustMarshal(&response)

	delegationQuery := stakingtypes.QueryDelegatorDelegationsRequest{DelegatorAddr: zone.DelegationAddress.Address, Pagination: &query.PageRequest{Limit: uint64(len(zone.Validators))}}
	bz := cdc.MustMarshal(&delegationQuery)

	err := keeper.DelegationsCallback(app.InterchainstakingKeeper, ctx, data, icqtypes.Query{ChainId: s.chainB.ChainID, Request: bz})

	s.Require().NoError(err)

	delegationRequests := 0
	for _, query := range app.InterchainQueryKeeper.AllQueries(ctx) {
		if query.CallbackId == "delegation" {
			delegationRequests++
		}
	}

	s.Require().Equal(0, delegationRequests)
	s.Require().Equal(3, len(app.InterchainstakingKeeper.GetAllDelegations(ctx, &zone)))
}

func TestDelegationsCallbackAllPresentOneChange(t *testing.T) {
	s := new(KeeperTestSuite)
	s.SetT(t)
	s.SetupTest()
	s.setupTestZones()

	app := s.GetQuicksilverApp(s.chainA)
	app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
	ctx := s.chainA.GetContext()
	cdc := app.IBCKeeper.Codec()

	zone, found := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
	s.Require().True(found)

	vals := s.GetQuicksilverApp(s.chainB).StakingKeeper.GetAllValidators(s.chainB.GetContext())
	delegationA := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[0].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationB := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[1].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationC := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[2].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}

	app.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationA)
	app.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationB)
	app.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationC)

	response := stakingtypes.QueryDelegatorDelegationsResponse{DelegationResponses: []stakingtypes.DelegationResponse{
		{Delegation: stakingtypes.Delegation{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[0].OperatorAddress, Shares: sdk.NewDec(1000)}, Balance: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))},
		{Delegation: stakingtypes.Delegation{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[1].OperatorAddress, Shares: sdk.NewDec(2000)}, Balance: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(2000))},
		{Delegation: stakingtypes.Delegation{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[2].OperatorAddress, Shares: sdk.NewDec(1000)}, Balance: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))},
	}}

	data := cdc.MustMarshal(&response)

	delegationQuery := stakingtypes.QueryDelegatorDelegationsRequest{DelegatorAddr: zone.DelegationAddress.Address, Pagination: &query.PageRequest{Limit: uint64(len(zone.Validators))}}
	bz := cdc.MustMarshal(&delegationQuery)

	err := keeper.DelegationsCallback(app.InterchainstakingKeeper, ctx, data, icqtypes.Query{ChainId: s.chainB.ChainID, Request: bz})

	s.Require().NoError(err)

	delegationRequests := 0
	for _, query := range app.InterchainQueryKeeper.AllQueries(ctx) {
		if query.CallbackId == "delegation" {
			delegationRequests++
		}
	}

	s.Require().Equal(1, delegationRequests)
	s.Require().Equal(3, len(app.InterchainstakingKeeper.GetAllDelegations(ctx, &zone)))
}

func TestDelegationsCallbackOneMissing(t *testing.T) {
	s := new(KeeperTestSuite)
	s.SetT(t)
	s.SetupTest()
	s.setupTestZones()

	app := s.GetQuicksilverApp(s.chainA)
	app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
	ctx := s.chainA.GetContext()
	cdc := app.IBCKeeper.Codec()

	zone, found := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
	s.Require().True(found)

	vals := s.GetQuicksilverApp(s.chainB).StakingKeeper.GetAllValidators(s.chainB.GetContext())
	delegationA := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[0].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationB := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[1].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationC := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[2].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}

	app.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationA)
	app.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationB)
	app.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationC)

	response := stakingtypes.QueryDelegatorDelegationsResponse{DelegationResponses: []stakingtypes.DelegationResponse{
		{Delegation: stakingtypes.Delegation{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[0].OperatorAddress, Shares: sdk.NewDec(1000)}, Balance: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))},
		{Delegation: stakingtypes.Delegation{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[1].OperatorAddress, Shares: sdk.NewDec(1000)}, Balance: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))},
	}}

	data := cdc.MustMarshal(&response)

	delegationQuery := stakingtypes.QueryDelegatorDelegationsRequest{DelegatorAddr: zone.DelegationAddress.Address, Pagination: &query.PageRequest{Limit: uint64(len(zone.Validators))}}
	bz := cdc.MustMarshal(&delegationQuery)

	err := keeper.DelegationsCallback(app.InterchainstakingKeeper, ctx, data, icqtypes.Query{ChainId: s.chainB.ChainID, Request: bz})

	s.Require().NoError(err)

	delegationRequests := 0
	for _, query := range app.InterchainQueryKeeper.AllQueries(ctx) {
		if query.CallbackId == "delegation" {
			delegationRequests++
		}
	}

	s.Require().Equal(1, delegationRequests)                                             // callback for 'missing' delegation.
	s.Require().Equal(3, len(app.InterchainstakingKeeper.GetAllDelegations(ctx, &zone))) // new delegation doesn't get removed until the callback.
}

func TestDelegationsCallbackOneAdditional(t *testing.T) {
	s := new(KeeperTestSuite)
	s.SetT(t)
	s.SetupTest()
	s.setupTestZones()

	app := s.GetQuicksilverApp(s.chainA)
	app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
	ctx := s.chainA.GetContext()
	cdc := app.IBCKeeper.Codec()

	zone, found := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
	s.Require().True(found)

	vals := s.GetQuicksilverApp(s.chainB).StakingKeeper.GetAllValidators(s.chainB.GetContext())
	delegationA := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[0].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationB := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[1].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationC := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[2].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}

	app.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationA)
	app.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationB)
	app.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationC)

	response := stakingtypes.QueryDelegatorDelegationsResponse{DelegationResponses: []stakingtypes.DelegationResponse{
		{Delegation: stakingtypes.Delegation{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[0].OperatorAddress, Shares: sdk.NewDec(1000)}, Balance: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))},
		{Delegation: stakingtypes.Delegation{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[1].OperatorAddress, Shares: sdk.NewDec(1000)}, Balance: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))},
		{Delegation: stakingtypes.Delegation{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[2].OperatorAddress, Shares: sdk.NewDec(1000)}, Balance: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))},
		{Delegation: stakingtypes.Delegation{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[3].OperatorAddress, Shares: sdk.NewDec(1000)}, Balance: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))},
	}}

	data := cdc.MustMarshal(&response)

	delegationQuery := stakingtypes.QueryDelegatorDelegationsRequest{DelegatorAddr: zone.DelegationAddress.Address, Pagination: &query.PageRequest{Limit: uint64(len(zone.Validators))}}
	bz := cdc.MustMarshal(&delegationQuery)

	err := keeper.DelegationsCallback(app.InterchainstakingKeeper, ctx, data, icqtypes.Query{ChainId: s.chainB.ChainID, Request: bz})

	s.Require().NoError(err)

	delegationRequests := 0
	for _, query := range app.InterchainQueryKeeper.AllQueries(ctx) {
		if query.CallbackId == "delegation" {
			delegationRequests++
		}
	}

	s.Require().Equal(1, delegationRequests)
	s.Require().Equal(3, len(app.InterchainstakingKeeper.GetAllDelegations(ctx, &zone))) // new delegation doesn't get added until the end
}

func TestDelegationCallbackNew(t *testing.T) {
	s := new(KeeperTestSuite)
	s.SetT(t)
	s.SetupTest()
	s.setupTestZones()

	app := s.GetQuicksilverApp(s.chainA)
	app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
	ctx := s.chainA.GetContext()
	cdc := app.IBCKeeper.Codec()

	zone, found := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
	s.Require().True(found)

	vals := s.GetQuicksilverApp(s.chainB).StakingKeeper.GetAllValidators(s.chainB.GetContext())
	delegationA := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[0].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationB := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[1].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationC := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[2].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}

	app.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationA)
	app.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationB)
	app.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationC)

	response := stakingtypes.Delegation{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[3].OperatorAddress, Shares: sdk.NewDec(1000)}

	data := cdc.MustMarshal(&response)

	delAddr, err := utils.AccAddressFromBech32(zone.DelegationAddress.Address, "")
	s.Require().NoError(err)
	valAddr, err := utils.ValAddressFromBech32(vals[3].OperatorAddress, "")
	s.Require().NoError(err)
	bz := stakingtypes.GetDelegationKey(delAddr, valAddr)

	err = keeper.DelegationCallback(app.InterchainstakingKeeper, ctx, data, icqtypes.Query{ChainId: s.chainB.ChainID, Request: bz})
	s.Require().NoError(err)

	s.Require().Equal(4, len(app.InterchainstakingKeeper.GetAllDelegations(ctx, &zone)))
}

func TestDelegationCallbackUpdate(t *testing.T) {
	s := new(KeeperTestSuite)
	s.SetT(t)
	s.SetupTest()
	s.setupTestZones()

	app := s.GetQuicksilverApp(s.chainA)
	app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
	ctx := s.chainA.GetContext()
	cdc := app.IBCKeeper.Codec()

	zone, found := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
	s.Require().True(found)

	vals := s.GetQuicksilverApp(s.chainB).StakingKeeper.GetAllValidators(s.chainB.GetContext())
	delegationA := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[0].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationB := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[1].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationC := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[2].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}

	app.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationA)
	app.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationB)
	app.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationC)

	response := stakingtypes.Delegation{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[2].OperatorAddress, Shares: sdk.NewDec(2000)}

	data := cdc.MustMarshal(&response)

	delAddr, err := utils.AccAddressFromBech32(zone.DelegationAddress.Address, "")
	s.Require().NoError(err)
	valAddr, err := utils.ValAddressFromBech32(vals[3].OperatorAddress, "")
	s.Require().NoError(err)
	bz := stakingtypes.GetDelegationKey(delAddr, valAddr)

	err = keeper.DelegationCallback(app.InterchainstakingKeeper, ctx, data, icqtypes.Query{ChainId: s.chainB.ChainID, Request: bz})
	s.Require().NoError(err)

	s.Require().Equal(3, len(app.InterchainstakingKeeper.GetAllDelegations(ctx, &zone)))
}

func TestDelegationCallbackNoOp(t *testing.T) {
	s := new(KeeperTestSuite)
	s.SetT(t)
	s.SetupTest()
	s.setupTestZones()

	app := s.GetQuicksilverApp(s.chainA)
	app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
	ctx := s.chainA.GetContext()
	cdc := app.IBCKeeper.Codec()

	zone, found := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
	s.Require().True(found)

	vals := s.GetQuicksilverApp(s.chainB).StakingKeeper.GetAllValidators(s.chainB.GetContext())
	delegationA := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[0].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationB := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[1].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationC := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[2].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}

	app.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationA)
	app.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationB)
	app.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationC)

	response := stakingtypes.Delegation{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[2].OperatorAddress, Shares: sdk.NewDec(1000)}

	data := cdc.MustMarshal(&response)

	delAddr, err := utils.AccAddressFromBech32(zone.DelegationAddress.Address, "")
	s.Require().NoError(err)
	valAddr, err := utils.ValAddressFromBech32(vals[3].OperatorAddress, "")
	s.Require().NoError(err)
	bz := stakingtypes.GetDelegationKey(delAddr, valAddr)

	err = keeper.DelegationCallback(app.InterchainstakingKeeper, ctx, data, icqtypes.Query{ChainId: s.chainB.ChainID, Request: bz})
	s.Require().NoError(err)

	s.Require().Equal(3, len(app.InterchainstakingKeeper.GetAllDelegations(ctx, &zone)))
}

func TestDelegationCallbackRemove(t *testing.T) {
	s := new(KeeperTestSuite)
	s.SetT(t)
	s.SetupTest()
	s.setupTestZones()

	app := s.GetQuicksilverApp(s.chainA)
	app.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
	ctx := s.chainA.GetContext()
	cdc := app.IBCKeeper.Codec()

	zone, found := app.InterchainstakingKeeper.GetZone(ctx, s.chainB.ChainID)
	s.Require().True(found)

	vals := s.GetQuicksilverApp(s.chainB).StakingKeeper.GetAllValidators(s.chainB.GetContext())
	delegationA := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[0].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationB := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[1].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationC := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[2].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}

	app.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationA)
	app.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationB)
	app.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationC)

	response := stakingtypes.Delegation{}

	data := cdc.MustMarshal(&response)

	delAddr, err := utils.AccAddressFromBech32(zone.DelegationAddress.Address, "")
	s.Require().NoError(err)
	valAddr, err := utils.ValAddressFromBech32(vals[3].OperatorAddress, "")
	s.Require().NoError(err)
	bz := stakingtypes.GetDelegationKey(delAddr, valAddr)

	err = keeper.DelegationCallback(app.InterchainstakingKeeper, ctx, data, icqtypes.Query{ChainId: s.chainB.ChainID, Request: bz})
	s.Require().NoError(err)

	delegationRequests := 0
	for _, query := range app.InterchainQueryKeeper.AllQueries(ctx) {
		if query.CallbackId == "delegation" {
			delegationRequests++
		}
	}

	s.Require().Equal(3, len(app.InterchainstakingKeeper.GetAllDelegations(ctx, &zone)))
}

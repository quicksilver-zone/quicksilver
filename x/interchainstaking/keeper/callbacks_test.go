package keeper_test

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/math"
	ics23 "github.com/confio/ics23/go"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/types/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	lightclienttypes "github.com/cosmos/ibc-go/v5/modules/light-clients/07-tendermint/types"
	"github.com/ingenuity-build/quicksilver/app"
	"github.com/ingenuity-build/quicksilver/utils"
	icqtypes "github.com/ingenuity-build/quicksilver/x/interchainquery/types"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/keeper"
	icstypes "github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
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
			expected:  &icstypes.Validator{ValoperAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), DelegatorShares: sdk.NewDec(2000), Score: sdk.ZeroDec(), Status: "BOND_STATUS_BONDED", LiquidShares: sdk.ZeroDec(), ValidatorBondShares: sdk.ZeroDec()},
		},
		{
			name:      "valid - +2000 tokens/shares",
			validator: stakingtypes.Validator{OperatorAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", Jailed: false, Status: stakingtypes.Bonded, Tokens: sdk.NewInt(4000), DelegatorShares: sdk.NewDec(4000), Commission: stakingtypes.NewCommission(sdk.MustNewDecFromStr("0.2"), sdk.MustNewDecFromStr("0.2"), sdk.MustNewDecFromStr("0.2"))},
			expected:  &icstypes.Validator{ValoperAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(4000), DelegatorShares: sdk.NewDec(4000), Score: sdk.ZeroDec(), Status: "BOND_STATUS_BONDED", LiquidShares: sdk.ZeroDec(), ValidatorBondShares: sdk.ZeroDec()},
		},
		{
			name:      "valid - inc. commission",
			validator: stakingtypes.Validator{OperatorAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", Jailed: false, Status: stakingtypes.Bonded, Tokens: sdk.NewInt(2000), DelegatorShares: sdk.NewDec(2000), Commission: stakingtypes.NewCommission(sdk.MustNewDecFromStr("0.5"), sdk.MustNewDecFromStr("0.2"), sdk.MustNewDecFromStr("0.2"))},
			expected:  &icstypes.Validator{ValoperAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", CommissionRate: sdk.MustNewDecFromStr("0.5"), VotingPower: sdk.NewInt(2000), DelegatorShares: sdk.NewDec(2000), Score: sdk.ZeroDec(), Status: "BOND_STATUS_BONDED", LiquidShares: sdk.ZeroDec(), ValidatorBondShares: sdk.ZeroDec()},
		},
		{
			name:      "valid - new validator",
			validator: stakingtypes.Validator{OperatorAddress: newVal, Jailed: false, Status: stakingtypes.Bonded, Tokens: sdk.NewInt(3000), DelegatorShares: sdk.NewDec(3050), Commission: stakingtypes.NewCommission(sdk.MustNewDecFromStr("0.25"), sdk.MustNewDecFromStr("0.2"), sdk.MustNewDecFromStr("0.2"))},
			expected:  &icstypes.Validator{ValoperAddress: newVal, CommissionRate: sdk.MustNewDecFromStr("0.25"), VotingPower: sdk.NewInt(3000), DelegatorShares: sdk.NewDec(3050), Score: sdk.ZeroDec(), Status: "BOND_STATUS_BONDED", LiquidShares: sdk.ZeroDec(), ValidatorBondShares: sdk.ZeroDec()},
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

func decodeBase64NoErr(str string) []byte {
	decoded, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		fmt.Println("Error decoding: ", str)
		panic(err)
	}
	return decoded
}

func (suite *KeeperTestSuite) TestDepositLsmTxCallback() {
	suite.Run("Deposit transaction successful", func() {
		suite.SetupTest()
		suite.setupTestZones()

		// setup quicksilver test app
		quicksilver := suite.GetQuicksilverApp(suite.chainA)
		quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()

		// get chainA context
		ctx := suite.chainA.GetContext()

		// get zone chainB context
		zone, _ := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
		zone.DepositAddress.IncrementBalanceWaitgroup()
		zone.WithdrawalAddress.IncrementBalanceWaitgroup()
		// add the validator from the gaiatest-1 network to our registered zone. This is required for LSM deposit as the tokenised share denom is checked against known validators.
		zone.Validators = append(zone.Validators, &icstypes.Validator{
			ValoperAddress:      "cosmosvaloper1gg7w8w2y9jfv76a2yyahe42y09g9ry2raa5rqf",
			CommissionRate:      sdk.NewDecWithPrec(1, 1),
			DelegatorShares:     sdk.MustNewDecFromStr("4235376641.000000000000000000"),
			VotingPower:         sdk.NewInt(4235376641),
			Status:              "BOND_STATUS_BONDED",
			Jailed:              false,
			Tombstoned:          false,
			ValidatorBondShares: sdk.MustNewDecFromStr("1000000.000000000000000000"),
			LiquidShares:        sdk.MustNewDecFromStr("4234076641.000000000000000000"),
		})
		// override the DepositAddress to match that of the chain where the fixture was captured.
		zone.DepositAddress.Address = "cosmos1avvehf3npvn6weyxtvyu7mhwwvjryzw69g43tq0nl80wqjglr6hse5mcz4"
		quicksilver.InterchainstakingKeeper.SetZone(ctx, &zone)

		// create tx fixture - this was taken from a live v1.2 chain (lfg-1 <-> gaiatest-1) hence the need to override client and consensus states, to match the source.
		payload := icqtypes.GetTxWithProofResponse{}
		payloadBytes := decodeBase64NoErr(txFixtureLsm)

		err := quicksilver.InterchainstakingKeeper.GetCodec().Unmarshal(payloadBytes, &payload)
		// update payload header to ensure we can validate it.
		payload.Header.Header.Time = ctx.BlockTime()
		suite.NoError(err)
		// cheat, and set the client state and consensus state for 07-tendermint-0 to match the incoming header.
		quicksilver.IBCKeeper.ClientKeeper.SetClientState(ctx, "07-tendermint-0", lightclienttypes.NewClientState("gaiatest-1", lightclienttypes.DefaultTrustLevel, time.Hour, time.Hour, time.Second*50, payload.Header.TrustedHeight, []*ics23.ProofSpec{}, []string{}, false, false))
		quicksilver.IBCKeeper.ClientKeeper.SetClientConsensusState(ctx, "07-tendermint-0", payload.Header.TrustedHeight, payload.Header.ConsensusState())

		requestData := tx.GetTxRequest{
			// hash of tx in `txFixture`
			Hash: "b1f1852d322328f6b8d8cacd180df2b1cbbd3dd64536c9ecbf1c896a15f6217a",
		}
		// check receipt does not exist created.
		_, found := quicksilver.InterchainstakingKeeper.GetReceipt(ctx, keeper.GetReceiptKey(zone.ChainId, requestData.Hash))

		suite.False(found)

		resDataBz, err := quicksilver.AppCodec().Marshal(&requestData)
		suite.NoError(err)

		// trigger the callback
		err = keeper.DepositTx(quicksilver.InterchainstakingKeeper, ctx, payloadBytes, icqtypes.Query{ChainId: suite.chainB.ChainID, Request: resDataBz})

		suite.NoError(err)

		// expect quick1a2zht8x2j0dqvuejr8pxpu7due3qmk405vakg9 to have 5000 uqatoms now!
		addrBytes, _ := utils.AccAddressFromBech32("cosmos1a2zht8x2j0dqvuejr8pxpu7due3qmk40lgdy3h", "")
		newBalance := quicksilver.BankKeeper.GetAllBalances(ctx, addrBytes)
		suite.Equal(newBalance.AmountOf("uqatom"), math.NewInt(5000))

		// check receipt was created.
		receipt, found := quicksilver.InterchainstakingKeeper.GetReceipt(ctx, keeper.GetReceiptKey(zone.ChainId, requestData.Hash))

		suite.True(found)

		sdkTx, err := keeper.TxDecoder(quicksilver.InterchainstakingKeeper.GetCodec())(payload.Proof.Data)
		suite.NoError(err)

		authTx, _ := sdkTx.(*tx.Tx)

		// validate receipt matches source / hash / amount
		var msg sdk.Msg
		quicksilver.InterchainstakingKeeper.GetCodec().UnpackAny(authTx.Body.Messages[0], &msg)
		sendmsg, _ := msg.(*banktypes.MsgSend)
		suite.Equal(receipt.Sender, sendmsg.FromAddress)
		suite.Equal(receipt.Txhash, requestData.Hash)
		bt := ctx.BlockTime()
		suite.Equal(receipt.FirstSeen, &bt)
		suite.Equal(receipt.Amount, sendmsg.Amount)

		// resubmitting the tx should not fail - it is silently ignored as we have now seen it before - but check the recipient balance has not changed.
		err = keeper.DepositTx(quicksilver.InterchainstakingKeeper, ctx, payloadBytes, icqtypes.Query{ChainId: suite.chainB.ChainID, Request: resDataBz})

		suite.NoError(err)

		nowBalance := quicksilver.BankKeeper.GetAllBalances(ctx, addrBytes)
		suite.Equal(nowBalance.AmountOf("uqatom"), math.NewInt(5000))
	})
}

func (suite *KeeperTestSuite) TestDepositTxCallback() {
	suite.Run("Deposit transaction successful", func() {
		suite.SetupTest()
		suite.setupTestZones()

		// setup quicksilver test app
		quicksilver := suite.GetQuicksilverApp(suite.chainA)
		quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()

		// get chainA context
		ctx := suite.chainA.GetContext()

		// get zone chainB context
		zone, _ := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
		zone.DepositAddress.IncrementBalanceWaitgroup()
		zone.WithdrawalAddress.IncrementBalanceWaitgroup()
		// add the validator from the gaiatest-1 network to our registered zone. This is required for LSM deposit as the tokenised share denom is checked against known validators.
		zone.Validators = append(zone.Validators, &icstypes.Validator{
			ValoperAddress:      "cosmosvaloper1gg7w8w2y9jfv76a2yyahe42y09g9ry2raa5rqf",
			CommissionRate:      sdk.NewDecWithPrec(1, 1),
			DelegatorShares:     sdk.MustNewDecFromStr("4235376641.000000000000000000"),
			VotingPower:         sdk.NewInt(4235376641),
			Status:              "BOND_STATUS_BONDED",
			Jailed:              false,
			Tombstoned:          false,
			ValidatorBondShares: sdk.MustNewDecFromStr("1000000.000000000000000000"),
			LiquidShares:        sdk.MustNewDecFromStr("4234076641.000000000000000000"),
		})
		// override the DepositAddress to match that of the chain where the fixture was captured.
		zone.DepositAddress.Address = "cosmos1d2jrh4gj66smxns6xfv8mdd4keef5ek97knl9gw9skkzryjyhxjsywlfhm"
		quicksilver.InterchainstakingKeeper.SetZone(ctx, &zone)

		// create tx fixture - this was taken from a live v1.2 chain (lfg-1 <-> gaiatest-1) hence the need to override client and consensus states, to match the source.
		payload := icqtypes.GetTxWithProofResponse{}
		payloadBytes := decodeBase64NoErr(txFixture)

		err := quicksilver.InterchainstakingKeeper.GetCodec().Unmarshal(payloadBytes, &payload)
		// update payload header to ensure we can validate it.
		payload.Header.Header.Time = ctx.BlockTime()
		suite.NoError(err)
		// cheat, and set the client state and consensus state for 07-tendermint-0 to match the incoming header.
		quicksilver.IBCKeeper.ClientKeeper.SetClientState(ctx, "07-tendermint-0", lightclienttypes.NewClientState("gaiatest-1", lightclienttypes.DefaultTrustLevel, time.Hour, time.Hour, time.Second*50, payload.Header.TrustedHeight, []*ics23.ProofSpec{}, []string{}, false, false))
		quicksilver.IBCKeeper.ClientKeeper.SetClientConsensusState(ctx, "07-tendermint-0", payload.Header.TrustedHeight, payload.Header.ConsensusState())

		requestData := tx.GetTxRequest{
			// hash of tx in `txFixture`
			Hash: "fa7b199bfa3877d2f438ad7802a6b92cddde5e812f5620f1db735b7a90439938",
		}
		// check receipt does not exist created.
		_, found := quicksilver.InterchainstakingKeeper.GetReceipt(ctx, keeper.GetReceiptKey(zone.ChainId, requestData.Hash))

		suite.False(found)

		resDataBz, err := quicksilver.AppCodec().Marshal(&requestData)
		suite.NoError(err)

		// trigger the callback
		err = keeper.DepositTx(quicksilver.InterchainstakingKeeper, ctx, payloadBytes, icqtypes.Query{ChainId: suite.chainB.ChainID, Request: resDataBz})

		suite.NoError(err)

		// expect cosmos1a2zht8x2j0dqvuejr8pxpu7due3qmk40lgdy3h to have 500000 uqatoms now!
		addrBytes, _ := utils.AccAddressFromBech32("cosmos1a2zht8x2j0dqvuejr8pxpu7due3qmk40lgdy3h", "")
		newBalance := quicksilver.BankKeeper.GetAllBalances(ctx, addrBytes)
		suite.Equal(newBalance.AmountOf("uqatom"), math.NewInt(500000))

		// check receipt was created.
		receipt, found := quicksilver.InterchainstakingKeeper.GetReceipt(ctx, keeper.GetReceiptKey(zone.ChainId, requestData.Hash))

		suite.True(found)

		sdkTx, err := keeper.TxDecoder(quicksilver.InterchainstakingKeeper.GetCodec())(payload.Proof.Data)
		suite.NoError(err)

		authTx, _ := sdkTx.(*tx.Tx)

		// validate receipt matches source / hash / amount
		var msg sdk.Msg
		quicksilver.InterchainstakingKeeper.GetCodec().UnpackAny(authTx.Body.Messages[0], &msg)
		sendmsg, _ := msg.(*banktypes.MsgSend)
		suite.Equal(receipt.Sender, sendmsg.FromAddress)
		suite.Equal(receipt.Txhash, requestData.Hash)
		bt := ctx.BlockTime()
		suite.Equal(receipt.FirstSeen, &bt)
		suite.Equal(receipt.Amount, sendmsg.Amount)

		// resubmitting the tx should not fail - it is silently ignored as we have now seen it before - but check the recipient balance has not changed.
		err = keeper.DepositTx(quicksilver.InterchainstakingKeeper, ctx, payloadBytes, icqtypes.Query{ChainId: suite.chainB.ChainID, Request: resDataBz})

		suite.NoError(err)

		nowBalance := quicksilver.BankKeeper.GetAllBalances(ctx, addrBytes)
		suite.Equal(nowBalance.AmountOf("uqatom"), math.NewInt(500000))
	})
}

func (suite *KeeperTestSuite) TestDepositLsmTxCallbackFailOnNonMatchingValidator() {
	suite.Run("Deposit transaction successful", func() {
		suite.SetupTest()
		suite.setupTestZones()

		// setup quicksilver test app
		quicksilver := suite.GetQuicksilverApp(suite.chainA)
		quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()

		// get chainA context
		ctx := suite.chainA.GetContext()

		// get zone chainB context
		zone, _ := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
		zone.DepositAddress.IncrementBalanceWaitgroup()
		zone.WithdrawalAddress.IncrementBalanceWaitgroup()

		// override the DepositAddress to match that of the chain where the fixture was captured.
		zone.DepositAddress.Address = "cosmos1avvehf3npvn6weyxtvyu7mhwwvjryzw69g43tq0nl80wqjglr6hse5mcz4"
		quicksilver.InterchainstakingKeeper.SetZone(ctx, &zone)

		// create tx fixture - this was taken from a live v1.2 chain (lfg-1 <-> gaiatest-1) hence the need to override client and consensus states, to match the source.
		payload := icqtypes.GetTxWithProofResponse{}
		payloadBytes := decodeBase64NoErr(txFixtureLsm)

		err := quicksilver.InterchainstakingKeeper.GetCodec().Unmarshal(payloadBytes, &payload)
		// update payload header to ensure we can validate it.
		payload.Header.Header.Time = ctx.BlockTime()
		suite.NoError(err)
		// cheat, and set the client state and consensus state for 07-tendermint-0 to match the incoming header.
		quicksilver.IBCKeeper.ClientKeeper.SetClientState(ctx, "07-tendermint-0", lightclienttypes.NewClientState("gaiatest-1", lightclienttypes.DefaultTrustLevel, time.Hour, time.Hour, time.Second*50, payload.Header.TrustedHeight, []*ics23.ProofSpec{}, []string{}, false, false))
		quicksilver.IBCKeeper.ClientKeeper.SetClientConsensusState(ctx, "07-tendermint-0", payload.Header.TrustedHeight, payload.Header.ConsensusState())

		requestData := tx.GetTxRequest{
			// hash of tx in `txFixture`
			Hash: "b1f1852d322328f6b8d8cacd180df2b1cbbd3dd64536c9ecbf1c896a15f6217a",
		}
		// check receipt does not exist created.
		_, found := quicksilver.InterchainstakingKeeper.GetReceipt(ctx, keeper.GetReceiptKey(zone.ChainId, requestData.Hash))

		suite.False(found)

		resDataBz, err := quicksilver.AppCodec().Marshal(&requestData)
		suite.NoError(err)

		// trigger the callback
		err = keeper.DepositTx(quicksilver.InterchainstakingKeeper, ctx, payloadBytes, icqtypes.Query{ChainId: suite.chainB.ChainID, Request: resDataBz})

		suite.ErrorContains(err, "unable to validate coins. Ignoring.")

		// expect quick1a2zht8x2j0dqvuejr8pxpu7due3qmk405vakg9 to have 0 uqatoms, as the deposit failed.
		addrBytes, _ := utils.AccAddressFromBech32("cosmos1a2zht8x2j0dqvuejr8pxpu7due3qmk40lgdy3h", "")
		newBalance := quicksilver.BankKeeper.GetAllBalances(ctx, addrBytes)
		suite.Equal(newBalance.AmountOf("uqatom"), math.NewInt(0))

		_, found = quicksilver.InterchainstakingKeeper.GetReceipt(ctx, keeper.GetReceiptKey(zone.ChainId, requestData.Hash))

		suite.False(found)
	})
}

const (
	txFixtureLsm = "GsEDCiCFDDobCzFK2Vf0BXcgdEycLSdJL8IP7PEVWKelDQeJ3xL2AgrXAQrUAQocL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBKzAQotY29zbW9zMWEyemh0OHgyajBkcXZ1ZWpyOHB4cHU3ZHVlM3FtazQwbGdkeTNoEkFjb3Ntb3MxYXZ2ZWhmM25wdm42d2V5eHR2eXU3bWh3d3Zqcnl6dzY5ZzQzdHEwbmw4MHdxamdscjZoc2U1bWN6NBo/Cjdjb3Ntb3N2YWxvcGVyMWdnN3c4dzJ5OWpmdjc2YTJ5eWFoZTQyeTA5ZzlyeTJyYWE1cnFmLzE0EgQ1MDAwElgKUApGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQLaGco86x6BgxaGOBf/rgbHMEyZzECi+5in9DJ31ln/0BIECgIIARgoEgQQwJoMGkAtbKm5mTCs2SJzFZL5UKaFbKascEfSLtLFX4w9H/iLKXVqia/1REtynG8yLW374PPGFRplDo62C3SrhSBSLETgGiQIARoghQw6GwsxStlX9AV3IHRMnC0nSS/CD+zxFVinpQ0Hid8i2wYK0AQKkgMKAggLEgpnYWlhdGVzdC0xGJjdDiIMCPHf2agGEODypL8CKkgKIFvrRJTTqdEJ0eh/bm+bNFIMSX7ad1Uz9FX2u8acwNOAEiQIARIgqdknqwXY2NKl/r0A/JEd6hFCVr+E+xoDP5xqjTdMzkkyIFTqmUpOcyiALxE9GyyJ8B0qHyYAXdEyebrP+zlYCVe/OiCFDDobCzFK2Vf0BXcgdEycLSdJL8IP7PEVWKelDQeJ30IgLdUPCAh3Ii0/aGdGLRM24PsOqJJsvS6jPy3hstJUQ0RKIC3VDwgIdyItP2hnRi0TNuD7DqiSbL0uoz8t4bLSVENEUiAEgJG8fdwoP3e/v5HXPETaWMPfipy8hnQF2Lfz2q2iL1ogxyvS5b5sdsYoCMUEDDELSqvtajtVi8Tix+aShLESfBdiIOOwxEKY/BwUmvv0yJlvuSQnrkHkZJuTTKSVmRt4UrhVaiDjsMRCmPwcFJr79MiZb7kkJ65B5GSbk0yklZkbeFK4VXIUeQRs3t3nFppZq/OiJ+/f0AsW+twSuAEImN0OGkgKIOEOuKl3gvM4+gGbzlmy63IKY27HPnTJ6rszQyUuZAwPEiQIARIgiO0gt0gzcxTIEfFhpxf+XrKDoSnwZ9/HXl9XavCfS7UiaAgCEhR5BGze3ecWmlmr86In79/QCxb63BoMCPbf2agGEJiC0c0CIkBUNOaucBUZko0uikQApp2uWUJQ/zAtwTr5PRWlVS5/wFJYMGBSDNh5EEWY4FTclhTHLV2aMyyH5pfH6L0fr50CEn4KPQoUeQRs3t3nFppZq/OiJ+/f0AsW+twSIgogx4ew5LC25gOeUAdpun5LhBSfIBHUbK7Zjyzn8VRr1ZwYhCESPQoUeQRs3t3nFppZq/OiJ+/f0AsW+twSIgogx4ew5LC25gOeUAdpun5LhBSfIBHUbK7Zjyzn8VRr1ZwYhCEaBggBEKvdDiJ+Cj0KFHkEbN7d5xaaWavzoifv39ALFvrcEiIKIMeHsOSwtuYDnlAHabp+S4QUnyAR1Gyu2Y8s5/FUa9WcGIQhEj0KFHkEbN7d5xaaWavzoifv39ALFvrcEiIKIMeHsOSwtuYDnlAHabp+S4QUnyAR1Gyu2Y8s5/FUa9WcGIQh"
	txFixture    = "GpEDCiCLUGKqmJoWFGAjKS1WTXAEkU48Kmq7MiB5rsPW08bLqhLGAgqnAQqkAQocL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBKDAQotY29zbW9zMWEyemh0OHgyajBkcXZ1ZWpyOHB4cHU3ZHVlM3FtazQwbGdkeTNoEkFjb3Ntb3MxZDJqcmg0Z2o2NnNteG5zNnhmdjhtZGQ0a2VlZjVlazk3a25sOWd3OXNra3pyeWp5aHhqc3l3bGZobRoPCgV1YXRvbRIGNTAwMDAwElgKUApGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQLaGco86x6BgxaGOBf/rgbHMEyZzECi+5in9DJ31ln/0BIECgIIARgrEgQQwJoMGkD0/1DW4n4Fp1JZyWDtlWBmi9+ulHrLioDyvQ/4NNLiVSUAj9x4ljCUNwlSzpPtykfjjnGT7IyByWnKB0bGayDWGiQIARogi1BiqpiaFhRgIyktVk1wBJFOPCpquzIgea7D1tPGy6oi2wYK0AQKkgMKAggLEgpnYWlhdGVzdC0xGKrnECIMCNq7+6gGELjv0fEBKkgKICeLs3J8bIJCpWaee12QDGfgsDqmwpDoxQStliWR9bFSEiQIARIgXQsgwP+G56ZGtMeQE1n+8KuZOODcJC75Q6dui0ARvYAyIBu01YZQurtiiOsKiKE/e5CGuD1ioTthG76thsmL10SeOiCLUGKqmJoWFGAjKS1WTXAEkU48Kmq7MiB5rsPW08bLqkIgr8VHiDIZrvPjyaybOEWM23OI1PkRdC+9XJhIjJRBNK1KIK/FR4gyGa7z48msmzhFjNtziNT5EXQvvVyYSIyUQTStUiAEgJG8fdwoP3e/v5HXPETaWMPfipy8hnQF2Lfz2q2iL1ogm6BBj4GRVW4wgJp9qZfWiClAzSc8nzvFbVjT3LGc1PBiIOOwxEKY/BwUmvv0yJlvuSQnrkHkZJuTTKSVmRt4UrhVaiDjsMRCmPwcFJr79MiZb7kkJ65B5GSbk0yklZkbeFK4VXIUeQRs3t3nFppZq/OiJ+/f0AsW+twSuAEIqucQGkgKIB5737NG8FYvnQW6/urw4FNMaM+9CzIhy1MLzQk1/p6WEiQIARIgNltfvzoTATg0D0mHHtrROQIgWFM0QqVxA88cst3U28IiaAgCEhR5BGze3ecWmlmr86In79/QCxb63BoMCN+7+6gGEPjly/sBIkDpPn0WzYyqh6Xx8Bru5+EaA4XFsEsfO6mrXMrZABOgmbrRqHyGcd5wNj2ddC7mj52Ls03KuAsxvWItEYeJLvQGEn4KPQoUeQRs3t3nFppZq/OiJ+/f0AsW+twSIgogx4ew5LC25gOeUAdpun5LhBSfIBHUbK7Zjyzn8VRr1ZwYqyESPQoUeQRs3t3nFppZq/OiJ+/f0AsW+twSIgogx4ew5LC25gOeUAdpun5LhBSfIBHUbK7Zjyzn8VRr1ZwYqyEaBggBEMnnECJ+Cj0KFHkEbN7d5xaaWavzoifv39ALFvrcEiIKIMeHsOSwtuYDnlAHabp+S4QUnyAR1Gyu2Y8s5/FUa9WcGKshEj0KFHkEbN7d5xaaWavzoifv39ALFvrcEiIKIMeHsOSwtuYDnlAHabp+S4QUnyAR1Gyu2Y8s5/FUa9WcGKsh"
)

package keeper_test

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/types/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	channeltypes "github.com/cosmos/ibc-go/v5/modules/core/04-channel/types"

	"github.com/quicksilver-zone/quicksilver/app"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	icqtypes "github.com/quicksilver-zone/quicksilver/x/interchainquery/types"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/keeper"
	icstypes "github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

// ValSetCallback

const (
	delegationQueryCallbackID = "delegation"
	storeStakingKey           = "store/staking/key"
)

func (suite *KeeperTestSuite) TestHandleValsetCallback() {
	newVal := addressutils.GenerateValAddressForTest()

	tests := []struct {
		name   string
		valset func(in stakingtypes.Validators) stakingtypes.QueryValidatorsResponse
		checks func(require *require.Assertions, ctx sdk.Context, quicksilver *app.Quicksilver, in stakingtypes.Validators)
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
			checks: func(require *require.Assertions, ctx sdk.Context, quicksilver *app.Quicksilver, in stakingtypes.Validators) {
				foundQuery := false
				_, addr, _ := bech32.DecodeAndConvert(in[0].OperatorAddress)
				data := stakingtypes.GetValidatorKey(addr)
				for _, i := range quicksilver.InterchainQueryKeeper.AllQueries(ctx) {
					if i.QueryType == storeStakingKey && bytes.Equal(i.Request, data) {
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
			checks: func(require *require.Assertions, ctx sdk.Context, quicksilver *app.Quicksilver, in stakingtypes.Validators) {
				foundQuery := false
				foundQuery2 := false
				_, addr, _ := bech32.DecodeAndConvert(in[1].OperatorAddress)
				data := stakingtypes.GetValidatorKey(addr)
				_, addr2, _ := bech32.DecodeAndConvert(in[2].OperatorAddress)
				data2 := stakingtypes.GetValidatorKey(addr2)
				for _, i := range quicksilver.InterchainQueryKeeper.AllQueries(ctx) {
					if i.QueryType == storeStakingKey && bytes.Equal(i.Request, data) {
						foundQuery = true
					}
					if i.QueryType == storeStakingKey && bytes.Equal(i.Request, data2) {
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
			checks: func(require *require.Assertions, ctx sdk.Context, quicksilver *app.Quicksilver, in stakingtypes.Validators) {
				foundQuery := false
				_, addr, _ := bech32.DecodeAndConvert(in[0].OperatorAddress)
				data := stakingtypes.GetValidatorKey(addr)
				for _, i := range quicksilver.InterchainQueryKeeper.AllQueries(ctx) {
					if i.QueryType == storeStakingKey && bytes.Equal(i.Request, data) {
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
			checks: func(require *require.Assertions, ctx sdk.Context, quicksilver *app.Quicksilver, in stakingtypes.Validators) {
				foundQuery := false
				foundQuery2 := false
				_, addr, _ := bech32.DecodeAndConvert(in[1].OperatorAddress)
				data := stakingtypes.GetValidatorKey(addr)
				_, addr2, _ := bech32.DecodeAndConvert(in[2].OperatorAddress)
				data2 := stakingtypes.GetValidatorKey(addr2)
				for _, i := range quicksilver.InterchainQueryKeeper.AllQueries(ctx) {
					if i.QueryType == storeStakingKey && bytes.Equal(i.Request, data) {
						foundQuery = true
					}
					if i.QueryType == storeStakingKey && bytes.Equal(i.Request, data2) {
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
			checks: func(require *require.Assertions, ctx sdk.Context, quicksilver *app.Quicksilver, in stakingtypes.Validators) {
				foundQuery := false
				foundQuery2 := false
				_, addr, _ := bech32.DecodeAndConvert(in[1].OperatorAddress)
				data := stakingtypes.GetValidatorKey(addr)
				_, addr2, _ := bech32.DecodeAndConvert(in[2].OperatorAddress)
				data2 := stakingtypes.GetValidatorKey(addr2)
				for _, i := range quicksilver.InterchainQueryKeeper.AllQueries(ctx) {
					if i.QueryType == storeStakingKey && bytes.Equal(i.Request, data) {
						foundQuery = true
					}
					if i.QueryType == storeStakingKey && bytes.Equal(i.Request, data2) {
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
			checks: func(require *require.Assertions, ctx sdk.Context, quicksilver *app.Quicksilver, in stakingtypes.Validators) {
				foundQuery := false
				foundQuery2 := false
				_, addr, _ := bech32.DecodeAndConvert(in[0].OperatorAddress)
				data := stakingtypes.GetValidatorKey(addr)
				_, addr2, _ := bech32.DecodeAndConvert(in[2].OperatorAddress)
				data2 := stakingtypes.GetValidatorKey(addr2)
				for _, i := range quicksilver.InterchainQueryKeeper.AllQueries(ctx) {
					if i.QueryType == storeStakingKey && bytes.Equal(i.Request, data) {
						foundQuery = true
					}
					if i.QueryType == storeStakingKey && bytes.Equal(i.Request, data2) {
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
				val := in[0]
				val.OperatorAddress = newVal.String()
				in = append(in, val)
				return stakingtypes.QueryValidatorsResponse{Validators: in}
			},
			checks: func(require *require.Assertions, ctx sdk.Context, quicksilver *app.Quicksilver, in stakingtypes.Validators) {
				foundQuery := false
				data := stakingtypes.GetValidatorKey(newVal)
				for _, i := range quicksilver.InterchainQueryKeeper.AllQueries(ctx) {
					if i.QueryType == storeStakingKey && bytes.Equal(i.Request, data) {
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
		// 	checks: func(require *require.Assertions, ctx sdk.Context, app *quicksilver.Quicksilver, in stakingtypes.Validators) {
		// 		foundQuery := false
		// 		_, addr, _ := bech32.DecodeAndConvert(in[0].OperatorAddress)
		// 		data := stakingtypes.GetValidatorKey(addr)
		// 		for _, i := range quicksilver.InterchainQueryKeeper.AllQueries(ctx) {
		// 			if i.QueryType == storeStakingKey && bytes.Equal(i.Request, data) {
		// 				foundQuery = true
		// 			}
		// 		}
		// 		require.True(foundQuery)
		// 	},
		// },
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			suite.SetupTest()
			suite.setupTestZones()

			quicksilver := suite.GetQuicksilverApp(suite.chainA)
			quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
			ctx := suite.chainA.GetContext()

			chainBVals := suite.GetQuicksilverApp(suite.chainB).StakingKeeper.GetValidators(suite.chainB.GetContext(), 300)

			queryResp := test.valset(chainBVals)
			bz, err := quicksilver.AppCodec().Marshal(&queryResp)
			suite.NoError(err)

			err = keeper.ValsetCallback(quicksilver.InterchainstakingKeeper, ctx, bz, icqtypes.Query{ChainId: suite.chainB.ChainID})
			suite.NoError(err)
			// valset callback doesn't actually update validators, but does emit icq callbacks.
			test.checks(suite.Require(), ctx, quicksilver, chainBVals)
		})
	}
}

func (suite *KeeperTestSuite) TestHandleValsetCallbackBadChain() {
	suite.Run("bad chain", func() {
		suite.SetupTest()
		suite.setupTestZones()

		quicksilver := suite.GetQuicksilverApp(suite.chainA)
		quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := suite.chainA.GetContext()

		queryResp := stakingtypes.QueryValidatorsResponse{Validators: []stakingtypes.Validator{}}
		bz, err := quicksilver.AppCodec().Marshal(&queryResp)
		suite.NoError(err)

		err = keeper.ValsetCallback(quicksilver.InterchainstakingKeeper, ctx, bz, icqtypes.Query{ChainId: "badchain"})
		// this should bail on a non-matching chain id.
		suite.Error(err)
	})
}

func (suite *KeeperTestSuite) TestHandleValsetCallbackNilValset() {
	suite.Run("nil valset", func() {
		suite.SetupTest()
		suite.setupTestZones()

		quicksilver := suite.GetQuicksilverApp(suite.chainA)
		quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := suite.chainA.GetContext()

		queryResp := stakingtypes.QueryValidatorsResponse{Validators: []stakingtypes.Validator{}}
		bz, err := quicksilver.AppCodec().Marshal(&queryResp)
		suite.NoError(err)

		err = keeper.ValsetCallback(quicksilver.InterchainstakingKeeper, ctx, bz, icqtypes.Query{ChainId: suite.chainB.ChainID})
		// this should error on unmarshalling an empty slice, which is not a valid response here.
		suite.Error(err)
	})
}

func (suite *KeeperTestSuite) TestHandleValsetCallbackInvalidResponse() {
	suite.Run("bad payload type", func() {
		suite.SetupTest()
		suite.setupTestZones()

		quicksilver := suite.GetQuicksilverApp(suite.chainA)
		quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := suite.chainA.GetContext()

		queryReq := stakingtypes.QueryValidatorsRequest{Status: stakingtypes.BondStatusBonded}
		bz, err := quicksilver.AppCodec().Marshal(&queryReq)
		suite.NoError(err)

		err = keeper.ValsetCallback(quicksilver.InterchainstakingKeeper, ctx, bz, icqtypes.Query{ChainId: suite.chainB.ChainID})
		// this should error on unmarshalling an empty slice, which is not a valid response here.
		suite.Error(err)
	})
}

func (suite *KeeperTestSuite) TestHandleValidatorCallbackBadChain() {
	suite.Run("bad chain", func() {
		suite.SetupTest()
		suite.setupTestZones()

		quicksilver := suite.GetQuicksilverApp(suite.chainA)
		quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := suite.chainA.GetContext()

		queryResp := stakingtypes.QueryValidatorsResponse{Validators: []stakingtypes.Validator{}}
		bz, err := quicksilver.AppCodec().Marshal(&queryResp)
		suite.NoError(err)

		err = keeper.ValidatorCallback(quicksilver.InterchainstakingKeeper, ctx, bz, icqtypes.Query{ChainId: "badchain"})
		// this should bail on a non-matching chain id.
		suite.Error(err)
	})
}

func (suite *KeeperTestSuite) TestHandleValidatorCallbackNilValue() {
	suite.Run("empty value", func() {
		suite.SetupTest()
		suite.setupTestZones()

		quicksilver := suite.GetQuicksilverApp(suite.chainA)
		quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := suite.chainA.GetContext()

		bz := []byte{}

		err := keeper.ValidatorCallback(quicksilver.InterchainstakingKeeper, ctx, bz, icqtypes.Query{ChainId: suite.chainB.ChainID})
		// this should error on unmarshalling an empty slice, which is not a valid response here.
		suite.Error(err)
	})
}

func (suite *KeeperTestSuite) TestHandleValidatorCallback() {
	newVal := addressutils.GenerateAddressForTestWithPrefix("cosmosvaloper")
	zone := icstypes.Zone{ConnectionId: "connection-0", ChainId: "cosmoshub-4", AccountPrefix: "cosmos", LocalDenom: "uqatom", BaseDenom: "uatom", Is_118: true}
	err := suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetValidator(suite.chainA.GetContext(), zone.ChainId, icstypes.Validator{ValoperAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), DelegatorShares: sdk.NewDec(2000)})
	suite.NoError(err)

	err = suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetValidator(suite.chainA.GetContext(), zone.ChainId, icstypes.Validator{ValoperAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), DelegatorShares: sdk.NewDec(2000)})
	suite.NoError(err)

	err = suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetValidator(suite.chainA.GetContext(), zone.ChainId, icstypes.Validator{ValoperAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), DelegatorShares: sdk.NewDec(2000)})
	suite.NoError(err)

	err = suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetValidator(suite.chainA.GetContext(), zone.ChainId, icstypes.Validator{ValoperAddress: "cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), DelegatorShares: sdk.NewDec(2000)})
	suite.NoError(err)

	err = suite.GetQuicksilverApp(suite.chainA).InterchainstakingKeeper.SetValidator(suite.chainA.GetContext(), zone.ChainId, icstypes.Validator{ValoperAddress: "cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), DelegatorShares: sdk.NewDec(2000)})
	suite.NoError(err)

	tests := []struct {
		name      string
		validator stakingtypes.Validator
		expected  icstypes.Validator
	}{
		{
			name:      "valid - no-op",
			validator: stakingtypes.Validator{OperatorAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", Jailed: false, Status: stakingtypes.Bonded, Tokens: sdk.NewInt(2000), DelegatorShares: sdk.NewDec(2000), Commission: stakingtypes.NewCommission(sdk.MustNewDecFromStr("0.2"), sdk.MustNewDecFromStr("0.2"), sdk.MustNewDecFromStr("0.2"))},
			expected:  icstypes.Validator{ValoperAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), DelegatorShares: sdk.NewDec(2000), Score: sdk.ZeroDec(), Status: "BOND_STATUS_BONDED"},
		},
		{
			name:      "valid - +2000 tokens/shares",
			validator: stakingtypes.Validator{OperatorAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", Jailed: false, Status: stakingtypes.Bonded, Tokens: sdk.NewInt(4000), DelegatorShares: sdk.NewDec(4000), Commission: stakingtypes.NewCommission(sdk.MustNewDecFromStr("0.2"), sdk.MustNewDecFromStr("0.2"), sdk.MustNewDecFromStr("0.2"))},
			expected:  icstypes.Validator{ValoperAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(4000), DelegatorShares: sdk.NewDec(4000), Score: sdk.ZeroDec(), Status: "BOND_STATUS_BONDED"},
		},
		{
			name:      "valid - inc. commission",
			validator: stakingtypes.Validator{OperatorAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", Jailed: false, Status: stakingtypes.Bonded, Tokens: sdk.NewInt(2000), DelegatorShares: sdk.NewDec(2000), Commission: stakingtypes.NewCommission(sdk.MustNewDecFromStr("0.5"), sdk.MustNewDecFromStr("0.2"), sdk.MustNewDecFromStr("0.2"))},
			expected:  icstypes.Validator{ValoperAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", CommissionRate: sdk.MustNewDecFromStr("0.5"), VotingPower: sdk.NewInt(2000), DelegatorShares: sdk.NewDec(2000), Score: sdk.ZeroDec(), Status: "BOND_STATUS_BONDED"},
		},
		{
			name:      "valid - new validator",
			validator: stakingtypes.Validator{OperatorAddress: newVal, Jailed: false, Status: stakingtypes.Bonded, Tokens: sdk.NewInt(3000), DelegatorShares: sdk.NewDec(3050), Commission: stakingtypes.NewCommission(sdk.MustNewDecFromStr("0.25"), sdk.MustNewDecFromStr("0.2"), sdk.MustNewDecFromStr("0.2"))},
			expected:  icstypes.Validator{ValoperAddress: newVal, CommissionRate: sdk.MustNewDecFromStr("0.25"), VotingPower: sdk.NewInt(3000), DelegatorShares: sdk.NewDec(3050), Score: sdk.ZeroDec(), Status: "BOND_STATUS_BONDED"},
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			test := test
			suite.SetupTest()
			suite.setupTestZones()

			quicksilver := suite.GetQuicksilverApp(suite.chainA)
			quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
			ctx := suite.chainA.GetContext()

			quicksilver.InterchainstakingKeeper.SetZone(ctx, &zone)

			bz, err := quicksilver.AppCodec().Marshal(&test.validator)
			suite.NoError(err)

			err = keeper.ValidatorCallback(quicksilver.InterchainstakingKeeper, ctx, bz, icqtypes.Query{ChainId: zone.ChainId})
			suite.NoError(err)

			zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, zone.ChainId)
			suite.True(found)

			valAddrBytes, err := addressutils.ValAddressFromBech32(test.expected.ValoperAddress, zone.GetValoperPrefix())
			suite.NoError(err)
			val, found := quicksilver.InterchainstakingKeeper.GetValidator(ctx, zone.ChainId, valAddrBytes)
			suite.True(found)
			suite.Equal(test.expected, val)
		})
	}
}

func (suite *KeeperTestSuite) TestHandleValidatorCallbackJailedWithSlashing() {
	completion := time.Now().UTC().Add(time.Hour)

	tests := []struct {
		name               string
		validator          func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) *stakingtypes.Validator
		expected           func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) *icstypes.Validator
		withdrawal         func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) icstypes.WithdrawalRecord
		expectedWithdrawal func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) icstypes.WithdrawalRecord
	}{
		{
			name: "jailed; single distribution",
			validator: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) *stakingtypes.Validator {
				vals := qs.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)
				return &stakingtypes.Validator{OperatorAddress: vals[0].ValoperAddress, Jailed: true, Status: stakingtypes.Bonded, Tokens: vals[0].VotingPower.Mul(sdk.NewInt(19)).Quo(sdk.NewInt(20)), DelegatorShares: vals[0].DelegatorShares, Commission: stakingtypes.NewCommission(vals[0].CommissionRate, sdk.MustNewDecFromStr("0.5"), sdk.MustNewDecFromStr("0.5"))}
			},
			expected: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) *icstypes.Validator {
				vals := qs.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)
				return &icstypes.Validator{ValoperAddress: vals[0].ValoperAddress, CommissionRate: vals[0].CommissionRate, VotingPower: vals[0].VotingPower.Mul(sdk.NewInt(19)).Quo(sdk.NewInt(20)), DelegatorShares: vals[0].DelegatorShares, Score: sdk.ZeroDec(), Status: "BOND_STATUS_BONDED", Jailed: true}
			},

			withdrawal: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) icstypes.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)
				return icstypes.WithdrawalRecord{
					ChainId:   suite.chainB.ChainID,
					Delegator: zone.DelegationAddress.Address,
					Distribution: []*icstypes.Distribution{
						{
							Valoper: vals[0].ValoperAddress,
							Amount:  1000,
						},
					},
					Recipient:      user1.String(),
					Amount:         sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))),
					BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1000)),
					Txhash:         "1613D2E8FBF7C7294A4D2247B55EE89FB22FC68C62D61050B944F1191DF092BD",
					Status:         icstypes.WithdrawStatusUnbond,
					CompletionTime: completion,
				}
			},
			expectedWithdrawal: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) icstypes.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)
				return icstypes.WithdrawalRecord{
					ChainId:   suite.chainB.ChainID,
					Delegator: zone.DelegationAddress.Address,
					Distribution: []*icstypes.Distribution{
						{
							Valoper: vals[0].ValoperAddress,
							Amount:  950,
						},
					},
					Recipient:      user1.String(),
					Amount:         sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(950))),
					BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1000)),
					Txhash:         "1613D2E8FBF7C7294A4D2247B55EE89FB22FC68C62D61050B944F1191DF092BD",
					Status:         icstypes.WithdrawStatusUnbond,
					CompletionTime: completion,
				}
			},
		},
		{
			name: "jailed; multi distribution",
			validator: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) *stakingtypes.Validator {
				vals := qs.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)
				return &stakingtypes.Validator{OperatorAddress: vals[0].ValoperAddress, Jailed: true, Status: stakingtypes.Bonded, Tokens: vals[0].VotingPower.Mul(sdk.NewInt(19)).Quo(sdk.NewInt(20)), DelegatorShares: vals[0].DelegatorShares, Commission: stakingtypes.NewCommission(vals[0].CommissionRate, sdk.MustNewDecFromStr("0.5"), sdk.MustNewDecFromStr("0.5"))}
			},
			expected: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) *icstypes.Validator {
				vals := qs.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)
				return &icstypes.Validator{ValoperAddress: vals[0].ValoperAddress, CommissionRate: vals[0].CommissionRate, VotingPower: vals[0].VotingPower.Mul(sdk.NewInt(19)).Quo(sdk.NewInt(20)), DelegatorShares: vals[0].DelegatorShares, Score: sdk.ZeroDec(), Status: "BOND_STATUS_BONDED", Jailed: true}
			},

			withdrawal: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) icstypes.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)
				return icstypes.WithdrawalRecord{
					ChainId:   suite.chainB.ChainID,
					Delegator: zone.DelegationAddress.Address,
					Distribution: []*icstypes.Distribution{
						{
							Valoper: vals[0].ValoperAddress,
							Amount:  500,
						},
						{
							Valoper: vals[1].ValoperAddress,
							Amount:  500,
						},
					},
					Recipient:      user1.String(),
					Amount:         sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))),
					BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1000)),
					Txhash:         "1613D2E8FBF7C7294A4D2247B55EE89FB22FC68C62D61050B944F1191DF092BD",
					Status:         icstypes.WithdrawStatusUnbond,
					CompletionTime: completion,
				}
			},
			expectedWithdrawal: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) icstypes.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)
				return icstypes.WithdrawalRecord{
					ChainId:   suite.chainB.ChainID,
					Delegator: zone.DelegationAddress.Address,
					Distribution: []*icstypes.Distribution{
						{
							Valoper: vals[0].ValoperAddress,
							Amount:  475,
						},
						{
							Valoper: vals[1].ValoperAddress,
							Amount:  500,
						},
					},
					Recipient:      user1.String(),
					Amount:         sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(975))),
					BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1000)),
					Txhash:         "1613D2E8FBF7C7294A4D2247B55EE89FB22FC68C62D61050B944F1191DF092BD",
					Status:         icstypes.WithdrawStatusUnbond,
					CompletionTime: completion,
				}
			},
		},
		{
			name: "jailed; multi distribution, unrelated validators - no-op",
			validator: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) *stakingtypes.Validator {
				vals := qs.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)
				return &stakingtypes.Validator{OperatorAddress: vals[0].ValoperAddress, Jailed: true, Status: stakingtypes.Bonded, Tokens: vals[0].VotingPower.Mul(sdk.NewInt(19)).Quo(sdk.NewInt(20)), DelegatorShares: vals[0].DelegatorShares, Commission: stakingtypes.NewCommission(vals[0].CommissionRate, sdk.MustNewDecFromStr("0.5"), sdk.MustNewDecFromStr("0.5"))}
			},
			expected: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) *icstypes.Validator {
				vals := qs.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)
				return &icstypes.Validator{ValoperAddress: vals[0].ValoperAddress, CommissionRate: vals[0].CommissionRate, VotingPower: vals[0].VotingPower.Mul(sdk.NewInt(19)).Quo(sdk.NewInt(20)), DelegatorShares: vals[0].DelegatorShares, Score: sdk.ZeroDec(), Status: "BOND_STATUS_BONDED", Jailed: true}
			},

			withdrawal: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) icstypes.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)
				return icstypes.WithdrawalRecord{
					ChainId:   suite.chainB.ChainID,
					Delegator: zone.DelegationAddress.Address,
					Distribution: []*icstypes.Distribution{
						{
							Valoper: vals[1].ValoperAddress,
							Amount:  500,
						},
						{
							Valoper: vals[2].ValoperAddress,
							Amount:  500,
						},
					},
					Recipient:      user1.String(),
					Amount:         sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))),
					BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1000)),
					Txhash:         "1613D2E8FBF7C7294A4D2247B55EE89FB22FC68C62D61050B944F1191DF092BD",
					Status:         icstypes.WithdrawStatusUnbond,
					CompletionTime: completion,
				}
			},
			expectedWithdrawal: func(ctx sdk.Context, qs *app.Quicksilver, zone icstypes.Zone) icstypes.WithdrawalRecord {
				vals := qs.InterchainstakingKeeper.GetValidators(ctx, zone.ChainId)
				return icstypes.WithdrawalRecord{
					ChainId:   suite.chainB.ChainID,
					Delegator: zone.DelegationAddress.Address,
					Distribution: []*icstypes.Distribution{
						{
							Valoper: vals[1].ValoperAddress,
							Amount:  500,
						},
						{
							Valoper: vals[2].ValoperAddress,
							Amount:  500,
						},
					},
					Recipient:      user1.String(),
					Amount:         sdk.NewCoins(sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))),
					BurnAmount:     sdk.NewCoin(zone.LocalDenom, sdk.NewInt(1000)),
					Txhash:         "1613D2E8FBF7C7294A4D2247B55EE89FB22FC68C62D61050B944F1191DF092BD",
					Status:         icstypes.WithdrawStatusUnbond,
					CompletionTime: completion,
				}
			},
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			suite.SetupTest()
			suite.setupTestZones()

			quicksilver := suite.GetQuicksilverApp(suite.chainA)
			quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
			ctx := suite.chainA.GetContext()

			zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
			suite.True(found)

			quicksilver.InterchainstakingKeeper.SetWithdrawalRecord(ctx, test.withdrawal(ctx, quicksilver, zone))

			bz, err := quicksilver.AppCodec().Marshal(test.validator(ctx, quicksilver, zone))
			suite.NoError(err)

			err = keeper.ValidatorCallback(quicksilver.InterchainstakingKeeper, ctx, bz, icqtypes.Query{ChainId: zone.ChainId})
			suite.NoError(err)

			wr, found := quicksilver.InterchainstakingKeeper.GetWithdrawalRecord(ctx, suite.chainB.ChainID, test.withdrawal(ctx, quicksilver, zone).Txhash, test.withdrawal(ctx, quicksilver, zone).Status)
			suite.True(found)
			suite.Equal(test.expectedWithdrawal(ctx, quicksilver, zone), wr)
		})
	}
}

func (suite *KeeperTestSuite) TestHandleRewardsCallbackBadChain() {
	suite.Run("bad chain", func() {
		suite.SetupTest()
		suite.setupTestZones()

		quicksilver := suite.GetQuicksilverApp(suite.chainA)
		quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := suite.chainA.GetContext()

		queryRes := distrtypes.QueryDelegationTotalRewardsResponse{}
		bz, err := quicksilver.AppCodec().Marshal(&queryRes)
		suite.NoError(err)

		err = keeper.RewardsCallback(quicksilver.InterchainstakingKeeper, ctx, bz, icqtypes.Query{ChainId: "badchain"})
		// this should bail on a non-matching chain id.
		suite.Error(err)
	})
}

func (suite *KeeperTestSuite) TestHandleRewardsEmptyRequestCallback() {
	suite.Run("empty request", func() {
		suite.SetupTest()
		suite.setupTestZones()

		quicksilver := suite.GetQuicksilverApp(suite.chainA)
		quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := suite.chainA.GetContext()

		queryReq := distrtypes.QueryDelegationTotalRewardsRequest{}
		bz, err := quicksilver.AppCodec().Marshal(&queryReq)
		suite.NoError(err)

		err = keeper.RewardsCallback(quicksilver.InterchainstakingKeeper, ctx, bz, icqtypes.Query{ChainId: suite.chainB.ChainID})
		// this should fail because the waitgroup becomes negative.
		suite.Errorf(err, "attempted to unmarshal zero length byte slice (2)")
	})
}

func (suite *KeeperTestSuite) TestHandleRewardsCallbackNonDelegator() {
	suite.Run("valid response, bad delegator", func() {
		suite.SetupTest()
		suite.setupTestZones()

		quicksilver := suite.GetQuicksilverApp(suite.chainA)
		quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := suite.chainA.GetContext()

		zone, _ := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
		zone.WithdrawalWaitgroup++
		quicksilver.InterchainstakingKeeper.SetZone(ctx, &zone)

		user := addressutils.GenerateAccAddressForTest()

		queryReq := distrtypes.QueryDelegationTotalRewardsRequest{
			DelegatorAddress: user.String(),
		}

		response := distrtypes.QueryDelegationTotalRewardsResponse{
			Rewards: []distrtypes.DelegationDelegatorReward{
				{ValidatorAddress: suite.chainB.Vals.Validators[0].String(), Reward: sdk.NewDecCoins(sdk.NewDecCoin("uatom", sdk.NewInt((1000))))},
			},
			Total: sdk.NewDecCoins(sdk.NewDecCoin("uatom", sdk.NewInt((1000)))),
		}
		reqbz, err := quicksilver.AppCodec().Marshal(&queryReq)
		suite.NoError(err)

		respbz, err := quicksilver.AppCodec().Marshal(&response)
		suite.NoError(err)
		err = keeper.RewardsCallback(quicksilver.InterchainstakingKeeper, ctx, respbz, icqtypes.Query{ChainId: suite.chainB.ChainID, Request: reqbz})
		//
		suite.Errorf(err, "failed attempting to withdraw rewards from non-delegation account")
	})
}

func (suite *KeeperTestSuite) TestHandleRewardsCallbackEmptyResponse() {
	suite.Run("empty response", func() {
		suite.SetupTest()
		suite.setupTestZones()

		quicksilver := suite.GetQuicksilverApp(suite.chainA)
		quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := suite.chainA.GetContext()

		zone, _ := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
		zone.WithdrawalWaitgroup++
		quicksilver.InterchainstakingKeeper.SetZone(ctx, &zone)

		queryReq := distrtypes.QueryDelegationTotalRewardsRequest{
			DelegatorAddress: zone.DelegationAddress.Address,
		}

		response := distrtypes.QueryDelegationTotalRewardsResponse{}
		reqbz, err := quicksilver.AppCodec().Marshal(&queryReq)
		suite.NoError(err)

		respbz, err := quicksilver.AppCodec().Marshal(&response)
		suite.NoError(err)
		err = keeper.RewardsCallback(quicksilver.InterchainstakingKeeper, ctx, respbz, icqtypes.Query{ChainId: suite.chainB.ChainID, Request: reqbz})
		//
		suite.NoError(err)
	})
}

func (suite *KeeperTestSuite) TestHandleValideRewardsCallback() {
	suite.Run("valid response, negative waitgroup", func() {
		suite.SetupTest()
		suite.setupTestZones()

		quicksilver := suite.GetQuicksilverApp(suite.chainA)
		quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := suite.chainA.GetContext()

		zone, _ := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
		zone.WithdrawalWaitgroup++
		quicksilver.InterchainstakingKeeper.SetZone(ctx, &zone)

		queryReq := distrtypes.QueryDelegationTotalRewardsRequest{
			DelegatorAddress: zone.DelegationAddress.Address,
		}

		response := distrtypes.QueryDelegationTotalRewardsResponse{
			Rewards: []distrtypes.DelegationDelegatorReward{
				{ValidatorAddress: suite.chainB.Vals.Validators[0].String(), Reward: sdk.NewDecCoins(sdk.NewDecCoin("uatom", sdk.NewInt((1000))))},
			},
			Total: sdk.NewDecCoins(sdk.NewDecCoin("uatom", sdk.NewInt((1000)))),
		}
		reqbz, err := quicksilver.AppCodec().Marshal(&queryReq)
		suite.NoError(err)

		respbz, err := quicksilver.AppCodec().Marshal(&response)
		suite.NoError(err)
		err = keeper.RewardsCallback(quicksilver.InterchainstakingKeeper, ctx, respbz, icqtypes.Query{ChainId: suite.chainB.ChainID, Request: reqbz})
		//
		suite.NoError(err)
	})
}

func (suite *KeeperTestSuite) TestHandleDistributeRewardsCallback() {
	suite.SetupTest()
	suite.setupTestZones()
	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	gaia := suite.GetQuicksilverApp(suite.chainB)
	quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()

	ctxA := suite.chainA.GetContext()
	ctxB := suite.chainB.GetContext()
	vals := gaia.StakingKeeper.GetAllValidators(ctxB)

	zone, _ := quicksilver.InterchainstakingKeeper.GetZone(ctxA, suite.chainB.ChainID)
	params := quicksilver.InterchainstakingKeeper.GetParams(ctxA)
	commisionRate := sdk.MustNewDecFromStr("0.2")
	params.CommissionRate = commisionRate
	quicksilver.InterchainstakingKeeper.SetParams(ctxA, params)

	prevRedemptionRate := zone.RedemptionRate
	tests := []struct {
		name            string
		zoneSetup       func()
		connectionSetup func() string
		responseMsg     func() []byte
		queryMsg        icqtypes.Query
		check           func()
		pass            bool
	}{
		{
			name: "valid case with positive rewards and -5% < delta < 2%",
			zoneSetup: func() {
				balances := sdk.NewCoins(
					sdk.NewCoin(
						zone.LocalDenom,
						math.NewInt(100_000_000),
					),
				)
				err := quicksilver.MintKeeper.MintCoins(ctxA, balances)
				suite.NoError(err)
				qAssetAmount := quicksilver.BankKeeper.GetSupply(ctxA, zone.LocalDenom)
				suite.Equal(balances[0], qAssetAmount)

				delegation := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[0].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(100_000_000))}
				quicksilver.InterchainstakingKeeper.SetDelegation(ctxA, &zone, delegation)
			},
			connectionSetup: func() string {
				channelID := quicksilver.IBCKeeper.ChannelKeeper.GenerateChannelIdentifier(ctxA)
				quicksilver.IBCKeeper.ChannelKeeper.SetChannel(ctxA, icstypes.TransferPort, channelID, channeltypes.Channel{State: channeltypes.OPEN, Ordering: channeltypes.ORDERED, Counterparty: channeltypes.Counterparty{PortId: icstypes.TransferPort, ChannelId: channelID}, ConnectionHops: []string{suite.path.EndpointA.ConnectionID}})
				quicksilver.IBCKeeper.ChannelKeeper.SetNextSequenceSend(ctxA, icstypes.TransferPort, channelID, 1)
				return channelID
			},
			responseMsg: func() []byte {
				balances := sdk.NewCoins(
					sdk.NewCoin(
						zone.BaseDenom,
						math.NewInt(1_000_000),
					),
				)

				response := banktypes.QueryAllBalancesResponse{
					Balances: balances,
				}
				respbz, err := quicksilver.AppCodec().Marshal(&response)
				suite.NoError(err)
				return respbz
			},
			queryMsg: icqtypes.Query{ChainId: suite.chainB.ChainID},
			check: func() {
				zone, _ := quicksilver.InterchainstakingKeeper.GetZone(ctxA, suite.chainB.ChainID)
				redemptionRate := zone.RedemptionRate
				ratio := sdk.MustNewDecFromStr("1.008")
				suite.Equal(ratio.Mul(prevRedemptionRate), redemptionRate)
			},
			pass: true,
		},
		{
			name: "valid case with no rewards",
			zoneSetup: func() {},
			connectionSetup: func() string {
				channelID := quicksilver.IBCKeeper.ChannelKeeper.GenerateChannelIdentifier(ctxA)
				quicksilver.IBCKeeper.ChannelKeeper.SetChannel(ctxA, icstypes.TransferPort, channelID, channeltypes.Channel{State: channeltypes.OPEN, Ordering: channeltypes.ORDERED, Counterparty: channeltypes.Counterparty{PortId: icstypes.TransferPort, ChannelId: channelID}, ConnectionHops: []string{suite.path.EndpointA.ConnectionID}})
				quicksilver.IBCKeeper.ChannelKeeper.SetNextSequenceSend(ctxA, icstypes.TransferPort, channelID, 1)
				return channelID
			},
			responseMsg: func() []byte {
				response := banktypes.QueryAllBalancesResponse{
					Balances: sdk.Coins{},
				}
				respbz, err := quicksilver.AppCodec().Marshal(&response)
				suite.NoError(err)
				return respbz
			},
			queryMsg: icqtypes.Query{ChainId: suite.chainB.ChainID},
			check: func() {
				zone, _ := quicksilver.InterchainstakingKeeper.GetZone(ctxA, suite.chainB.ChainID)
				redemptionRate := zone.RedemptionRate

				suite.Equal(prevRedemptionRate, redemptionRate)
			},
			pass: true,
		},
		{
			name:      "invalid host zone",
			zoneSetup: func() {},
			connectionSetup: func() string {
				channelID := quicksilver.IBCKeeper.ChannelKeeper.GenerateChannelIdentifier(ctxA)
				quicksilver.IBCKeeper.ChannelKeeper.SetChannel(ctxA, icstypes.TransferPort, channelID, channeltypes.Channel{State: channeltypes.OPEN, Ordering: channeltypes.ORDERED, Counterparty: channeltypes.Counterparty{PortId: icstypes.TransferPort, ChannelId: channelID}, ConnectionHops: []string{suite.path.EndpointA.ConnectionID}})
				quicksilver.IBCKeeper.ChannelKeeper.SetNextSequenceSend(ctxA, icstypes.TransferPort, channelID, 1)
				return channelID
			},
			responseMsg: func() []byte {
				balances := sdk.NewCoins(
					sdk.NewCoin(
						zone.BaseDenom,
						math.NewInt(10_000_000),
					),
				)

				response := banktypes.QueryAllBalancesResponse{
					Balances: balances,
				}
				respbz, err := quicksilver.AppCodec().Marshal(&response)
				suite.NoError(err)
				return respbz
			},
			queryMsg: icqtypes.Query{ChainId: ""},
			check:    func() {},
			pass:     false,
		},
		{
			name:      "invalid response",
			zoneSetup: func() {},
			connectionSetup: func() string {
				channelID := quicksilver.IBCKeeper.ChannelKeeper.GenerateChannelIdentifier(ctxA)
				quicksilver.IBCKeeper.ChannelKeeper.SetChannel(ctxA, icstypes.TransferPort, channelID, channeltypes.Channel{State: channeltypes.OPEN, Ordering: channeltypes.ORDERED, Counterparty: channeltypes.Counterparty{PortId: icstypes.TransferPort, ChannelId: channelID}, ConnectionHops: []string{suite.path.EndpointA.ConnectionID}})
				quicksilver.IBCKeeper.ChannelKeeper.SetNextSequenceSend(ctxA, icstypes.TransferPort, channelID, 1)
				return channelID
			},
			responseMsg: func() []byte {
				balance := sdk.NewCoin(
					zone.BaseDenom,
					math.NewInt(10_000_000),
				)
				respbz, err := quicksilver.AppCodec().Marshal(&balance)
				suite.NoError(err)
				return respbz
			},
			queryMsg: icqtypes.Query{ChainId: suite.chainB.ChainID},
			check:    func() {},
			pass:     false,
		},
		{
			name:      "no connection setup",
			zoneSetup: func() {},
			connectionSetup: func() string {
				return ""
			},
			responseMsg: func() []byte {
				balances := sdk.NewCoins(
					sdk.NewCoin(
						zone.BaseDenom,
						math.NewInt(10_000_000),
					),
				)

				response := banktypes.QueryAllBalancesResponse{
					Balances: balances,
				}
				respbz, err := quicksilver.AppCodec().Marshal(&response)
				suite.NoError(err)
				return respbz
			},
			queryMsg: icqtypes.Query{ChainId: suite.chainB.ChainID},
			check:    func() {},
			pass:     false,
		},
	}
	for _, test := range tests {
		suite.Run(test.name, func() {
			fmt.Println("redemption rate: ", zone.RedemptionRate)

			// Send coin to withdrawal address
			balances := sdk.NewCoins(
				sdk.NewCoin(
					zone.BaseDenom,
					math.NewInt(10_000_000),
				),
			)
			err := gaia.MintKeeper.MintCoins(ctxB, balances)
			suite.NoError(err)
			addr, err := addressutils.AccAddressFromBech32(zone.WithdrawalAddress.Address, "")
			suite.NoError(err)
			err = gaia.BankKeeper.SendCoinsFromModuleToAccount(ctxB, minttypes.ModuleName, addr, balances)
			suite.NoError(err)

			test.zoneSetup()
			channelID := test.connectionSetup()

			respbz := test.responseMsg()
			err = keeper.DistributeRewardsFromWithdrawAccount(quicksilver.InterchainstakingKeeper, ctxA, respbz, test.queryMsg)

			if test.pass {
				suite.NoError(err)
			} else {
				suite.Error(err)
			}

			test.check()
			zone, _ = quicksilver.InterchainstakingKeeper.GetZone(ctxA, suite.chainB.ChainID)
			fmt.Println("redemption rate: ", zone.RedemptionRate)

			commitments := quicksilver.IBCKeeper.ChannelKeeper.GetAllPacketCommitments(ctxA)
			fmt.Println("commitments: ", commitments[0])

			channel, found := quicksilver.IBCKeeper.ChannelKeeper.GetChannel(ctxA, icstypes.TransferPort, channelID)
			if found {
				channel.State = channeltypes.CLOSED
				quicksilver.IBCKeeper.ChannelKeeper.SetChannel(ctxA, icstypes.TransferPort, channelID, channel)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestAllBalancesCallback() {
	suite.Run("all balances non-zero)", func() {
		suite.SetupTest()
		suite.setupTestZones()

		quicksilver := suite.GetQuicksilverApp(suite.chainA)
		quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := suite.chainA.GetContext()

		zone, _ := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)

		queryReq := banktypes.QueryAllBalancesRequest{
			Address: zone.DepositAddress.Address,
		}
		reqbz, err := quicksilver.AppCodec().Marshal(&queryReq)
		suite.NoError(err)

		response := banktypes.QueryAllBalancesResponse{Balances: sdk.NewCoins(sdk.NewCoin("uqck", sdk.OneInt()))}
		respbz, err := quicksilver.AppCodec().Marshal(&response)
		suite.NoError(err)

		err = keeper.AllBalancesCallback(quicksilver.InterchainstakingKeeper, ctx, respbz, icqtypes.Query{ChainId: suite.chainB.ChainID, Request: reqbz})
		suite.NoError(err)

		// refetch zone
		zone, _ = quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
		suite.Equal(uint32(1), zone.DepositAddress.BalanceWaitgroup)

		_, addr, err := bech32.DecodeAndConvert(zone.DepositAddress.Address)
		suite.NoError(err)
		data := banktypes.CreateAccountBalancesPrefix(addr)

		// check a ICQ request was made
		found := false
		quicksilver.InterchainQueryKeeper.IterateQueries(ctx, func(index int64, queryInfo icqtypes.Query) (stop bool) {
			if queryInfo.ChainId == zone.ChainId &&
				queryInfo.ConnectionId == zone.ConnectionId &&
				queryInfo.QueryType == icstypes.BankStoreKey &&
				bytes.Equal(queryInfo.Request, append(data, []byte(response.Balances[0].GetDenom())...)) {
				found = true
				return true
			}
			return false
		})
		suite.True(found)
	})
}

func (suite *KeeperTestSuite) TestAllBalancesCallbackWithExistingWg() {
	suite.Run("all balances non-zero)", func() {
		suite.SetupTest()
		suite.setupTestZones()

		quicksilver := suite.GetQuicksilverApp(suite.chainA)
		quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := suite.chainA.GetContext()

		zone, _ := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
		zone.DepositAddress.BalanceWaitgroup = 2
		quicksilver.InterchainstakingKeeper.SetZone(ctx, &zone)

		queryReq := banktypes.QueryAllBalancesRequest{
			Address: zone.DepositAddress.Address,
		}
		reqbz, err := quicksilver.AppCodec().Marshal(&queryReq)
		suite.NoError(err)

		response := banktypes.QueryAllBalancesResponse{Balances: sdk.NewCoins(sdk.NewCoin("uqck", sdk.OneInt()))}
		respbz, err := quicksilver.AppCodec().Marshal(&response)
		suite.NoError(err)

		err = keeper.AllBalancesCallback(quicksilver.InterchainstakingKeeper, ctx, respbz, icqtypes.Query{ChainId: suite.chainB.ChainID, Request: reqbz})
		suite.NoError(err)

		// refetch zone
		zone, _ = quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
		suite.Equal(uint32(1), zone.DepositAddress.BalanceWaitgroup)

		_, addr, err := bech32.DecodeAndConvert(zone.DepositAddress.Address)
		suite.NoError(err)
		data := banktypes.CreateAccountBalancesPrefix(addr)

		// check a ICQ request was made
		found := false
		quicksilver.InterchainQueryKeeper.IterateQueries(ctx, func(index int64, queryInfo icqtypes.Query) (stop bool) {
			if queryInfo.ChainId == zone.ChainId &&
				queryInfo.ConnectionId == zone.ConnectionId &&
				queryInfo.QueryType == icstypes.BankStoreKey &&
				bytes.Equal(queryInfo.Request, append(data, []byte(response.Balances[0].GetDenom())...)) {
				found = true
				return true
			}
			return false
		})
		suite.True(found)
	})
}

// tests where we have an existing balance and that balance is now reported as zero.
// we expect that an icq query will be emitted to assert with proof that the balance
// is now zero.
func (suite *KeeperTestSuite) TestAllBalancesCallbackExistingBalanceNowNil() {
	suite.Run("existing balance - now zero - deposit", func() {
		suite.SetupTest()
		suite.setupTestZones()

		quicksilver := suite.GetQuicksilverApp(suite.chainA)
		quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := suite.chainA.GetContext()

		zone, _ := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
		zone.DepositAddress.Balance = sdk.NewCoins(sdk.NewCoin("uqck", sdk.OneInt()))
		quicksilver.InterchainstakingKeeper.SetZone(ctx, &zone)

		queryReq := banktypes.QueryAllBalancesRequest{
			Address: zone.DepositAddress.Address,
		}
		reqbz, err := quicksilver.AppCodec().Marshal(&queryReq)
		suite.NoError(err)

		response := banktypes.QueryAllBalancesResponse{Balances: sdk.Coins{}}
		respbz, err := quicksilver.AppCodec().Marshal(&response)
		suite.NoError(err)

		err = keeper.AllBalancesCallback(quicksilver.InterchainstakingKeeper, ctx, respbz, icqtypes.Query{ChainId: suite.chainB.ChainID, Request: reqbz})
		suite.NoError(err)

		// refetch zone
		zone, _ = quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
		suite.Equal(uint32(1), zone.DepositAddress.BalanceWaitgroup)

		_, addr, err := bech32.DecodeAndConvert(zone.DepositAddress.Address)
		suite.NoError(err)
		data := banktypes.CreateAccountBalancesPrefix(addr)

		// check a ICQ request was made
		found := false
		quicksilver.InterchainQueryKeeper.IterateQueries(ctx, func(index int64, queryInfo icqtypes.Query) (stop bool) {
			if queryInfo.ChainId == zone.ChainId &&
				queryInfo.ConnectionId == zone.ConnectionId &&
				queryInfo.QueryType == icstypes.BankStoreKey &&
				bytes.Equal(queryInfo.Request, append(data, []byte("uqck")...)) {
				found = true
				return true
			}
			return false
		})
		suite.True(found)
	})

	suite.Run("existing balance - now zero - withdrawal", func() {
		suite.SetupTest()
		suite.setupTestZones()

		quicksilver := suite.GetQuicksilverApp(suite.chainA)
		quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := suite.chainA.GetContext()

		zone, _ := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
		zone.WithdrawalAddress.Balance = sdk.NewCoins(sdk.NewCoin("uqck", sdk.OneInt()))
		quicksilver.InterchainstakingKeeper.SetZone(ctx, &zone)

		queryReq := banktypes.QueryAllBalancesRequest{
			Address: zone.WithdrawalAddress.Address,
		}
		reqbz, err := quicksilver.AppCodec().Marshal(&queryReq)
		suite.NoError(err)

		response := banktypes.QueryAllBalancesResponse{Balances: sdk.Coins{}}
		respbz, err := quicksilver.AppCodec().Marshal(&response)
		suite.NoError(err)

		err = keeper.AllBalancesCallback(quicksilver.InterchainstakingKeeper, ctx, respbz, icqtypes.Query{ChainId: suite.chainB.ChainID, Request: reqbz})
		suite.NoError(err)

		// refetch zone
		zone, _ = quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
		suite.Equal(uint32(1), zone.WithdrawalAddress.BalanceWaitgroup)

		_, addr, err := bech32.DecodeAndConvert(zone.WithdrawalAddress.Address)
		suite.NoError(err)
		data := banktypes.CreateAccountBalancesPrefix(addr)

		// check a ICQ request was made
		found := false
		quicksilver.InterchainQueryKeeper.IterateQueries(ctx, func(index int64, queryInfo icqtypes.Query) (stop bool) {
			if queryInfo.ChainId == zone.ChainId &&
				queryInfo.ConnectionId == zone.ConnectionId &&
				queryInfo.QueryType == icstypes.BankStoreKey &&
				bytes.Equal(queryInfo.Request, append(data, []byte("uqck")...)) {
				found = true
				return true
			}
			return false
		})
		suite.True(found)
	})
}

func (suite *KeeperTestSuite) TestAllBalancesCallbackMulti() {
	suite.Run("all balances non-zero)", func() {
		suite.SetupTest()
		suite.setupTestZones()

		quicksilver := suite.GetQuicksilverApp(suite.chainA)
		quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := suite.chainA.GetContext()

		zone, _ := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)

		queryReq := banktypes.QueryAllBalancesRequest{
			Address: zone.DepositAddress.Address,
		}
		reqbz, err := quicksilver.AppCodec().Marshal(&queryReq)
		suite.NoError(err)

		response := banktypes.QueryAllBalancesResponse{Balances: sdk.NewCoins(sdk.NewCoin("uqck", sdk.OneInt()), sdk.NewCoin("stake", sdk.OneInt()))}
		respbz, err := quicksilver.AppCodec().Marshal(&response)
		suite.NoError(err)

		err = keeper.AllBalancesCallback(quicksilver.InterchainstakingKeeper, ctx, respbz, icqtypes.Query{ChainId: suite.chainB.ChainID, Request: reqbz})
		suite.NoError(err)

		// refetch zone
		zone, _ = quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
		suite.Equal(uint32(2), zone.DepositAddress.BalanceWaitgroup)

		_, addr, err := bech32.DecodeAndConvert(zone.DepositAddress.Address)
		suite.NoError(err)
		data := banktypes.CreateAccountBalancesPrefix(addr)

		// check a ICQ request was made for each denom
		for _, coin := range response.Balances {
			found := false
			quicksilver.InterchainQueryKeeper.IterateQueries(ctx, func(index int64, queryInfo icqtypes.Query) (stop bool) {
				if queryInfo.ChainId == zone.ChainId &&
					queryInfo.ConnectionId == zone.ConnectionId &&
					queryInfo.QueryType == icstypes.BankStoreKey &&
					bytes.Equal(queryInfo.Request, append(data, []byte(coin.GetDenom())...)) {
					found = true
					return true
				}
				return false
			})
			suite.True(found)
		}
	})
}

func (suite *KeeperTestSuite) TestAccountBalanceCallback() {
	suite.Run("account balance", func() {
		suite.SetupTest()
		suite.setupTestZones()

		quicksilver := suite.GetQuicksilverApp(suite.chainA)
		quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := suite.chainA.GetContext()

		zone, _ := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
		zone.DepositAddress.IncrementBalanceWaitgroup()
		zone.WithdrawalAddress.IncrementBalanceWaitgroup()
		quicksilver.InterchainstakingKeeper.SetZone(ctx, &zone)

		response := sdk.NewCoin("qck", sdk.NewInt(10))
		respbz, err := quicksilver.AppCodec().Marshal(&response)
		suite.NoError(err)

		for _, addr := range []string{zone.DepositAddress.Address, zone.WithdrawalAddress.Address} {
			accAddr, err := sdk.AccAddressFromBech32(addr)
			suite.NoError(err)
			data := append(banktypes.CreateAccountBalancesPrefix(accAddr), []byte("qck")...)

			err = keeper.AccountBalanceCallback(quicksilver.InterchainstakingKeeper, ctx, respbz, icqtypes.Query{ChainId: suite.chainB.ChainID, Request: data})
			suite.NoError(err)
		}
	})
}

func (suite *KeeperTestSuite) TestAccountBalance046Callback() {
	suite.Run("account balance", func() {
		suite.SetupTest()
		suite.setupTestZones()

		quicksilver := suite.GetQuicksilverApp(suite.chainA)
		quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := suite.chainA.GetContext()

		zone, _ := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
		zone.DepositAddress.IncrementBalanceWaitgroup()
		zone.WithdrawalAddress.IncrementBalanceWaitgroup()
		quicksilver.InterchainstakingKeeper.SetZone(ctx, &zone)

		response := sdk.NewInt(10)

		respbz, err := response.Marshal()
		suite.NoError(err)

		for _, addr := range []string{zone.DepositAddress.Address, zone.WithdrawalAddress.Address} {
			accAddr, err := sdk.AccAddressFromBech32(addr)
			suite.NoError(err)
			data := append(banktypes.CreateAccountBalancesPrefix(accAddr), []byte("qck")...)

			err = keeper.AccountBalanceCallback(quicksilver.InterchainstakingKeeper, ctx, respbz, icqtypes.Query{ChainId: suite.chainB.ChainID, Request: data})
			suite.NoError(err)
		}
	})
}

func (suite *KeeperTestSuite) TestAccountBalanceCallbackMismatch() {
	suite.Run("account balance", func() {
		suite.SetupTest()
		suite.setupTestZones()

		quicksilver := suite.GetQuicksilverApp(suite.chainA)
		quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := suite.chainA.GetContext()

		zone, _ := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
		zone.DepositAddress.IncrementBalanceWaitgroup()
		zone.WithdrawalAddress.IncrementBalanceWaitgroup()
		quicksilver.InterchainstakingKeeper.SetZone(ctx, &zone)

		response := sdk.NewCoin("qck", sdk.NewInt(10))
		respbz, err := quicksilver.AppCodec().Marshal(&response)
		suite.NoError(err)

		for _, addr := range []string{zone.DepositAddress.Address, zone.WithdrawalAddress.Address} {
			accAddr, err := sdk.AccAddressFromBech32(addr)
			suite.NoError(err)
			data := append(banktypes.CreateAccountBalancesPrefix(accAddr), []byte("stake")...)

			err = keeper.AccountBalanceCallback(quicksilver.InterchainstakingKeeper, ctx, respbz, icqtypes.Query{ChainId: suite.chainB.ChainID, Request: data})
			suite.ErrorContains(err, "received coin denom qck does not match requested denom stake")
		}
	})
}

func (suite *KeeperTestSuite) TestAccountBalanceCallbackNil() {
	suite.Run("account balance", func() {
		suite.SetupTest()
		suite.setupTestZones()

		quicksilver := suite.GetQuicksilverApp(suite.chainA)
		quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := suite.chainA.GetContext()

		zone, _ := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
		zone.DepositAddress.IncrementBalanceWaitgroup()
		zone.WithdrawalAddress.IncrementBalanceWaitgroup()
		quicksilver.InterchainstakingKeeper.SetZone(ctx, &zone)

		var response *sdk.Coin
		respbz, err := quicksilver.AppCodec().Marshal(response)
		suite.NoError(err)

		for _, addr := range []string{zone.DepositAddress.Address, zone.WithdrawalAddress.Address} {
			accAddr, err := sdk.AccAddressFromBech32(addr)
			suite.NoError(err)
			data := append(banktypes.CreateAccountBalancesPrefix(accAddr), []byte("stake")...)

			err = keeper.AccountBalanceCallback(quicksilver.InterchainstakingKeeper, ctx, respbz, icqtypes.Query{ChainId: suite.chainB.ChainID, Request: data})
			suite.NoError(err)
		}
	})
}

// Ensures that a fuzz vector which resulted in a crash of ValidatorReq.Pagination crashing
// doesn't creep back up. Please see https://github.com/quicksilver-zone/quicksilver-incognito/issues/82
func TestValsetCallbackNilValidatorReqPagination(t *testing.T) {
	suite := new(KeeperTestSuite)
	suite.SetT(t)
	suite.SetupTest()
	suite.setupTestZones()

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
	ctx := suite.chainA.GetContext()

	data := []byte("\x12\"\n 00000000000000000000000000000000")
	err := keeper.ValsetCallback(quicksilver.InterchainstakingKeeper, ctx, data, icqtypes.Query{ChainId: suite.chainB.ChainID})
	suite.NoError(err)
}

func TestDelegationsCallbackAllPresentNoChange(t *testing.T) {
	suite := new(KeeperTestSuite)
	suite.SetT(t)
	suite.SetupTest()
	suite.setupTestZones()

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
	ctx := suite.chainA.GetContext()
	cdc := quicksilver.IBCKeeper.Codec()

	zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.True(found)

	vals := suite.GetQuicksilverApp(suite.chainB).StakingKeeper.GetAllValidators(suite.chainB.GetContext())
	delegationA := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[0].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationB := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[1].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationC := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[2].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}

	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationA)
	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationB)
	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationC)

	response := stakingtypes.QueryDelegatorDelegationsResponse{DelegationResponses: []stakingtypes.DelegationResponse{
		{Delegation: stakingtypes.Delegation{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[0].OperatorAddress, Shares: sdk.NewDec(1000)}, Balance: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))},
		{Delegation: stakingtypes.Delegation{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[1].OperatorAddress, Shares: sdk.NewDec(1000)}, Balance: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))},
		{Delegation: stakingtypes.Delegation{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[2].OperatorAddress, Shares: sdk.NewDec(1000)}, Balance: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))},
	}}

	data := cdc.MustMarshal(&response)

	delegationQuery := stakingtypes.QueryDelegatorDelegationsRequest{DelegatorAddr: zone.DelegationAddress.Address, Pagination: &query.PageRequest{Limit: uint64(len(zone.Validators))}}
	bz := cdc.MustMarshal(&delegationQuery)

	err := keeper.DelegationsCallback(quicksilver.InterchainstakingKeeper, ctx, data, icqtypes.Query{ChainId: suite.chainB.ChainID, Request: bz})

	suite.NoError(err)

	delegationRequests := 0
	for _, query := range quicksilver.InterchainQueryKeeper.AllQueries(ctx) {
		if query.CallbackId == delegationQueryCallbackID {
			delegationRequests++
		}
	}

	suite.Equal(0, delegationRequests)
	suite.Equal(3, len(quicksilver.InterchainstakingKeeper.GetAllDelegations(ctx, &zone)))
}

func TestDelegationsCallbackAllPresentOneChange(t *testing.T) {
	suite := new(KeeperTestSuite)
	suite.SetT(t)
	suite.SetupTest()
	suite.setupTestZones()

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
	ctx := suite.chainA.GetContext()
	cdc := quicksilver.IBCKeeper.Codec()

	zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.True(found)

	vals := suite.GetQuicksilverApp(suite.chainB).StakingKeeper.GetAllValidators(suite.chainB.GetContext())
	delegationA := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[0].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationB := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[1].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationC := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[2].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}

	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationA)
	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationB)
	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationC)

	response := stakingtypes.QueryDelegatorDelegationsResponse{DelegationResponses: []stakingtypes.DelegationResponse{
		{Delegation: stakingtypes.Delegation{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[0].OperatorAddress, Shares: sdk.NewDec(1000)}, Balance: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))},
		{Delegation: stakingtypes.Delegation{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[1].OperatorAddress, Shares: sdk.NewDec(2000)}, Balance: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(2000))},
		{Delegation: stakingtypes.Delegation{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[2].OperatorAddress, Shares: sdk.NewDec(1000)}, Balance: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))},
	}}

	data := cdc.MustMarshal(&response)

	delegationQuery := stakingtypes.QueryDelegatorDelegationsRequest{DelegatorAddr: zone.DelegationAddress.Address, Pagination: &query.PageRequest{Limit: uint64(len(zone.Validators))}}
	bz := cdc.MustMarshal(&delegationQuery)

	err := keeper.DelegationsCallback(quicksilver.InterchainstakingKeeper, ctx, data, icqtypes.Query{ChainId: suite.chainB.ChainID, Request: bz})

	suite.NoError(err)

	delegationRequests := 0
	for _, query := range quicksilver.InterchainQueryKeeper.AllQueries(ctx) {
		if query.CallbackId == delegationQueryCallbackID {
			delegationRequests++
		}
	}

	suite.Equal(1, delegationRequests)
	suite.Equal(3, len(quicksilver.InterchainstakingKeeper.GetAllDelegations(ctx, &zone)))
}

func TestDelegationsCallbackOneMissing(t *testing.T) {
	suite := new(KeeperTestSuite)
	suite.SetT(t)
	suite.SetupTest()
	suite.setupTestZones()

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
	ctx := suite.chainA.GetContext()
	cdc := quicksilver.IBCKeeper.Codec()

	zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.True(found)

	vals := suite.GetQuicksilverApp(suite.chainB).StakingKeeper.GetAllValidators(suite.chainB.GetContext())
	delegationA := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[0].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationB := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[1].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationC := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[2].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}

	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationA)
	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationB)
	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationC)

	response := stakingtypes.QueryDelegatorDelegationsResponse{DelegationResponses: []stakingtypes.DelegationResponse{
		{Delegation: stakingtypes.Delegation{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[0].OperatorAddress, Shares: sdk.NewDec(1000)}, Balance: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))},
		{Delegation: stakingtypes.Delegation{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[1].OperatorAddress, Shares: sdk.NewDec(1000)}, Balance: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))},
	}}

	data := cdc.MustMarshal(&response)

	delegationQuery := stakingtypes.QueryDelegatorDelegationsRequest{DelegatorAddr: zone.DelegationAddress.Address, Pagination: &query.PageRequest{Limit: uint64(len(zone.Validators))}}
	bz := cdc.MustMarshal(&delegationQuery)

	err := keeper.DelegationsCallback(quicksilver.InterchainstakingKeeper, ctx, data, icqtypes.Query{ChainId: suite.chainB.ChainID, Request: bz})

	suite.NoError(err)

	delegationRequests := 0
	for _, query := range quicksilver.InterchainQueryKeeper.AllQueries(ctx) {
		if query.CallbackId == delegationQueryCallbackID {
			delegationRequests++
		}
	}

	suite.Equal(1, delegationRequests)                                                     // callback for 'missing' delegation.
	suite.Equal(3, len(quicksilver.InterchainstakingKeeper.GetAllDelegations(ctx, &zone))) // new delegation doesn't get removed until the callback.
}

func TestDelegationsCallbackOneAdditional(t *testing.T) {
	suite := new(KeeperTestSuite)
	suite.SetT(t)
	suite.SetupTest()
	suite.setupTestZones()

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
	ctx := suite.chainA.GetContext()
	cdc := quicksilver.IBCKeeper.Codec()

	zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.True(found)

	vals := suite.GetQuicksilverApp(suite.chainB).StakingKeeper.GetAllValidators(suite.chainB.GetContext())
	delegationA := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[0].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationB := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[1].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationC := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[2].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}

	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationA)
	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationB)
	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationC)

	response := stakingtypes.QueryDelegatorDelegationsResponse{DelegationResponses: []stakingtypes.DelegationResponse{
		{Delegation: stakingtypes.Delegation{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[0].OperatorAddress, Shares: sdk.NewDec(1000)}, Balance: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))},
		{Delegation: stakingtypes.Delegation{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[1].OperatorAddress, Shares: sdk.NewDec(1000)}, Balance: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))},
		{Delegation: stakingtypes.Delegation{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[2].OperatorAddress, Shares: sdk.NewDec(1000)}, Balance: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))},
		{Delegation: stakingtypes.Delegation{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[3].OperatorAddress, Shares: sdk.NewDec(1000)}, Balance: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))},
	}}

	data := cdc.MustMarshal(&response)

	delegationQuery := stakingtypes.QueryDelegatorDelegationsRequest{DelegatorAddr: zone.DelegationAddress.Address, Pagination: &query.PageRequest{Limit: uint64(len(zone.Validators))}}
	bz := cdc.MustMarshal(&delegationQuery)

	err := keeper.DelegationsCallback(quicksilver.InterchainstakingKeeper, ctx, data, icqtypes.Query{ChainId: suite.chainB.ChainID, Request: bz})

	suite.NoError(err)

	delegationRequests := 0
	for _, query := range quicksilver.InterchainQueryKeeper.AllQueries(ctx) {
		if query.CallbackId == delegationQueryCallbackID {
			delegationRequests++
		}
	}

	suite.Equal(1, delegationRequests)
	suite.Equal(3, len(quicksilver.InterchainstakingKeeper.GetAllDelegations(ctx, &zone))) // new delegation doesn't get added until the end
}

func TestDelegationCallbackNew(t *testing.T) {
	suite := new(KeeperTestSuite)
	suite.SetT(t)
	suite.SetupTest()
	suite.setupTestZones()

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
	ctx := suite.chainA.GetContext()
	cdc := quicksilver.IBCKeeper.Codec()

	zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.True(found)

	vals := suite.GetQuicksilverApp(suite.chainB).StakingKeeper.GetAllValidators(suite.chainB.GetContext())
	delegationA := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[0].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationB := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[1].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationC := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[2].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}

	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationA)
	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationB)
	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationC)

	response := stakingtypes.Delegation{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[3].OperatorAddress, Shares: sdk.NewDec(1000)}

	data := cdc.MustMarshal(&response)

	delAddr, err := addressutils.AccAddressFromBech32(zone.DelegationAddress.Address, "")
	suite.NoError(err)
	valAddr, err := addressutils.ValAddressFromBech32(vals[3].OperatorAddress, "")
	suite.NoError(err)
	bz := stakingtypes.GetDelegationKey(delAddr, valAddr)

	err = keeper.DelegationCallback(quicksilver.InterchainstakingKeeper, ctx, data, icqtypes.Query{ChainId: suite.chainB.ChainID, Request: bz})
	suite.NoError(err)

	suite.Equal(4, len(quicksilver.InterchainstakingKeeper.GetAllDelegations(ctx, &zone)))
}

func TestDelegationCallbackUpdate(t *testing.T) {
	suite := new(KeeperTestSuite)
	suite.SetT(t)
	suite.SetupTest()
	suite.setupTestZones()

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
	ctx := suite.chainA.GetContext()
	cdc := quicksilver.IBCKeeper.Codec()

	zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.True(found)

	vals := suite.GetQuicksilverApp(suite.chainB).StakingKeeper.GetAllValidators(suite.chainB.GetContext())
	delegationA := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[0].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationB := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[1].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationC := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[2].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}

	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationA)
	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationB)
	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationC)

	response := stakingtypes.Delegation{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[2].OperatorAddress, Shares: sdk.NewDec(2000)}

	data := cdc.MustMarshal(&response)

	delAddr, err := addressutils.AccAddressFromBech32(zone.DelegationAddress.Address, "")
	suite.NoError(err)
	valAddr, err := addressutils.ValAddressFromBech32(vals[3].OperatorAddress, "")
	suite.NoError(err)
	bz := stakingtypes.GetDelegationKey(delAddr, valAddr)

	err = keeper.DelegationCallback(quicksilver.InterchainstakingKeeper, ctx, data, icqtypes.Query{ChainId: suite.chainB.ChainID, Request: bz})
	suite.NoError(err)

	suite.Equal(3, len(quicksilver.InterchainstakingKeeper.GetAllDelegations(ctx, &zone)))
}

func TestDelegationCallbackNoOp(t *testing.T) {
	suite := new(KeeperTestSuite)
	suite.SetT(t)
	suite.SetupTest()
	suite.setupTestZones()

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
	ctx := suite.chainA.GetContext()
	cdc := quicksilver.IBCKeeper.Codec()

	zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.True(found)

	vals := suite.GetQuicksilverApp(suite.chainB).StakingKeeper.GetAllValidators(suite.chainB.GetContext())
	delegationA := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[0].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationB := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[1].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationC := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[2].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}

	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationA)
	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationB)
	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationC)

	response := stakingtypes.Delegation{DelegatorAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[2].OperatorAddress, Shares: sdk.NewDec(1000)}

	data := cdc.MustMarshal(&response)

	delAddr, err := addressutils.AccAddressFromBech32(zone.DelegationAddress.Address, "")
	suite.NoError(err)
	valAddr, err := addressutils.ValAddressFromBech32(vals[3].OperatorAddress, "")
	suite.NoError(err)
	bz := stakingtypes.GetDelegationKey(delAddr, valAddr)
	ctx = suite.chainA.GetContext()
	err = keeper.DelegationCallback(quicksilver.InterchainstakingKeeper, ctx, data, icqtypes.Query{ChainId: suite.chainB.ChainID, Request: bz})
	suite.NoError(err)

	suite.Equal(3, len(quicksilver.InterchainstakingKeeper.GetAllDelegations(ctx, &zone)))
}

func TestDelegationCallbackRemove(t *testing.T) {
	suite := new(KeeperTestSuite)
	suite.SetT(t)
	suite.SetupTest()
	suite.setupTestZones()

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
	ctx := suite.chainA.GetContext()
	cdc := quicksilver.IBCKeeper.Codec()

	zone, found := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
	suite.True(found)

	vals := suite.GetQuicksilverApp(suite.chainB).StakingKeeper.GetAllValidators(suite.chainB.GetContext())
	delegationA := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[0].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationB := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[1].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}
	delegationC := icstypes.Delegation{DelegationAddress: zone.DelegationAddress.Address, ValidatorAddress: vals[2].OperatorAddress, Amount: sdk.NewCoin(zone.BaseDenom, sdk.NewInt(1000))}

	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationA)
	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationB)
	quicksilver.InterchainstakingKeeper.SetDelegation(ctx, &zone, delegationC)

	response := stakingtypes.Delegation{}

	data := cdc.MustMarshal(&response)

	delAddr, err := addressutils.AccAddressFromBech32(zone.DelegationAddress.Address, "")
	suite.NoError(err)
	valAddr, err := addressutils.ValAddressFromBech32(vals[3].OperatorAddress, "")
	suite.NoError(err)
	bz := stakingtypes.GetDelegationKey(delAddr, valAddr)

	err = keeper.DelegationCallback(quicksilver.InterchainstakingKeeper, ctx, data, icqtypes.Query{ChainId: suite.chainB.ChainID, Request: bz})
	suite.NoError(err)

	delegationRequests := 0
	for _, query := range quicksilver.InterchainQueryKeeper.AllQueries(ctx) {
		if query.CallbackId == delegationQueryCallbackID {
			delegationRequests++
		}
	}

	suite.Equal(3, len(quicksilver.InterchainstakingKeeper.GetAllDelegations(ctx, &zone)))
}

func TestDepositIntervalCallback(t *testing.T) {
	suite := new(KeeperTestSuite)
	suite.SetT(t)
	suite.SetupTest()
	suite.setupTestZones()

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
	ctx := suite.chainA.GetContext()

	data, err := base64.StdEncoding.DecodeString(depositTxFixture)
	res := tx.GetTxsEventResponse{}
	quicksilver.InterchainQueryKeeper.IBCKeeper.Codec().MustUnmarshal(data, &res)
	suite.NoError(err)

	err = keeper.DepositIntervalCallback(quicksilver.InterchainstakingKeeper, ctx, data, icqtypes.Query{ChainId: suite.chainB.ChainID})
	suite.NoError(err)
	txQueryCount := 0
	for _, query := range quicksilver.InterchainQueryKeeper.AllQueries(ctx) {
		if query.QueryType == "tendermint.Tx" {
			txQueryCount++
		}
	}
	// check we have the correct number (29) tendermint.Tx ICQ requests from the above payload.
	suite.Equal(int(res.Pagination.Total), txQueryCount)
}

func TestDepositIntervalCallbackWithExistingTxs(t *testing.T) {
	suite := new(KeeperTestSuite)
	suite.SetT(t)
	suite.SetupTest()
	suite.setupTestZones()

	quicksilver := suite.GetQuicksilverApp(suite.chainA)
	quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
	ctx := suite.chainA.GetContext()

	data, err := base64.StdEncoding.DecodeString(depositTxFixture)
	res := tx.GetTxsEventResponse{}
	quicksilver.InterchainQueryKeeper.IBCKeeper.Codec().MustUnmarshal(data, &res)
	suite.NoError(err)
	var msg banktypes.MsgSend
	_ = quicksilver.InterchainQueryKeeper.IBCKeeper.Codec().UnpackAny(res.TxResponses[0].Tx, msg)

	msgA := msg
	txrA := res.TxResponses[0]
	quicksilver.InterchainstakingKeeper.SetReceipt(ctx, icstypes.Receipt{ChainId: suite.chainB.ChainID, Sender: msgA.FromAddress, Txhash: txrA.TxHash, Amount: msgA.Amount})

	_ = quicksilver.InterchainQueryKeeper.IBCKeeper.Codec().UnpackAny(res.TxResponses[1].Tx, msg)

	msgB := msg
	txrB := res.TxResponses[1]
	quicksilver.InterchainstakingKeeper.SetReceipt(ctx, icstypes.Receipt{ChainId: suite.chainB.ChainID, Sender: msgB.FromAddress, Txhash: txrB.TxHash, Amount: msgB.Amount})

	_ = quicksilver.InterchainQueryKeeper.IBCKeeper.Codec().UnpackAny(res.TxResponses[2].Tx, msg)

	msgC := msg
	txrC := res.TxResponses[2]
	quicksilver.InterchainstakingKeeper.SetReceipt(ctx, icstypes.Receipt{ChainId: suite.chainB.ChainID, Sender: msgC.FromAddress, Txhash: txrC.TxHash, Amount: msgC.Amount})

	err = keeper.DepositIntervalCallback(quicksilver.InterchainstakingKeeper, ctx, data, icqtypes.Query{ChainId: suite.chainB.ChainID})
	suite.NoError(err)
	txQueryCount := 0
	for _, query := range quicksilver.InterchainQueryKeeper.AllQueries(ctx) {
		if query.QueryType == "tendermint.Tx" {
			txQueryCount++
		}
	}
	// check we have the correct number (29 minus - 3 receipts = 26) tendermint.Tx ICQ requests from the above payload.
	suite.Equal(int(res.Pagination.Total)-3, txQueryCount)
}

func (suite *KeeperTestSuite) TestDelegationAccountBalanceCallback() {
	suite.Run("account balance", func() {
		suite.SetupTest()
		suite.setupTestZones()

		quicksilver := suite.GetQuicksilverApp(suite.chainA)
		quicksilver.InterchainstakingKeeper.CallbackHandler().RegisterCallbacks()
		ctx := suite.chainA.GetContext()

		zone, _ := quicksilver.InterchainstakingKeeper.GetZone(ctx, suite.chainB.ChainID)
		zone.DepositAddress.IncrementBalanceWaitgroup()
		zone.WithdrawalAddress.IncrementBalanceWaitgroup()
		quicksilver.InterchainstakingKeeper.SetZone(ctx, &zone)

		response := sdk.NewCoin("qck", sdk.NewInt(10))
		respbz, err := quicksilver.AppCodec().Marshal(&response)
		suite.NoError(err)

		delAddr := zone.DelegationAddress.Address

		accAddr, err := addressutils.AccAddressFromBech32(delAddr, "cosmos")
		suite.NoError(err)

		data := append(banktypes.CreateAccountBalancesPrefix(accAddr), []byte("qck")...)

		err = keeper.DelegationAccountBalanceCallback(quicksilver.InterchainstakingKeeper, ctx, respbz, icqtypes.Query{ChainId: suite.chainB.ChainID, Request: data})

		suite.NoError(err)
	})
}

// keep depositTxFixture at the foot of the file, so it's not in the way!
var depositTxFixture = "CtgCCqkBCqYBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEoUBCixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhoTCgZ1c3RhcnMSCTEwMDAwMDAwMBJoClEKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiECIKszyzzlrimoxi4OygKjfL4J70aNA7PbGYBYKXZ4+9sSBAoCCH8YmQESEwoNCgZ1c3RhcnMSAzg3OBCtrQUaQElgAsjyeRVdZzUjUhdBHh7YZbn4eGgdUIF6g0shVueYZkgLNust5TNRieP3QyqhN46Bevs45oDqYM9NA1hnbHAKxgMKmQIKpAEKHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQSgwEKLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrEkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGhEKBnVzdGFycxIHNTAwMDAwMBJwTWdDaUtUNVJwWFJpT1RrNEhZWm1pZnp5VnhUZ01nWlB6aWFHMFhGRmF2V1lqT1pYNjBySDZMUEFNZ3NNYi9RSzJiS1NtbXBrV0JORWsyMVB5dldOTWdrNWtjb2pWdUxNK0JSUi9ENDRBNkFOL3NKWRJmCk4KRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiECS/kk9fP+wtzFOPU1Kt81wEQW0sIod7vo4A7TJYfspkISBAoCCH8SFAoOCgZ1c3RhcnMSBDUwMDAQwJoMGkB9SOr5/nZ/L1x+Tx4RACxndhtvzTt0PR85lYiHXRntvmZBDkJvbirHPBjMEllXkQR8R7snwGFBjbXFbBn70SC5CqAECvACCqYBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEoUBCixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhoTCgZ1c3RhcnMSCTUwMDAwMDAwMBLEAUhBWlB6aWFHMFhGRmF2V1lqT1pYNjBySDZMUEFIQUVPemtSNlJqY1NOa2o5NnE4RTRYeDBUdmluSEFNVXRzdDk4M0NOUzJqcm5OZjNzcFlNY293WUhCZzN2cHJDSkdYSkk2dWFuZC9lY2tsVUltZmJIQjY1QXlEYkdMUkp5RW9sTWlwdzY0UlNIYS9qSEIydUkwUVVMTVd4eUhLUDdoNkpKdW1ac3ZEZkhCMWhsaEdjMlVJUWZkNmk5WGZreFpWa1ZvaU4SaQpRCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAiCrM8s85a4pqMYuDsoCo3y+Ce9GjQOz2xmAWCl2ePvbEgQKAgh/GJwBEhQKDgoGdXN0YXJzEgQyMDAwEMCaDBpA2ev2D9ExtmR2V6z5HZQpVLIpoEpARnIf839V1TfYK4ZrjBjsU6yWAj4DOM44Gjk+uq4uUUdCpZ3gJQ93lgMUqQrzAgrGAQqlAQocL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBKEAQosc3RhcnMxZjZnOWd1eWV5emd6amM5bDh3ZzR4bDV4MHJ2eGRkZXdkcWp2N3YSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24aEgoGdXN0YXJzEggyMDAwMDAwMBIceUdISUdocVdUazFYWkNvS2VjSTdSOFRWK0M0MhJmCk4KRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiEDl4VfAbna07ch0CHtt6W3h6bErwHQsTASBwUShp8TOewSBAoCCH8SFAoOCgZ1c3RhcnMSBDUwMDAQwJoMGkBZEGK6gbUXO6d1pZRcvi3m9lb33cTw+AoMtDqP8WnWjFSXv+KdK6EWEQj1yYC7o/rJe5+dTwCAbuijciewNWApCukDCrkCCqcBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEoYBCixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhoUCgZ1c3RhcnMSCjE1MDAwMDAwMDASjAFLQVpQemlhRzBYRkZhdldZak9aWDYwckg2TFBBS0F5M05BM1loMWNPRUlmY1FENVB2SjY5QVpUZUtDTDNua0NaS1VUUmMzWXlycythZ1lwa1Z0SEhLQ3BlNHNHREsraWI1VHlCRFZFZUNUeGxySzNvS0RENGYybm5ZME1GUkJXSGl2ZnIrSWJldEpzZhJpClEKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiECIKszyzzlrimoxi4OygKjfL4J70aNA7PbGYBYKXZ4+9sSBAoCCH8YnQESFAoOCgZ1c3RhcnMSBDIwMDAQwJoMGkAGPkcAPFYfAiXW+VJ5a9/e6pMpV/PA38uGuqnelI6ZMn2u4WT79slE1abCINSdeWWrM56xiZ8qVBQLmpr/MXhdCp0ECu4CCqQBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEoMBCixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhoRCgZ1c3RhcnMSBzEwMDAwMDASxAFIQUNpS1Q1UnBYUmlPVGs0SFlabWlmenlWeFRnSEF4UjJuQ0MrVy9OdmxRREhwb3drQ1pZeDJEakhBazVrY29qVnVMTStCUlIvRDQ0QTZBTi9zSllIQUcxWkEvOTBWMnhGZXVjZkZEVzNlR29FYmJjSEF5M05BM1loMWNPRUlmY1FENVB2SjY5QVpUZUhCMWhsaEdjMlVJUWZkNmk5WGZreFpWa1ZvaU5IQjY1QXlEYkdMUkp5RW9sTWlwdzY0UlNIYS9qEmgKUApGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQJL+ST18/7C3MU49TUq3zXARBbSwih3u+jgDtMlh+ymQhIECgIIfxgBEhQKDgoGdXN0YXJzEgQ1MDAwEMCaDBpAgKoGTvvPuCkhiE+cl05IHS8826MatiXvLO7lqoY4gghDLH3GSQbbIThHkyJcxI4gYbmW3zdxfOR2ojI1kChyZgqsAwr9AQqkAQocL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBKDAQosc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24aEQoGdXN0YXJzEgcyMDAwMDAwElRRZ0VPemtSNlJqY1NOa2o5NnE4RTRYeDBUdmluUWdHMVpBLzkwVjJ4RmV1Y2ZGRFczZUdvRWJiY1Fnc01iL1FLMmJLU21tcGtXQk5FazIxUHl2V04SaApQCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAkv5JPXz/sLcxTj1NSrfNcBEFtLCKHe76OAO0yWH7KZCEgQKAgh/GAISFAoOCgZ1c3RhcnMSBDUwMDAQwJoMGkDkvqyEUs83s/QX+8FovUM3y+dndq2N5ZKoBC91aecleGFBIgtvof6IRXoLiVipWGUQGSL1QMmj8KpTBUjBvb3fCqwDCv0BCqQBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEoMBCixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhoRCgZ1c3RhcnMSBzEwMDAwMDASVFFnRU96a1I2UmpjU05rajk2cThFNFh4MFR2aW5RZ01VdHN0OTgzQ05TMmpybk5mM3NwWU1jb3dZUWd4UjJuQ0MrVy9OdmxRREhwb3drQ1pZeDJEahJoClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiECS/kk9fP+wtzFOPU1Kt81wEQW0sIod7vo4A7TJYfspkISBAoCCH8YAxIUCg4KBnVzdGFycxIENTAwMBDAmgwaQKwZk4AcLxb7jo3lVoUMjpLMwyWIc2VLuuft6H4fPFWQAqZjIOSUhjQO8HhhCSH2PLEJe3/nwy7NAvcXhqVefKUK9AIKxQEKpAEKHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQSgwEKLHN0YXJzMWY2ZzlndXlleXpnempjOWw4d2c0eGw1eDBydnhkZGV3ZHFqdjd2EkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGhEKBnVzdGFycxIHNzAwMDAwMBIceUFDaUtUNVJwWFJpT1RrNEhZWm1pZnp5VnhUZxJoClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiEDl4VfAbna07ch0CHtt6W3h6bErwHQsTASBwUShp8TOewSBAoCCH8YARIUCg4KBnVzdGFycxIENTAwMBDAmgwaQLtDxquZKz0SAf3FFz0M4fSm08zRgEOaXRoHlwmAt2rjNT+pYTclthlYTxzSB5Im9poORFxf4isVUDxxPtVWAeAKrAMK/QEKpAEKHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQSgwEKLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrEkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGhEKBnVzdGFycxIHMTUwMDAwMBJUUWdFT3prUjZSamNTTmtqOTZxOEU0WHgwVHZpblFnTVV0c3Q5ODNDTlMyanJuTmYzc3BZTWNvd1lRZ3NNYi9RSzJiS1NtbXBrV0JORWsyMVB5dldOEmgKUApGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQJL+ST18/7C3MU49TUq3zXARBbSwih3u+jgDtMlh+ymQhIECgIIfxgEEhQKDgoGdXN0YXJzEgQ1MDAwEMCaDBpASWwyaYuO37B8LpcL+FSkA+OHelV651eyGmuwTUliVINJr3RZJCSUEoaVsNvaOrSQstubidUOspXDoRhkWQVznAqFBArVAgqnAQocL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBKGAQosc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24aFAoGdXN0YXJzEgoxMDAwMDAwMDAwEqgBSVUzRk1iRWtZcWlWZWhrNFFPOFVpMDRkZXRKUElWTThjaXZqQ1JXZVg0S1MyNmhBZ2hSRi82SVBJVllBa2NqekRjUUZSZk5kVkNJWUQvYmZMcnpTSVYyZlQwTVdHVWlTMTYrYlo2MFRKMmxFaVo5UElXQmNXcXMyajFwOHRHa2x5RXpwallpRzRoemlJV0hJR2hxV1RrMVhaQ29LZWNJN1I4VFYrQzQyEmkKUQpGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQIgqzPLPOWuKajGLg7KAqN8vgnvRo0Ds9sZgFgpdnj72xIECgIIfxieARIUCg4KBnVzdGFycxIEMjAwMBDAmgwaQO0yEz+HtMpMj3bwA/+Z/LlnpNdtxbVkMzFM8LJX3HWFRQVD/iAP7LfptF6Y0IkZMlsd/U7fu3uuUN1mOzo/flAKrAMK/QEKpAEKHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQSgwEKLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrEkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGhEKBnVzdGFycxIHMzc1NzQ5MBJUUWdFT3prUjZSamNTTmtqOTZxOEU0WHgwVHZpblFnc01iL1FLMmJLU21tcGtXQk5FazIxUHl2V05RZ3hSMm5DQytXL052bFFESHBvd2tDWll4MkRqEmgKUApGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQJL+ST18/7C3MU49TUq3zXARBbSwih3u+jgDtMlh+ymQhIECgIIfxgFEhQKDgoGdXN0YXJzEgQ1MDAwEMCaDBpAfKidESTdgazihhU5qH0QNwxfES4G0oXtyNVASv0VhgEBWfZCUi9Htj8sTKb9amLFoBMt6IbCBo6iGCJhSMGUQgqQAwrhAQqkAQocL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBKDAQosc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24aEQoGdXN0YXJzEgcxNzY4NTQwEjhaQUVPemtSNlJqY1NOa2o5NnE4RTRYeDBUdmluWkFNVXRzdDk4M0NOUzJqcm5OZjNzcFlNY293WRJoClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiECS/kk9fP+wtzFOPU1Kt81wEQW0sIod7vo4A7TJYfspkISBAoCCH8YBhIUCg4KBnVzdGFycxIENTAwMBDAmgwaQB2zAkk+TzGcmQ+7hRlibIY61B/w0/3f5sIr7+Rf7kDifli4iQAQdO8VT9Q8941A+gH7oVKf7AcmhRqg+tfitCkK9AIKxQEKpAEKHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQSgwEKLHN0YXJzMXZ6OTJ3ZjY3a3NkbnNtY2pldWU2eDJ6anNmbGRwOWc5eThmcXk3EkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGhEKBnVzdGFycxIHMjAwMDAwMBIceUNSd1lZcyt0dXk1RmoxYmFxUVVhVWtRZFVVUBJoClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiECyaWQWoNXRbTeWTP0/toyUqhJDytSq23LR1+EqX9mVXQSBAoCCH8YARIUCg4KBnVzdGFycxIENTAwMBDAmgwaQHi83fAXkMvx4zJbvTGOH03wNx96oCz7rGAg1RP6vxgCZ1Z6wwpmKAgk560uzF7/R+LKhjT76fn+NGDRz7i9KrMKyAMKmQIKpAEKHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQSgwEKLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrEkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGhEKBnVzdGFycxIHMjAwMDAwMBJwTWdFT3prUjZSamNTTmtqOTZxOEU0WHgwVHZpbk1nTVV0c3Q5ODNDTlMyanJuTmYzc3BZTWNvd1lNZ3hSMm5DQytXL052bFFESHBvd2tDWll4MkRqTWdrNWtjb2pWdUxNK0JSUi9ENDRBNkFOL3NKWRJoClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiECS/kk9fP+wtzFOPU1Kt81wEQW0sIod7vo4A7TJYfspkISBAoCCH8YBxIUCg4KBnVzdGFycxIENTAwMBDAmgwaQHEoktRv5aSuK/JEbl4uxv8aetsFVaViK1eCjaYidQVXLMrbyu6T7mOwt3ILL9drKWpILfrhcwZpBVwu3grPaykKgQQK0gIKpAEKHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQSgwEKLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrEkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGhEKBnVzdGFycxIHMzAwMDAwMBKoAUlRRU96a1I2UmpjU05rajk2cThFNFh4MFR2aW5JUVpQemlhRzBYRkZhdldZak9aWDYwckg2TFBBSVFrNWtjb2pWdUxNK0JSUi9ENDRBNkFOL3NKWUlSTVRUQUdWSVNkcTBOZ0lBMGxCOVE4b2NNcjJJUmczdnByQ0pHWEpJNnVhbmQvZWNrbFVJbWZiSVI2NUF5RGJHTFJKeUVvbE1pcHc2NFJTSGEvahJoClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiECS/kk9fP+wtzFOPU1Kt81wEQW0sIod7vo4A7TJYfspkISBAoCCH8YCBIUCg4KBnVzdGFycxIENTAwMBDAmgwaQN1uUMSWW8LBqyDcT2/vK7n7jfKzEzWCGcC4aOrxsAKKCYGUcp9E8b78+HJhb975F5Ivv9AT1bvpUfmDSv5SRC0KuQQKigMKpAEKHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQSgwEKLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrEkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGhEKBnVzdGFycxIHMTAwMDAwMBLgAUdmWFVhTXZjclVGaGdmQ3grZmVOUUIxcThia1NHZjIxM2FHajl6UVUwbkhvSXZVbTVaMnU5amVBR2ZTV21LQjBpSmFKUy9qNTJXdXhrM2o1QmQyZ0dmTm5US3A2ZHA0UkVOa0hpOUNOOUI2a0pMdk9HZkNVU2MvN1BsaTVCRXlEb3FpVjNxVjlVZVBkR2VqbWRDMnhrcS9GNFZML01wcWZqZUppKzYra0dlZllqZjFyTmNjOTVSeW5vU2ZOUU1lTXhZNWxHZUc4eEJRM3FIUndvOWxYNkVuNThUZi9uRVRuEmgKUApGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQJL+ST18/7C3MU49TUq3zXARBbSwih3u+jgDtMlh+ymQhIECgIIfxgJEhQKDgoGdXN0YXJzEgQ1MDAwEMCaDBpArf/TzsLQ1o9oGq3A+Er8YE/F7IK4EYghN2T5YQfL+ztCfS/yAofujS8um+MIYYydCG0lXDvSn2F3ForSLJJjZQr0AgrFAQqkAQocL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBKDAQosc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24aEQoGdXN0YXJzEgcxMDAwMDAwEhx5QUNpS1Q1UnBYUmlPVGs0SFlabWlmenlWeFRnEmgKUApGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQJL+ST18/7C3MU49TUq3zXARBbSwih3u+jgDtMlh+ymQhIECgIIfxgKEhQKDgoGdXN0YXJzEgQ1MDAwEMCaDBpAX/m74idZXxB536AazZXWeYrE2Rin6wvU9ZH9Uw+SYmF7u8qg6bIlaLbPmM0frnZtVMNrztBzITc+gQLxj1+5ewr0AgrFAQqkAQocL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBKDAQosc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24aEQoGdXN0YXJzEgcxMDAwMDAwEhx5QUVXVXNMa2swT3N1SXk1TFRCb1lnV0RabnI1EmgKUApGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQJL+ST18/7C3MU49TUq3zXARBbSwih3u+jgDtMlh+ymQhIECgIIfxgLEhQKDgoGdXN0YXJzEgQ1MDAwEMCaDBpAjYRqlmiuIb/xb4KmZe/AO3Bj/m1dc0xYLeGjqgmurKkbQa5X3bDV5O1igDT4NnjdOpooUcrUVEHQdm4oAuhu9wr0AgrFAQqkAQocL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBKDAQosc3RhcnMxZjZnOWd1eWV5emd6amM5bDh3ZzR4bDV4MHJ2eGRkZXdkcWp2N3YSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24aEQoGdXN0YXJzEgcxMDAwMDAwEhx5QU1VdHN0OTgzQ05TMmpybk5mM3NwWU1jb3dZEmgKUApGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQOXhV8BudrTtyHQIe23pbeHpsSvAdCxMBIHBRKGnxM57BIECgIIfxgCEhQKDgoGdXN0YXJzEgQyMDAwEMCaDBpA8QvKmmxkxC2ZHQOZDzgM6oK7rtfNKUzt+91sZF4cJr4f3JmaxZigKuTXf20fU1Bi1mi2SkqlsH7CBtZ6DOmabAqvAwr/AQqmAQocL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBKFAQosc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24aEwoGdXN0YXJzEgkxMDAwMDAwMDASVFFnQ2lLVDVScFhSaU9UazRIWVptaWZ6eVZ4VGdRZ0VPemtSNlJqY1NOa2o5NnE4RTRYeDBUdmluUWdaUHppYUcwWEZGYXZXWWpPWlg2MHJINkxQQRJpClEKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiECIKszyzzlrimoxi4OygKjfL4J70aNA7PbGYBYKXZ4+9sSBAoCCH8YnwESFAoOCgZ1c3RhcnMSBDIwMDAQwJoMGkB3y1Ai+bOvEfswOf5xyIW814pqv4zzZLMUp8XkqPWhFjVvY1nGJ2tQBeeOQRHEXfdqeqRkSTcIXy65xdsKLsmYCq8DCv8BCqYBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEoUBCixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhoTCgZ1c3RhcnMSCTUwMDAwMDAwMBJUUWdNVXRzdDk4M0NOUzJqcm5OZjNzcFlNY293WVFnRU96a1I2UmpjU05rajk2cThFNFh4MFR2aW5RZ0NpS1Q1UnBYUmlPVGs0SFlabWlmenlWeFRnEmkKUQpGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQIgqzPLPOWuKajGLg7KAqN8vgnvRo0Ds9sZgFgpdnj72xIECgIIfxigARIUCg4KBnVzdGFycxIEMjAwMBDAmgwaQHFlgSYl1NJpA8gRdfEn/YG4a009bTfc1jehhMU6gtpOLFsgt2EJtgcqECLYfw/Mw3xYCUbj4JiXw9gkf6XeMFoKyQMKmgIKpQEKHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQShAEKLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrEkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGhIKBnVzdGFycxIIMTAwMDAwMDAScE1nQ2lLVDVScFhSaU9UazRIWVptaWZ6eVZ4VGdNZ01VdHN0OTgzQ05TMmpybk5mM3NwWU1jb3dZTWdFV1VzTGtrME9zdUl5NUxUQm9ZZ1dEWm5yNU1neFIybkNDK1cvTnZsUURIcG93a0NaWXgyRGoSaApQCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAkv5JPXz/sLcxTj1NSrfNcBEFtLCKHe76OAO0yWH7KZCEgQKAgh/GAwSFAoOCgZ1c3RhcnMSBDUwMDAQwJoMGkAooVTvbdztaLPbRHGzAYSzL3d4g+E80M/rBo7PdZywDzng7zEaobm00ZQmZobUjZTI16QJ1w32C14dZ3/MgHmICrkECooDCqQBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEoMBCixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhoRCgZ1c3RhcnMSBzUwMDAwMDAS4AFHUUNpS1Q1UnBYUmlPVGs0SFlabWlmenlWeFRnR1FNVXRzdDk4M0NOUzJqcm5OZjNzcFlNY293WUdReFIybkNDK1cvTnZsUURIcG93a0NaWXgyRGpHUXNNYi9RSzJiS1NtbXBrV0JORWsyMVB5dldOR1FrNWtjb2pWdUxNK0JSUi9ENDRBNkFOL3NKWUdRRzFaQS85MFYyeEZldWNmRkRXM2VHb0ViYmNHUUVPemtSNlJqY1NOa2o5NnE4RTRYeDBUdmluR1FFV1VzTGtrME9zdUl5NUxUQm9ZZ1dEWm5yNRJoClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiECS/kk9fP+wtzFOPU1Kt81wEQW0sIod7vo4A7TJYfspkISBAoCCH8YDRIUCg4KBnVzdGFycxIENTAwMBDAmgwaQHgpyuZv3WEpwaz44m5SpbcMeu4AgJfNmBJGxeb1lEZiCecGY9cDeQp+C9MhPm2yFiEOkHFcoEHIBEOyOMDHzQ0KrAMK/QEKpAEKHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQSgwEKLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrEkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGhEKBnVzdGFycxIHMTAwMDAwMBJUUWdFV1VzTGtrME9zdUl5NUxUQm9ZZ1dEWm5yNVFnTVV0c3Q5ODNDTlMyanJuTmYzc3BZTWNvd1lRZ3NNYi9RSzJiS1NtbXBrV0JORWsyMVB5dldOEmgKUApGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQJL+ST18/7C3MU49TUq3zXARBbSwih3u+jgDtMlh+ymQhIECgIIfxgOEhQKDgoGdXN0YXJzEgQ1MDAwEMCaDBpALKLOH8yygw0RUY5rS1DlVGDNKWbArRbqqvEiHpS4JMFL5Fnvs5wenCyNsx1LxqUGR5o1jPXxSbHWuyDLqsmBzwqwAwqAAgqnAQocL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBKGAQosc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24aFAoGdXN0YXJzEgoxMDAwMDAwMDAwElRRZ0VPemtSNlJqY1NOa2o5NnE4RTRYeDBUdmluUWdaUHppYUcwWEZGYXZXWWpPWlg2MHJINkxQQVFnc01iL1FLMmJLU21tcGtXQk5FazIxUHl2V04SaQpRCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAiCrM8s85a4pqMYuDsoCo3y+Ce9GjQOz2xmAWCl2ePvbEgQKAgh/GKEBEhQKDgoGdXN0YXJzEgQyMDAwEMCaDBpAGQcY+LdPb2CtxjNjWhKozbJVRyFjSIXxJ6cjUDPzm3gj/kBJOJvs7lnY3es5VlrQyhZobhWt5a+LfgoO+fJxGQrLAwqbAgqmAQocL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBKFAQosc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24aEwoGdXN0YXJzEgkyNTAwMDAwMDAScE1nTVV0c3Q5ODNDTlMyanJuTmYzc3BZTWNvd1lNZ1pQemlhRzBYRkZhdldZak9aWDYwckg2TFBBTWd4UjJuQ0MrVy9OdmxRREhwb3drQ1pZeDJEak1oZzN2cHJDSkdYSkk2dWFuZC9lY2tsVUltZmISaQpRCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAiCrM8s85a4pqMYuDsoCo3y+Ce9GjQOz2xmAWCl2ePvbEgQKAgh/GKIBEhQKDgoGdXN0YXJzEgQyMDAwEMCaDBpAc97zykfLc6m9n25wXZuELmQmDasv72o5Qu58JzRbRn4uuekjQGEuVK9b2NHKXzdftIhOCBwEZiQkGcbOZaydwAqRAwrhAQqkAQocL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBKDAQosc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24aEQoGdXN0YXJzEgcxMDAwMDAwEjhaQUVPemtSNlJqY1NOa2o5NnE4RTRYeDBUdmluWkFNVXRzdDk4M0NOUzJqcm5OZjNzcFlNY293WRJpClEKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiECIKszyzzlrimoxi4OygKjfL4J70aNA7PbGYBYKXZ4+9sSBAoCCH8YowESFAoOCgZ1c3RhcnMSBDIwMDAQwJoMGkA8yOC8B2WorzoEjn/r6kaWZqCmfGD8FDbKT43EHR9/+mgWIsqsAZlJ/xNDYjtoQjQLoNU+XedIiE2lCzvgh2TtCukDCrkCCqcBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEoYBCixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhoUCgZ1c3RhcnMSCjEyNzEwMDAwMDASjAFLQU1VdHN0OTgzQ05TMmpybk5mM3NwWU1jb3dZS0FzTWIvUUsyYktTbW1wa1dCTkVrMjFQeXZXTktCTVRUQUdWSVNkcTBOZ0lBMGxCOVE4b2NNcjJLQmczdnByQ0pHWEpJNnVhbmQvZWNrbFVJbWZiS0F4UjJuQ0MrVy9OdmxRREhwb3drQ1pZeDJEahJpClEKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiECIKszyzzlrimoxi4OygKjfL4J70aNA7PbGYBYKXZ4+9sSBAoCCH8YpAESFAoOCgZ1c3RhcnMSBDIwMDAQwJoMGkCMO+pyquEGUdm+bxmGGO/6dulN+1TX9dtS4xIfbCUd1XuOxrD5scCmesiJK/oOnb1Z69cY93KhDOfKjxM+/bkwEpYYCJSvgwMSQENBQkI4QTMzRkFBM0NCQUQ2NDYxREFGODAxM0NBQjU1N0YzOUIzNTBCMkRFM0ZENUYzQUQzNTcyQjVGMzBBOUUqQDBBMUUwQTFDMkY2MzZGNzM2RDZGNzMyRTYyNjE2RTZCMkU3NjMxNjI2NTc0NjEzMTJFNEQ3MzY3NTM2NTZFNjQyjwZbeyJldmVudHMiOlt7InR5cGUiOiJjb2luX3JlY2VpdmVkIiwiYXR0cmlidXRlcyI6W3sia2V5IjoicmVjZWl2ZXIiLCJ2YWx1ZSI6InN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24ifSx7ImtleSI6ImFtb3VudCIsInZhbHVlIjoiMTAwMDAwMDAwdXN0YXJzIn1dfSx7InR5cGUiOiJjb2luX3NwZW50IiwiYXR0cmlidXRlcyI6W3sia2V5Ijoic3BlbmRlciIsInZhbHVlIjoic3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgifSx7ImtleSI6ImFtb3VudCIsInZhbHVlIjoiMTAwMDAwMDAwdXN0YXJzIn1dfSx7InR5cGUiOiJtZXNzYWdlIiwiYXR0cmlidXRlcyI6W3sia2V5IjoiYWN0aW9uIiwidmFsdWUiOiIvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kIn0seyJrZXkiOiJzZW5kZXIiLCJ2YWx1ZSI6InN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4In0seyJrZXkiOiJtb2R1bGUiLCJ2YWx1ZSI6ImJhbmsifV19LHsidHlwZSI6InRyYW5zZmVyIiwiYXR0cmlidXRlcyI6W3sia2V5IjoicmVjaXBpZW50IiwidmFsdWUiOiJzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuIn0seyJrZXkiOiJzZW5kZXIiLCJ2YWx1ZSI6InN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4In0seyJrZXkiOiJhbW91bnQiLCJ2YWx1ZSI6IjEwMDAwMDAwMHVzdGFycyJ9XX1dfV06hgQaeAoNY29pbl9yZWNlaXZlZBJMCghyZWNlaXZlchJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhIZCgZhbW91bnQSDzEwMDAwMDAwMHVzdGFycxpgCgpjb2luX3NwZW50EjcKB3NwZW5kZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4EhkKBmFtb3VudBIPMTAwMDAwMDAwdXN0YXJzGnkKB21lc3NhZ2USJgoGYWN0aW9uEhwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEjYKBnNlbmRlchIsc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgSDgoGbW9kdWxlEgRiYW5rGqwBCgh0cmFuc2ZlchJNCglyZWNpcGllbnQSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24SNgoGc2VuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBIZCgZhbW91bnQSDzEwMDAwMDAwMHVzdGFyc0itrQVQh7EEWvICChUvY29zbW9zLnR4LnYxYmV0YTEuVHgS2AIKqQEKpgEKHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQShQEKLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4EkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGhMKBnVzdGFycxIJMTAwMDAwMDAwEmgKUQpGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQIgqzPLPOWuKajGLg7KAqN8vgnvRo0Ds9sZgFgpdnj72xIECgIIfxiZARITCg0KBnVzdGFycxIDODc4EK2tBRpASWACyPJ5FV1nNSNSF0EeHthlufh4aB1QgXqDSyFW55hmSAs26y3lM1GJ4/dDKqE3joF6+zjmgOpgz00DWGdscGIUMjAyMy0wMS0wN1QxMDoyNjoxNVpqXgoKY29pbl9zcGVudBI5CgdzcGVuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBgBEhUKBmFtb3VudBIJODc4dXN0YXJzGAFqYgoNY29pbl9yZWNlaXZlZBI6CghyZWNlaXZlchIsc3RhcnMxN3hwZnZha20yYW1nOTYyeWxzNmY4NHoza2VsbDhjNWx5OTVhcXYYARIVCgZhbW91bnQSCTg3OHVzdGFycxgBapgBCgh0cmFuc2ZlchI7CglyZWNpcGllbnQSLHN0YXJzMTd4cGZ2YWttMmFtZzk2MnlsczZmODR6M2tlbGw4YzVseTk1YXF2GAESOAoGc2VuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBgBEhUKBmFtb3VudBIJODc4dXN0YXJzGAFqQwoHbWVzc2FnZRI4CgZzZW5kZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4GAFqVQoCdHgSEgoDZmVlEgk4Nzh1c3RhcnMYARI7CglmZWVfcGF5ZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4GAFqQwoCdHgSPQoHYWNjX3NlcRIwc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgvMTUzGAFqbQoCdHgSZwoJc2lnbmF0dXJlElhTV0FDeVBKNUZWMW5OU05TRjBFZUh0aGx1Zmg0YUIxUWdYcURTeUZXNTVobVNBczI2eTNsTTFHSjQvZERLcUUzam9GNit6am1nT3BnejAwRFdHZHNjQT09GAFqMwoHbWVzc2FnZRIoCgZhY3Rpb24SHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQYAWpkCgpjb2luX3NwZW50EjkKB3NwZW5kZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4GAESGwoGYW1vdW50Eg8xMDAwMDAwMDB1c3RhcnMYAWp8Cg1jb2luX3JlY2VpdmVkEk4KCHJlY2VpdmVyEkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGAESGwoGYW1vdW50Eg8xMDAwMDAwMDB1c3RhcnMYAWqyAQoIdHJhbnNmZXISTwoJcmVjaXBpZW50EkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGAESOAoGc2VuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBgBEhsKBmFtb3VudBIPMTAwMDAwMDAwdXN0YXJzGAFqQwoHbWVzc2FnZRI4CgZzZW5kZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4GAFqGwoHbWVzc2FnZRIQCgZtb2R1bGUSBGJhbmsYARL0GAiitIMDEkBCNUMwMTFFREQ4QkQ3OUE3MzNGNUNCNEE4MUUyNDRFNUZGMjUyMzdBRTE1OUYyMjk2QjBCMUEzREVDMjFCNEI5KkAwQTFFMEExQzJGNjM2RjczNkQ2RjczMkU2MjYxNkU2QjJFNzYzMTYyNjU3NDYxMzEyRTRENzM2NzUzNjU2RTY0MokGW3siZXZlbnRzIjpbeyJ0eXBlIjoiY29pbl9yZWNlaXZlZCIsImF0dHJpYnV0ZXMiOlt7ImtleSI6InJlY2VpdmVyIiwidmFsdWUiOiJzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuIn0seyJrZXkiOiJhbW91bnQiLCJ2YWx1ZSI6IjUwMDAwMDB1c3RhcnMifV19LHsidHlwZSI6ImNvaW5fc3BlbnQiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJzcGVuZGVyIiwidmFsdWUiOiJzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5ayJ9LHsia2V5IjoiYW1vdW50IiwidmFsdWUiOiI1MDAwMDAwdXN0YXJzIn1dfSx7InR5cGUiOiJtZXNzYWdlIiwiYXR0cmlidXRlcyI6W3sia2V5IjoiYWN0aW9uIiwidmFsdWUiOiIvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kIn0seyJrZXkiOiJzZW5kZXIiLCJ2YWx1ZSI6InN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrIn0seyJrZXkiOiJtb2R1bGUiLCJ2YWx1ZSI6ImJhbmsifV19LHsidHlwZSI6InRyYW5zZmVyIiwiYXR0cmlidXRlcyI6W3sia2V5IjoicmVjaXBpZW50IiwidmFsdWUiOiJzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuIn0seyJrZXkiOiJzZW5kZXIiLCJ2YWx1ZSI6InN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrIn0seyJrZXkiOiJhbW91bnQiLCJ2YWx1ZSI6IjUwMDAwMDB1c3RhcnMifV19XX1dOoAEGnYKDWNvaW5fcmVjZWl2ZWQSTAoIcmVjZWl2ZXISQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24SFwoGYW1vdW50Eg01MDAwMDAwdXN0YXJzGl4KCmNvaW5fc3BlbnQSNwoHc3BlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsSFwoGYW1vdW50Eg01MDAwMDAwdXN0YXJzGnkKB21lc3NhZ2USJgoGYWN0aW9uEhwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEjYKBnNlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsSDgoGbW9kdWxlEgRiYW5rGqoBCgh0cmFuc2ZlchJNCglyZWNpcGllbnQSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24SNgoGc2VuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axIXCgZhbW91bnQSDTUwMDAwMDB1c3RhcnNIwJoMUPTqBFrgAwoVL2Nvc21vcy50eC52MWJldGExLlR4EsYDCpkCCqQBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEoMBCixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhoRCgZ1c3RhcnMSBzUwMDAwMDAScE1nQ2lLVDVScFhSaU9UazRIWVptaWZ6eVZ4VGdNZ1pQemlhRzBYRkZhdldZak9aWDYwckg2TFBBTWdzTWIvUUsyYktTbW1wa1dCTkVrMjFQeXZXTk1nazVrY29qVnVMTStCUlIvRDQ0QTZBTi9zSlkSZgpOCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAkv5JPXz/sLcxTj1NSrfNcBEFtLCKHe76OAO0yWH7KZCEgQKAgh/EhQKDgoGdXN0YXJzEgQ1MDAwEMCaDBpAfUjq+f52fy9cfk8eEQAsZ3Ybb807dD0fOZWIh10Z7b5mQQ5Cb24qxzwYzBJZV5EEfEe7J8BhQY21xWwZ+9EguWIUMjAyMy0wMS0wN1QxMTozMDo0OVpqXwoKY29pbl9zcGVudBI5CgdzcGVuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axgBEhYKBmFtb3VudBIKNTAwMHVzdGFycxgBamMKDWNvaW5fcmVjZWl2ZWQSOgoIcmVjZWl2ZXISLHN0YXJzMTd4cGZ2YWttMmFtZzk2MnlsczZmODR6M2tlbGw4YzVseTk1YXF2GAESFgoGYW1vdW50Ego1MDAwdXN0YXJzGAFqmQEKCHRyYW5zZmVyEjsKCXJlY2lwaWVudBIsc3RhcnMxN3hwZnZha20yYW1nOTYyeWxzNmY4NHoza2VsbDhjNWx5OTVhcXYYARI4CgZzZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrGAESFgoGYW1vdW50Ego1MDAwdXN0YXJzGAFqQwoHbWVzc2FnZRI4CgZzZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrGAFqVgoCdHgSEwoDZmVlEgo1MDAwdXN0YXJzGAESOwoJZmVlX3BheWVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axgBakEKAnR4EjsKB2FjY19zZXESLnN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrLzAYAWptCgJ0eBJnCglzaWduYXR1cmUSWGZVanErZjUyZnk5Y2ZrOGVFUUFzWjNZYmI4MDdkRDBmT1pXSWgxMFo3YjVtUVE1Q2IyNHF4endZekJKWlY1RUVmRWU3SjhCaFFZMjF4V3daKzlFZ3VRPT0YAWozCgdtZXNzYWdlEigKBmFjdGlvbhIcL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBgBamIKCmNvaW5fc3BlbnQSOQoHc3BlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYARIZCgZhbW91bnQSDTUwMDAwMDB1c3RhcnMYAWp6Cg1jb2luX3JlY2VpdmVkEk4KCHJlY2VpdmVyEkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGAESGQoGYW1vdW50Eg01MDAwMDAwdXN0YXJzGAFqsAEKCHRyYW5zZmVyEk8KCXJlY2lwaWVudBJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhgBEjgKBnNlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYARIZCgZhbW91bnQSDTUwMDAwMDB1c3RhcnMYAWpDCgdtZXNzYWdlEjgKBnNlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYAWobCgdtZXNzYWdlEhAKBm1vZHVsZRIEYmFuaxgBEuIZCKa0gwMSQEFGRjI4NkE4RTBCQjMxNTIzRTA1ODJBNzZCRjA5MTg3RENCRjk2QUVEOENGRThBMDVCNzVBRjZFRjA4OTBGQ0YqQDBBMUUwQTFDMkY2MzZGNzM2RDZGNzMyRTYyNjE2RTZCMkU3NjMxNjI2NTc0NjEzMTJFNEQ3MzY3NTM2NTZFNjQyjwZbeyJldmVudHMiOlt7InR5cGUiOiJjb2luX3JlY2VpdmVkIiwiYXR0cmlidXRlcyI6W3sia2V5IjoicmVjZWl2ZXIiLCJ2YWx1ZSI6InN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24ifSx7ImtleSI6ImFtb3VudCIsInZhbHVlIjoiNTAwMDAwMDAwdXN0YXJzIn1dfSx7InR5cGUiOiJjb2luX3NwZW50IiwiYXR0cmlidXRlcyI6W3sia2V5Ijoic3BlbmRlciIsInZhbHVlIjoic3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgifSx7ImtleSI6ImFtb3VudCIsInZhbHVlIjoiNTAwMDAwMDAwdXN0YXJzIn1dfSx7InR5cGUiOiJtZXNzYWdlIiwiYXR0cmlidXRlcyI6W3sia2V5IjoiYWN0aW9uIiwidmFsdWUiOiIvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kIn0seyJrZXkiOiJzZW5kZXIiLCJ2YWx1ZSI6InN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4In0seyJrZXkiOiJtb2R1bGUiLCJ2YWx1ZSI6ImJhbmsifV19LHsidHlwZSI6InRyYW5zZmVyIiwiYXR0cmlidXRlcyI6W3sia2V5IjoicmVjaXBpZW50IiwidmFsdWUiOiJzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuIn0seyJrZXkiOiJzZW5kZXIiLCJ2YWx1ZSI6InN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4In0seyJrZXkiOiJhbW91bnQiLCJ2YWx1ZSI6IjUwMDAwMDAwMHVzdGFycyJ9XX1dfV06hgQaeAoNY29pbl9yZWNlaXZlZBJMCghyZWNlaXZlchJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhIZCgZhbW91bnQSDzUwMDAwMDAwMHVzdGFycxpgCgpjb2luX3NwZW50EjcKB3NwZW5kZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4EhkKBmFtb3VudBIPNTAwMDAwMDAwdXN0YXJzGnkKB21lc3NhZ2USJgoGYWN0aW9uEhwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEjYKBnNlbmRlchIsc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgSDgoGbW9kdWxlEgRiYW5rGqwBCgh0cmFuc2ZlchJNCglyZWNpcGllbnQSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24SNgoGc2VuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBIZCgZhbW91bnQSDzUwMDAwMDAwMHVzdGFyc0jAmgxQv8AEWroEChUvY29zbW9zLnR4LnYxYmV0YTEuVHgSoAQK8AIKpgEKHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQShQEKLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4EkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGhMKBnVzdGFycxIJNTAwMDAwMDAwEsQBSEFaUHppYUcwWEZGYXZXWWpPWlg2MHJINkxQQUhBRU96a1I2UmpjU05rajk2cThFNFh4MFR2aW5IQU1VdHN0OTgzQ05TMmpybk5mM3NwWU1jb3dZSEJnM3ZwckNKR1hKSTZ1YW5kL2Vja2xVSW1mYkhCNjVBeURiR0xSSnlFb2xNaXB3NjRSU0hhL2pIQjJ1STBRVUxNV3h5SEtQN2g2Skp1bVpzdkRmSEIxaGxoR2MyVUlRZmQ2aTlYZmt4WlZrVm9pThJpClEKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiECIKszyzzlrimoxi4OygKjfL4J70aNA7PbGYBYKXZ4+9sSBAoCCH8YnAESFAoOCgZ1c3RhcnMSBDIwMDAQwJoMGkDZ6/YP0TG2ZHZXrPkdlClUsimgSkBGch/zf1XVN9grhmuMGOxTrJYCPgM4zjgaOT66ri5RR0KlneAlD3eWAxSpYhQyMDIzLTAxLTA3VDExOjMxOjEzWmpfCgpjb2luX3NwZW50EjkKB3NwZW5kZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4GAESFgoGYW1vdW50EgoyMDAwdXN0YXJzGAFqYwoNY29pbl9yZWNlaXZlZBI6CghyZWNlaXZlchIsc3RhcnMxN3hwZnZha20yYW1nOTYyeWxzNmY4NHoza2VsbDhjNWx5OTVhcXYYARIWCgZhbW91bnQSCjIwMDB1c3RhcnMYAWqZAQoIdHJhbnNmZXISOwoJcmVjaXBpZW50EixzdGFyczE3eHBmdmFrbTJhbWc5NjJ5bHM2Zjg0ejNrZWxsOGM1bHk5NWFxdhgBEjgKBnNlbmRlchIsc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgYARIWCgZhbW91bnQSCjIwMDB1c3RhcnMYAWpDCgdtZXNzYWdlEjgKBnNlbmRlchIsc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgYAWpWCgJ0eBITCgNmZWUSCjIwMDB1c3RhcnMYARI7CglmZWVfcGF5ZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4GAFqQwoCdHgSPQoHYWNjX3NlcRIwc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgvMTU2GAFqbQoCdHgSZwoJc2lnbmF0dXJlElgyZXYyRDlFeHRtUjJWNno1SFpRcFZMSXBvRXBBUm5JZjgzOVYxVGZZSzRacmpCanNVNnlXQWo0RE9NNDRHamsrdXE0dVVVZENwWjNnSlE5M2xnTVVxUT09GAFqMwoHbWVzc2FnZRIoCgZhY3Rpb24SHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQYAWpkCgpjb2luX3NwZW50EjkKB3NwZW5kZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4GAESGwoGYW1vdW50Eg81MDAwMDAwMDB1c3RhcnMYAWp8Cg1jb2luX3JlY2VpdmVkEk4KCHJlY2VpdmVyEkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGAESGwoGYW1vdW50Eg81MDAwMDAwMDB1c3RhcnMYAWqyAQoIdHJhbnNmZXISTwoJcmVjaXBpZW50EkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGAESOAoGc2VuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBgBEhsKBmFtb3VudBIPNTAwMDAwMDAwdXN0YXJzGAFqQwoHbWVzc2FnZRI4CgZzZW5kZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4GAFqGwoHbWVzc2FnZRIQCgZtb2R1bGUSBGJhbmsYARKqGAjBtIMDEkAzRDYwNjMxQkFEMkE4MkQ0REI4RUY0Q0FFOTFCNjVDMTkxRjk1RjQ0REMyNUNCNkFGNUU3OEVBODg4MDE1MTdGKkAwQTFFMEExQzJGNjM2RjczNkQ2RjczMkU2MjYxNkU2QjJFNzYzMTYyNjU3NDYxMzEyRTRENzM2NzUzNjU2RTY0MowGW3siZXZlbnRzIjpbeyJ0eXBlIjoiY29pbl9yZWNlaXZlZCIsImF0dHJpYnV0ZXMiOlt7ImtleSI6InJlY2VpdmVyIiwidmFsdWUiOiJzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuIn0seyJrZXkiOiJhbW91bnQiLCJ2YWx1ZSI6IjIwMDAwMDAwdXN0YXJzIn1dfSx7InR5cGUiOiJjb2luX3NwZW50IiwiYXR0cmlidXRlcyI6W3sia2V5Ijoic3BlbmRlciIsInZhbHVlIjoic3RhcnMxZjZnOWd1eWV5emd6amM5bDh3ZzR4bDV4MHJ2eGRkZXdkcWp2N3YifSx7ImtleSI6ImFtb3VudCIsInZhbHVlIjoiMjAwMDAwMDB1c3RhcnMifV19LHsidHlwZSI6Im1lc3NhZ2UiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJhY3Rpb24iLCJ2YWx1ZSI6Ii9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQifSx7ImtleSI6InNlbmRlciIsInZhbHVlIjoic3RhcnMxZjZnOWd1eWV5emd6amM5bDh3ZzR4bDV4MHJ2eGRkZXdkcWp2N3YifSx7ImtleSI6Im1vZHVsZSIsInZhbHVlIjoiYmFuayJ9XX0seyJ0eXBlIjoidHJhbnNmZXIiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJyZWNpcGllbnQiLCJ2YWx1ZSI6InN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24ifSx7ImtleSI6InNlbmRlciIsInZhbHVlIjoic3RhcnMxZjZnOWd1eWV5emd6amM5bDh3ZzR4bDV4MHJ2eGRkZXdkcWp2N3YifSx7ImtleSI6ImFtb3VudCIsInZhbHVlIjoiMjAwMDAwMDB1c3RhcnMifV19XX1dOoMEGncKDWNvaW5fcmVjZWl2ZWQSTAoIcmVjZWl2ZXISQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24SGAoGYW1vdW50Eg4yMDAwMDAwMHVzdGFycxpfCgpjb2luX3NwZW50EjcKB3NwZW5kZXISLHN0YXJzMWY2ZzlndXlleXpnempjOWw4d2c0eGw1eDBydnhkZGV3ZHFqdjd2EhgKBmFtb3VudBIOMjAwMDAwMDB1c3RhcnMaeQoHbWVzc2FnZRImCgZhY3Rpb24SHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQSNgoGc2VuZGVyEixzdGFyczFmNmc5Z3V5ZXl6Z3pqYzlsOHdnNHhsNXgwcnZ4ZGRld2RxanY3dhIOCgZtb2R1bGUSBGJhbmsaqwEKCHRyYW5zZmVyEk0KCXJlY2lwaWVudBJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhI2CgZzZW5kZXISLHN0YXJzMWY2ZzlndXlleXpnempjOWw4d2c0eGw1eDBydnhkZGV3ZHFqdjd2EhgKBmFtb3VudBIOMjAwMDAwMDB1c3RhcnNIwJoMUKvlBFqNAwoVL2Nvc21vcy50eC52MWJldGExLlR4EvMCCsYBCqUBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEoQBCixzdGFyczFmNmc5Z3V5ZXl6Z3pqYzlsOHdnNHhsNXgwcnZ4ZGRld2RxanY3dhJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhoSCgZ1c3RhcnMSCDIwMDAwMDAwEhx5R0hJR2hxV1RrMVhaQ29LZWNJN1I4VFYrQzQyEmYKTgpGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQOXhV8BudrTtyHQIe23pbeHpsSvAdCxMBIHBRKGnxM57BIECgIIfxIUCg4KBnVzdGFycxIENTAwMBDAmgwaQFkQYrqBtRc7p3WllFy+Leb2VvfdxPD4Cgy0Oo/xadaMVJe/4p0roRYRCPXJgLuj+sl7n51PAIBu6KNyJ7A1YCliFDIwMjMtMDEtMDdUMTE6MzM6NTFaal8KCmNvaW5fc3BlbnQSOQoHc3BlbmRlchIsc3RhcnMxZjZnOWd1eWV5emd6amM5bDh3ZzR4bDV4MHJ2eGRkZXdkcWp2N3YYARIWCgZhbW91bnQSCjUwMDB1c3RhcnMYAWpjCg1jb2luX3JlY2VpdmVkEjoKCHJlY2VpdmVyEixzdGFyczE3eHBmdmFrbTJhbWc5NjJ5bHM2Zjg0ejNrZWxsOGM1bHk5NWFxdhgBEhYKBmFtb3VudBIKNTAwMHVzdGFycxgBapkBCgh0cmFuc2ZlchI7CglyZWNpcGllbnQSLHN0YXJzMTd4cGZ2YWttMmFtZzk2MnlsczZmODR6M2tlbGw4YzVseTk1YXF2GAESOAoGc2VuZGVyEixzdGFyczFmNmc5Z3V5ZXl6Z3pqYzlsOHdnNHhsNXgwcnZ4ZGRld2RxanY3dhgBEhYKBmFtb3VudBIKNTAwMHVzdGFycxgBakMKB21lc3NhZ2USOAoGc2VuZGVyEixzdGFyczFmNmc5Z3V5ZXl6Z3pqYzlsOHdnNHhsNXgwcnZ4ZGRld2RxanY3dhgBalYKAnR4EhMKA2ZlZRIKNTAwMHVzdGFycxgBEjsKCWZlZV9wYXllchIsc3RhcnMxZjZnOWd1eWV5emd6amM5bDh3ZzR4bDV4MHJ2eGRkZXdkcWp2N3YYAWpBCgJ0eBI7CgdhY2Nfc2VxEi5zdGFyczFmNmc5Z3V5ZXl6Z3pqYzlsOHdnNHhsNXgwcnZ4ZGRld2RxanY3di8wGAFqbQoCdHgSZwoJc2lnbmF0dXJlElhXUkJpdW9HMUZ6dW5kYVdVWEw0dDV2Wlc5OTNFOFBnS0RMUTZqL0ZwMW94VWw3L2luU3VoRmhFSTljbUF1NlA2eVh1Zm5VOEFnRzdvbzNJbnNEVmdLUT09GAFqMwoHbWVzc2FnZRIoCgZhY3Rpb24SHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQYAWpjCgpjb2luX3NwZW50EjkKB3NwZW5kZXISLHN0YXJzMWY2ZzlndXlleXpnempjOWw4d2c0eGw1eDBydnhkZGV3ZHFqdjd2GAESGgoGYW1vdW50Eg4yMDAwMDAwMHVzdGFycxgBansKDWNvaW5fcmVjZWl2ZWQSTgoIcmVjZWl2ZXISQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24YARIaCgZhbW91bnQSDjIwMDAwMDAwdXN0YXJzGAFqsQEKCHRyYW5zZmVyEk8KCXJlY2lwaWVudBJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhgBEjgKBnNlbmRlchIsc3RhcnMxZjZnOWd1eWV5emd6amM5bDh3ZzR4bDV4MHJ2eGRkZXdkcWp2N3YYARIaCgZhbW91bnQSDjIwMDAwMDAwdXN0YXJzGAFqQwoHbWVzc2FnZRI4CgZzZW5kZXISLHN0YXJzMWY2ZzlndXlleXpnempjOWw4d2c0eGw1eDBydnhkZGV3ZHFqdjd2GAFqGwoHbWVzc2FnZRIQCgZtb2R1bGUSBGJhbmsYARK0GQirtYMDEkBEMjU2NDE5NUI2QjIzOTQxRkI4OTI3NjQzQ0QxQ0IyMzU1OEUzRDAzMURDRjc0RDcxNzZCNTQ1NzlCMDBDQjRGKkAwQTFFMEExQzJGNjM2RjczNkQ2RjczMkU2MjYxNkU2QjJFNzYzMTYyNjU3NDYxMzEyRTRENzM2NzUzNjU2RTY0MpIGW3siZXZlbnRzIjpbeyJ0eXBlIjoiY29pbl9yZWNlaXZlZCIsImF0dHJpYnV0ZXMiOlt7ImtleSI6InJlY2VpdmVyIiwidmFsdWUiOiJzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuIn0seyJrZXkiOiJhbW91bnQiLCJ2YWx1ZSI6IjE1MDAwMDAwMDB1c3RhcnMifV19LHsidHlwZSI6ImNvaW5fc3BlbnQiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJzcGVuZGVyIiwidmFsdWUiOiJzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOCJ9LHsia2V5IjoiYW1vdW50IiwidmFsdWUiOiIxNTAwMDAwMDAwdXN0YXJzIn1dfSx7InR5cGUiOiJtZXNzYWdlIiwiYXR0cmlidXRlcyI6W3sia2V5IjoiYWN0aW9uIiwidmFsdWUiOiIvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kIn0seyJrZXkiOiJzZW5kZXIiLCJ2YWx1ZSI6InN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4In0seyJrZXkiOiJtb2R1bGUiLCJ2YWx1ZSI6ImJhbmsifV19LHsidHlwZSI6InRyYW5zZmVyIiwiYXR0cmlidXRlcyI6W3sia2V5IjoicmVjaXBpZW50IiwidmFsdWUiOiJzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuIn0seyJrZXkiOiJzZW5kZXIiLCJ2YWx1ZSI6InN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4In0seyJrZXkiOiJhbW91bnQiLCJ2YWx1ZSI6IjE1MDAwMDAwMDB1c3RhcnMifV19XX1dOokEGnkKDWNvaW5fcmVjZWl2ZWQSTAoIcmVjZWl2ZXISQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24SGgoGYW1vdW50EhAxNTAwMDAwMDAwdXN0YXJzGmEKCmNvaW5fc3BlbnQSNwoHc3BlbmRlchIsc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgSGgoGYW1vdW50EhAxNTAwMDAwMDAwdXN0YXJzGnkKB21lc3NhZ2USJgoGYWN0aW9uEhwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEjYKBnNlbmRlchIsc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgSDgoGbW9kdWxlEgRiYW5rGq0BCgh0cmFuc2ZlchJNCglyZWNpcGllbnQSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24SNgoGc2VuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBIaCgZhbW91bnQSEDE1MDAwMDAwMDB1c3RhcnNIwJoMUIS8BFqDBAoVL2Nvc21vcy50eC52MWJldGExLlR4EukDCrkCCqcBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEoYBCixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhoUCgZ1c3RhcnMSCjE1MDAwMDAwMDASjAFLQVpQemlhRzBYRkZhdldZak9aWDYwckg2TFBBS0F5M05BM1loMWNPRUlmY1FENVB2SjY5QVpUZUtDTDNua0NaS1VUUmMzWXlycythZ1lwa1Z0SEhLQ3BlNHNHREsraWI1VHlCRFZFZUNUeGxySzNvS0RENGYybm5ZME1GUkJXSGl2ZnIrSWJldEpzZhJpClEKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiECIKszyzzlrimoxi4OygKjfL4J70aNA7PbGYBYKXZ4+9sSBAoCCH8YnQESFAoOCgZ1c3RhcnMSBDIwMDAQwJoMGkAGPkcAPFYfAiXW+VJ5a9/e6pMpV/PA38uGuqnelI6ZMn2u4WT79slE1abCINSdeWWrM56xiZ8qVBQLmpr/MXhdYhQyMDIzLTAxLTA3VDExOjQ0OjE2WmpfCgpjb2luX3NwZW50EjkKB3NwZW5kZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4GAESFgoGYW1vdW50EgoyMDAwdXN0YXJzGAFqYwoNY29pbl9yZWNlaXZlZBI6CghyZWNlaXZlchIsc3RhcnMxN3hwZnZha20yYW1nOTYyeWxzNmY4NHoza2VsbDhjNWx5OTVhcXYYARIWCgZhbW91bnQSCjIwMDB1c3RhcnMYAWqZAQoIdHJhbnNmZXISOwoJcmVjaXBpZW50EixzdGFyczE3eHBmdmFrbTJhbWc5NjJ5bHM2Zjg0ejNrZWxsOGM1bHk5NWFxdhgBEjgKBnNlbmRlchIsc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgYARIWCgZhbW91bnQSCjIwMDB1c3RhcnMYAWpDCgdtZXNzYWdlEjgKBnNlbmRlchIsc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgYAWpWCgJ0eBITCgNmZWUSCjIwMDB1c3RhcnMYARI7CglmZWVfcGF5ZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4GAFqQwoCdHgSPQoHYWNjX3NlcRIwc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgvMTU3GAFqbQoCdHgSZwoJc2lnbmF0dXJlElhCajVIQUR4V0h3SWwxdmxTZVd2ZjN1cVRLVmZ6d04vTGhycXAzcFNPbVRKOXJ1RmsrL2JKUk5XbXdpRFVuWGxscXpPZXNZbWZLbFFVQzVxYS96RjRYUT09GAFqMwoHbWVzc2FnZRIoCgZhY3Rpb24SHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQYAWplCgpjb2luX3NwZW50EjkKB3NwZW5kZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4GAESHAoGYW1vdW50EhAxNTAwMDAwMDAwdXN0YXJzGAFqfQoNY29pbl9yZWNlaXZlZBJOCghyZWNlaXZlchJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhgBEhwKBmFtb3VudBIQMTUwMDAwMDAwMHVzdGFycxgBarMBCgh0cmFuc2ZlchJPCglyZWNpcGllbnQSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24YARI4CgZzZW5kZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4GAESHAoGYW1vdW50EhAxNTAwMDAwMDAwdXN0YXJzGAFqQwoHbWVzc2FnZRI4CgZzZW5kZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4GAFqGwoHbWVzc2FnZRIQCgZtb2R1bGUSBGJhbmsYARLLGQixtYMDEkA1RDkzRDUyRUZDNkM0MDJEMTMwRTlDMzRCM0QyOUNCNkFDMkM0NkZCQjhBNTNCNEZGRTVCRTE5QUE5M0ZFMUVGKkAwQTFFMEExQzJGNjM2RjczNkQ2RjczMkU2MjYxNkU2QjJFNzYzMTYyNjU3NDYxMzEyRTRENzM2NzUzNjU2RTY0MokGW3siZXZlbnRzIjpbeyJ0eXBlIjoiY29pbl9yZWNlaXZlZCIsImF0dHJpYnV0ZXMiOlt7ImtleSI6InJlY2VpdmVyIiwidmFsdWUiOiJzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuIn0seyJrZXkiOiJhbW91bnQiLCJ2YWx1ZSI6IjEwMDAwMDB1c3RhcnMifV19LHsidHlwZSI6ImNvaW5fc3BlbnQiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJzcGVuZGVyIiwidmFsdWUiOiJzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5ayJ9LHsia2V5IjoiYW1vdW50IiwidmFsdWUiOiIxMDAwMDAwdXN0YXJzIn1dfSx7InR5cGUiOiJtZXNzYWdlIiwiYXR0cmlidXRlcyI6W3sia2V5IjoiYWN0aW9uIiwidmFsdWUiOiIvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kIn0seyJrZXkiOiJzZW5kZXIiLCJ2YWx1ZSI6InN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrIn0seyJrZXkiOiJtb2R1bGUiLCJ2YWx1ZSI6ImJhbmsifV19LHsidHlwZSI6InRyYW5zZmVyIiwiYXR0cmlidXRlcyI6W3sia2V5IjoicmVjaXBpZW50IiwidmFsdWUiOiJzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuIn0seyJrZXkiOiJzZW5kZXIiLCJ2YWx1ZSI6InN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrIn0seyJrZXkiOiJhbW91bnQiLCJ2YWx1ZSI6IjEwMDAwMDB1c3RhcnMifV19XX1dOoAEGnYKDWNvaW5fcmVjZWl2ZWQSTAoIcmVjZWl2ZXISQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24SFwoGYW1vdW50Eg0xMDAwMDAwdXN0YXJzGl4KCmNvaW5fc3BlbnQSNwoHc3BlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsSFwoGYW1vdW50Eg0xMDAwMDAwdXN0YXJzGnkKB21lc3NhZ2USJgoGYWN0aW9uEhwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEjYKBnNlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsSDgoGbW9kdWxlEgRiYW5rGqoBCgh0cmFuc2ZlchJNCglyZWNpcGllbnQSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24SNgoGc2VuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axIXCgZhbW91bnQSDTEwMDAwMDB1c3RhcnNIwJoMUM++BFq3BAoVL2Nvc21vcy50eC52MWJldGExLlR4Ep0ECu4CCqQBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEoMBCixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhoRCgZ1c3RhcnMSBzEwMDAwMDASxAFIQUNpS1Q1UnBYUmlPVGs0SFlabWlmenlWeFRnSEF4UjJuQ0MrVy9OdmxRREhwb3drQ1pZeDJEakhBazVrY29qVnVMTStCUlIvRDQ0QTZBTi9zSllIQUcxWkEvOTBWMnhGZXVjZkZEVzNlR29FYmJjSEF5M05BM1loMWNPRUlmY1FENVB2SjY5QVpUZUhCMWhsaEdjMlVJUWZkNmk5WGZreFpWa1ZvaU5IQjY1QXlEYkdMUkp5RW9sTWlwdzY0UlNIYS9qEmgKUApGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQJL+ST18/7C3MU49TUq3zXARBbSwih3u+jgDtMlh+ymQhIECgIIfxgBEhQKDgoGdXN0YXJzEgQ1MDAwEMCaDBpAgKoGTvvPuCkhiE+cl05IHS8826MatiXvLO7lqoY4gghDLH3GSQbbIThHkyJcxI4gYbmW3zdxfOR2ojI1kChyZmIUMjAyMy0wMS0wN1QxMTo0NDo1MlpqXwoKY29pbl9zcGVudBI5CgdzcGVuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axgBEhYKBmFtb3VudBIKNTAwMHVzdGFycxgBamMKDWNvaW5fcmVjZWl2ZWQSOgoIcmVjZWl2ZXISLHN0YXJzMTd4cGZ2YWttMmFtZzk2MnlsczZmODR6M2tlbGw4YzVseTk1YXF2GAESFgoGYW1vdW50Ego1MDAwdXN0YXJzGAFqmQEKCHRyYW5zZmVyEjsKCXJlY2lwaWVudBIsc3RhcnMxN3hwZnZha20yYW1nOTYyeWxzNmY4NHoza2VsbDhjNWx5OTVhcXYYARI4CgZzZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrGAESFgoGYW1vdW50Ego1MDAwdXN0YXJzGAFqQwoHbWVzc2FnZRI4CgZzZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrGAFqVgoCdHgSEwoDZmVlEgo1MDAwdXN0YXJzGAESOwoJZmVlX3BheWVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axgBakEKAnR4EjsKB2FjY19zZXESLnN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrLzEYAWptCgJ0eBJnCglzaWduYXR1cmUSWGdLb0dUdnZQdUNraGlFK2NsMDVJSFM4ODI2TWF0aVh2TE83bHFvWTRnZ2hETEgzR1NRYmJJVGhIa3lKY3hJNGdZYm1XM3pkeGZPUjJvakkxa0NoeVpnPT0YAWozCgdtZXNzYWdlEigKBmFjdGlvbhIcL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBgBamIKCmNvaW5fc3BlbnQSOQoHc3BlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYARIZCgZhbW91bnQSDTEwMDAwMDB1c3RhcnMYAWp6Cg1jb2luX3JlY2VpdmVkEk4KCHJlY2VpdmVyEkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGAESGQoGYW1vdW50Eg0xMDAwMDAwdXN0YXJzGAFqsAEKCHRyYW5zZmVyEk8KCXJlY2lwaWVudBJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhgBEjgKBnNlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYARIZCgZhbW91bnQSDTEwMDAwMDB1c3RhcnMYAWpDCgdtZXNzYWdlEjgKBnNlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYAWobCgdtZXNzYWdlEhAKBm1vZHVsZRIEYmFuaxgBEtoYCMK2gwMSQEU3REMwMkM5RkIyQzVEN0Q5M0VEOUY0OEIwNTU0MEYxRDFCRTVCMjA5OTk1QkQyM0FDOTQ4RTE0RjU5ODc2NTMqQDBBMUUwQTFDMkY2MzZGNzM2RDZGNzMyRTYyNjE2RTZCMkU3NjMxNjI2NTc0NjEzMTJFNEQ3MzY3NTM2NTZFNjQyiQZbeyJldmVudHMiOlt7InR5cGUiOiJjb2luX3JlY2VpdmVkIiwiYXR0cmlidXRlcyI6W3sia2V5IjoicmVjZWl2ZXIiLCJ2YWx1ZSI6InN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24ifSx7ImtleSI6ImFtb3VudCIsInZhbHVlIjoiMjAwMDAwMHVzdGFycyJ9XX0seyJ0eXBlIjoiY29pbl9zcGVudCIsImF0dHJpYnV0ZXMiOlt7ImtleSI6InNwZW5kZXIiLCJ2YWx1ZSI6InN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrIn0seyJrZXkiOiJhbW91bnQiLCJ2YWx1ZSI6IjIwMDAwMDB1c3RhcnMifV19LHsidHlwZSI6Im1lc3NhZ2UiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJhY3Rpb24iLCJ2YWx1ZSI6Ii9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQifSx7ImtleSI6InNlbmRlciIsInZhbHVlIjoic3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsifSx7ImtleSI6Im1vZHVsZSIsInZhbHVlIjoiYmFuayJ9XX0seyJ0eXBlIjoidHJhbnNmZXIiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJyZWNpcGllbnQiLCJ2YWx1ZSI6InN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24ifSx7ImtleSI6InNlbmRlciIsInZhbHVlIjoic3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsifSx7ImtleSI6ImFtb3VudCIsInZhbHVlIjoiMjAwMDAwMHVzdGFycyJ9XX1dfV06gAQadgoNY29pbl9yZWNlaXZlZBJMCghyZWNlaXZlchJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhIXCgZhbW91bnQSDTIwMDAwMDB1c3RhcnMaXgoKY29pbl9zcGVudBI3CgdzcGVuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axIXCgZhbW91bnQSDTIwMDAwMDB1c3RhcnMaeQoHbWVzc2FnZRImCgZhY3Rpb24SHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQSNgoGc2VuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axIOCgZtb2R1bGUSBGJhbmsaqgEKCHRyYW5zZmVyEk0KCXJlY2lwaWVudBJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhI2CgZzZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrEhcKBmFtb3VudBINMjAwMDAwMHVzdGFyc0jAmgxQz7QEWsYDChUvY29zbW9zLnR4LnYxYmV0YTEuVHgSrAMK/QEKpAEKHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQSgwEKLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrEkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGhEKBnVzdGFycxIHMjAwMDAwMBJUUWdFT3prUjZSamNTTmtqOTZxOEU0WHgwVHZpblFnRzFaQS85MFYyeEZldWNmRkRXM2VHb0ViYmNRZ3NNYi9RSzJiS1NtbXBrV0JORWsyMVB5dldOEmgKUApGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQJL+ST18/7C3MU49TUq3zXARBbSwih3u+jgDtMlh+ymQhIECgIIfxgCEhQKDgoGdXN0YXJzEgQ1MDAwEMCaDBpA5L6shFLPN7P0F/vBaL1DN8vnZ3atjeWSqAQvdWnnJXhhQSILb6H+iEV6C4lYqVhlEBki9UDJo/CqUwVIwb2932IUMjAyMy0wMS0wN1QxMTo1OTowNVpqXwoKY29pbl9zcGVudBI5CgdzcGVuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axgBEhYKBmFtb3VudBIKNTAwMHVzdGFycxgBamMKDWNvaW5fcmVjZWl2ZWQSOgoIcmVjZWl2ZXISLHN0YXJzMTd4cGZ2YWttMmFtZzk2MnlsczZmODR6M2tlbGw4YzVseTk1YXF2GAESFgoGYW1vdW50Ego1MDAwdXN0YXJzGAFqmQEKCHRyYW5zZmVyEjsKCXJlY2lwaWVudBIsc3RhcnMxN3hwZnZha20yYW1nOTYyeWxzNmY4NHoza2VsbDhjNWx5OTVhcXYYARI4CgZzZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrGAESFgoGYW1vdW50Ego1MDAwdXN0YXJzGAFqQwoHbWVzc2FnZRI4CgZzZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrGAFqVgoCdHgSEwoDZmVlEgo1MDAwdXN0YXJzGAESOwoJZmVlX3BheWVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axgBakEKAnR4EjsKB2FjY19zZXESLnN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrLzIYAWptCgJ0eBJnCglzaWduYXR1cmUSWDVMNnNoRkxQTjdQMEYvdkJhTDFETjh2blozYXRqZVdTcUFRdmRXbm5KWGhoUVNJTGI2SCtpRVY2QzRsWXFWaGxFQmtpOVVESm8vQ3FVd1ZJd2IyOTN3PT0YAWozCgdtZXNzYWdlEigKBmFjdGlvbhIcL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBgBamIKCmNvaW5fc3BlbnQSOQoHc3BlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYARIZCgZhbW91bnQSDTIwMDAwMDB1c3RhcnMYAWp6Cg1jb2luX3JlY2VpdmVkEk4KCHJlY2VpdmVyEkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGAESGQoGYW1vdW50Eg0yMDAwMDAwdXN0YXJzGAFqsAEKCHRyYW5zZmVyEk8KCXJlY2lwaWVudBJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhgBEjgKBnNlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYARIZCgZhbW91bnQSDTIwMDAwMDB1c3RhcnMYAWpDCgdtZXNzYWdlEjgKBnNlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYAWobCgdtZXNzYWdlEhAKBm1vZHVsZRIEYmFuaxgBEtoYCMq2gwMSQENEQTA2NjRBRTk4MzMxRjMxMjg5Njc4MEYzRjdFNkFDNjZGMjZDODNGN0VCMzMwQ0MwNDA4RDE5MjI1MDA5RkIqQDBBMUUwQTFDMkY2MzZGNzM2RDZGNzMyRTYyNjE2RTZCMkU3NjMxNjI2NTc0NjEzMTJFNEQ3MzY3NTM2NTZFNjQyiQZbeyJldmVudHMiOlt7InR5cGUiOiJjb2luX3JlY2VpdmVkIiwiYXR0cmlidXRlcyI6W3sia2V5IjoicmVjZWl2ZXIiLCJ2YWx1ZSI6InN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24ifSx7ImtleSI6ImFtb3VudCIsInZhbHVlIjoiMTAwMDAwMHVzdGFycyJ9XX0seyJ0eXBlIjoiY29pbl9zcGVudCIsImF0dHJpYnV0ZXMiOlt7ImtleSI6InNwZW5kZXIiLCJ2YWx1ZSI6InN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrIn0seyJrZXkiOiJhbW91bnQiLCJ2YWx1ZSI6IjEwMDAwMDB1c3RhcnMifV19LHsidHlwZSI6Im1lc3NhZ2UiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJhY3Rpb24iLCJ2YWx1ZSI6Ii9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQifSx7ImtleSI6InNlbmRlciIsInZhbHVlIjoic3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsifSx7ImtleSI6Im1vZHVsZSIsInZhbHVlIjoiYmFuayJ9XX0seyJ0eXBlIjoidHJhbnNmZXIiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJyZWNpcGllbnQiLCJ2YWx1ZSI6InN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24ifSx7ImtleSI6InNlbmRlciIsInZhbHVlIjoic3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsifSx7ImtleSI6ImFtb3VudCIsInZhbHVlIjoiMTAwMDAwMHVzdGFycyJ9XX1dfV06gAQadgoNY29pbl9yZWNlaXZlZBJMCghyZWNlaXZlchJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhIXCgZhbW91bnQSDTEwMDAwMDB1c3RhcnMaXgoKY29pbl9zcGVudBI3CgdzcGVuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axIXCgZhbW91bnQSDTEwMDAwMDB1c3RhcnMaeQoHbWVzc2FnZRImCgZhY3Rpb24SHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQSNgoGc2VuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axIOCgZtb2R1bGUSBGJhbmsaqgEKCHRyYW5zZmVyEk0KCXJlY2lwaWVudBJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhI2CgZzZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrEhcKBmFtb3VudBINMTAwMDAwMHVzdGFyc0jAmgxQgrUEWsYDChUvY29zbW9zLnR4LnYxYmV0YTEuVHgSrAMK/QEKpAEKHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQSgwEKLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrEkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGhEKBnVzdGFycxIHMTAwMDAwMBJUUWdFT3prUjZSamNTTmtqOTZxOEU0WHgwVHZpblFnTVV0c3Q5ODNDTlMyanJuTmYzc3BZTWNvd1lRZ3hSMm5DQytXL052bFFESHBvd2tDWll4MkRqEmgKUApGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQJL+ST18/7C3MU49TUq3zXARBbSwih3u+jgDtMlh+ymQhIECgIIfxgDEhQKDgoGdXN0YXJzEgQ1MDAwEMCaDBpArBmTgBwvFvuOjeVWhQyOkszDJYhzZUu65+3ofh88VZACpmMg5JSGNA7weGEJIfY8sQl7f+fDLs0C9xeGpV58pWIUMjAyMy0wMS0wN1QxMTo1OTo1MVpqXwoKY29pbl9zcGVudBI5CgdzcGVuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axgBEhYKBmFtb3VudBIKNTAwMHVzdGFycxgBamMKDWNvaW5fcmVjZWl2ZWQSOgoIcmVjZWl2ZXISLHN0YXJzMTd4cGZ2YWttMmFtZzk2MnlsczZmODR6M2tlbGw4YzVseTk1YXF2GAESFgoGYW1vdW50Ego1MDAwdXN0YXJzGAFqmQEKCHRyYW5zZmVyEjsKCXJlY2lwaWVudBIsc3RhcnMxN3hwZnZha20yYW1nOTYyeWxzNmY4NHoza2VsbDhjNWx5OTVhcXYYARI4CgZzZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrGAESFgoGYW1vdW50Ego1MDAwdXN0YXJzGAFqQwoHbWVzc2FnZRI4CgZzZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrGAFqVgoCdHgSEwoDZmVlEgo1MDAwdXN0YXJzGAESOwoJZmVlX3BheWVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axgBakEKAnR4EjsKB2FjY19zZXESLnN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrLzMYAWptCgJ0eBJnCglzaWduYXR1cmUSWHJCbVRnQnd2RnZ1T2plVldoUXlPa3N6REpZaHpaVXU2NSszb2ZoODhWWkFDcG1NZzVKU0dOQTd3ZUdFSklmWThzUWw3ZitmRExzMEM5eGVHcFY1OHBRPT0YAWozCgdtZXNzYWdlEigKBmFjdGlvbhIcL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBgBamIKCmNvaW5fc3BlbnQSOQoHc3BlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYARIZCgZhbW91bnQSDTEwMDAwMDB1c3RhcnMYAWp6Cg1jb2luX3JlY2VpdmVkEk4KCHJlY2VpdmVyEkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGAESGQoGYW1vdW50Eg0xMDAwMDAwdXN0YXJzGAFqsAEKCHRyYW5zZmVyEk8KCXJlY2lwaWVudBJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhgBEjgKBnNlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYARIZCgZhbW91bnQSDTEwMDAwMDB1c3RhcnMYAWpDCgdtZXNzYWdlEjgKBnNlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYAWobCgdtZXNzYWdlEhAKBm1vZHVsZRIEYmFuaxgBEqIYCPLGgwMSQDcwOUMzRTMwNTUyQzY2NzJDNkI1Mzk2NTIzRDgxNTQzODQyMTZDMUYwQkI2MTI0QzJCOUNEOENGNzhCRDY4RkEqQDBBMUUwQTFDMkY2MzZGNzM2RDZGNzMyRTYyNjE2RTZCMkU3NjMxNjI2NTc0NjEzMTJFNEQ3MzY3NTM2NTZFNjQyiQZbeyJldmVudHMiOlt7InR5cGUiOiJjb2luX3JlY2VpdmVkIiwiYXR0cmlidXRlcyI6W3sia2V5IjoicmVjZWl2ZXIiLCJ2YWx1ZSI6InN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24ifSx7ImtleSI6ImFtb3VudCIsInZhbHVlIjoiNzAwMDAwMHVzdGFycyJ9XX0seyJ0eXBlIjoiY29pbl9zcGVudCIsImF0dHJpYnV0ZXMiOlt7ImtleSI6InNwZW5kZXIiLCJ2YWx1ZSI6InN0YXJzMWY2ZzlndXlleXpnempjOWw4d2c0eGw1eDBydnhkZGV3ZHFqdjd2In0seyJrZXkiOiJhbW91bnQiLCJ2YWx1ZSI6IjcwMDAwMDB1c3RhcnMifV19LHsidHlwZSI6Im1lc3NhZ2UiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJhY3Rpb24iLCJ2YWx1ZSI6Ii9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQifSx7ImtleSI6InNlbmRlciIsInZhbHVlIjoic3RhcnMxZjZnOWd1eWV5emd6amM5bDh3ZzR4bDV4MHJ2eGRkZXdkcWp2N3YifSx7ImtleSI6Im1vZHVsZSIsInZhbHVlIjoiYmFuayJ9XX0seyJ0eXBlIjoidHJhbnNmZXIiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJyZWNpcGllbnQiLCJ2YWx1ZSI6InN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24ifSx7ImtleSI6InNlbmRlciIsInZhbHVlIjoic3RhcnMxZjZnOWd1eWV5emd6amM5bDh3ZzR4bDV4MHJ2eGRkZXdkcWp2N3YifSx7ImtleSI6ImFtb3VudCIsInZhbHVlIjoiNzAwMDAwMHVzdGFycyJ9XX1dfV06gAQadgoNY29pbl9yZWNlaXZlZBJMCghyZWNlaXZlchJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhIXCgZhbW91bnQSDTcwMDAwMDB1c3RhcnMaXgoKY29pbl9zcGVudBI3CgdzcGVuZGVyEixzdGFyczFmNmc5Z3V5ZXl6Z3pqYzlsOHdnNHhsNXgwcnZ4ZGRld2RxanY3dhIXCgZhbW91bnQSDTcwMDAwMDB1c3RhcnMaeQoHbWVzc2FnZRImCgZhY3Rpb24SHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQSNgoGc2VuZGVyEixzdGFyczFmNmc5Z3V5ZXl6Z3pqYzlsOHdnNHhsNXgwcnZ4ZGRld2RxanY3dhIOCgZtb2R1bGUSBGJhbmsaqgEKCHRyYW5zZmVyEk0KCXJlY2lwaWVudBJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhI2CgZzZW5kZXISLHN0YXJzMWY2ZzlndXlleXpnempjOWw4d2c0eGw1eDBydnhkZGV3ZHFqdjd2EhcKBmFtb3VudBINNzAwMDAwMHVzdGFyc0jAmgxQybAEWo4DChUvY29zbW9zLnR4LnYxYmV0YTEuVHgS9AIKxQEKpAEKHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQSgwEKLHN0YXJzMWY2ZzlndXlleXpnempjOWw4d2c0eGw1eDBydnhkZGV3ZHFqdjd2EkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGhEKBnVzdGFycxIHNzAwMDAwMBIceUFDaUtUNVJwWFJpT1RrNEhZWm1pZnp5VnhUZxJoClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiEDl4VfAbna07ch0CHtt6W3h6bErwHQsTASBwUShp8TOewSBAoCCH8YARIUCg4KBnVzdGFycxIENTAwMBDAmgwaQLtDxquZKz0SAf3FFz0M4fSm08zRgEOaXRoHlwmAt2rjNT+pYTclthlYTxzSB5Im9poORFxf4isVUDxxPtVWAeBiFDIwMjMtMDEtMDdUMTU6MjQ6NTdaal8KCmNvaW5fc3BlbnQSOQoHc3BlbmRlchIsc3RhcnMxZjZnOWd1eWV5emd6amM5bDh3ZzR4bDV4MHJ2eGRkZXdkcWp2N3YYARIWCgZhbW91bnQSCjUwMDB1c3RhcnMYAWpjCg1jb2luX3JlY2VpdmVkEjoKCHJlY2VpdmVyEixzdGFyczE3eHBmdmFrbTJhbWc5NjJ5bHM2Zjg0ejNrZWxsOGM1bHk5NWFxdhgBEhYKBmFtb3VudBIKNTAwMHVzdGFycxgBapkBCgh0cmFuc2ZlchI7CglyZWNpcGllbnQSLHN0YXJzMTd4cGZ2YWttMmFtZzk2MnlsczZmODR6M2tlbGw4YzVseTk1YXF2GAESOAoGc2VuZGVyEixzdGFyczFmNmc5Z3V5ZXl6Z3pqYzlsOHdnNHhsNXgwcnZ4ZGRld2RxanY3dhgBEhYKBmFtb3VudBIKNTAwMHVzdGFycxgBakMKB21lc3NhZ2USOAoGc2VuZGVyEixzdGFyczFmNmc5Z3V5ZXl6Z3pqYzlsOHdnNHhsNXgwcnZ4ZGRld2RxanY3dhgBalYKAnR4EhMKA2ZlZRIKNTAwMHVzdGFycxgBEjsKCWZlZV9wYXllchIsc3RhcnMxZjZnOWd1eWV5emd6amM5bDh3ZzR4bDV4MHJ2eGRkZXdkcWp2N3YYAWpBCgJ0eBI7CgdhY2Nfc2VxEi5zdGFyczFmNmc5Z3V5ZXl6Z3pqYzlsOHdnNHhsNXgwcnZ4ZGRld2RxanY3di8xGAFqbQoCdHgSZwoJc2lnbmF0dXJlElh1MFBHcTVrclBSSUIvY1VYUFF6aDlLYlR6TkdBUTVwZEdnZVhDWUMzYXVNMVA2bGhOeVcyR1ZoUEhOSUhraWIybWc1RVhGL2lLeFZRUEhFKzFWWUI0QT09GAFqMwoHbWVzc2FnZRIoCgZhY3Rpb24SHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQYAWpiCgpjb2luX3NwZW50EjkKB3NwZW5kZXISLHN0YXJzMWY2ZzlndXlleXpnempjOWw4d2c0eGw1eDBydnhkZGV3ZHFqdjd2GAESGQoGYW1vdW50Eg03MDAwMDAwdXN0YXJzGAFqegoNY29pbl9yZWNlaXZlZBJOCghyZWNlaXZlchJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhgBEhkKBmFtb3VudBINNzAwMDAwMHVzdGFycxgBarABCgh0cmFuc2ZlchJPCglyZWNpcGllbnQSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24YARI4CgZzZW5kZXISLHN0YXJzMWY2ZzlndXlleXpnempjOWw4d2c0eGw1eDBydnhkZGV3ZHFqdjd2GAESGQoGYW1vdW50Eg03MDAwMDAwdXN0YXJzGAFqQwoHbWVzc2FnZRI4CgZzZW5kZXISLHN0YXJzMWY2ZzlndXlleXpnempjOWw4d2c0eGw1eDBydnhkZGV3ZHFqdjd2GAFqGwoHbWVzc2FnZRIQCgZtb2R1bGUSBGJhbmsYARLaGAinx4MDEkAyQUM4MzA4Qjc1MTA1Mjc5N0M3RDc3MDI5ODczQzc5MDhGM0Q0NjA0NTk5NEExMDNDMkI2RUM4NEYxODhDQkNFKkAwQTFFMEExQzJGNjM2RjczNkQ2RjczMkU2MjYxNkU2QjJFNzYzMTYyNjU3NDYxMzEyRTRENzM2NzUzNjU2RTY0MokGW3siZXZlbnRzIjpbeyJ0eXBlIjoiY29pbl9yZWNlaXZlZCIsImF0dHJpYnV0ZXMiOlt7ImtleSI6InJlY2VpdmVyIiwidmFsdWUiOiJzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuIn0seyJrZXkiOiJhbW91bnQiLCJ2YWx1ZSI6IjE1MDAwMDB1c3RhcnMifV19LHsidHlwZSI6ImNvaW5fc3BlbnQiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJzcGVuZGVyIiwidmFsdWUiOiJzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5ayJ9LHsia2V5IjoiYW1vdW50IiwidmFsdWUiOiIxNTAwMDAwdXN0YXJzIn1dfSx7InR5cGUiOiJtZXNzYWdlIiwiYXR0cmlidXRlcyI6W3sia2V5IjoiYWN0aW9uIiwidmFsdWUiOiIvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kIn0seyJrZXkiOiJzZW5kZXIiLCJ2YWx1ZSI6InN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrIn0seyJrZXkiOiJtb2R1bGUiLCJ2YWx1ZSI6ImJhbmsifV19LHsidHlwZSI6InRyYW5zZmVyIiwiYXR0cmlidXRlcyI6W3sia2V5IjoicmVjaXBpZW50IiwidmFsdWUiOiJzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuIn0seyJrZXkiOiJzZW5kZXIiLCJ2YWx1ZSI6InN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrIn0seyJrZXkiOiJhbW91bnQiLCJ2YWx1ZSI6IjE1MDAwMDB1c3RhcnMifV19XX1dOoAEGnYKDWNvaW5fcmVjZWl2ZWQSTAoIcmVjZWl2ZXISQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24SFwoGYW1vdW50Eg0xNTAwMDAwdXN0YXJzGl4KCmNvaW5fc3BlbnQSNwoHc3BlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsSFwoGYW1vdW50Eg0xNTAwMDAwdXN0YXJzGnkKB21lc3NhZ2USJgoGYWN0aW9uEhwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEjYKBnNlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsSDgoGbW9kdWxlEgRiYW5rGqoBCgh0cmFuc2ZlchJNCglyZWNpcGllbnQSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24SNgoGc2VuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axIXCgZhbW91bnQSDTE1MDAwMDB1c3RhcnNIwJoMUM+0BFrGAwoVL2Nvc21vcy50eC52MWJldGExLlR4EqwDCv0BCqQBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEoMBCixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhoRCgZ1c3RhcnMSBzE1MDAwMDASVFFnRU96a1I2UmpjU05rajk2cThFNFh4MFR2aW5RZ01VdHN0OTgzQ05TMmpybk5mM3NwWU1jb3dZUWdzTWIvUUsyYktTbW1wa1dCTkVrMjFQeXZXThJoClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiECS/kk9fP+wtzFOPU1Kt81wEQW0sIod7vo4A7TJYfspkISBAoCCH8YBBIUCg4KBnVzdGFycxIENTAwMBDAmgwaQElsMmmLjt+wfC6XC/hUpAPjh3pVeudXshprsE1JYlSDSa90WSQklBKGlbDb2jq0kLLbm4nVDrKVw6EYZFkFc5xiFDIwMjMtMDEtMDdUMTU6MzA6MTNaal8KCmNvaW5fc3BlbnQSOQoHc3BlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYARIWCgZhbW91bnQSCjUwMDB1c3RhcnMYAWpjCg1jb2luX3JlY2VpdmVkEjoKCHJlY2VpdmVyEixzdGFyczE3eHBmdmFrbTJhbWc5NjJ5bHM2Zjg0ejNrZWxsOGM1bHk5NWFxdhgBEhYKBmFtb3VudBIKNTAwMHVzdGFycxgBapkBCgh0cmFuc2ZlchI7CglyZWNpcGllbnQSLHN0YXJzMTd4cGZ2YWttMmFtZzk2MnlsczZmODR6M2tlbGw4YzVseTk1YXF2GAESOAoGc2VuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axgBEhYKBmFtb3VudBIKNTAwMHVzdGFycxgBakMKB21lc3NhZ2USOAoGc2VuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axgBalYKAnR4EhMKA2ZlZRIKNTAwMHVzdGFycxgBEjsKCWZlZV9wYXllchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYAWpBCgJ0eBI7CgdhY2Nfc2VxEi5zdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5ay80GAFqbQoCdHgSZwoJc2lnbmF0dXJlElhTV3d5YVl1TzM3QjhMcGNMK0ZTa0ErT0hlbFY2NTFleUdtdXdUVWxpVklOSnIzUlpKQ1NVRW9hVnNOdmFPclNRc3R1YmlkVU9zcFhEb1Joa1dRVnpuQT09GAFqMwoHbWVzc2FnZRIoCgZhY3Rpb24SHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQYAWpiCgpjb2luX3NwZW50EjkKB3NwZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrGAESGQoGYW1vdW50Eg0xNTAwMDAwdXN0YXJzGAFqegoNY29pbl9yZWNlaXZlZBJOCghyZWNlaXZlchJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhgBEhkKBmFtb3VudBINMTUwMDAwMHVzdGFycxgBarABCgh0cmFuc2ZlchJPCglyZWNpcGllbnQSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24YARI4CgZzZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrGAESGQoGYW1vdW50Eg0xNTAwMDAwdXN0YXJzGAFqQwoHbWVzc2FnZRI4CgZzZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrGAFqGwoHbWVzc2FnZRIQCgZtb2R1bGUSBGJhbmsYARLQGQi6x4MDEkBENUQ2QTNDMEQ1RUNCMDQ4MTEzQTI1NjI4RTA0MDUyOEQwMURBRDEzNEY1MkVCMjM0RkE4RjQxMTYzMjY2QkIyKkAwQTFFMEExQzJGNjM2RjczNkQ2RjczMkU2MjYxNkU2QjJFNzYzMTYyNjU3NDYxMzEyRTRENzM2NzUzNjU2RTY0MpIGW3siZXZlbnRzIjpbeyJ0eXBlIjoiY29pbl9yZWNlaXZlZCIsImF0dHJpYnV0ZXMiOlt7ImtleSI6InJlY2VpdmVyIiwidmFsdWUiOiJzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuIn0seyJrZXkiOiJhbW91bnQiLCJ2YWx1ZSI6IjEwMDAwMDAwMDB1c3RhcnMifV19LHsidHlwZSI6ImNvaW5fc3BlbnQiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJzcGVuZGVyIiwidmFsdWUiOiJzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOCJ9LHsia2V5IjoiYW1vdW50IiwidmFsdWUiOiIxMDAwMDAwMDAwdXN0YXJzIn1dfSx7InR5cGUiOiJtZXNzYWdlIiwiYXR0cmlidXRlcyI6W3sia2V5IjoiYWN0aW9uIiwidmFsdWUiOiIvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kIn0seyJrZXkiOiJzZW5kZXIiLCJ2YWx1ZSI6InN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4In0seyJrZXkiOiJtb2R1bGUiLCJ2YWx1ZSI6ImJhbmsifV19LHsidHlwZSI6InRyYW5zZmVyIiwiYXR0cmlidXRlcyI6W3sia2V5IjoicmVjaXBpZW50IiwidmFsdWUiOiJzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuIn0seyJrZXkiOiJzZW5kZXIiLCJ2YWx1ZSI6InN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4In0seyJrZXkiOiJhbW91bnQiLCJ2YWx1ZSI6IjEwMDAwMDAwMDB1c3RhcnMifV19XX1dOokEGnkKDWNvaW5fcmVjZWl2ZWQSTAoIcmVjZWl2ZXISQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24SGgoGYW1vdW50EhAxMDAwMDAwMDAwdXN0YXJzGmEKCmNvaW5fc3BlbnQSNwoHc3BlbmRlchIsc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgSGgoGYW1vdW50EhAxMDAwMDAwMDAwdXN0YXJzGnkKB21lc3NhZ2USJgoGYWN0aW9uEhwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEjYKBnNlbmRlchIsc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgSDgoGbW9kdWxlEgRiYW5rGq0BCgh0cmFuc2ZlchJNCglyZWNpcGllbnQSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24SNgoGc2VuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBIaCgZhbW91bnQSEDEwMDAwMDAwMDB1c3RhcnNIwJoMUM++BFqfBAoVL2Nvc21vcy50eC52MWJldGExLlR4EoUECtUCCqcBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEoYBCixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhoUCgZ1c3RhcnMSCjEwMDAwMDAwMDASqAFJVTNGTWJFa1lxaVZlaGs0UU84VWkwNGRldEpQSVZNOGNpdmpDUldlWDRLUzI2aEFnaFJGLzZJUElWWUFrY2p6RGNRRlJmTmRWQ0lZRC9iZkxyelNJVjJmVDBNV0dVaVMxNitiWjYwVEoybEVpWjlQSVdCY1dxczJqMXA4dEdrbHlFenBqWWlHNGh6aUlXSElHaHFXVGsxWFpDb0tlY0k3UjhUVitDNDISaQpRCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAiCrM8s85a4pqMYuDsoCo3y+Ce9GjQOz2xmAWCl2ePvbEgQKAgh/GJ4BEhQKDgoGdXN0YXJzEgQyMDAwEMCaDBpA7TITP4e0ykyPdvAD/5n8uWek123FtWQzMUzwslfcdYVFBUP+IA/st+m0XpjQiRkyWx39Tt+7e65Q3WY7Oj9+UGIUMjAyMy0wMS0wN1QxNTozMjowNFpqXwoKY29pbl9zcGVudBI5CgdzcGVuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBgBEhYKBmFtb3VudBIKMjAwMHVzdGFycxgBamMKDWNvaW5fcmVjZWl2ZWQSOgoIcmVjZWl2ZXISLHN0YXJzMTd4cGZ2YWttMmFtZzk2MnlsczZmODR6M2tlbGw4YzVseTk1YXF2GAESFgoGYW1vdW50EgoyMDAwdXN0YXJzGAFqmQEKCHRyYW5zZmVyEjsKCXJlY2lwaWVudBIsc3RhcnMxN3hwZnZha20yYW1nOTYyeWxzNmY4NHoza2VsbDhjNWx5OTVhcXYYARI4CgZzZW5kZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4GAESFgoGYW1vdW50EgoyMDAwdXN0YXJzGAFqQwoHbWVzc2FnZRI4CgZzZW5kZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4GAFqVgoCdHgSEwoDZmVlEgoyMDAwdXN0YXJzGAESOwoJZmVlX3BheWVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBgBakMKAnR4Ej0KB2FjY19zZXESMHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4LzE1OBgBam0KAnR4EmcKCXNpZ25hdHVyZRJYN1RJVFA0ZTB5a3lQZHZBRC81bjh1V2VrMTIzRnRXUXpNVXp3c2xmY2RZVkZCVVArSUEvc3QrbTBYcGpRaVJreVd4MzlUdCs3ZTY1UTNXWTdPajkrVUE9PRgBajMKB21lc3NhZ2USKAoGYWN0aW9uEhwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kGAFqZQoKY29pbl9zcGVudBI5CgdzcGVuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBgBEhwKBmFtb3VudBIQMTAwMDAwMDAwMHVzdGFycxgBan0KDWNvaW5fcmVjZWl2ZWQSTgoIcmVjZWl2ZXISQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24YARIcCgZhbW91bnQSEDEwMDAwMDAwMDB1c3RhcnMYAWqzAQoIdHJhbnNmZXISTwoJcmVjaXBpZW50EkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGAESOAoGc2VuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBgBEhwKBmFtb3VudBIQMTAwMDAwMDAwMHVzdGFycxgBakMKB21lc3NhZ2USOAoGc2VuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBgBahsKB21lc3NhZ2USEAoGbW9kdWxlEgRiYW5rGAES2hgIxMqDAxJAMjQ0RjdDRkNENkY1RDY3MkJFNjM1NzM2NDhCMDQ1Rjc2QzFFRkJGRkJBMjk1MUFCMEQ4QzBFMDJCRTlBMkUwRSpAMEExRTBBMUMyRjYzNkY3MzZENkY3MzJFNjI2MTZFNkIyRTc2MzE2MjY1NzQ2MTMxMkU0RDczNjc1MzY1NkU2NDKJBlt7ImV2ZW50cyI6W3sidHlwZSI6ImNvaW5fcmVjZWl2ZWQiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJyZWNlaXZlciIsInZhbHVlIjoic3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbiJ9LHsia2V5IjoiYW1vdW50IiwidmFsdWUiOiIzNzU3NDkwdXN0YXJzIn1dfSx7InR5cGUiOiJjb2luX3NwZW50IiwiYXR0cmlidXRlcyI6W3sia2V5Ijoic3BlbmRlciIsInZhbHVlIjoic3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsifSx7ImtleSI6ImFtb3VudCIsInZhbHVlIjoiMzc1NzQ5MHVzdGFycyJ9XX0seyJ0eXBlIjoibWVzc2FnZSIsImF0dHJpYnV0ZXMiOlt7ImtleSI6ImFjdGlvbiIsInZhbHVlIjoiL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZCJ9LHsia2V5Ijoic2VuZGVyIiwidmFsdWUiOiJzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5ayJ9LHsia2V5IjoibW9kdWxlIiwidmFsdWUiOiJiYW5rIn1dfSx7InR5cGUiOiJ0cmFuc2ZlciIsImF0dHJpYnV0ZXMiOlt7ImtleSI6InJlY2lwaWVudCIsInZhbHVlIjoic3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbiJ9LHsia2V5Ijoic2VuZGVyIiwidmFsdWUiOiJzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5ayJ9LHsia2V5IjoiYW1vdW50IiwidmFsdWUiOiIzNzU3NDkwdXN0YXJzIn1dfV19XTqABBp2Cg1jb2luX3JlY2VpdmVkEkwKCHJlY2VpdmVyEkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuEhcKBmFtb3VudBINMzc1NzQ5MHVzdGFycxpeCgpjb2luX3NwZW50EjcKB3NwZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrEhcKBmFtb3VudBINMzc1NzQ5MHVzdGFycxp5CgdtZXNzYWdlEiYKBmFjdGlvbhIcL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBI2CgZzZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrEg4KBm1vZHVsZRIEYmFuaxqqAQoIdHJhbnNmZXISTQoJcmVjaXBpZW50EkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuEjYKBnNlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsSFwoGYW1vdW50Eg0zNzU3NDkwdXN0YXJzSMCaDFDPtARaxgMKFS9jb3Ntb3MudHgudjFiZXRhMS5UeBKsAwr9AQqkAQocL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBKDAQosc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24aEQoGdXN0YXJzEgczNzU3NDkwElRRZ0VPemtSNlJqY1NOa2o5NnE4RTRYeDBUdmluUWdzTWIvUUsyYktTbW1wa1dCTkVrMjFQeXZXTlFneFIybkNDK1cvTnZsUURIcG93a0NaWXgyRGoSaApQCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAkv5JPXz/sLcxTj1NSrfNcBEFtLCKHe76OAO0yWH7KZCEgQKAgh/GAUSFAoOCgZ1c3RhcnMSBDUwMDAQwJoMGkB8qJ0RJN2BrOKGFTmofRA3DF8RLgbShe3I1UBK/RWGAQFZ9kJSL0e2PyxMpv1qYsWgEy3ohsIGjqIYImFIwZRCYhQyMDIzLTAxLTA3VDE2OjEwOjUwWmpfCgpjb2luX3NwZW50EjkKB3NwZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrGAESFgoGYW1vdW50Ego1MDAwdXN0YXJzGAFqYwoNY29pbl9yZWNlaXZlZBI6CghyZWNlaXZlchIsc3RhcnMxN3hwZnZha20yYW1nOTYyeWxzNmY4NHoza2VsbDhjNWx5OTVhcXYYARIWCgZhbW91bnQSCjUwMDB1c3RhcnMYAWqZAQoIdHJhbnNmZXISOwoJcmVjaXBpZW50EixzdGFyczE3eHBmdmFrbTJhbWc5NjJ5bHM2Zjg0ejNrZWxsOGM1bHk5NWFxdhgBEjgKBnNlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYARIWCgZhbW91bnQSCjUwMDB1c3RhcnMYAWpDCgdtZXNzYWdlEjgKBnNlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYAWpWCgJ0eBITCgNmZWUSCjUwMDB1c3RhcnMYARI7CglmZWVfcGF5ZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrGAFqQQoCdHgSOwoHYWNjX3NlcRIuc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsvNRgBam0KAnR4EmcKCXNpZ25hdHVyZRJYZktpZEVTVGRnYXppaGhVNXFIMFFOd3hmRVM0RzBvWHR5TlZBU3YwVmhnRUJXZlpDVWk5SHRqOHNUS2I5YW1MRm9CTXQ2SWJDQm82aUdDSmhTTUdVUWc9PRgBajMKB21lc3NhZ2USKAoGYWN0aW9uEhwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kGAFqYgoKY29pbl9zcGVudBI5CgdzcGVuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axgBEhkKBmFtb3VudBINMzc1NzQ5MHVzdGFycxgBanoKDWNvaW5fcmVjZWl2ZWQSTgoIcmVjZWl2ZXISQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24YARIZCgZhbW91bnQSDTM3NTc0OTB1c3RhcnMYAWqwAQoIdHJhbnNmZXISTwoJcmVjaXBpZW50EkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGAESOAoGc2VuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axgBEhkKBmFtb3VudBINMzc1NzQ5MHVzdGFycxgBakMKB21lc3NhZ2USOAoGc2VuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axgBahsKB21lc3NhZ2USEAoGbW9kdWxlEgRiYW5rGAESvhgIs9WDAxJAMTIwODlFQTI4NTU1NjE0NzhFQzA0QjVGQTczOTg1ODVFOTlCMkUyRDYzNUE0ODI2RUY3NTJBQUM3NkMxREE1NCpAMEExRTBBMUMyRjYzNkY3MzZENkY3MzJFNjI2MTZFNkIyRTc2MzE2MjY1NzQ2MTMxMkU0RDczNjc1MzY1NkU2NDKJBlt7ImV2ZW50cyI6W3sidHlwZSI6ImNvaW5fcmVjZWl2ZWQiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJyZWNlaXZlciIsInZhbHVlIjoic3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbiJ9LHsia2V5IjoiYW1vdW50IiwidmFsdWUiOiIxNzY4NTQwdXN0YXJzIn1dfSx7InR5cGUiOiJjb2luX3NwZW50IiwiYXR0cmlidXRlcyI6W3sia2V5Ijoic3BlbmRlciIsInZhbHVlIjoic3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsifSx7ImtleSI6ImFtb3VudCIsInZhbHVlIjoiMTc2ODU0MHVzdGFycyJ9XX0seyJ0eXBlIjoibWVzc2FnZSIsImF0dHJpYnV0ZXMiOlt7ImtleSI6ImFjdGlvbiIsInZhbHVlIjoiL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZCJ9LHsia2V5Ijoic2VuZGVyIiwidmFsdWUiOiJzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5ayJ9LHsia2V5IjoibW9kdWxlIiwidmFsdWUiOiJiYW5rIn1dfSx7InR5cGUiOiJ0cmFuc2ZlciIsImF0dHJpYnV0ZXMiOlt7ImtleSI6InJlY2lwaWVudCIsInZhbHVlIjoic3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbiJ9LHsia2V5Ijoic2VuZGVyIiwidmFsdWUiOiJzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5ayJ9LHsia2V5IjoiYW1vdW50IiwidmFsdWUiOiIxNzY4NTQwdXN0YXJzIn1dfV19XTqABBp2Cg1jb2luX3JlY2VpdmVkEkwKCHJlY2VpdmVyEkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuEhcKBmFtb3VudBINMTc2ODU0MHVzdGFycxpeCgpjb2luX3NwZW50EjcKB3NwZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrEhcKBmFtb3VudBINMTc2ODU0MHVzdGFycxp5CgdtZXNzYWdlEiYKBmFjdGlvbhIcL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBI2CgZzZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrEg4KBm1vZHVsZRIEYmFuaxqqAQoIdHJhbnNmZXISTQoJcmVjaXBpZW50EkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuEjYKBnNlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsSFwoGYW1vdW50Eg0xNzY4NTQwdXN0YXJzSMCaDFCCswRaqgMKFS9jb3Ntb3MudHgudjFiZXRhMS5UeBKQAwrhAQqkAQocL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBKDAQosc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24aEQoGdXN0YXJzEgcxNzY4NTQwEjhaQUVPemtSNlJqY1NOa2o5NnE4RTRYeDBUdmluWkFNVXRzdDk4M0NOUzJqcm5OZjNzcFlNY293WRJoClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiECS/kk9fP+wtzFOPU1Kt81wEQW0sIod7vo4A7TJYfspkISBAoCCH8YBhIUCg4KBnVzdGFycxIENTAwMBDAmgwaQB2zAkk+TzGcmQ+7hRlibIY61B/w0/3f5sIr7+Rf7kDifli4iQAQdO8VT9Q8941A+gH7oVKf7AcmhRqg+tfitCliFDIwMjMtMDEtMDdUMTg6Mjc6NDRaal8KCmNvaW5fc3BlbnQSOQoHc3BlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYARIWCgZhbW91bnQSCjUwMDB1c3RhcnMYAWpjCg1jb2luX3JlY2VpdmVkEjoKCHJlY2VpdmVyEixzdGFyczE3eHBmdmFrbTJhbWc5NjJ5bHM2Zjg0ejNrZWxsOGM1bHk5NWFxdhgBEhYKBmFtb3VudBIKNTAwMHVzdGFycxgBapkBCgh0cmFuc2ZlchI7CglyZWNpcGllbnQSLHN0YXJzMTd4cGZ2YWttMmFtZzk2MnlsczZmODR6M2tlbGw4YzVseTk1YXF2GAESOAoGc2VuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axgBEhYKBmFtb3VudBIKNTAwMHVzdGFycxgBakMKB21lc3NhZ2USOAoGc2VuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axgBalYKAnR4EhMKA2ZlZRIKNTAwMHVzdGFycxgBEjsKCWZlZV9wYXllchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYAWpBCgJ0eBI7CgdhY2Nfc2VxEi5zdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5ay82GAFqbQoCdHgSZwoJc2lnbmF0dXJlElhIYk1DU1Q1UE1aeVpEN3VGR1dKc2hqclVIL0RUL2QvbXdpdnY1Ri91UU9KK1dMaUpBQkIwN3hWUDFEejNqVUQ2QWZ1aFVwL3NCeWFGR3FENjErSzBLUT09GAFqMwoHbWVzc2FnZRIoCgZhY3Rpb24SHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQYAWpiCgpjb2luX3NwZW50EjkKB3NwZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrGAESGQoGYW1vdW50Eg0xNzY4NTQwdXN0YXJzGAFqegoNY29pbl9yZWNlaXZlZBJOCghyZWNlaXZlchJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhgBEhkKBmFtb3VudBINMTc2ODU0MHVzdGFycxgBarABCgh0cmFuc2ZlchJPCglyZWNpcGllbnQSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24YARI4CgZzZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrGAESGQoGYW1vdW50Eg0xNzY4NTQwdXN0YXJzGAFqQwoHbWVzc2FnZRI4CgZzZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrGAFqGwoHbWVzc2FnZRIQCgZtb2R1bGUSBGJhbmsYARKiGAju4oMDEkBBNkM0NjZFNDM5NUZBNUQzMTM4QzA5NjcwOUE2NEUxNjc1NENDNjNFRTI5QzVCMDM4Rjc4OTFBMTlCRTYxQTQyKkAwQTFFMEExQzJGNjM2RjczNkQ2RjczMkU2MjYxNkU2QjJFNzYzMTYyNjU3NDYxMzEyRTRENzM2NzUzNjU2RTY0MokGW3siZXZlbnRzIjpbeyJ0eXBlIjoiY29pbl9yZWNlaXZlZCIsImF0dHJpYnV0ZXMiOlt7ImtleSI6InJlY2VpdmVyIiwidmFsdWUiOiJzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuIn0seyJrZXkiOiJhbW91bnQiLCJ2YWx1ZSI6IjIwMDAwMDB1c3RhcnMifV19LHsidHlwZSI6ImNvaW5fc3BlbnQiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJzcGVuZGVyIiwidmFsdWUiOiJzdGFyczF2ejkyd2Y2N2tzZG5zbWNqZXVlNngyempzZmxkcDlnOXk4ZnF5NyJ9LHsia2V5IjoiYW1vdW50IiwidmFsdWUiOiIyMDAwMDAwdXN0YXJzIn1dfSx7InR5cGUiOiJtZXNzYWdlIiwiYXR0cmlidXRlcyI6W3sia2V5IjoiYWN0aW9uIiwidmFsdWUiOiIvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kIn0seyJrZXkiOiJzZW5kZXIiLCJ2YWx1ZSI6InN0YXJzMXZ6OTJ3ZjY3a3NkbnNtY2pldWU2eDJ6anNmbGRwOWc5eThmcXk3In0seyJrZXkiOiJtb2R1bGUiLCJ2YWx1ZSI6ImJhbmsifV19LHsidHlwZSI6InRyYW5zZmVyIiwiYXR0cmlidXRlcyI6W3sia2V5IjoicmVjaXBpZW50IiwidmFsdWUiOiJzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuIn0seyJrZXkiOiJzZW5kZXIiLCJ2YWx1ZSI6InN0YXJzMXZ6OTJ3ZjY3a3NkbnNtY2pldWU2eDJ6anNmbGRwOWc5eThmcXk3In0seyJrZXkiOiJhbW91bnQiLCJ2YWx1ZSI6IjIwMDAwMDB1c3RhcnMifV19XX1dOoAEGnYKDWNvaW5fcmVjZWl2ZWQSTAoIcmVjZWl2ZXISQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24SFwoGYW1vdW50Eg0yMDAwMDAwdXN0YXJzGl4KCmNvaW5fc3BlbnQSNwoHc3BlbmRlchIsc3RhcnMxdno5MndmNjdrc2Ruc21jamV1ZTZ4Mnpqc2ZsZHA5Zzl5OGZxeTcSFwoGYW1vdW50Eg0yMDAwMDAwdXN0YXJzGnkKB21lc3NhZ2USJgoGYWN0aW9uEhwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEjYKBnNlbmRlchIsc3RhcnMxdno5MndmNjdrc2Ruc21jamV1ZTZ4Mnpqc2ZsZHA5Zzl5OGZxeTcSDgoGbW9kdWxlEgRiYW5rGqoBCgh0cmFuc2ZlchJNCglyZWNpcGllbnQSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24SNgoGc2VuZGVyEixzdGFyczF2ejkyd2Y2N2tzZG5zbWNqZXVlNngyempzZmxkcDlnOXk4ZnF5NxIXCgZhbW91bnQSDTIwMDAwMDB1c3RhcnNIwJoMUN2vBFqOAwoVL2Nvc21vcy50eC52MWJldGExLlR4EvQCCsUBCqQBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEoMBCixzdGFyczF2ejkyd2Y2N2tzZG5zbWNqZXVlNngyempzZmxkcDlnOXk4ZnF5NxJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhoRCgZ1c3RhcnMSBzIwMDAwMDASHHlDUndZWXMrdHV5NUZqMWJhcVFVYVVrUWRVVVASaApQCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAsmlkFqDV0W03lkz9P7aMlKoSQ8rUqtty0dfhKl/ZlV0EgQKAgh/GAESFAoOCgZ1c3RhcnMSBDUwMDAQwJoMGkB4vN3wF5DL8eMyW70xjh9N8DcfeqAs+6xgINUT+r8YAmdWesMKZigIJOetLsxe/0fiyoY0++n5/jRg0c+4vSqzYhQyMDIzLTAxLTA3VDIxOjE3OjAyWmpfCgpjb2luX3NwZW50EjkKB3NwZW5kZXISLHN0YXJzMXZ6OTJ3ZjY3a3NkbnNtY2pldWU2eDJ6anNmbGRwOWc5eThmcXk3GAESFgoGYW1vdW50Ego1MDAwdXN0YXJzGAFqYwoNY29pbl9yZWNlaXZlZBI6CghyZWNlaXZlchIsc3RhcnMxN3hwZnZha20yYW1nOTYyeWxzNmY4NHoza2VsbDhjNWx5OTVhcXYYARIWCgZhbW91bnQSCjUwMDB1c3RhcnMYAWqZAQoIdHJhbnNmZXISOwoJcmVjaXBpZW50EixzdGFyczE3eHBmdmFrbTJhbWc5NjJ5bHM2Zjg0ejNrZWxsOGM1bHk5NWFxdhgBEjgKBnNlbmRlchIsc3RhcnMxdno5MndmNjdrc2Ruc21jamV1ZTZ4Mnpqc2ZsZHA5Zzl5OGZxeTcYARIWCgZhbW91bnQSCjUwMDB1c3RhcnMYAWpDCgdtZXNzYWdlEjgKBnNlbmRlchIsc3RhcnMxdno5MndmNjdrc2Ruc21jamV1ZTZ4Mnpqc2ZsZHA5Zzl5OGZxeTcYAWpWCgJ0eBITCgNmZWUSCjUwMDB1c3RhcnMYARI7CglmZWVfcGF5ZXISLHN0YXJzMXZ6OTJ3ZjY3a3NkbnNtY2pldWU2eDJ6anNmbGRwOWc5eThmcXk3GAFqQQoCdHgSOwoHYWNjX3NlcRIuc3RhcnMxdno5MndmNjdrc2Ruc21jamV1ZTZ4Mnpqc2ZsZHA5Zzl5OGZxeTcvMRgBam0KAnR4EmcKCXNpZ25hdHVyZRJYZUx6ZDhCZVF5L0hqTWx1OU1ZNGZUZkEzSDNxZ0xQdXNZQ0RWRS9xL0dBSm5WbnJEQ21Zb0NDVG5yUzdNWHY5SDRzcUdOUHZwK2Y0MFlOSFB1TDBxc3c9PRgBajMKB21lc3NhZ2USKAoGYWN0aW9uEhwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kGAFqYgoKY29pbl9zcGVudBI5CgdzcGVuZGVyEixzdGFyczF2ejkyd2Y2N2tzZG5zbWNqZXVlNngyempzZmxkcDlnOXk4ZnF5NxgBEhkKBmFtb3VudBINMjAwMDAwMHVzdGFycxgBanoKDWNvaW5fcmVjZWl2ZWQSTgoIcmVjZWl2ZXISQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24YARIZCgZhbW91bnQSDTIwMDAwMDB1c3RhcnMYAWqwAQoIdHJhbnNmZXISTwoJcmVjaXBpZW50EkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGAESOAoGc2VuZGVyEixzdGFyczF2ejkyd2Y2N2tzZG5zbWNqZXVlNngyempzZmxkcDlnOXk4ZnF5NxgBEhkKBmFtb3VudBINMjAwMDAwMHVzdGFycxgBakMKB21lc3NhZ2USOAoGc2VuZGVyEixzdGFyczF2ejkyd2Y2N2tzZG5zbWNqZXVlNngyempzZmxkcDlnOXk4ZnF5NxgBahsKB21lc3NhZ2USEAoGbW9kdWxlEgRiYW5rGAES9hgIws2EAxJARjIxM0MyMUExMzFCNjIwREJGMDBERTQwOTY1MEUwMjJGMzAzMUVFMDJCQzI2MjhDNjQ2Mjk3MjgwNDgwMzRDNypAMEExRTBBMUMyRjYzNkY3MzZENkY3MzJFNjI2MTZFNkIyRTc2MzE2MjY1NzQ2MTMxMkU0RDczNjc1MzY1NkU2NDKJBlt7ImV2ZW50cyI6W3sidHlwZSI6ImNvaW5fcmVjZWl2ZWQiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJyZWNlaXZlciIsInZhbHVlIjoic3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbiJ9LHsia2V5IjoiYW1vdW50IiwidmFsdWUiOiIyMDAwMDAwdXN0YXJzIn1dfSx7InR5cGUiOiJjb2luX3NwZW50IiwiYXR0cmlidXRlcyI6W3sia2V5Ijoic3BlbmRlciIsInZhbHVlIjoic3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsifSx7ImtleSI6ImFtb3VudCIsInZhbHVlIjoiMjAwMDAwMHVzdGFycyJ9XX0seyJ0eXBlIjoibWVzc2FnZSIsImF0dHJpYnV0ZXMiOlt7ImtleSI6ImFjdGlvbiIsInZhbHVlIjoiL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZCJ9LHsia2V5Ijoic2VuZGVyIiwidmFsdWUiOiJzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5ayJ9LHsia2V5IjoibW9kdWxlIiwidmFsdWUiOiJiYW5rIn1dfSx7InR5cGUiOiJ0cmFuc2ZlciIsImF0dHJpYnV0ZXMiOlt7ImtleSI6InJlY2lwaWVudCIsInZhbHVlIjoic3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbiJ9LHsia2V5Ijoic2VuZGVyIiwidmFsdWUiOiJzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5ayJ9LHsia2V5IjoiYW1vdW50IiwidmFsdWUiOiIyMDAwMDAwdXN0YXJzIn1dfV19XTqABBp2Cg1jb2luX3JlY2VpdmVkEkwKCHJlY2VpdmVyEkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuEhcKBmFtb3VudBINMjAwMDAwMHVzdGFycxpeCgpjb2luX3NwZW50EjcKB3NwZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrEhcKBmFtb3VudBINMjAwMDAwMHVzdGFycxp5CgdtZXNzYWdlEiYKBmFjdGlvbhIcL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBI2CgZzZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrEg4KBm1vZHVsZRIEYmFuaxqqAQoIdHJhbnNmZXISTQoJcmVjaXBpZW50EkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuEjYKBnNlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsSFwoGYW1vdW50Eg0yMDAwMDAwdXN0YXJzSMCaDFCvtwRa4gMKFS9jb3Ntb3MudHgudjFiZXRhMS5UeBLIAwqZAgqkAQocL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBKDAQosc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24aEQoGdXN0YXJzEgcyMDAwMDAwEnBNZ0VPemtSNlJqY1NOa2o5NnE4RTRYeDBUdmluTWdNVXRzdDk4M0NOUzJqcm5OZjNzcFlNY293WU1neFIybkNDK1cvTnZsUURIcG93a0NaWXgyRGpNZ2s1a2NvalZ1TE0rQlJSL0Q0NEE2QU4vc0pZEmgKUApGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQJL+ST18/7C3MU49TUq3zXARBbSwih3u+jgDtMlh+ymQhIECgIIfxgHEhQKDgoGdXN0YXJzEgQ1MDAwEMCaDBpAcSiS1G/lpK4r8kRuXi7G/xp62wVVpWIrV4KNpiJ1BVcsytvK7pPuY7C3cgsv12spakgt+uFzBmkFXC7eCs9rKWIUMjAyMy0wMS0wOFQxOTozODowMVpqXwoKY29pbl9zcGVudBI5CgdzcGVuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axgBEhYKBmFtb3VudBIKNTAwMHVzdGFycxgBamMKDWNvaW5fcmVjZWl2ZWQSOgoIcmVjZWl2ZXISLHN0YXJzMTd4cGZ2YWttMmFtZzk2MnlsczZmODR6M2tlbGw4YzVseTk1YXF2GAESFgoGYW1vdW50Ego1MDAwdXN0YXJzGAFqmQEKCHRyYW5zZmVyEjsKCXJlY2lwaWVudBIsc3RhcnMxN3hwZnZha20yYW1nOTYyeWxzNmY4NHoza2VsbDhjNWx5OTVhcXYYARI4CgZzZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrGAESFgoGYW1vdW50Ego1MDAwdXN0YXJzGAFqQwoHbWVzc2FnZRI4CgZzZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrGAFqVgoCdHgSEwoDZmVlEgo1MDAwdXN0YXJzGAESOwoJZmVlX3BheWVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axgBakEKAnR4EjsKB2FjY19zZXESLnN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrLzcYAWptCgJ0eBJnCglzaWduYXR1cmUSWGNTaVMxRy9scEs0cjhrUnVYaTdHL3hwNjJ3VlZwV0lyVjRLTnBpSjFCVmNzeXR2SzdwUHVZN0MzY2dzdjEyc3Bha2d0K3VGekJta0ZYQzdlQ3M5cktRPT0YAWozCgdtZXNzYWdlEigKBmFjdGlvbhIcL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBgBamIKCmNvaW5fc3BlbnQSOQoHc3BlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYARIZCgZhbW91bnQSDTIwMDAwMDB1c3RhcnMYAWp6Cg1jb2luX3JlY2VpdmVkEk4KCHJlY2VpdmVyEkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGAESGQoGYW1vdW50Eg0yMDAwMDAwdXN0YXJzGAFqsAEKCHRyYW5zZmVyEk8KCXJlY2lwaWVudBJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhgBEjgKBnNlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYARIZCgZhbW91bnQSDTIwMDAwMDB1c3RhcnMYAWpDCgdtZXNzYWdlEjgKBnNlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYAWobCgdtZXNzYWdlEhAKBm1vZHVsZRIEYmFuaxgBEq8ZCMKChQMSQDcyRjhCNUNDMUY4NDU3Q0Y0OThDNEEyNTVDMzY2MDBFMDUyMzNEOTY0QjI1QkQ3QjgzNjA5MUUxQkM5N0M3MDEqQDBBMUUwQTFDMkY2MzZGNzM2RDZGNzMyRTYyNjE2RTZCMkU3NjMxNjI2NTc0NjEzMTJFNEQ3MzY3NTM2NTZFNjQyiQZbeyJldmVudHMiOlt7InR5cGUiOiJjb2luX3JlY2VpdmVkIiwiYXR0cmlidXRlcyI6W3sia2V5IjoicmVjZWl2ZXIiLCJ2YWx1ZSI6InN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24ifSx7ImtleSI6ImFtb3VudCIsInZhbHVlIjoiMzAwMDAwMHVzdGFycyJ9XX0seyJ0eXBlIjoiY29pbl9zcGVudCIsImF0dHJpYnV0ZXMiOlt7ImtleSI6InNwZW5kZXIiLCJ2YWx1ZSI6InN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrIn0seyJrZXkiOiJhbW91bnQiLCJ2YWx1ZSI6IjMwMDAwMDB1c3RhcnMifV19LHsidHlwZSI6Im1lc3NhZ2UiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJhY3Rpb24iLCJ2YWx1ZSI6Ii9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQifSx7ImtleSI6InNlbmRlciIsInZhbHVlIjoic3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsifSx7ImtleSI6Im1vZHVsZSIsInZhbHVlIjoiYmFuayJ9XX0seyJ0eXBlIjoidHJhbnNmZXIiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJyZWNpcGllbnQiLCJ2YWx1ZSI6InN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24ifSx7ImtleSI6InNlbmRlciIsInZhbHVlIjoic3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsifSx7ImtleSI6ImFtb3VudCIsInZhbHVlIjoiMzAwMDAwMHVzdGFycyJ9XX1dfV06gAQadgoNY29pbl9yZWNlaXZlZBJMCghyZWNlaXZlchJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhIXCgZhbW91bnQSDTMwMDAwMDB1c3RhcnMaXgoKY29pbl9zcGVudBI3CgdzcGVuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axIXCgZhbW91bnQSDTMwMDAwMDB1c3RhcnMaeQoHbWVzc2FnZRImCgZhY3Rpb24SHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQSNgoGc2VuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axIOCgZtb2R1bGUSBGJhbmsaqgEKCHRyYW5zZmVyEk0KCXJlY2lwaWVudBJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhI2CgZzZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrEhcKBmFtb3VudBINMzAwMDAwMHVzdGFyc0jAmgxQy7sEWpsEChUvY29zbW9zLnR4LnYxYmV0YTEuVHgSgQQK0gIKpAEKHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQSgwEKLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrEkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGhEKBnVzdGFycxIHMzAwMDAwMBKoAUlRRU96a1I2UmpjU05rajk2cThFNFh4MFR2aW5JUVpQemlhRzBYRkZhdldZak9aWDYwckg2TFBBSVFrNWtjb2pWdUxNK0JSUi9ENDRBNkFOL3NKWUlSTVRUQUdWSVNkcTBOZ0lBMGxCOVE4b2NNcjJJUmczdnByQ0pHWEpJNnVhbmQvZWNrbFVJbWZiSVI2NUF5RGJHTFJKeUVvbE1pcHc2NFJTSGEvahJoClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiECS/kk9fP+wtzFOPU1Kt81wEQW0sIod7vo4A7TJYfspkISBAoCCH8YCBIUCg4KBnVzdGFycxIENTAwMBDAmgwaQN1uUMSWW8LBqyDcT2/vK7n7jfKzEzWCGcC4aOrxsAKKCYGUcp9E8b78+HJhb975F5Ivv9AT1bvpUfmDSv5SRC1iFDIwMjMtMDEtMDlUMDY6NDQ6MTFaal8KCmNvaW5fc3BlbnQSOQoHc3BlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYARIWCgZhbW91bnQSCjUwMDB1c3RhcnMYAWpjCg1jb2luX3JlY2VpdmVkEjoKCHJlY2VpdmVyEixzdGFyczE3eHBmdmFrbTJhbWc5NjJ5bHM2Zjg0ejNrZWxsOGM1bHk5NWFxdhgBEhYKBmFtb3VudBIKNTAwMHVzdGFycxgBapkBCgh0cmFuc2ZlchI7CglyZWNpcGllbnQSLHN0YXJzMTd4cGZ2YWttMmFtZzk2MnlsczZmODR6M2tlbGw4YzVseTk1YXF2GAESOAoGc2VuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axgBEhYKBmFtb3VudBIKNTAwMHVzdGFycxgBakMKB21lc3NhZ2USOAoGc2VuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axgBalYKAnR4EhMKA2ZlZRIKNTAwMHVzdGFycxgBEjsKCWZlZV9wYXllchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYAWpBCgJ0eBI7CgdhY2Nfc2VxEi5zdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5ay84GAFqbQoCdHgSZwoJc2lnbmF0dXJlElgzVzVReEpaYndzR3JJTnhQYis4cnVmdU44ck1UTllJWndMaG82dkd3QW9vSmdaUnluMFR4dnZ6NGNtRnYzdmtYa2krLzBCUFZ1K2xSK1lOSy9sSkVMUT09GAFqMwoHbWVzc2FnZRIoCgZhY3Rpb24SHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQYAWpiCgpjb2luX3NwZW50EjkKB3NwZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrGAESGQoGYW1vdW50Eg0zMDAwMDAwdXN0YXJzGAFqegoNY29pbl9yZWNlaXZlZBJOCghyZWNlaXZlchJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhgBEhkKBmFtb3VudBINMzAwMDAwMHVzdGFycxgBarABCgh0cmFuc2ZlchJPCglyZWNpcGllbnQSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24YARI4CgZzZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrGAESGQoGYW1vdW50Eg0zMDAwMDAwdXN0YXJzGAFqQwoHbWVzc2FnZRI4CgZzZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrGAFqGwoHbWVzc2FnZRIQCgZtb2R1bGUSBGJhbmsYARLnGQiqhYUDEkA3MDlFOUE2RjVGRUMzNzhEQkIwNjQ0OTNDMjcxQTg1RTIwQjY3OTczREIwNkY4MjFDOEJFNkU3Qzg5MzdCNDhBKkAwQTFFMEExQzJGNjM2RjczNkQ2RjczMkU2MjYxNkU2QjJFNzYzMTYyNjU3NDYxMzEyRTRENzM2NzUzNjU2RTY0MokGW3siZXZlbnRzIjpbeyJ0eXBlIjoiY29pbl9yZWNlaXZlZCIsImF0dHJpYnV0ZXMiOlt7ImtleSI6InJlY2VpdmVyIiwidmFsdWUiOiJzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuIn0seyJrZXkiOiJhbW91bnQiLCJ2YWx1ZSI6IjEwMDAwMDB1c3RhcnMifV19LHsidHlwZSI6ImNvaW5fc3BlbnQiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJzcGVuZGVyIiwidmFsdWUiOiJzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5ayJ9LHsia2V5IjoiYW1vdW50IiwidmFsdWUiOiIxMDAwMDAwdXN0YXJzIn1dfSx7InR5cGUiOiJtZXNzYWdlIiwiYXR0cmlidXRlcyI6W3sia2V5IjoiYWN0aW9uIiwidmFsdWUiOiIvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kIn0seyJrZXkiOiJzZW5kZXIiLCJ2YWx1ZSI6InN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrIn0seyJrZXkiOiJtb2R1bGUiLCJ2YWx1ZSI6ImJhbmsifV19LHsidHlwZSI6InRyYW5zZmVyIiwiYXR0cmlidXRlcyI6W3sia2V5IjoicmVjaXBpZW50IiwidmFsdWUiOiJzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuIn0seyJrZXkiOiJzZW5kZXIiLCJ2YWx1ZSI6InN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrIn0seyJrZXkiOiJhbW91bnQiLCJ2YWx1ZSI6IjEwMDAwMDB1c3RhcnMifV19XX1dOoAEGnYKDWNvaW5fcmVjZWl2ZWQSTAoIcmVjZWl2ZXISQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24SFwoGYW1vdW50Eg0xMDAwMDAwdXN0YXJzGl4KCmNvaW5fc3BlbnQSNwoHc3BlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsSFwoGYW1vdW50Eg0xMDAwMDAwdXN0YXJzGnkKB21lc3NhZ2USJgoGYWN0aW9uEhwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEjYKBnNlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsSDgoGbW9kdWxlEgRiYW5rGqoBCgh0cmFuc2ZlchJNCglyZWNpcGllbnQSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24SNgoGc2VuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axIXCgZhbW91bnQSDTEwMDAwMDB1c3RhcnNIwJoMUNG/BFrTBAoVL2Nvc21vcy50eC52MWJldGExLlR4ErkECooDCqQBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEoMBCixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhoRCgZ1c3RhcnMSBzEwMDAwMDAS4AFHZlhVYU12Y3JVRmhnZkN4K2ZlTlFCMXE4YmtTR2YyMTNhR2o5elFVMG5Ib0l2VW01WjJ1OWplQUdmU1dtS0IwaUphSlMvajUyV3V4azNqNUJkMmdHZk5uVEtwNmRwNFJFTmtIaTlDTjlCNmtKTHZPR2ZDVVNjLzdQbGk1QkV5RG9xaVYzcVY5VWVQZEdlam1kQzJ4a3EvRjRWTC9NcHFmamVKaSs2K2tHZWZZamYxck5jYzk1Unlub1NmTlFNZU14WTVsR2VHOHhCUTNxSFJ3bzlsWDZFbjU4VGYvbkVUbhJoClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiECS/kk9fP+wtzFOPU1Kt81wEQW0sIod7vo4A7TJYfspkISBAoCCH8YCRIUCg4KBnVzdGFycxIENTAwMBDAmgwaQK3/087C0NaPaBqtwPhK/GBPxeyCuBGIITdk+WEHy/s7Qn0v8gKH7o0vLpvjCGGMnQhtJVw70p9hdxaK0iySY2ViFDIwMjMtMDEtMDlUMDc6MTk6MzFaal8KCmNvaW5fc3BlbnQSOQoHc3BlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYARIWCgZhbW91bnQSCjUwMDB1c3RhcnMYAWpjCg1jb2luX3JlY2VpdmVkEjoKCHJlY2VpdmVyEixzdGFyczE3eHBmdmFrbTJhbWc5NjJ5bHM2Zjg0ejNrZWxsOGM1bHk5NWFxdhgBEhYKBmFtb3VudBIKNTAwMHVzdGFycxgBapkBCgh0cmFuc2ZlchI7CglyZWNpcGllbnQSLHN0YXJzMTd4cGZ2YWttMmFtZzk2MnlsczZmODR6M2tlbGw4YzVseTk1YXF2GAESOAoGc2VuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axgBEhYKBmFtb3VudBIKNTAwMHVzdGFycxgBakMKB21lc3NhZ2USOAoGc2VuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axgBalYKAnR4EhMKA2ZlZRIKNTAwMHVzdGFycxgBEjsKCWZlZV9wYXllchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYAWpBCgJ0eBI7CgdhY2Nfc2VxEi5zdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5ay85GAFqbQoCdHgSZwoJc2lnbmF0dXJlElhyZi9UenNMUTFvOW9HcTNBK0VyOFlFL0Y3SUs0RVlnaE4yVDVZUWZMK3p0Q2ZTL3lBb2Z1alM4dW0rTUlZWXlkQ0cwbFhEdlNuMkYzRm9yU0xKSmpaUT09GAFqMwoHbWVzc2FnZRIoCgZhY3Rpb24SHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQYAWpiCgpjb2luX3NwZW50EjkKB3NwZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrGAESGQoGYW1vdW50Eg0xMDAwMDAwdXN0YXJzGAFqegoNY29pbl9yZWNlaXZlZBJOCghyZWNlaXZlchJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhgBEhkKBmFtb3VudBINMTAwMDAwMHVzdGFycxgBarABCgh0cmFuc2ZlchJPCglyZWNpcGllbnQSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24YARI4CgZzZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrGAESGQoGYW1vdW50Eg0xMDAwMDAwdXN0YXJzGAFqQwoHbWVzc2FnZRI4CgZzZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrGAFqGwoHbWVzc2FnZRIQCgZtb2R1bGUSBGJhbmsYARKjGAjXpIUDEkBEREY1MDU1NkU0NjIxNjQyMDZGMTY2Q0FBREQ5RUM5M0ZENzExNDE4Q0NFNEZCNTZBODQ4RTY2RkExMTc1OTUzKkAwQTFFMEExQzJGNjM2RjczNkQ2RjczMkU2MjYxNkU2QjJFNzYzMTYyNjU3NDYxMzEyRTRENzM2NzUzNjU2RTY0MokGW3siZXZlbnRzIjpbeyJ0eXBlIjoiY29pbl9yZWNlaXZlZCIsImF0dHJpYnV0ZXMiOlt7ImtleSI6InJlY2VpdmVyIiwidmFsdWUiOiJzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuIn0seyJrZXkiOiJhbW91bnQiLCJ2YWx1ZSI6IjEwMDAwMDB1c3RhcnMifV19LHsidHlwZSI6ImNvaW5fc3BlbnQiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJzcGVuZGVyIiwidmFsdWUiOiJzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5ayJ9LHsia2V5IjoiYW1vdW50IiwidmFsdWUiOiIxMDAwMDAwdXN0YXJzIn1dfSx7InR5cGUiOiJtZXNzYWdlIiwiYXR0cmlidXRlcyI6W3sia2V5IjoiYWN0aW9uIiwidmFsdWUiOiIvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kIn0seyJrZXkiOiJzZW5kZXIiLCJ2YWx1ZSI6InN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrIn0seyJrZXkiOiJtb2R1bGUiLCJ2YWx1ZSI6ImJhbmsifV19LHsidHlwZSI6InRyYW5zZmVyIiwiYXR0cmlidXRlcyI6W3sia2V5IjoicmVjaXBpZW50IiwidmFsdWUiOiJzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuIn0seyJrZXkiOiJzZW5kZXIiLCJ2YWx1ZSI6InN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrIn0seyJrZXkiOiJhbW91bnQiLCJ2YWx1ZSI6IjEwMDAwMDB1c3RhcnMifV19XX1dOoAEGnYKDWNvaW5fcmVjZWl2ZWQSTAoIcmVjZWl2ZXISQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24SFwoGYW1vdW50Eg0xMDAwMDAwdXN0YXJzGl4KCmNvaW5fc3BlbnQSNwoHc3BlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsSFwoGYW1vdW50Eg0xMDAwMDAwdXN0YXJzGnkKB21lc3NhZ2USJgoGYWN0aW9uEhwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEjYKBnNlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsSDgoGbW9kdWxlEgRiYW5rGqoBCgh0cmFuc2ZlchJNCglyZWNpcGllbnQSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24SNgoGc2VuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axIXCgZhbW91bnQSDTEwMDAwMDB1c3RhcnNIwJoMUOewBFqOAwoVL2Nvc21vcy50eC52MWJldGExLlR4EvQCCsUBCqQBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEoMBCixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhoRCgZ1c3RhcnMSBzEwMDAwMDASHHlBQ2lLVDVScFhSaU9UazRIWVptaWZ6eVZ4VGcSaApQCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAkv5JPXz/sLcxTj1NSrfNcBEFtLCKHe76OAO0yWH7KZCEgQKAgh/GAoSFAoOCgZ1c3RhcnMSBDUwMDAQwJoMGkBf+bviJ1lfEHnfoBrNldZ5isTZGKfrC9T1kf1TD5JiYXu7yqDpsiVots+YzR+udm1Uw2vO0HMhNz6BAvGPX7l7YhQyMDIzLTAxLTA5VDEzOjU0OjAwWmpfCgpjb2luX3NwZW50EjkKB3NwZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrGAESFgoGYW1vdW50Ego1MDAwdXN0YXJzGAFqYwoNY29pbl9yZWNlaXZlZBI6CghyZWNlaXZlchIsc3RhcnMxN3hwZnZha20yYW1nOTYyeWxzNmY4NHoza2VsbDhjNWx5OTVhcXYYARIWCgZhbW91bnQSCjUwMDB1c3RhcnMYAWqZAQoIdHJhbnNmZXISOwoJcmVjaXBpZW50EixzdGFyczE3eHBmdmFrbTJhbWc5NjJ5bHM2Zjg0ejNrZWxsOGM1bHk5NWFxdhgBEjgKBnNlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYARIWCgZhbW91bnQSCjUwMDB1c3RhcnMYAWpDCgdtZXNzYWdlEjgKBnNlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYAWpWCgJ0eBITCgNmZWUSCjUwMDB1c3RhcnMYARI7CglmZWVfcGF5ZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrGAFqQgoCdHgSPAoHYWNjX3NlcRIvc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsvMTAYAWptCgJ0eBJnCglzaWduYXR1cmUSWFgvbTc0aWRaWHhCNTM2QWF6WlhXZVlyRTJSaW42d3ZVOVpIOVV3K1NZbUY3dThxZzZiSWxhTGJQbU0wZnJuWnRWTU5yenRCeklUYytnUUx4ajErNWV3PT0YAWozCgdtZXNzYWdlEigKBmFjdGlvbhIcL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBgBamIKCmNvaW5fc3BlbnQSOQoHc3BlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYARIZCgZhbW91bnQSDTEwMDAwMDB1c3RhcnMYAWp6Cg1jb2luX3JlY2VpdmVkEk4KCHJlY2VpdmVyEkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGAESGQoGYW1vdW50Eg0xMDAwMDAwdXN0YXJzGAFqsAEKCHRyYW5zZmVyEk8KCXJlY2lwaWVudBJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhgBEjgKBnNlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYARIZCgZhbW91bnQSDTEwMDAwMDB1c3RhcnMYAWpDCgdtZXNzYWdlEjgKBnNlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYAWobCgdtZXNzYWdlEhAKBm1vZHVsZRIEYmFuaxgBEqMYCOKkhQMSQDZEQkVEODhBNDc2NTAxMEVCQzgzM0I5RjcxMUE1RjgxNzNFNzA5NjA5RjA3QTk5RjkwQzlBMUJFMUQyQjRCNjkqQDBBMUUwQTFDMkY2MzZGNzM2RDZGNzMyRTYyNjE2RTZCMkU3NjMxNjI2NTc0NjEzMTJFNEQ3MzY3NTM2NTZFNjQyiQZbeyJldmVudHMiOlt7InR5cGUiOiJjb2luX3JlY2VpdmVkIiwiYXR0cmlidXRlcyI6W3sia2V5IjoicmVjZWl2ZXIiLCJ2YWx1ZSI6InN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24ifSx7ImtleSI6ImFtb3VudCIsInZhbHVlIjoiMTAwMDAwMHVzdGFycyJ9XX0seyJ0eXBlIjoiY29pbl9zcGVudCIsImF0dHJpYnV0ZXMiOlt7ImtleSI6InNwZW5kZXIiLCJ2YWx1ZSI6InN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrIn0seyJrZXkiOiJhbW91bnQiLCJ2YWx1ZSI6IjEwMDAwMDB1c3RhcnMifV19LHsidHlwZSI6Im1lc3NhZ2UiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJhY3Rpb24iLCJ2YWx1ZSI6Ii9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQifSx7ImtleSI6InNlbmRlciIsInZhbHVlIjoic3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsifSx7ImtleSI6Im1vZHVsZSIsInZhbHVlIjoiYmFuayJ9XX0seyJ0eXBlIjoidHJhbnNmZXIiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJyZWNpcGllbnQiLCJ2YWx1ZSI6InN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24ifSx7ImtleSI6InNlbmRlciIsInZhbHVlIjoic3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsifSx7ImtleSI6ImFtb3VudCIsInZhbHVlIjoiMTAwMDAwMHVzdGFycyJ9XX1dfV06gAQadgoNY29pbl9yZWNlaXZlZBJMCghyZWNlaXZlchJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhIXCgZhbW91bnQSDTEwMDAwMDB1c3RhcnMaXgoKY29pbl9zcGVudBI3CgdzcGVuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axIXCgZhbW91bnQSDTEwMDAwMDB1c3RhcnMaeQoHbWVzc2FnZRImCgZhY3Rpb24SHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQSNgoGc2VuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axIOCgZtb2R1bGUSBGJhbmsaqgEKCHRyYW5zZmVyEk0KCXJlY2lwaWVudBJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhI2CgZzZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrEhcKBmFtb3VudBINMTAwMDAwMHVzdGFyc0jAmgxQ0rAEWo4DChUvY29zbW9zLnR4LnYxYmV0YTEuVHgS9AIKxQEKpAEKHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQSgwEKLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrEkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGhEKBnVzdGFycxIHMTAwMDAwMBIceUFFV1VzTGtrME9zdUl5NUxUQm9ZZ1dEWm5yNRJoClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiECS/kk9fP+wtzFOPU1Kt81wEQW0sIod7vo4A7TJYfspkISBAoCCH8YCxIUCg4KBnVzdGFycxIENTAwMBDAmgwaQI2EapZoriG/8W+CpmXvwDtwY/5tXXNMWC3ho6oJrqypG0GuV92w1eTtYoA0+DZ43TqaKFHK1FRB0HZuKALobvdiFDIwMjMtMDEtMDlUMTM6NTU6MDZaal8KCmNvaW5fc3BlbnQSOQoHc3BlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYARIWCgZhbW91bnQSCjUwMDB1c3RhcnMYAWpjCg1jb2luX3JlY2VpdmVkEjoKCHJlY2VpdmVyEixzdGFyczE3eHBmdmFrbTJhbWc5NjJ5bHM2Zjg0ejNrZWxsOGM1bHk5NWFxdhgBEhYKBmFtb3VudBIKNTAwMHVzdGFycxgBapkBCgh0cmFuc2ZlchI7CglyZWNpcGllbnQSLHN0YXJzMTd4cGZ2YWttMmFtZzk2MnlsczZmODR6M2tlbGw4YzVseTk1YXF2GAESOAoGc2VuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axgBEhYKBmFtb3VudBIKNTAwMHVzdGFycxgBakMKB21lc3NhZ2USOAoGc2VuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axgBalYKAnR4EhMKA2ZlZRIKNTAwMHVzdGFycxgBEjsKCWZlZV9wYXllchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYAWpCCgJ0eBI8CgdhY2Nfc2VxEi9zdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5ay8xMRgBam0KAnR4EmcKCXNpZ25hdHVyZRJYallScWxtaXVJYi94YjRLbVplL0FPM0JqL20xZGMweFlMZUdqcWdtdXJLa2JRYTVYM2JEVjVPMWlnRFQ0Tm5qZE9wb29VY3JVVkVIUWRtNG9BdWh1OXc9PRgBajMKB21lc3NhZ2USKAoGYWN0aW9uEhwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kGAFqYgoKY29pbl9zcGVudBI5CgdzcGVuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axgBEhkKBmFtb3VudBINMTAwMDAwMHVzdGFycxgBanoKDWNvaW5fcmVjZWl2ZWQSTgoIcmVjZWl2ZXISQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24YARIZCgZhbW91bnQSDTEwMDAwMDB1c3RhcnMYAWqwAQoIdHJhbnNmZXISTwoJcmVjaXBpZW50EkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGAESOAoGc2VuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axgBEhkKBmFtb3VudBINMTAwMDAwMHVzdGFycxgBakMKB21lc3NhZ2USOAoGc2VuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axgBahsKB21lc3NhZ2USEAoGbW9kdWxlEgRiYW5rGAESohgI0aqFAxJAOTdFQjA4QTYyMUE0MzI2MjJBQkE3RENFQkY3RjY4MzQ3QTQ3NjVERDM2NzFCMkMzNkY1OTc3OTJFMEU5RUQxNSpAMEExRTBBMUMyRjYzNkY3MzZENkY3MzJFNjI2MTZFNkIyRTc2MzE2MjY1NzQ2MTMxMkU0RDczNjc1MzY1NkU2NDKJBlt7ImV2ZW50cyI6W3sidHlwZSI6ImNvaW5fcmVjZWl2ZWQiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJyZWNlaXZlciIsInZhbHVlIjoic3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbiJ9LHsia2V5IjoiYW1vdW50IiwidmFsdWUiOiIxMDAwMDAwdXN0YXJzIn1dfSx7InR5cGUiOiJjb2luX3NwZW50IiwiYXR0cmlidXRlcyI6W3sia2V5Ijoic3BlbmRlciIsInZhbHVlIjoic3RhcnMxZjZnOWd1eWV5emd6amM5bDh3ZzR4bDV4MHJ2eGRkZXdkcWp2N3YifSx7ImtleSI6ImFtb3VudCIsInZhbHVlIjoiMTAwMDAwMHVzdGFycyJ9XX0seyJ0eXBlIjoibWVzc2FnZSIsImF0dHJpYnV0ZXMiOlt7ImtleSI6ImFjdGlvbiIsInZhbHVlIjoiL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZCJ9LHsia2V5Ijoic2VuZGVyIiwidmFsdWUiOiJzdGFyczFmNmc5Z3V5ZXl6Z3pqYzlsOHdnNHhsNXgwcnZ4ZGRld2RxanY3diJ9LHsia2V5IjoibW9kdWxlIiwidmFsdWUiOiJiYW5rIn1dfSx7InR5cGUiOiJ0cmFuc2ZlciIsImF0dHJpYnV0ZXMiOlt7ImtleSI6InJlY2lwaWVudCIsInZhbHVlIjoic3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbiJ9LHsia2V5Ijoic2VuZGVyIiwidmFsdWUiOiJzdGFyczFmNmc5Z3V5ZXl6Z3pqYzlsOHdnNHhsNXgwcnZ4ZGRld2RxanY3diJ9LHsia2V5IjoiYW1vdW50IiwidmFsdWUiOiIxMDAwMDAwdXN0YXJzIn1dfV19XTqABBp2Cg1jb2luX3JlY2VpdmVkEkwKCHJlY2VpdmVyEkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuEhcKBmFtb3VudBINMTAwMDAwMHVzdGFycxpeCgpjb2luX3NwZW50EjcKB3NwZW5kZXISLHN0YXJzMWY2ZzlndXlleXpnempjOWw4d2c0eGw1eDBydnhkZGV3ZHFqdjd2EhcKBmFtb3VudBINMTAwMDAwMHVzdGFycxp5CgdtZXNzYWdlEiYKBmFjdGlvbhIcL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBI2CgZzZW5kZXISLHN0YXJzMWY2ZzlndXlleXpnempjOWw4d2c0eGw1eDBydnhkZGV3ZHFqdjd2Eg4KBm1vZHVsZRIEYmFuaxqqAQoIdHJhbnNmZXISTQoJcmVjaXBpZW50EkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuEjYKBnNlbmRlchIsc3RhcnMxZjZnOWd1eWV5emd6amM5bDh3ZzR4bDV4MHJ2eGRkZXdkcWp2N3YSFwoGYW1vdW50Eg0xMDAwMDAwdXN0YXJzSMCaDFDJsARajgMKFS9jb3Ntb3MudHgudjFiZXRhMS5UeBL0AgrFAQqkAQocL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBKDAQosc3RhcnMxZjZnOWd1eWV5emd6amM5bDh3ZzR4bDV4MHJ2eGRkZXdkcWp2N3YSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24aEQoGdXN0YXJzEgcxMDAwMDAwEhx5QU1VdHN0OTgzQ05TMmpybk5mM3NwWU1jb3dZEmgKUApGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQOXhV8BudrTtyHQIe23pbeHpsSvAdCxMBIHBRKGnxM57BIECgIIfxgCEhQKDgoGdXN0YXJzEgQyMDAwEMCaDBpA8QvKmmxkxC2ZHQOZDzgM6oK7rtfNKUzt+91sZF4cJr4f3JmaxZigKuTXf20fU1Bi1mi2SkqlsH7CBtZ6DOmabGIUMjAyMy0wMS0wOVQxNTowODo0NFpqXwoKY29pbl9zcGVudBI5CgdzcGVuZGVyEixzdGFyczFmNmc5Z3V5ZXl6Z3pqYzlsOHdnNHhsNXgwcnZ4ZGRld2RxanY3dhgBEhYKBmFtb3VudBIKMjAwMHVzdGFycxgBamMKDWNvaW5fcmVjZWl2ZWQSOgoIcmVjZWl2ZXISLHN0YXJzMTd4cGZ2YWttMmFtZzk2MnlsczZmODR6M2tlbGw4YzVseTk1YXF2GAESFgoGYW1vdW50EgoyMDAwdXN0YXJzGAFqmQEKCHRyYW5zZmVyEjsKCXJlY2lwaWVudBIsc3RhcnMxN3hwZnZha20yYW1nOTYyeWxzNmY4NHoza2VsbDhjNWx5OTVhcXYYARI4CgZzZW5kZXISLHN0YXJzMWY2ZzlndXlleXpnempjOWw4d2c0eGw1eDBydnhkZGV3ZHFqdjd2GAESFgoGYW1vdW50EgoyMDAwdXN0YXJzGAFqQwoHbWVzc2FnZRI4CgZzZW5kZXISLHN0YXJzMWY2ZzlndXlleXpnempjOWw4d2c0eGw1eDBydnhkZGV3ZHFqdjd2GAFqVgoCdHgSEwoDZmVlEgoyMDAwdXN0YXJzGAESOwoJZmVlX3BheWVyEixzdGFyczFmNmc5Z3V5ZXl6Z3pqYzlsOHdnNHhsNXgwcnZ4ZGRld2RxanY3dhgBakEKAnR4EjsKB2FjY19zZXESLnN0YXJzMWY2ZzlndXlleXpnempjOWw4d2c0eGw1eDBydnhkZGV3ZHFqdjd2LzIYAWptCgJ0eBJnCglzaWduYXR1cmUSWDhRdkttbXhreEMyWkhRT1pEemdNNm9LN3J0Zk5LVXp0Kzkxc1pGNGNKcjRmM0ptYXhaaWdLdVRYZjIwZlUxQmkxbWkyU2txbHNIN0NCdFo2RE9tYWJBPT0YAWozCgdtZXNzYWdlEigKBmFjdGlvbhIcL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBgBamIKCmNvaW5fc3BlbnQSOQoHc3BlbmRlchIsc3RhcnMxZjZnOWd1eWV5emd6amM5bDh3ZzR4bDV4MHJ2eGRkZXdkcWp2N3YYARIZCgZhbW91bnQSDTEwMDAwMDB1c3RhcnMYAWp6Cg1jb2luX3JlY2VpdmVkEk4KCHJlY2VpdmVyEkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGAESGQoGYW1vdW50Eg0xMDAwMDAwdXN0YXJzGAFqsAEKCHRyYW5zZmVyEk8KCXJlY2lwaWVudBJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhgBEjgKBnNlbmRlchIsc3RhcnMxZjZnOWd1eWV5emd6amM5bDh3ZzR4bDV4MHJ2eGRkZXdkcWp2N3YYARIZCgZhbW91bnQSDTEwMDAwMDB1c3RhcnMYAWpDCgdtZXNzYWdlEjgKBnNlbmRlchIsc3RhcnMxZjZnOWd1eWV5emd6amM5bDh3ZzR4bDV4MHJ2eGRkZXdkcWp2N3YYAWobCgdtZXNzYWdlEhAKBm1vZHVsZRIEYmFuaxgBEvEYCJCuhQMSQDI1OTU0NDY4QzVGODA3RkY2NjM0OTk2N0U2QTgwQUZBNjBBOEUzNUQzNzhCMkZBMjY3RjlDQzc2NkEyM0QxNjgqQDBBMUUwQTFDMkY2MzZGNzM2RDZGNzMyRTYyNjE2RTZCMkU3NjMxNjI2NTc0NjEzMTJFNEQ3MzY3NTM2NTZFNjQyjwZbeyJldmVudHMiOlt7InR5cGUiOiJjb2luX3JlY2VpdmVkIiwiYXR0cmlidXRlcyI6W3sia2V5IjoicmVjZWl2ZXIiLCJ2YWx1ZSI6InN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24ifSx7ImtleSI6ImFtb3VudCIsInZhbHVlIjoiMTAwMDAwMDAwdXN0YXJzIn1dfSx7InR5cGUiOiJjb2luX3NwZW50IiwiYXR0cmlidXRlcyI6W3sia2V5Ijoic3BlbmRlciIsInZhbHVlIjoic3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgifSx7ImtleSI6ImFtb3VudCIsInZhbHVlIjoiMTAwMDAwMDAwdXN0YXJzIn1dfSx7InR5cGUiOiJtZXNzYWdlIiwiYXR0cmlidXRlcyI6W3sia2V5IjoiYWN0aW9uIiwidmFsdWUiOiIvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kIn0seyJrZXkiOiJzZW5kZXIiLCJ2YWx1ZSI6InN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4In0seyJrZXkiOiJtb2R1bGUiLCJ2YWx1ZSI6ImJhbmsifV19LHsidHlwZSI6InRyYW5zZmVyIiwiYXR0cmlidXRlcyI6W3sia2V5IjoicmVjaXBpZW50IiwidmFsdWUiOiJzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuIn0seyJrZXkiOiJzZW5kZXIiLCJ2YWx1ZSI6InN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4In0seyJrZXkiOiJhbW91bnQiLCJ2YWx1ZSI6IjEwMDAwMDAwMHVzdGFycyJ9XX1dfV06hgQaeAoNY29pbl9yZWNlaXZlZBJMCghyZWNlaXZlchJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhIZCgZhbW91bnQSDzEwMDAwMDAwMHVzdGFycxpgCgpjb2luX3NwZW50EjcKB3NwZW5kZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4EhkKBmFtb3VudBIPMTAwMDAwMDAwdXN0YXJzGnkKB21lc3NhZ2USJgoGYWN0aW9uEhwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEjYKBnNlbmRlchIsc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgSDgoGbW9kdWxlEgRiYW5rGqwBCgh0cmFuc2ZlchJNCglyZWNpcGllbnQSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24SNgoGc2VuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBIZCgZhbW91bnQSDzEwMDAwMDAwMHVzdGFyc0jAmgxQorcEWskDChUvY29zbW9zLnR4LnYxYmV0YTEuVHgSrwMK/wEKpgEKHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQShQEKLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4EkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGhMKBnVzdGFycxIJMTAwMDAwMDAwElRRZ0NpS1Q1UnBYUmlPVGs0SFlabWlmenlWeFRnUWdFT3prUjZSamNTTmtqOTZxOEU0WHgwVHZpblFnWlB6aWFHMFhGRmF2V1lqT1pYNjBySDZMUEESaQpRCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAiCrM8s85a4pqMYuDsoCo3y+Ce9GjQOz2xmAWCl2ePvbEgQKAgh/GJ8BEhQKDgoGdXN0YXJzEgQyMDAwEMCaDBpAd8tQIvmzrxH7MDn+cciFvNeKar+M82SzFKfF5Kj1oRY1b2NZxidrUAXnjkERxF33anqkZEk3CF8uucXbCi7JmGIUMjAyMy0wMS0wOVQxNTo1MjozOFpqXwoKY29pbl9zcGVudBI5CgdzcGVuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBgBEhYKBmFtb3VudBIKMjAwMHVzdGFycxgBamMKDWNvaW5fcmVjZWl2ZWQSOgoIcmVjZWl2ZXISLHN0YXJzMTd4cGZ2YWttMmFtZzk2MnlsczZmODR6M2tlbGw4YzVseTk1YXF2GAESFgoGYW1vdW50EgoyMDAwdXN0YXJzGAFqmQEKCHRyYW5zZmVyEjsKCXJlY2lwaWVudBIsc3RhcnMxN3hwZnZha20yYW1nOTYyeWxzNmY4NHoza2VsbDhjNWx5OTVhcXYYARI4CgZzZW5kZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4GAESFgoGYW1vdW50EgoyMDAwdXN0YXJzGAFqQwoHbWVzc2FnZRI4CgZzZW5kZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4GAFqVgoCdHgSEwoDZmVlEgoyMDAwdXN0YXJzGAESOwoJZmVlX3BheWVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBgBakMKAnR4Ej0KB2FjY19zZXESMHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4LzE1ORgBam0KAnR4EmcKCXNpZ25hdHVyZRJYZDh0UUl2bXpyeEg3TURuK2NjaUZ2TmVLYXIrTTgyU3pGS2ZGNUtqMW9SWTFiMk5aeGlkclVBWG5qa0VSeEYzM2FucWtaRWszQ0Y4dXVjWGJDaTdKbUE9PRgBajMKB21lc3NhZ2USKAoGYWN0aW9uEhwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kGAFqZAoKY29pbl9zcGVudBI5CgdzcGVuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBgBEhsKBmFtb3VudBIPMTAwMDAwMDAwdXN0YXJzGAFqfAoNY29pbl9yZWNlaXZlZBJOCghyZWNlaXZlchJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhgBEhsKBmFtb3VudBIPMTAwMDAwMDAwdXN0YXJzGAFqsgEKCHRyYW5zZmVyEk8KCXJlY2lwaWVudBJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhgBEjgKBnNlbmRlchIsc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgYARIbCgZhbW91bnQSDzEwMDAwMDAwMHVzdGFycxgBakMKB21lc3NhZ2USOAoGc2VuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBgBahsKB21lc3NhZ2USEAoGbW9kdWxlEgRiYW5rGAES8RgI95eGAxJARDhBRjVDRDlDOTEzM0Q2OUY5NUMxMjdFOEU5RUE1MTM5REY1ODhDNTY5RjJDMzczRURCODkzQTY1N0EwMTlGMipAMEExRTBBMUMyRjYzNkY3MzZENkY3MzJFNjI2MTZFNkIyRTc2MzE2MjY1NzQ2MTMxMkU0RDczNjc1MzY1NkU2NDKPBlt7ImV2ZW50cyI6W3sidHlwZSI6ImNvaW5fcmVjZWl2ZWQiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJyZWNlaXZlciIsInZhbHVlIjoic3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbiJ9LHsia2V5IjoiYW1vdW50IiwidmFsdWUiOiI1MDAwMDAwMDB1c3RhcnMifV19LHsidHlwZSI6ImNvaW5fc3BlbnQiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJzcGVuZGVyIiwidmFsdWUiOiJzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOCJ9LHsia2V5IjoiYW1vdW50IiwidmFsdWUiOiI1MDAwMDAwMDB1c3RhcnMifV19LHsidHlwZSI6Im1lc3NhZ2UiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJhY3Rpb24iLCJ2YWx1ZSI6Ii9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQifSx7ImtleSI6InNlbmRlciIsInZhbHVlIjoic3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgifSx7ImtleSI6Im1vZHVsZSIsInZhbHVlIjoiYmFuayJ9XX0seyJ0eXBlIjoidHJhbnNmZXIiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJyZWNpcGllbnQiLCJ2YWx1ZSI6InN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24ifSx7ImtleSI6InNlbmRlciIsInZhbHVlIjoic3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgifSx7ImtleSI6ImFtb3VudCIsInZhbHVlIjoiNTAwMDAwMDAwdXN0YXJzIn1dfV19XTqGBBp4Cg1jb2luX3JlY2VpdmVkEkwKCHJlY2VpdmVyEkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuEhkKBmFtb3VudBIPNTAwMDAwMDAwdXN0YXJzGmAKCmNvaW5fc3BlbnQSNwoHc3BlbmRlchIsc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgSGQoGYW1vdW50Eg81MDAwMDAwMDB1c3RhcnMaeQoHbWVzc2FnZRImCgZhY3Rpb24SHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQSNgoGc2VuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBIOCgZtb2R1bGUSBGJhbmsarAEKCHRyYW5zZmVyEk0KCXJlY2lwaWVudBJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhI2CgZzZW5kZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4EhkKBmFtb3VudBIPNTAwMDAwMDAwdXN0YXJzSMCaDFCitwRayQMKFS9jb3Ntb3MudHgudjFiZXRhMS5UeBKvAwr/AQqmAQocL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBKFAQosc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24aEwoGdXN0YXJzEgk1MDAwMDAwMDASVFFnTVV0c3Q5ODNDTlMyanJuTmYzc3BZTWNvd1lRZ0VPemtSNlJqY1NOa2o5NnE4RTRYeDBUdmluUWdDaUtUNVJwWFJpT1RrNEhZWm1pZnp5VnhUZxJpClEKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiECIKszyzzlrimoxi4OygKjfL4J70aNA7PbGYBYKXZ4+9sSBAoCCH8YoAESFAoOCgZ1c3RhcnMSBDIwMDAQwJoMGkBxZYEmJdTSaQPIEXXxJ/2BuGtNPW033NY3oYTFOoLaTixbILdhCbYHKhAi2H8PzMN8WAlG4+CYl8PYJH+l3jBaYhQyMDIzLTAxLTEwVDE0OjAxOjIwWmpfCgpjb2luX3NwZW50EjkKB3NwZW5kZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4GAESFgoGYW1vdW50EgoyMDAwdXN0YXJzGAFqYwoNY29pbl9yZWNlaXZlZBI6CghyZWNlaXZlchIsc3RhcnMxN3hwZnZha20yYW1nOTYyeWxzNmY4NHoza2VsbDhjNWx5OTVhcXYYARIWCgZhbW91bnQSCjIwMDB1c3RhcnMYAWqZAQoIdHJhbnNmZXISOwoJcmVjaXBpZW50EixzdGFyczE3eHBmdmFrbTJhbWc5NjJ5bHM2Zjg0ejNrZWxsOGM1bHk5NWFxdhgBEjgKBnNlbmRlchIsc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgYARIWCgZhbW91bnQSCjIwMDB1c3RhcnMYAWpDCgdtZXNzYWdlEjgKBnNlbmRlchIsc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgYAWpWCgJ0eBITCgNmZWUSCjIwMDB1c3RhcnMYARI7CglmZWVfcGF5ZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4GAFqQwoCdHgSPQoHYWNjX3NlcRIwc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgvMTYwGAFqbQoCdHgSZwoJc2lnbmF0dXJlElhjV1dCSmlYVTBta0R5QkYxOFNmOWdiaHJUVDF0Tjl6V042R0V4VHFDMms0c1d5QzNZUW0yQnlvUUl0aC9EOHpEZkZnSlJ1UGdtSmZEMkNSL3BkNHdXZz09GAFqMwoHbWVzc2FnZRIoCgZhY3Rpb24SHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQYAWpkCgpjb2luX3NwZW50EjkKB3NwZW5kZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4GAESGwoGYW1vdW50Eg81MDAwMDAwMDB1c3RhcnMYAWp8Cg1jb2luX3JlY2VpdmVkEk4KCHJlY2VpdmVyEkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGAESGwoGYW1vdW50Eg81MDAwMDAwMDB1c3RhcnMYAWqyAQoIdHJhbnNmZXISTwoJcmVjaXBpZW50EkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGAESOAoGc2VuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBgBEhsKBmFtb3VudBIPNTAwMDAwMDAwdXN0YXJzGAFqQwoHbWVzc2FnZRI4CgZzZW5kZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4GAFqGwoHbWVzc2FnZRIQCgZtb2R1bGUSBGJhbmsYARKBGQjyl4gDEkAyRjlGMDE2Q0FEMDcxRjZGN0ExQUU0OEIyN0FCMDY0QzhBQkYxMzU5MDEwRTdBRTg4NjJERkI4QjdEMzA2RjBBKkAwQTFFMEExQzJGNjM2RjczNkQ2RjczMkU2MjYxNkU2QjJFNzYzMTYyNjU3NDYxMzEyRTRENzM2NzUzNjU2RTY0MowGW3siZXZlbnRzIjpbeyJ0eXBlIjoiY29pbl9yZWNlaXZlZCIsImF0dHJpYnV0ZXMiOlt7ImtleSI6InJlY2VpdmVyIiwidmFsdWUiOiJzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuIn0seyJrZXkiOiJhbW91bnQiLCJ2YWx1ZSI6IjEwMDAwMDAwdXN0YXJzIn1dfSx7InR5cGUiOiJjb2luX3NwZW50IiwiYXR0cmlidXRlcyI6W3sia2V5Ijoic3BlbmRlciIsInZhbHVlIjoic3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsifSx7ImtleSI6ImFtb3VudCIsInZhbHVlIjoiMTAwMDAwMDB1c3RhcnMifV19LHsidHlwZSI6Im1lc3NhZ2UiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJhY3Rpb24iLCJ2YWx1ZSI6Ii9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQifSx7ImtleSI6InNlbmRlciIsInZhbHVlIjoic3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsifSx7ImtleSI6Im1vZHVsZSIsInZhbHVlIjoiYmFuayJ9XX0seyJ0eXBlIjoidHJhbnNmZXIiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJyZWNpcGllbnQiLCJ2YWx1ZSI6InN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24ifSx7ImtleSI6InNlbmRlciIsInZhbHVlIjoic3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsifSx7ImtleSI6ImFtb3VudCIsInZhbHVlIjoiMTAwMDAwMDB1c3RhcnMifV19XX1dOoMEGncKDWNvaW5fcmVjZWl2ZWQSTAoIcmVjZWl2ZXISQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24SGAoGYW1vdW50Eg4xMDAwMDAwMHVzdGFycxpfCgpjb2luX3NwZW50EjcKB3NwZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrEhgKBmFtb3VudBIOMTAwMDAwMDB1c3RhcnMaeQoHbWVzc2FnZRImCgZhY3Rpb24SHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQSNgoGc2VuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axIOCgZtb2R1bGUSBGJhbmsaqwEKCHRyYW5zZmVyEk0KCXJlY2lwaWVudBJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhI2CgZzZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrEhgKBmFtb3VudBIOMTAwMDAwMDB1c3RhcnNIwJoMULm3BFrjAwoVL2Nvc21vcy50eC52MWJldGExLlR4EskDCpoCCqUBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEoQBCixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhoSCgZ1c3RhcnMSCDEwMDAwMDAwEnBNZ0NpS1Q1UnBYUmlPVGs0SFlabWlmenlWeFRnTWdNVXRzdDk4M0NOUzJqcm5OZjNzcFlNY293WU1nRVdVc0xrazBPc3VJeTVMVEJvWWdXRFpucjVNZ3hSMm5DQytXL052bFFESHBvd2tDWll4MkRqEmgKUApGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQJL+ST18/7C3MU49TUq3zXARBbSwih3u+jgDtMlh+ymQhIECgIIfxgMEhQKDgoGdXN0YXJzEgQ1MDAwEMCaDBpAKKFU723c7Wiz20RxswGEsy93eIPhPNDP6waOz3WcsA854O8xGqG5tNGUJmaG1I2UyNekCdcN9gteHWd/zIB5iGIUMjAyMy0wMS0xMlQxOTo1MjoyM1pqXwoKY29pbl9zcGVudBI5CgdzcGVuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axgBEhYKBmFtb3VudBIKNTAwMHVzdGFycxgBamMKDWNvaW5fcmVjZWl2ZWQSOgoIcmVjZWl2ZXISLHN0YXJzMTd4cGZ2YWttMmFtZzk2MnlsczZmODR6M2tlbGw4YzVseTk1YXF2GAESFgoGYW1vdW50Ego1MDAwdXN0YXJzGAFqmQEKCHRyYW5zZmVyEjsKCXJlY2lwaWVudBIsc3RhcnMxN3hwZnZha20yYW1nOTYyeWxzNmY4NHoza2VsbDhjNWx5OTVhcXYYARI4CgZzZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrGAESFgoGYW1vdW50Ego1MDAwdXN0YXJzGAFqQwoHbWVzc2FnZRI4CgZzZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrGAFqVgoCdHgSEwoDZmVlEgo1MDAwdXN0YXJzGAESOwoJZmVlX3BheWVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axgBakIKAnR4EjwKB2FjY19zZXESL3N0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrLzEyGAFqbQoCdHgSZwoJc2lnbmF0dXJlElhLS0ZVNzIzYzdXaXoyMFJ4c3dHRXN5OTNlSVBoUE5EUDZ3YU96M1djc0E4NTRPOHhHcUc1dE5HVUptYUcxSTJVeU5la0NkY045Z3RlSFdkL3pJQjVpQT09GAFqMwoHbWVzc2FnZRIoCgZhY3Rpb24SHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQYAWpjCgpjb2luX3NwZW50EjkKB3NwZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrGAESGgoGYW1vdW50Eg4xMDAwMDAwMHVzdGFycxgBansKDWNvaW5fcmVjZWl2ZWQSTgoIcmVjZWl2ZXISQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24YARIaCgZhbW91bnQSDjEwMDAwMDAwdXN0YXJzGAFqsQEKCHRyYW5zZmVyEk8KCXJlY2lwaWVudBJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhgBEjgKBnNlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYARIaCgZhbW91bnQSDjEwMDAwMDAwdXN0YXJzGAFqQwoHbWVzc2FnZRI4CgZzZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrGAFqGwoHbWVzc2FnZRIQCgZtb2R1bGUSBGJhbmsYARLoGQiVmIgDEkA4NEFBNjBDNzU4NTI4NDk3QTYzNjkzMzE3MjdCMDY4Q0E2Q0VGMDZCRTdDRTZBODNEMTFCQ0M5RjNEQzRBNkE5KkAwQTFFMEExQzJGNjM2RjczNkQ2RjczMkU2MjYxNkU2QjJFNzYzMTYyNjU3NDYxMzEyRTRENzM2NzUzNjU2RTY0MokGW3siZXZlbnRzIjpbeyJ0eXBlIjoiY29pbl9yZWNlaXZlZCIsImF0dHJpYnV0ZXMiOlt7ImtleSI6InJlY2VpdmVyIiwidmFsdWUiOiJzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuIn0seyJrZXkiOiJhbW91bnQiLCJ2YWx1ZSI6IjUwMDAwMDB1c3RhcnMifV19LHsidHlwZSI6ImNvaW5fc3BlbnQiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJzcGVuZGVyIiwidmFsdWUiOiJzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5ayJ9LHsia2V5IjoiYW1vdW50IiwidmFsdWUiOiI1MDAwMDAwdXN0YXJzIn1dfSx7InR5cGUiOiJtZXNzYWdlIiwiYXR0cmlidXRlcyI6W3sia2V5IjoiYWN0aW9uIiwidmFsdWUiOiIvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kIn0seyJrZXkiOiJzZW5kZXIiLCJ2YWx1ZSI6InN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrIn0seyJrZXkiOiJtb2R1bGUiLCJ2YWx1ZSI6ImJhbmsifV19LHsidHlwZSI6InRyYW5zZmVyIiwiYXR0cmlidXRlcyI6W3sia2V5IjoicmVjaXBpZW50IiwidmFsdWUiOiJzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuIn0seyJrZXkiOiJzZW5kZXIiLCJ2YWx1ZSI6InN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrIn0seyJrZXkiOiJhbW91bnQiLCJ2YWx1ZSI6IjUwMDAwMDB1c3RhcnMifV19XX1dOoAEGnYKDWNvaW5fcmVjZWl2ZWQSTAoIcmVjZWl2ZXISQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24SFwoGYW1vdW50Eg01MDAwMDAwdXN0YXJzGl4KCmNvaW5fc3BlbnQSNwoHc3BlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsSFwoGYW1vdW50Eg01MDAwMDAwdXN0YXJzGnkKB21lc3NhZ2USJgoGYWN0aW9uEhwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEjYKBnNlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsSDgoGbW9kdWxlEgRiYW5rGqoBCgh0cmFuc2ZlchJNCglyZWNpcGllbnQSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24SNgoGc2VuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axIXCgZhbW91bnQSDTUwMDAwMDB1c3RhcnNIwJoMUJnABFrTBAoVL2Nvc21vcy50eC52MWJldGExLlR4ErkECooDCqQBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEoMBCixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhoRCgZ1c3RhcnMSBzUwMDAwMDAS4AFHUUNpS1Q1UnBYUmlPVGs0SFlabWlmenlWeFRnR1FNVXRzdDk4M0NOUzJqcm5OZjNzcFlNY293WUdReFIybkNDK1cvTnZsUURIcG93a0NaWXgyRGpHUXNNYi9RSzJiS1NtbXBrV0JORWsyMVB5dldOR1FrNWtjb2pWdUxNK0JSUi9ENDRBNkFOL3NKWUdRRzFaQS85MFYyeEZldWNmRkRXM2VHb0ViYmNHUUVPemtSNlJqY1NOa2o5NnE4RTRYeDBUdmluR1FFV1VzTGtrME9zdUl5NUxUQm9ZZ1dEWm5yNRJoClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiECS/kk9fP+wtzFOPU1Kt81wEQW0sIod7vo4A7TJYfspkISBAoCCH8YDRIUCg4KBnVzdGFycxIENTAwMBDAmgwaQHgpyuZv3WEpwaz44m5SpbcMeu4AgJfNmBJGxeb1lEZiCecGY9cDeQp+C9MhPm2yFiEOkHFcoEHIBEOyOMDHzQ1iFDIwMjMtMDEtMTJUMTk6NTU6NTNaal8KCmNvaW5fc3BlbnQSOQoHc3BlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYARIWCgZhbW91bnQSCjUwMDB1c3RhcnMYAWpjCg1jb2luX3JlY2VpdmVkEjoKCHJlY2VpdmVyEixzdGFyczE3eHBmdmFrbTJhbWc5NjJ5bHM2Zjg0ejNrZWxsOGM1bHk5NWFxdhgBEhYKBmFtb3VudBIKNTAwMHVzdGFycxgBapkBCgh0cmFuc2ZlchI7CglyZWNpcGllbnQSLHN0YXJzMTd4cGZ2YWttMmFtZzk2MnlsczZmODR6M2tlbGw4YzVseTk1YXF2GAESOAoGc2VuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axgBEhYKBmFtb3VudBIKNTAwMHVzdGFycxgBakMKB21lc3NhZ2USOAoGc2VuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axgBalYKAnR4EhMKA2ZlZRIKNTAwMHVzdGFycxgBEjsKCWZlZV9wYXllchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYAWpCCgJ0eBI8CgdhY2Nfc2VxEi9zdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5ay8xMxgBam0KAnR4EmcKCXNpZ25hdHVyZRJYZUNuSzVtL2RZU25CclBqaWJsS2x0d3g2N2dDQWw4MllFa2JGNXZXVVJtSUo1d1pqMXdONUNuNEwweUUrYmJJV0lRNlFjVnlnUWNnRVE3STR3TWZORFE9PRgBajMKB21lc3NhZ2USKAoGYWN0aW9uEhwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kGAFqYgoKY29pbl9zcGVudBI5CgdzcGVuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axgBEhkKBmFtb3VudBINNTAwMDAwMHVzdGFycxgBanoKDWNvaW5fcmVjZWl2ZWQSTgoIcmVjZWl2ZXISQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24YARIZCgZhbW91bnQSDTUwMDAwMDB1c3RhcnMYAWqwAQoIdHJhbnNmZXISTwoJcmVjaXBpZW50EkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGAESOAoGc2VuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axgBEhkKBmFtb3VudBINNTAwMDAwMHVzdGFycxgBakMKB21lc3NhZ2USOAoGc2VuZGVyEixzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5axgBahsKB21lc3NhZ2USEAoGbW9kdWxlEgRiYW5rGAES2xgI7pqIAxJAOTRBMTIzNzcwMzZENTU4MzEyQjhBNzczRjQ1QjNCRDAzRTc4QTIxQTcxOUJDMUJCMzlCRUNBRDFBNTQ4QjJGNypAMEExRTBBMUMyRjYzNkY3MzZENkY3MzJFNjI2MTZFNkIyRTc2MzE2MjY1NzQ2MTMxMkU0RDczNjc1MzY1NkU2NDKJBlt7ImV2ZW50cyI6W3sidHlwZSI6ImNvaW5fcmVjZWl2ZWQiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJyZWNlaXZlciIsInZhbHVlIjoic3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbiJ9LHsia2V5IjoiYW1vdW50IiwidmFsdWUiOiIxMDAwMDAwdXN0YXJzIn1dfSx7InR5cGUiOiJjb2luX3NwZW50IiwiYXR0cmlidXRlcyI6W3sia2V5Ijoic3BlbmRlciIsInZhbHVlIjoic3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsifSx7ImtleSI6ImFtb3VudCIsInZhbHVlIjoiMTAwMDAwMHVzdGFycyJ9XX0seyJ0eXBlIjoibWVzc2FnZSIsImF0dHJpYnV0ZXMiOlt7ImtleSI6ImFjdGlvbiIsInZhbHVlIjoiL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZCJ9LHsia2V5Ijoic2VuZGVyIiwidmFsdWUiOiJzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5ayJ9LHsia2V5IjoibW9kdWxlIiwidmFsdWUiOiJiYW5rIn1dfSx7InR5cGUiOiJ0cmFuc2ZlciIsImF0dHJpYnV0ZXMiOlt7ImtleSI6InJlY2lwaWVudCIsInZhbHVlIjoic3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbiJ9LHsia2V5Ijoic2VuZGVyIiwidmFsdWUiOiJzdGFyczE5NTRxOWFwYXdyNmtnOGV6NHVreDhqeXVheGFrejd5ZTIydHR5ayJ9LHsia2V5IjoiYW1vdW50IiwidmFsdWUiOiIxMDAwMDAwdXN0YXJzIn1dfV19XTqABBp2Cg1jb2luX3JlY2VpdmVkEkwKCHJlY2VpdmVyEkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuEhcKBmFtb3VudBINMTAwMDAwMHVzdGFycxpeCgpjb2luX3NwZW50EjcKB3NwZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrEhcKBmFtb3VudBINMTAwMDAwMHVzdGFycxp5CgdtZXNzYWdlEiYKBmFjdGlvbhIcL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBI2CgZzZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrEg4KBm1vZHVsZRIEYmFuaxqqAQoIdHJhbnNmZXISTQoJcmVjaXBpZW50EkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuEjYKBnNlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsSFwoGYW1vdW50Eg0xMDAwMDAwdXN0YXJzSMCaDFDPtARaxgMKFS9jb3Ntb3MudHgudjFiZXRhMS5UeBKsAwr9AQqkAQocL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBKDAQosc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24aEQoGdXN0YXJzEgcxMDAwMDAwElRRZ0VXVXNMa2swT3N1SXk1TFRCb1lnV0RabnI1UWdNVXRzdDk4M0NOUzJqcm5OZjNzcFlNY293WVFnc01iL1FLMmJLU21tcGtXQk5FazIxUHl2V04SaApQCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAkv5JPXz/sLcxTj1NSrfNcBEFtLCKHe76OAO0yWH7KZCEgQKAgh/GA4SFAoOCgZ1c3RhcnMSBDUwMDAQwJoMGkAsos4fzLKDDRFRjmtLUOVUYM0pZsCtFuqq8SIelLgkwUvkWe+znB6cLI2zHUvGpQZHmjWM9fFJsda7IMuqyYHPYhQyMDIzLTAxLTEyVDIwOjI5OjU2WmpfCgpjb2luX3NwZW50EjkKB3NwZW5kZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrGAESFgoGYW1vdW50Ego1MDAwdXN0YXJzGAFqYwoNY29pbl9yZWNlaXZlZBI6CghyZWNlaXZlchIsc3RhcnMxN3hwZnZha20yYW1nOTYyeWxzNmY4NHoza2VsbDhjNWx5OTVhcXYYARIWCgZhbW91bnQSCjUwMDB1c3RhcnMYAWqZAQoIdHJhbnNmZXISOwoJcmVjaXBpZW50EixzdGFyczE3eHBmdmFrbTJhbWc5NjJ5bHM2Zjg0ejNrZWxsOGM1bHk5NWFxdhgBEjgKBnNlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYARIWCgZhbW91bnQSCjUwMDB1c3RhcnMYAWpDCgdtZXNzYWdlEjgKBnNlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYAWpWCgJ0eBITCgNmZWUSCjUwMDB1c3RhcnMYARI7CglmZWVfcGF5ZXISLHN0YXJzMTk1NHE5YXBhd3I2a2c4ZXo0dWt4OGp5dWF4YWt6N3llMjJ0dHlrGAFqQgoCdHgSPAoHYWNjX3NlcRIvc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsvMTQYAWptCgJ0eBJnCglzaWduYXR1cmUSWExLTE9IOHl5Z3cwUlVZNXJTMURsVkdETktXYkFyUmJxcXZFaUhwUzRKTUZMNUZudnM1d2VuQ3lOc3gxTHhxVUdSNW8xalBYeFNiSFd1eURMcXNtQnp3PT0YAWozCgdtZXNzYWdlEigKBmFjdGlvbhIcL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBgBamIKCmNvaW5fc3BlbnQSOQoHc3BlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYARIZCgZhbW91bnQSDTEwMDAwMDB1c3RhcnMYAWp6Cg1jb2luX3JlY2VpdmVkEk4KCHJlY2VpdmVyEkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGAESGQoGYW1vdW50Eg0xMDAwMDAwdXN0YXJzGAFqsAEKCHRyYW5zZmVyEk8KCXJlY2lwaWVudBJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhgBEjgKBnNlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYARIZCgZhbW91bnQSDTEwMDAwMDB1c3RhcnMYAWpDCgdtZXNzYWdlEjgKBnNlbmRlchIsc3RhcnMxOTU0cTlhcGF3cjZrZzhlejR1a3g4anl1YXhha3o3eWUyMnR0eWsYAWobCgdtZXNzYWdlEhAKBm1vZHVsZRIEYmFuaxgBEvsYCIObiQMSQDU3NkYwOUYxREMxNURDQUIwRDE4OTA5MTM4Q0Y3Q0IzOTU2RDg0RUE3OTYzMThDNTg3MTQyNjk1MjA0NjgzQjQqQDBBMUUwQTFDMkY2MzZGNzM2RDZGNzMyRTYyNjE2RTZCMkU3NjMxNjI2NTc0NjEzMTJFNEQ3MzY3NTM2NTZFNjQykgZbeyJldmVudHMiOlt7InR5cGUiOiJjb2luX3JlY2VpdmVkIiwiYXR0cmlidXRlcyI6W3sia2V5IjoicmVjZWl2ZXIiLCJ2YWx1ZSI6InN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24ifSx7ImtleSI6ImFtb3VudCIsInZhbHVlIjoiMTAwMDAwMDAwMHVzdGFycyJ9XX0seyJ0eXBlIjoiY29pbl9zcGVudCIsImF0dHJpYnV0ZXMiOlt7ImtleSI6InNwZW5kZXIiLCJ2YWx1ZSI6InN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4In0seyJrZXkiOiJhbW91bnQiLCJ2YWx1ZSI6IjEwMDAwMDAwMDB1c3RhcnMifV19LHsidHlwZSI6Im1lc3NhZ2UiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJhY3Rpb24iLCJ2YWx1ZSI6Ii9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQifSx7ImtleSI6InNlbmRlciIsInZhbHVlIjoic3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgifSx7ImtleSI6Im1vZHVsZSIsInZhbHVlIjoiYmFuayJ9XX0seyJ0eXBlIjoidHJhbnNmZXIiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJyZWNpcGllbnQiLCJ2YWx1ZSI6InN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24ifSx7ImtleSI6InNlbmRlciIsInZhbHVlIjoic3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgifSx7ImtleSI6ImFtb3VudCIsInZhbHVlIjoiMTAwMDAwMDAwMHVzdGFycyJ9XX1dfV06iQQaeQoNY29pbl9yZWNlaXZlZBJMCghyZWNlaXZlchJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhIaCgZhbW91bnQSEDEwMDAwMDAwMDB1c3RhcnMaYQoKY29pbl9zcGVudBI3CgdzcGVuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBIaCgZhbW91bnQSEDEwMDAwMDAwMDB1c3RhcnMaeQoHbWVzc2FnZRImCgZhY3Rpb24SHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQSNgoGc2VuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBIOCgZtb2R1bGUSBGJhbmsarQEKCHRyYW5zZmVyEk0KCXJlY2lwaWVudBJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhI2CgZzZW5kZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4EhoKBmFtb3VudBIQMTAwMDAwMDAwMHVzdGFyc0jAmgxQyrcEWsoDChUvY29zbW9zLnR4LnYxYmV0YTEuVHgSsAMKgAIKpwEKHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQShgEKLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4EkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGhQKBnVzdGFycxIKMTAwMDAwMDAwMBJUUWdFT3prUjZSamNTTmtqOTZxOEU0WHgwVHZpblFnWlB6aWFHMFhGRmF2V1lqT1pYNjBySDZMUEFRZ3NNYi9RSzJiS1NtbXBrV0JORWsyMVB5dldOEmkKUQpGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQIgqzPLPOWuKajGLg7KAqN8vgnvRo0Ds9sZgFgpdnj72xIECgIIfxihARIUCg4KBnVzdGFycxIEMjAwMBDAmgwaQBkHGPi3T29grcYzY1oSqM2yVUchY0iF8SenI1Az85t4I/5ASTib7O5Z2N3rOVZa0MoWaG4VreWvi34KDvnycRliFDIwMjMtMDEtMTNUMjM6MTg6MjZaal8KCmNvaW5fc3BlbnQSOQoHc3BlbmRlchIsc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgYARIWCgZhbW91bnQSCjIwMDB1c3RhcnMYAWpjCg1jb2luX3JlY2VpdmVkEjoKCHJlY2VpdmVyEixzdGFyczE3eHBmdmFrbTJhbWc5NjJ5bHM2Zjg0ejNrZWxsOGM1bHk5NWFxdhgBEhYKBmFtb3VudBIKMjAwMHVzdGFycxgBapkBCgh0cmFuc2ZlchI7CglyZWNpcGllbnQSLHN0YXJzMTd4cGZ2YWttMmFtZzk2MnlsczZmODR6M2tlbGw4YzVseTk1YXF2GAESOAoGc2VuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBgBEhYKBmFtb3VudBIKMjAwMHVzdGFycxgBakMKB21lc3NhZ2USOAoGc2VuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBgBalYKAnR4EhMKA2ZlZRIKMjAwMHVzdGFycxgBEjsKCWZlZV9wYXllchIsc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgYAWpDCgJ0eBI9CgdhY2Nfc2VxEjBzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOC8xNjEYAWptCgJ0eBJnCglzaWduYXR1cmUSWEdRY1krTGRQYjJDdHhqTmpXaEtvemJKVlJ5RmpTSVh4SjZjalVEUHptM2dqL2tCSk9KdnM3bG5ZM2VzNVZsclF5aFpvYmhXdDVhK0xmZ29PK2ZKeEdRPT0YAWozCgdtZXNzYWdlEigKBmFjdGlvbhIcL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBgBamUKCmNvaW5fc3BlbnQSOQoHc3BlbmRlchIsc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgYARIcCgZhbW91bnQSEDEwMDAwMDAwMDB1c3RhcnMYAWp9Cg1jb2luX3JlY2VpdmVkEk4KCHJlY2VpdmVyEkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGAESHAoGYW1vdW50EhAxMDAwMDAwMDAwdXN0YXJzGAFqswEKCHRyYW5zZmVyEk8KCXJlY2lwaWVudBJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhgBEjgKBnNlbmRlchIsc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgYARIcCgZhbW91bnQSEDEwMDAwMDAwMDB1c3RhcnMYAWpDCgdtZXNzYWdlEjgKBnNlbmRlchIsc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgYAWobCgdtZXNzYWdlEhAKBm1vZHVsZRIEYmFuaxgBEo0ZCOGciQMSQDU5OTFDMjkyNkYxRDQ3MkNFMjYwRDA0MkM1QUZBRUQ0RjE0Mzk5NDUxNDY1RUEzRjVFNjRCMzBFNDM3NEY5QzQqQDBBMUUwQTFDMkY2MzZGNzM2RDZGNzMyRTYyNjE2RTZCMkU3NjMxNjI2NTc0NjEzMTJFNEQ3MzY3NTM2NTZFNjQyjwZbeyJldmVudHMiOlt7InR5cGUiOiJjb2luX3JlY2VpdmVkIiwiYXR0cmlidXRlcyI6W3sia2V5IjoicmVjZWl2ZXIiLCJ2YWx1ZSI6InN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24ifSx7ImtleSI6ImFtb3VudCIsInZhbHVlIjoiMjUwMDAwMDAwdXN0YXJzIn1dfSx7InR5cGUiOiJjb2luX3NwZW50IiwiYXR0cmlidXRlcyI6W3sia2V5Ijoic3BlbmRlciIsInZhbHVlIjoic3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgifSx7ImtleSI6ImFtb3VudCIsInZhbHVlIjoiMjUwMDAwMDAwdXN0YXJzIn1dfSx7InR5cGUiOiJtZXNzYWdlIiwiYXR0cmlidXRlcyI6W3sia2V5IjoiYWN0aW9uIiwidmFsdWUiOiIvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kIn0seyJrZXkiOiJzZW5kZXIiLCJ2YWx1ZSI6InN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4In0seyJrZXkiOiJtb2R1bGUiLCJ2YWx1ZSI6ImJhbmsifV19LHsidHlwZSI6InRyYW5zZmVyIiwiYXR0cmlidXRlcyI6W3sia2V5IjoicmVjaXBpZW50IiwidmFsdWUiOiJzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuIn0seyJrZXkiOiJzZW5kZXIiLCJ2YWx1ZSI6InN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4In0seyJrZXkiOiJhbW91bnQiLCJ2YWx1ZSI6IjI1MDAwMDAwMHVzdGFycyJ9XX1dfV06hgQaeAoNY29pbl9yZWNlaXZlZBJMCghyZWNlaXZlchJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhIZCgZhbW91bnQSDzI1MDAwMDAwMHVzdGFycxpgCgpjb2luX3NwZW50EjcKB3NwZW5kZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4EhkKBmFtb3VudBIPMjUwMDAwMDAwdXN0YXJzGnkKB21lc3NhZ2USJgoGYWN0aW9uEhwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEjYKBnNlbmRlchIsc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgSDgoGbW9kdWxlEgRiYW5rGqwBCgh0cmFuc2ZlchJNCglyZWNpcGllbnQSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24SNgoGc2VuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBIZCgZhbW91bnQSDzI1MDAwMDAwMHVzdGFyc0jAmgxQurkEWuUDChUvY29zbW9zLnR4LnYxYmV0YTEuVHgSywMKmwIKpgEKHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQShQEKLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4EkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGhMKBnVzdGFycxIJMjUwMDAwMDAwEnBNZ01VdHN0OTgzQ05TMmpybk5mM3NwWU1jb3dZTWdaUHppYUcwWEZGYXZXWWpPWlg2MHJINkxQQU1neFIybkNDK1cvTnZsUURIcG93a0NaWXgyRGpNaGczdnByQ0pHWEpJNnVhbmQvZWNrbFVJbWZiEmkKUQpGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQIgqzPLPOWuKajGLg7KAqN8vgnvRo0Ds9sZgFgpdnj72xIECgIIfxiiARIUCg4KBnVzdGFycxIEMjAwMBDAmgwaQHPe88pHy3OpvZ9ucF2bhC5kJg2rL+9qOULufCc0W0Z+LrnpI0BhLlSvW9jRyl83X7SITggcBGYkJBnGzmWsncBiFDIwMjMtMDEtMTNUMjM6NDA6MDdaal8KCmNvaW5fc3BlbnQSOQoHc3BlbmRlchIsc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgYARIWCgZhbW91bnQSCjIwMDB1c3RhcnMYAWpjCg1jb2luX3JlY2VpdmVkEjoKCHJlY2VpdmVyEixzdGFyczE3eHBmdmFrbTJhbWc5NjJ5bHM2Zjg0ejNrZWxsOGM1bHk5NWFxdhgBEhYKBmFtb3VudBIKMjAwMHVzdGFycxgBapkBCgh0cmFuc2ZlchI7CglyZWNpcGllbnQSLHN0YXJzMTd4cGZ2YWttMmFtZzk2MnlsczZmODR6M2tlbGw4YzVseTk1YXF2GAESOAoGc2VuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBgBEhYKBmFtb3VudBIKMjAwMHVzdGFycxgBakMKB21lc3NhZ2USOAoGc2VuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBgBalYKAnR4EhMKA2ZlZRIKMjAwMHVzdGFycxgBEjsKCWZlZV9wYXllchIsc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgYAWpDCgJ0eBI9CgdhY2Nfc2VxEjBzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOC8xNjIYAWptCgJ0eBJnCglzaWduYXR1cmUSWGM5N3p5a2ZMYzZtOW4yNXdYWnVFTG1RbURhc3Y3Mm81UXU1OEp6UmJSbjR1dWVralFHRXVWSzliMk5IS1h6ZGZ0SWhPQ0J3RVppUWtHY2JPWmF5ZHdBPT0YAWozCgdtZXNzYWdlEigKBmFjdGlvbhIcL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBgBamQKCmNvaW5fc3BlbnQSOQoHc3BlbmRlchIsc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgYARIbCgZhbW91bnQSDzI1MDAwMDAwMHVzdGFycxgBanwKDWNvaW5fcmVjZWl2ZWQSTgoIcmVjZWl2ZXISQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24YARIbCgZhbW91bnQSDzI1MDAwMDAwMHVzdGFycxgBarIBCgh0cmFuc2ZlchJPCglyZWNpcGllbnQSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24YARI4CgZzZW5kZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4GAESGwoGYW1vdW50Eg8yNTAwMDAwMDB1c3RhcnMYAWpDCgdtZXNzYWdlEjgKBnNlbmRlchIsc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgYAWobCgdtZXNzYWdlEhAKBm1vZHVsZRIEYmFuaxgBEsEYCPOciQMSQDQyODQ3NkU0QzU2ODE1MjM4NEY2MzlBRUYzMjI5OTFFNDMwMTlDMUM3MEY5N0QwQzEzNjJFMUM5OURCNEQ1QzIqQDBBMUUwQTFDMkY2MzZGNzM2RDZGNzMyRTYyNjE2RTZCMkU3NjMxNjI2NTc0NjEzMTJFNEQ3MzY3NTM2NTZFNjQyiQZbeyJldmVudHMiOlt7InR5cGUiOiJjb2luX3JlY2VpdmVkIiwiYXR0cmlidXRlcyI6W3sia2V5IjoicmVjZWl2ZXIiLCJ2YWx1ZSI6InN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24ifSx7ImtleSI6ImFtb3VudCIsInZhbHVlIjoiMTAwMDAwMHVzdGFycyJ9XX0seyJ0eXBlIjoiY29pbl9zcGVudCIsImF0dHJpYnV0ZXMiOlt7ImtleSI6InNwZW5kZXIiLCJ2YWx1ZSI6InN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4In0seyJrZXkiOiJhbW91bnQiLCJ2YWx1ZSI6IjEwMDAwMDB1c3RhcnMifV19LHsidHlwZSI6Im1lc3NhZ2UiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJhY3Rpb24iLCJ2YWx1ZSI6Ii9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQifSx7ImtleSI6InNlbmRlciIsInZhbHVlIjoic3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgifSx7ImtleSI6Im1vZHVsZSIsInZhbHVlIjoiYmFuayJ9XX0seyJ0eXBlIjoidHJhbnNmZXIiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJyZWNpcGllbnQiLCJ2YWx1ZSI6InN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24ifSx7ImtleSI6InNlbmRlciIsInZhbHVlIjoic3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgifSx7ImtleSI6ImFtb3VudCIsInZhbHVlIjoiMTAwMDAwMHVzdGFycyJ9XX1dfV06gAQadgoNY29pbl9yZWNlaXZlZBJMCghyZWNlaXZlchJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhIXCgZhbW91bnQSDTEwMDAwMDB1c3RhcnMaXgoKY29pbl9zcGVudBI3CgdzcGVuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBIXCgZhbW91bnQSDTEwMDAwMDB1c3RhcnMaeQoHbWVzc2FnZRImCgZhY3Rpb24SHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQSNgoGc2VuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBIOCgZtb2R1bGUSBGJhbmsaqgEKCHRyYW5zZmVyEk0KCXJlY2lwaWVudBJAc3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbhI2CgZzZW5kZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4EhcKBmFtb3VudBINMTAwMDAwMHVzdGFyc0jAmgxQhbUEWqsDChUvY29zbW9zLnR4LnYxYmV0YTEuVHgSkQMK4QEKpAEKHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQSgwEKLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4EkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGhEKBnVzdGFycxIHMTAwMDAwMBI4WkFFT3prUjZSamNTTmtqOTZxOEU0WHgwVHZpblpBTVV0c3Q5ODNDTlMyanJuTmYzc3BZTWNvd1kSaQpRCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAiCrM8s85a4pqMYuDsoCo3y+Ce9GjQOz2xmAWCl2ePvbEgQKAgh/GKMBEhQKDgoGdXN0YXJzEgQyMDAwEMCaDBpAPMjgvAdlqK86BI5/6+pGlmagpnxg/BQ2yk+NxB0ff/poFiLKrAGZSf8TQ2I7aEI0C6DVPl3nSIhNpQs74Idk7WIUMjAyMy0wMS0xM1QyMzo0MTo1MlpqXwoKY29pbl9zcGVudBI5CgdzcGVuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBgBEhYKBmFtb3VudBIKMjAwMHVzdGFycxgBamMKDWNvaW5fcmVjZWl2ZWQSOgoIcmVjZWl2ZXISLHN0YXJzMTd4cGZ2YWttMmFtZzk2MnlsczZmODR6M2tlbGw4YzVseTk1YXF2GAESFgoGYW1vdW50EgoyMDAwdXN0YXJzGAFqmQEKCHRyYW5zZmVyEjsKCXJlY2lwaWVudBIsc3RhcnMxN3hwZnZha20yYW1nOTYyeWxzNmY4NHoza2VsbDhjNWx5OTVhcXYYARI4CgZzZW5kZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4GAESFgoGYW1vdW50EgoyMDAwdXN0YXJzGAFqQwoHbWVzc2FnZRI4CgZzZW5kZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4GAFqVgoCdHgSEwoDZmVlEgoyMDAwdXN0YXJzGAESOwoJZmVlX3BheWVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBgBakMKAnR4Ej0KB2FjY19zZXESMHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4LzE2MxgBam0KAnR4EmcKCXNpZ25hdHVyZRJYUE1qZ3ZBZGxxSzg2Qkk1LzYrcEdsbWFncG54Zy9CUTJ5aytOeEIwZmYvcG9GaUxLckFHWlNmOFRRMkk3YUVJMEM2RFZQbDNuU0loTnBRczc0SWRrN1E9PRgBajMKB21lc3NhZ2USKAoGYWN0aW9uEhwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kGAFqYgoKY29pbl9zcGVudBI5CgdzcGVuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBgBEhkKBmFtb3VudBINMTAwMDAwMHVzdGFycxgBanoKDWNvaW5fcmVjZWl2ZWQSTgoIcmVjZWl2ZXISQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24YARIZCgZhbW91bnQSDTEwMDAwMDB1c3RhcnMYAWqwAQoIdHJhbnNmZXISTwoJcmVjaXBpZW50EkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGAESOAoGc2VuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBgBEhkKBmFtb3VudBINMTAwMDAwMHVzdGFycxgBakMKB21lc3NhZ2USOAoGc2VuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBgBahsKB21lc3NhZ2USEAoGbW9kdWxlEgRiYW5rGAEStBkIj8mLAxJANUJFRTY0QzNEMEEzMDYwMjU0Q0VDQkE1NUMxM0FFOEY1NUNBRkJCNTZENDZBRTQzQTBFRjQ2OEQzODQzMDc4NCpAMEExRTBBMUMyRjYzNkY3MzZENkY3MzJFNjI2MTZFNkIyRTc2MzE2MjY1NzQ2MTMxMkU0RDczNjc1MzY1NkU2NDKSBlt7ImV2ZW50cyI6W3sidHlwZSI6ImNvaW5fcmVjZWl2ZWQiLCJhdHRyaWJ1dGVzIjpbeyJrZXkiOiJyZWNlaXZlciIsInZhbHVlIjoic3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbiJ9LHsia2V5IjoiYW1vdW50IiwidmFsdWUiOiIxMjcxMDAwMDAwdXN0YXJzIn1dfSx7InR5cGUiOiJjb2luX3NwZW50IiwiYXR0cmlidXRlcyI6W3sia2V5Ijoic3BlbmRlciIsInZhbHVlIjoic3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgifSx7ImtleSI6ImFtb3VudCIsInZhbHVlIjoiMTI3MTAwMDAwMHVzdGFycyJ9XX0seyJ0eXBlIjoibWVzc2FnZSIsImF0dHJpYnV0ZXMiOlt7ImtleSI6ImFjdGlvbiIsInZhbHVlIjoiL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZCJ9LHsia2V5Ijoic2VuZGVyIiwidmFsdWUiOiJzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOCJ9LHsia2V5IjoibW9kdWxlIiwidmFsdWUiOiJiYW5rIn1dfSx7InR5cGUiOiJ0cmFuc2ZlciIsImF0dHJpYnV0ZXMiOlt7ImtleSI6InJlY2lwaWVudCIsInZhbHVlIjoic3RhcnMxNms5cWtxNTdrcHdjbnphd2Q4dTB1dGw2dTJ6aDVtcjJkejdxcDN3eTd5d3g5bTd4a2RhcW5xNW1zbiJ9LHsia2V5Ijoic2VuZGVyIiwidmFsdWUiOiJzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOCJ9LHsia2V5IjoiYW1vdW50IiwidmFsdWUiOiIxMjcxMDAwMDAwdXN0YXJzIn1dfV19XTqJBBp5Cg1jb2luX3JlY2VpdmVkEkwKCHJlY2VpdmVyEkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuEhoKBmFtb3VudBIQMTI3MTAwMDAwMHVzdGFycxphCgpjb2luX3NwZW50EjcKB3NwZW5kZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4EhoKBmFtb3VudBIQMTI3MTAwMDAwMHVzdGFycxp5CgdtZXNzYWdlEiYKBmFjdGlvbhIcL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBI2CgZzZW5kZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4Eg4KBm1vZHVsZRIEYmFuaxqtAQoIdHJhbnNmZXISTQoJcmVjaXBpZW50EkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuEjYKBnNlbmRlchIsc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgSGgoGYW1vdW50EhAxMjcxMDAwMDAwdXN0YXJzSMCaDFDmuwRagwQKFS9jb3Ntb3MudHgudjFiZXRhMS5UeBLpAwq5AgqnAQocL2Nvc21vcy5iYW5rLnYxYmV0YTEuTXNnU2VuZBKGAQosc3RhcnMxNngwM3djcDM3a3g1ZThlaGNranh2d2NnazlqMGNxbmg4cWx1ZTgSQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24aFAoGdXN0YXJzEgoxMjcxMDAwMDAwEowBS0FNVXRzdDk4M0NOUzJqcm5OZjNzcFlNY293WUtBc01iL1FLMmJLU21tcGtXQk5FazIxUHl2V05LQk1UVEFHVklTZHEwTmdJQTBsQjlROG9jTXIyS0JnM3ZwckNKR1hKSTZ1YW5kL2Vja2xVSW1mYktBeFIybkNDK1cvTnZsUURIcG93a0NaWXgyRGoSaQpRCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAiCrM8s85a4pqMYuDsoCo3y+Ce9GjQOz2xmAWCl2ePvbEgQKAgh/GKQBEhQKDgoGdXN0YXJzEgQyMDAwEMCaDBpAjDvqcqrhBlHZvm8Zhhjv+nbpTftU1/XbUuMSH2wlHdV7jsaw+bHApnrIiSv6Dp29WevXGPdyoQznyo8TPv25MGIUMjAyMy0wMS0xNlQxNDoyOTowNFpqXwoKY29pbl9zcGVudBI5CgdzcGVuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBgBEhYKBmFtb3VudBIKMjAwMHVzdGFycxgBamMKDWNvaW5fcmVjZWl2ZWQSOgoIcmVjZWl2ZXISLHN0YXJzMTd4cGZ2YWttMmFtZzk2MnlsczZmODR6M2tlbGw4YzVseTk1YXF2GAESFgoGYW1vdW50EgoyMDAwdXN0YXJzGAFqmQEKCHRyYW5zZmVyEjsKCXJlY2lwaWVudBIsc3RhcnMxN3hwZnZha20yYW1nOTYyeWxzNmY4NHoza2VsbDhjNWx5OTVhcXYYARI4CgZzZW5kZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4GAESFgoGYW1vdW50EgoyMDAwdXN0YXJzGAFqQwoHbWVzc2FnZRI4CgZzZW5kZXISLHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4GAFqVgoCdHgSEwoDZmVlEgoyMDAwdXN0YXJzGAESOwoJZmVlX3BheWVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBgBakMKAnR4Ej0KB2FjY19zZXESMHN0YXJzMTZ4MDN3Y3AzN2t4NWU4ZWhja2p4dndjZ2s5ajBjcW5oOHFsdWU4LzE2NBgBam0KAnR4EmcKCXNpZ25hdHVyZRJYakR2cWNxcmhCbEhadm04WmhoanYrbmJwVGZ0VTEvWGJVdU1TSDJ3bEhkVjdqc2F3K2JIQXBucklpU3Y2RHAyOVdldlhHUGR5b1F6bnlvOFRQdjI1TUE9PRgBajMKB21lc3NhZ2USKAoGYWN0aW9uEhwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kGAFqZQoKY29pbl9zcGVudBI5CgdzcGVuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBgBEhwKBmFtb3VudBIQMTI3MTAwMDAwMHVzdGFycxgBan0KDWNvaW5fcmVjZWl2ZWQSTgoIcmVjZWl2ZXISQHN0YXJzMTZrOXFrcTU3a3B3Y256YXdkOHUwdXRsNnUyemg1bXIyZHo3cXAzd3k3eXd4OW03eGtkYXFucTVtc24YARIcCgZhbW91bnQSEDEyNzEwMDAwMDB1c3RhcnMYAWqzAQoIdHJhbnNmZXISTwoJcmVjaXBpZW50EkBzdGFyczE2azlxa3E1N2twd2NuemF3ZDh1MHV0bDZ1MnpoNW1yMmR6N3FwM3d5N3l3eDltN3hrZGFxbnE1bXNuGAESOAoGc2VuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBgBEhwKBmFtb3VudBIQMTI3MTAwMDAwMHVzdGFycxgBakMKB21lc3NhZ2USOAoGc2VuZGVyEixzdGFyczE2eDAzd2NwMzdreDVlOGVoY2tqeHZ3Y2drOWowY3FuaDhxbHVlOBgBahsKB21lc3NhZ2USEAoGbW9kdWxlEgRiYW5rGAEaAhAd"

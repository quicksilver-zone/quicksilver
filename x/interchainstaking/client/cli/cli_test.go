package cli_test

import (
	"fmt"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"
	tmcli "github.com/tendermint/tendermint/libs/cli"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/quicksilver-zone/quicksilver/app"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/client/cli"
	"github.com/quicksilver-zone/quicksilver/x/interchainstaking/types"
)

const (
	testDelegator = "cosmos1r2dthxctqzhwg299e7aaeqwfkgcc9hg8n9scjg"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
	zones   []types.Zone
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	// Use baseURL to make API HTTP requests or use val.RPCClient to make direct
	// Tendermint RPC calls. (from testutil/network godocs)

	s.cfg = app.DefaultConfig()

	updateGenesisConfigState := func(moduleName string, moduleState proto.Message) {
		buf, err := s.cfg.Codec.MarshalJSON(moduleState)
		s.Require().NoError(err)
		s.cfg.GenesisState[moduleName] = buf
	}

	zone := types.Zone{
		ConnectionId:                 "connection-0",
		ChainId:                      "cosmoshub-4",
		DepositAddress:               nil,
		WithdrawalAddress:            nil,
		PerformanceAddress:           nil,
		DelegationAddress:            nil,
		AccountPrefix:                "cosmos",
		LocalDenom:                   "uqatom",
		BaseDenom:                    "uatom",
		RedemptionRate:               sdk.ZeroDec(),
		LastRedemptionRate:           sdk.ZeroDec(),
		Validators:                   nil,
		AggregateIntent:              types.ValidatorIntents{},
		MultiSend:                    false,
		LiquidityModule:              false,
		WithdrawalWaitgroup:          0,
		IbcNextValidatorsHash:        nil,
		ValidatorSelectionAllocation: 0,
		HoldingsAllocation:           0,
		LastEpochHeight:              0,
		Tvl:                          sdk.ZeroDec(),
		UnbondingPeriod:              0,
		MessagesPerTx:                0,
		Is_118:                       true,
		DustThreshold:                math.NewInt(1000000),
	}

	zone.Validators = append(zone.Validators,
		&types.Validator{ValoperAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), DelegatorShares: sdk.NewDec(2000), Score: sdk.ZeroDec(), ValidatorBondShares: sdk.ZeroDec(), LiquidShares: sdk.ZeroDec()},
		&types.Validator{ValoperAddress: "cosmosvaloper156gqf9837u7d4c4678yt3rl4ls9c5vuursrrzf", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), DelegatorShares: sdk.NewDec(2000), Score: sdk.ZeroDec(), ValidatorBondShares: sdk.ZeroDec(), LiquidShares: sdk.ZeroDec()},
		&types.Validator{ValoperAddress: "cosmosvaloper14lultfckehtszvzw4ehu0apvsr77afvyju5zzy", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), DelegatorShares: sdk.NewDec(2000), Score: sdk.ZeroDec(), ValidatorBondShares: sdk.ZeroDec(), LiquidShares: sdk.ZeroDec()},
		&types.Validator{ValoperAddress: "cosmosvaloper1a3yjj7d3qnx4spgvjcwjq9cw9snrrrhu5h6jll", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), DelegatorShares: sdk.NewDec(2000), Score: sdk.ZeroDec(), ValidatorBondShares: sdk.ZeroDec(), LiquidShares: sdk.ZeroDec()},
		&types.Validator{ValoperAddress: "cosmosvaloper1z8zjv3lntpwxua0rtpvgrcwl0nm0tltgpgs6l7", CommissionRate: sdk.MustNewDecFromStr("0.2"), VotingPower: sdk.NewInt(2000), DelegatorShares: sdk.NewDec(2000), Score: sdk.ZeroDec(), ValidatorBondShares: sdk.ZeroDec(), LiquidShares: sdk.ZeroDec()},
	)

	delegationIntents := []types.DelegatorIntentsForZone{
		{
			ChainId: zone.ChainId,
			DelegationIntent: []*types.DelegatorIntent{
				{
					Delegator: testDelegator,
					Intents: types.ValidatorIntents{
						{
							ValoperAddress: zone.Validators[0].ValoperAddress,
							Weight:         sdk.NewDec(1),
						},
						{
							ValoperAddress: zone.Validators[1].ValoperAddress,
							Weight:         sdk.NewDec(1),
						},
					},
				},
			},
		},
	}
	// setup basic genesis state
	newGenesis := types.DefaultGenesis()
	newGenesis.Zones = []types.Zone{zone}
	newGenesis.DelegatorIntents = delegationIntents
	updateGenesisConfigState(types.ModuleName, newGenesis)
	s.zones = []types.Zone{zone}

	net, err := network.New(s.T(), s.T().TempDir(), s.cfg)
	s.Require().NoError(err)
	s.network = net

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestGetCmdZones() {
	val := s.network.Validators[0]

	tests := []struct {
		name      string
		args      []string
		expectErr bool
		respType  *types.QueryZonesResponse
		expected  *types.QueryZonesResponse
	}{
		{
			"valid",
			[]string{},
			false,
			&types.QueryZonesResponse{},
			&types.QueryZonesResponse{
				Zones: s.zones,
			},
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			clientCtx := val.ClientCtx

			runFlags := []string{
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			}

			args := tt.args
			args = append(args, runFlags...)
			cmd := cli.GetCmdZones()

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
			if tt.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), tt.respType), out.String())
				for i, zone := range tt.respType.Zones {
					s.Require().True(s.ZonesEqual(tt.expected.Zones[i], zone))
				}
			}
		})
	}
}

func (s *IntegrationTestSuite) ZonesEqual(zoneA, zoneB types.Zone) bool {
	s.T().Helper()

	s.Require().Equal(zoneA.BaseDenom, zoneB.BaseDenom)
	s.Require().Equal(zoneA.AggregateIntent, zoneB.AggregateIntent)
	s.Require().Equal(zoneA.ChainId, zoneB.ChainId)
	s.Require().Equal(zoneA.AccountPrefix, zoneB.AccountPrefix)
	s.Require().Equal(zoneA.DelegationAddress, zoneB.DelegationAddress)
	s.Require().Equal(zoneA.DepositAddress, zoneB.DepositAddress)
	s.Require().Equal(zoneA.ConnectionId, zoneB.ConnectionId)
	s.Require().Equal(zoneA.Decimals, zoneB.Decimals)
	s.Require().Equal(zoneA.DepositsEnabled, zoneB.DepositsEnabled)
	s.Require().Equal(zoneA.PerformanceAddress, zoneA.PerformanceAddress)
	s.Require().Equal(zoneA.WithdrawalWaitgroup, zoneB.WithdrawalWaitgroup)
	s.Require().Equal(zoneA.WithdrawalAddress, zoneB.WithdrawalAddress)
	s.Require().Equal(zoneA.HoldingsAllocation, zoneB.HoldingsAllocation)
	s.Require().Equal(zoneA.RedemptionRate, zoneB.RedemptionRate)
	s.Require().Equal(zoneA.ReturnToSender, zoneB.ReturnToSender)
	s.Require().Equal(zoneA.LastEpochHeight, zoneB.LastEpochHeight)
	s.Require().Equal(zoneA.LiquidityModule, zoneB.LiquidityModule)
	s.Require().Equal(zoneA.LocalDenom, zoneB.LocalDenom)
	s.Require().Equal(zoneA.DustThreshold, zoneB.DustThreshold)
	for i := range zoneA.Validators {
		s.Require().Equal(zoneA.Validators[i], zoneB.Validators[i])
	}
	s.Require().Equal(zoneA.Is_118, zoneB.Is_118)

	return true
}

func (s *IntegrationTestSuite) TestGetDelegatorIntentCmd() {
	val := s.network.Validators[0]
	tests := []struct {
		name      string
		args      []string
		expectErr bool
		respType  proto.Message
		expected  proto.Message
	}{
		{
			"no args",
			[]string{},
			true,
			&types.QueryDelegatorIntentResponse{},
			&types.QueryDelegatorIntentResponse{},
		},
		{
			"empty args",
			[]string{"", ""},
			true,
			&types.QueryDelegatorIntentResponse{},
			&types.QueryDelegatorIntentResponse{},
		},
		{
			"invalid chainID",
			[]string{"boguschainid", ""},
			true,
			&types.QueryDelegatorIntentResponse{},
			&types.QueryDelegatorIntentResponse{},
		},
		{
			"valid",
			[]string{s.zones[0].ChainId, testDelegator},
			false,
			&types.QueryDelegatorIntentResponse{},
			&types.QueryDelegatorIntentResponse{
				Intent: &types.DelegatorIntent{
					Delegator: testDelegator,
					Intents: types.ValidatorIntents{
						{
							ValoperAddress: s.zones[0].Validators[0].ValoperAddress,
							Weight:         sdk.NewDec(1),
						},
						{
							ValoperAddress: s.zones[0].Validators[1].ValoperAddress,
							Weight:         sdk.NewDec(1),
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			clientCtx := val.ClientCtx
			runFlags := []string{
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			}

			args := tt.args
			args = append(args, runFlags...)
			cmd := cli.GetDelegatorIntentCmd()

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
			if tt.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), tt.respType), out.String())
				s.Require().Equal(tt.expected.String(), tt.respType.String())
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetDepositAccountCmd() {
	val := s.network.Validators[0]

	tests := []struct {
		name      string
		args      []string
		expectErr bool
		respType  proto.Message
		expected  proto.Message
	}{
		{
			"no args",
			[]string{},
			true,
			&types.QueryDepositAccountForChainResponse{},
			&types.QueryDepositAccountForChainResponse{},
		},
		{
			"empty args",
			[]string{""},
			true,
			&types.QueryDepositAccountForChainResponse{},
			&types.QueryDepositAccountForChainResponse{},
		},
		{
			"invalid chainID",
			[]string{"boguschainid"},
			true,
			&types.QueryDepositAccountForChainResponse{},
			&types.QueryDepositAccountForChainResponse{},
		},
		/* {
			"valid",
			[]string{s.cfg.ChainID},
			false,
			&types.QueryDepositAccountForChainResponse{},
			&types.QueryDepositAccountForChainResponse{},
		}, */
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			clientCtx := val.ClientCtx

			cmd := cli.GetDepositAccountCmd()

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tt.args)
			if tt.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), tt.respType), out.String())
				s.Require().Equal(tt.expected, tt.respType)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetSignalIntentTxCmd() {
	val := s.network.Validators[0]

	tests := []struct {
		name         string
		args         []string
		expectErr    bool
		expectedCode uint32
		respType     proto.Message
	}{
		{
			"no args",
			[]string{},
			true,
			0,
			&sdk.TxResponse{},
		},
		{
			"empty args",
			[]string{"", ""},
			true,
			0,
			&sdk.TxResponse{},
		},
		{
			"invalid delegation_intent arg_format",
			[]string{s.network.Config.ChainID, "intents"},
			true,
			0,
			&sdk.TxResponse{},
		},
		{
			"invalid delegation_intent content",
			[]string{s.network.Config.ChainID, "0.0cosmos1valoper1xxxxxxxxx,0.1cosmosvaloper1yyyyyyyyy,1.1cosmosvaloper1zzzzzzzzz"},
			true,
			0,
			&sdk.TxResponse{},
		},
		{
			"invalid delegation_intent valoperAddress",
			[]string{s.network.Config.ChainID, "0.3A12UEL5L,0.3a12uel5l,0.1notok1ezyfcl"},
			true,
			0,
			&sdk.TxResponse{},
		},
		{
			"invalid delegation_intent weightOverrun",
			[]string{s.network.Config.ChainID, "0.4A12UEL5L,0.3a12uel5l,0.4abcdef1qpzry9x8gf2tvdw0s3jn54khce6mua7lmqqqxw"},
			true,
			0,
			&sdk.TxResponse{},
		},
		{
			"invalid delegation_intent weightUnderrun",
			[]string{s.network.Config.ChainID, "0.3A12UEL5L,0.3a12uel5l,0.3abcdef1qpzry9x8gf2tvdw0s3jn54khce6mua7lmqqqxw"},
			true,
			0,
			&sdk.TxResponse{},
		},
		{
			"invalid delegation_intent maxWeightOverrun",
			[]string{s.network.Config.ChainID, "0.3A12UEL5L,0.3a12uel5l,1.4abcdef1qpzry9x8gf2tvdw0s3jn54khce6mua7lmqqqxw"},
			true,
			0,
			&sdk.TxResponse{},
		},
		{
			"invalid from_address",
			[]string{s.network.Config.ChainID, "0.3A12UEL5L,0.3a12uel5l,0.4abcdef1qpzry9x8gf2tvdw0s3jn54khce6mua7lmqqqxw", "--from", "bogusAddress"},
			true,
			0,
			&sdk.TxResponse{},
		},
		{
			"invalid chainID",
			[]string{
				"boguschainid",
				"0.3A12UEL5L,0.3a12uel5l,0.4abcdef1qpzry9x8gf2tvdw0s3jn54khce6mua7lmqqqxw",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
			},
			true,
			0,
			&sdk.TxResponse{},
		},
		/* {
			"valid",
			[]string{
				s.network.Config.ChainID,
				"0.3A12UEL5L,0.3a12uel5l,0.4abcdef1qpzry9x8gf2tvdw0s3jn54khce6mua7lmqqqxw",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
			},
			false,
			0,
			&sdk.TxResponse{},
		}, */
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			clientCtx := val.ClientCtx

			runFlags := []string{
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=true", flags.FlagDryRun),
			}

			args := tt.args
			args = append(args, runFlags...)
			cmd := cli.GetSignalIntentTxCmd()

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
			if tt.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), tt.respType), out.String())
				txResp := tt.respType.(*sdk.TxResponse)
				s.Require().Equal(tt.expectedCode, txResp.Code, fmt.Sprintf("%v\n", txResp))
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetWithdrawalRecordsCmd() {
	val := s.network.Validators[0]

	tests := []struct {
		name      string
		args      []string
		expectErr bool
		respType  proto.Message
		expected  proto.Message
	}{
		{
			"valid",
			[]string{},
			false,
			&types.QueryWithdrawalRecordsResponse{},
			&types.QueryWithdrawalRecordsResponse{
				Withdrawals: []types.WithdrawalRecord{},
			},
		},
		{
			"invalid",
			[]string{"bogus"},
			true,
			&types.QueryWithdrawalRecordsResponse{},
			&types.QueryWithdrawalRecordsResponse{},
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			clientCtx := val.ClientCtx

			runFlags := []string{
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			}
			tt.args = append(tt.args, runFlags...)
			cmd := cli.GetWithdrawalRecordsCmd()

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tt.args)
			if tt.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), tt.respType), out.String())
				s.Require().Equal(tt.expected, tt.respType)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetUserWithdrawalRecords() {
	val := s.network.Validators[0]

	tests := []struct {
		name      string
		args      []string
		expectErr bool
		respType  proto.Message
		expected  proto.Message
	}{
		{
			"no args",
			[]string{},
			true,
			&types.QueryWithdrawalRecordsResponse{},
			&types.QueryWithdrawalRecordsResponse{},
		},
		{
			"empty args",
			[]string{""},
			true,
			&types.QueryWithdrawalRecordsResponse{},
			&types.QueryWithdrawalRecordsResponse{},
		},
		{
			"invalid chainID",
			[]string{"boguschainid"},
			true,
			&types.QueryWithdrawalRecordsResponse{},
			&types.QueryWithdrawalRecordsResponse{},
		},
		{
			"valid",
			[]string{
				"cosmos1r2dthxctqzhwg299e7aaeqwfkgcc9hg8n9scjg",
			},
			false,
			&types.QueryWithdrawalRecordsResponse{},
			&types.QueryWithdrawalRecordsResponse{
				Withdrawals: []types.WithdrawalRecord{},
			},
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			clientCtx := val.ClientCtx

			runFlags := []string{
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			}
			tt.args = append(tt.args, runFlags...)
			cmd := cli.GetUserWithdrawalRecordsCmd()

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tt.args)
			if tt.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), tt.respType), out.String())
				s.Require().Equal(tt.expected, tt.respType)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetZoneWithdrawalRecordsCmd() {
	val := s.network.Validators[0]

	tests := []struct {
		name      string
		args      []string
		expectErr bool
		respType  proto.Message
		expected  proto.Message
	}{
		{
			"no args",
			[]string{},
			true,
			&types.QueryWithdrawalRecordsResponse{},
			&types.QueryWithdrawalRecordsResponse{},
		},
		{
			"empty args",
			[]string{""},
			true,
			&types.QueryWithdrawalRecordsResponse{},
			&types.QueryWithdrawalRecordsResponse{},
		},
		{
			"invalid chainID",
			[]string{"boguschainid"},
			true,
			&types.QueryWithdrawalRecordsResponse{},
			&types.QueryWithdrawalRecordsResponse{},
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			clientCtx := val.ClientCtx

			runFlags := []string{
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=true", flags.FlagDryRun),
			}

			tt.args = append(tt.args, runFlags...)

			cmd := cli.GetZoneWithdrawalRecordsCmd()

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tt.args)
			if tt.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), tt.respType), out.String())
				s.Require().Equal(tt.expected, tt.respType)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetZoneRedelegationRecordsCmd() {
	val := s.network.Validators[0]

	tests := []struct {
		name      string
		args      []string
		expectErr bool
		respType  proto.Message
		expected  proto.Message
	}{
		{
			"no args",
			[]string{},
			true,
			&types.QueryRedelegationRecordsResponse{},
			&types.QueryRedelegationRecordsResponse{},
		},
		{
			"empty args",
			[]string{""},
			false,
			&types.QueryRedelegationRecordsResponse{
				Pagination: &query.PageResponse{},
			},
			&types.QueryRedelegationRecordsResponse{
				Redelegations: []types.RedelegationRecord{},
				Pagination: &query.PageResponse{
					Total: 0,
				},
			},
		},
		{
			"invalid chainID",
			[]string{"boguschainid"},
			false,
			&types.QueryRedelegationRecordsResponse{
				Pagination: &query.PageResponse{},
			},
			&types.QueryRedelegationRecordsResponse{
				Redelegations: []types.RedelegationRecord{},
				Pagination: &query.PageResponse{
					Total: 0,
				},
			},
		},
		{
			"valid",
			[]string{s.cfg.ChainID},
			false,
			&types.QueryRedelegationRecordsResponse{
				Pagination: &query.PageResponse{},
			},
			&types.QueryRedelegationRecordsResponse{
				Redelegations: []types.RedelegationRecord{},
				Pagination: &query.PageResponse{
					Total: 0,
				},
			},
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			clientCtx := val.ClientCtx

			runFlags := []string{
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			}
			tt.args = append(tt.args, runFlags...)

			cmd := cli.GetZoneRedelegationRecordsCmd()

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tt.args)
			if tt.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), tt.respType), out.String())
				s.Require().Equal(tt.expected, tt.respType)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetZoneValidatorsCmd() {
	val := s.network.Validators[0]

	tests := []struct {
		name      string
		args      []string
		expectErr bool
		respType  proto.Message
		expected  proto.Message
	}{
		{
			"no args",
			[]string{},
			true,
			&types.QueryZoneValidatorsResponse{},
			&types.QueryZoneValidatorsResponse{},
		},
		{
			"empty args",
			[]string{""},
			true,
			&types.QueryZoneValidatorsResponse{},
			&types.QueryZoneValidatorsResponse{},
		},
		{
			"invalid chainID",
			[]string{"boguschainid"},
			false,
			&types.QueryZoneValidatorsResponse{},
			&types.QueryZoneValidatorsResponse{
				Validators: []types.Validator{},
				Pagination: &query.PageResponse{
					Total: 0,
				},
			},
		},
		// Not work because of the validator is not set in the store(SetValidatorForZone)
		// {
		// 	"valid",
		// 	[]string{s.zones[0].ChainId},
		// 	false,
		// 	&types.QueryZoneValidatorsResponse{},
		// 	&types.QueryZoneValidatorsResponse{
		// 		// Convert from []*types.Validator to []types.Validator
		// 		Validators: func() []types.Validator {
		// 			var validators []types.Validator
		// 			for _, val := range s.zones[0].Validators {
		// 				validators = append(validators, *val)
		// 			}
		// 			return validators
		// 		}(),
		// 		Pagination: &query.PageResponse{
		// 			Total: 5,
		// 		},
		// 	},
		// },
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			clientCtx := val.ClientCtx

			runFlags := []string{
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			}
			tt.args = append(tt.args, runFlags...)

			cmd := cli.GetZoneValidatorsCmd()
			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tt.args)
			if tt.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), tt.respType), out.String())
				s.Require().Equal(tt.expected, tt.respType)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetTxStatusCmd() {
	val := s.network.Validators[0]

	tests := []struct {
		name      string
		args      []string
		expectErr bool
		respType  proto.Message
	}{
		{
			"no args",
			[]string{},
			true,
			&sdk.TxResponse{},
		},
		{
			"empty args",
			[]string{""},
			true,
			&sdk.TxResponse{},
		},
		{
			"only one arg",
			[]string{"boguschainid"},
			true,
			&sdk.TxResponse{},
		},
		// TODO: Add a test for a valid tx hash
		// {
		// 	"valid",
		// 	[]string{"chainid", "0"},
		// 	false,
		// 	&sdk.TxResponse{},
		// },
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			clientCtx := val.ClientCtx

			runFlags := []string{
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			}
			tt.args = append(tt.args, runFlags...)

			cmd := cli.GetTxStatusCmd()
			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tt.args)
			if tt.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), tt.respType), out.String())
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetReceiptsCmd() {
	val := s.network.Validators[0]

	tests := []struct {
		name        string
		args        []string
		expectErr   bool
		respType    proto.Message
		expectedLen int
	}{
		{
			"no args",
			[]string{},
			true,
			&types.QueryReceiptsResponse{},
			0,
		},
		{
			"empty args",
			[]string{""},
			true,
			&types.QueryReceiptsResponse{},
			0,
		},
		{
			"invalid chainID",
			[]string{"boguschainid"},
			true,
			&types.QueryReceiptsResponse{},
			0,
		},
		{
			"valid",
			[]string{s.zones[0].ChainId},
			false,
			&types.QueryReceiptsResponse{},
			1,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			clientCtx := val.ClientCtx

			runFlags := []string{
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			}
			tt.args = append(tt.args, runFlags...)

			cmd := cli.GetReceiptsCmd()
			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tt.args)
			if tt.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), tt.respType), out.String())
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetUnbondingRecordsCmd() {
	val := s.network.Validators[0]

	tests := []struct {
		name      string
		args      []string
		expectErr bool
		respType  proto.Message
	}{
		{
			"no args",
			[]string{},
			true,
			&types.QueryUnbondingRecordsResponse{},
		},
		{
			"valid",
			[]string{s.zones[0].ChainId},
			false,
			&types.QueryUnbondingRecordsResponse{},
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			clientCtx := val.ClientCtx

			runFlags := []string{
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			}
			tt.args = append(tt.args, runFlags...)

			cmd := cli.GetUnbondingRecordsCmd()
			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tt.args)
			if tt.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), tt.respType), out.String())
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetZoneCmd() {
	val := s.network.Validators[0]

	tests := []struct {
		name      string
		args      []string
		expectErr bool
		respType  proto.Message
	}{
		{
			"no args",
			[]string{},
			true,
			&types.QueryZoneResponse{},
		},
		{
			"empty args",
			[]string{""},
			true,
			&types.QueryZoneResponse{},
		},
		{
			"invalid chainID",
			[]string{"boguschainid"},
			true,
			&types.QueryZoneResponse{},
		},
		{
			"valid",
			[]string{s.zones[0].ChainId},
			false,
			&types.QueryZoneResponse{
				Zone:  s.zones[0],
				Stats: &types.Statistics{},
			},
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			clientCtx := val.ClientCtx

			runFlags := []string{
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			}
			tt.args = append(tt.args, runFlags...)

			cmd := cli.GetZoneCmd()
			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tt.args)
			if tt.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), tt.respType), out.String())

			}
		})
	}
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

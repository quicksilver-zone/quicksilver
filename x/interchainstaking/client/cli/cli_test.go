package cli_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/gogo/protobuf/proto"
	"github.com/ingenuity-build/quicksilver/app"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/client/cli"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
	"github.com/stretchr/testify/suite"
	tmcli "github.com/tendermint/tendermint/libs/cli"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	// Use baseURL to make API HTTP requests or use val.RPCClient to make direct
	// Tendermint RPC calls. (from testutil/network godocs)

	s.cfg = app.DefaultConfig()
	var err error
	s.network, err = network.New(s.T(), s.T().TempDir(), s.cfg)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestGetCmdZonesInfos() {
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
			&types.QueryRegisteredZonesInfoResponse{},
			&types.QueryRegisteredZonesInfoResponse{
				Pagination: &query.PageResponse{},
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		s.Run(tt.name, func() {
			clientCtx := val.ClientCtx

			flags := []string{
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			}
			args := append(tt.args, flags...)

			cmd := cli.GetCmdZonesInfos()

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
			if tt.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), tt.respType), out.String())
				s.Require().Equal(tt.expected.String(), tt.respType.String(), out.String())
			}
		})
	}
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
			"invalid chainid",
			[]string{"boguschainid", ""},
			true,
			&types.QueryDelegatorIntentResponse{},
			&types.QueryDelegatorIntentResponse{},
		},
		/*{
			"valid",
			[]string{s.cfg.ChainID, ""},
			false,
			&types.QueryDelegatorIntentResponse{},
			&types.QueryDelegatorIntentResponse{},
		},*/
	}
	for _, tt := range tests {
		tt := tt

		s.Run(tt.name, func() {
			clientCtx := val.ClientCtx

			flags := []string{
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			}
			args := append(tt.args, flags...)

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

/*func (s *IntegrationTestSuite) TestGetDepositAccountCmd() {
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
			"invalid chainid",
			[]string{"boguschainid"},
			true,
			&types.QueryDepositAccountForChainResponse{},
			&types.QueryDepositAccountForChainResponse{},
		},
		{
			"valid",
			[]string{s.cfg.ChainID},
			false,
			&types.QueryDepositAccountForChainResponse{},
			&types.QueryDepositAccountForChainResponse{},
		},
	}
	for _, tt := range tests {
		tt := tt

		s.Run(tt.name, func() {
			clientCtx := val.ClientCtx

			flags := []string{
				//fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			}
			args := append(tt.args, flags...)

			cmd := cli.GetDepositAccountCmd()

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
			if tt.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), tt.respType), out.String())
				s.Require().Equal(tt.expected, tt.respType)
			}
		})
	}
}*/

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
			"invalid chain_id",
			[]string{
				"boguschainid",
				"0.3A12UEL5L,0.3a12uel5l,0.4abcdef1qpzry9x8gf2tvdw0s3jn54khce6mua7lmqqqxw",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
			},
			true,
			0,
			&sdk.TxResponse{},
		},
		/*{
			"valid",
			[]string{
				s.network.Config.ChainID,
				"0.3A12UEL5L,0.3a12uel5l,0.4abcdef1qpzry9x8gf2tvdw0s3jn54khce6mua7lmqqqxw",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
			},
			false,
			0,
			&sdk.TxResponse{},
		},*/
	}
	for _, tt := range tests {
		tt := tt

		s.Run(tt.name, func() {
			clientCtx := val.ClientCtx

			flags := []string{
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=true", flags.FlagDryRun),
			}
			args := append(tt.args, flags...)

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

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

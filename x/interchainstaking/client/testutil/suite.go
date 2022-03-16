package testutil

import (
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/client/cli"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func NewIntegrationTestSuite(cfg network.Config) *IntegrationTestSuite {
	return &IntegrationTestSuite{cfg: cfg}
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	s.network = network.New(s.T(), s.cfg)

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
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
			[]string{"chainid", "intents"},
			true,
			0,
			&sdk.TxResponse{},
		},
		{
			"invalid delegation_intent content",
			[]string{"chainid", "0.0cosmos1valoper1xxxxxxxxx,0.1cosmosvaloper1yyyyyyyyy,1.1cosmosvaloper1zzzzzzzzz"},
			true,
			0,
			&sdk.TxResponse{},
		},
		{
			"invalid delegation_intent valoperAddress",
			[]string{"chainid", "0.3A12UEL5L,0.3a12uel5l,0.1notok1ezyfcl"},
			true,
			0,
			&sdk.TxResponse{},
		},
		{
			"invalid delegation_intent weightOverrun",
			[]string{"chainid", "0.4A12UEL5L,0.3a12uel5l,0.4abcdef1qpzry9x8gf2tvdw0s3jn54khce6mua7lmqqqxw"},
			true,
			0,
			&sdk.TxResponse{},
		},
		{
			"invalid delegation_intent weightUnderrun",
			[]string{"chainid", "0.3A12UEL5L,0.3a12uel5l,0.3abcdef1qpzry9x8gf2tvdw0s3jn54khce6mua7lmqqqxw"},
			true,
			0,
			&sdk.TxResponse{},
		},
		{
			"invalid delegation_intent maxWeightOverrun",
			[]string{"chainid", "0.3A12UEL5L,0.3a12uel5l,1.4abcdef1qpzry9x8gf2tvdw0s3jn54khce6mua7lmqqqxw"},
			true,
			0,
			&sdk.TxResponse{},
		},
		{
			"invalid from_address",
			[]string{"chainid", "0.3A12UEL5L,0.3a12uel5l,0.4abcdef1qpzry9x8gf2tvdw0s3jn54khce6mua7lmqqqxw", "--from", "bogusAddress"},
			true,
			0,
			&sdk.TxResponse{},
		},
		{
			"valid",
			[]string{"chainid", "0.3A12UEL5L,0.3a12uel5l,0.4abcdef1qpzry9x8gf2tvdw0s3jn54khce6mua7lmqqqxw", "--from", val.Address.String()},
			false,
			0,
			&sdk.TxResponse{},
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			clientCtx := val.ClientCtx

			flags := []string{}
			args := append(tt.args, flags...)

			cmd := cli.GetSignalIntentTxCmd()

			bz, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
			if tt.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)

				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), tt.respType), bz.String())
				txResp := tt.respType.(*sdk.TxResponse)
				s.Require().Equal(tt.expectedCode, txResp.Code)
			}
		})
	}
}

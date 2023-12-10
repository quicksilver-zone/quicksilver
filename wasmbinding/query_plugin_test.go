package wasmbinding_test

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	"github.com/golang/protobuf/proto" // nolint:staticcheck
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/quicksilver-zone/quicksilver/app"
	"github.com/quicksilver-zone/quicksilver/wasmbinding"
	epochtypes "github.com/quicksilver-zone/quicksilver/x/epochs/types"
)

type StargateTestSuite struct {
	suite.Suite

	ctx sdk.Context
	app *app.Quicksilver
}

func (s *StargateTestSuite) SetupTest() {
	s.app = app.Setup(s.T(), false)
	s.ctx = s.app.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "quicksilver-1", Time: time.Now().UTC()})
}

func TestStargateTestSuite(t *testing.T) {
	suite.Run(t, new(StargateTestSuite))
}

func (s *StargateTestSuite) TestStargateQuerier() {
	testCases := []struct {
		name                   string
		testSetup              func()
		path                   string
		requestData            func() []byte
		responseProtoStruct    interface{}
		expectedQuerierError   bool
		expectedUnMarshalError bool
		resendRequest          bool
	}{
		{
			name: "happy path",
			path: "/quicksilver.epochs.v1.Query/EpochInfos",
			requestData: func() []byte {
				epochrequest := epochtypes.QueryEpochsInfoRequest{}
				bz, err := proto.Marshal(&epochrequest)
				s.Require().NoError(err)
				return bz
			},
			responseProtoStruct: &epochtypes.QueryEpochsInfoResponse{},
		},
		{
			name: "invalid query router route",
			testSetup: func() {
				wasmbinding.SetWhitelistedQuery("invalid/query/router/route", &epochtypes.QueryEpochsInfoRequest{})
			},
			path: "invalid/query/router/route",
			requestData: func() []byte {
				return []byte{}
			},
			expectedQuerierError: true,
		},
		{
			name: "unmatching path and data in request",
			path: "/quicksilver.epochs.v1.Query/EpochInfos",
			requestData: func() []byte {
				epochrequest := epochtypes.QueryCurrentEpochRequest{}
				bz, err := proto.Marshal(&epochrequest)
				s.Require().NoError(err)
				return bz
			},
			responseProtoStruct:    &epochtypes.QueryCurrentEpochResponse{},
			expectedUnMarshalError: true,
		},
		{
			name: "error in unmarshalling response",
			// set up whitelist with wrong data
			testSetup: func() {
				wasmbinding.SetWhitelistedQuery("/quicksilver.epochs.v1.Query/EpochInfos",
					&banktypes.QueryAllBalancesResponse{})
			},
			path: "/quicksilver.epochs.v1.Query/EpochInfos",
			requestData: func() []byte {
				return []byte{}
			},
			responseProtoStruct:  &epochtypes.QueryCurrentEpochResponse{},
			expectedQuerierError: true,
		},
		{
			name: "error in grpc querier",
			// set up whitelist with wrong data
			testSetup: func() {
				wasmbinding.SetWhitelistedQuery("/cosmos.bank.v1beta1.Query/AllBalances", &banktypes.QueryAllBalancesRequest{})
			},
			path: "/cosmos.bank.v1beta1.Query/AllBalances",
			requestData: func() []byte {
				bankrequest := banktypes.QueryAllBalancesRequest{}
				bz, err := proto.Marshal(&bankrequest)
				s.Require().NoError(err)
				return bz
			},
			responseProtoStruct:  &banktypes.QueryAllBalancesRequest{},
			expectedQuerierError: true,
		},
		// TODO: errors in wrong query in state machine
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			s.SetupTest()
			if tc.testSetup != nil {
				tc.testSetup()
			}

			stargateQuerier := wasmbinding.StargateQuerier(*s.app.GRPCQueryRouter(), s.app.AppCodec())
			stargateRequest := &wasmvmtypes.StargateQuery{
				Path: tc.path,
				Data: tc.requestData(),
			}
			stargateResponse, err := stargateQuerier(s.ctx, stargateRequest)
			if tc.expectedQuerierError {
				s.Require().Error(err)
				return
			}

			s.Require().NoError(err)

			protoResponse, ok := tc.responseProtoStruct.(proto.Message)
			s.Require().True(ok)

			// test correctness by unmarshalling json response into proto struct
			err = s.app.AppCodec().UnmarshalJSON(stargateResponse, protoResponse)
			if tc.expectedUnMarshalError {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(protoResponse)
			}

			if tc.resendRequest {
				stargateQuerier = wasmbinding.StargateQuerier(*s.app.GRPCQueryRouter(), s.app.AppCodec())
				stargateRequest = &wasmvmtypes.StargateQuery{
					Path: tc.path,
					Data: tc.requestData(),
				}
				resendResponse, err := stargateQuerier(s.ctx, stargateRequest)
				s.Require().NoError(err)
				s.Require().Equal(stargateResponse, resendResponse)
			}
		})
	}
}

func (s *StargateTestSuite) TestConvertProtoToJsonMarshal() {
	testCases := []struct {
		name                  string
		queryPath             string
		protoResponseStruct   codec.ProtoMarshaler
		originalResponse      string
		expectedProtoResponse codec.ProtoMarshaler
		expectedError         bool
	}{
		{
			name:                "successful conversion from proto response to json marshalled response",
			queryPath:           "/cosmos.bank.v1beta1.Query/AllBalances",
			originalResponse:    "0a090a036261721202333012050a03666f6f",
			protoResponseStruct: &banktypes.QueryAllBalancesResponse{},
			expectedProtoResponse: &banktypes.QueryAllBalancesResponse{
				Balances: sdk.NewCoins(sdk.NewCoin("bar", sdk.NewInt(30))),
				Pagination: &query.PageResponse{
					NextKey: []byte("foo"),
				},
			},
		},
		{
			name:                "invalid proto response struct",
			queryPath:           "/cosmos.bank.v1beta1.Query/AllBalances",
			originalResponse:    "0a090a036261721202333012050a03666f6f",
			protoResponseStruct: &epochtypes.QueryCurrentEpochResponse{},
			expectedError:       true,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			s.SetupTest()

			originalVersionBz, err := hex.DecodeString(tc.originalResponse)
			s.Require().NoError(err)

			jsonMarshalledResponse, err := wasmbinding.ConvertProtoToJSONMarshal(tc.protoResponseStruct, originalVersionBz, s.app.AppCodec())
			if tc.expectedError {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			// check response by json marshalling proto response into json response manually
			jsonMarshalExpectedResponse, err := s.app.AppCodec().MarshalJSON(tc.expectedProtoResponse)
			s.Require().NoError(err)
			s.Require().Equal(jsonMarshalledResponse, jsonMarshalExpectedResponse)
		})
	}
}

// TestDeterministicJsonMarshal tests that we get deterministic JSON marshalled response upon
// proto struct update in the state machine.
func (s *StargateTestSuite) TestDeterministicJsonMarshal() {
	testCases := []struct {
		name                string
		testSetup           func()
		originalResponsebz  []byte
		updatedResponsebz   []byte
		queryPath           string
		responseProtoStruct interface{}
		expectedProto       func() proto.Message
	}{
		/*
		   * Origin Response
		   * balances:<denom:"bar" amount:"30" > pagination:<next_key:"foo" >
		   * New Version Response
		   * The binary built from the proto response with additional field address
		   * balances:<denom:"bar" amount:"30" > pagination:<next_key:"foo" > address:"cosmos1j6j5tsquq2jlw2af7l3xekyaq7zg4l8jsufu78"
		   // Origin proto
		   message QueryAllBalancesResponse {
		  	// balances is the balances of all the coins.
		  	repeated cosmos.base.v1beta1.Coin balances = 1
		  	[(gogoproto.nullable) = false, (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coins"];
		  	// pagination defines the pagination in the response.
		  	cosmos.base.query.v1beta1.PageResponse pagination = 2;
		  }
		  // Updated proto
		  message QueryAllBalancesResponse {
		  	// balances is the balances of all the coins.
		  	repeated cosmos.base.v1beta1.Coin balances = 1
		  	[(gogoproto.nullable) = false, (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coins"];
		  	// pagination defines the pagination in the response.
		  	cosmos.base.query.v1beta1.PageResponse pagination = 2;
		  	// address is the address to query all balances for.
		  	string address = 3;
		  }
		*/
		{
			"Query All Balances",
			func() {
				wasmbinding.SetWhitelistedQuery("/cosmos.bank.v1beta1.Query/AllBalances", &banktypes.QueryAllBalancesResponse{})
			},
			[]byte{10, 9, 10, 3, 98, 97, 114, 18, 2, 51, 48, 18, 5, 10, 3, 102, 111, 111},
			[]byte{
				10, 9, 10, 3, 98, 97, 114, 18, 2, 51, 48, 18, 5, 10, 3, 102, 111, 111, 26, 45, 99, 111, 115, 109, 111, 115, 49, 106,
				54, 106, 53, 116, 115, 113, 117, 113, 50, 106, 108, 119, 50, 97, 102, 55, 108, 51, 120, 101, 107, 121, 97, 113, 55, 122, 103,
				52, 108, 56, 106, 115, 117, 102, 117, 55, 56,
			},
			"/cosmos.bank.v1beta1.Query/AllBalances",
			&banktypes.QueryAllBalancesResponse{},
			func() proto.Message {
				return &banktypes.QueryAllBalancesResponse{
					Balances: sdk.NewCoins(sdk.NewCoin("bar", sdk.NewInt(30))),
					Pagination: &query.PageResponse{
						NextKey: []byte("foo"),
					},
				}
			},
		},
		/*
		  // Origin proto
		  message QueryAccountResponse {
		    // account defines the account of the corresponding address.
		    google.protobuf.Any account = 1 [(cosmos_proto.accepts_interface) = "AccountI"];
		  }
		  // Updated proto
		  message QueryAccountResponse {
		    // account defines the account of the corresponding address.
		    google.protobuf.Any account = 1 [(cosmos_proto.accepts_interface) = "AccountI"];
		    // address is the address to query for.
		  	string address = 2;
		  }
		*/
		{
			"Query Account",
			nil,
			[]byte{
				10, 83, 10, 32, 47, 99, 111, 115, 109, 111, 115, 46, 97, 117, 116, 104, 46, 118, 49, 98, 101, 116, 97, 49, 46, 66, 97, 115,
				101, 65, 99, 99, 111, 117, 110, 116, 18, 47, 10, 45, 99, 111, 115, 109, 111, 115, 49, 102, 56, 117, 120, 117, 108, 116, 110, 56,
				115, 113, 122, 104, 122, 110, 114, 115, 122, 51, 113, 55, 55, 120, 119, 97, 113, 117, 104, 103, 114, 115, 103, 54, 106, 121, 118, 102, 121,
			},
			[]byte{
				10, 83, 10, 32, 47, 99, 111, 115, 109, 111, 115, 46, 97, 117, 116, 104, 46, 118, 49, 98, 101, 116, 97, 49, 46, 66, 97, 115,
				101, 65, 99, 99, 111, 117, 110, 116, 18, 47, 10, 45, 99, 111, 115, 109, 111, 115, 49, 102, 56, 117, 120, 117, 108, 116, 110, 56,
				115, 113, 122, 104, 122, 110, 114, 115, 122, 51, 113, 55, 55, 120, 119, 97, 113, 117, 104, 103, 114, 115, 103, 54, 106, 121, 118, 102, 121,
				18, 45, 99, 111, 115, 109, 111, 115, 49, 102, 56, 117, 120, 117, 108, 116, 110, 56, 115, 113, 122, 104, 122, 110, 114, 115, 122, 51, 113, 55,
				55, 120, 119, 97, 113, 117, 104, 103, 114, 115, 103, 54, 106, 121, 118, 102, 121,
			},
			"/cosmos.auth.v1beta1.Query/Account",
			&authtypes.QueryAccountResponse{},
			func() proto.Message {
				account := authtypes.BaseAccount{
					Address: "cosmos1f8uxultn8sqzhznrsz3q77xwaquhgrsg6jyvfy",
				}
				accountResponse, err := codectypes.NewAnyWithValue(&account)
				s.Require().NoError(err)
				return &authtypes.QueryAccountResponse{
					Account: accountResponse,
				}
			},
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			s.SetupTest()

			if tc.testSetup != nil {
				tc.testSetup()
			}

			binding, err := wasmbinding.GetWhitelistedQuery(tc.queryPath)
			s.Require().Nil(err)

			s.Require().NoError(err)
			jsonMarshalledOriginalBz, err := wasmbinding.ConvertProtoToJSONMarshal(binding, tc.originalResponsebz, s.app.AppCodec())
			s.Require().NoError(err)

			jsonMarshalledUpdatedBz, err := wasmbinding.ConvertProtoToJSONMarshal(binding, tc.updatedResponsebz, s.app.AppCodec())
			s.Require().NoError(err)

			// json marshalled bytes should be the same since we use the same proto struct for unmarshalling
			s.Require().Equal(jsonMarshalledOriginalBz, jsonMarshalledUpdatedBz)

			// raw build also make same result
			jsonMarshalExpectedResponse, err := s.app.AppCodec().MarshalJSON(tc.expectedProto())
			s.Require().NoError(err)
			s.Require().Equal(jsonMarshalledUpdatedBz, jsonMarshalExpectedResponse)
		})
	}
}

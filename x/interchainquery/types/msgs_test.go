package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	ibctesting "github.com/cosmos/ibc-go/v5/testing"

	"github.com/quicksilver-zone/quicksilver/app"
	"github.com/quicksilver-zone/quicksilver/utils/addressutils"
	"github.com/quicksilver-zone/quicksilver/x/interchainquery/keeper"
	"github.com/quicksilver-zone/quicksilver/x/interchainquery/types"
)

var (
	coordinator *ibctesting.Coordinator
	chainA      *ibctesting.TestChain
	chainB      *ibctesting.TestChain
	path        *ibctesting.Path

	testAddress sdk.AccAddress = addressutils.GenerateAccAddressForTest()
)

func init() {
	ibctesting.DefaultTestingAppInit = app.SetupTestingApp
}

func GetSimApp(chain *ibctesting.TestChain) *app.Quicksilver {
	quicksilver, ok := chain.App.(*app.Quicksilver)
	if !ok {
		panic("not quicksilver app")
	}

	return quicksilver
}

func newSimAppPath(chainA, chainB *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointB.ChannelConfig.PortID = ibctesting.TransferPort

	return path
}

func TestMsgSubmitQueryResponse(t *testing.T) {
	coordinator = ibctesting.NewCoordinator(t, 2)
	chainA = coordinator.GetChain(ibctesting.GetChainID(1))
	chainB = coordinator.GetChain(ibctesting.GetChainID(2))
	path = newSimAppPath(chainA, chainB)
	coordinator.SetupConnections(path)

	bondedQuery := stakingtypes.QueryValidatorsRequest{Status: stakingtypes.BondStatusBonded}
	bz, err := bondedQuery.Marshal()
	require.NoError(t, err)

	qvr := stakingtypes.QueryValidatorsResponse{
		Validators: GetSimApp(chainB).StakingKeeper.GetBondedValidatorsByPower(chainB.GetContext()),
	}

	msg := types.MsgSubmitQueryResponse{
		ChainId:     chainB.ChainID + "-N",
		QueryId:     keeper.GenerateQueryHash(path.EndpointB.ConnectionID, chainB.ChainID, "cosmos.staking.v1beta1.Query/Validators", bz, ""),
		Result:      GetSimApp(chainB).AppCodec().MustMarshalJSON(&qvr),
		Height:      chainB.CurrentHeader.Height,
		FromAddress: testAddress.String(),
	}

	require.NoError(t, msg.ValidateBasic())
	require.Equal(t, types.RouterKey, msg.Route())
	require.Equal(t, types.TypeMsgSubmitQueryResponse, msg.Type())
	require.Equal(t, testAddress.String(), msg.GetSigners()[0].String())
}

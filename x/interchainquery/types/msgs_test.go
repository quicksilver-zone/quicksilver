package types_test

import (
	"testing"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctesting "github.com/cosmos/ibc-go/v5/testing"
	"github.com/stretchr/testify/require"

	"github.com/ingenuity-build/quicksilver/app"
	"github.com/ingenuity-build/quicksilver/x/interchainquery/keeper"
	"github.com/ingenuity-build/quicksilver/x/interchainquery/types"
)

const TestOwnerAddress = "osmo1ad2w96t58m0drxcgkf5unm9tg3f9adhjw3m5j5"

var (
	coordinator *ibctesting.Coordinator
	chainA      *ibctesting.TestChain
	chainB      *ibctesting.TestChain
	path        *ibctesting.Path
)

func init() {
	ibctesting.DefaultTestingAppInit = app.SetupTestingApp
}

func GetSimApp(chain *ibctesting.TestChain) *app.Quicksilver {
	app, ok := chain.App.(*app.Quicksilver)
	if !ok {
		panic("not quicksilver app")
	}

	return app
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
		FromAddress: TestOwnerAddress,
	}

	require.NoError(t, msg.ValidateBasic())
	require.Equal(t, types.RouterKey, msg.Route())
	require.Equal(t, types.TypeMsgSubmitQueryResponse, msg.Type())
	require.Equal(t, TestOwnerAddress, msg.GetSigners()[0].String())
}

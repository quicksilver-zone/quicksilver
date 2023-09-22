package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/quicksilver-zone/quicksilver/app"
	"github.com/quicksilver-zone/quicksilver/x/epochs/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app         *app.Quicksilver
	ctx         sdk.Context
	queryClient types.QueryClient
}

// Test helpers.

func (s *KeeperTestSuite) DoSetupTest(t *testing.T) {
	t.Helper()

	checkTx := false
	s.app = app.Setup(t, checkTx)

	s.ctx = s.app.BaseApp.NewContext(false, tmproto.Header{})

	queryHelper := baseapp.NewQueryServerTestHelper(s.ctx, s.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, &s.app.EpochsKeeper)
	s.queryClient = types.NewQueryClient(queryHelper)
}

func (s *KeeperTestSuite) SetupTest() {
	s.DoSetupTest(s.T())
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
